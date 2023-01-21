package mysql

import (
	"fmt"
	"log"
	"strings"
	"time"

	// nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMySQLFlexibleServerConfigurationSet() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMySQLFlexibleServerConfigurationSetCreate,
		Update: resourceMySQLFlexibleServerConfigurationSetCreate,
		Read:   resourceMySQLFlexibleServerConfigurationSetRead,
		Delete: resourceMySQLFlexibleServerConfigurationSetDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			return nil
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{

			"resource_group_name": commonschema.ResourceGroupName(),

			"server_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ServerName,
			},

			"config_map": {
				Type: pluginsdk.TypeMap,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
				Required: true,
				ForceNew: false,
			},
		},
	}
}

func resourceMySQLFlexibleServerConfigurationSetCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	return nil
}

func resourceMySQLFlexibleServerConfigurationSetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.FlexibleServerConfigurationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	serverID := strings.TrimSuffix(d.Id(), "/configurationSet/set")
	id, err := parse.FlexibleServerID(serverID)
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serverName := id.Name
	resp, err := client.ListByServer(ctx, resourceGroup, serverName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response().Response) {
			log.Printf("[WARN] %s was not found", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	configMap := make(map[string]interface{})
	configs := resp.Values()
	for _, conf := range configs {
		key := conf.Name
		value := conf.ConfigurationProperties.Value
		configMap[*key] = *value
	}

	d.Set("config_map", configMap)

	d.Set("server_name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	return nil
}

func resourceMySQLFlexibleServerConfigurationSetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	return nil
}
