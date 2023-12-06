# assemble-and-tag-action
Properly assembles and tags your Golang based GitHub Action

---

A GitHub Action for publishing Golang based Actions! It's designed to act on new releases, and updates the tag with any compiled Go binaries. The process looks like this:

- Scans the `./bin` directory for compiled files.
- Force pushes the `action.yml`, a shim and the above files to the release's tag.
- Force pushes to the major version tag (ex: `v1.0.0` -> `v1`)
- Force pushes to the major minor version tag (ex: `v1.0.0` -> `v1.0`)

This action is meant to work with https://github.com/sethvargo/go-githubactions.

## Usage

```yaml
name: Publish

on:
  release:
    types: [published, edited]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v3
        with:
          ref: ${{ github.event.release.tag_name }}
      - name: Setup Golang
        uses: actions/setup-go@v4

      - name: Run tests
        run: |
          go test -count=1 -race -timeout=10m ./...

      - name: Build artifacts
        env:
          CGO_ENABLED: 0
        run: |
          go build -ldflags="-s -w" -o bin/main-${go env GOOS}-${go env GOARCH} ./cmd/...

      - uses: voidlock/assemble-and-tag-action@v1
        env:
          GITHUB_TOKEN: ${{ github.token }}
```

