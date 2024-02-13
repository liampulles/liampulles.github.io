GOBIN := $(shell go env GOBIN)

clean:
	rm -rf _site_gen
	rm -rf _site

pre-commit: ${GOBIN}/minify clean
	$(MAKE) -C htmlgen install
	htmlgen -output=_site_gen
	cp -r static/* _site_gen
	minify -r -o _site/ _site_gen/
	cp -r images _site/images

${GOBIN}/minify:
	go install github.com/tdewolff/minify/cmd/minify@latest