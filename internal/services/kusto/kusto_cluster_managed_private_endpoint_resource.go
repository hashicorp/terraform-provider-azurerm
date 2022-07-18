package kusto

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"time"
)

func resourceKustoClusterManagedPrivateEndpoint() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKustoClusterManagedPrivateEndpointCreateUpdate,
		Read:   resourceKustoClusterManagedPrivateEndpointRead,
		Delete: resourceKustoClusterManagedPrivateEndpointDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ManagedPrivateEndpointID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{},
	}
}

func resourceKustoClusterManagedPrivateEndpointDelete(data *schema.ResourceData, i interface{}) error {

}

func resourceKustoClusterManagedPrivateEndpointRead(data *schema.ResourceData, i interface{}) error {

}

func resourceKustoClusterManagedPrivateEndpointCreateUpdate(data *schema.ResourceData, i interface{}) error {

}
