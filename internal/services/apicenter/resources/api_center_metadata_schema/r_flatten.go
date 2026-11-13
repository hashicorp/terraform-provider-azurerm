// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_metadata_schema

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/metadataschemas"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/services"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *metadataschemas.MetadataSchemaId, model *metadataschemas.MetadataSchema) error {
	if model == nil {
		return fmt.Errorf("retrieving %s: got nil model", *id)
	}
	state := ApiCenterMetadataSchemaModel{
		Name:      id.MetadataSchemaName,
		ServiceId: services.NewServiceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName).ID(),
	}

	if props := model.Properties; props != nil {
		state.MetadataAssignment = flattenMetadataAssignmentList(props.AssignedTo)
		state.Schema = props.Schema
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}
	return metadata.Encode(&state)
}

func flattenMetadataAssignmentList(input *[]metadataschemas.MetadataAssignment) []MetadataAssignment {
	if input == nil {
		return []MetadataAssignment{}
	}
	result := make([]MetadataAssignment, 0, len(*input))
	for _, v := range *input {
		result = append(result, MetadataAssignment{
			Deprecated: pointer.From(v.Deprecated),
			Entity:     pointer.FromEnum(v.Entity),
			Required:   pointer.From(v.Required),
		})
	}
	return result
}
