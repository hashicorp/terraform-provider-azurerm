// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_cassandra_cluster

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2023-04-15/managedcassandras"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCassandraClusterRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.ManagedCassandraClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := managedcassandras.ParseCassandraClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.CassandraClustersGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Error reading %q - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading %q: %+v", id, err)
	}

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("name", id.CassandraClusterName)
	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			d.Set("delegated_management_subnet_id", props.DelegatedManagementSubnetId)
			d.Set("authentication_method", pointer.FromEnum(props.AuthenticationMethod))
			d.Set("repair_enabled", props.RepairEnabled)
			d.Set("version", props.CassandraVersion)
			d.Set("hours_between_backups", props.HoursBetweenBackups)

			if err := d.Set("client_certificate_pems", flattenCassandraClusterCertificate(props.ClientCertificates)); err != nil {
				return fmt.Errorf("setting `client_certificate_pems`: %+v", err)
			}

			if err := d.Set("external_gossip_certificate_pems", flattenCassandraClusterCertificate(props.ExternalGossipCertificates)); err != nil {
				return fmt.Errorf("setting `external_gossip_certificate_pems`: %+v", err)
			}

			if err := d.Set("external_seed_node_ip_addresses", flattenCassandraClusterExternalSeedNode(props.ExternalSeedNodes)); err != nil {
				return fmt.Errorf("setting `external_seed_node_ip_addresses`: %+v", err)
			}
		}

		if v := model.Identity; v != nil {
			if err := d.Set("identity", flattenCassandraClusterIdentity(v)); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	// The "default_admin_password" is not returned in GET response, hence setting it from config.
	d.Set("default_admin_password", d.Get("default_admin_password").(string))
	return nil
}

func cosmosdbCassandraClusterStateRefreshFunc(ctx context.Context, client *managedcassandras.ManagedCassandrasClient, id managedcassandras.CassandraClusterId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		res, err := client.CassandraClustersGet(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("polling for %s: %+v", id, err)
		}

		if model := res.Model; model != nil {
			if model.Properties != nil && model.Properties.ProvisioningState != nil {
				return res, string(*model.Properties.ProvisioningState), nil
			}
		}
		return nil, "", errors.New("unable to read provisioning state")
	}
}
