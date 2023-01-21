package mariadb

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mariadb/2018-06-01/configurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mariadb/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMariaDbConfigurationSet() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMariaDbConfigurationSetCreateUpdate,
		Update: resourceMariaDbConfigurationSetCreateUpdate,
		Read:   resourceMariaDbConfigurationSetRead,
		Delete: resourceMariaDbConfigurationSetDelete,

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

func resourceMariaDbConfigurationSetCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	return nil
}

func resourceMariaDbConfigurationSetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MariaDB.ConfigurationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	serverID := strings.TrimSuffix(d.Id(), "/configurationSet/set")

	id, err := configurations.ParseServerID(serverID)
	if err != nil {
		return err
	}

	resp, err := client.ListByServer(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[WARN] %s was not found", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	configMap := make(map[string]interface{})
	configs := resp.Model.Value
	for _, conf := range *configs {
		key := conf.Name
		value := conf.Properties.Value
		configMap[*key] = *value
	}
	d.Set("config_map", configMap)
	d.Set("server_name", id.ServerName)
	d.Set("resource_group_name", id.ResourceGroupName)

	return nil
}

func resourceMariaDbConfigurationSetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	return nil
}
