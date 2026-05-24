"""Example usage of the tree package."""

import ast

from tree import get


def main():
    out = get("examples/example")

    print(out.name)                        # example
    print(out.leafs[0].name)               # somefile.py
    print(out.leafs[0].path)               # examples/example/somefile.py
    print(out.nodes[0].name)               # example/subexample
    print(out.nodes[0].leafs[0].name)      # someotherfile.py
    print(out.nodes[0].leafs[0].path)      # examples/example/subexample/someotherfile.py

    out.leafs[0].ast()

    for node in ast.walk(out.leafs[0].syntax_tree):
        if isinstance(node, ast.Assign):
            for target in node.targets:
                if isinstance(target, ast.Name):
                    print(target.id)       # x

    # example
    # |       somefile.py
    # example/subexample
    # |       someotherfile.py
    out.print()


if __name__ == "__main__":
    main()
