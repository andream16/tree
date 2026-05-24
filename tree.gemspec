# frozen_string_literal: true

Gem::Specification.new do |spec|
  spec.name          = "tree"
  spec.version       = "1.0.0"
  spec.authors       = ["andream16"]
  spec.summary       = "Simple Ruby project to tree structure. Supports only .rb files."
  spec.license       = "MIT"
  spec.files         = Dir["lib/**/*.rb"]
  spec.require_paths = ["lib"]

  spec.required_ruby_version = ">= 3.0"

  spec.add_dependency "parser", "~> 3.0"

  spec.add_development_dependency "rspec", "~> 3.12"
  spec.add_development_dependency "rake", "~> 13.0"
end
