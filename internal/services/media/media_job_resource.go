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
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-07-01/encodings"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/media/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMediaJob() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMediaJobCreate,
		Read:   resourceMediaJobRead,
		Update: resourceMediaJobUpdate,
		Delete: resourceMediaJobDelete,

		DeprecationMessage: azureMediaRetirementMessage,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := encodings.ParseJobID(id)
			return err
		}),

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.JobV0ToV1{},
		}),
		SchemaVersion: 1,

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9(_)]{1,128}$"),
					"Job name must be 1 - 128 characters long, can contain letters, numbers, underscores, and hyphens (but the first and last character must be a letter or number).",
				),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"media_services_account_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-z0-9]{3,24}$"),
					"Media Services Account name must be 3 - 24 characters long, contain only lowercase letters and numbers.",
				),
			},

			"transform_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9(_)]{1,128}$"),
					"Transform name must be 1 - 128 characters long, can contain letters, numbers, underscores, and hyphens (but the first and last character must be a letter or number).",
				),
			},

			"input_asset": {
				Type:     pluginsdk.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile("^[-a-zA-Z0-9]{1,128}$"),
								"Asset name must be 1 - 128 characters long, contain only letters, hyphen and numbers.",
							),
						},
						"label": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"output_asset": {
				Type:     pluginsdk.TypeList,
				Required: true,
				ForceNew: true,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile("^[-a-zA-Z0-9]{1,128}$"),
								"Asset name must be 1 - 128 characters long, contain only letters, hyphen and numbers.",
							),
						},
						"label": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"priority": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(encodings.PriorityNormal),
				ValidateFunc: validation.StringInSlice([]string{
					string(encodings.PriorityHigh),
					string(encodings.PriorityNormal),
					string(encodings.PriorityLow),
				}, false),
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceMediaJobCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20220701Client.Encodings
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := encodings.NewJobID(subscriptionId, d.Get("resource_group_name").(string), d.Get("media_services_account_name").(string), d.Get("transform_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.JobsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_media_job", id.ID())
		}
	}

	payload := encodings.Job{
		Properties: &encodings.JobProperties{
			Description: utils.String(d.Get("description").(string)),
		},
	}

	if v, ok := d.GetOk("priority"); ok {
		payload.Properties.Priority = pointer.To(encodings.Priority(v.(string)))
	}

	if v, ok := d.GetOk("input_asset"); ok {
		payload.Properties.Input = expandInputAsset(v.([]interface{}))
	}

	if v, ok := d.GetOk("output_asset"); ok {
		outputAssets, err := expandOutputAssets(v.([]interface{}))
		if err != nil {
			return err
		}
		payload.Properties.Outputs = *outputAssets
	}

	if _, err := client.JobsCreate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceMediaJobRead(d, meta)
}

func resourceMediaJobRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20220701Client.Encodings
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := encodings.ParseJobID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.JobsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.JobName)
	d.Set("transform_name", id.TransformName)
	d.Set("media_services_account_name", id.MediaServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("description", props.Description)

			priority := ""
			if props.Priority != nil {
				priority = string(*props.Priority)
			}
			d.Set("priority", priority)

			inputAsset, err := flattenInputAsset(props.Input)
			if err != nil {
				return err
			}
			if err = d.Set("input_asset", inputAsset); err != nil {
				return fmt.Errorf("flattening `input_asset`: %s", err)
			}

			outputAssets, err := flattenOutputAssets(props.Outputs)
			if err != nil {
				return fmt.Errorf("flattening `output_asset`: %s", err)
			}
			if err = d.Set("output_asset", outputAssets); err != nil {
				return fmt.Errorf("setting `output_asset`: %s", err)
			}
		}
	}

	return nil
}

func resourceMediaJobUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20220701Client.Encodings
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := encodings.ParseJobID(d.Id())
	if err != nil {
		return err
	}

	// TODO: refactor this to use
	payload := encodings.Job{
		Properties: &encodings.JobProperties{
			Description: utils.String(d.Get("description").(string)),
		},
	}

	if v, ok := d.GetOk("priority"); ok {
		payload.Properties.Priority = pointer.To(encodings.Priority(v.(string)))
	}

	if v, ok := d.GetOk("input_asset"); ok {
		inputAsset := expandInputAsset(v.([]interface{}))
		payload.Properties.Input = inputAsset
	}

	if v, ok := d.GetOk("output_asset"); ok {
		outputAssets, err := expandOutputAssets(v.([]interface{}))
		if err != nil {
			return err
		}
		payload.Properties.Outputs = *outputAssets
	}

	if _, err := client.JobsUpdate(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceMediaJobRead(d, meta)
}

func resourceMediaJobDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20220701Client.Encodings
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := encodings.ParseJobID(d.Id())
	if err != nil {
		return err
	}

	// Cancel the job before we attempt to delete it.
	if _, err := client.JobsCancelJob(ctx, *id); err != nil {
		return fmt.Errorf("cancelling %s: %+v", *id, err)
	}

	if _, err := client.JobsDelete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandInputAsset(input []interface{}) encodings.JobInput {
	inputAsset := input[0].(map[string]interface{})
	return &encodings.JobInputAsset{
		AssetName: inputAsset["name"].(string),
		Label:     utils.String(inputAsset["label"].(string)),
	}
}

func flattenInputAsset(input encodings.JobInput) ([]interface{}, error) {
	if input == nil {
		return make([]interface{}, 0), nil
	}

	asset, ok := input.(encodings.JobInputAsset)
	if !ok {
		return nil, fmt.Errorf("Unexpected type for Input Asset. Currently only JobInputAsset is supported.")
	}

	label := ""
	if asset.Label != nil {
		label = *asset.Label
	}

	return []interface{}{
		map[string]interface{}{
			"name":  asset.AssetName,
			"label": label,
		},
	}, nil
}

func expandOutputAssets(input []interface{}) (*[]encodings.JobOutput, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("Job must contain at least one output_asset.")
	}

	outputAssets := make([]encodings.JobOutput, 0)
	for _, output := range input {
		v := output.(map[string]interface{})

		outputAssets = append(outputAssets, encodings.JobOutputAsset{
			AssetName: v["name"].(string),
			Label:     utils.String(v["label"].(string)),
		})
	}

	return &outputAssets, nil
}

func flattenOutputAssets(input []encodings.JobOutput) ([]interface{}, error) {
	outputAssets := make([]interface{}, 0)
	for _, output := range input {
		outputAssetJob, ok := output.(encodings.JobOutputAsset)
		if !ok {
			return nil, fmt.Errorf("unexpected type for output_asset. Currently only JobOutputAsset is supported.")
		}

		label := ""
		if outputAssetJob.Label != nil {
			label = *outputAssetJob.Label
		}

		outputAssets = append(outputAssets, map[string]interface{}{
			"name":  outputAssetJob.AssetName,
			"label": label,
		})
	}
	return outputAssets, nil
}
