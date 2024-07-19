package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2023-06-01-preview/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2023-06-01-preview/virtualendpoints"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/custompollers"
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
			Create: pluginsdk.DefaultTimeout(10 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(10 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(10 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := virtualendpoints.ParseVirtualEndpointID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:        pluginsdk.TypeString,
				Description: "The name of the Virtual Endpoint",
				ForceNew:    true,
				Required:    true,
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
				Description:  "The Resource ID of the *Replica* Postgres Flexible Server this should be associated with",
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
	sourceServer := d.Get("source_server_id").(string)
	replicaServer := d.Get("replica_server_id").(string)
	virtualEndpointType := d.Get("type").(string)

	sourceServerId, err := servers.ParseFlexibleServerID(sourceServer)
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

	// This API can be a bit flaky if the same named resource is created/destroyed quickly
	// usually waiting a minute or two before redeploying is enough to resolve the conflict
	if err = client.CreateThenPoll(ctx, id, virtualendpoints.VirtualEndpointResource{
		Name: &name,
		Properties: &virtualendpoints.VirtualEndpointResourceProperties{
			EndpointType: (*virtualendpoints.VirtualEndpointType)(&virtualEndpointType),
			Members:      &[]string{replicaServerId.FlexibleServerName},
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

	if err := d.Set("name", id.VirtualEndpointName); err != nil {
		return fmt.Errorf("setting `name`: %+v", err)
	}

	flexibleServerId := servers.NewFlexibleServerID(id.SubscriptionId, id.ResourceGroupName, id.FlexibleServerName)
	if err := d.Set("source_server_id", flexibleServerId.ID()); err != nil {
		return fmt.Errorf("setting `source_server_id`: %+v", err)
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

	if model := resp.Model; model != nil && model.Properties != nil {
		// Model.Properties.Members should be a tuple => [source_server, replication_server]
		if resp.Model.Properties.Members != nil {
			replicateServerId := servers.NewFlexibleServerID(id.SubscriptionId, id.ResourceGroupName, (*resp.Model.Properties.Members)[1])
			if err := d.Set("replica_server_id", replicateServerId.ID()); err != nil {
				return fmt.Errorf("setting `replica_server_id`: %+v", err)
			}
		} else {
			// if members list is nil, this is an endpoint that was previously deleted
			log.Printf("[INFO] Postgresql Flexible Server Endpoint %q was previously deleted - removing from state", d.Id())
			d.SetId("")
			return nil
		}
	}

	return nil
}

func resourcePostgresqlFlexibleServerVirtualEndpointUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.VirtualEndpointClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualendpoints.ParseVirtualEndpointID(d.Id())
	if err != nil {
		return err
	}

	replicaServer := d.Get("replica_server_id").(string)
	virtualEndpointType := d.Get("type").(string)

	replicaServerId, err := servers.ParseFlexibleServerID(replicaServer)
	if err != nil {
		return err
	}

	locks.ByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)
	defer locks.UnlockByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)

	if err := client.UpdateThenPoll(ctx, *id, virtualendpoints.VirtualEndpointResourceForPatch{
		Properties: &virtualendpoints.VirtualEndpointResourceProperties{
			EndpointType: (*virtualendpoints.VirtualEndpointType)(&virtualEndpointType),
			Members:      &[]string{replicaServerId.FlexibleServerName},
		},
	}); err != nil {
		return fmt.Errorf("updating %q: %+v", id, err)
	}

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

	if err := DeletePostgresFlexibileServerVirtualEndpoint(ctx, client, id); err != nil {
		return err
	}

	return nil
}

// exposed so we can access from tests
func DeletePostgresFlexibileServerVirtualEndpoint(ctx context.Context, client *virtualendpoints.VirtualEndpointsClient, id *virtualendpoints.VirtualEndpointId) error {
	deletePoller := custompollers.NewPostgresFlexibleServerVirtualEndpointDeletePoller(client, *id)
	poller := pollers.NewPoller(deletePoller, 5*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err := poller.PollUntilDone(ctx); err != nil {
		return err
	}
	return nil
}
