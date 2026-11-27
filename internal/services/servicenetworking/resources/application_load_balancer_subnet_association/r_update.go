// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_load_balancer_subnet_association

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2025-01-01/associationsinterface"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	servicenetworkingClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicenetworking/client"
)

func (t Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := servicenetworkingClient.From(metadata.Client).AssociationsInterface

			var config AssociationResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding %v", err)
			}

			id, err := associationsinterface.ParseAssociationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// Although `AssociationSubnetUpdate` defined in the SDK contains the `subnetId`, while per testing it can not be updated
			// Tracked on https://github.com/Azure/azure-rest-api-specs/issues/26657
			associationUpdate := associationsinterface.AssociationUpdate{}

			if metadata.ResourceData.HasChange("tags") {
				associationUpdate.Tags = tags.Expand(config.Tags)
			}

			if _, err = client.Update(ctx, *id, associationUpdate); err != nil {
				return fmt.Errorf("updating %s: %v", *id, err)
			}

			return nil
		},
	}
}
