// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package confidentialledger

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/confidentialledger/2022-05-13/confidentialledger"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/confidentialledger/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceConfidentialLedger() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceConfidentialLedgerCreate,
		Read:   resourceConfidentialLedgerRead,
		Update: resourceConfidentialLedgerUpdate,
		Delete: resourceConfidentialLedgerDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := confidentialledger.ParseLedgerID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			// Required
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ConfidentialLedgerName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"ledger_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(confidentialledger.LedgerTypePrivate),
					string(confidentialledger.LedgerTypePublic),
				}, false),
			},

			"azuread_based_service_principal": {
				// this is Required since if none are specified then the calling SP gets added
				Type:     pluginsdk.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"ledger_role_name": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(confidentialledger.LedgerRoleNameAdministrator),
								string(confidentialledger.LedgerRoleNameContributor),
								string(confidentialledger.LedgerRoleNameReader),
							}, false),
						},
						"principal_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
						"tenant_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
					},
				},
			},

			// Optional
			"certificate_based_security_principal": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"ledger_role_name": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(confidentialledger.LedgerRoleNameAdministrator),
								string(confidentialledger.LedgerRoleNameContributor),
								string(confidentialledger.LedgerRoleNameReader),
							}, false),
						},
						"pem_public_key": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"tags": commonschema.Tags(),

			"identity_service_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"ledger_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceConfidentialLedgerCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ConfidentialLedger.ConfidentialLedgerClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := confidentialledger.NewLedgerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.LedgerGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_confidential_ledger", id.ID())
	}

	aadBasedUsers := expandAADBasedSecurityPrincipal(d.Get("azuread_based_service_principal").([]interface{}))
	certBasedUsers := expandCertBasedSecurityPrincipal(d.Get("certificate_based_security_principal").([]interface{}))
	ledgerType := confidentialledger.LedgerType(d.Get("ledger_type").(string))
	location := location.Normalize(d.Get("location").(string))
	parameters := confidentialledger.ConfidentialLedger{
		Location: utils.String(location),
		Properties: &confidentialledger.LedgerProperties{
			AadBasedSecurityPrincipals:  aadBasedUsers,
			CertBasedSecurityPrincipals: certBasedUsers,
			LedgerType:                  &ledgerType,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.LedgerCreateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("error creating %s: %+v", id.ID(), err)
	}

	d.SetId(id.ID())
	return resourceConfidentialLedgerRead(d, meta)
}

func resourceConfidentialLedgerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ConfidentialLedger.ConfidentialLedgerClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := confidentialledger.ParseLedgerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.LedgerGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.LedgerName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			if err := d.Set("azuread_based_service_principal", flattenAADBasedSecurityPrincipal(props.AadBasedSecurityPrincipals)); err != nil {
				return fmt.Errorf("setting `azuread_based_service_principal`: %+v", err)
			}
			if err := d.Set("certificate_based_security_principal", flattenCertBasedSecurityPrincipal(props.CertBasedSecurityPrincipals)); err != nil {
				return fmt.Errorf("setting `certificate_based_security_principal`: %+v", err)
			}

			ledgerType := ""
			if props.LedgerType != nil {
				ledgerType = string(*props.LedgerType)
			}
			d.Set("ledger_type", ledgerType)

			d.Set("ledger_endpoint", props.LedgerUri)
			d.Set("identity_service_endpoint", props.IdentityServiceUri)
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}

func resourceConfidentialLedgerUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ConfidentialLedger.ConfidentialLedgerClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := confidentialledger.ParseLedgerID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.LedgerGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving existing %s: %+v", *id, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("retrieving existing %s: model was nil", *id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving existing %s: `properties` was nil", *id)
	}

	ledger := confidentialledger.ConfidentialLedger{
		Location: existing.Model.Location,
		Properties: &confidentialledger.LedgerProperties{
			AadBasedSecurityPrincipals:  existing.Model.Properties.AadBasedSecurityPrincipals,
			CertBasedSecurityPrincipals: existing.Model.Properties.CertBasedSecurityPrincipals,
			LedgerType:                  existing.Model.Properties.LedgerType,
		},
		Tags: existing.Model.Tags,
	}

	if d.HasChange("azuread_based_service_principal") {
		aadBasedUsers := expandAADBasedSecurityPrincipal(d.Get("azuread_based_service_principal").([]interface{}))
		ledger.Properties.AadBasedSecurityPrincipals = aadBasedUsers
	}

	if d.HasChange("certificate_based_security_principal") {
		certBasedUsers := expandCertBasedSecurityPrincipal(d.Get("certificate_based_security_principal").([]interface{}))
		ledger.Properties.CertBasedSecurityPrincipals = certBasedUsers
	}

	if d.HasChange("tags") {
		ledger.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.LedgerUpdateThenPoll(ctx, *id, ledger); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceConfidentialLedgerRead(d, meta)
}

func resourceConfidentialLedgerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ConfidentialLedger.ConfidentialLedgerClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := confidentialledger.ParseLedgerID(d.Id())
	if err != nil {
		return err
	}

	if err := client.LedgerDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandAADBasedSecurityPrincipal(input []interface{}) *[]confidentialledger.AADBasedSecurityPrincipal {
	output := make([]confidentialledger.AADBasedSecurityPrincipal, 0)

	for _, item := range input {
		v := item.(map[string]interface{})
		ledgerRoleName := confidentialledger.LedgerRoleName(v["ledger_role_name"].(string))
		principalId := v["principal_id"].(string)
		tenantId := v["tenant_id"].(string)

		result := confidentialledger.AADBasedSecurityPrincipal{
			LedgerRoleName: &ledgerRoleName,
			PrincipalId:    utils.String(principalId),
			TenantId:       utils.String(tenantId),
		}

		output = append(output, result)
	}

	return &output
}

func expandCertBasedSecurityPrincipal(input []interface{}) *[]confidentialledger.CertBasedSecurityPrincipal {
	output := make([]confidentialledger.CertBasedSecurityPrincipal, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		ledgerRoleName := confidentialledger.LedgerRoleName(v["ledger_role_name"].(string))
		output = append(output, confidentialledger.CertBasedSecurityPrincipal{
			Cert:           utils.String(v["pem_public_key"].(string)),
			LedgerRoleName: &ledgerRoleName,
		})
	}

	return &output
}

func flattenAADBasedSecurityPrincipal(input *[]confidentialledger.AADBasedSecurityPrincipal) []interface{} {
	output := make([]interface{}, 0)
	if input == nil {
		return output
	}

	for _, item := range *input {
		ledgerRoleName := ""
		if item.LedgerRoleName != nil {
			ledgerRoleName = string(*item.LedgerRoleName)
		}

		principalId := ""
		if item.PrincipalId != nil {
			principalId = *item.PrincipalId
		}

		tenantId := ""
		if item.TenantId != nil {
			tenantId = *item.TenantId
		}

		output = append(output, map[string]interface{}{
			"ledger_role_name": ledgerRoleName,
			"principal_id":     principalId,
			"tenant_id":        tenantId,
		})
	}

	return output
}

func flattenCertBasedSecurityPrincipal(input *[]confidentialledger.CertBasedSecurityPrincipal) []interface{} {
	output := make([]interface{}, 0)
	if input == nil {
		return []interface{}{}
	}

	for _, item := range *input {
		pemPublicKey := ""
		if item.Cert != nil {
			pemPublicKey = *item.Cert
		}

		ledgerRoleName := ""
		if item.LedgerRoleName != nil {
			ledgerRoleName = string(*item.LedgerRoleName)
		}

		output = append(output, map[string]interface{}{
			"ledger_role_name": ledgerRoleName,
			"pem_public_key":   pemPublicKey,
		})
	}

	return output
}
