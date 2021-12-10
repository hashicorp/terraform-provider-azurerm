package logic

import (
	"encoding/json"
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

func resourceLogicAppIntegrationAccountAgreement() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogicAppIntegrationAccountAgreementCreateUpdate,
		Read:   resourceLogicAppIntegrationAccountAgreementRead,
		Update: resourceLogicAppIntegrationAccountAgreementCreateUpdate,
		Delete: resourceLogicAppIntegrationAccountAgreementDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.IntegrationAccountAgreementID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IntegrationAccountAgreementName(),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"integration_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IntegrationAccountName(),
			},

			"agreement_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(logic.AgreementTypeAS2),
					string(logic.AgreementTypeX12),
					string(logic.AgreementTypeEdifact),
				}, false),
			},

			"content": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
			},

			"guest_identity": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
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

			"guest_partner_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.IntegrationAccountPartnerName(),
			},

			"host_identity": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
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

			"host_partner_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.IntegrationAccountPartnerName(),
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

func resourceLogicAppIntegrationAccountAgreementCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Logic.IntegrationAccountAgreementClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewIntegrationAccountAgreementID(subscriptionId, d.Get("resource_group_name").(string), d.Get("integration_account_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.IntegrationAccountName, id.AgreementName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_logic_app_integration_account_agreement", id.ID())
		}
	}

	agreementContent := logic.AgreementContent{}
	content := d.Get("content").(string)
	if err := json.Unmarshal([]byte(content), &agreementContent); err != nil {
		return fmt.Errorf("parsing JSON: %+v", err)
	}

	parameters := logic.IntegrationAccountAgreement{
		IntegrationAccountAgreementProperties: &logic.IntegrationAccountAgreementProperties{
			AgreementType: logic.AgreementType(d.Get("agreement_type").(string)),
			GuestIdentity: expandIntegrationAccountAgreementBusinessIdentity(d.Get("guest_identity").([]interface{})),
			GuestPartner:  utils.String(d.Get("guest_partner_name").(string)),
			HostIdentity:  expandIntegrationAccountAgreementBusinessIdentity(d.Get("host_identity").([]interface{})),
			HostPartner:   utils.String(d.Get("host_partner_name").(string)),
			Content:       &agreementContent,
		},
	}

	if v, ok := d.GetOk("metadata"); ok {
		metadata := v.(map[string]interface{})
		parameters.IntegrationAccountAgreementProperties.Metadata = &metadata
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.IntegrationAccountName, id.AgreementName, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogicAppIntegrationAccountAgreementRead(d, meta)
}

func resourceLogicAppIntegrationAccountAgreementRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountAgreementClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IntegrationAccountAgreementID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.IntegrationAccountName, id.AgreementName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.AgreementName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("integration_account_name", id.IntegrationAccountName)

	if props := resp.IntegrationAccountAgreementProperties; props != nil {
		d.Set("agreement_type", props.AgreementType)
		d.Set("guest_partner_name", props.GuestPartner)
		d.Set("host_partner_name", props.HostPartner)

		if props.Content != nil {
			content, err := json.Marshal(props.Content)
			if err != nil {
				return err
			}
			d.Set("content", string(content))
		}

		if err := d.Set("guest_identity", flattenIntegrationAccountAgreementBusinessIdentity(props.GuestIdentity)); err != nil {
			return fmt.Errorf("setting `guest_identity`: %+v", err)
		}

		if err := d.Set("host_identity", flattenIntegrationAccountAgreementBusinessIdentity(props.HostIdentity)); err != nil {
			return fmt.Errorf("setting `host_identity`: %+v", err)
		}

		if props.Metadata != nil {
			metadata := props.Metadata.(map[string]interface{})
			d.Set("metadata", metadata)
		}
	}
	return nil
}

func resourceLogicAppIntegrationAccountAgreementDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountAgreementClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IntegrationAccountAgreementID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.IntegrationAccountName, id.AgreementName); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandIntegrationAccountAgreementBusinessIdentity(input []interface{}) *logic.BusinessIdentity {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	return &logic.BusinessIdentity{
		Qualifier: utils.String(v["qualifier"].(string)),
		Value:     utils.String(v["value"].(string)),
	}
}

func flattenIntegrationAccountAgreementBusinessIdentity(input *logic.BusinessIdentity) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var qualifier string
	if input.Qualifier != nil {
		qualifier = *input.Qualifier
	}

	var value string
	if input.Value != nil {
		value = *input.Value
	}

	return []interface{}{
		map[string]interface{}{
			"qualifier": qualifier,
			"value":     value,
		},
	}
}
