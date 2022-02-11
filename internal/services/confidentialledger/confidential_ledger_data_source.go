package confidentialledger

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/confidentialledger/sdk/2021-05-13-preview/confidentialledger"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/confidentialledger/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceConfidentialLedger() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceConfidentialLedgerRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		// This should match the Schema in resourceConfidentialLedger
		Schema: map[string]*pluginsdk.Schema{
			"aad_based_security_principals": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"principal_id": {
							Type:      pluginsdk.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"tenant_id": {
							Type:      pluginsdk.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"ledger_role_name": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Administrator",
								"Contributor",
								"Reader",
							}, false),
						},
					},
				},
			},

			"cert_based_security_principals": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"cert": {
							Type:      pluginsdk.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"ledger_role_name": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Administrator",
								"Contributor",
								"Reader",
							}, false),
						},
					},
				},
			},

			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ConfidentialLedgerName,
			},

			"ledger_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Public",
					"Private",
				}, false),
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"tags": commonschema.Tags(),
		},
	}
}

func dataSourceConfidentialLedgerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	// Warning!! This is not a duplicate of resourceConfidentialLedgerRead.

	client := meta.(*clients.Client).ConfidentialLedger.ConfidentialLedgereClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	// One difference - the request does not come in with an id.
	resourceId := confidentialledger.NewLedgerID(subscriptionId, resourceGroup, name)

	resp, err := client.LedgerGet(ctx, resourceId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", resourceId.ID())
		}

		return fmt.Errorf("error retrieving %s: %+v", resourceId.ID(), err)
	}

	// Another difference - the id must be set.
	d.SetId(resourceId.ID())

	d.Set("name", resourceId.LedgerName)
	d.Set("resource_group_name", resourceId.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("ledger_type", string(*model.Properties.LedgerType))
		d.Set("location", location.Normalize(model.Location))
		d.Set("tags", model.Tags)

		aadBasedUsers, err := flattenConfidentialLedgerAADBasedSecurityPrincipal(model.Properties.AadBasedSecurityPrincipals)
		if err != nil {
			return fmt.Errorf("error retrieving AAD-based users for %s: %+v", resourceId.ID(), err)
		}

		certBasedUsers, err := flattenConfidentialLedgerCertBasedSecurityPrincipal(model.Properties.CertBasedSecurityPrincipals)
		if err != nil {
			return fmt.Errorf("error retrieving cert-based users for %s: %+v", resourceId.ID(), err)
		}

		d.Set("aad_based_security_principals", aadBasedUsers)
		d.Set("cert_based_security_principals", certBasedUsers)
	}

	return nil
}
