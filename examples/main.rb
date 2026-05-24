#!/usr/bin/env ruby
# frozen_string_literal: true

require_relative "../lib/tree"

node = Tree::Node.new
result = Tree.get("examples/example", node)

puts result.name                   # example
puts result.leafs[0].name          # somefile.rb
puts result.leafs[0].path          # examples/example/somefile.rb
puts result.nodes[0].name          # subexample
puts result.nodes[0].leafs[0].name # someotherfile.rb
puts result.nodes[0].leafs[0].path # examples/example/subexample/someotherfile.rb

# Parse AST
result.leafs[0].ast
puts result.leafs[0].syntax_tree

# Pretty print
result.print
