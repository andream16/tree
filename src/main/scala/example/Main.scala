package example

import tree.{Tree as ScalaTree, *}
import scala.meta.*

@main def run(): Unit =
  val result = ScalaTree.get("examples/example", Some(Node()))
  result match
    case Left(err) =>
      throw new RuntimeException(err.message)
    case Right(out) =>
      println(out.name)                        // example
      println(out.leafs.head.name)             // somefile.scala
      println(out.leafs.head.path)             // examples/example/somefile.scala
      println(out.nodes.head.name)             // example/subexample
      println(out.nodes.head.leafs.head.name)  // someotherfile.scala
      println(out.nodes.head.leafs.head.path)  // examples/example/subexample/someotherfile.scala

      // Parse the AST of the first leaf
      out.leafs.head.ast() match
        case Left(err) =>
          throw new RuntimeException(err.message)
        case Right(leaf) =>
          leaf.syntaxTree.foreach { tree =>
            tree.traverse {
              case name: Term.Name => println(name.value)
              case lit: Lit        => println(lit.syntax)
            }
          }

      // Pretty print the tree structure
      out.printTree()
