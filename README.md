# bun starter kit

[![build workflow](https://ecampus-be/actions/workflows/build.yml/badge.svg)](https://ecampus-be/actions)

Bun starter kit consists of:

- [bunrouter](https://github.com/uptrace/bunrouter) is an extremely fast and flexible router.
- [bun](https://github.com/uptrace/bun)
- Hooks to initialize the app.
- CLI to run HTTP server and migrations, for example, `go run cmd/bun/*.go db help`.
- [example](ecampus) package that shows how to load fixtures and test handlers.

You can also check [bun-realworld-app](https://github.com/go-bun/bun-realworld-app) which is a JSON
API built with Bun starter kit.

## Quickstart

To start using this kit, clone the repo:

```shell
git clone https://ecampus-be.git
```

Make sure you have correct information in `app/config/test.yaml` and then run migrations (database
must exist before running):

```shell
go run cmd/bun/main.go -env=dev db init
go run cmd/bun/main.go -env=dev db migrate
```

To start the server:

```shell
go run cmd/bun/main.go -env=dev runserver
```

Then run the tests in [example](ecampus) package:

```shell
cd ecampus
go test
```

See [documentation](https://bun.uptrace.dev/guide/starter-kit.html) for more info.
