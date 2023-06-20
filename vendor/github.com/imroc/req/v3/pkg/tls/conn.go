package tls

import (
	"crypto/tls"
	"net"
)

// Conn is the recommended interface for the connection
// returned by the DailTLS function (Client.SetDialTLS,
// Transport.DialTLSContext), so that the TLS handshake negotiation
// can automatically decide whether to use HTTP2 or HTTP1 (ALPN).
// If this interface is not implemented, HTTP1 will be used by default.
type Conn interface {
	net.Conn
	// ConnectionState returns basic TLS details about the connection.
	ConnectionState() tls.ConnectionState
	// Handshake runs the client or server handshake
	// protocol if it has not yet been run.
	//
	// Most uses of this package need not call Handshake explicitly: the
	// first Read or Write will call it automatically.
	//
	// For control over canceling or setting a timeout on a handshake, use
	// HandshakeContext or the Dialer's DialContext method instead.
	Handshake() error
}
