package resource

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/confidentialledger/sdk/2021-05-13-preview/confidentialledger"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

// "fmt"
// "log"
// "time"

// "github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
// "github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
// "github.com/hashicorp/terraform-provider-azurerm/internal/clients"
// "github.com/hashicorp/terraform-provider-azurerm/internal/location"
// "github.com/hashicorp/terraform-provider-azurerm/internal/services/confidentialledger/sdk/2021-05-13-preview/confidentialledger"
// "github.com/hashicorp/terraform-provider-azurerm/internal/services/confidentialledger/validate"
// "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
// "github.com/hashicorp/terraform-provider-azurerm/internal/tags"
// "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
// "github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
// "github.com/hashicorp/terraform-provider-azurerm/utils"

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

		// Schema: map[string]*pluginsdk.Schema{
		// 	"name": {
		// 		Type:         pluginsdk.TypeString,
		// 		Required:     true,
		// 		ForceNew:     true,
		// 		ValidateFunc: validate.ConfidentialLedgerID,
		// 	},

		// 	"location": azure.SchemaLocation(),

		// 	"identity": a{}.Schema(),

		// 	// the API changed and now returns the rg in lowercase
		// 	// revert when https://github.com/Azure/azure-sdk-for-go/issues/6606 is fixed
		// 	"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

		// 	"sku": {
		// 		Type:     pluginsdk.TypeString,
		// 		Optional: true,
		// 		Default:  "free",
		// 		ValidateFunc: validation.StringInSlice([]string{
		// 			"free",
		// 			"standard",
		// 		}, false),
		// 	},

		// 	"endpoint": {
		// 		Type:     pluginsdk.TypeString,
		// 		Computed: true,
		// 	},

		// 	"primary_read_key": {
		// 		Type:     pluginsdk.TypeList,
		// 		Computed: true,
		// 		Elem: &pluginsdk.Resource{
		// 			Schema: map[string]*pluginsdk.Schema{
		// 				"id": {
		// 					Type:      pluginsdk.TypeString,
		// 					Computed:  true,
		// 					Sensitive: true,
		// 				},
		// 				"secret": {
		// 					Type:      pluginsdk.TypeString,
		// 					Computed:  true,
		// 					Sensitive: true,
		// 				},
		// 				"connection_string": {
		// 					Type:      pluginsdk.TypeString,
		// 					Computed:  true,
		// 					Sensitive: true,
		// 				},
		// 			},
		// 		},
		// 	},

		// 	"secondary_read_key": {
		// 		Type:     pluginsdk.TypeList,
		// 		Computed: true,
		// 		Elem: &pluginsdk.Resource{
		// 			Schema: map[string]*pluginsdk.Schema{
		// 				"id": {
		// 					Type:      pluginsdk.TypeString,
		// 					Computed:  true,
		// 					Sensitive: true,
		// 				},
		// 				"secret": {
		// 					Type:      pluginsdk.TypeString,
		// 					Computed:  true,
		// 					Sensitive: true,
		// 				},
		// 				"connection_string": {
		// 					Type:      pluginsdk.TypeString,
		// 					Computed:  true,
		// 					Sensitive: true,
		// 				},
		// 			},
		// 		},
		// 	},

		// 	"primary_write_key": {
		// 		Type:     pluginsdk.TypeList,
		// 		Computed: true,
		// 		Elem: &pluginsdk.Resource{
		// 			Schema: map[string]*pluginsdk.Schema{
		// 				"id": {
		// 					Type:      pluginsdk.TypeString,
		// 					Computed:  true,
		// 					Sensitive: true,
		// 				},
		// 				"secret": {
		// 					Type:      pluginsdk.TypeString,
		// 					Computed:  true,
		// 					Sensitive: true,
		// 				},
		// 				"connection_string": {
		// 					Type:      pluginsdk.TypeString,
		// 					Computed:  true,
		// 					Sensitive: true,
		// 				},
		// 			},
		// 		},
		// 	},

		// 	"secondary_write_key": {
		// 		Type:     pluginsdk.TypeList,
		// 		Computed: true,
		// 		Elem: &pluginsdk.Resource{
		// 			Schema: map[string]*pluginsdk.Schema{
		// 				"id": {
		// 					Type:      pluginsdk.TypeString,
		// 					Computed:  true,
		// 					Sensitive: true,
		// 				},
		// 				"secret": {
		// 					Type:      pluginsdk.TypeString,
		// 					Computed:  true,
		// 					Sensitive: true,
		// 				},
		// 				"connection_string": {
		// 					Type:      pluginsdk.TypeString,
		// 					Computed:  true,
		// 					Sensitive: true,
		// 				},
		// 			},
		// 		},
		// 	},

		// 	"tags": commonschema.Tags(),
		// },
	}
}

func resourceConfidentialLedgerCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ConfidentialLedger.ConfidentialLedgereClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM App Configuration creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	resourceId := confidentialledger.NewLedgerID(subscriptionId, resourceGroup, name)
	existing, err := client.LedgerGet(ctx, resourceId)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", resourceId, err)
		}
	}
	// if !response.WasNotFound(existing.HttpResponse) {
	// 	return tf.ImportAsExistsError("azurerm_app_configuration", resourceId.ID())
	// }

	// parameters := configurationstores.ConfigurationStore{
	// 	Location: azure.NormalizeLocation(d.Get("location").(string)),
	// 	Sku: configurationstores.Sku{
	// 		Name: d.Get("sku").(string),
	// 	},
	// 	Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	// }

	// identity, err := expandAppConfigurationIdentity(d.Get("identity").([]interface{}))
	// if err != nil {
	// 	return fmt.Errorf("expanding `identity`: %+v", err)
	// }
	// parameters.Identity = identity

	// if err := client.CreateThenPoll(ctx, resourceId, parameters); err != nil {
	// 	return fmt.Errorf("creating %s: %+v", resourceId, err)
	// }

	// d.SetId(resourceId.ID())
	// return resourceAppConfigurationRead(d, meta)

	return nil
}
