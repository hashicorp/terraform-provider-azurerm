// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package CdnFrontDoorsecretparams

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CdnFrontDoorSecretParameters struct {
	TypeName   cdn.TypeBasicSecretParameters
	ConfigName string
}

type CdnFrontDoorSecretMappings struct {
	UrlSigningKey                     CdnFrontDoorSecretParameters
	ManagedCertificate                CdnFrontDoorSecretParameters
	CustomerCertificate               CdnFrontDoorSecretParameters
	AzureFirstPartyManagedCertificate CdnFrontDoorSecretParameters
}

func InitializeCdnFrontDoorSecretMappings() *CdnFrontDoorSecretMappings {
	m := new(CdnFrontDoorSecretMappings)

	m.UrlSigningKey = CdnFrontDoorSecretParameters{
		TypeName:   cdn.TypeBasicSecretParametersTypeURLSigningKey,
		ConfigName: "url_signing_key",
	}

	m.ManagedCertificate = CdnFrontDoorSecretParameters{
		TypeName:   cdn.TypeBasicSecretParametersTypeManagedCertificate,
		ConfigName: "managed_certificate",
	}

	m.CustomerCertificate = CdnFrontDoorSecretParameters{
		TypeName:   cdn.TypeBasicSecretParametersTypeCustomerCertificate,
		ConfigName: "customer_certificate",
	}

	m.AzureFirstPartyManagedCertificate = CdnFrontDoorSecretParameters{
		TypeName:   cdn.TypeBasicSecretParametersTypeAzureFirstPartyManagedCertificate,
		ConfigName: "azure_first_party_managed_certificate",
	}

	return m
}

func ExpandCdnFrontDoorCustomerCertificateParameters(ctx context.Context, input []interface{}, clients *clients.Client) (cdn.BasicSecretParameters, error) {
	m := InitializeCdnFrontDoorSecretMappings()
	item := input[0].(map[string]interface{})

	certificateBaseURL := item["key_vault_certificate_id"].(string)
	certificateId, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(certificateBaseURL)
	if err != nil {
		return nil, fmt.Errorf("retrieving the Key Vault Certificate Resource ID from the Key Vault Certificate Base URL %q: %s", certificateBaseURL, err)
	}

	var useLatest bool
	if certificateId.Version == "" {
		useLatest = true
	}

	subscriptionId := commonids.NewSubscriptionID(clients.Account.SubscriptionId)
	keyVaultBaseId, err := clients.KeyVault.KeyVaultIDFromBaseUrl(ctx, subscriptionId, certificateId.KeyVaultBaseUrl)
	if err != nil {
		return nil, fmt.Errorf("retrieving the Key Vault Resource ID from the Key Vault Base URL %q: %s", certificateId.KeyVaultBaseUrl, err)
	}

	if keyVaultBaseId == nil {
		return nil, fmt.Errorf("unexpected nil Key Vault Resource ID retrieved from the Key Vault Base URL %q", certificateId.KeyVaultBaseUrl)
	}

	keyVaultId, err := commonids.ParseKeyVaultID(*keyVaultBaseId)
	if err != nil {
		return nil, err
	}

	secretSource := keyVaultParse.NewSecretVersionlessID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroupName, keyVaultId.VaultName, certificateId.Name)

	customerCertificate := &cdn.CustomerCertificateParameters{
		Type: m.CustomerCertificate.TypeName,
		SecretSource: &cdn.ResourceReference{
			ID: utils.String(secretSource.ID()),
		},
		UseLatestVersion: utils.Bool(useLatest),
	}

	if !useLatest {
		customerCertificate.SecretVersion = utils.String(certificateId.Version)
	}

	return customerCertificate, nil
}
