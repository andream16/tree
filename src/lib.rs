use std::fs;
use std::io;
use std::path::Path;
use std::sync::mpsc;
use std::thread;

pub use syn::File as SynFile;

const RS_EXT: &str = "rs";

/// Represents a directory that can contain sub-directories (Nodes) and Rust source files (Leafs).
#[derive(Debug, Clone, PartialEq)]
pub struct Node {
    pub name: String,
    pub nodes: Vec<Node>,
    pub leafs: Vec<Leaf>,
}

/// Represents a single `.rs` file within a Node.
pub struct Leaf {
    pub name: String,
    pub path: String,
    pub syntax_tree: Option<SynFile>,
}

impl std::fmt::Debug for Leaf {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        f.debug_struct("Leaf")
            .field("name", &self.name)
            .field("path", &self.path)
            .field("syntax_tree", &self.syntax_tree.as_ref().map(|_| "..."))
            .finish()
    }
}

impl Clone for Leaf {
    fn clone(&self) -> Self {
        Leaf {
            name: self.name.clone(),
            path: self.path.clone(),
            // syn::File is not Clone; re-parse if needed
            syntax_tree: None,
        }
    }
}

impl PartialEq for Leaf {
    fn eq(&self, other: &Self) -> bool {
        self.name == other.name && self.path == other.path
    }
}

#[derive(Debug)]
pub enum Error {
    EmptyPath,
    Io(io::Error),
    Parse(String),
}

// SAFETY: io::Error is Send, String is Send.
unsafe impl Send for Error {}

impl std::fmt::Display for Error {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Error::EmptyPath => write!(f, "empty path"),
            Error::Io(e) => write!(f, "io error: {e}"),
            Error::Parse(e) => write!(f, "parse error: {e}"),
        }
    }
}

impl std::error::Error for Error {}

impl From<io::Error> for Error {
    fn from(e: io::Error) -> Self {
        Error::Io(e)
    }
}

impl From<syn::Error> for Error {
    fn from(e: syn::Error) -> Self {
        Error::Parse(e.to_string())
    }
}

/// Recursively walks `path` and builds a tree of Rust source files.
///
/// Traverses directories concurrently, collecting `.rs` files as Leafs
/// and subdirectories as child Nodes.
pub fn get(path: &str) -> Result<Node, Error> {
    if path.is_empty() {
        return Err(Error::EmptyPath);
    }

    let dir = Path::new("./").join(path);
    let entries = fs::read_dir(&dir)?;

    let mut nodes_to_visit: Vec<(String, String)> = Vec::new();
    let mut leafs: Vec<Leaf> = Vec::new();

    for entry in entries {
        let entry = entry?;
        let file_type = entry.file_type()?;
        let file_name = entry.file_name().to_string_lossy().to_string();

        if file_type.is_file() {
            if Path::new(&file_name)
                .extension()
                .is_some_and(|ext| ext == RS_EXT)
            {
                leafs.push(Leaf {
                    name: file_name.clone(),
                    path: format!("{}/{}", path, file_name),
                    syntax_tree: None,
                });
            }
        } else if file_type.is_dir() {
            let child_path = format!("{}/{}", path, file_name);
            nodes_to_visit.push((file_name, child_path));
        }
    }

    // Sort for deterministic output
    leafs.sort_by(|a, b| a.name.cmp(&b.name));
    nodes_to_visit.sort_by(|a, b| a.0.cmp(&b.0));

    // Recurse into subdirectories concurrently using scoped threads.
    // Node/Leaf contain Option<SynFile> which isn't Send, but during tree
    // construction syntax_tree is always None, so we use unsafe Send wrapper.
    struct SendNode(Node);
    // SAFETY: During get(), syntax_tree is always None. syn::File (which is !Send
    // due to proc_macro2::Span) is never present during concurrent traversal.
    unsafe impl Send for SendNode {}

    let (tx, rx) = mpsc::channel::<(String, Result<SendNode, Error>)>();

    let handles: Vec<_> = nodes_to_visit
        .into_iter()
        .map(|(name, child_path)| {
            let tx = tx.clone();
            thread::spawn(move || {
                let result = get(&child_path).map(SendNode);
                tx.send((name, result)).ok();
            })
        })
        .collect();

    drop(tx);

    let mut child_nodes: Vec<(String, Node)> = Vec::new();
    for (name, result) in rx {
        let node = result?.0;
        child_nodes.push((name, node));
    }

    for handle in handles {
        handle.join().ok();
    }

    // Sort child nodes for deterministic output
    child_nodes.sort_by(|a, b| a.0.cmp(&b.0));

    let node_name = current_package(path);

    Ok(Node {
        name: node_name,
        nodes: child_nodes.into_iter().map(|(_, n)| n).collect(),
        leafs,
    })
}

impl Leaf {
    /// Parses the Rust source file and populates `syntax_tree`.
    pub fn ast(&mut self) -> Result<(), Error> {
        let content = fs::read_to_string(&self.path)?;
        if content.is_empty() {
            return Ok(());
        }
        let syntax = syn::parse_file(&content)?;
        self.syntax_tree = Some(syntax);
        Ok(())
    }
}

impl Node {
    /// Pretty-prints the tree structure.
    pub fn print(&self) {
        self.print_inner("");
    }

    fn print_inner(&self, prefix: &str) {
        println!("{}{}", prefix, self.name);
        for leaf in &self.leafs {
            println!("{}|\t{}", prefix, leaf.name);
        }
        for node in &self.nodes {
            node.print_inner(prefix);
        }
    }
}

fn current_package(path: &str) -> String {
    match path.find('/') {
        Some(i) => path[i + 1..].to_string(),
        None => path.to_string(),
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_get_empty_path() {
        let result = get("");
        assert!(matches!(result, Err(Error::EmptyPath)));
    }

    #[test]
    fn test_get_nonexistent_path() {
        let result = get("nonexistent_path");
        assert!(matches!(result, Err(Error::Io(_))));
    }

    #[test]
    fn test_get_testdata() {
        let result = get("testdata").unwrap();

        assert_eq!(result.name, "testdata");
        assert_eq!(result.leafs.len(), 1);
        assert_eq!(result.leafs[0].name, "somefile.rs");
        assert_eq!(result.leafs[0].path, "testdata/somefile.rs");

        assert_eq!(result.nodes.len(), 1);
        assert_eq!(result.nodes[0].name, "somedir");
        assert_eq!(result.nodes[0].leafs.len(), 2);
        assert_eq!(result.nodes[0].leafs[0].name, "somefile.rs");
        assert_eq!(result.nodes[0].leafs[0].path, "testdata/somedir/somefile.rs");
        assert_eq!(result.nodes[0].leafs[1].name, "someotherfile.rs");
        assert_eq!(
            result.nodes[0].leafs[1].path,
            "testdata/somedir/someotherfile.rs"
        );
    }

    #[test]
    fn test_leaf_ast_file_not_found() {
        let mut leaf = Leaf {
            name: String::new(),
            path: String::new(),
            syntax_tree: None,
        };
        let result = leaf.ast();
        assert!(matches!(result, Err(Error::Io(_))));
    }

    #[test]
    fn test_leaf_ast_empty_file() {
        let mut leaf = Leaf {
            name: "somefile.rs".to_string(),
            path: "testdata/somedir/somefile.rs".to_string(),
            syntax_tree: None,
        };
        let result = leaf.ast();
        assert!(result.is_ok());
        assert!(leaf.syntax_tree.is_none());
    }

    #[test]
    fn test_leaf_ast_valid_file() {
        let mut leaf = Leaf {
            name: "somefile.rs".to_string(),
            path: "testdata/somefile.rs".to_string(),
            syntax_tree: None,
        };
        let result = leaf.ast();
        assert!(result.is_ok());
        assert!(leaf.syntax_tree.is_some());
    }
}
