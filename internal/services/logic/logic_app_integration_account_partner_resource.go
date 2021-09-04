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

func resourceLogicAppIntegrationAccountPartner() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogicAppIntegrationAccountPartnerCreateUpdate,
		Read:   resourceLogicAppIntegrationAccountPartnerRead,
		Update: resourceLogicAppIntegrationAccountPartnerCreateUpdate,
		Delete: resourceLogicAppIntegrationAccountPartnerDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.IntegrationAccountPartnerID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IntegrationAccountPartnerName(),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"integration_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IntegrationAccountName(),
			},

			"business_identity": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"qualifier": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.IntegrationAccountPartnerBusinessIdentityQualifier(),
						},

						"value": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.IntegrationAccountPartnerBusinessIdentityValue(),
						},
					},
				},
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

func resourceLogicAppIntegrationAccountPartnerCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Logic.IntegrationAccountPartnerClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewIntegrationAccountPartnerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("integration_account_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.IntegrationAccountName, id.PartnerName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_logic_app_integration_account_partner", id.ID())
		}
	}

	parameters := logic.IntegrationAccountPartner{
		IntegrationAccountPartnerProperties: &logic.IntegrationAccountPartnerProperties{
			Content: &logic.PartnerContent{
				B2b: &logic.B2BPartnerContent{
					BusinessIdentities: expandIntegrationAccountPartnerBusinessIdentity(d.Get("business_identity").(*pluginsdk.Set).List()),
				},
			},
			PartnerType: logic.PartnerTypeB2B,
		},
	}

	if v, ok := d.GetOk("metadata"); ok {
		metadata, _ := pluginsdk.ExpandJsonFromString(v.(string))
		parameters.IntegrationAccountPartnerProperties.Metadata = metadata
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.IntegrationAccountName, id.PartnerName, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogicAppIntegrationAccountPartnerRead(d, meta)
}

func resourceLogicAppIntegrationAccountPartnerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountPartnerClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IntegrationAccountPartnerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.IntegrationAccountName, id.PartnerName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.PartnerName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("integration_account_name", id.IntegrationAccountName)

	if props := resp.IntegrationAccountPartnerProperties; props != nil {
		if props.Content != nil && props.Content.B2b != nil && props.Content.B2b.BusinessIdentities != nil {
			if err := d.Set("business_identity", flattenIntegrationAccountPartnerBusinessIdentity(props.Content.B2b.BusinessIdentities)); err != nil {
				return fmt.Errorf("setting `business_identity`: %+v", err)
			}
		}

		if props.Metadata != nil {
			metadataValue := props.Metadata.(map[string]interface{})
			metadataStr, _ := pluginsdk.FlattenJsonToString(metadataValue)
			d.Set("metadata", metadataStr)
		}
	}

	return nil
}

func resourceLogicAppIntegrationAccountPartnerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountPartnerClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IntegrationAccountPartnerID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.IntegrationAccountName, id.PartnerName); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandIntegrationAccountPartnerBusinessIdentity(input []interface{}) *[]logic.BusinessIdentity {
	results := make([]logic.BusinessIdentity, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		results = append(results, logic.BusinessIdentity{
			Qualifier: utils.String(v["qualifier"].(string)),
			Value:     utils.String(v["value"].(string)),
		})
	}

	return &results
}

func flattenIntegrationAccountPartnerBusinessIdentity(input *[]logic.BusinessIdentity) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var qualifier string
		if item.Qualifier != nil {
			qualifier = *item.Qualifier
		}

		var value string
		if item.Value != nil {
			value = *item.Value
		}

		results = append(results, map[string]interface{}{
			"qualifier": qualifier,
			"value":     value,
		})
	}

	return results
}
