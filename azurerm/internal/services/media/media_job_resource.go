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

func resourceMediaJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceMediaJobCreate,
		Read:   resourceMediaJobRead,
		Update: resourceMediaJobUpdate,
		Delete: resourceMediaJobDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.JobID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9(_)]{1,128}$"),
					"Job name must be 1 - 128 characters long, can contain letters, numbers, underscores, and hyphens (but the first and last character must be a letter or number).",
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

			"transform_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9(_)]{1,128}$"),
					"Transform name must be 1 - 128 characters long, can contain letters, numbers, underscores, and hyphens (but the first and last character must be a letter or number).",
				),
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"priority": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(media.High), string(media.Normal), string(media.Low),
				}, true),
			},

			"input_asset": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"asset_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile("^[-a-zA-Z0-9]{1,128}$"),
								"Asset name must be 1 - 128 characters long, contain only letters, hyphen and numbers.",
							),
						},
					},
				},
			},
			"output_asset": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"asset_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile("^[-a-zA-Z0-9]{1,128}$"),
								"Asset name must be 1 - 128 characters long, contain only letters, hyphen and numbers.",
							),
						},
					},
				},
			},
		},
	}
}

func resourceMediaJobCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.JobsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	jobName := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("media_services_account_name").(string)
	transformName := d.Get("transform_name").(string)
	description := d.Get("description").(string)

	parameters := media.Job{
		JobProperties: &media.JobProperties{
			Description: utils.String(description),
		},
	}

	if v, ok := d.GetOk("priority"); ok {
		parameters.Priority = media.Priority(v.(string))
	}

	if v, ok := d.GetOk("input_asset"); ok {
		inputAsset, err := expandInputAsset(v.([]interface{}))
		if err != nil {
			return err
		}
		parameters.JobProperties.Input = inputAsset
	}

	if v, ok := d.GetOk("output_asset"); ok {
		outputAssets, err := expandOutputAssets(v.([]interface{}))
		if err != nil {
			return err
		}
		parameters.JobProperties.Outputs = outputAssets
	}

	if _, err := client.Create(ctx, resourceGroup, accountName, transformName, jobName, parameters); err != nil {
		return fmt.Errorf("Error creating Job %q in Media Services Account %q (Resource Group %q): %+v", jobName, accountName, resourceGroup, err)
	}

	job, err := client.Get(ctx, resourceGroup, accountName, transformName, jobName)
	if err != nil {
		return fmt.Errorf("Error retrieving Job %q from Media Services Account %q (Resource Group %q): %+v", jobName, accountName, resourceGroup, err)
	}

	d.SetId(*job.ID)

	return resourceMediaJobRead(d, meta)
}

func resourceMediaJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.JobsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := parse.JobID(d.Id())
	if err != nil {
		return err
	}
	resp, err := client.Get(ctx, id.ResourceGroup, id.MediaserviceName, id.TransformName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Job %q was not found in Media Services Account %q and Resource Group %q - removing from state", id.Name, id.MediaserviceName, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Job %q in Media Services Account %q (Resource Group %q): %+v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("media_services_account_name", id.MediaserviceName)
	d.Set("transform_name", id.TransformName)

	if props := resp.JobProperties; props != nil {
		if description := props.Description; description != nil {
			d.Set("description", description)
		}
		d.Set("priority", string(props.Priority))
		inputAsset := flattenInputAsset(props.Input)
		if err := d.Set("input_asset", inputAsset); err != nil {
			return fmt.Errorf("Error flattening `input_asset`: %s", err)
		}

		outputAssets := flattenOutputAssets(props.Outputs)
		if err := d.Set("output_asset", outputAssets); err != nil {
			return fmt.Errorf("Error flattening `output_asset`: %s", err)
		}
	}
	return nil
}

func resourceMediaJobUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.JobsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	jobName := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("media_services_account_name").(string)
	transformName := d.Get("transform_name").(string)
	description := d.Get("description").(string)

	parameters := media.Job{
		JobProperties: &media.JobProperties{
			Description: utils.String(description),
		},
	}

	if v, ok := d.GetOk("priority"); ok {
		parameters.Priority = media.Priority(v.(string))
	}

	if v, ok := d.GetOk("input_asset"); ok {
		inputAsset, err := expandInputAsset(v.([]interface{}))
		if err != nil {
			return err
		}
		parameters.JobProperties.Input = inputAsset
	}

	if v, ok := d.GetOk("output_asset"); ok {
		outputAssets, err := expandOutputAssets(v.([]interface{}))
		if err != nil {
			return err
		}
		parameters.JobProperties.Outputs = outputAssets
	}

	if _, err := client.Update(ctx, resourceGroup, accountName, transformName, jobName, parameters); err != nil {
		return fmt.Errorf("Error creating Job %q in Media Services Account %q (Resource Group %q): %+v", jobName, accountName, resourceGroup, err)
	}

	job, err := client.Get(ctx, resourceGroup, accountName, transformName, jobName)
	if err != nil {
		return fmt.Errorf("Error retrieving Job %q from Media Services Account %q (Resource Group %q): %+v", jobName, accountName, resourceGroup, err)
	}

	d.SetId(*job.ID)

	return resourceMediaJobRead(d, meta)
}

func resourceMediaJobDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.JobsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.JobID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.MediaserviceName, id.TransformName, id.Name)
	if err != nil {
		if response.WasNotFound(resp.Response) {
			return nil
		}
		return fmt.Errorf("Error deleting Job %q in Media Services Account %q (Resource Group %q): %+v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
	}

	return nil
}

func expandInputAsset(input []interface{}) (media.BasicJobInput, error) {
	inputAsset := input[0].(map[string]interface{})
	assetName := inputAsset["asset_name"].(string)
	return &media.JobInputAsset{
		AssetName: utils.String(assetName),
	}, nil
}

func flattenInputAsset(input media.BasicJobInput) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	asset, _ := input.AsJobInputAsset()
	result := make(map[string]interface{})
	if asset.AssetName != nil {
		result["asset_name"] = *asset.AssetName
	}
	return []interface{}{result}
}

func expandOutputAssets(outputs []interface{}) (*[]media.BasicJobOutput, error) {
	if len(outputs) == 0 {
		return nil, fmt.Errorf("Job must contain at least one output_asset.")
	}
	outputAssets := make([]media.BasicJobOutput, len(outputs))
	for index, output := range outputs {
		outputAsset := output.(map[string]interface{})
		assetName := outputAsset["asset_name"].(string)

		jobOutputAsset := media.JobOutputAsset{
			AssetName: utils.String(assetName),
		}
		outputAssets[index] = jobOutputAsset
	}

	return &outputAssets, nil
}

func flattenOutputAssets(outputs *[]media.BasicJobOutput) []interface{} {
	if outputs == nil || len(*outputs) == 0 {
		return []interface{}{}
	}

	outputAssets := make([]interface{}, len(*outputs))
	for i, output := range *outputs {
		outputAsset := make(map[string]interface{})
		outputAssetJob, _ := output.AsJobOutputAsset()
		if outputAssetJob.AssetName != nil {
			outputAsset["asset_name"] = outputAssetJob.AssetName
		}
		if outputAsset != nil {
			outputAssets[i] = outputAsset
		}
	}

	return outputAssets
}
