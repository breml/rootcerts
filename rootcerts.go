//go:generate go run generate_data.go

// Package rootcerts provides an embedded copy of the "Mozilla Included CA
// Certificate List" (https://wiki.mozilla.org/CA/Included_Certificates),
// more specifically the "PEM of Root Certificates in Mozilla's Root Store with
// the Websites (TLS/SSL) Trust Bit Enabled"
// (https://ccadb-public.secure.force.com/mozilla/IncludedRootsPEMTxt?TrustBitsInclude=Websites).
// The "Mozilla Included CA Certificate List" is maintained as part of the
// Common CA Database effort (https://golang.org/pkg/crypto/x509/).
// If this package is imported anywhere in the program, then if the crypto/x509
// package cannot find the system certificate pool, it will use this embedded
// information.
//
// Additionally, the usage of this embedded information can be forced by setting
// the the environment variable `GO_ROOTCERTS_ENABLE=1` while  running a
// program, which includes this package.
//
// Importing this package will increase the size of a program by about 250 KB.
//
// This package should normally be imported by a program's main package, not by
// a library. Libraries normally shouldn't decide whether to include the
// "Mozilla Included CA Certificate List" in a program.
package rootcerts

import (
	"crypto/x509"
	"os"
	_ "unsafe" // for go:linkname
)

const forceEnableEnvVar = "GO_ROOTCERTS_ENABLE"

//go:linkname systemRoots crypto/x509.systemRoots
var systemRoots *x509.CertPool

func init() {
	// Ensure x509.SystemCertPool is executed once
	x509.SystemCertPool() // nolint: errcheck

	if systemRoots != nil && len(systemRoots.Subjects()) > 0 && os.Getenv(forceEnableEnvVar) != "1" {
		return
	}

	roots := x509.NewCertPool()
	d := data
	roots.AppendCertsFromPEM([]byte(d))
	systemRoots = roots
}
