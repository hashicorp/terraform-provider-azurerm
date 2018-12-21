package azurerm

import (
	"fmt"
	"log"
	"net/http"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"

	"github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2015-05-01/insights"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmApplicationInsightsAPIKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApplicationInsightsAPIKeyCreate,
		Read:   resourceArmApplicationInsightsAPIKeyRead,
		Delete: resourceArmApplicationInsightsAPIKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"api_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"application_insights_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"read_permissions": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Set:      schema.HashString,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validate.LowercaseString,
				},
			},

			"write_permissions": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Set:      schema.HashString,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validate.LowercaseString,
				},
			},
		},
	}
}

func resourceArmApplicationInsightsAPIKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appInsightsAPIKeyClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM Application Insights API key creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	appInsightsName := d.Get("application_insights_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, appInsightsName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Application Insights API key %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_application_insights_api_key", *existing.ID)
		}
	}

	apiKeyProperties := insights.APIKeyRequest{
		Name:                  &name,
		LinkedReadProperties:  azure.ExpandApplicationInsightsAPIKeyLinkedProperties(d.Get("read_permissions").(*schema.Set), client.SubscriptionID, resGroup, appInsightsName),
		LinkedWriteProperties: azure.ExpandApplicationInsightsAPIKeyLinkedProperties(d.Get("write_permissions").(*schema.Set), client.SubscriptionID, resGroup, appInsightsName),
	}

	result, err := client.Create(ctx, resGroup, appInsightsName, apiKeyProperties)
	if err != nil {
		return fmt.Errorf("Error creating Application Insights API key %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if result.Response.StatusCode != http.StatusOK {
		return fmt.Errorf("Error creating Application Insights API key %q (Resource Group %q): %+v", name, resGroup, result.Response)
	}

	if result.APIKey == nil {
		return fmt.Errorf("Error creating Application Insights API key %q (Resource Group %q): got empty API key", name, resGroup)
	}

	d.SetId(*result.ID)

	// API key can only retrieved at key creation
	d.Set("api_key", result.APIKey)

	return resourceArmApplicationInsightsAPIKeyRead(d, meta)
}

func resourceArmApplicationInsightsAPIKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appInsightsAPIKeyClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Reading AzureRM Application Insights API key '%s'", id)

	resGroup := id.ResourceGroup
	appInsightsName := id.Path["components"]
	keyID := id.Path["apikeys"]

	resp, err := client.Get(ctx, resGroup, appInsightsName, keyID)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM Application Insights API key '%s': %+v", keyID, err)
	}

	d.Set("resource_group_name", resGroup)
	d.Set("application_insights_name", appInsightsName)

	d.Set("name", resp.Name)
	readProps := azure.FlattenApplicationInsightsAPIKeyLinkedProperties(resp.LinkedReadProperties)
	if err := d.Set("read_permissions", readProps); err != nil {
		return fmt.Errorf("Error flattening `read_permissions `: %s", err)
	}
	writeProps := azure.FlattenApplicationInsightsAPIKeyLinkedProperties(resp.LinkedWriteProperties)
	if err := d.Set("write_permissions", writeProps); err != nil {
		return fmt.Errorf("Error flattening `write_permissions `: %s", err)
	}

	return nil
}

func resourceArmApplicationInsightsAPIKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appInsightsAPIKeyClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	appInsightsName := id.Path["components"]
	keyID := id.Path["apikeys"]

	log.Printf("[DEBUG] Deleting AzureRM Application Insights API key '%s' (resource group '%s')", keyID, resGroup)

	resp, err := client.Delete(ctx, resGroup, appInsightsName, keyID)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("Error issuing AzureRM delete request for Application Insights API key '%s': %+v", keyID, err)
	}

	return nil
}
