// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package media

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/streamingpoliciesandstreaminglocators"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/media/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/media/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMediaStreamingLocator() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMediaStreamingLocatorCreate,
		Read:   resourceMediaStreamingLocatorRead,
		Delete: resourceMediaStreamingLocatorDelete,

		DeprecationMessage: azureMediaRetirementMessage,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := streamingpoliciesandstreaminglocators.ParseStreamingLocatorID(id)
			return err
		}),

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.StreamingLocatorV0ToV1{},
		}),
		SchemaVersion: 1,

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9(_)]{1,128}$"),
					"Streaming Locator name must be 1 - 128 characters long, can contain letters, numbers, underscores, and hyphens (but the first and last character must be a letter or number).",
				),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"media_services_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AccountName,
			},

			"asset_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9]{1,128}$"),
					"Asset name must be 1 - 128 characters long, contain only letters, hyphen and numbers.",
				),
			},

			"streaming_policy_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"alternative_media_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// lintignore:XS003
			"content_key": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"content_key_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.IsUUID,
						},

						"label_reference_in_streaming_policy": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"policy_name": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(streamingpoliciesandstreaminglocators.StreamingLocatorContentKeyTypeCommonEncryptionCbcs),
								string(streamingpoliciesandstreaminglocators.StreamingLocatorContentKeyTypeCommonEncryptionCenc),
								string(streamingpoliciesandstreaminglocators.StreamingLocatorContentKeyTypeEnvelopeEncryption),
							}, false),
						},

						"value": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"default_content_key_policy_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"end_time": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"filter_names": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"start_time": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"streaming_locator_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
		},
	}
}

func resourceMediaStreamingLocatorCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20220801Client.StreamingPoliciesAndStreamingLocators
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := streamingpoliciesandstreaminglocators.NewStreamingLocatorID(subscriptionID, d.Get("resource_group_name").(string), d.Get("media_services_account_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.StreamingLocatorsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_media_streaming_locator", id.ID())
		}
	}

	payload := streamingpoliciesandstreaminglocators.StreamingLocator{
		Properties: &streamingpoliciesandstreaminglocators.StreamingLocatorProperties{
			AssetName:           d.Get("asset_name").(string),
			StreamingPolicyName: d.Get("streaming_policy_name").(string),
		},
	}

	if alternativeMediaID, ok := d.GetOk("alternative_media_id"); ok {
		payload.Properties.AlternativeMediaId = utils.String(alternativeMediaID.(string))
	}

	if contentKeys, ok := d.GetOk("content_key"); ok {
		payload.Properties.ContentKeys = expandContentKeys(contentKeys.([]interface{}))
	}

	if defaultContentKeyPolicyName, ok := d.GetOk("default_content_key_policy_name"); ok {
		payload.Properties.DefaultContentKeyPolicyName = utils.String(defaultContentKeyPolicyName.(string))
	}

	if endTimeRaw, ok := d.GetOk("end_time"); ok {
		if endTimeRaw.(string) != "" {
			endTime, err := time.Parse(time.RFC3339, endTimeRaw.(string))
			if err != nil {
				return err
			}
			payload.Properties.SetEndTimeAsTime(endTime)
		}
	}

	if filters, ok := d.GetOk("filter_names"); ok {
		payload.Properties.Filters = utils.ExpandStringSlice(filters.([]interface{}))
	}

	if startTimeRaw, ok := d.GetOk("start_time"); ok {
		if startTimeRaw.(string) != "" {
			startTime, err := time.Parse(time.RFC3339, startTimeRaw.(string))
			if err != nil {
				return err
			}
			payload.Properties.SetStartTimeAsTime(startTime)
		}
	}

	if idRaw, ok := d.GetOk("streaming_locator_id"); ok {
		payload.Properties.StreamingLocatorId = pointer.To(idRaw.(string))
	}

	if _, err := client.StreamingLocatorsCreate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceMediaStreamingLocatorRead(d, meta)
}

func resourceMediaStreamingLocatorRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20220801Client.StreamingPoliciesAndStreamingLocators
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := streamingpoliciesandstreaminglocators.ParseStreamingLocatorID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.StreamingLocatorsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.StreamingLocatorName)
	d.Set("media_services_account_name", id.MediaServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("asset_name", props.AssetName)
			d.Set("streaming_policy_name", props.StreamingPolicyName)
			d.Set("alternative_media_id", props.AlternativeMediaId)
			d.Set("default_content_key_policy_name", props.DefaultContentKeyPolicyName)

			contentKeys := flattenContentKeys(props.ContentKeys)
			if err := d.Set("content_key", contentKeys); err != nil {
				return fmt.Errorf("flattening `content_key`: %s", err)
			}

			endTime := ""
			if props.EndTime != nil {
				t, err := props.GetEndTimeAsTime()
				if err != nil {
					return fmt.Errorf("parsing EndTime: %+v", err)
				}
				endTime = t.Format(time.RFC3339)
			}
			d.Set("end_time", endTime)
			d.Set("filter_names", utils.FlattenStringSlice(props.Filters))

			startTime := ""
			if props.StartTime != nil {
				t, err := props.GetStartTimeAsTime()
				if err != nil {
					return fmt.Errorf("parsing StartTime: %+v", err)
				}
				startTime = t.Format(time.RFC3339)
			}
			d.Set("start_time", startTime)

			d.Set("streaming_locator_id", props.StreamingLocatorId)
		}
	}

	return nil
}

func resourceMediaStreamingLocatorDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20220801Client.StreamingPoliciesAndStreamingLocators
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := streamingpoliciesandstreaminglocators.ParseStreamingLocatorID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.StreamingLocatorsDelete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandContentKeys(input []interface{}) *[]streamingpoliciesandstreaminglocators.StreamingLocatorContentKey {
	results := make([]streamingpoliciesandstreaminglocators.StreamingLocatorContentKey, 0)

	for _, contentKeyRaw := range input {
		if contentKeyRaw == nil {
			continue
		}
		contentKey := contentKeyRaw.(map[string]interface{})

		streamingLocatorContentKey := streamingpoliciesandstreaminglocators.StreamingLocatorContentKey{}

		if contentKey["content_key_id"] != nil {
			streamingLocatorContentKey.Id = contentKey["content_key_id"].(string)
		}

		if contentKey["label_reference_in_streaming_policy"] != nil {
			streamingLocatorContentKey.LabelReferenceInStreamingPolicy = utils.String(contentKey["label_reference_in_streaming_policy"].(string))
		}

		if contentKey["policy_name"] != nil {
			streamingLocatorContentKey.PolicyName = utils.String(contentKey["policy_name"].(string))
		}

		if contentKey["type"] != nil {
			streamingLocatorContentKey.Type = pointer.To(streamingpoliciesandstreaminglocators.StreamingLocatorContentKeyType(contentKey["type"].(string)))
		}

		if contentKey["value"] != nil {
			streamingLocatorContentKey.Value = utils.String(contentKey["value"].(string))
		}

		results = append(results, streamingLocatorContentKey)
	}

	return &results
}

func flattenContentKeys(input *[]streamingpoliciesandstreaminglocators.StreamingLocatorContentKey) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, contentKey := range *input {
		labelReferenceInStreamingPolicy := ""
		if contentKey.LabelReferenceInStreamingPolicy != nil {
			labelReferenceInStreamingPolicy = *contentKey.LabelReferenceInStreamingPolicy
		}

		policyName := ""
		if contentKey.PolicyName != nil {
			policyName = *contentKey.PolicyName
		}

		contentKeyType := ""
		if contentKey.Type != nil {
			contentKeyType = string(*contentKey.Type)
		}

		value := ""
		if contentKey.Value != nil {
			value = *contentKey.Value
		}

		results = append(results, map[string]interface{}{
			"content_key_id":                      contentKey.Id,
			"label_reference_in_streaming_policy": labelReferenceInStreamingPolicy,
			"policy_name":                         policyName,
			"type":                                contentKeyType,
			"value":                               value,
		})
	}

	return results
}
