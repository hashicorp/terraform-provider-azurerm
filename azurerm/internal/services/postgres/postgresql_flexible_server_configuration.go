package postgres

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/postgresql/mgmt/2020-02-14-preview/postgresqlflexibleservers"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourcePostgresqlFlexibleServerConfiguration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFlexibleServerConfigurationUpdate,
		Read:   resourceFlexibleServerConfigurationRead,
		Update: resourceFlexibleServerConfigurationUpdate,
		Delete: resourceFlexibleServerConfigurationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FlexibleServerConfigurationID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"server_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FlexibleServerName,
			},

			"value": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},
		},
	}
}

func resourceFlexibleServerConfigurationUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.FlexibleServersConfigurationsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Postgresql Flexible Server configuration creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	serverName := d.Get("server_name").(string)

	props := postgresqlflexibleservers.Configuration{
		ConfigurationProperties: &postgresqlflexibleservers.ConfigurationProperties{
			Value:  utils.String(d.Get("value").(string)),
			Source: utils.String("user-override"),
		},
	}

	future, err := client.Update(ctx, resGroup, serverName, name, props)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return err
	}

	resp, err := client.Get(ctx, resGroup, serverName, name)
	if err != nil {
		return err
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Azure Postgresql Flexible Server configuration %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*resp.ID)

	return resourceFlexibleServerConfigurationRead(d, meta)
}

func resourceFlexibleServerConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.FlexibleServersConfigurationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FlexibleServerConfigurationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.FlexibleServerName, id.ConfigurationName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] Azure Postgresql Flexible Server configuration '%s' was not found (resource group '%s')", id.ConfigurationName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure Postgresql Flexible Server configuration %s: %+v", id.ConfigurationName, err)
	}

	d.Set("name", id.ConfigurationName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("server_name", id.FlexibleServerName)

	if props := resp.ConfigurationProperties; props != nil {
		d.Set("value", props.Value)
	}

	return nil
}

func resourceFlexibleServerConfigurationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.FlexibleServersConfigurationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FlexibleServerConfigurationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.FlexibleServerName, id.ConfigurationName)
	if err != nil {
		return fmt.Errorf("Error retrieving Azure Postgresql Flexible Server configuration '%s': %+v", id.ConfigurationName, err)
	}

	props := postgresqlflexibleservers.Configuration{
		ConfigurationProperties: &postgresqlflexibleservers.ConfigurationProperties{
			Value:  resp.DefaultValue,
			Source: utils.String("user-override"),
		},
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.FlexibleServerName, id.ConfigurationName, props)
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
