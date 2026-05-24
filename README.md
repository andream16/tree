# tree [![CircleCI](https://circleci.com/gh/andream16/tree.svg?style=svg)](https://circleci.com/gh/andream16/tree) [![MIT Licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/andream16/tree/master/LICENSE)

Simple Python project-to-tree structure. Supports only `.py` files.

## Install

```bash
pip install -e ".[dev]"
```

## Usage

```python
import ast

from tree import get

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
```

## Test

```bash
make test
```
