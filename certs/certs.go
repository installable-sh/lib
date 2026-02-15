package certs

import (
	"crypto/x509"
	_ "embed"
	"encoding/pem"

	"github.com/installable-sh/lib/log"
)

//go:embed ca-certificates.crt
var CACerts []byte

// CertPool returns a certificate pool that combines system certificates
// with any embedded certificates. If system certs aren't available
// (e.g., in scratch containers), only embedded certs are used.
// Debug output is controlled by the logger's debug level.
func CertPool(logger log.DebugLogger) (*x509.CertPool, error) {
	// Start with system certs if available
	pool, err := x509.SystemCertPool()
	if err != nil || pool == nil {
		logger.Debugf("System cert pool unavailable: %v, using empty pool", err)
		pool = x509.NewCertPool()
	} else {
		logger.Debugf("Loaded system certificate pool")
	}

	// Append embedded certs if present
	if len(CACerts) > 0 {
		count := countPEMCerts(CACerts)
		logger.Debugf("Appending %d embedded CA certificates (%d bytes)", count, len(CACerts))
		pool.AppendCertsFromPEM(CACerts)
	} else {
		logger.Debugf("No embedded CA certificates")
	}

	return pool, nil
}

// countPEMCerts counts the number of certificates in PEM data.
func countPEMCerts(data []byte) int {
	count := 0
	for len(data) > 0 {
		var block *pem.Block
		block, data = pem.Decode(data)
		if block == nil {
			break
		}
		if block.Type == "CERTIFICATE" {
			count++
		}
	}
	return count
}
