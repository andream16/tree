ThisBuild / scalaVersion := "3.3.1"
ThisBuild / organization := "com.github.andream16"
ThisBuild / version      := "0.1.0"

lazy val root = (project in file("."))
  .settings(
    name := "tree",
    libraryDependencies ++= Seq(
      "org.scalameta" %% "scalameta" % "4.17.0",
      "org.scalatest" %% "scalatest" % "3.2.17" % Test
    )
  )
