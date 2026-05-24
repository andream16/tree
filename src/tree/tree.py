"""Simple Python project-to-tree structure. Supports only .py files."""

from __future__ import annotations

import ast as _ast
import os
from dataclasses import dataclass, field
from typing import Optional

PY_EXT = ".py"


class TreeError(Exception):
    pass


class PathError(TreeError):
    pass


class NodeError(TreeError):
    pass


@dataclass
class Leaf:
    """Represents a Python file."""

    name: str = ""
    path: str = ""
    syntax_tree: Optional[_ast.Module] = field(default=None, repr=False)

    def ast(self) -> None:
        """Parse the file and set syntax_tree."""
        with open(self.path, "r") as f:
            content = f.read()
        if not content:
            return
        self.syntax_tree = _ast.parse(content, filename=self.name)


@dataclass
class Node:
    """Represents a package (directory) that can contain sub-packages."""

    name: str = ""
    nodes: list[Node] = field(default_factory=list)
    leafs: list[Leaf] = field(default_factory=list)

    def print(self, _indent: int = 0) -> None:
        """Pretty-print the project structure."""
        print(self.name)
        print("|", end="")
        for leaf in self.leafs:
            print(f"\t{leaf.name}")
        for node in self.nodes:
            node.print(_indent + 1)


def get(path: str, node: Optional[Node] = None) -> Node:
    """Return all Python files in the given path organized by directory.

    Args:
        path: Relative path to scan.
        node: Optional Node to populate. A new one is created if None.

    Returns:
        Populated Node representing the directory tree.

    Raises:
        NodeError: If node is explicitly passed as a non-Node type.
        PathError: If path is empty.
        FileNotFoundError: If path does not exist.
    """
    if node is None:
        node = Node()

    if not isinstance(node, Node):
        raise NodeError("node must be a Node instance")

    if not path:
        raise PathError("empty path")

    entries = os.listdir(os.path.join(".", path))

    node.name = _current_package(path)

    if not entries:
        return node

    dirs: list[Node] = []
    leafs: list[Leaf] = []

    for entry in sorted(entries):
        full = os.path.join(path, entry)
        if os.path.isfile(os.path.join(".", full)):
            _, ext = os.path.splitext(entry)
            if ext == PY_EXT:
                leafs.append(Leaf(name=entry))
        elif os.path.isdir(os.path.join(".", full)):
            dirs.append(Node(name=entry))

    for sub_node in dirs:
        child = get(path + "/" + sub_node.name, sub_node)
        node.nodes.append(child)

    for leaf in leafs:
        leaf.path = path + "/" + leaf.name
        node.leafs.append(leaf)

    return node


def _current_package(path: str) -> str:
    i = path.find("/")
    if i == -1:
        return path
    return path[i + 1:]
