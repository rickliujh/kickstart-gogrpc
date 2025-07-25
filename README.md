# kickstart-gogrpc
Boilerplate of GRPC server for go

### TODOs
- ~~sqlc~~
- ~~migrate~~
- ~~logr~~
- Otle
- AAA
- Error handler (middleware)

## Decisions

### No pre-commit hooks
The boilerplate do not provide pre-commit hooks of all sorts. Allow fast local commits is a major upside of DVCS. Clean commit can be achieved during pre-push or during pre-merge at PR level with CI scan.

See: [Reddit Discussion](https://www.reddit.com/r/rust/comments/13h78se/comment/jk3mgpx/?utm_source=share&utm_medium=web3x&utm_name=web3xcss&utm_term=1&utm_content=share_button)

## Guides

### Build
Two type of build:
1. `make build-dev`: static link but with debug info reserved
1. `make build-prod`: static link, omit the DWARF symbol, omit the symbol table `-ldflags="-w -s"`

### Dockerfile

#### Choose the type of build 
`docker build -t kickstart-go:v0.0.1 --build-arg BUILD_TYPE=prod .`

Allowed value: one of [prod, dev]

#### Specify name of the binary being built
`docker build -t kickstart-go:v0.0.1 --build-arg APP_NAME=my-app .`

> [!NOTE]
> In the final image the binary will always be called `app ` and located at `/`.
> To change the name of the binary in the final image you need:
> 1. change the destination of the copy path
> ```
> COPY --from=builder /go/src/app/bin/${APP_NAME} /my-app
> ```
> 2. change the endpoint
> ```
> ENTRYPOINT ["/my-app"]
> ```

#### Specify the image that will run the binary
`docker build -t kickstart-go:v0.0.1 --build-arg RUN_IMG=golang:alpine .`

Default value: scratch

### Logging
[Logging in Go with Slog](https://betterstack.com/community/guides/logging/logging-in-go/)

## References

- [grpcme | Simple go gRPC service template](https://github.com/mchmarny/grpcme)
