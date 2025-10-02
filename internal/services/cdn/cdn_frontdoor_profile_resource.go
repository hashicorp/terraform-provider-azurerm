// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorProfile() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontDoorProfileCreate,
		Read:   resourceCdnFrontDoorProfileRead,
		Update: resourceCdnFrontDoorProfileUpdate,
		Delete: resourceCdnFrontDoorProfileDelete,

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
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontDoorName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"response_timeout_seconds": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      120,
				ValidateFunc: validation.IntBetween(16, 240),
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(profiles.SkuNamePremiumAzureFrontDoor),
					string(profiles.SkuNameStandardAzureFrontDoor),
				}, false),
			},

			"log_scrubbing_rule": {
				Type:     pluginsdk.TypeSet,
				MaxItems: 3,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"match_variable": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice(
								profiles.PossibleValuesForScrubbingRuleEntryMatchVariable(),
								false),
						},
					},
				},
			},

			"tags": commonschema.Tags(),

			"resource_guid": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			// Verify that they are not downgrading the service from Premium SKU -> Standard SKU...
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v interface{}) error {
				oSku, nSku := diff.GetChange("sku_name")

				if oSku != "" {
					if oSku.(string) == string(profiles.SkuNamePremiumAzureFrontDoor) && nSku.(string) == string(profiles.SkuNameStandardAzureFrontDoor) {
						return fmt.Errorf("downgrading `sku_name` from `%s` to `%s` is not supported", profiles.SkuNamePremiumAzureFrontDoor, profiles.SkuNameStandardAzureFrontDoor)
					}
				}

				return nil
			}),
		),
	}
}

func resourceCdnFrontDoorProfileCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorProfilesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := profiles.NewProfileID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_cdn_frontdoor_profile", id.ID())
	}

	props := profiles.Profile{
		Location: location.Normalize("global"),
		Properties: &profiles.ProfileProperties{
			OriginResponseTimeoutSeconds: pointer.To(int64(d.Get("response_timeout_seconds").(int))),
		},
		Sku: profiles.Sku{
			Name: pointer.To(profiles.SkuName(d.Get("sku_name").(string))),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	props.Properties.LogScrubbing = expandCdnFrontDoorProfileLogScrubbing(d.Get("log_scrubbing_rule").(*pluginsdk.Set).List())

	if v, ok := d.GetOk("identity"); ok {
		i, err := identity.ExpandSystemAndUserAssignedMap(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}

		props.Identity = i
	}

	err = client.CreateThenPoll(ctx, id, props)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCdnFrontDoorProfileRead(d, meta)
}

func resourceCdnFrontDoorProfileRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorProfilesClient
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
		if skuName := model.Sku.Name; skuName != nil {
			d.Set("sku_name", string(pointer.From(skuName)))
		}

		identity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}

		if err := d.Set("identity", identity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if props := model.Properties; props != nil {
			d.Set("response_timeout_seconds", int(pointer.From(props.OriginResponseTimeoutSeconds)))

			// whilst this is returned in the API as FrontDoorID other resources refer to
			// this as the Resource GUID, so we will for consistency
			d.Set("resource_guid", pointer.From(props.FrontDoorId))

			if err := d.Set("log_scrubbing_rule", flattenCdnFrontDoorProfileLogScrubbingRules(props.LogScrubbing)); err != nil {
				return fmt.Errorf("setting `log_scrubbing_rule`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceCdnFrontDoorProfileUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorProfilesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := profiles.ParseProfileID(d.Id())
	if err != nil {
		return err
	}

	props := profiles.ProfileUpdateParameters{
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
		Properties: &profiles.ProfilePropertiesUpdateParameters{},
	}

	if d.HasChange("response_timeout_seconds") {
		props.Properties.OriginResponseTimeoutSeconds = pointer.To(int64(d.Get("response_timeout_seconds").(int)))
	}

	if d.HasChange("log_scrubbing_rule") {
		props.Properties.LogScrubbing = expandCdnFrontDoorProfileLogScrubbing(d.Get("log_scrubbing_rule").(*pluginsdk.Set).List())
	}

	if d.HasChange("identity") {
		i, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}

		props.Identity = i
	}

	err = client.UpdateThenPoll(ctx, pointer.From(id), props)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceCdnFrontDoorProfileRead(d, meta)
}

func resourceCdnFrontDoorProfileDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorProfilesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := profiles.ParseProfileID(d.Id())
	if err != nil {
		return err
	}

	err = client.DeleteThenPoll(ctx, pointer.From(id))
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandCdnFrontDoorProfileLogScrubbing(input []interface{}) *profiles.ProfileLogScrubbing {
	if len(input) == 0 {
		return &profiles.ProfileLogScrubbing{
			State: pointer.To(profiles.ProfileScrubbingStateDisabled),
		}
	}

	return &profiles.ProfileLogScrubbing{
		State:          pointer.To(profiles.ProfileScrubbingStateEnabled),
		ScrubbingRules: expandCdnFrontDoorProfileLogScrubbingRules(input),
	}
}

func expandCdnFrontDoorProfileLogScrubbingRules(input []interface{}) *[]profiles.ProfileScrubbingRules {
	if len(input) == 0 {
		return nil
	}

	scrubbingRules := make([]profiles.ProfileScrubbingRules, 0)

	for _, rule := range input {
		v := rule.(map[string]interface{})

		item := profiles.ProfileScrubbingRules{
			MatchVariable:         profiles.ScrubbingRuleEntryMatchVariable(v["match_variable"].(string)),
			Selector:              nil,
			SelectorMatchOperator: profiles.ScrubbingRuleEntryMatchOperatorEqualsAny, // EqualsAny is the only valid SelectorMatchOperator for log scrubbing in the Profile API
			State:                 pointer.To(profiles.ScrubbingRuleEntryStateEnabled),
		}

		scrubbingRules = append(scrubbingRules, item)
	}

	return &scrubbingRules
}

func flattenCdnFrontDoorProfileLogScrubbingRules(input *profiles.ProfileLogScrubbing) []interface{} {
	result := make([]interface{}, 0)

	if input == nil || pointer.From(input.State) == profiles.ProfileScrubbingStateDisabled {
		return result
	}

	if input.ScrubbingRules == nil {
		return result
	}

	for _, scrubbingRule := range *input.ScrubbingRules {
		if scrubbingRule.State != nil && pointer.From(scrubbingRule.State) == profiles.ScrubbingRuleEntryStateEnabled {
			result = append(result, map[string]interface{}{
				"match_variable": scrubbingRule.MatchVariable,
			})
		}
	}

	return result
}
