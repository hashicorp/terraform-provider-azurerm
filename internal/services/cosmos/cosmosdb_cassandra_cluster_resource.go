package cosmos

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-10-15/documentdb"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/attestation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/parse"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
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

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": commonschema.Location(),

			"delegated_management_subnet_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: networkValidate.SubnetID,
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
				Default:  string(documentdb.AuthenticationMethodCassandra),
				ValidateFunc: validation.StringInSlice([]string{
					string(documentdb.AuthenticationMethodNone),
					string(documentdb.AuthenticationMethodCassandra),
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

			"tags": tags.Schema(),
		},
	}
}

func resourceCassandraClusterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraClustersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroupName := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	id := parse.NewCassandraClusterID(subscriptionId, resourceGroupName, name)

	existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_cosmosdb_cassandra_cluster", id.ID())
	}

	expandedIdentity, err := expandCassandraClusterIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	body := documentdb.ClusterResource{
		Identity: expandedIdentity,
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Properties: &documentdb.ClusterResourceProperties{
			AuthenticationMethod:          documentdb.AuthenticationMethod(d.Get("authentication_method").(string)),
			CassandraVersion:              utils.String(d.Get("version").(string)),
			DelegatedManagementSubnetID:   utils.String(d.Get("delegated_management_subnet_id").(string)),
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

	future, err := client.CreateUpdate(ctx, id.ResourceGroup, id.Name, body)
	if err != nil {
		return fmt.Errorf("creating %q: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on create for %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCassandraClusterRead(d, meta)
}

func resourceCassandraClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraClustersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CassandraClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading %q - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading %q: %+v", id, err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	d.Set("name", id.Name)
	if props := resp.Properties; props != nil {
		if res := props; res != nil {
			d.Set("delegated_management_subnet_id", props.DelegatedManagementSubnetID)
			d.Set("authentication_method", string(props.AuthenticationMethod))
			d.Set("repair_enabled", props.RepairEnabled)
			d.Set("version", props.CassandraVersion)

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

	if v := resp.Identity; v != nil {
		if err := d.Set("identity", flattenCassandraClusterIdentity(v)); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}
	}

	// The "default_admin_password" is not returned in GET response, hence setting it from config.
	d.Set("default_admin_password", d.Get("default_admin_password").(string))
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceCassandraClusterUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraClustersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroupName := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	id := parse.NewCassandraClusterID(subscriptionId, resourceGroupName, name)

	expandedIdentity, err := expandCassandraClusterIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	body := documentdb.ClusterResource{
		Identity: expandedIdentity,
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Properties: &documentdb.ClusterResourceProperties{
			AuthenticationMethod:          documentdb.AuthenticationMethod(d.Get("authentication_method").(string)),
			CassandraVersion:              utils.String(d.Get("version").(string)),
			DelegatedManagementSubnetID:   utils.String(d.Get("delegated_management_subnet_id").(string)),
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
	_, err = client.CreateUpdate(ctx, id.ResourceGroup, id.Name, body)
	if err != nil {
		return fmt.Errorf("updating %q: %+v", id, err)
	}

	// Issue: https://github.com/Azure/azure-rest-api-specs/issues/19021
	// There is a long running issue on updating this resource.
	// The API cannot update the property after WaitForCompletionRef is returned.
	// It has to wait a while after that. Then the property can be updated successfully.
	stateConf := &pluginsdk.StateChangeConf{
		Delay:      1 * time.Minute,
		Pending:    []string{string(documentdb.ManagedCassandraProvisioningStateUpdating)},
		Target:     []string{string(documentdb.ManagedCassandraProvisioningStateSucceeded)},
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
	client := meta.(*clients.Client).Cosmos.CassandraClustersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CassandraClusterID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("deleting %q: %+v", id, err)
		}
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("waiting on delete future for %q: %+v", id, err)
	}

	return nil
}

func cosmosdbCassandraClusterStateRefreshFunc(ctx context.Context, client *documentdb.CassandraClustersClient, id parse.CassandraClusterId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return nil, "", fmt.Errorf("polling for %s: %+v", id, err)
		}

		if res.Properties != nil && res.Properties.ProvisioningState != "" {
			return res, string(res.Properties.ProvisioningState), nil
		}
		return nil, "", fmt.Errorf("unable to read provisioning state")
	}
}

func expandCassandraClusterIdentity(input []interface{}) (*documentdb.ManagedCassandraManagedServiceIdentity, error) {
	expanded, err := identity.ExpandSystemAssigned(input)
	if err != nil {
		return nil, err
	}

	return &documentdb.ManagedCassandraManagedServiceIdentity{
		Type: documentdb.ManagedCassandraResourceIdentityType(string(expanded.Type)),
	}, nil
}

func expandCassandraClusterCertificate(input []interface{}) *[]documentdb.Certificate {
	results := make([]documentdb.Certificate, 0)

	for _, pem := range input {
		result := documentdb.Certificate{
			Pem: utils.String(pem.(string)),
		}
		results = append(results, result)
	}

	return &results
}

func expandCassandraClusterExternalSeedNode(input []interface{}) *[]documentdb.SeedNode {
	results := make([]documentdb.SeedNode, 0)

	for _, ipAddress := range input {
		result := documentdb.SeedNode{
			IPAddress: utils.String(ipAddress.(string)),
		}
		results = append(results, result)
	}

	return &results
}

func flattenCassandraClusterCertificate(input *[]documentdb.Certificate) []interface{} {
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

func flattenCassandraClusterExternalSeedNode(input *[]documentdb.SeedNode) []interface{} {
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

func flattenCassandraClusterIdentity(input *documentdb.ManagedCassandraManagedServiceIdentity) []interface{} {
	var transform *identity.SystemAssigned

	if input != nil {
		transform = &identity.SystemAssigned{
			Type: identity.Type(string(input.Type)),
		}
		if input.PrincipalID != nil {
			transform.PrincipalId = *input.PrincipalID
		}
		if input.TenantID != nil {
			transform.TenantId = *input.TenantID
		}
	}

	return identity.FlattenSystemAssigned(transform)
}
