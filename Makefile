clean:
	rm -rf dist/
build:
	go build -o dist/ghas-org-enablement main.go
snapshot:
	goreleaser release --snapshot
release:
	goreleaser release