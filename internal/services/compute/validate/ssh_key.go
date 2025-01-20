// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"crypto/ed25519"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"
)

// SSHKey performs some basic validation on supplied SSH Keys - Encoded Signature and Key Size are evaluated
// Will require rework if/when other Key Types are supported
func SSHKey(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if strings.TrimSpace(v) == "" {
		return nil, []error{fmt.Errorf("expected %q to not be an empty string or whitespace", k)}
	}

	keyParts := strings.Fields(v)
	if len(keyParts) > 1 {
		byteStr, err := base64.StdEncoding.DecodeString(keyParts[1])
		if err != nil {
			return nil, []error{fmt.Errorf("decoding %q for public key data", k)}
		}
		pubKey, err := ssh.ParsePublicKey(byteStr)
		if err != nil {
			return nil, []error{fmt.Errorf("parsing %q as a public key object", k)}
		}

		switch pubKey.Type() {
		case ssh.KeyAlgoRSA:
			rsaPubKey, ok := pubKey.(ssh.CryptoPublicKey).CryptoPublicKey().(*rsa.PublicKey)
			if !ok {
				return nil, []error{fmt.Errorf("- could not retrieve the RSA public key from the SSH public key")}
			}
			rsaPubKeyBits := rsaPubKey.Size() * 8
			if rsaPubKeyBits < 2048 {
				return nil, []error{fmt.Errorf("- the provided RSA SSH key has %d bits. Only ssh-rsa keys with 2048 bits or higher are supported by Azure", rsaPubKeyBits)}
			}
		case ssh.KeyAlgoED25519:
			_, ok := pubKey.(ssh.CryptoPublicKey).CryptoPublicKey().(ed25519.PublicKey)
			if !ok {
				return nil, []error{fmt.Errorf("- could not retrieve the ED25519 public key from the SSH public key")}
			}
		default:
			return nil, []error{fmt.Errorf("- the provided %s SSH key is not supported. Only RSA and ED25519 SSH keys are supported by Azure", pubKey.Type())}
		}
	} else {
		return nil, []error{fmt.Errorf("%q is not a complete SSH2 Public Key", k)}
	}

	return warnings, errors
}
