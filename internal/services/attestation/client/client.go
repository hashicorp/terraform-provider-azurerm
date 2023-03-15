package client

import (
	"github.com/Azure/azure-sdk-for-go/services/attestation/2020-10-01/attestation"
	"github.com/hashicorp/go-azure-sdk/resource-manager/attestation/2020-10-01/attestationproviders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ProviderClient *attestationproviders.AttestationProvidersClient
	PolicyClient   *attestation.PolicyClient
}

func NewClient(o *common.ClientOptions) *Client {
	providerClient := attestationproviders.NewAttestationProvidersClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&providerClient.Client, o.ResourceManagerAuthorizer)

	policyClient := attestation.NewPolicyClient()
	o.ConfigureClient(&policyClient.Client, o.AttestationAuthorizer)

	return &Client{
		ProviderClient: &providerClient,
		PolicyClient:   &policyClient,
	}
}
