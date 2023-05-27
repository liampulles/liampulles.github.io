run: gen-proverbs bundle-install
	bundle exec jekyll serve --port 4002

pre-commit: gen-proverbs bundle-install
	bundle update
	bundle exec jekyll build

gen-proverbs:
	proverb-gen > proverbs.html

# You might need to run this a couple times for it to work
bundle-install:
	bundle install

setup-local:
	sudo apt-get install ruby-full
	sudo gem install bundler
