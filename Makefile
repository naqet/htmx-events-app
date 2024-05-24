run:
	npx tailwindcss -i ./root.css -o ./static/styles.css --minify
	templ generate
	go build -o ./tmp/main .
