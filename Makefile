GOBIN := $(shell go env GOBIN)
RESTORE_IMGS=$(wildcard static/images/restorations/*.png)
RESTORE_THUMBS=$(patsubst static/images/restorations/%.png,static/images/restorations-thumb/%.jpg,$(RESTORE_IMGS))

clean:
	rm -rf _site_gen
	rm -rf _site

pre-commit: ${GOBIN}/minify clean compress-restorations
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

compress-restorations: $(RESTORE_THUMBS)

static/images/restorations-thumb/%.jpg: static/images/restorations/%.png
	convert -geometry x1000 -strip -interlace Plane -quality 80 $^ $@

${GOBIN}/minify:
	go install github.com/tdewolff/minify/cmd/minify@latest