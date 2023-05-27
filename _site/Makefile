GOBIN := $(shell go env GOBIN)

run: gen-proverbs bundle-install
	bundle exec jekyll serve --port 4002 --draft

pre-commit: gen-proverbs bundle-install ${GOBIN}/minify
	bundle update
	bundle exec jekyll build
	minify -r -o _min_site/ _site/
	cp -r _min_site/* _site
	rm -rf _min_site

gen-proverbs:
	proverb-gen > proverbs.html

# You might need to run this a couple times for it to work
bundle-install:
	bundle install

setup-local:
	sudo apt-get install ruby-full
	sudo gem install bundler

${GOBIN}/minify:
	go install github.com/tdewolff/minify/cmd/minify@latest