default: test

.PHONY: install test lint

install:
	pip install -e ".[dev]"

test:
	pytest -v --tb=short

lint:
	python -m py_compile src/tree/tree.py