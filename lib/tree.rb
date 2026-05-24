# frozen_string_literal: true

require "parser/current"

module Tree
  RUBY_EXT = ".rb"

  class Error < StandardError; end
  class PathError < Error; end
  class NodeError < Error; end

  # Represents a package (directory) that can contain sub-packages.
  class Node
    attr_accessor :name, :nodes, :leafs

    def initialize(name: nil, nodes: [], leafs: [])
      @name = name
      @nodes = nodes
      @leafs = leafs
    end

    # Pretty prints the project structure.
    def print(io: $stdout)
      io.puts "#{@name}"
      io.print "|"
      @leafs.each { |l| io.puts "\t#{l.name}" }
      @nodes.each { |n| n.print(io: io) }
    end
  end

  # Represents a Ruby file.
  class Leaf
    attr_accessor :name, :path, :syntax_tree

    def initialize(name: nil, path: nil)
      @name = name
      @path = path
      @syntax_tree = nil
    end

    # Parses the file and sets syntax_tree to the AST.
    def ast
      content = File.read(@path)
      return if content.empty?

      @syntax_tree = Parser::CurrentRuby.parse(content)
    end
  end

  # Returns all Ruby files in the given path as a Node tree.
  def self.get(path, node)
    validate!(path, node)

    entries = Dir.entries("./#{path}").reject { |e| e.start_with?(".") }.sort
    node.name = current_package(path)

    return node if entries.empty?

    nodes, leafs = filter_ruby_files_dirs(entries, path)

    threads = nodes.map do |child_node|
      Thread.new { Tree.get("#{path}/#{child_node.name}", child_node) }
    end

    threads.each do |t|
      result = t.value
      node.nodes << result
    end

    leafs.each do |leaf|
      leaf.path = "#{path}/#{leaf.name}"
      node.leafs << leaf
    end

    node
  end

  class << self
    private

    def validate!(path, node)
      raise NodeError, "node can't be nil" if node.nil?
      raise PathError, "empty path" if path.nil? || path.empty?
    end

    def current_package(path)
      idx = path.index("/")
      return path if idx.nil?

      path[(idx + 1)..]
    end

    def filter_ruby_files_dirs(entries, path)
      leafs = []
      nodes = []

      entries.each do |entry|
        full_path = "./#{path}/#{entry}"
        if File.extname(entry) == RUBY_EXT
          leafs << Leaf.new(name: entry)
        elsif File.directory?(full_path)
          nodes << Node.new(name: entry)
        end
      end

      [nodes, leafs]
    end
  end
end
