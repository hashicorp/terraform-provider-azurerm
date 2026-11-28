// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_custom_certificate

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *webpubsub.CustomCertificateId, model *webpubsub.CustomCertificate) error {
	if model == nil {
		return fmt.Errorf("retrieving %s: got nil model", *id)
	}

	vaultBasedUri := model.Properties.KeyVaultBaseUri
	certName := model.Properties.KeyVaultSecretName
	certVersion := pointer.From(model.Properties.KeyVaultSecretVersion)
	nestedItem, err := keyvault.NewNestedItemID(vaultBasedUri, keyvault.NestedItemTypeCertificate, certName, certVersion)
	if err != nil {
		return err
	}

	state := CustomCertWebPubsubModel{
		Name:               id.CustomCertificateName,
		CustomCertId:       nestedItem.ID(),
		WebPubsubId:        webpubsub.NewWebPubSubID(id.SubscriptionId, id.ResourceGroupName, id.WebPubSubName).ID(),
		CertificateVersion: certVersion,
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
