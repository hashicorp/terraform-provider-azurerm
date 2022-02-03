package resource

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
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/confidentialledger/sdk/2021-05-13-preview/confidentialledger"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/confidentialledger/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceConfidentialLedger() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceConfidentialLedgerCreate,
		// Read:   resourceConfidentialLedgerRead,
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
			"aad_based_security_principals": {},

			"cert_based_security_principals": {},

			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ConfidentialLedgerID,
			},

			"ledger_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "free",
				ValidateFunc: validation.StringInSlice([]string{
					"Public",
					"Private",
				}, false),
			},

			"location": azure.SchemaLocation(),

			// the API changed and now returns the rg in lowercase
			// revert when https://github.com/Azure/azure-sdk-for-go/issues/6606 is fixed
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

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

	// TODO: Insert ledger properties..?
	parameters := confidentialledger.ConfidentialLedger{
		Location: azure.NormalizeLocation(d.Get("location").(string)),
		Name:     &resourceId.LedgerName,
		Properties: &confidentialledger.LedgerProperties{
			AadBasedSecurityPrincipals:  nil, // *[]AADBasedSecurityPrincipal
			CertBasedSecurityPrincipals: nil, // *[]CertBasedSecurityPrincipal
			LedgerType:                  nil, // *LedgerType
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.LedgerCreateThenPoll(ctx, resourceId, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", resourceId, err)
	}

	d.SetId(resourceId.ID())
	// TODO
	// return resourceAppConfigurationRead(d, meta)
	return nil
}
