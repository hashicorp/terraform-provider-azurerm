// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logic

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountagreements"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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
			_, err := integrationaccountagreements.ParseAgreementID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IntegrationAccountAgreementName(),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

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
					string(integrationaccountagreements.AgreementTypeASTwo),
					string(integrationaccountagreements.AgreementTypeXOneTwo),
					string(integrationaccountagreements.AgreementTypeEdifact),
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

	id := integrationaccountagreements.NewAgreementID(subscriptionId, d.Get("resource_group_name").(string), d.Get("integration_account_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_logic_app_integration_account_agreement", id.ID())
		}
	}

	agreementContent := integrationaccountagreements.AgreementContent{}
	content := d.Get("content").(string)
	if err := json.Unmarshal([]byte(content), &agreementContent); err != nil {
		return fmt.Errorf("parsing JSON: %+v", err)
	}

	parameters := integrationaccountagreements.IntegrationAccountAgreement{
		Properties: integrationaccountagreements.IntegrationAccountAgreementProperties{
			AgreementType: integrationaccountagreements.AgreementType(d.Get("agreement_type").(string)),
			GuestIdentity: expandIntegrationAccountAgreementBusinessIdentity(d.Get("guest_identity").([]interface{})),
			GuestPartner:  d.Get("guest_partner_name").(string),
			HostIdentity:  expandIntegrationAccountAgreementBusinessIdentity(d.Get("host_identity").([]interface{})),
			HostPartner:   d.Get("host_partner_name").(string),
			Content:       agreementContent,
		},
	}

	if v, ok := d.GetOk("metadata"); ok {
		parameters.Properties.Metadata = &v
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogicAppIntegrationAccountAgreementRead(d, meta)
}

func resourceLogicAppIntegrationAccountAgreementRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountAgreementClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationaccountagreements.ParseAgreementID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.AgreementName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("integration_account_name", id.IntegrationAccountName)

	if model := resp.Model; model != nil {
		props := model.Properties

		d.Set("agreement_type", props.AgreementType)
		d.Set("guest_partner_name", props.GuestPartner)
		d.Set("host_partner_name", props.HostPartner)

		content, err := json.Marshal(props.Content)
		if err != nil {
			return err
		}
		d.Set("content", string(content))

		if err := d.Set("guest_identity", flattenIntegrationAccountAgreementBusinessIdentity(props.GuestIdentity)); err != nil {
			return fmt.Errorf("setting `guest_identity`: %+v", err)
		}

		if err := d.Set("host_identity", flattenIntegrationAccountAgreementBusinessIdentity(props.HostIdentity)); err != nil {
			return fmt.Errorf("setting `host_identity`: %+v", err)
		}

		if props.Metadata != nil {
			d.Set("metadata", props.Metadata)
		}

	}

	return nil
}

func resourceLogicAppIntegrationAccountAgreementDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountAgreementClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationaccountagreements.ParseAgreementID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandIntegrationAccountAgreementBusinessIdentity(input []interface{}) integrationaccountagreements.BusinessIdentity {
	if len(input) == 0 {
		return integrationaccountagreements.BusinessIdentity{}
	}
	v := input[0].(map[string]interface{})

	return integrationaccountagreements.BusinessIdentity{
		Qualifier: v["qualifier"].(string),
		Value:     v["value"].(string),
	}
}

func flattenIntegrationAccountAgreementBusinessIdentity(input integrationaccountagreements.BusinessIdentity) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"qualifier": input.Qualifier,
			"value":     input.Value,
		},
	}
}
