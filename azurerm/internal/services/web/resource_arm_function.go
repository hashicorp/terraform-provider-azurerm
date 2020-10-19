package web

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2019-08-01/web"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	webValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmFunction() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmFunctionCreate,
		Read:   resourceArmFunctionRead,
		Update: resourceArmFunctionCreate,
		Delete: resourceArmFunctionDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ParseFunctionID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: webValidate.AppServiceName,
			},

			"app_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"config": {
				Type:         schema.TypeString,
				Required:     true,
				StateFunc:    azure.NormalizeJson,
				ValidateFunc: validation.All(validation.StringIsNotEmpty, validation.StringIsJSON),
			},

			"files": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"test_data": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"href": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"script_href": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"script_root_href": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secrets_file_href": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"config_href": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"test_data_href": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"invoke_url_template": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"language": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmFunctionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Function creation.")

	name := d.Get("name").(string)
	functionAppName := d.Get("app_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.GetFunction(ctx, resourceGroup, functionAppName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Function %q (Resource Group %q, Function App %q): %s", name, resourceGroup, functionAppName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_function", *existing.ID)
		}
	}

	config, err := expandFunctionConfig(d.Get("config"))
	if err != nil {
		return err
	}
	files := expandFunctionFiles(d.Get("files"))
	testData := d.Get("test_data").(string)

	functionEnvelope := web.FunctionEnvelope{
		FunctionEnvelopeProperties: &web.FunctionEnvelopeProperties{
			Files:    files,
			Config:   config,
			TestData: &testData,
		},
	}

	log.Printf("[INFO] preparing arguments for AzureRM Function creation with Properties: %+v.", functionEnvelope)

	createFuture, err := client.CreateFunction(ctx, resourceGroup, functionAppName, name, functionEnvelope)
	if err != nil {
		return fmt.Errorf("creating/updating Function %q (Resource Group %q): %s", name, resourceGroup, err)
	}

	err = createFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("waiting for Function %q (Resource Group %q) to become available: %s", name, resourceGroup, err)
	}

	read, err := client.GetFunction(ctx, resourceGroup, functionAppName, name)
	if err != nil {
		return fmt.Errorf("retrieving Function %q (Resource Group %q): %s", name, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("reading Function %s (Resource Group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmFunctionRead(d, meta)
}

func resourceArmFunctionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseFunctionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetFunction(ctx, id.ResourceGroup, id.FunctionAppName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Function %q (Resource Group %q, Function App %q) was not found - removing from state", id.Name, id.ResourceGroup, id.FunctionAppName)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on AzureRM Function %q: %+v", id.Name, err)
	}

	d.Set("name", id.Name)
	d.Set("app_name", id.FunctionAppName)
	d.Set("resource_group_name", id.ResourceGroup)
	config, err := flattenFunctionConfig(resp.Config)
	if err != nil {
		return err
	}
	d.Set("config", config)
	d.Set("files", resp.Files)
	d.Set("test_data", resp.TestData)
	d.Set("type", resp.Type)
	d.Set("language", resp.Language)
	d.Set("href", resp.Href)
	d.Set("script_href", resp.ScriptHref)
	d.Set("script_root_href", resp.ScriptRootPathHref)
	d.Set("secrets_file_href", resp.SecretsFileHref)
	d.Set("config_href", resp.ConfigHref)
	d.Set("test_data_href", resp.TestDataHref)
	d.Set("invoke_url_template", resp.InvokeURLTemplate)

	return nil
}

func resourceArmFunctionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseFunctionID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting Function %q (Resource Group %q, Function App %q)", id.Name, id.ResourceGroup, id.FunctionAppName)

	resp, err := client.DeleteFunction(ctx, id.ResourceGroup, id.FunctionAppName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return err
		}
	}

	return nil
}

func expandFunctionFiles(input interface{}) map[string]*string {
	files := input.(map[string]interface{})

	return utils.ExpandMapStringPtrString(files)
}

func expandFunctionConfig(input interface{}) (map[string]interface{}, error) {
	cfg := input.(string)

	var config map[string]interface{}
	if err := json.Unmarshal([]byte(cfg), &config); err != nil {
		return nil, fmt.Errorf("parsing Function config: %s", err)
	}

	return config, nil
}

func flattenFunctionConfig(input interface{}) (*string, error) {
	config, err := json.Marshal(input)

	if err != nil {
		return nil, fmt.Errorf("serializing config to JSON: %+v", err)
	}

	return utils.String(string(config)), nil
}
