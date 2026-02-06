package certs

import (
	"crypto/x509"
	_ "embed"
)

//go:embed ca-certificates.crt
var CACerts []byte

// CertPool returns a certificate pool that combines system certificates
// with any embedded certificates. If system certs aren't available
// (e.g., in scratch containers), only embedded certs are used.
func CertPool() (*x509.CertPool, error) {
	// Start with system certs if available
	pool, err := x509.SystemCertPool()
	if err != nil || pool == nil {
		pool = x509.NewCertPool()
	}

	// Append embedded certs if present
	if len(CACerts) > 0 {
		pool.AppendCertsFromPEM(CACerts)
	}

	return pool, nil
}
