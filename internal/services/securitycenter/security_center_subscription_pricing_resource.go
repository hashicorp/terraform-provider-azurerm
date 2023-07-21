// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package securitycenter

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security" // nolint: staticcheck
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
						"name:": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"AgentlessDiscoveryForKubernetes",
								"OnUploadMalwareScanning",
								"SensitiveDataDiscovery",
								"ContainerRegistriesVulnerabilityAssessments",
							}, false),
						},
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},
						"additional_properties": {
							Type:     pluginsdk.TypeMap,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
						"operation_status": {
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"code": {
										Computed: true,
										Type:     pluginsdk.TypeString,
									},
									"message": {
										Computed: true,
										Type:     pluginsdk.TypeString,
									},
								},
							},
						},
					},
				},
				Set: func(i interface{}) int {
					m := i.(map[string]interface{})
					return pluginsdk.HashString(m["name"])
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

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if existing.Model != nil && existing.Model.Properties != nil && existing.Model.Properties.PricingTier != pricings_v2023_01_01.PricingTierFree {
			return fmt.Errorf("the pricing tier of this subscription is not Free \r %+v", tf.ImportAsExistsError("azurerm_security_center_subscription_pricing", id.ID()))
		}
	}

	if v, ok := d.GetOk("subplan"); ok {
		pricing.Properties.SubPlan = utils.String(v.(string))
	}

	if v, ok := d.GetOk("extension"); ok {
		pricing.Properties.Extensions = expandPricingExtensions(v.(*pluginsdk.Set).List())
	}

	if _, err := client.Update(ctx, id, pricing); err != nil {
		return fmt.Errorf("setting %s: %+v", id, err)
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
			if properties.Extensions != nil {
				err = d.Set("extension", flattenPricingExtensions(properties.Extensions))
				if err != nil {
					return fmt.Errorf("setting `extension`: %+v", err)
				}
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

func expandPricingExtensions(extensionInput []interface{}) *[]pricings_v2023_01_01.Extension {
	if len(extensionInput) <= 0 {
		return nil
	}

	extensions := make([]pricings_v2023_01_01.Extension, len(extensionInput))

	for i, v := range extensionInput {
		raw := v.(map[string]interface{})
		extensions[i].Name = raw["name"].(string)
		extensions[i].IsEnabled = pricings_v2023_01_01.IsEnabledTrue
		if !raw["enabled"].(bool) {
			extensions[i].IsEnabled = pricings_v2023_01_01.IsEnabledFalse
		}
		if vap, ok := raw["additional_properties"]; ok {
			extensions[i].AdditionalExtensionProperties = &vap
		}
	}

	return &extensions
}

func flattenPricingExtensions(extensions *[]pricings_v2023_01_01.Extension) []interface{} {
	if extensions == nil || len(*extensions) <= 0 {
		return []interface{}{}
	}

	data := make([]interface{}, len(*extensions))
	for i, v := range *extensions {
		item := make(map[string]interface{})
		item["name"] = v.Name
		item["enabled"] = v.IsEnabled == pricings_v2023_01_01.IsEnabledTrue
		if v.AdditionalExtensionProperties != nil {
			item["additional_properties"] = v.AdditionalExtensionProperties
		}
		if v.OperationStatus != nil {
			status := make(map[string]interface{})
			status["code"] = v.OperationStatus.Code
			status["message"] = v.OperationStatus.Message
			item["operation_status"] = status
		}
		data[i] = item
	}

	return data
}
