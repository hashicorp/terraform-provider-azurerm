// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_dev_box_definition

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/devboxdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.DevBoxDefinitions

			id, err := devboxdefinitions.ParseDevCenterDevBoxDefinitionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model DevCenterDevBoxDefinitionResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := devboxdefinitions.DevBoxDefinitionUpdate{
				Properties: &devboxdefinitions.DevBoxDefinitionUpdateProperties{},
			}

			if metadata.ResourceData.HasChange("image_reference_id") {
				parameters.Properties.ImageReference = &devboxdefinitions.ImageReference{
					Id: pointer.To(model.ImageReferenceId),
				}
			}

			if metadata.ResourceData.HasChange("hibernate_support_enabled") {
				hs := devboxdefinitions.HibernateSupportDisabled
				if model.HibernateSupport {
					hs = devboxdefinitions.HibernateSupportEnabled
				}
				parameters.Properties.HibernateSupport = pointer.To(hs)
			}

			if metadata.ResourceData.HasChange("sku_name") {
				parameters.Properties.Sku = expandDevCenterDevBoxDefinitionSku(model.SkuName)
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = pointer.To(model.Tags)
			}

			if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
