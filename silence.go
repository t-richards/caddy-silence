package caddyhttp

import (
	"bufio"
	"net"
	"sync"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

// from crypto/tls.
type recordType uint8

const (
	recordTypeChangeCipherSpec recordType = 20
	recordTypeAlert            recordType = 21
	recordTypeHandshake        recordType = 22
	recordTypeApplicationData  recordType = 23
)

func init() {
	caddy.RegisterModule(HTTPSilenceWrapper{})
}

// HTTPSilenceWrapper tells Caddy to shut up.
type HTTPSilenceWrapper struct{}

func (HTTPSilenceWrapper) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "caddy.listeners.silence",
		New: func() caddy.Module { return new(HTTPSilenceWrapper) },
	}
}

func (h *HTTPSilenceWrapper) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	return nil
}

func (h *HTTPSilenceWrapper) WrapListener(l net.Listener) net.Listener {
	return &httpSilenceListener{l}
}

// httpSilenceListener is listener that checks the first byte
// of the request when the server is intended to accept HTTPS requests,
// to immediately close any misdirected HTTP requests.
type httpSilenceListener struct {
	net.Listener
}

// Accept waits for and returns the next connection to the listener,
// wrapping it with a httpSilenceConn.
func (l *httpSilenceListener) Accept() (net.Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}

	return &httpSilenceConn{
		Conn: c,
		r:    bufio.NewReader(c),
	}, nil
}

type httpSilenceConn struct {
	net.Conn
	once sync.Once
	r    *bufio.Reader
}

// Read tries to peek at the first byte of the request.
// If it does not look like a valid TLS record header, we close the connection.
func (c *httpSilenceConn) Read(p []byte) (int, error) {
	c.once.Do(func() {
		firstBytes, err := c.r.Peek(1)
		if err != nil {
			return
		}

		// If the request looks like HTTP, then we silence it.
		if !validTLSRecordHeader(recordType(firstBytes[0])) {
			c.Conn.Close()
			return
		}
	})

	return c.r.Read(p)
}

// validTLSRecordHeader reports whether the given record type
// exists in the set of valid TLS record types.
func validTLSRecordHeader(r recordType) bool {
	switch r {
	case recordTypeChangeCipherSpec, recordTypeAlert, recordTypeHandshake, recordTypeApplicationData:
		return true
	}
	return false
}

var (
	_ caddy.ListenerWrapper = (*HTTPSilenceWrapper)(nil)
	_ caddyfile.Unmarshaler = (*HTTPSilenceWrapper)(nil)
)
