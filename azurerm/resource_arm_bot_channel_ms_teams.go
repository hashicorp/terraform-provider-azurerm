package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/botservice/mgmt/2018-07-12/botservice"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmBotChannelMsTeams() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmBotChannelMsTeamsCreate,
		Read:   resourceArmBotChannelMsTeamsRead,
		Delete: resourceArmBotChannelMsTeamsDelete,
		Update: resourceArmBotChannelMsTeamsUpdate,

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
			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"bot_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"calling_web_hook": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.ValidateBotMSTeamsCallingWebHook(),
			},

			"enable_calling": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceArmBotChannelMsTeamsCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Bot.ChannelClient
	ctx, cancel := timeouts.ForCreate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	botName := d.Get("bot_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, string(botservice.ChannelNameMsTeamsChannel), botName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of creating Channel MsTeams for Bot %q (Resource Group %q): %+v", resourceGroup, botName, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_bot_channel_ms_teams", *existing.ID)
		}
	}

	channel := botservice.BotChannel{
		Properties: botservice.MsTeamsChannel{
			Properties: &botservice.MsTeamsChannelProperties{
				EnableCalling:  utils.Bool(d.Get("enable_calling").(bool)),
				CallingWebHook: utils.String(d.Get("calling_web_hook").(string)),
				IsEnabled:      utils.Bool(true),
			},
			ChannelName: botservice.ChannelNameMsTeamsChannel1,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if _, err := client.Create(ctx, resourceGroup, botName, botservice.ChannelNameMsTeamsChannel, channel); err != nil {
		return fmt.Errorf("Error issuing create request for Channel MsTeams for Bot %q (Resource Group %q): %+v", resourceGroup, botName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, botName, string(botservice.ChannelNameMsTeamsChannel))
	if err != nil {
		return fmt.Errorf("Error making get request for Channel MsTeams for Bot %q (Resource Group %q): %+v", resourceGroup, botName, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Channel MsTeams for Bot %q (Resource Group %q): %+v", resourceGroup, botName, err)
	}

	d.SetId(*resp.ID)

	return resourceArmBotChannelMsTeamsRead(d, meta)
}

func resourceArmBotChannelMsTeamsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Bot.ChannelClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	botName := id.Path["botServices"]
	resp, err := client.Get(ctx, id.ResourceGroup, botName, string(botservice.ChannelNameMsTeamsChannel))
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Channel MsTeams for Bot %q (Resource Group %q) was not found - removing from state!", id.ResourceGroup, botName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Channel MsTeams for Bot %q (Resource Group %q): %+v", id.ResourceGroup, botName, err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", resp.Location)
	d.Set("bot_name", botName)

	if props := resp.Properties; props != nil {
		if channel, ok := props.AsMsTeamsChannel(); ok {
			if channelProps := channel.Properties; channelProps != nil {
				d.Set("calling_web_hook", channelProps.CallingWebHook)
				d.Set("enable_calling", channelProps.EnableCalling)
			}
		}
	}

	return nil
}

func resourceArmBotChannelMsTeamsUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Bot.ChannelClient
	ctx, cancel := timeouts.ForUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	botName := d.Get("bot_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	channel := botservice.BotChannel{
		Properties: botservice.MsTeamsChannel{
			Properties: &botservice.MsTeamsChannelProperties{
				EnableCalling:  utils.Bool(d.Get("enable_calling").(bool)),
				CallingWebHook: utils.String(d.Get("calling_web_hook").(string)),
				IsEnabled:      utils.Bool(true),
			},
			ChannelName: botservice.ChannelNameMsTeamsChannel1,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if _, err := client.Update(ctx, resourceGroup, botName, botservice.ChannelNameMsTeamsChannel, channel); err != nil {
		return fmt.Errorf("Error issuing create request for Channel MsTeams for Bot %q (Resource Group %q): %+v", resourceGroup, botName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, botName, string(botservice.ChannelNameMsTeamsChannel))
	if err != nil {
		return fmt.Errorf("Error making get request for Channel MsTeams for Bot %q (Resource Group %q): %+v", resourceGroup, botName, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Channel MsTeams for Bot %q (Resource Group %q): %+v", resourceGroup, botName, err)
	}

	d.SetId(*resp.ID)

	return resourceArmBotChannelMsTeamsRead(d, meta)
}

func resourceArmBotChannelMsTeamsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Bot.ChannelClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	botName := id.Path["botServices"]

	resp, err := client.Delete(ctx, id.ResourceGroup, botName, string(botservice.ChannelNameMsTeamsChannel))
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting Channel MsTeams for Bot %q (Resource Group %q): %+v", id.ResourceGroup, botName, err)
		}
	}

	return nil
}
