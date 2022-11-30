// This is a Snyk internal resource type.
// It aggregates the values from individual Configuration resources into one map.

package postgres

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/configurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourcePostgreSQLConfigurationSet() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePostgreSQLConfigurationSetCreateUpdate,
		Update: resourcePostgreSQLConfigurationSetCreateUpdate,
		Read:   resourcePostgreSQLConfigurationSetRead,
		Delete: resourcePostgreSQLConfigurationSetDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := configurations.ParseConfigurationID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"config_map": {
				Type: pluginsdk.TypeMap,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
				Required: true,
				ForceNew: false,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"server_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ServerName,
			},
		},
	}
}

func resourcePostgreSQLConfigurationSetCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	return nil
}

func resourcePostgreSQLConfigurationSetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ConfigurationsClient
	// Manually creating the context instead of using timeouts.ForRead. The timeouts in
	// the schema above do not get used here. The way that this ResourceData gets
	// constructed by the plugin SDK, all the timeouts are set to nil, resulting in a
	// 450-second timeout (this timeout might be on Azure's side). Creating the context
	// this way reliably gets us a 30-second timeout.
	ctx, cancel := context.WithTimeout(meta.(*clients.Client).StopContext, 30*time.Second)

	defer cancel()

	serverID := strings.TrimSuffix(d.Id(), "/configurationSet/set")
	id, err := configurations.ParseServerID(serverID)
	if err != nil {
		return err
	}

	resp, err := client.ListByServer(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[WARN] PostgreSQL Server %q was not found (resource group %q)", id.ServerName, id.ResourceGroupName)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making List Configuration request on Azure PostgreSQL Server %q (Resource Group %q): %+v", id.ServerName, id.ResourceGroupName, err)
	}

	d.Set("server_name", id.ServerName)
	d.Set("resource_group_name", id.ResourceGroupName)

	configMap := make(map[string]interface{})
	configs := resp.Model.Value
	for _, conf := range *configs {
		key := conf.Name
		value := conf.Properties.Value
		configMap[*key] = *value
	}
	d.Set("config_map", configMap)

	return nil
}

func resourcePostgreSQLConfigurationSetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	return nil
}
