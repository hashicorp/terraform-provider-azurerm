// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_management_private_link_association

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-05-01/privatelinkassociation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.PrivateLinkAssociationClient

			var config ResourceManagementPrivateLinkAssociationResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			managementGroupId, err := commonids.ParseManagementGroupID(config.ManagementGroupId)
			if err != nil {
				return err
			}

			var name string
			if config.Name != "" {
				name = config.Name
			}

			if name == "" {
				name = uuid.New().String()
			}

			id := privatelinkassociation.NewPrivateLinkAssociationID(managementGroupId.GroupId, name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			payload := privatelinkassociation.PrivateLinkAssociationObject{
				Properties: &privatelinkassociation.PrivateLinkAssociationProperties{
					PrivateLink:         pointer.To(config.ResourceManagementPrivateLinkId),
					PublicNetworkAccess: r.expandPublicNetworkAccess(config.PublicNetworkAccessEnabled),
				},
			}

			if _, err := client.Put(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r Resource) expandPublicNetworkAccess(input bool) *privatelinkassociation.PublicNetworkAccessOptions {
	output := privatelinkassociation.PublicNetworkAccessOptionsEnabled
	if !input {
		output = privatelinkassociation.PublicNetworkAccessOptionsDisabled
	}

	return &output
}
