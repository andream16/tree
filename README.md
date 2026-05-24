# tree

Recursively builds a tree structure of Rust source files in a project directory. Supports concurrent directory traversal and AST parsing via `syn`.

## Usage

Add to your `Cargo.toml`:

```toml
[dependencies]
tree = { git = "https://github.com/andream16/tree" }
```

### Example

```rust
use tree::get;

fn main() {
    let mut root = get("src").expect("failed to build tree");

    println!("{}", root.name);           // src
    println!("{}", root.leafs[0].name);  // lib.rs

    // Parse the AST
    root.leafs[0].ast().expect("failed to parse");
    if let Some(ref syntax) = root.leafs[0].syntax_tree {
        println!("Items: {}", syntax.items.len());
    }

    // Pretty-print the tree
    root.print();
}
```

## Building

```sh
make build
```

## Testing

```sh
make test
```

## License

MIT
