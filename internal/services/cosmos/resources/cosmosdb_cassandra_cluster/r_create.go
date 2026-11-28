// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_cassandra_cluster

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2023-04-15/managedcassandras"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCassandraClusterCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.ManagedCassandraClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroupName := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	id := managedcassandras.NewCassandraClusterID(subscriptionId, resourceGroupName, name)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.CassandraClustersGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_cosmosdb_cassandra_cluster", id.ID())
		}
	}

	expandedIdentity, err := expandCassandraClusterIdentity(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	authenticationMethod := managedcassandras.AuthenticationMethod(d.Get("authentication_method").(string))

	body := managedcassandras.ClusterResource{
		Identity: expandedIdentity,
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &managedcassandras.ClusterResourceProperties{
			AuthenticationMethod:          &authenticationMethod,
			CassandraVersion:              pointer.To(d.Get("version").(string)),
			DelegatedManagementSubnetId:   pointer.To(d.Get("delegated_management_subnet_id").(string)),
			HoursBetweenBackups:           pointer.To(int64(d.Get("hours_between_backups").(int))),
			InitialCassandraAdminPassword: pointer.To(d.Get("default_admin_password").(string)),
			RepairEnabled:                 pointer.To(d.Get("repair_enabled").(bool)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if v, ok := d.GetOk("client_certificate_pems"); ok {
		body.Properties.ClientCertificates = expandCassandraClusterCertificate(v.([]any))
	}

	if v, ok := d.GetOk("external_gossip_certificate_pems"); ok {
		body.Properties.ExternalGossipCertificates = expandCassandraClusterCertificate(v.([]any))
	}

	if v, ok := d.GetOk("external_seed_node_ip_addresses"); ok {
		body.Properties.ExternalSeedNodes = expandCassandraClusterExternalSeedNode(v.([]any))
	}

	if err = client.CassandraClustersCreateUpdateCallbackThenPoll(ctx, id, body, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCassandraClusterRead(d, meta)
}
