// This is a Snyk internal resource type.
// It aggregates the values from individual Configuration resources into one map.

package mysql

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMySQLConfigurationSet() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMySQLConfigurationSetCreate,
		Read:   resourceMySQLConfigurationSetRead,
		Delete: resourceMySQLConfigurationSetDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ConfigurationID(id)
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

func resourceMySQLConfigurationSetCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	return nil
}

func resourceMySQLConfigurationSetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.ConfigurationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// id, err := parse.ConfigurationID(serverID)
	serverID := strings.TrimSuffix(d.Id(), "/configurationSet/set")
	id, err := parse.ServerID(serverID)
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serverName := id.Name
	resp, err := client.ListByServer(ctx, resourceGroup, serverName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] MySQL Server %q was not found (resource group %q)", serverName, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making List Configuration request on Azure MySQL Server %q (Resource Group %q): %+v", serverName, resourceGroup, err)
	}

	configMap := make(map[string]interface{})
	configs := resp.Value
	for _, conf := range *configs {
		key := conf.Name
		value := conf.ConfigurationProperties.Value
		configMap[*key] = *value
	}

	d.Set("server_name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("config_map", configMap)

	return nil
}

func resourceMySQLConfigurationSetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	return nil
}
