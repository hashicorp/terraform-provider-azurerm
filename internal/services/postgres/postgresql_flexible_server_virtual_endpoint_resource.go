package postgres

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2023-06-01-preview/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2023-06-01-preview/virtualendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

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
				Type:         pluginsdk.TypeString,
				Description:  "The name of the Virtual Endpoint",
				ForceNew:     true,
				Required:     true,
				ValidateFunc: virtualendpoints.ValidateVirtualEndpointID,
			},
			"source_server_id": {
				Type:         pluginsdk.TypeString,
				Description:  "The Resource ID of the *Source* Postgres Flexible Server this should be associated with",
				ForceNew:     true,
				Required:     true,
				ValidateFunc: servers.ValidateFlexibleServerID,
			},
			"replica_server_id": {
				Type:         pluginsdk.TypeString,
				Description:  "The Resource ID of the *Source* Postgres Flexible Server this should be associated with",
				ForceNew:     true,
				Required:     true,
				ValidateFunc: servers.ValidateFlexibleServerID,
			},
			"type": {
				Type:         pluginsdk.TypeString,
				Description:  "The type of Virtual Endpoint",
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringInSlice(virtualendpoints.PossibleValuesForVirtualEndpointType(), true),
			},
		},
	}
}

func resourcePostgresqlFlexibleServerVirtualEndpointCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.VirtualEndpointClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	flexibleServer := d.Get("source_server_id").(string)
	replicaServer := d.Get("replica_server_id").(string)
	virtualEndpointType := d.Get("type").(string)

	sourceServerId, err := servers.ParseFlexibleServerID(flexibleServer)
	if err != nil {
		return err
	}

	replicaServerId, err := servers.ParseFlexibleServerID(replicaServer)
	if err != nil {
		return err
	}

	id := virtualendpoints.NewVirtualEndpointID(sourceServerId.SubscriptionId, sourceServerId.ResourceGroupName, sourceServerId.FlexibleServerName, name)

	locks.ByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)
	defer locks.UnlockByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)

	if err = client.CreateThenPoll(ctx, id, virtualendpoints.VirtualEndpointResource{
		Name: &name,
		Properties: &virtualendpoints.VirtualEndpointResourceProperties{
			EndpointType: (*virtualendpoints.VirtualEndpointType)(&virtualEndpointType),
			Members:      &[]string{replicaServerId.FlexibleServerName}, // TODO: Can we pass multiple at once?
		},
	}); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return nil
}

func resourcePostgresqlFlexibleServerVirtualEndpointRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.VirtualEndpointClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualendpoints.ParseVirtualEndpointID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Postgresql Flexible Server Endpoint %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if err := d.Set("name", id.FlexibleServerName); err != nil {
		return fmt.Errorf("setting `name`: %+v", err)
	}

	flexibleServerId, err := servers.ParseFlexibleServerID(d.Id())
	if err != nil {
		return err
	}

	if err := d.Set("source_server_id", flexibleServerId.ID()); err != nil {
		return fmt.Errorf("setting `source_server_id`: %+v", err)
	}

	if model := resp.Model; model != nil {
		if err := d.Set("replica_server_id", (*resp.Model.Properties.Members)[0]); err != nil {
			return fmt.Errorf("setting `replica_server_id`: %+v", err)
		} //TODO: This should be more resiliant

		if err := d.Set("type", resp.Model.Type); err != nil {
			return fmt.Errorf("setting `type`: %+v", err)
		}
	}

	return nil
}

func resourcePostgresqlFlexibleServerVirtualEndpointUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	return nil
}

func resourcePostgresqlFlexibleServerVirtualEndpointDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.VirtualEndpointClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualendpoints.ParseVirtualEndpointID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)
	defer locks.UnlockByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}

	return nil
}
