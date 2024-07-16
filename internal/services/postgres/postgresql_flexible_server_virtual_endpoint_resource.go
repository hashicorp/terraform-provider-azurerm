package postgres

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2023-06-01-preview/virtualendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var postgresqlFlexibleServerVirtualEndpointResourceName = "azurerm_postgresql_flexible_server_virtual_endpoint"

func resourcePostgresqlFlexibleServerVirtualEndpoint() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePostgresqlFlexibleServerVirtualEndpointCreate,
		Read:   resourcePostgresqlFlexibleServerVirtualEndpointRead,
		Update: resourcePostgresqlFlexibleServerVirtualEndpointUpdate,
		Delete: resourcePostgresqlFlexibleServerVirtualEndpointDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(1 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(1 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(1 * time.Hour),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := virtualendpoints.ParseVirtualEndpointID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
	}
}

func resourcePostgresqlFlexibleServerVirtualEndpointCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	return nil
}

func resourcePostgresqlFlexibleServerVirtualEndpointRead(d *pluginsdk.ResourceData, meta interface{}) error {
	return nil
}

func resourcePostgresqlFlexibleServerVirtualEndpointUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	return nil
}

func resourcePostgresqlFlexibleServerVirtualEndpointDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	return nil
}
