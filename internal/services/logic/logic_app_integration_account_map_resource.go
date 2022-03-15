package logic

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/logic/mgmt/2019-05-01/logic"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceLogicAppIntegrationAccountMap() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogicAppIntegrationAccountMapCreateUpdate,
		Read:   resourceLogicAppIntegrationAccountMapRead,
		Update: resourceLogicAppIntegrationAccountMapCreateUpdate,
		Delete: resourceLogicAppIntegrationAccountMapDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.IntegrationAccountMapID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IntegrationAccountMapName(),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"integration_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IntegrationAccountName(),
			},

			"content": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"map_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(logic.MapTypeXslt),
					string(logic.MapTypeXslt20),
					string(logic.MapTypeXslt30),
					string(logic.MapTypeLiquid),
				}, false),
			},

			"metadata": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func resourceLogicAppIntegrationAccountMapCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Logic.IntegrationAccountMapClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewIntegrationAccountMapID(subscriptionId, d.Get("resource_group_name").(string), d.Get("integration_account_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.IntegrationAccountName, id.MapName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_logic_app_integration_account_map", id.ID())
		}
	}

	parameters := logic.IntegrationAccountMap{
		IntegrationAccountMapProperties: &logic.IntegrationAccountMapProperties{
			MapType: logic.MapType(d.Get("map_type").(string)),
			Content: utils.String(d.Get("content").(string)),
		},
	}

	if parameters.IntegrationAccountMapProperties.MapType == logic.MapTypeLiquid {
		parameters.IntegrationAccountMapProperties.ContentType = utils.String("text/plain")
	} else {
		parameters.IntegrationAccountMapProperties.ContentType = utils.String("application/xml")
	}

	if v, ok := d.GetOk("metadata"); ok {
		metadata := v.(map[string]interface{})
		parameters.IntegrationAccountMapProperties.Metadata = &metadata
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.IntegrationAccountName, id.MapName, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogicAppIntegrationAccountMapRead(d, meta)
}

func resourceLogicAppIntegrationAccountMapRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountMapClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IntegrationAccountMapID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.IntegrationAccountName, id.MapName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.MapName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("integration_account_name", id.IntegrationAccountName)

	if props := resp.IntegrationAccountMapProperties; props != nil {
		d.Set("map_type", props.MapType)
		d.Set("content", d.Get("content").(string))

		if props.Metadata != nil {
			metadata := props.Metadata.(map[string]interface{})
			d.Set("metadata", metadata)
		}
	}

	return nil
}

func resourceLogicAppIntegrationAccountMapDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountMapClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IntegrationAccountMapID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.IntegrationAccountName, id.MapName); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
