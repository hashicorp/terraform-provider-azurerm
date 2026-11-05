// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automanage_configuration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automanage.ConfigurationProfilesClient

			id, err := configurationprofiles.ParseConfigurationProfileID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			properties := configurationprofiles.ConfigurationProfile{
				Location: location.Normalize(model.Location),
				Properties: &configurationprofiles.ConfigurationProfileProperties{
					Configuration: expandConfigurationProfile(model),
				},
				Tags: pointer.To(model.Tags),
			}

			if _, err := client.CreateOrUpdate(ctx, *id, properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
