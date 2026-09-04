// Copyright IBM Corp. 2014, 2025
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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2023-01-01/pricings"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSecurityCenterSubscriptionPricing() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSecurityCenterSubscriptionPricingCreate,
		Read:   resourceSecurityCenterSubscriptionPricingRead,
		Update: resourceSecurityCenterSubscriptionPricingUpdate,
		Delete: resourceSecurityCenterSubscriptionPricingDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := pricings.ParsePricingID(id)
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
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInEnumSlice(security.PossiblePricingTierValues(), false),
			},

			"resource_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "VirtualMachines",
				ValidateFunc: validation.StringInSlice([]string{
					"AI",
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
				ForceNew: true,
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

func resourceSecurityCenterSubscriptionPricingCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.PricingClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id := pricings.NewPricingID(subscriptionId, d.Get("resource_type").(string))

	// Lock on subscription ID to prevent concurrent pricing updates across different resource types,
	// as the API only allows one pricing update per subscription at a time.
	locks.ByID(commonids.NewSubscriptionID(id.SubscriptionId).ID())
	defer locks.UnlockByID(commonids.NewSubscriptionID(id.SubscriptionId).ID())

	pricing := pricings.Pricing{
		Properties: &pricings.PricingProperties{
			PricingTier: pricings.PricingTier(d.Get("tier").(string)),
		},
	}

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		apiResponse, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(apiResponse.HttpResponse) {
				return fmt.Errorf("checking for presence of apiResponse %s: %+v", id, err)
			}
		}

		if err == nil && apiResponse.Model != nil && apiResponse.Model.Properties != nil && apiResponse.Model.Properties.PricingTier != pricings.PricingTierFree {
			return fmt.Errorf("the pricing tier of this subscription is not Free \r %+v", tf.ImportAsExistsError("azurerm_security_center_subscription_pricing", id.ID()))
		}
	}

	apiResponse, err := client.Get(ctx, id)
	if err != nil && !response.WasNotFound(apiResponse.HttpResponse) {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	extensionsStatusFromBackend := make([]pricings.Extension, 0)
	if err == nil && apiResponse.Model != nil && apiResponse.Model.Properties != nil {
		if apiResponse.Model.Properties.Extensions != nil {
			extensionsStatusFromBackend = *apiResponse.Model.Properties.Extensions
		}
	}

	if vSub, okSub := d.GetOk("subplan"); okSub {
		pricing.Properties.SubPlan = pointer.To(vSub.(string))
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

	// can not set any extension for free tier in the same request.
	if pricing.Properties.PricingTier == pricings.PricingTierStandard {
		pricing.Properties.Extensions = expandSecurityCenterSubscriptionPricingExtensions(realCfgExtensions, &extensionsStatusFromBackend)
	}

	if len(realCfgExtensions) > 0 && pricing.Properties.PricingTier == pricings.PricingTierFree {
		return fmt.Errorf("extensions cannot be enabled when using free tier")
	}

	pollerType := custompollers.NewSubscriptionPricingUpdatePoller(client, id, pricing)
	poller := pollers.NewPoller(pollerType, 5*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("setting %s: %+v", id, err)
	}

	// the extensions from backend might vary after pricing tier changed.
	updateResponse := pollerType.Response
	if updateResponse.Model != nil && updateResponse.Model.Properties != nil && updateResponse.Model.Properties.Extensions != nil {
		extensionsStatusFromBackend = *updateResponse.Model.Properties.Extensions
	}

	pricing.Properties.Extensions = expandSecurityCenterSubscriptionPricingExtensions(realCfgExtensions, &extensionsStatusFromBackend)

	poller = pollers.NewPoller(custompollers.NewSubscriptionPricingUpdatePoller(client, id, pricing), 5*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceSecurityCenterSubscriptionPricingRead(d, meta)
}

func resourceSecurityCenterSubscriptionPricingUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.PricingClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := pricings.ParsePricingID(d.Id())
	if err != nil {
		return err
	}

	// Lock on subscription ID to prevent concurrent pricing updates across different resource types,
	// as the API only allows one pricing update per subscription at a time.
	locks.ByID(commonids.NewSubscriptionID(id.SubscriptionId).ID())
	defer locks.UnlockByID(commonids.NewSubscriptionID(id.SubscriptionId).ID())

	apiResponse, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	update := pricings.Pricing{
		Properties: &pricings.PricingProperties{
			PricingTier: pricings.PricingTier(d.Get("tier").(string)),
		},
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

	if len(realCfgExtensions) > 0 && update.Properties.PricingTier == pricings.PricingTierFree {
		return fmt.Errorf("extensions cannot be enabled when using free tier")
	}

	extensionsStatusFromBackend := make([]pricings.Extension, 0)
	currentlyFreeTier := false
	if apiResponse.Model != nil && apiResponse.Model.Properties != nil {
		if apiResponse.Model.Properties.Extensions != nil {
			extensionsStatusFromBackend = *apiResponse.Model.Properties.Extensions
		}

		currentlyFreeTier = apiResponse.Model.Properties.PricingTier == pricings.PricingTierFree
	}

	// Update from `free` tier to `Standard`, we need to update it to `standard` tier first without extensions
	// Then do an additional update for the `extensions`
	// The service will enable extensions in the first update
	requiredAdditionalUpdate := false
	if update.Properties.PricingTier == pricings.PricingTierStandard {
		update.Properties.Extensions = expandSecurityCenterSubscriptionPricingExtensions(realCfgExtensions, &extensionsStatusFromBackend)
		requiredAdditionalUpdate = currentlyFreeTier
	}

	pollerType := custompollers.NewSubscriptionPricingUpdatePoller(client, *id, update)
	poller := pollers.NewPoller(pollerType, 5*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("setting %s: %+v", id, err)
	}

	// The extensions list from backend might vary after `tier` changed, thus we need to retrieve it again.
	updateResponse := pollerType.Response
	if updateResponse.Model != nil && updateResponse.Model.Properties != nil {
		if updateResponse.Model.Properties.Extensions != nil {
			extensionsStatusFromBackend = *updateResponse.Model.Properties.Extensions
		}
	}

	if requiredAdditionalUpdate {
		update.Properties.Extensions = expandSecurityCenterSubscriptionPricingExtensions(realCfgExtensions, &extensionsStatusFromBackend)
		poller = pollers.NewPoller(custompollers.NewSubscriptionPricingUpdatePoller(client, *id, update), 5*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
		if err := poller.PollUntilDone(ctx); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())
	return resourceSecurityCenterSubscriptionPricingRead(d, meta)
}

func resourceSecurityCenterSubscriptionPricingRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

func resourceSecurityCenterSubscriptionPricingDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.PricingClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := pricings.ParsePricingID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing %s: %+v", d.Id(), err)
	}

	// Lock on subscription ID to prevent concurrent pricing updates across different resource types,
	// as the API only allows one pricing update per subscription at a time.
	locks.ByID(commonids.NewSubscriptionID(id.SubscriptionId).ID())
	defer locks.UnlockByID(commonids.NewSubscriptionID(id.SubscriptionId).ID())

	pricing := pricings.Pricing{
		Properties: &pricings.PricingProperties{
			PricingTier: pricings.PricingTierFree,
		},
	}

	poller := pollers.NewPoller(custompollers.NewSubscriptionPricingUpdatePoller(client, *id, pricing), 5*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("setting %s: %+v", id, err)
	}

	log.Printf("[DEBUG] Security Center Subscription deletion invocation")
	return nil
}

func expandSecurityCenterSubscriptionPricingExtensions(inputList []interface{}, extensionsStatusFromBackend *[]pricings.Extension) *[]pricings.Extension {
	extensionStatuses := map[string]bool{}
	extensionProperties := make(map[string]interface{})

	outputList := make([]pricings.Extension, 0, len(inputList))
	if extensionsStatusFromBackend != nil {
		for _, backendExtension := range *extensionsStatusFromBackend {
			// set the default value to false, then turn on the extension that appear in the template
			extensionStatuses[backendExtension.Name] = false
			if backendExtension.AdditionalExtensionProperties != nil {
				extensionProperties[backendExtension.Name] = *backendExtension.AdditionalExtensionProperties
			}
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
		isEnabled := pricings.IsEnabledFalse
		if toBeEnabled {
			isEnabled = pricings.IsEnabledTrue
		}
		output := pricings.Extension{
			Name:      extensionName,
			IsEnabled: isEnabled,
		}

		// The service will return HTTP 500 if the payload contains extensionProperties and `IsEnabled==false`
		// `AdditionalProperties of Extension 'xxx' can't be updated while the extension is disabled (IsEnabled = False)`
		if vAdditional, ok := extensionProperties[extensionName]; ok && toBeEnabled {
			props, _ := vAdditional.(*interface{})
			p := (*props).(map[string]interface{})
			output.AdditionalExtensionProperties = pointer.To(p)
		}

		outputList = append(outputList, output)
	}

	return &outputList
}

func flattenExtensions(inputList *[]pricings.Extension) []interface{} {
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
