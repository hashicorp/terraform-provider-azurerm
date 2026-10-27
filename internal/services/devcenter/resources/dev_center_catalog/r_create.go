// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_catalog

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/catalogs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model DevCenterCatalogsResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			client := metadata.Client.DevCenter.V20250201.Catalogs
			subscriptionId := metadata.Client.Account.SubscriptionId
			devCenterId, err := catalogs.ParseDevCenterID(model.DevCenterID)
			if err != nil {
				return fmt.Errorf("parsing dev center id: %+v", err)
			}
			devCenterName := devCenterId.DevCenterName
			id := catalogs.NewDevCenterCatalogID(subscriptionId, model.ResourceGroupName, devCenterName, model.Name)

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

			catalogProperties := catalogs.Catalog{
				Properties: &catalogs.CatalogProperties{
					AdoGit: expandCatalogProperties(model.CatalogAdoGit),
					GitHub: expandCatalogProperties(model.CatalogGitHub),
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, catalogProperties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandCatalogProperties(input []CatalogPropertiesModel) *catalogs.GitCatalog {
	if len(input) == 0 {
		return nil
	}

	return &catalogs.GitCatalog{
		Uri:              pointer.To(input[0].URI),
		Branch:           pointer.To(input[0].Branch),
		SecretIdentifier: pointer.To(input[0].KeyVaultKeyUrl),
		Path:             pointer.To(input[0].Path),
	}
}
