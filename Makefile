default: test

install:
	bundle install

test:
	bundle exec rspec --format documentation