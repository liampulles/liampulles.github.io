run: gen-proverbs
	bundle exec jekyll serve

pre-commit: gen-proverbs
	bundle install
	bundle update
	bundle exec jekyll build

gen-proverbs:
	proverb-gen > proverbs.html

setup-local:
	sudo apt-get install ruby-full
	sudo gem install bundler