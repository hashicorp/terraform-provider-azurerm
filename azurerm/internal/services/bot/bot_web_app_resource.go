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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/bot/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmBotWebApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmBotWebAppCreate,
		Read:   resourceArmBotWebAppRead,
		Update: resourceArmBotWebAppUpdate,
		Delete: resourceArmBotWebAppDelete,
		Importer: azSchema.ValidateResourceIDPriorToImportThen(func(id string) error {
			_, err := parse.BotServiceID(id)
			return err
		}, func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
			client := meta.(*clients.Client).Bot.BotClient
			ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
			defer cancel()

			id, err := parse.BotServiceID(d.Id())
			if err != nil {
				return nil, err
			}

			resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return nil, fmt.Errorf("Web App Bot %q was not found in Resource Group %q", id.Name, id.ResourceGroup)
				}

				return nil, fmt.Errorf("retrieving Web App Bot %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
			}
			if resp.Kind != botservice.KindSdk {
				return nil, fmt.Errorf("Bot %q (Resource Group %q) was not a Web App - got %q", id.Name, id.ResourceGroup, string(resp.Kind))
			}

			return []*schema.ResourceData{d}, nil
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

			"luis_app_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsUUID,
				},
			},

			"luis_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmBotWebAppCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.BotClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId := parse.NewBotServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceId.ResourceGroup, resourceId.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of creating Web App Bot %q (Resource Group %q): %+v", resourceId.Name, resourceId.ResourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_bot_web_app", resourceId.ID(""))
		}
	}

	displayName := d.Get("display_name").(string)
	if displayName == "" {
		displayName = resourceId.Name
	}

	bot := botservice.Bot{
		Properties: &botservice.BotProperties{
			DisplayName:                       utils.String(displayName),
			Endpoint:                          utils.String(d.Get("endpoint").(string)),
			MsaAppID:                          utils.String(d.Get("microsoft_app_id").(string)),
			DeveloperAppInsightKey:            utils.String(d.Get("developer_app_insights_key").(string)),
			DeveloperAppInsightsAPIKey:        utils.String(d.Get("developer_app_insights_api_key").(string)),
			DeveloperAppInsightsApplicationID: utils.String(d.Get("developer_app_insights_application_id").(string)),
			LuisAppIds:                        utils.ExpandStringSlice(d.Get("luis_app_ids").([]interface{})),
			LuisKey:                           utils.String(d.Get("luis_key").(string)),
		},
		Location: utils.String(d.Get("location").(string)),
		Sku: &botservice.Sku{
			Name: botservice.SkuName(d.Get("sku").(string)),
		},
		Kind: botservice.KindSdk,
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.Create(ctx, resourceId.ResourceGroup, resourceId.Name, bot); err != nil {
		return fmt.Errorf("creating Web App Bot %q (Resource Group %q): %+v", resourceId.Name, resourceId.ResourceGroup, err)
	}

	d.SetId(resourceId.ID(""))
	return resourceArmBotWebAppRead(d, meta)
}

func resourceArmBotWebAppRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.BotClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Web App Bot %q was not found in Resource Group %q - removing from state", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Web App Bot %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if sku := resp.Sku; sku != nil {
		d.Set("sku", string(sku.Name))
	}

	if props := resp.Properties; props != nil {
		d.Set("microsoft_app_id", props.MsaAppID)
		d.Set("endpoint", props.Endpoint)
		d.Set("display_name", props.DisplayName)
		d.Set("developer_app_insights_key", props.DeveloperAppInsightKey)
		d.Set("developer_app_insights_application_id", props.DeveloperAppInsightsApplicationID)
		d.Set("luis_app_ids", props.LuisAppIds)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmBotWebAppUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.BotClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotServiceID(d.Id())
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	if displayName == "" {
		displayName = id.Name
	}

	bot := botservice.Bot{
		Properties: &botservice.BotProperties{
			DisplayName:                       utils.String(displayName),
			Endpoint:                          utils.String(d.Get("endpoint").(string)),
			MsaAppID:                          utils.String(d.Get("microsoft_app_id").(string)),
			DeveloperAppInsightKey:            utils.String(d.Get("developer_app_insights_key").(string)),
			DeveloperAppInsightsAPIKey:        utils.String(d.Get("developer_app_insights_api_key").(string)),
			DeveloperAppInsightsApplicationID: utils.String(d.Get("developer_app_insights_application_id").(string)),
			LuisAppIds:                        utils.ExpandStringSlice(d.Get("luis_app_ids").([]interface{})),
			LuisKey:                           utils.String(d.Get("luis_key").(string)),
		},
		Location: utils.String(d.Get("location").(string)),
		Sku: &botservice.Sku{
			Name: botservice.SkuName(d.Get("sku").(string)),
		},
		Kind: botservice.KindSdk,
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.Name, bot); err != nil {
		return fmt.Errorf("updating Web App Bot %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return resourceArmBotWebAppRead(d, meta)
}

func resourceArmBotWebAppDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.BotClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("deleting Web App Bot %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}
