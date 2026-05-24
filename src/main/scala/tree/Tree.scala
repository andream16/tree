package tree

import java.io.File
import java.nio.file.{Files, Paths}
import scala.concurrent.{Await, Future}
import scala.concurrent.ExecutionContext.Implicits.global
import scala.concurrent.duration.Duration
import scala.meta.*

val ScalaExt = ".scala"

/** Represents a package (directory) that can contain sub-packages. */
case class Node(
    var name: String = "",
    var nodes: List[Node] = Nil,
    var leafs: List[Leaf] = Nil
)

/** Represents a Scala source file. */
case class Leaf(
    name: String = "",
    path: String = "",
    var syntaxTree: Option[Source] = None
)

/** Errors used for validation. */
case class TreeError(message: String) extends Exception(message)

object Tree:

  private val errPath = "path"
  private val errNode = "node"

  /** Returns all Scala files in the given path as a tree of Nodes and Leafs.
    *
    * Equivalent to the Go `Get` function - recursively walks directories,
    * collecting `.scala` files as Leafs and subdirectories as child Nodes.
    */
  def get(path: String, node: Option[Node]): Either[TreeError, Node] =
    validate(path, node) match
      case Some(err) => Left(err)
      case None =>
        val n = node.get
        val dir = new File("./" + path)

        if !dir.exists() || !dir.isDirectory then
          Left(TreeError(s"path not found: $path"))
        else
          val files = dir.listFiles()
          if files == null then Left(TreeError(s"cannot read directory: $path"))
          else
            n.name = currentPackage(path)

            if files.isEmpty then Right(n)
            else
              val (dirs, scalaFiles) = filterScalaFilesDirs(files.toList)

              // Process subdirectories concurrently
              val futures = dirs.map { childNode =>
                Future {
                  get(path + "/" + childNode.name, Some(childNode))
                }
              }

              val results = futures.map(f => Await.result(f, Duration.Inf))

              val error = results.collectFirst { case Left(err) => err }
              error match
                case Some(err) => Left(err)
                case None =>
                  results.foreach {
                    case Right(child) => n.nodes = n.nodes :+ child
                    case _            => ()
                  }

                  // Attach leafs with their full paths
                  for leaf <- scalaFiles do
                    n.leafs = n.leafs :+ leaf.copy(path = path + "/" + leaf.name)

                  Right(n)

  private def validate(path: String, node: Option[Node]): Option[TreeError] =
    if node.isEmpty then Some(TreeError(s"node can't be nil: $errNode"))
    else if path.isEmpty then Some(TreeError(s"empty path: $errPath"))
    else None

  private def currentPackage(path: String): String =
    val i = path.indexOf('/')
    if i == -1 then path
    else path.substring(i + 1)

  private def filterScalaFilesDirs(files: List[File]): (List[Node], List[Leaf]) =
    val leafs = files.collect {
      case f if f.isFile && f.getName.endsWith(ScalaExt) =>
        Leaf(name = f.getName)
    }
    val nodes = files.collect {
      case f if f.isDirectory =>
        Node(name = f.getName)
    }
    (nodes, leafs)

end Tree

extension (leaf: Leaf)
  /** Parses the Scala source file and sets the syntaxTree field.
    * Equivalent to the Go `Ast()` method.
    */
  def ast(): Either[TreeError, Leaf] =
    val path = Paths.get(leaf.path)
    if !Files.exists(path) then
      return Left(TreeError(s"file not found: ${leaf.path}"))

    val content = new String(Files.readAllBytes(path))
    if content.isEmpty then return Right(leaf)

    content.parse[Source] match
      case Parsed.Success(tree) =>
        val updated = leaf.copy(syntaxTree = Some(tree))
        Right(updated)
      case Parsed.Error(pos, msg, _) =>
        Left(TreeError(s"parse error in ${leaf.name}: $msg"))

extension (node: Node)
  /** Pretty prints the project structure.
    * Equivalent to the Go `Print()` method.
    */
  def printTree(): Unit =
    println(node.name)
    scala.Predef.print("|")
    for l <- node.leafs do
      println("\t" + l.name)
    for child <- node.nodes do
      child.printTree()
