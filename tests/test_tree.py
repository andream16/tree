"""Tests for the tree module."""

import ast
import os

import pytest

from tree import Node, Leaf, get
from tree.tree import PathError, NodeError


class TestGet:
    def test_raises_on_non_node(self):
        with pytest.raises(NodeError):
            get("somepath", "not a node")  # type: ignore

    def test_raises_on_empty_path(self):
        with pytest.raises(PathError):
            get("", Node())

    def test_raises_on_missing_path(self):
        with pytest.raises(FileNotFoundError):
            get("nonexistent_path", Node())

    def test_returns_expected_result(self):
        out = get("testdata_py", Node())

        assert out.name == "testdata_py"
        assert len(out.nodes) == 1
        assert out.nodes[0].name == "somedir"
        assert len(out.nodes[0].leafs) == 2

        leaf_names = sorted(l.name for l in out.nodes[0].leafs)
        assert leaf_names == ["somefile.py", "someotherfile.py"]

        assert len(out.leafs) == 1
        assert out.leafs[0].name == "somefile.py"
        assert out.leafs[0].path == "testdata_py/somefile.py"

    def test_default_node_creation(self):
        out = get("testdata_py")
        assert out.name == "testdata_py"
        assert isinstance(out, Node)


class TestLeafAst:
    def test_raises_on_missing_file(self):
        leaf = Leaf(name="missing.py", path="nonexistent/missing.py")
        with pytest.raises(FileNotFoundError):
            leaf.ast()

    def test_returns_early_on_empty_file(self):
        leaf = Leaf(
            name="somefile.py",
            path="testdata_py/somedir/somefile.py",
        )
        leaf.ast()
        assert leaf.syntax_tree is None

    def test_parses_valid_file(self):
        leaf = Leaf(
            name="somefile.py",
            path="testdata_py/somefile.py",
        )
        leaf.ast()
        assert leaf.syntax_tree is not None
        assert isinstance(leaf.syntax_tree, ast.Module)


class TestNodePrint:
    def test_print(self, capsys):
        node = Node(
            name="root",
            leafs=[Leaf(name="a.py"), Leaf(name="b.py")],
            nodes=[
                Node(name="sub", leafs=[Leaf(name="c.py")])
            ],
        )
        node.print()
        captured = capsys.readouterr()
        assert "root" in captured.out
        assert "a.py" in captured.out
        assert "b.py" in captured.out
        assert "sub" in captured.out
        assert "c.py" in captured.out
