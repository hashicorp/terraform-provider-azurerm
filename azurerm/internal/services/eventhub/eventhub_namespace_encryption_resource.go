package eventhub

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/eventhub/mgmt/2018-01-01-preview/eventhub"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceEventHubNamespaceEncryption() *schema.Resource {
	return &schema.Resource{
		Create: resourceEventHubNamespaceEncryptionCreateUpdate,
		Read:   resourceEventHubNamespaceEncryptionRead,
		Update: resourceEventHubNamespaceEncryptionCreateUpdate,
		Delete: resourceEventHubNamespaceEncryptionDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.NamespaceID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"namespace_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NamespaceID,
			},
			"key_vault_uri": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsURLWithHTTPS,
			},
			"key_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceEventHubNamespaceEncryptionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.NamespacesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM EventHub Namespace update for BYOK Encryption.")

	namespaceId, _ := parse.NamespaceID(d.Get("namespace_id").(string))

	if d.IsNewResource() {
		namespace, err := client.Get(ctx, namespaceId.ResourceGroup, namespaceId.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(namespace.Response) {
				return fmt.Errorf("Error checking for presence of existing EventHub Namespace %q (Resource Group %q): %s", namespaceId.Name, namespaceId.ResourceGroup, err)
			}
		}

		keyVersion := ""
		if v, ok := d.GetOk("key_version"); ok {
			keyVersion = v.(string)
		}

		kvProperties := make([]eventhub.KeyVaultProperties, 0)

		kvProperties = append(kvProperties, eventhub.KeyVaultProperties{
			KeyVaultURI: utils.String(d.Get("key_vault_uri").(string)),
			KeyName:     utils.String(d.Get("key_name").(string)),
			KeyVersion:  utils.String(keyVersion),
		})

		updatedNamespace := namespace
		updatedNamespace.EHNamespaceProperties.Encryption = &eventhub.Encryption{
			KeyVaultProperties: &kvProperties,
			KeySource:          eventhub.MicrosoftKeyVault,
		}

		future, err := client.CreateOrUpdate(ctx, namespaceId.ResourceGroup, namespaceId.Name, updatedNamespace)
		if err != nil {
			return fmt.Errorf("enabling BYOK encryption for EventHub Namespace: %+v", err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("enabling BYOK encryption for EventHub Namespace: %+v", err)
		}

		if d.IsNewResource() {
			d.SetId(namespaceId.ID())
		}

		return resourceEventHubNamespaceEncryptionRead(d, meta)
	}

	return fmt.Errorf("encryption settings for EventHub Namespace can not be modified after creation")
}

func resourceEventHubNamespaceEncryptionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.NamespacesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NamespaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] EventHub Namespace %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading EventHub Namespace %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("namespace_id", id.ID())

	if encryptionProps := resp.EHNamespaceProperties.Encryption; encryptionProps != nil {
		if encryptionProps.KeyVaultProperties != nil && len(*encryptionProps.KeyVaultProperties) != 0 {
			kvProps := *encryptionProps.KeyVaultProperties
			// Namespaces client strips trailing slash from the URI, we want it since KV client returns it
			d.Set("key_vault_uri", fmt.Sprintf("%s/", *kvProps[0].KeyVaultURI))
			d.Set("key_name", *kvProps[0].KeyName)
			d.Set("key_version", *kvProps[0].KeyVersion)
		}
	}

	return nil
}

func resourceEventHubNamespaceEncryptionDelete(d *schema.ResourceData, meta interface{}) error {
	return nil // can not be deleted
}
