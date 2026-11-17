// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventhub_consumer_group

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/consumergroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *consumergroups.ConsumerGroupId, model *consumergroups.ConsumerGroup) error {
	state := ConsumerGroupObject{
		Name:              id.ConsumerGroupName,
		NamespaceName:     id.NamespaceName,
		EventHubName:      id.EventhubName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil && model.Properties != nil {
		state.UserMetadata = pointer.From(model.Properties.UserMetadata)
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
