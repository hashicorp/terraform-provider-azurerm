package media

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/media/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"

	"github.com/Azure/azure-sdk-for-go/services/mediaservices/mgmt/2020-05-01/media"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/gofrs/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/media/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceMediaStreamingLocator() *schema.Resource {
	return &schema.Resource{
		Create: resourceMediaStreamingLocatorCreate,
		Read:   resourceMediaStreamingLocatorRead,
		Delete: resourceMediaStreamingLocatorDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.StreamingLocatorID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9(_)]{1,128}$"),
					"Steraming Locator name must be 1 - 128 characters long, can contain letters, numbers, underscores, and hyphens (but the first and last character must be a letter or number).",
				),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"media_services_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AccountName,
			},

			"asset_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9]{1,128}$"),
					"Asset name must be 1 - 128 characters long, contain only letters, hyphen and numbers.",
				),
			},

			"streaming_policy_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"alternative_media_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			//lintignore:XS003
			"content_key": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content_key_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.IsUUID,
						},

						"label_reference_in_streaming_policy": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"policy_name": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(media.StreamingLocatorContentKeyTypeCommonEncryptionCbcs),
								string(media.StreamingLocatorContentKeyTypeCommonEncryptionCenc),
								string(media.StreamingLocatorContentKeyTypeEnvelopeEncryption),
							}, false),
						},

						"value": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"default_content_key_policy_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"end_time": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"start_time": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"streaming_locator_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
		},
	}
}

func resourceMediaStreamingLocatorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.StreamingLocatorsClient
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceID := parse.NewStreamingLocatorID(subscriptionID, d.Get("resource_group_name").(string), d.Get("media_services_account_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceID.ResourceGroup, resourceID.MediaserviceName, resourceID.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Media Job %q (Media Service account %q) (ResourceGroup %q): %s", resourceID.Name, resourceID.MediaserviceName, resourceID.ResourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_media_streaming_locator", *existing.ID)
		}
	}

	parameters := media.StreamingLocator{
		StreamingLocatorProperties: &media.StreamingLocatorProperties{
			AssetName:           utils.String(d.Get("asset_name").(string)),
			StreamingPolicyName: utils.String(d.Get("streaming_policy_name").(string)),
		},
	}

	if alternativeMediaID, ok := d.GetOk("alternative_media_id"); ok {
		parameters.StreamingLocatorProperties.AlternativeMediaID = utils.String(alternativeMediaID.(string))
	}

	if contentKeys, ok := d.GetOk("content_key"); ok {
		parameters.StreamingLocatorProperties.ContentKeys = expandContentKeys(contentKeys.([]interface{}))
	}

	if defaultContentKeyPolicyName, ok := d.GetOk("default_content_key_policy_name"); ok {
		parameters.StreamingLocatorProperties.DefaultContentKeyPolicyName = utils.String(defaultContentKeyPolicyName.(string))
	}

	if endTimeRaw, ok := d.GetOk("end_time"); ok {
		if endTimeRaw.(string) != "" {
			endTime, err := date.ParseTime(time.RFC3339, endTimeRaw.(string))
			if err != nil {
				return err
			}
			parameters.StreamingLocatorProperties.EndTime = &date.Time{
				Time: endTime,
			}
		}
	}

	if startTimeRaw, ok := d.GetOk("start_time"); ok {
		if startTimeRaw.(string) != "" {
			startTime, err := date.ParseTime(time.RFC3339, startTimeRaw.(string))
			if err != nil {
				return err
			}
			parameters.StreamingLocatorProperties.StartTime = &date.Time{
				Time: startTime,
			}
		}
	}

	if idRaw, ok := d.GetOk("streaming_locator_id"); ok {
		id := uuid.FromStringOrNil(idRaw.(string))
		parameters.StreamingLocatorProperties.StreamingLocatorID = &id
	}

	if _, err := client.Create(ctx, resourceID.ResourceGroup, resourceID.MediaserviceName, resourceID.Name, parameters); err != nil {
		return fmt.Errorf("Error creating Streaming Locator %q in Media Services Account %q (Resource Group %q): %+v", resourceID.Name, resourceID.MediaserviceName, resourceID.ResourceGroup, err)
	}

	d.SetId(resourceID.ID())

	return resourceMediaStreamingLocatorRead(d, meta)
}

func resourceMediaStreamingLocatorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.StreamingLocatorsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StreamingLocatorID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Streaming Locator %q was not found in Media Services Account %q and Resource Group %q - removing from state", id.Name, id.MediaserviceName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Streaming Locator %q in Media Services Account %q (Resource Group %q): %+v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("media_services_account_name", id.MediaserviceName)

	if props := resp.StreamingLocatorProperties; props != nil {
		d.Set("asset_name", props.AssetName)
		d.Set("streaming_policy_name", props.StreamingPolicyName)
		d.Set("alternative_media_id", props.AlternativeMediaID)
		d.Set("default_content_key_policy_name", props.DefaultContentKeyPolicyName)

		contentKeys := flattenContentKeys(resp.ContentKeys)
		if err := d.Set("content_key", contentKeys); err != nil {
			return fmt.Errorf("Error flattening `content_key`: %s", err)
		}

		endTime := ""
		if props.EndTime != nil {
			endTime = props.EndTime.Format(time.RFC3339)
		}
		d.Set("end_time", endTime)

		startTime := ""
		if props.StartTime != nil {
			startTime = props.StartTime.Format(time.RFC3339)
		}
		d.Set("start_time", startTime)

		id := ""
		if props.StreamingLocatorID != nil {
			id = props.StreamingLocatorID.String()
		}
		d.Set("streaming_locator_id", id)
	}

	return nil
}

func resourceMediaStreamingLocatorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.StreamingLocatorsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StreamingLocatorID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, id.ResourceGroup, id.MediaserviceName, id.Name); err != nil {
		return fmt.Errorf("Error deleting Streaming Locator %q in Media Services Account %q (Resource Group %q): %+v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
	}

	return nil
}

func expandContentKeys(input []interface{}) *[]media.StreamingLocatorContentKey {
	results := make([]media.StreamingLocatorContentKey, 0)

	for _, contentKeyRaw := range input {
		if contentKeyRaw == nil {
			continue
		}
		contentKey := contentKeyRaw.(map[string]interface{})

		streamingLocatorContentKey := media.StreamingLocatorContentKey{}

		if contentKey["content_key_id"] != nil {
			id := uuid.FromStringOrNil(contentKey["content_key_id"].(string))
			streamingLocatorContentKey.ID = &id
		}

		if contentKey["label_reference_in_streaming_policy"] != nil {
			streamingLocatorContentKey.LabelReferenceInStreamingPolicy = utils.String(contentKey["label_reference_in_streaming_policy"].(string))
		}

		if contentKey["policy_name"] != nil {
			streamingLocatorContentKey.PolicyName = utils.String(contentKey["policy_name"].(string))
		}

		if contentKey["type"] != nil {
			streamingLocatorContentKey.Type = media.StreamingLocatorContentKeyType(contentKey["type"].(string))
		}

		if contentKey["value"] != nil {
			streamingLocatorContentKey.Value = utils.String(contentKey["value"].(string))
		}

		results = append(results, streamingLocatorContentKey)
	}

	return &results
}

func flattenContentKeys(input *[]media.StreamingLocatorContentKey) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, contentKey := range *input {
		id := ""
		if contentKey.ID != nil {
			id = contentKey.ID.String()
		}

		labelReferenceInStreamingPolicy := ""
		if contentKey.LabelReferenceInStreamingPolicy != nil {
			labelReferenceInStreamingPolicy = *contentKey.LabelReferenceInStreamingPolicy
		}

		policyName := ""
		if contentKey.PolicyName != nil {
			policyName = *contentKey.PolicyName
		}

		value := ""
		if contentKey.Value != nil {
			value = *contentKey.Value
		}

		results = append(results, map[string]interface{}{
			"content_key_id":                      id,
			"label_reference_in_streaming_policy": labelReferenceInStreamingPolicy,
			"policy_name":                         policyName,
			"type":                                string(contentKey.Type),
			"value":                               value,
		})
	}

	return results
}
