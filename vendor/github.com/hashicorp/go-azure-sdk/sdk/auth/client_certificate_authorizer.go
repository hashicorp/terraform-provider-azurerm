// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"software.sslmate.com/src/go-pkcs12"
)

type ClientCertificateAuthorizerOptions struct {
	// Environment is the Azure environment/cloud being targeted
	Environment environments.Environment

	// Api describes the Azure API being used
	Api environments.Api

	// TenantId is the tenant to authenticate against
	TenantId string

	// AuxTenantIds lists additional tenants to authenticate against, currently only
	// used for Resource Manager when auxiliary tenants are needed.
	// e.g. https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/authenticate-multi-tenant
	AuxTenantIds []string

	// ClientId is the client ID used when authenticating
	ClientId string

	// Pkcs12Data is the binary PKCS#12 archive data containing the certificate and private key
	Pkcs12Data []byte

	// Pkcs12Path is a path to a binary PKCS#12 archive on the filesystem
	Pkcs12Path string

	// Pkcs12Pass is the challenge passphrase to decrypt the PKCS#12 archive
	Pkcs12Pass string
}

// NewClientCertificateAuthorizer returns an authorizer which uses client certificate authentication.
func NewClientCertificateAuthorizer(ctx context.Context, options ClientCertificateAuthorizerOptions) (Authorizer, error) {
	if len(options.Pkcs12Data) == 0 {
		var err error
		options.Pkcs12Data, err = os.ReadFile(options.Pkcs12Path)
		if err != nil {
			return nil, fmt.Errorf("could not read PKCS#12 archive at %q: %s", options.Pkcs12Path, err)
		}
	}

	// we aren't interested in the issuer chain, but we use the DecodeChain method to parse them out in case they are present
	key, cert, _, err := pkcs12.DecodeChain(options.Pkcs12Data, options.Pkcs12Pass)
	if err != nil {
		return nil, fmt.Errorf("could not decode PKCS#12 archive: %s", err)
	}

	scope, err := environments.Scope(options.Api)
	if err != nil {
		return nil, fmt.Errorf("determining scope for %q: %+v", options.Api.Name(), err)
	}

	conf := clientCredentialsConfig{
		Environment:        options.Environment,
		TenantID:           options.TenantId,
		AuxiliaryTenantIDs: options.AuxTenantIds,
		ClientID:           options.ClientId,
		PrivateKey:         key,
		Certificate:        cert,
		Scopes: []string{
			*scope,
		},
	}
	return conf.TokenSource(ctx, clientCredentialsAssertionType)
}
