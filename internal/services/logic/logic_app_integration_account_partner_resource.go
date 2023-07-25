// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logic

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountpartners"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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
			_, err := integrationaccountpartners.ParsePartnerID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IntegrationAccountPartnerName(),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

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

	id := integrationaccountpartners.NewPartnerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("integration_account_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_logic_app_integration_account_partner", id.ID())
		}
	}

	parameters := integrationaccountpartners.IntegrationAccountPartner{
		Properties: integrationaccountpartners.IntegrationAccountPartnerProperties{
			Content: integrationaccountpartners.PartnerContent{
				B2b: &integrationaccountpartners.B2BPartnerContent{
					BusinessIdentities: expandIntegrationAccountPartnerBusinessIdentity(d.Get("business_identity").(*pluginsdk.Set).List()),
				},
			},
			PartnerType: integrationaccountpartners.PartnerTypeBTwoB,
		},
	}

	if v, ok := d.GetOk("metadata"); ok {
		parameters.Properties.Metadata = &v
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogicAppIntegrationAccountPartnerRead(d, meta)
}

func resourceLogicAppIntegrationAccountPartnerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountPartnerClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationaccountpartners.ParsePartnerID(d.Id())
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

	d.Set("name", id.PartnerName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("integration_account_name", id.IntegrationAccountName)

	if model := resp.Model; model != nil {
		props := model.Properties
		if props.Content.B2b != nil && props.Content.B2b.BusinessIdentities != nil {
			if err := d.Set("business_identity", flattenIntegrationAccountPartnerBusinessIdentity(props.Content.B2b.BusinessIdentities)); err != nil {
				return fmt.Errorf("setting `business_identity`: %+v", err)
			}
		}

		if props.Metadata != nil {
			d.Set("metadata", props.Metadata)
		}
	}

	return nil
}

func resourceLogicAppIntegrationAccountPartnerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountPartnerClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationaccountpartners.ParsePartnerID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandIntegrationAccountPartnerBusinessIdentity(input []interface{}) *[]integrationaccountpartners.BusinessIdentity {
	results := make([]integrationaccountpartners.BusinessIdentity, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		results = append(results, integrationaccountpartners.BusinessIdentity{
			Qualifier: v["qualifier"].(string),
			Value:     v["value"].(string),
		})
	}

	return &results
}

func flattenIntegrationAccountPartnerBusinessIdentity(input *[]integrationaccountpartners.BusinessIdentity) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, map[string]interface{}{
			"qualifier": item.Qualifier,
			"value":     item.Value,
		})
	}

	return results
}
