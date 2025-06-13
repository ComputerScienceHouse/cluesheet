## Contributing

1. Install go (on arch: `pacman -S go`)
2. In the root directory of the repo, run `podman compose up -d`
3. Obtain [golang migrate](https://github.com/golang-migrate/migrate). You may need to download the CLI directly from a release
4. In `backend`, run `migrate -source file://$(pwd)/migrations -database postgres://postgres:test@localhost:5432/postgres?sslmode=disable up`, on windows run `migrate -path ./migrations -database postgres://postgres:test@localhost:5432/postgres?sslmode=disable up`
5. In `backend`, run `go run .`

### Version tagging
If compiled with `-buildvcs`, traces and logs will include a commit hash. You can keep this set by running `go env -w GOFLAGS=-buildvcs`

### Adding new dependencies

`go mod tidy`

### Formatting

`gofmt -w -s .`

### Creating new DB migrations

Create a new file with an incremented prefix number, a human readable string, and ending with `.up.sql`, with a corresponding `.down.sql` that has the same number prefix to revert your change. Apply these with golang-migrate
