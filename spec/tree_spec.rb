# frozen_string_literal: true

require_relative "../lib/tree"

RSpec.describe Tree do
  describe ".get" do
    it "raises NodeError when node is nil" do
      expect { Tree.get("somepath", nil) }.to raise_error(Tree::NodeError)
    end

    it "raises PathError when path is empty" do
      node = Tree::Node.new
      expect { Tree.get("", node) }.to raise_error(Tree::PathError)
    end

    it "raises Errno::ENOENT when path does not exist" do
      node = Tree::Node.new
      expect { Tree.get("nonexistent_path", node) }.to raise_error(Errno::ENOENT)
    end

    it "returns expected tree structure" do
      node = Tree::Node.new
      result = Tree.get("testdata", node)

      expect(result.name).to eq("testdata")

      expect(result.leafs.length).to eq(1)
      expect(result.leafs[0].name).to eq("somefile.rb")
      expect(result.leafs[0].path).to eq("testdata/somefile.rb")

      expect(result.nodes.length).to eq(1)
      expect(result.nodes[0].name).to eq("somedir")
      expect(result.nodes[0].leafs.length).to eq(2)
      expect(result.nodes[0].leafs.map(&:name)).to contain_exactly("somefile.rb", "someotherfile.rb")
      expect(result.nodes[0].leafs.map(&:path)).to contain_exactly(
        "testdata/somedir/somefile.rb",
        "testdata/somedir/someotherfile.rb"
      )
    end
  end

  describe Tree::Leaf do
    describe "#ast" do
      it "raises Errno::ENOENT when file does not exist" do
        leaf = Tree::Leaf.new(name: "missing.rb", path: "missing.rb")
        expect { leaf.ast }.to raise_error(Errno::ENOENT)
      end

      it "returns nil for an empty file" do
        leaf = Tree::Leaf.new(name: "somefile.rb", path: "testdata/somedir/somefile.rb")
        expect(leaf.ast).to be_nil
        expect(leaf.syntax_tree).to be_nil
      end

      it "parses a valid Ruby file into an AST" do
        leaf = Tree::Leaf.new(name: "somefile.rb", path: "testdata/somefile.rb")
        leaf.ast

        expect(leaf.syntax_tree).not_to be_nil
        expect(leaf.syntax_tree).to be_a(Parser::AST::Node)
      end
    end
  end

  describe Tree::Node do
    describe "#print" do
      it "outputs the tree structure" do
        node = Tree::Node.new
        Tree.get("testdata", node)

        output = StringIO.new
        node.print(io: output)

        text = output.string
        expect(text).to include("testdata")
        expect(text).to include("somefile.rb")
        expect(text).to include("somedir")
      end
    end
  end
end
