GOBIN := $(shell go env GOBIN)

pre-commit: ${GOBIN}/minify
	$(MAKE) -C htmlgen install
	htmlgen -output=_site_gen
	rm -rf _site
	minify -r -o _site/ _site_gen/
	cp -r images _site/images
	cp -r static/* _site

${GOBIN}/minify:
	go install github.com/tdewolff/minify/cmd/minify@latest