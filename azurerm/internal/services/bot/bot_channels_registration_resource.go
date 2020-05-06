package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/botservice/mgmt/2018-07-12/botservice"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmBotChannelsRegistration() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmBotChannelsRegistrationCreate,
		Read:   resourceArmBotChannelsRegistrationRead,
		Delete: resourceArmBotChannelsRegistrationDelete,
		Update: resourceArmBotChannelsRegistrationUpdate,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"sku": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(botservice.F0),
					string(botservice.S1),
				}, false),
			},

			"microsoft_app_id": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},

			"display_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"developer_app_insights_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsUUID,
			},

			"developer_app_insights_api_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"developer_app_insights_application_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsUUID,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmBotChannelsRegistrationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.BotClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of creating Bot Channels Registration %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_bot_channels_registration", *existing.ID)
		}
	}

	displayName := d.Get("display_name").(string)
	if displayName == "" {
		displayName = name
	}

	bot := botservice.Bot{
		Properties: &botservice.BotProperties{
			DisplayName:                       utils.String(displayName),
			Endpoint:                          utils.String(d.Get("endpoint").(string)),
			MsaAppID:                          utils.String(d.Get("microsoft_app_id").(string)),
			DeveloperAppInsightKey:            utils.String(d.Get("developer_app_insights_key").(string)),
			DeveloperAppInsightsAPIKey:        utils.String(d.Get("developer_app_insights_api_key").(string)),
			DeveloperAppInsightsApplicationID: utils.String(d.Get("developer_app_insights_application_id").(string)),
		},
		Location: utils.String(d.Get("location").(string)),
		Sku: &botservice.Sku{
			Name: botservice.SkuName(d.Get("sku").(string)),
		},
		Kind: botservice.KindBot,
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.Create(ctx, resourceGroup, name, bot); err != nil {
		return fmt.Errorf("Error issuing create request for Bot Channels Registration %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error making get request for Bot Channels Registration %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Bot Channels Registration %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceArmBotChannelsRegistrationRead(d, meta)
}

func resourceArmBotChannelsRegistrationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.BotClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	name := id.Path["botServices"]

	resp, err := client.Get(ctx, id.ResourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Bot Channels Registration %q (Resource Group %q) - removing from state", name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Bot Channels Registration %q (Resource Group %q): %+v", name, id.ResourceGroup, err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("name", resp.Name)
	d.Set("location", resp.Location)

	if sku := resp.Sku; sku != nil {
		d.Set("sku", string(sku.Name))
	}

	if props := resp.Properties; props != nil {
		d.Set("microsoft_app_id", props.MsaAppID)
		d.Set("endpoint", props.Endpoint)
		d.Set("display_name", props.DisplayName)
		d.Set("developer_app_insights_key", props.DeveloperAppInsightKey)
		d.Set("developer_app_insights_application_id", props.DeveloperAppInsightsApplicationID)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmBotChannelsRegistrationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.BotClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	t := d.Get("tags").(map[string]interface{})
	displayName := d.Get("display_name").(string)
	if displayName == "" {
		displayName = name
	}

	bot := botservice.Bot{
		Properties: &botservice.BotProperties{
			DisplayName:                       utils.String(displayName),
			Endpoint:                          utils.String(d.Get("endpoint").(string)),
			MsaAppID:                          utils.String(d.Get("microsoft_app_id").(string)),
			DeveloperAppInsightKey:            utils.String(d.Get("developer_app_insights_key").(string)),
			DeveloperAppInsightsAPIKey:        utils.String(d.Get("developer_app_insights_api_key").(string)),
			DeveloperAppInsightsApplicationID: utils.String(d.Get("developer_app_insights_application_id").(string)),
		},
		Location: utils.String(d.Get("location").(string)),
		Sku: &botservice.Sku{
			Name: botservice.SkuName(d.Get("sku").(string)),
		},
		Kind: botservice.KindBot,
		Tags: tags.Expand(t),
	}

	if _, err := client.Update(ctx, resourceGroup, name, bot); err != nil {
		return fmt.Errorf("Error issuing update request for Bot Channels Registration %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error making get request for Bot Channels Registration %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Bot Channels Registration %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceArmBotChannelsRegistrationRead(d, meta)
}

func resourceArmBotChannelsRegistrationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.BotClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	name := id.Path["botServices"]

	resp, err := client.Delete(ctx, id.ResourceGroup, name)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting Bot Channels Registration %q (Resource Group %q): %+v", name, id.ResourceGroup, err)
		}
	}

	return nil
}
