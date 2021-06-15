GIT_VER := $(shell git describe --tags)

packages:
	$(MAKE) clean
	cd cmd/tenco && gox -os="linux darwin" -arch="amd64" -output="../../dist/{{.Dir}}-${GIT_VER}_{{.OS}}_{{.Arch}}"
	cd dist && find . -name "*${GIT_VER}*" -type f -exec zip {}.zip {} \;

release:
	ghr -u mix3 $(GIT_VER) dist/

clean:
	rm -rf dist/
