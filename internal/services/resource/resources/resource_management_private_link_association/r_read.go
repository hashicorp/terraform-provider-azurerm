// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_management_private_link_association

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-05-01/privatelinkassociation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.PrivateLinkAssociationClient
			schema := ResourceManagementPrivateLinkAssociationResourceSchema{}

			id, err := privatelinkassociation.ParsePrivateLinkAssociationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			schema.ManagementGroupId = commonids.NewManagementGroupID(id.GroupId).ID()
			schema.Name = id.PlaId

			if model := resp.Model; model != nil {
				if prop := model.Properties; prop != nil {
					schema.PublicNetworkAccessEnabled = r.flattenPublicNetworkAccess(prop.PublicNetworkAccess)
					schema.TenantID = pointer.From(prop.TenantID)
					schema.ResourceManagementPrivateLinkId = pointer.From(prop.PrivateLink)
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r Resource) flattenPublicNetworkAccess(input *privatelinkassociation.PublicNetworkAccessOptions) bool {
	if input == nil || *input == privatelinkassociation.PublicNetworkAccessOptionsEnabled {
		return true
	}

	return false
}

type ResourceManagementPrivateLinkAssociationResourceSchema struct {
	ManagementGroupId               string `tfschema:"management_group_id"`
	Name                            string `tfschema:"name"`
	ResourceManagementPrivateLinkId string `tfschema:"resource_management_private_link_id"`
	PublicNetworkAccessEnabled      bool   `tfschema:"public_network_access_enabled"`
	TenantID                        string `tfschema:"tenant_id"`
}
