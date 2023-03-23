package confidentialledger

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/confidentialledger/2022-05-13/confidentialledger"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/confidentialledger/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceConfidentialLedger() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceConfidentialLedgerRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			// Required
			"name": {
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
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"principal_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
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
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"pem_public_key": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"ledger_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"location": commonschema.LocationComputed(),

			"tags": commonschema.TagsDataSource(),

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

func dataSourceConfidentialLedgerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ConfidentialLedger.ConfidentialLedgerClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
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
