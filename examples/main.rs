use tree::{get, Node};

fn main() {
    let mut root: Node = get("examples/example").expect("failed to get tree");

    println!("Name: {}", root.name);
    println!("Leaf: {}", root.leafs[0].name);
    println!("Path: {}", root.leafs[0].path);
    println!("Sub-node: {}", root.nodes[0].name);
    println!("Sub-leaf: {}", root.nodes[0].leafs[0].name);
    println!("Sub-path: {}", root.nodes[0].leafs[0].path);

    // Parse the AST of the first leaf
    root.leafs[0].ast().expect("failed to parse AST");
    if let Some(ref syntax) = root.leafs[0].syntax_tree {
        println!("Items in {}: {}", root.leafs[0].name, syntax.items.len());
    }

    println!();
    root.print();
}
