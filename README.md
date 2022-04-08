# caddy-silence

Tells Caddy to shut up.

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
	servers :8443 {
		listener_wrappers {
			silence
			tls
		}
	}
}
```

## License

Apache 2.
