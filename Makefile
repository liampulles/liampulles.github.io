GOBIN := $(shell go env GOBIN)

# Use secrets
include secrets

clean:
	rm -rf _site_gen
	rm -rf _site

pre-commit: ${GOBIN}/minify clean
	$(MAKE) -C htmlgen install
	htmlgen -output=_site_gen
	cp -r static_minable/* _site_gen
	minify -r -o _site/ _site_gen/
	cp -r static/* _site

watch:
	while true; do \
		$(MAKE) pre-commit; \
		inotifywait -qre close_write htmlgen; \
	done

${GOBIN}/minify:
	go install github.com/tdewolff/minify/cmd/minify@latest