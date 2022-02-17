package cdn

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	tagsHelper "github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"

	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceFrontdoorProfile() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFrontdoorProfileCreate,
		Read:   resourceFrontdoorProfileRead,
		Update: resourceFrontdoorProfileUpdate,
		Delete: resourceFrontdoorProfileDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := profiles.ParseProfileID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"frontdoor_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentity(),

			"kind": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"location": azure.SchemaLocation(),

			"origin_response_timeout_seconds": {
				Type:         pluginsdk.TypeInt,
				ForceNew:     true,
				Optional:     true,
				Default:      120,
				ValidateFunc: validation.IntBetween(16, 240),
			},

			"provisioning_state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"resource_state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
				ForceNew: true,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(profiles.SkuNameCustomVerizon),
					string(profiles.SkuNamePremiumAzureFrontDoor),
					string(profiles.SkuNamePremiumVerizon),
					string(profiles.SkuNameStandardAkamai),
					string(profiles.SkuNameStandardAvgBandWidthChinaCdn),
					string(profiles.SkuNameStandardAzureFrontDoor),
					string(profiles.SkuNameStandardChinaCdn),
					string(profiles.SkuNameStandardMicrosoft),
					string(profiles.SkuNameStandardNineFiveFiveBandWidthChinaCdn),
					string(profiles.SkuNameStandardPlusAvgBandWidthChinaCdn),
					string(profiles.SkuNameStandardPlusChinaCdn),
					string(profiles.SkuNameStandardPlusNineFiveFiveBandWidthChinaCdn),
					string(profiles.SkuNameStandardVerizon),
				}, false),
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceFrontdoorProfileCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorProfileClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	sdkId := profiles.NewProfileID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	id := parse.NewFrontdoorProfileID(subscriptionId, sdkId.ResourceGroupName, sdkId.ProfileName)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, sdkId)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_frontdoor_profile", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location"))

	identity, err := expandSystemAndUserAssignedIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	props := profiles.Profile{
		Location: location,
		Properties: &profiles.ProfileProperties{
			Identity:                     identity,
			OriginResponseTimeoutSeconds: utils.Int64(int64(d.Get("origin_response_timeout_seconds").(int))),
		},
		Sku:  *expandProfileSku(d.Get("sku_name").(string)),
		Tags: tagsHelper.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateThenPoll(ctx, sdkId, props); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceFrontdoorProfileRead(d, meta)
}

func resourceFrontdoorProfileRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorProfileClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := profiles.ParseProfileID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ProfileName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		if props := model.Properties; props != nil {
			d.Set("frontdoor_id", props.FrontDoorId)
			d.Set("kind", model.Kind)
			d.Set("identity", flattenSystemAndUserAssignedIdentity(props.Identity))
			d.Set("origin_response_timeout_seconds", props.OriginResponseTimeoutSeconds)
			d.Set("provisioning_state", props.ProvisioningState)
			d.Set("resource_state", props.ResourceState)
		}

		if err := d.Set("sku_name", flattenProfileSku(&model.Sku)); err != nil {
			return fmt.Errorf("setting `sku_name`: %+v", err)
		}

		if err := tags.FlattenAndSet(d, ConvertFrontdoorProfileTags(model.Tags)); err != nil {
			return err
		}
	}
	return nil
}

func resourceFrontdoorProfileUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorProfileClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := profiles.ParseProfileID(d.Id())
	if err != nil {
		return err
	}

	props := profiles.ProfileUpdateParameters{
		Tags: tagsHelper.Expand(d.Get("tags").(map[string]interface{})),
	}
	if err := client.UpdateThenPoll(ctx, *id, props); err != nil {

		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceFrontdoorProfileRead(d, meta)
}

func resourceFrontdoorProfileDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorProfileClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := profiles.ParseProfileID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {

		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	return nil
}

func expandProfileSku(input string) *profiles.Sku {
	if len(input) == 0 || input == "" {
		return nil
	}

	nameValue := profiles.SkuName(input)
	return &profiles.Sku{
		Name: &nameValue,
	}
}

func expandSystemAndUserAssignedIdentity(input []interface{}) (*identity.SystemAndUserAssignedMap, error) {
	expanded, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	out := identity.SystemAndUserAssignedMap{
		Type: expanded.Type,
	}

	if len(expanded.IdentityIds) > 0 {
		for k := range expanded.IdentityIds {
			out.IdentityIds[k] = identity.UserAssignedIdentityDetails{
				// intentionally empty
			}
		}
	}

	return &out, nil
}

func flattenProfileSku(input *profiles.Sku) string {
	result := ""
	if input == nil {
		return result
	}

	if input.Name != nil {
		result = string(*input.Name)
	}

	return result
}

func flattenSystemAndUserAssignedIdentity(input *identity.SystemAndUserAssignedMap) []interface{} {
	var transform *identity.SystemAndUserAssignedMap

	if input != nil {
		transform = &identity.SystemAndUserAssignedMap{
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
			Type:        identity.Type(string(input.Type)),
		}

		if input.PrincipalId != "" {
			transform.PrincipalId = input.PrincipalId
		}

		if input.TenantId != "" {
			transform.TenantId = input.TenantId
		}

		for k, v := range input.IdentityIds {
			transform.IdentityIds[k] = identity.UserAssignedIdentityDetails{
				ClientId:    v.ClientId,
				PrincipalId: v.PrincipalId,
			}
		}
	}

	out, err := identity.FlattenSystemAndUserAssignedMap(transform)
	if err != nil {
		return []interface{}{}
	}

	return *out
}
