// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package web

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAppServicePlan() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppServicePlanCreateUpdate,
		Read:   resourceAppServicePlanRead,
		Update: resourceAppServicePlanCreateUpdate,
		Delete: resourceAppServicePlanDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AppServicePlanID(id)
			return err
		}),

		DeprecationMessage: "The `azurerm_app_service_plan` resource has been superseded by the `azurerm_service_plan` resource. Whilst this resource will continue to be available in the 2.x and 3.x releases it is feature-frozen for compatibility purposes, will no longer receive any updates and will be removed in a future major release of the Azure Provider.",

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		SchemaVersion: 1,

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.AppServicePlanV0toV1{},
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AppServicePlanName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"kind": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "Windows",
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					// @tombuildsstuff: I believe `app` is the older representation of `Windows`
					// thus we need to support it to be able to import resources without recreating them.
					// @jcorioland: new SKU and kind 'xenon' have been added for Windows Containers support
					// https://azure.microsoft.com/en-us/blog/announcing-the-public-preview-of-windows-container-support-in-azure-app-service/
					"App",
					"elastic",
					"FunctionApp",
					"Linux",
					"Windows",
					"xenon",
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"sku": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"tier": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"size": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"capacity": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			// / AppServicePlanProperties
			"app_service_environment_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"per_site_scaling": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"reserved": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"maximum_elastic_worker_count": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},

			"maximum_number_of_workers": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"is_xenon": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"zone_redundant": {
				Type:     pluginsdk.TypeBool,
				ForceNew: true,
				Optional: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceAppServicePlanCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicePlansClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM App Service Plan creation.")

	id := parse.NewAppServicePlanID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ServerFarmName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_app_service_plan", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	kind := d.Get("kind").(string)
	t := d.Get("tags").(map[string]interface{})

	sku := expandAppServicePlanSku(d)
	properties := &web.AppServicePlanProperties{}

	isXenon := d.Get("is_xenon").(bool)
	properties.IsXenon = &isXenon

	if kind == "xenon" && !isXenon {
		return fmt.Errorf("creating or updating %s: when kind is set to xenon, is_xenon property should be set to true", id)
	}

	appServicePlan := web.AppServicePlan{
		Location:                 &location,
		Kind:                     &kind,
		Sku:                      &sku,
		Tags:                     tags.Expand(t),
		AppServicePlanProperties: properties,
	}

	if v := d.Get("app_service_environment_id").(string); v != "" {
		appServicePlan.AppServicePlanProperties.HostingEnvironmentProfile = &web.HostingEnvironmentProfile{
			ID: utils.String(v),
		}
	}

	if v := d.Get("per_site_scaling").(bool); v {
		appServicePlan.AppServicePlanProperties.PerSiteScaling = utils.Bool(v)
	}

	reserved := d.Get("reserved").(bool)
	if strings.EqualFold(kind, "Linux") && !reserved {
		return fmt.Errorf("`reserved` has to be set to true when kind is set to `Linux`")
	}

	if strings.EqualFold(kind, "Windows") && reserved {
		return fmt.Errorf("`reserved` has to be set to false when kind is set to `Windows`")
	}

	if v := d.Get("maximum_elastic_worker_count").(int); v > 0 {
		appServicePlan.AppServicePlanProperties.MaximumElasticWorkerCount = utils.Int32(int32(v))
	}

	if v := d.Get("zone_redundant").(bool); v {
		appServicePlan.AppServicePlanProperties.ZoneRedundant = utils.Bool(v)
	}

	if reserved {
		appServicePlan.AppServicePlanProperties.Reserved = utils.Bool(reserved)
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerFarmName, appServicePlan)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the create/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceAppServicePlanRead(d, meta)
}

func resourceAppServicePlanRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicePlansClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServicePlanID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServerFarmName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Service Plan %q was not found in Resource Group %q - removnig from state!", id.ServerFarmName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on App Service Plan %q (Resource Group %q): %+v", id.ServerFarmName, id.ResourceGroup, err)
	}

	// A 404 doesn't error from the app service plan sdk so we'll add this check here to catch resource not found responses
	// TODO This block can be removed if https://github.com/Azure/azure-sdk-for-go/issues/5407 gets addressed.
	if utils.ResponseWasNotFound(resp.Response) {
		log.Printf("[DEBUG] App Service Plan %q was not found in Resource Group %q - removing from state!", id.ServerFarmName, id.ResourceGroup)
		d.SetId("")
		return nil
	}

	d.Set("name", id.ServerFarmName)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	d.Set("kind", resp.Kind)

	if props := resp.AppServicePlanProperties; props != nil {
		appServiceEnvironmentId := ""
		if props.HostingEnvironmentProfile != nil && props.HostingEnvironmentProfile.ID != nil {
			appServiceEnvironmentId = *props.HostingEnvironmentProfile.ID
		}
		d.Set("app_service_environment_id", appServiceEnvironmentId)

		maximumNumberOfWorkers := 0
		if props.MaximumNumberOfWorkers != nil {
			maximumNumberOfWorkers = int(*props.MaximumNumberOfWorkers)
		}
		d.Set("maximum_number_of_workers", maximumNumberOfWorkers)

		maximumElasticWorkerCount := 0
		if props.MaximumElasticWorkerCount != nil {
			maximumElasticWorkerCount = int(*props.MaximumElasticWorkerCount)
		}
		d.Set("maximum_elastic_worker_count", maximumElasticWorkerCount)

		d.Set("per_site_scaling", props.PerSiteScaling)
		d.Set("reserved", props.Reserved)
		d.Set("is_xenon", props.IsXenon)
		d.Set("zone_redundant", props.ZoneRedundant)
	}

	if err := d.Set("sku", flattenAppServicePlanSku(resp.Sku)); err != nil {
		return fmt.Errorf("setting `sku`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceAppServicePlanDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicePlansClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServicePlanID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting App Service Plan %q (Resource Group %q)", id.ServerFarmName, id.ResourceGroup)

	resp, err := client.Delete(ctx, id.ResourceGroup, id.ServerFarmName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting App Service Plan %q (Resource Group %q): %+v", id.ServerFarmName, id.ResourceGroup, err)
		}
	}

	return nil
}

func expandAppServicePlanSku(d *pluginsdk.ResourceData) web.SkuDescription {
	configs := d.Get("sku").([]interface{})
	config := configs[0].(map[string]interface{})

	tier := config["tier"].(string)
	size := config["size"].(string)

	sku := web.SkuDescription{
		Name: utils.String(size),
		Tier: utils.String(tier),
		Size: utils.String(size),
	}

	if v, ok := config["capacity"]; ok {
		capacity := v.(int)
		sku.Capacity = utils.Int32(int32(capacity))
	}

	return sku
}

func flattenAppServicePlanSku(input *web.SkuDescription) []interface{} {
	outputs := make([]interface{}, 0)
	if input == nil {
		return outputs
	}

	output := make(map[string]interface{}, 2)

	if input.Tier != nil {
		output["tier"] = *input.Tier
	}

	if input.Size != nil {
		output["size"] = *input.Size
	}

	if input.Capacity != nil {
		output["capacity"] = *input.Capacity
	}

	outputs = append(outputs, output)

	return outputs
}
