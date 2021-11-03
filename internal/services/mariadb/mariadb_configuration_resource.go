package mariadb

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mariadb/mgmt/2018-06-01/mariadb"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mariadb/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mariadb/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMariaDbConfiguration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMariaDbConfigurationCreateUpdate,
		Read:   resourceMariaDbConfigurationRead,
		Delete: resourceMariaDbConfigurationDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.MariaDBConfigurationID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"server_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ServerName,
			},

			"value": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceMariaDbConfigurationCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MariaDB.ConfigurationsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM MariaDb Configuration creation.")
	id := parse.NewMariaDBConfigurationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("server_name").(string), d.Get("name").(string))
	value := d.Get("value").(string)

	properties := mariadb.Configuration{
		ConfigurationProperties: &mariadb.ConfigurationProperties{
			Value: utils.String(value),
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerName, id.ConfigurationName, properties)
	if err != nil {
		return fmt.Errorf("issuing create/update request for %s: %v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for create/update of %s: %v", id, err)
	}

	d.SetId(id.ID())

	return resourceMariaDbConfigurationRead(d, meta)
}

func resourceMariaDbConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MariaDB.ConfigurationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MariaDBConfigurationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.ConfigurationName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] %s was not found", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	d.Set("name", id.ConfigurationName)
	d.Set("server_name", id.ServerName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("value", resp.ConfigurationProperties.Value)

	return nil
}

func resourceMariaDbConfigurationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MariaDB.ConfigurationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MariaDBConfigurationID(d.Id())
	if err != nil {
		return err
	}

	// "delete" = resetting this to the default value
	resp, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.ConfigurationName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	properties := mariadb.Configuration{
		ConfigurationProperties: &mariadb.ConfigurationProperties{
			// we can alternatively set `source: "system-default"`
			Value: resp.DefaultValue,
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerName, id.ConfigurationName, properties)
	if err != nil {
		return err
	}

	return future.WaitForCompletionRef(ctx, client.Client)
}
