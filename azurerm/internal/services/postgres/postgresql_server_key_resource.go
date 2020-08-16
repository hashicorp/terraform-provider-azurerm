package postgres

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2020-01-01/postgresql"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPostgreSQLServerKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPostgreSQLServerKeyCreate,
		Read:   resourceArmPostgreSQLServerKeyRead,
		Delete: resourceArmPostgreSQLServerKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			"resource_group_name": azure.SchemaResourceGroupName(),

			"server_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.PostgresServerServerName,
			},

			"key_url": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.PostgresServerServerKeyUrl,
			},

			"key_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceArmPostgreSQLServerKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ServerKeysClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resGroup := d.Get("resource_group_name").(string)
	serverName := d.Get("server_name").(string)
	keyURL := d.Get("key_url").(string)
	name := getServerKeyNameFromUri(keyURL)
	keyType := d.Get("key_type").(string)

	log.Printf("[INFO] Creating key name to match required pattern.")

	if features.ShouldResourcesBeImported() {
		existing, err := client.Get(ctx, resGroup, serverName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing server key %s (resource group %s) ID", name, resGroup)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_postgresql_server_key", *existing.ID)
		}
	}

	properties := postgresql.ServerKey{
		ServerKeyProperties: &postgresql.ServerKeyProperties{
			ServerKeyType: utils.String(keyType),
			URI:           utils.String(keyURL),
		},
	}

	future, err := client.CreateOrUpdate(ctx, serverName, name, properties, resGroup)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, serverName, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read PostgreSQL server key %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmPostgreSQLServerKeyRead(d, meta)
}

func resourceArmPostgreSQLServerKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ServerKeysClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	serverName := id.Path["servers"]
	name := id.Path["keys"]

	resp, err := client.Get(ctx, resGroup, serverName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] PostgreSQL server key '%s' was not found (resource group '%s')", name, resGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure PostgreSQL server key %s: %+v", name, err)
	}

	d.Set("resource_group_name", resGroup)
	d.Set("server_name", serverName)
	d.Set("key_type", resp.ServerKeyType)
	d.Set("key_url", resp.URI)

	return nil
}

func resourceArmPostgreSQLServerKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ServerKeysClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	serverName := id.Path["servers"]
	name := id.Path["keys"]

	future, err := client.Delete(ctx, resGroup, serverName, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return err
	}

	return nil
}

func getServerKeyNameFromUri(uri string) string {
	log.Printf("[INFO] preparing arguments for AzureRM PostgreSQL Server Key creation.")

	vault := strings.Split(strings.Replace(uri, "https://", "", 1), ".")[0]
	log.Printf("[DEBUG] Extracted vault name: %s", vault)

	key := strings.Split(strings.Replace(uri, "https://", "", 1), "/")[2]
	log.Printf("[DEBUG] Extracted key name: %s", key)

	version := strings.Split(strings.Replace(uri, "https://", "", 1), "/")[3]
	log.Printf("[DEBUG] Extracted version: %s", version)

	return fmt.Sprintf("%s_%s_%s", vault, key, version)
}
