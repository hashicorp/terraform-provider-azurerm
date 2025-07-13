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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

			"log_scrubbing": {
				Type:     pluginsdk.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"scrubbing_rule": {
							Type:     pluginsdk.TypeList,
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
			// Validate log scrubbing configuration
			func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
				if logScrubbingRaw, ok := diff.GetOk("log_scrubbing"); ok {
					logScrubbing := logScrubbingRaw.([]interface{})
					if len(logScrubbing) > 0 && logScrubbing[0] == nil {
						return fmt.Errorf("when the `log_scrubbing` block is defined, at least one `scrubbing_rule` must be specified")
					}

					// Validate scrubbing rules for duplicates
					if len(logScrubbing) > 0 {
						if config, ok := logScrubbing[0].(map[string]interface{}); ok {
							if rulesRaw, ok := config["scrubbing_rule"].([]interface{}); ok {
								matchVariables := make(map[string]bool)
								for i, rule := range rulesRaw {
									if ruleMap, ok := rule.(map[string]interface{}); ok {
										if matchVar, ok := ruleMap["match_variable"].(string); ok && matchVar != "" {
											if matchVariables[matchVar] {
												return fmt.Errorf("duplicate `%s` rule found in `log_scrubbing.0.scrubbing_rule.%d.match_variable`", matchVar, i)
											}
											matchVariables[matchVar] = true
										}
									}
								}
							}
						}
					}
				}

				return nil
			},
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

	// Always set log scrubbing - if block is omitted, it will be set to disabled
	logScrubbing := expandCdnFrontDoorProfileLogScrubbing(d.Get("log_scrubbing").([]interface{}))
	props.Properties.LogScrubbing = logScrubbing

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

			if err := d.Set("log_scrubbing", flattenCdnFrontDoorProfileLogScrubbing(props.LogScrubbing)); err != nil {
				return fmt.Errorf("setting `log_scrubbing`: %+v", err)
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

	if d.HasChange("log_scrubbing") {
		// Always set log scrubbing - if block is omitted, it will be set to disabled
		logScrubbing := expandCdnFrontDoorProfileLogScrubbing(d.Get("log_scrubbing").([]interface{}))
		props.Properties.LogScrubbing = logScrubbing
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
		// When log_scrubbing block is not configured, set to disabled
		policyDisabled := profiles.ProfileScrubbingStateDisabled
		return &profiles.ProfileLogScrubbing{
			State:          &policyDisabled,
			ScrubbingRules: nil,
		}
	}

	inputRaw := input[0].(map[string]interface{})

	// When log_scrubbing block is present, always enable it
	policyEnabled := profiles.ProfileScrubbingStateEnabled
	scrubbingRules := expandCdnFrontDoorProfileScrubbingRules(inputRaw["scrubbing_rule"].([]interface{}))

	return &profiles.ProfileLogScrubbing{
		State:          &policyEnabled,
		ScrubbingRules: scrubbingRules,
	}
}

func expandCdnFrontDoorProfileScrubbingRules(input []interface{}) *[]profiles.ProfileScrubbingRules {
	if len(input) == 0 {
		return nil
	}

	scrubbingRules := make([]profiles.ProfileScrubbingRules, 0)

	for _, rule := range input {
		v := rule.(map[string]interface{})
		var item profiles.ProfileScrubbingRules

		// Rule is enabled by default when present (following "None" pattern)
		ruleEnabled := profiles.ScrubbingRuleEntryStateEnabled
		item.State = &ruleEnabled
		item.MatchVariable = profiles.ScrubbingRuleEntryMatchVariable(v["match_variable"].(string))

		// EqualsAny is the only valid SelectorMatchOperator for log scrubbing in the Profile API
		// so we will hardcode it here
		item.SelectorMatchOperator = "EqualsAny"

		scrubbingRules = append(scrubbingRules, item)
	}

	return &scrubbingRules
}

func flattenCdnFrontDoorProfileLogScrubbing(input *profiles.ProfileLogScrubbing) []interface{} {
	if input == nil || pointer.From(input.State) == profiles.ProfileScrubbingStateDisabled {
		// When log scrubbing is disabled, return empty list to represent omitted configuration
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})
	result["scrubbing_rule"] = flattenCdnFrontDoorProfileScrubbingRules(input.ScrubbingRules)

	return []interface{}{result}
}

func flattenCdnFrontDoorProfileScrubbingRules(scrubbingRules *[]profiles.ProfileScrubbingRules) interface{} {
	result := make([]interface{}, 0)

	if scrubbingRules == nil || len(*scrubbingRules) == 0 {
		return result
	}

	for _, scrubbingRule := range *scrubbingRules {
		if scrubbingRule.State != nil && pointer.From(scrubbingRule.State) == profiles.ScrubbingRuleEntryStateEnabled {
			item := map[string]interface{}{}
			item["match_variable"] = scrubbingRule.MatchVariable
			result = append(result, item)
		}
	}

	return result
}
