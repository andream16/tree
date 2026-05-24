package tree

import org.scalatest.flatspec.AnyFlatSpec
import org.scalatest.matchers.should.Matchers
import org.scalatest.EitherValues

class TreeSpec extends AnyFlatSpec with Matchers with EitherValues:

  "Tree.get" should "fail because node is None" in {
    val result = Tree.get("somepath", None)
    result.isLeft shouldBe true
    result.left.value.message should include("node")
  }

  it should "fail because path is empty" in {
    val result = Tree.get("", Some(Node()))
    result.isLeft shouldBe true
    result.left.value.message should include("path")
  }

  it should "fail because path is not found" in {
    val result = Tree.get("nonexistent", Some(Node()))
    result.isLeft shouldBe true
    result.left.value.message should include("path not found")
  }

  it should "return expected result for testdata" in {
    val result = Tree.get("src/test/resources/testdata", Some(Node()))
    result.isRight shouldBe true

    val node = result.value
    node.name shouldBe "test/resources/testdata"

    // Should have one leaf (somefile.scala) at root level
    node.leafs.map(_.name) should contain("somefile.scala")

    // Should have one subdirectory node
    node.nodes should have size 1
    // currentPackage strips the first path component
    node.nodes.head.name shouldBe "test/resources/testdata/somedir"

    // The subdirectory should contain two scala files
    val somedir = node.nodes.head
    somedir.leafs.map(_.name) should contain allOf ("somefile.scala", "someotherfile.scala")

    // Verify paths are set correctly
    val rootLeaf = node.leafs.find(_.name == "somefile.scala").get
    rootLeaf.path shouldBe "src/test/resources/testdata/somefile.scala"

    val subdirLeafs = somedir.leafs
    subdirLeafs.foreach { l =>
      l.path should startWith("src/test/resources/testdata/somedir/")
    }
  }

  "Leaf.ast" should "fail because file does not exist" in {
    val leaf = Leaf(name = "missing.scala", path = "nonexistent.scala")
    val result = leaf.ast()
    result.isLeft shouldBe true
    result.left.value.message should include("file not found")
  }

  it should "return early because file is empty" in {
    val leaf = Leaf(
      name = "somefile.scala",
      path = "src/test/resources/testdata/somedir/somefile.scala"
    )
    val result = leaf.ast()
    result.isRight shouldBe true
    result.value.syntaxTree shouldBe None
  }

  it should "return a valid syntax tree" in {
    val leaf = Leaf(
      name = "somefile.scala",
      path = "src/test/resources/testdata/somefile.scala"
    )
    val result = leaf.ast()
    result.isRight shouldBe true
    result.value.syntaxTree shouldBe defined
  }
