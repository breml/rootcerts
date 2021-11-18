// Package embedded makes available the "Mozilla Included CA Certificate List"
// without any side-effects (unlike package rootcerts).
package embedded

// MozillaCACertificatesPEM returns "Mozilla Included CA Certificate List"
// (https://wiki.mozilla.org/CA/Included_Certificates) in PEM format.
//
// Use of these certificates is governed by Mozilla Public License 2.0
// that can be found in the LICENSE.certificates file.
func MozillaCACertificatesPEM() string {
	return data
}
