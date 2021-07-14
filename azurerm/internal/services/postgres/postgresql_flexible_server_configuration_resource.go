package postgres

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2021-06-01/postgresqlflexibleservers"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
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
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"server_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FlexibleServerID,
			},

			"value": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceFlexibleServerConfigurationUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Postgres.FlexibleServersConfigurationsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Postgresql Flexible Server configuration creation.")

	name := d.Get("name").(string)
	serverId, err := parse.FlexibleServerID(d.Get("server_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewFlexibleServerConfigurationID(subscriptionId, serverId.ResourceGroup, serverId.Name, name)

	props := postgresqlflexibleservers.Configuration{
		ConfigurationProperties: &postgresqlflexibleservers.ConfigurationProperties{
			Value:  utils.String(d.Get("value").(string)),
			Source: utils.String("user-override"),
		},
	}

	future, err := client.Update(ctx, serverId.ResourceGroup, serverId.Name, name, props)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for create/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceFlexibleServerConfigurationRead(d, meta)
}

func resourceFlexibleServerConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
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
			log.Printf("[WARN] %s was not found, removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request for %s: %+v", id, err)
	}

	d.Set("name", id.ConfigurationName)
	d.Set("server_id", parse.NewFlexibleServerID(subscriptionId, id.ResourceGroup, id.FlexibleServerName).ID())

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
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	props := postgresqlflexibleservers.Configuration{
		ConfigurationProperties: &postgresqlflexibleservers.ConfigurationProperties{
			Value:  resp.DefaultValue,
			Source: utils.String("user-override"),
		},
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.FlexibleServerName, id.ConfigurationName, props)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
	}

	return nil
}
