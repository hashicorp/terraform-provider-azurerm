// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/attestation/2020-10-01/attestationproviders"
	authWrapper "github.com/hashicorp/go-azure-sdk/sdk/auth/autorest"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/tombuildsstuff/kermit/sdk/attestation/2022-08-01/attestation"
)

type Client struct {
	ProviderClient *attestationproviders.AttestationProvidersClient

	o *common.ClientOptions
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	providerClient, err := attestationproviders.NewAttestationProvidersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Providers client: %+v", err)
	}
	o.Configure(providerClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ProviderClient: providerClient,
		o:              o,
	}, nil
}

// DataPlaneEndpointForProvider returns the Data Plane endpoint for the specified Attestation Provider
func (c *Client) DataPlaneEndpointForProvider(ctx context.Context, id attestationproviders.AttestationProvidersId) (*string, error) {
	// NOTE: there's potential performance benefits to caching this, but since this is used in a single resource for now it's not really needed
	existing, err := c.ProviderClient.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	dataPlaneUri := ""
	if model := existing.Model; model != nil && model.Properties != nil && model.Properties.AttestUri != nil {
		dataPlaneUri = *model.Properties.AttestUri
	}
	if dataPlaneUri == "" {
		return nil, fmt.Errorf("retrieving %s: unable to determining the Data Plane URI `model.Properties.AttestUri` was nil", id)
	}
	return &dataPlaneUri, nil
}

// DataPlaneClientWithEndpoint returns a DataPlaneClient for the given Attestation Provider Data Plane endpoint
func (c *Client) DataPlaneClientWithEndpoint(endpoint string) (*attestation.PolicyClient, error) {
	// the endpoint is in the format `https://acctestapzllwo64ym0.uks.attest.azure.net`
	// however the authorization token is needed for `https://attest.azure.net`, so we'll want
	// to compute that (since it varies per environment)
	uri, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as a URI: %+v", endpoint, err)
	}
	// trim off the first two segments (name and dc)
	segments := strings.Split(uri.Host, ".")
	if len(segments) >= 2 {
		segments = segments[2:]
	}
	authTokenUri := fmt.Sprintf("https://%s/", strings.Join(segments, "."))
	domainSuffix, ok := c.o.Environment.Attestation.DomainSuffix()
	if !ok {
		return nil, fmt.Errorf("building Authorizer for %q: domain suffix for Attestation service could not be determined", endpoint)
	}
	api := environments.AttestationAPI(authTokenUri, *domainSuffix)
	auth, err := c.o.Authorizers.AuthorizerFunc(api)
	if err != nil {
		return nil, fmt.Errorf("building Authorizer for %q: %+v", endpoint, err)
	}

	policyClient := attestation.NewPolicyClient()
	policyClient.RetryAttempts = 5
	c.o.ConfigureClient(&policyClient.Client, authWrapper.AutorestAuthorizer(auth))
	return &policyClient, nil
}
