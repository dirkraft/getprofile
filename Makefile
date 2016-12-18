.PHONY: all publish-dev clean develop

all: dist/getprofile.darwin.amd64 dist/getprofile.linux.386 dist/getprofile.linux.amd64 dist/getprofile.windows.amd64.exe

dist/getprofile.darwin.amd64:
	mkdir -p dist/
	GOOS=darwin GOARCH=amd64 go build -o dist/getprofile.darwin.amd64 getprofile.go

dist/getprofile.linux.386:
	mkdir -p dist/
	GOOS=linux GOARCH=386 go build -o dist/getprofile.linux.386 getprofile.go

dist/getprofile.linux.amd64:
	mkdir -p dist/
	GOOS=linux GOARCH=amd64 go build -o dist/getprofile.linux.amd64 getprofile.go

dist/getprofile.windows.amd64.exe:
	mkdir -p dist/
	GOOS=windows GOARCH=amd64 go build -o dist/getprofile.windows.amd64.exe getprofile.go

publish-dev: all
	scripts/unpublish-dev.sh
	scripts/publish-dev.sh dist/getprofile.darwin.amd64
	scripts/publish-dev.sh dist/getprofile.linux.386
	scripts/publish-dev.sh dist/getprofile.linux.amd64
	scripts/publish-dev.sh dist/getprofile.windows.amd64.exe
	scripts/update-dev-tag.sh

clean:
	rm -rf dist/

develop:
	go get ./...
