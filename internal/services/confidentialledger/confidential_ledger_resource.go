package resource

import (
	"fmt"
	"log"
	"strconv"
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

	numAadBasedUsers, err := strconv.Atoi(d.Get("aad_based_security_principals.#").(string))
	if err != nil {
		return fmt.Errorf("Could not convert 'aad_based_security_principals.#' = '%s' to integer", d.Get("aad_based_security_principals.#"))
	}

	aadBasedUsers := make([]confidentialledger.AADBasedSecurityPrincipal, numAadBasedUsers)
	for i := 0; i < numAadBasedUsers; i++ {
		tempData := d.Get(fmt.Sprintf("aad_based_security_principals.%d", i)).(pluginsdk.ResourceData)
		ledgerRoleName := tempData.Get("ledger_role_name").(confidentialledger.LedgerRoleName)

		aadBasedUsers = append(aadBasedUsers, confidentialledger.AADBasedSecurityPrincipal{
			LedgerRoleName: &ledgerRoleName,
			PrincipalId:    nil,
			TenantId:       nil,
		})
	}

	numCertBasedUsers, err := strconv.Atoi(d.Get("cert_based_security_principals.#").(string))
	if err != nil {
		return fmt.Errorf("Could not convert 'aad_based_security_principals.#' = '%s' to integer", d.Get("aad_based_security_principals.#"))
	}

	certBasedUsers := make([]confidentialledger.CertBasedSecurityPrincipal, numCertBasedUsers)
	for i := 0; i < numAadBasedUsers; i++ {
		certBasedUsers = append(certBasedUsers, confidentialledger.CertBasedSecurityPrincipal{
			Cert:           nil,
			LedgerRoleName: nil,
		})
	}

	// TODO: Insert ledger properties..?
	parameters := confidentialledger.ConfidentialLedger{
		Location: azure.NormalizeLocation(d.Get("location").(string)),
		Name:     &resourceId.LedgerName,
		Properties: &confidentialledger.LedgerProperties{
			AadBasedSecurityPrincipals:  &aadBasedUsers,
			CertBasedSecurityPrincipals: &certBasedUsers,
			LedgerType:                  d.Get("ledger_type").(*confidentialledger.LedgerType),
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
