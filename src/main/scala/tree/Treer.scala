package tree

/** Trait for third-party mocking, equivalent to the Go Treer interface. */
trait Treer:
  def get(path: String, node: Option[Node]): Either[TreeError, Node]
  def ast(leaf: Leaf): Either[TreeError, Leaf]
