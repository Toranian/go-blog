all: build
sass:
	@(sass --watch web/static/scss/style.scss:web/static/css/style.css --style compressed)

.PHONY: all sass
