// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package securitycenter

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	pricings_v2023_01_01 "github.com/hashicorp/go-azure-sdk/resource-manager/security/2023-01-01/pricings"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSecurityCenterSubscriptionPricing() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSecurityCenterSubscriptionPricingUpdate,
		Read:   resourceSecurityCenterSubscriptionPricingRead,
		Update: resourceSecurityCenterSubscriptionPricingUpdate,
		Delete: resourceSecurityCenterSubscriptionPricingDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.PricingID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SubscriptionPricingV0ToV1{},
		}),

		Schema: map[string]*pluginsdk.Schema{
			"tier": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(security.PricingTierFree),
					string(security.PricingTierStandard),
				}, false),
			},
			"resource_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "VirtualMachines",
				ValidateFunc: validation.StringInSlice([]string{
					"Api",
					"AppServices",
					"ContainerRegistry",
					"KeyVaults",
					"KubernetesService",
					"SqlServers",
					"SqlServerVirtualMachines",
					"StorageAccounts",
					"VirtualMachines",
					"Arm",
					"Dns",
					"OpenSourceRelationalDatabases",
					"Containers",
					"CosmosDbs",
					"CloudPosture",
				}, false),
			},
			"subplan": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
			"extension": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotWhiteSpace,
						},
						"additional_extension_properties": {
							Type:     pluginsdk.TypeMap,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotWhiteSpace,
							},
						},
					},
				},
			},
		},
	}
}

func resourceSecurityCenterSubscriptionPricingUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.PricingClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := pricings_v2023_01_01.NewPricingID(subscriptionId, d.Get("resource_type").(string))
	pricing := pricings_v2023_01_01.Pricing{
		Properties: &pricings_v2023_01_01.PricingProperties{
			PricingTier: pricings_v2023_01_01.PricingTier(d.Get("tier").(string)),
		},
	}

	apiResponse, err := client.Get(ctx, id)
	if d.IsNewResource() {
		if err != nil {
			if !response.WasNotFound(apiResponse.HttpResponse) {
				return fmt.Errorf("checking for presence of apiResponse %s: %+v", id, err)
			}
		}

		if err == nil && apiResponse.Model != nil && apiResponse.Model.Properties != nil && apiResponse.Model.Properties.PricingTier != pricings_v2023_01_01.PricingTierFree {
			return fmt.Errorf("the pricing tier of this subscription is not Free \r %+v", tf.ImportAsExistsError("azurerm_security_center_subscription_pricing", id.ID()))
		}
	}

	extensionsStatusFromBackend := make([]pricings_v2023_01_01.Extension, 0)
	isCurrentlyInFree := false
	if err == nil && apiResponse.Model != nil && apiResponse.Model.Properties != nil {
		if apiResponse.Model.Properties.Extensions != nil {
			extensionsStatusFromBackend = *apiResponse.Model.Properties.Extensions
		}

		if apiResponse.Model.Properties.PricingTier == pricings_v2023_01_01.PricingTierFree {
			isCurrentlyInFree = true
		}
	}

	if vSub, okSub := d.GetOk("subplan"); okSub {
		pricing.Properties.SubPlan = utils.String(vSub.(string))
	}

	// When the state file contains an `extension` with `additional_extension_properties`
	// But the tf config does not, `d.Get("extension")` will contain a zero element.
	// Tracked by https://github.com/hashicorp/terraform-plugin-sdk/issues/1248
	realCfgExtensions := make([]interface{}, 0)
	for _, e := range d.Get("extension").(*pluginsdk.Set).List() {
		v := e.(map[string]interface{})
		if v["name"] != "" {
			realCfgExtensions = append(realCfgExtensions, e)
		}
	}

	if d.HasChange("extension") {
		// can not set extensions for free tier
		if pricing.Properties.PricingTier == pricings_v2023_01_01.PricingTierStandard {
			extensions := expandSecurityCenterSubscriptionPricingExtensions(realCfgExtensions, &extensionsStatusFromBackend)
			pricing.Properties.Extensions = extensions
		}
	}

	if len(realCfgExtensions) > 0 && pricing.Properties.PricingTier == pricings_v2023_01_01.PricingTierFree {
		return fmt.Errorf("extensions cannot be enabled when using free tier")
	}

	updateResponse, updateErr := client.Update(ctx, id, pricing)
	if updateErr != nil {
		return fmt.Errorf("setting %s: %+v", id, updateErr)
	}

	if updateErr == nil && updateResponse.Model != nil && updateResponse.Model.Properties != nil {
		if updateResponse.Model.Properties.Extensions != nil {
			extensionsStatusFromBackend = *updateResponse.Model.Properties.Extensions
		}
	}

	// after turning on the bundle, we have now the extensions list
	if d.IsNewResource() || isCurrentlyInFree {
		extensions := expandSecurityCenterSubscriptionPricingExtensions(realCfgExtensions, &extensionsStatusFromBackend)
		pricing.Properties.Extensions = extensions
		_, updateErr := client.Update(ctx, id, pricing)
		if err != nil {
			if updateErr != nil {
				return fmt.Errorf("setting %s: %+v", id, updateErr)
			}
		}
	}

	d.SetId(id.ID())
	return resourceSecurityCenterSubscriptionPricingRead(d, meta)
}

func resourceSecurityCenterSubscriptionPricingRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.PricingClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := pricings_v2023_01_01.ParsePricingID(d.Id())
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
			err = d.Set("extension", flattenExtensions(properties.Extensions))
			if err != nil {
				return fmt.Errorf("setting `extension`: %+v", err)
			}
		}
	}

	return nil
}

func resourceSecurityCenterSubscriptionPricingDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.PricingClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := pricings_v2023_01_01.ParsePricingID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing %s: %+v", d.Id(), err)
	}

	pricing := pricings_v2023_01_01.Pricing{
		Properties: &pricings_v2023_01_01.PricingProperties{
			PricingTier: pricings_v2023_01_01.PricingTierFree,
		},
	}

	if _, err := client.Update(ctx, *id, pricing); err != nil {
		return fmt.Errorf("setting %s: %+v", id, err)
	}

	log.Printf("[DEBUG] Security Center Subscription deletion invocation")
	return nil
}

func expandSecurityCenterSubscriptionPricingExtensions(inputList []interface{}, extensionsStatusFromBackend *[]pricings_v2023_01_01.Extension) *[]pricings_v2023_01_01.Extension {
	extensionStatuses := map[string]bool{}
	extensionProperties := make(map[string]interface{})
	var outputList []pricings_v2023_01_01.Extension

	if extensionsStatusFromBackend != nil {
		for _, backendExtension := range *extensionsStatusFromBackend {
			// set the default value to false, then turn on the extension that appear in the template
			extensionStatuses[backendExtension.Name] = false
		}
	}

	// set any extension in the template to be true
	for _, v := range inputList {
		input := v.(map[string]interface{})
		if input["name"] == "" {
			continue
		}
		extensionStatuses[input["name"].(string)] = true
		if vAdditional, ok := input["additional_extension_properties"]; ok {
			extensionProperties[input["name"].(string)] = &vAdditional
		}
	}

	for extensionName, toBeEnabled := range extensionStatuses {

		isEnabled := pricings_v2023_01_01.IsEnabledFalse
		if toBeEnabled {
			isEnabled = pricings_v2023_01_01.IsEnabledTrue
		}
		output := pricings_v2023_01_01.Extension{
			Name:                          extensionName,
			IsEnabled:                     isEnabled,
			AdditionalExtensionProperties: pointer.To(extensionProperties),
		}
		outputList = append(outputList, output)
	}

	return &outputList
}

func flattenExtensions(inputList *[]pricings_v2023_01_01.Extension) []interface{} {
	outputList := make([]interface{}, 0)

	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		// only keep enabled extensions
		if !strings.EqualFold(string(input.IsEnabled), "true") {
			continue
		}

		output := map[string]interface{}{
			"name": input.Name,
		}
		if input.AdditionalExtensionProperties != nil {
			output["additional_extension_properties"] = *input.AdditionalExtensionProperties
		}

		outputList = append(outputList, output)
	}

	return outputList
}
