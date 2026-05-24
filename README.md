# tree [![MIT Licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/andream16/tree/master/LICENSE)

Simple Ruby project to tree structure. It supports only Ruby files.

## Installation

Add to your Gemfile:

```ruby
gem "tree", git: "https://github.com/andream16/tree.git"
```

Then run `bundle install`.

## Usage

```ruby
require "tree"

node = Tree::Node.new
result = Tree.get("examples/example", node)

puts result.name                   # example
puts result.leafs[0].name          # somefile.rb
puts result.leafs[0].path          # examples/example/somefile.rb
puts result.nodes[0].name          # subexample
puts result.nodes[0].leafs[0].name # someotherfile.rb
puts result.nodes[0].leafs[0].path # examples/example/subexample/someotherfile.rb

# Parse AST for a leaf
result.leafs[0].ast
puts result.leafs[0].syntax_tree   # (module (const nil :Example) nil)

# Pretty print the tree
result.print
# example
# |	somefile.rb
# subexample
# |	someotherfile.rb
```

## Testing

```
make install
make test
```
