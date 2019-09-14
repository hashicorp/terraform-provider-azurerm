package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/botservice/mgmt/2018-07-12/botservice"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmBotConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmBotConnectionCreate,
		Read:   resourceArmBotConnectionRead,
		Update: resourceArmBotConnectionUpdate,
		Delete: resourceArmBotConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"bot_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"service_provider_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"client_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"client_secret": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"scopes": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmBotConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).bot.ConnectionClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	botName := d.Get("bot_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name, botName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of creating Bot Connection %q (Resource Group %q / Bot %q): %+v", name, resourceGroup, botName, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_bot_connection", *existing.ID)
		}
	}

	connection := botservice.ConnectionSetting{
		Properties: &botservice.ConnectionSettingProperties{
			ServiceProviderDisplayName: utils.String(d.Get("service_provider_name").(string)),
			ClientID:                   utils.String(d.Get("client_id").(string)),
			ClientSecret:               utils.String(d.Get("client_secret").(string)),
			Scopes:                     utils.String(d.Get("scopes").(string)),
			Parameters:                 expandAzureRMBotConnectionParameters(d.Get("parameters").(map[string]interface{})),
		},
		Kind:     botservice.KindBot,
		Location: utils.String(d.Get("location").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.Create(ctx, resourceGroup, botName, name, connection); err != nil {
		return fmt.Errorf("Error issuing create request for creating Bot Connection %q (Resource Group %q / Bot %q): %+v", name, resourceGroup, botName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, botName, name)
	if err != nil {
		return fmt.Errorf("Error making get request for Bot Connection %q (Resource Group %q / Bot %q): %+v", name, resourceGroup, botName, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Bot Connection %q (Resource Group %q / Bot %q): %+v", name, resourceGroup, botName, err)
	}

	d.SetId(*resp.ID)

	return resourceArmBotConnectionRead(d, meta)
}

func resourceArmBotConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).bot.ConnectionClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	botName := id.Path["botServices"]
	name := id.Path["connections"]

	resp, err := client.Get(ctx, id.ResourceGroup, botName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Bot Connection %q (Resource Group %q / Bot %q)", name, id.ResourceGroup, botName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Bot Connection %q (Resource Group %q / Bot %q): %+v", name, id.ResourceGroup, botName, err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("name", name)
	d.Set("bot_name", botName)
	d.Set("location", resp.Location)

	if props := resp.Properties; props != nil {
		d.Set("client_id", props.ClientID)
		d.Set("scopes", props.Scopes)
		if err := d.Set("parameters", flattenAzureRMBotConnectionParameters(props.Parameters)); err != nil {
			return fmt.Errorf("Error setting `parameters`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmBotConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).bot.ConnectionClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	botName := d.Get("bot_name").(string)

	connection := botservice.ConnectionSetting{
		Properties: &botservice.ConnectionSettingProperties{
			ServiceProviderDisplayName: utils.String(d.Get("service_provider_name").(string)),
			ClientID:                   utils.String(d.Get("client_id").(string)),
			ClientSecret:               utils.String(d.Get("client_secret").(string)),
			Scopes:                     utils.String(d.Get("scopes").(string)),
			Parameters:                 expandAzureRMBotConnectionParameters(d.Get("parameters").(map[string]interface{})),
		},
		Kind:     botservice.KindBot,
		Location: utils.String(d.Get("location").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.Update(ctx, resourceGroup, botName, name, connection); err != nil {
		return fmt.Errorf("Error issuing update request for creating Bot Connection %q (Resource Group %q / Bot %q): %+v", name, resourceGroup, botName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, botName, name)
	if err != nil {
		return fmt.Errorf("Error making get request for Bot Connection %q (Resource Group %q / Bot %q): %+v", name, resourceGroup, botName, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Bot Connection %q (Resource Group %q / Bot %q): %+v", name, resourceGroup, botName, err)
	}

	d.SetId(*resp.ID)

	return resourceArmBotConnectionRead(d, meta)

}

func resourceArmBotConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).bot.ConnectionClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	botName := id.Path["botServices"]
	name := id.Path["connections"]

	resp, err := client.Delete(ctx, id.ResourceGroup, botName, name)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting Bot Connection %q (Resource Group %q / Bot %q): %+v", name, id.ResourceGroup, botName, err)
		}
	}

	return nil
}

func expandAzureRMBotConnectionParameters(input map[string]interface{}) *[]botservice.ConnectionSettingParameter {
	output := make([]botservice.ConnectionSettingParameter, 0)

	for k, v := range input {
		output = append(output, botservice.ConnectionSettingParameter{
			Key:   utils.String(k),
			Value: utils.String(v.(string)),
		})
	}
	return &output
}

func flattenAzureRMBotConnectionParameters(input *[]botservice.ConnectionSettingParameter) map[string]interface{} {
	output := make(map[string]interface{})
	if input == nil {
		return output
	}

	for _, parameter := range *input {
		if key := parameter.Key; key != nil {
			// We disregard the clientSecret and clientId as one is sensitive and the other is returned in the ClientId attribute.
			if *key != "clientSecret" && *key != "clientId" && *key != "scopes" {
				if value := parameter.Value; value != nil {
					output[*key] = *value
				}
			}
		}
	}

	return output
}
