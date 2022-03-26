package cdn

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	track1 "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCdnFrontdoorProfile() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontdoorProfileCreate,
		Read:   resourceCdnFrontdoorProfileRead,
		Update: resourceCdnFrontdoorProfileUpdate,
		Delete: resourceCdnFrontdoorProfileDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FrontdoorProfileID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"cdn_frontdoor_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"response_timeout_seconds": {
				Type:         pluginsdk.TypeInt,
				ForceNew:     true,
				Optional:     true,
				Default:      120,
				ValidateFunc: validation.IntBetween(16, 240),
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
				ForceNew: true,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(track1.SkuNameCustomVerizon),
					string(track1.SkuNamePremiumAzureFrontDoor),
					string(track1.SkuNamePremiumVerizon),
					string(track1.SkuNameStandardAkamai),
					string(track1.SkuNameStandardAvgBandWidthChinaCdn),
					string(track1.SkuNameStandardAzureFrontDoor),
					string(track1.SkuNameStandardChinaCdn),
					string(track1.SkuNameStandardMicrosoft),
					string(track1.SkuNameStandard955BandWidthChinaCdn),
					string(track1.SkuNameStandardPlusAvgBandWidthChinaCdn),
					string(track1.SkuNameStandardPlusChinaCdn),
					string(track1.SkuNameStandardPlus955BandWidthChinaCdn),
					string(track1.SkuNameStandardVerizon),
				}, false),
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceCdnFrontdoorProfileCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorProfileClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewFrontdoorProfileID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_cdn_frontdoor_profile", id.ID())
		}
	}

	// Can only be Global for all Frontdoors
	location := azure.NormalizeLocation("global")

	props := track1.Profile{
		Location: utils.String(location),
		ProfileProperties: &track1.ProfileProperties{
			OriginResponseTimeoutSeconds: utils.Int32(int32(d.Get("response_timeout_seconds").(int))),
		},
		Sku:  expandProfileSku(d.Get("sku_name").(string)),
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.ProfileName, props)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCdnFrontdoorProfileRead(d, meta)
}

func resourceCdnFrontdoorProfileRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorProfileClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontdoorProfileID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ProfileName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.ProfileProperties; props != nil {
		d.Set("cdn_frontdoor_id", props.FrontDoorID)
		d.Set("response_timeout_seconds", props.OriginResponseTimeoutSeconds)
	}

	if err := d.Set("sku_name", flattenProfileSku(resp.Sku)); err != nil {
		return fmt.Errorf("setting `sku_name`: %+v", err)
	}

	if err := tags.FlattenAndSet(d, resp.Tags); err != nil {
		return err
	}

	return nil
}

func resourceCdnFrontdoorProfileUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorProfileClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontdoorProfileID(d.Id())
	if err != nil {
		return err
	}

	props := track1.ProfileUpdateParameters{
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.ProfileName, props)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the update of %s: %+v", *id, err)
	}

	return resourceCdnFrontdoorProfileRead(d, meta)
}

func resourceCdnFrontdoorProfileDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorProfileClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontdoorProfileID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ProfileName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return nil
}

func expandProfileSku(input string) *track1.Sku {
	if len(input) == 0 || input == "" {
		return nil
	}

	return &track1.Sku{
		Name: track1.SkuName(input),
	}
}

func flattenProfileSku(input *track1.Sku) string {
	result := ""
	if input == nil {
		return result
	}

	if input.Name != "" {
		result = string(input.Name)
	}

	return result
}
