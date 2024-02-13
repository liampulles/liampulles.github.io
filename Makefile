GOBIN := $(shell go env GOBIN)

run: gen-proverbs bundle-install
	-fuser -k 4001/tcp
	bundle exec jekyll serve --port 4001 --draft

pre-commit: gen-proverbs bundle-install ${GOBIN}/minify
	bundle update
	bundle exec jekyll build
	minify -r -o _min_site/ _site/
	cp -r _min_site/* _site
	rm -rf _min_site

htmlgen-test:
	$(MAKE) -C htmlgen install
	htmlgen-cli -output=_site_test
	rm -rf _site_test_min
	minify -r -o _site_test_min/ _site_test/

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