package confidentialledger

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/confidentialledger/sdk/2021-05-13-preview/confidentialledger"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/confidentialledger/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceConfidentialLedger() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceConfidentialLedgerCreate,
		Read:   resourceConfidentialLedgerRead,
		// Update: resourceConfidentialLedgerUpdate,
		// Delete: resourceConfidentialLedgerDelete,

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
			"aad_based_security_principals": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"principal_id": {
							Type:      pluginsdk.TypeString,
							Sensitive: true,
						},
						"tenant_id": {
							Type:      pluginsdk.TypeString,
							Sensitive: true,
						},
						"ledger_role_name": {
							Type: pluginsdk.TypeString,
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
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"cert": {
							Type:      pluginsdk.TypeString,
							Sensitive: true,
						},
						"ledger_role_name": {
							Type: pluginsdk.TypeString,
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
				ValidateFunc: validate.ConfidentialLedgerID,
			},

			"ledger_type": {
				Type: pluginsdk.TypeString,
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

func resourceConfidentialLedgerCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ConfidentialLedger.ConfidentialLedgereClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM App Configuration creation.")

	ledgerName := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	resourceId := confidentialledger.NewLedgerID(subscriptionId, resourceGroup, ledgerName)
	existing, err := client.LedgerGet(ctx, resourceId)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", resourceId.ID(), err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_confidential_ledger", resourceId.ID())
	}

	aadBasedUsers := expandConfidentialLedgerAADBasedSecurityPrincipal(d.Get("aad_based_security_principals").([]interface{}))
	certBasedUsers := expandConfidentialLedgerCertBasedSecurityPrincipal(d.Get("cert_based_security_principals").([]interface{}))

	// TODO: Insert ledger properties..?
	parameters := confidentialledger.ConfidentialLedger{
		Location: azure.NormalizeLocation(d.Get("location").(string)),
		Name:     &resourceId.LedgerName,
		Properties: &confidentialledger.LedgerProperties{
			AadBasedSecurityPrincipals:  aadBasedUsers,
			CertBasedSecurityPrincipals: certBasedUsers,
			LedgerType:                  d.Get("ledger_type").(*confidentialledger.LedgerType),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.LedgerCreateThenPoll(ctx, resourceId, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", resourceId, err)
	}

	d.SetId(resourceId.ID())
	return resourceConfidentialLedgerRead(d, meta)
}

func resourceConfidentialLedgerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ConfidentialLedger.ConfidentialLedgereClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := confidentialledger.ParseLedgerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.LedgerGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.LedgerName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("ledger_type", model.Properties.LedgerType)
		d.Set("location", location.Normalize(model.Location))
		d.Set("tags", model.Tags)

		aadBasedUsers, err := flattenConfidentialLedgerAADBasedSecurityPrincipal(model.Properties.AadBasedSecurityPrincipals)
		if err != nil {
			return fmt.Errorf("Error retrieving AAD-based users for %s: %+v", *id, err)
		}

		certBasedUsers, err := flattenConfidentialLedgerCertBasedSecurityPrincipal(model.Properties.CertBasedSecurityPrincipals)
		if err != nil {
			return fmt.Errorf("Error retrieving cert-based users for %s: %+v", *id, err)
		}

		d.Set("aad_based_security_principals", aadBasedUsers)
		d.Set("cert_based_security_principals", certBasedUsers)
	}

	return nil
}
