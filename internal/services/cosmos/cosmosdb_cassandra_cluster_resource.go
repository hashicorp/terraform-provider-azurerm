// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2023-04-15/managedcassandras"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/attestation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCassandraCluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCassandraClusterCreate,
		Read:   resourceCassandraClusterRead,
		Update: resourceCassandraClusterUpdate,
		Delete: resourceCassandraClusterDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.CassandraClusterID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"delegated_management_subnet_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateSubnetID,
			},

			"default_admin_password": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"authentication_method": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(managedcassandras.AuthenticationMethodCassandra),
				ValidateFunc: validation.StringInSlice([]string{
					string(managedcassandras.AuthenticationMethodNone),
					string(managedcassandras.AuthenticationMethodCassandra),
				}, false),
			},

			"client_certificate_pems": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validate.IsCert,
				},
			},

			"external_gossip_certificate_pems": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validate.IsCert,
				},
			},

			"external_seed_node_ip_addresses": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.IsIPv4Address,
				},
			},

			"hours_between_backups": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Default:  24,
			},

			"identity": commonschema.SystemAssignedIdentityOptional(),

			"repair_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"version": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "3.11",
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"3.11",
					"4.0",
				}, false),
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceCassandraClusterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.ManagedCassandraClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroupName := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	id := managedcassandras.NewCassandraClusterID(subscriptionId, resourceGroupName, name)

	existing, err := client.CassandraClustersGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_cosmosdb_cassandra_cluster", id.ID())
	}

	expandedIdentity, err := expandCassandraClusterIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	authenticationMethod := managedcassandras.AuthenticationMethod(d.Get("authentication_method").(string))

	body := managedcassandras.ClusterResource{
		Identity: expandedIdentity,
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Properties: &managedcassandras.ClusterResourceProperties{
			AuthenticationMethod:          &authenticationMethod,
			CassandraVersion:              utils.String(d.Get("version").(string)),
			DelegatedManagementSubnetId:   utils.String(d.Get("delegated_management_subnet_id").(string)),
			HoursBetweenBackups:           utils.Int64(int64(d.Get("hours_between_backups").(int))),
			InitialCassandraAdminPassword: utils.String(d.Get("default_admin_password").(string)),
			RepairEnabled:                 utils.Bool(d.Get("repair_enabled").(bool)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("client_certificate_pems"); ok {
		body.Properties.ClientCertificates = expandCassandraClusterCertificate(v.([]interface{}))
	}

	if v, ok := d.GetOk("external_gossip_certificate_pems"); ok {
		body.Properties.ExternalGossipCertificates = expandCassandraClusterCertificate(v.([]interface{}))
	}

	if v, ok := d.GetOk("external_seed_node_ip_addresses"); ok {
		body.Properties.ExternalSeedNodes = expandCassandraClusterExternalSeedNode(v.([]interface{}))
	}

	err = client.CassandraClustersCreateUpdateThenPoll(ctx, id, body)
	if err != nil {
		return fmt.Errorf("creating %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCassandraClusterRead(d, meta)
}

func resourceCassandraClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			if res := props; res != nil {
				d.Set("delegated_management_subnet_id", props.DelegatedManagementSubnetId)
				d.Set("authentication_method", string(pointer.From(props.AuthenticationMethod)))
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

func resourceCassandraClusterUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.ManagedCassandraClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroupName := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	id := managedcassandras.NewCassandraClusterID(subscriptionId, resourceGroupName, name)

	expandedIdentity, err := expandCassandraClusterIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	authenticationMethod := managedcassandras.AuthenticationMethod(d.Get("authentication_method").(string))

	body := managedcassandras.ClusterResource{
		Identity: expandedIdentity,
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Properties: &managedcassandras.ClusterResourceProperties{
			AuthenticationMethod:          &authenticationMethod,
			CassandraVersion:              utils.String(d.Get("version").(string)),
			DelegatedManagementSubnetId:   utils.String(d.Get("delegated_management_subnet_id").(string)),
			HoursBetweenBackups:           utils.Int64(int64(d.Get("hours_between_backups").(int))),
			InitialCassandraAdminPassword: utils.String(d.Get("default_admin_password").(string)),
			RepairEnabled:                 utils.Bool(d.Get("repair_enabled").(bool)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("client_certificate_pems"); ok {
		body.Properties.ClientCertificates = expandCassandraClusterCertificate(v.([]interface{}))
	}

	if v, ok := d.GetOk("external_gossip_certificate_pems"); ok {
		body.Properties.ExternalGossipCertificates = expandCassandraClusterCertificate(v.([]interface{}))
	}

	if v, ok := d.GetOk("external_seed_node_ip_addresses"); ok {
		body.Properties.ExternalSeedNodes = expandCassandraClusterExternalSeedNode(v.([]interface{}))
	}

	// Though there is update method but Service API complains it isn't implemented
	err = client.CassandraClustersCreateUpdateThenPoll(ctx, id, body)
	if err != nil {
		return fmt.Errorf("updating %q: %+v", id, err)
	}

	// Issue: https://github.com/Azure/azure-rest-api-specs/issues/19021
	// There is a long running issue on updating this resource.
	// The API cannot update the property after WaitForCompletionRef is returned.
	// It has to wait a while after that. Then the property can be updated successfully.
	stateConf := &pluginsdk.StateChangeConf{
		Delay:      1 * time.Minute,
		Pending:    []string{string(managedcassandras.ManagedCassandraProvisioningStateUpdating)},
		Target:     []string{string(managedcassandras.ManagedCassandraProvisioningStateSucceeded)},
		Refresh:    cosmosdbCassandraClusterStateRefreshFunc(ctx, client, id),
		MinTimeout: 15 * time.Second,
		Timeout:    d.Timeout(pluginsdk.TimeoutUpdate),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", id, err)
	}

	return resourceCassandraClusterRead(d, meta)
}

func resourceCassandraClusterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.ManagedCassandraClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := managedcassandras.ParseCassandraClusterID(d.Id())
	if err != nil {
		return err
	}

	err = client.CassandraClustersDeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}

	return nil
}

func cosmosdbCassandraClusterStateRefreshFunc(ctx context.Context, client *managedcassandras.ManagedCassandrasClient, id managedcassandras.CassandraClusterId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.CassandraClustersGet(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("polling for %s: %+v", id, err)
		}

		if model := res.Model; model != nil {
			if model.Properties != nil && model.Properties.ProvisioningState != nil {
				return res, string(*model.Properties.ProvisioningState), nil
			}
		}
		return nil, "", fmt.Errorf("unable to read provisioning state")
	}
}

func expandCassandraClusterIdentity(input []interface{}) (*identity.SystemAssigned, error) {
	expanded, err := identity.ExpandSystemAssigned(input)
	if err != nil {
		return nil, err
	}

	return &identity.SystemAssigned{
		Type: expanded.Type,
	}, nil
}

func expandCassandraClusterCertificate(input []interface{}) *[]managedcassandras.Certificate {
	results := make([]managedcassandras.Certificate, 0)

	for _, pem := range input {
		result := managedcassandras.Certificate{
			Pem: utils.String(pem.(string)),
		}
		results = append(results, result)
	}

	return &results
}

func expandCassandraClusterExternalSeedNode(input []interface{}) *[]managedcassandras.SeedNode {
	results := make([]managedcassandras.SeedNode, 0)

	for _, ipAddress := range input {
		result := managedcassandras.SeedNode{
			IPAddress: utils.String(ipAddress.(string)),
		}
		results = append(results, result)
	}

	return &results
}

func flattenCassandraClusterCertificate(input *[]managedcassandras.Certificate) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var pem string
		if item.Pem != nil {
			pem = *item.Pem
		}

		results = append(results, pem)
	}

	return results
}

func flattenCassandraClusterExternalSeedNode(input *[]managedcassandras.SeedNode) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var ipAddress string
		if item.IPAddress != nil {
			ipAddress = *item.IPAddress
		}

		results = append(results, ipAddress)
	}

	return results
}

func flattenCassandraClusterIdentity(input *identity.SystemAssigned) []interface{} {
	var transform *identity.SystemAssigned

	if input != nil {
		transform = &identity.SystemAssigned{
			Type: input.Type,
		}
		if input.PrincipalId != "" {
			transform.PrincipalId = input.PrincipalId
		}
		if input.TenantId != "" {
			transform.TenantId = input.TenantId
		}
	}

	return identity.FlattenSystemAssigned(transform)
}
