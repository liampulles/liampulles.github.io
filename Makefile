clean:
	rm -rf _site_gen
	rm -rf _site

combine-js:
	cat _site_gen/maybe_pages.js _site_gen/script.js > _site_gen/temp.js
	rm _site_gen/maybe_pages.js _site_gen/script.js
	mv _site_gen/temp.js _site_gen/script.js

pre-commit: ${GOBIN}/minify clean
	$(MAKE) -C htmlgen install
	htmlgen -output=_site_gen
	cp -r static_minable/* _site_gen
	$(MAKE) combine-js
	minify -r -o _site/ _site_gen/
	cp -r static/* _site

watch:
	while true; do \
		$(MAKE) pre-commit; \
		inotifywait -qre close_write htmlgen; \
	done

# Place the restoration in the main folder, and use this to create a compressed thumb version.
# Might have to bump up your imagemagick policy: https://stackoverflow.com/a/53699200
static/images/restorations-thumb/%.jpg: %.png
	convert -geometry x1000 -strip -interlace Plane -quality 90 $^ $@

${GOBIN}/minify:
	go install github.com/tdewolff/minify/cmd/minify@latest