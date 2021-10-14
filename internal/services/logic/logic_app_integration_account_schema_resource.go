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

func resourceLogicAppIntegrationAccountSchema() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogicAppIntegrationAccountSchemaCreateUpdate,
		Read:   resourceLogicAppIntegrationAccountSchemaRead,
		Update: resourceLogicAppIntegrationAccountSchemaCreateUpdate,
		Delete: resourceLogicAppIntegrationAccountSchemaDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.IntegrationAccountSchemaID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IntegrationAccountSchemaName(),
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

			"file_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.IntegrationAccountSchemaFileName(),
			},

			"metadata": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
			},
		},
	}
}

func resourceLogicAppIntegrationAccountSchemaCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Logic.IntegrationAccountSchemaClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewIntegrationAccountSchemaID(subscriptionId, d.Get("resource_group_name").(string), d.Get("integration_account_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.IntegrationAccountName, id.SchemaName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_logic_app_integration_account_schema", id.ID())
		}
	}

	parameters := logic.IntegrationAccountSchema{
		IntegrationAccountSchemaProperties: &logic.IntegrationAccountSchemaProperties{
			SchemaType:  logic.SchemaTypeXML,
			Content:     utils.String(d.Get("content").(string)),
			ContentType: utils.String("application/xml"),
		},
	}

	if v, ok := d.GetOk("file_name"); ok {
		parameters.IntegrationAccountSchemaProperties.FileName = utils.String(v.(string))
	}

	if v, ok := d.GetOk("metadata"); ok {
		metadata, _ := pluginsdk.ExpandJsonFromString(v.(string))
		parameters.IntegrationAccountSchemaProperties.Metadata = metadata
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.IntegrationAccountName, id.SchemaName, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogicAppIntegrationAccountSchemaRead(d, meta)
}

func resourceLogicAppIntegrationAccountSchemaRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountSchemaClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IntegrationAccountSchemaID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.IntegrationAccountName, id.SchemaName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.SchemaName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("integration_account_name", id.IntegrationAccountName)

	if props := resp.IntegrationAccountSchemaProperties; props != nil {
		d.Set("content", d.Get("content").(string))
		d.Set("file_name", d.Get("file_name").(string))

		if props.Metadata != nil {
			metadataValue := props.Metadata.(map[string]interface{})
			metadataStr, _ := pluginsdk.FlattenJsonToString(metadataValue)
			d.Set("metadata", metadataStr)
		}
	}

	return nil
}

func resourceLogicAppIntegrationAccountSchemaDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountSchemaClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IntegrationAccountSchemaID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.IntegrationAccountName, id.SchemaName); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
