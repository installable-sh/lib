# lib

Shared Go libraries for [installable.sh](https://github.com/installable-sh).

## Installation

```bash
go get github.com/installable-sh/lib
```

## Packages

### certs

Embedded CA certificates from Alpine Linux, updated weekly.

```go
import "github.com/installable-sh/lib/certs"

pool := certs.CertPool()
```

### fetch

HTTP client with retries and embedded CA certificates.

```go
import "github.com/installable-sh/lib/fetch"

resp, err := fetch.Get(ctx, "https://example.com")
```

### shell

Shell script execution using mvdan.cc/sh.

```go
import "github.com/installable-sh/lib/shell"

err := shell.Run(ctx, "echo hello", os.Stdout, os.Stderr)
```

### version

Version information utilities.

```go
import "github.com/installable-sh/lib/version"

v := version.Get()
```

## License

Apache 2.0 - see [LICENSE](LICENSE)
