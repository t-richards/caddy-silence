# caddy-silence

Tells Caddy to shut up.

When plain HTTP requests are issued against HTTPS servers, the Go standard library will emit the following response:

```
HTTP/1.0 400 Bad Request

Client sent an HTTP request to an HTTPS server.
```

This module eliminates this response entirely by closing the connection when it detects that a plain HTTP request has been sent to an HTTPS listener.

## Getting started

Build:

```bash
export CADDY_VERSION=master
go install github.com/caddyserver/xcaddy/cmd/xcaddy@latest
xcaddy build --with=github.com/t-richards/caddy-silence
```

Configure:

```
# Global config block
{
	servers :443 {
		listener_wrappers {
			silence
			tls
		}
	}
}
```

## License

[Apache 2.0](./NOTICE).
