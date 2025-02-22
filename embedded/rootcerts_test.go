package embedded_test

import (
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"github.com/breml/rootcerts/embedded"
)

func parsePEM(pemCerts []byte) (certs []*x509.Certificate, err error) {
	for len(pemCerts) > 0 {
		var block *pem.Block
		block, pemCerts = pem.Decode(pemCerts)
		if block == nil {
			break
		}
		if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
			continue
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
		certs = append(certs, cert)
	}
	return
}

func checkRootCertsPEM(t *testing.T, pemCerts []byte, whenFail time.Time, whenWarn time.Time) (ok bool) {
	const warnEmoji = "\u26a0\ufe0f"
	// t.Logf("%#v %[1]x %x", warnEmoji, []rune(warnEmoji))
	now := time.Now()
	t.Logf("Checking certificate validity on %s...", whenFail)
	certs, err := parsePEM(pemCerts)
	if err != nil {
		t.Error(err)
		return false
	}

	roots := x509.NewCertPool()
	for _, cert := range certs {
		roots.AddCert(cert)
	}

	var minExpires time.Time
	ok = true
	for _, cert := range certs {
		name := cert.Subject.CommonName
		if name == "" {
			name = cert.Subject.String() + " (⚠️ missing CommonName)"
			if name == "" {
				name = cert.Issuer.String()
			}
		}

		if !cert.IsCA {
			t.Errorf("\u274C %s: not a certificate authority", name)
		}
		const keyUsageExpected = x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature
		if (cert.KeyUsage &^ keyUsageExpected) != 0 {
			t.Logf(warnEmoji+" %s: unexpected key usage %#x (expecting %#x, see constants at https://pkg.go.dev/crypto/x509#KeyUsage)", name, cert.KeyUsage, keyUsageExpected)
		}
		if minExpires.IsZero() || cert.NotAfter.Before(minExpires) {
			minExpires = cert.NotAfter
		}
		// Check that the certificate is valid now
		if cert.NotBefore.After(now) {
			t.Errorf("\u274C %s: fails NotBefore check: %s", name, cert.NotBefore)
			continue
		}
		// ... and that it will still be valid later
		if cert.NotAfter.Before(whenFail) {
			t.Errorf("\u274C %s: fails NotAfter check: %s", name, cert.NotAfter)
			continue
		} else if cert.NotAfter.Before(whenWarn) {
			t.Logf(warnEmoji+" %s: fails NotAfter check: %s", name, cert.NotAfter)
		}
		_, err := cert.Verify(x509.VerifyOptions{
			Roots:       roots,
			CurrentTime: whenFail,
		})
		if err != nil {
			t.Errorf("\u274C %s: %s", name, err)
			ok = false
		} else {
			t.Logf("\u2705 %s (expires: %s)", name, cert.NotAfter)
		}
	}
	if ok {
		t.Log("Success.")
		t.Logf("MinExpire: %s", minExpires)
	}
	return
}

func TestCerts(t *testing.T) {
	// Check that certificates will still be valid in 1 month, warn if invalid in 3 months
	checkRootCertsPEM(t, []byte(embedded.MozillaCACertificatesPEM()), time.Now().AddDate(0, 1, 0), time.Now().AddDate(0, 3, 0))

	// Should fail
	//checkRootCertsPEM(t, []byte(embedded.MozillaCACertificatesPEM()), time.Now().AddDate(20, 0, 0), time.Now().AddDate(30, 0, 0))
}
