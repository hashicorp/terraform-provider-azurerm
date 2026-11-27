// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_subscription_pricing

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2023-01-01/pricings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSecurityCenterSubscriptionPricingRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).SecurityCenter.PricingClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := pricings.ParsePricingID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("resource_type", id.PricingName)
	if resp.Model != nil {
		if properties := resp.Model.Properties; properties != nil {
			d.Set("tier", properties.PricingTier)
			d.Set("subplan", properties.SubPlan)
			if err = d.Set("extension", flattenExtensions(properties.Extensions)); err != nil {
				return fmt.Errorf("setting `extension`: %+v", err)
			}
		}
	}

	return nil
}

func flattenExtensions(inputList *[]pricings.Extension) []any {
	outputList := make([]any, 0)

	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		// only keep enabled extensions
		if !strings.EqualFold(string(input.IsEnabled), "true") {
			continue
		}

		output := map[string]any{
			"name": input.Name,
		}
		if input.AdditionalExtensionProperties != nil {
			output["additional_extension_properties"] = *input.AdditionalExtensionProperties
		}

		outputList = append(outputList, output)
	}

	return outputList
}
