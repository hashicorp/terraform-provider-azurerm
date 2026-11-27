// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package signalr_service_custom_domain

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2024-03-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *signalr.CustomDomainId, model *signalr.CustomDomain) error {
	state := CustomDomainSignalrServiceModel{
		Name:             id.CustomDomainName,
		SignalRServiceId: signalr.NewSignalRID(id.SubscriptionId, id.ResourceGroupName, id.SignalRName).ID(),
	}

	if model != nil {
		props := model.Properties
		state.DomainName = props.DomainName

		if props.CustomCertificate.Id != nil {
			signalrCustomCertificateID, err := signalr.ParseCustomCertificateIDInsensitively(*props.CustomCertificate.Id)
			if err != nil {
				return fmt.Errorf("parsing signalr custom cert id for %s: %+v", *id, err)
			}
			state.SignalrCustomCertificateId = signalrCustomCertificateID.ID()
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
