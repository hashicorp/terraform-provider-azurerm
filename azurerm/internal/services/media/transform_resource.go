package media

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mediaservices/mgmt/2020-05-01/media"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/media/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceTransform() *schema.Resource {
	return &schema.Resource{
		Create: resourceTransformCreateUpdate,
		Read:   resourceTransformRead,
		Update: resourceTransformCreateUpdate,
		Delete: resourceTransformDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.TransformID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[_-a-zA-Z0-9]{1,128}$"),
					"Media Services Account name must be 1 - 128 characters long, can contain letters, numbers, underscores, and hyphens (but the first and last character must be a letter or number).",
				),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"media_services_account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-z0-9]{3,24}$"),
					"Media Services Account name must be 3 - 24 characters long, contain only lowercase letters and numbers.",
				),
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"output": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"on_error_type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"ContinueJob", "StopProcessingJob",
							}, true),
						},
						"preset": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"BuiltInStandardEncoderPreset", "AudioAnalyzerPreset",
										}, true),
									},

									"preset_name": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"AACGoodQualityAudio", "AdaptiveStreaming",
											"ContentAwareEncoding", "ContentAwareEncodingExperimental",
											"CopyAllBitrateNonInterleaved", "H264MultipleBitrate1080p",
											"H264MultipleBitrate720p", "H264MultipleBitrateSD",
											"H264SingleBitrate1080p", "H264SingleBitrate720p", "H264SingleBitrateSD",
										}, true),
									},
								},
							},
						},
						"relative_priority": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"High", "Normal", "Low",
							}, true),
						},
					},
				},
			},
		},
	}
}

func resourceTransformCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.TransformsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	transformName := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("media_services_account_name").(string)
	description := d.Get("description").(string)

	parameters := media.Transform{
		TransformProperties: &media.TransformProperties{
			Description: utils.String(description),
		},
	}

	if v, ok := d.GetOk("output"); ok {
		parameters.Outputs = expandTransformOuputs(v.(*schema.Set).List())
	}

	_, err := client.CreateOrUpdate(ctx, resourceGroup, accountName, transformName, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Transform %q in Media Services Account %q (Resource Group %q): %+v", transformName, accountName, resourceGroup, err)
	}

	transform, err := client.Get(ctx, resourceGroup, accountName, transformName)
	if err != nil {
		return fmt.Errorf("Error retrieving Transform %q from Media Services Account %q (Resource Group %q): %+v", transformName, accountName, resourceGroup, err)
	}

	d.SetId(*transform.ID)

	return resourceTransformRead(d, meta)
}

func resourceTransformRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.TransformsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TransformID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Transform %q was not found in Media Services Account %q and Resource Group %q - removing from state", id.Name, id.MediaserviceName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Transform %q in Media Services Account %q (Resource Group %q): %+v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("media_services_account_name", id.MediaserviceName)

	if props := resp.TransformProperties; props != nil {
		if description := props.Description; description != nil {
			d.Set("description", description)
		}

		outputs := flattenTransformOutputs(props.Outputs)
		if err := d.Set("output", outputs); err != nil {
			return fmt.Errorf("Error flattening `output`: %s", err)
		}
	}

	return nil
}

func resourceTransformDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.TransformsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TransformID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)

	if err != nil {
		if response.WasNotFound(resp.Response) {
			return nil
		}
		return fmt.Errorf("Error deleting Transform %q in Media Services Account %q (Resource Group %q): %+v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
	}

	return nil
}

func expandTransformOuputs(input []interface{}) *[]media.TransformOutput {
	results := make([]media.TransformOutput, 0)

	for _, transformOuputRaw := range input {
		transform := transformOuputRaw.(map[string]interface{})

		onError := transform["on_error_type"].(string)
		relativePriority := transform["relative_priority"].(string)
		preset := expandPreset(transform["preset"].(*schema.Set).List())

		transformOuput := media.TransformOutput{
			OnError:          media.OnErrorType(onError),
			RelativePriority: media.Priority(relativePriority),
			Preset:           preset,
		}

		results = append(results, transformOuput)
	}

	return &results

}

func flattenTransformOutputs(input *[]media.TransformOutput) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, transformOuput := range *input {
		output := make(map[string]interface{})
		output["on_error_type"] = string(transformOuput.OnError)
		output["relative_priority"] = string(transformOuput.RelativePriority)
		results = append(results, output)
	}

	return results
}

func expandPreset(presets []interface{}) *media.Preset {
	preset := presets[0].(map[string]interface{})
	presetType := preset["type"].(string)
	switch presetType {
	case "BuiltInStandardEncoderPreset":
		presetName := preset["name"].(string)
		builtInPreset := media.BuiltInStandardEncoderPreset{
			PresetName: media.EncoderNamedPreset(presetName),
			OdataType:  media.OdataTypeMicrosoftMediaBuiltInStandardEncoderPreset,
		}
		preset, _ := builtInPreset.AsPreset()
		return preset
	}

	return nil
}
