// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_custom_domain

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *webpubsub.CustomDomainId, model *webpubsub.CustomDomain) error {
	state := CustomDomainWebPubsubModel{
		Name:        id.CustomDomainName,
		WebPubsubId: webpubsub.NewWebPubSubID(id.SubscriptionId, id.ResourceGroupName, id.WebPubSubName).ID(),
	}

	if model != nil {
		props := model.Properties
		state.DomainName = props.DomainName

		if props.CustomCertificate.Id != nil {
			webPubsubCustomCertificateID, err := webpubsub.ParseCustomCertificateIDInsensitively(*props.CustomCertificate.Id)
			if err != nil {
				return err
			}
			state.WebPubsubCustomCertificateId = webPubsubCustomCertificateID.ID()
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
