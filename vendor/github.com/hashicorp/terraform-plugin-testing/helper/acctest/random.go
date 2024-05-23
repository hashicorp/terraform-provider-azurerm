// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package acctest

import (
	"bytes"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"math/rand"
	"net/netip"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// Helpers for generating random tidbits for use in identifiers to prevent
// collisions in acceptance tests.

// RandInt generates a random integer
func RandInt() int {
	return rand.Int()
}

// RandomWithPrefix is used to generate a unique name with a prefix, for
// randomizing names in acceptance tests
func RandomWithPrefix(name string) string {
	return fmt.Sprintf("%s-%d", name, RandInt())
}

// RandIntRange returns a random integer between min (inclusive) and max (exclusive)
func RandIntRange(min int, max int) int {
	return rand.Intn(max-min) + min
}

// RandString generates a random alphanumeric string of the length specified
func RandString(strlen int) string {
	return RandStringFromCharSet(strlen, CharSetAlphaNum)
}

// RandStringFromCharSet generates a random string by selecting characters from
// the charset provided
func RandStringFromCharSet(strlen int, charSet string) string {
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = charSet[RandIntRange(0, len(charSet))]
	}
	return string(result)
}

// RandSSHKeyPair generates a random public and private SSH key pair.
//
// The public key is returned in OpenSSH authorized key format, for example:
//
//	ssh-rsa XXX comment
//
// The private key is RSA algorithm, 1024 bits, PEM encoded, and has no
// passphrase. Testing with different or stricter security requirements should
// use the standard library [crypto] and [golang.org/x/crypto/ssh] packages
// directly.
func RandSSHKeyPair(comment string) (string, string, error) {
	privateKey, privateKeyPEM, err := genPrivateKey()
	if err != nil {
		return "", "", err
	}

	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}
	keyMaterial := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(publicKey)))
	return fmt.Sprintf("%s %s", keyMaterial, comment), privateKeyPEM, nil
}

// RandTLSCert generates a self-signed TLS certificate with a newly created
// private key, and returns both the cert and the private key PEM encoded.
//
// The private key uses RSA algorithm, 1024 bits, and has no passphrase.
//
// The certificate expires in 24 hours, has a random serial number, and is
// set for Encipherment, Digital Signature, and Server Auth key usage.
// Only the organization name of the subject is configurable.
//
// Testing with different or stricter security requirements should
// use the standard library [crypto] and [golang.org/x/crypto] packages
// directly.
func RandTLSCert(orgName string) (string, string, error) {
	template := &x509.Certificate{
		SerialNumber: big.NewInt(int64(RandInt())),
		Subject: pkix.Name{
			Organization: []string{orgName},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	privateKey, privateKeyPEM, err := genPrivateKey()
	if err != nil {
		return "", "", err
	}

	cert, err := x509.CreateCertificate(crand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return "", "", err
	}

	certPEM, err := pemEncode(cert, "CERTIFICATE")
	if err != nil {
		return "", "", err
	}

	return certPEM, privateKeyPEM, nil
}

// RandIpAddress returns a random IP address in the specified CIDR block.
func RandIpAddress(s string) (string, error) {
	prefix, err := netip.ParsePrefix(s)
	if err != nil {
		return "", err
	}

	if prefix.IsSingleIP() {
		return prefix.Addr().String(), nil
	}

	// base address as byte slice
	prefixBytes, err := prefix.Masked().Addr().MarshalBinary()
	if err != nil {
		return "", err
	}

	// inverse mask (ones in the host bits) as byte slice
	inverseMaskBytes, err := inverseMask(prefix.Bits(), len(prefixBytes))
	if err != nil {
		return "", err
	}

	// the result starts life as 4 or 16 bytes of random data
	resultBytes := make([]byte, len(inverseMaskBytes))
	_, err = crand.Read(resultBytes)
	if err != nil {
		return "", err
	}

	// use the prefix and inverse mask to restore the network bits
	for i := range inverseMaskBytes {
		resultBytes[i] = (resultBytes[i] & inverseMaskBytes[i]) + prefixBytes[i]
	}

	result, ok := netip.AddrFromSlice(resultBytes)
	if !ok {
		return "", fmt.Errorf("unable to create random address from bytes: %#v", resultBytes)
	}

	return result.String(), nil
}

func genPrivateKey() (*rsa.PrivateKey, string, error) {
	privateKey, err := rsa.GenerateKey(crand.Reader, 1024)
	if err != nil {
		return nil, "", err
	}

	privateKeyPEM, err := pemEncode(x509.MarshalPKCS1PrivateKey(privateKey), "RSA PRIVATE KEY")
	if err != nil {
		return nil, "", err
	}

	return privateKey, privateKeyPEM, nil
}

func pemEncode(b []byte, block string) (string, error) {
	var buf bytes.Buffer
	pb := &pem.Block{Type: block, Bytes: b}
	if err := pem.Encode(&buf, pb); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func inverseMask(bits, byteLen int) ([]byte, error) {
	if bits > byteLen*8 {
		return nil, fmt.Errorf("cannot fit a %d-bit mask into %d bytes", bits, byteLen)
	}

	iBits := (byteLen * 8) - bits
	var result []byte
	for iBits > 0 {
		b := uint8((1 << iBits) - 1)
		result = append([]byte{b}, result...)
		iBits -= 8
	}

	return append(make([]byte, byteLen-len(result)), result...), nil
}

const (
	// CharSetAlphaNum is the alphanumeric character set for use with
	// RandStringFromCharSet
	CharSetAlphaNum = "abcdefghijklmnopqrstuvwxyz012346789"

	// CharSetAlpha is the alphabetical character set for use with
	// RandStringFromCharSet
	CharSetAlpha = "abcdefghijklmnopqrstuvwxyz"
)
