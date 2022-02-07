package logic

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/logic/mgmt/2019-05-01/logic"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceLogicAppIntegrationAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogicAppIntegrationAccountCreateUpdate,
		Read:   resourceLogicAppIntegrationAccountRead,
		Update: resourceLogicAppIntegrationAccountCreateUpdate,
		Delete: resourceLogicAppIntegrationAccountDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.IntegrationAccountID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IntegrationAccountName(),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(logic.IntegrationAccountSkuNameBasic),
					string(logic.IntegrationAccountSkuNameFree),
					string(logic.IntegrationAccountSkuNameStandard),
				}, false),
			},

			"integration_service_environment_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.IntegrationServiceEnvironmentID,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceLogicAppIntegrationAccountCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewIntegrationAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for present of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_logic_app_integration_account", id.ID())
		}
	}

	account := logic.IntegrationAccount{
		IntegrationAccountProperties: &logic.IntegrationAccountProperties{},
		Location:                     utils.String(location.Normalize(d.Get("location").(string))),
		Sku: &logic.IntegrationAccountSku{
			Name: logic.IntegrationAccountSkuName(d.Get("sku_name").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("integration_service_environment_id"); ok {
		account.IntegrationAccountProperties.IntegrationServiceEnvironment = &logic.ResourceReference{
			ID: utils.String(v.(string)),
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, account); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceLogicAppIntegrationAccountRead(d, meta)
}

func resourceLogicAppIntegrationAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IntegrationAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving Integration Account %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	d.Set("sku_name", string(resp.Sku.Name))

	if props := resp.IntegrationAccountProperties; props != nil {
		iseId := ""
		if props.IntegrationServiceEnvironment != nil && props.IntegrationServiceEnvironment.ID != nil {
			iseId = *props.IntegrationServiceEnvironment.ID
		}
		d.Set("integration_service_environment_id", iseId)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceLogicAppIntegrationAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IntegrationAccountID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.Name); err != nil {
		return fmt.Errorf("deleting Integration Account %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}
