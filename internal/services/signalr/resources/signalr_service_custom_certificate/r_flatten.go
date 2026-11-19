// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package signalr_service_custom_certificate

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2024-03-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *signalr.CustomCertificateId, cert *signalr.CustomCertificate) error {
	state := CustomCertSignalrServiceResourceModel{
		Name:             id.CustomCertificateName,
		SignalRServiceId: signalr.NewSignalRID(id.SubscriptionId, id.ResourceGroupName, id.SignalRName).ID(),
	}

	if cert != nil {
		certVersion := pointer.From(cert.Properties.KeyVaultSecretVersion)
		nestedItem, err := keyvault.NewNestedItemID(cert.Properties.KeyVaultBaseUri, keyvault.NestedItemTypeCertificate, cert.Properties.KeyVaultSecretName, certVersion)
		if err != nil {
			return err
		}
		state.CustomCertId = nestedItem.ID()
		state.CertificateVersion = certVersion
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
