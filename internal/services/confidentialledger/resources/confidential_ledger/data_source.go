// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package confidential_ledger

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/confidentialledger/2022-05-13/confidentialledger"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	confidentialledgerClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/confidentialledger/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/confidentialledger/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Gets information about an existing Confidential Ledger.",
		Read:        dataSourceConfidentialLedgerRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			// Required
			"name": {
				Description:  "The name of this Confidential Ledger.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ConfidentialLedgerName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			// Computed
			"azuread_based_service_principal": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"ledger_role_name": {
							Description: "The Ledger Role to grant this AzureAD Service Principal.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"principal_id": {
							Description: "The Principal ID of the AzureAD Service Principal.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"tenant_id": {
							Description: "The Tenant ID for this AzureAD Service Principal.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},

			"certificate_based_security_principal": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"ledger_role_name": {
							Description: "The Ledger Role to grant this Certificate Security Principal.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"pem_public_key": {
							Description: "The public key, in PEM format, of the certificate used by this identity to authenticate with the Confidential Ledger.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},

			"ledger_type": {
				Description: "The type of Confidential Ledger.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"location": commonschema.LocationComputed(),

			"tags": commonschema.TagsDataSource(),

			"identity_service_endpoint": {
				Description: "The Identity Service Endpoint for this Confidential Ledger.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},
			"ledger_endpoint": {
				Description: "The Endpoint for this Confidential Ledger.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceConfidentialLedgerRead(d *pluginsdk.ResourceData, meta any) error {
	client := confidentialledgerClient.From(meta).ConfidentialLedgerClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id := confidentialledger.NewLedgerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.LedgerGet(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.LedgerName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			if err := d.Set("azuread_based_service_principal", flattenAADBasedSecurityPrincipal(props.AadBasedSecurityPrincipals)); err != nil {
				return fmt.Errorf("setting `aad_based_security_principals`: %+v", err)
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
