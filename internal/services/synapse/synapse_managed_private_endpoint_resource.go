package synapse

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-10-15/documentdb"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2021-10-01/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2020-01-01/mysql"
	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2018-06-01-preview/sql"
	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/2019-06-01-preview/managedvirtualnetwork"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-09-01/storage"
	"github.com/Azure/azure-sdk-for-go/services/synapse/mgmt/2021-03-01/synapse"
	cognitive "github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2021-04-30/privateendpointconnections"
	mariadb "github.com/hashicorp/go-azure-sdk/resource-manager/mariadb/2018-06-01/privateendpointconnections"
	mariadbServers "github.com/hashicorp/go-azure-sdk/resource-manager/mariadb/2018-06-01/servers"
	postgresqlServers "github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/servers"
	postgresql "github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2018-06-01/privateendpointconnections"
	purview "github.com/hashicorp/go-azure-sdk/resource-manager/purview/2021-07-01/privateendpointconnection"
	search "github.com/hashicorp/go-azure-sdk/resource-manager/search/2020-03-13/privateendpointconnections"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	cosmosParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/parse"
	keyvaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	mysqlParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/parse"
	networkParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	sqlParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/sql/parse"
	storageParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const approvalDescription = "Auto-approved by Terraform"

func resourceSynapseManagedPrivateEndpoint() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSynapseManagedPrivateEndpointCreate,
		Read:   resourceSynapseManagedPrivateEndpointRead,
		Delete: resourceSynapseManagedPrivateEndpointDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ManagedPrivateEndpointID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"synapse_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WorkspaceID,
			},

			"target_resource_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"subresource_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: networkValidate.PrivateLinkSubResourceName,
			},

			"is_manual_connection": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSynapseManagedPrivateEndpointCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	workspaceClient := meta.(*clients.Client).Synapse.WorkspaceClient
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment

	workspaceId, err := parse.WorkspaceID(d.Get("synapse_workspace_id").(string))
	if err != nil {
		return err
	}

	targetResourceIdRaw := d.Get("target_resource_id").(string)
	targetResourceId, err := azure.ParseAzureResourceID(targetResourceIdRaw)
	if err != nil {
		return err
	}

	workspace, err := workspaceClient.Get(ctx, workspaceId.ResourceGroup, workspaceId.Name)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *workspaceId, err)
	}
	if workspace.WorkspaceProperties == nil || workspace.WorkspaceProperties.ManagedVirtualNetwork == nil {
		return fmt.Errorf("empty or nil `ManagedVirtualNetwork` for %s: %+v", *workspaceId, err)
	}

	id := parse.NewManagedPrivateEndpointID(workspaceId.SubscriptionId, workspaceId.ResourceGroup, workspaceId.Name, *workspace.WorkspaceProperties.ManagedVirtualNetwork, d.Get("name").(string))

	client, err := synapseClient.ManagedPrivateEndpointsClient(workspaceId.Name, environment.SynapseEndpointSuffix)
	if err != nil {
		return fmt.Errorf("building Client for %s: %+v", id, err)
	}

	// check exist
	existing, err := client.Get(ctx, id.ManagedVirtualNetworkName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of %s: %+v", id, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_synapse_managed_private_endpoint", id.ID())
	}

	managedPrivateEndpoint := managedvirtualnetwork.ManagedPrivateEndpoint{
		Properties: &managedvirtualnetwork.ManagedPrivateEndpointProperties{
			PrivateLinkResourceID: utils.String(targetResourceIdRaw),
			GroupID:               utils.String(d.Get("subresource_name").(string)),
		},
	}
	resp, err := client.Create(ctx, id.ManagedVirtualNetworkName, id.Name, managedPrivateEndpoint)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for %s", id)
	}

	timeout, _ := ctx.Deadline()

	stateConf := &pluginsdk.StateChangeConf{
		Pending:      []string{string(synapse.ProvisioningStateProvisioning)},
		Target:       []string{string(synapse.ProvisioningStateSucceeded)},
		Refresh:      managedPrivateEndpointProvisioningStateRefreshFunc(ctx, client, id),
		Timeout:      time.Until(timeout),
		PollInterval: 15 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of %s: %+v", id, err)
	}

	isManualConnection, ok := d.GetOk("is_manual_connection")
	if ok && !isManualConnection.(bool) {
		privateEndpointName := fmt.Sprintf("%s.%s", id.WorkspaceName, id.Name)

		autoApprovalFunctionMap := map[string]func(ctx context.Context, meta interface{}, id parse.ManagedPrivateEndpointId, resourceId, privateEndpointName string) error{
			"Microsoft.CognitiveServices": approveManagedPrivateEndpointForCognitive,
			"Microsoft.DBforMariaDB":      approveManagedPrivateEndpointForMariaDB,
			"Microsoft.DBforMySQL":        approveManagedPrivateEndpointForMySQL,
			"Microsoft.DBforPostgreSQL":   approveManagedPrivateEndpointForPostgreSQL,
			"Microsoft.DocumentDB":        approveManagedPrivateEndpointForCosmos,
			"Microsoft.KeyVault":          approveManagedPrivateEndpointForKeyVault,
			"Microsoft.Purview":           approveManagedPrivateEndpointForPurview,
			"Microsoft.Search":            approveManagedPrivateEndpointForSearch,
			"Microsoft.Sql":               approveManagedPrivateEndpointForSQL,
			"Microsoft.Storage":           approveManagedPrivateEndpointForStorage,
			"Microsoft.Synapse":           approveManagedPrivateEndpointForSynapse,
		}

		autoApprovalFunction, ok := autoApprovalFunctionMap[targetResourceId.Provider]
		if !ok {
			return fmt.Errorf("auto-approving %s: `%s` is not a supported provider", id, targetResourceId.Provider)
		}

		err = autoApprovalFunction(ctx, meta, id, targetResourceIdRaw, privateEndpointName)
		if err != nil {
			return err
		}

		stateConf := &pluginsdk.StateChangeConf{
			Pending:      []string{"Pending"},
			Target:       []string{"Approved"},
			Refresh:      managedPrivateEndpointConnectionStateStatusRefreshFunc(ctx, client, id),
			Timeout:      time.Until(timeout),
			PollInterval: 15 * time.Second,
		}
		if _, err = stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for connection state status of %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())
	return resourceSynapseManagedPrivateEndpointRead(d, meta)
}

func resourceSynapseManagedPrivateEndpointRead(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment

	id, err := parse.ManagedPrivateEndpointID(d.Id())
	if err != nil {
		return err
	}

	client, err := synapseClient.ManagedPrivateEndpointsClient(id.WorkspaceName, environment.SynapseEndpointSuffix)
	if err != nil {
		return fmt.Errorf("building Client for %s: %v", *id, err)
	}

	resp, err := client.Get(ctx, id.ManagedVirtualNetworkName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	workspaceId := parse.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)
	d.Set("synapse_workspace_id", workspaceId.ID())
	d.Set("name", id.Name)

	if props := resp.Properties; props != nil {
		d.Set("target_resource_id", props.PrivateLinkResourceID)
		d.Set("subresource_name", props.GroupID)
	}

	isManualConnection, ok := d.GetOk("is_manual_connection")
	d.Set("is_manual_connection", !ok || isManualConnection.(bool))

	return nil
}

func resourceSynapseManagedPrivateEndpointDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment

	id, err := parse.ManagedPrivateEndpointID(d.Id())
	if err != nil {
		return err
	}

	client, err := synapseClient.ManagedPrivateEndpointsClient(id.WorkspaceName, environment.SynapseEndpointSuffix)
	if err != nil {
		return fmt.Errorf("building Client for %s: %v", *id, err)
	}

	if _, err := client.Delete(ctx, id.ManagedVirtualNetworkName, id.Name); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func managedPrivateEndpointProvisioningStateRefreshFunc(ctx context.Context, client *managedvirtualnetwork.ManagedPrivateEndpointsClient, id parse.ManagedPrivateEndpointId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, id.ManagedVirtualNetworkName, id.Name)
		if err != nil {
			return nil, "", fmt.Errorf("polling for %s: %+v", id, err)
		}

		if resp.Properties == nil || resp.Properties.ProvisioningState == nil {
			return resp, "", fmt.Errorf("nil ProvisioningState returned for %s", id)
		}

		return resp, *resp.Properties.ProvisioningState, nil
	}
}

func managedPrivateEndpointConnectionStateStatusRefreshFunc(ctx context.Context, client *managedvirtualnetwork.ManagedPrivateEndpointsClient, id parse.ManagedPrivateEndpointId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, id.ManagedVirtualNetworkName, id.Name)
		if err != nil {
			return nil, "", fmt.Errorf("polling for %s: %+v", id, err)
		}

		if resp.Properties == nil || resp.Properties.ConnectionState == nil || resp.Properties.ConnectionState.Status == nil {
			return resp, "", fmt.Errorf("nil Status returned for connection state of %s", id)
		}

		return resp, *resp.Properties.ConnectionState.Status, nil
	}
}

func approveManagedPrivateEndpointForCognitive(ctx context.Context, meta interface{}, id parse.ManagedPrivateEndpointId, resourceId, privateEndpointName string) error {
	client := meta.(*clients.Client).Cognitive.PrivateEndpointConnectionsClient

	targetResourceId, err := cognitive.ParseAccountID(resourceId)
	if err != nil {
		return err
	}

	result, err := client.List(ctx, *targetResourceId)
	if err != nil {
		return err
	}

	var privateEndpointConnectionId *string
	for _, val := range *result.Model.Value {
		if val.Properties == nil || val.Properties.PrivateEndpoint == nil || val.Properties.PrivateEndpoint.Id == nil {
			continue
		}
		privateEndpointId, err := networkParse.PrivateEndpointID(*val.Properties.PrivateEndpoint.Id)
		if err != nil {
			continue
		}
		if privateEndpointName == privateEndpointId.Name {
			privateEndpointConnectionId = val.Id
			break
		}
	}
	if privateEndpointConnectionId == nil {
		return fmt.Errorf("finding private endpoint connection ID for %s: %+v", id, err)
	}

	parsedPrivateEndpointConnectionId, err := cognitive.ParsePrivateEndpointConnectionID(*privateEndpointConnectionId)
	if err != nil {
		return err
	}

	approvedStatus := cognitive.PrivateEndpointServiceConnectionStatusApproved
	input := cognitive.PrivateEndpointConnection{
		Properties: &cognitive.PrivateEndpointConnectionProperties{
			PrivateLinkServiceConnectionState: cognitive.PrivateLinkServiceConnectionState{
				Status:      &approvedStatus,
				Description: utils.String(approvalDescription),
			},
		},
	}
	err = client.CreateOrUpdateThenPoll(ctx, *parsedPrivateEndpointConnectionId, input)
	if err != nil {
		return fmt.Errorf("approving %s: %+v", id, err)
	}

	return nil
}

func approveManagedPrivateEndpointForCosmos(ctx context.Context, meta interface{}, id parse.ManagedPrivateEndpointId, resourceId, privateEndpointName string) error {
	client := meta.(*clients.Client).Cosmos.PrivateEndpointConnectionClient

	targetResourceId, err := cosmosParse.DatabaseAccountID(resourceId)
	if err != nil {
		return err
	}

	result, err := client.ListByDatabaseAccount(ctx, targetResourceId.ResourceGroup, targetResourceId.Name)
	if err != nil {
		return err
	}

	var privateEndpointConnectionName *string
	for _, val := range *result.Value {
		if val.PrivateEndpoint == nil || val.PrivateEndpoint.ID == nil {
			continue
		}
		privateEndpointId, err := networkParse.PrivateEndpointID(*val.PrivateEndpoint.ID)
		if err != nil {
			continue
		}
		if privateEndpointName == privateEndpointId.Name {
			privateEndpointConnectionName = val.Name
			break
		}
	}
	if privateEndpointConnectionName == nil {
		return fmt.Errorf("finding private endpoint connection name for %s: %+v", id, err)
	}

	parameters := documentdb.PrivateEndpointConnection{
		PrivateEndpointConnectionProperties: &documentdb.PrivateEndpointConnectionProperties{
			PrivateLinkServiceConnectionState: &documentdb.PrivateLinkServiceConnectionStateProperty{
				Status:      utils.String("Approved"),
				Description: utils.String(approvalDescription),
			},
		},
	}
	future, err := client.CreateOrUpdate(ctx, targetResourceId.ResourceGroup, targetResourceId.Name, *privateEndpointConnectionName, parameters)
	if err != nil {
		return fmt.Errorf("approving %s: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for approval of %s: %+v", id, err)
	}

	return nil
}

func approveManagedPrivateEndpointForKeyVault(ctx context.Context, meta interface{}, id parse.ManagedPrivateEndpointId, resourceId, privateEndpointName string) error {
	client := meta.(*clients.Client).KeyVault.PrivateEndpointConnectionsClient
	vaultsClient := meta.(*clients.Client).KeyVault.VaultsClient

	targetResourceId, err := keyvaultParse.VaultID(resourceId)
	if err != nil {
		return err
	}

	targetResource, err := vaultsClient.Get(ctx, targetResourceId.ResourceGroup, targetResourceId.Name)
	if err != nil {
		return err
	}

	var privateEndpointConnectionId *string
	for _, val := range *targetResource.Properties.PrivateEndpointConnections {
		if val.PrivateEndpoint == nil || val.PrivateEndpoint.ID == nil {
			continue
		}
		privateEndpointId, err := networkParse.PrivateEndpointID(*val.PrivateEndpoint.ID)
		if err != nil {
			continue
		}

		if privateEndpointName == privateEndpointId.Name {
			privateEndpointConnectionId = val.ID
			break
		}
	}
	if privateEndpointConnectionId == nil {
		return fmt.Errorf("finding private endpoint connection ID for %s: %+v", id, err)
	}

	parsedPrivateEndpointConnectionId, err := keyvaultParse.PrivateEndpointConnectionID(*privateEndpointConnectionId)
	if err != nil {
		return err
	}

	parameters := keyvault.PrivateEndpointConnection{
		PrivateEndpointConnectionProperties: &keyvault.PrivateEndpointConnectionProperties{
			PrivateLinkServiceConnectionState: &keyvault.PrivateLinkServiceConnectionState{
				Status:      keyvault.PrivateEndpointServiceConnectionStatusApproved,
				Description: utils.String(approvalDescription),
			},
		},
	}
	_, err = client.Put(ctx, targetResourceId.ResourceGroup, targetResourceId.Name, parsedPrivateEndpointConnectionId.Name, parameters)
	if err != nil {
		return fmt.Errorf("approving %s: %+v", id, err)
	}

	return nil
}

func approveManagedPrivateEndpointForPurview(ctx context.Context, meta interface{}, id parse.ManagedPrivateEndpointId, resourceId, privateEndpointName string) error {
	client := meta.(*clients.Client).Purview.PrivateEndpointConnectionsClient

	targetResourceId, err := purview.ParseAccountID(resourceId)
	if err != nil {
		return err
	}

	result, err := client.ListByAccount(ctx, *targetResourceId)
	if err != nil {
		return err
	}

	var privateEndpointConnectionId *string
	for _, val := range *result.Model {
		if val.Properties == nil || val.Properties.PrivateEndpoint == nil || val.Properties.PrivateEndpoint.Id == nil {
			continue
		}
		privateEndpointId, err := networkParse.PrivateEndpointID(*val.Properties.PrivateEndpoint.Id)
		if err != nil {
			continue
		}
		if privateEndpointName == privateEndpointId.Name {
			privateEndpointConnectionId = val.Id
			break
		}
	}
	if privateEndpointConnectionId == nil {
		return fmt.Errorf("finding private endpoint connection ID for %s: %+v", id, err)
	}

	parsedPrivateEndpointConnectionId, err := purview.ParsePrivateEndpointConnectionID(*privateEndpointConnectionId)
	if err != nil {
		return err
	}

	approvedStatus := purview.StatusApproved
	input := purview.PrivateEndpointConnection{
		Properties: &purview.PrivateEndpointConnectionProperties{
			PrivateLinkServiceConnectionState: &purview.PrivateLinkServiceConnectionState{
				Status:      &approvedStatus,
				Description: utils.String(approvalDescription),
			},
		},
	}
	err = client.CreateOrUpdateThenPoll(ctx, *parsedPrivateEndpointConnectionId, input)
	if err != nil {
		return fmt.Errorf("approving %s: %+v", id, err)
	}

	return nil
}

func approveManagedPrivateEndpointForMariaDB(ctx context.Context, meta interface{}, id parse.ManagedPrivateEndpointId, resourceId, privateEndpointName string) error {
	client := meta.(*clients.Client).MariaDB.PrivateEndpointConnectionClient

	targetResourceId, err := mariadb.ParseServerID(resourceId)
	if err != nil {
		return err
	}

	result, err := client.ListByServer(ctx, *targetResourceId)
	if err != nil {
		return err
	}

	var privateEndpointConnectionId *string
	for _, val := range *result.Model {
		if val.Properties == nil || val.Properties.PrivateEndpoint == nil || val.Properties.PrivateEndpoint.Id == nil {
			continue
		}
		privateEndpointId, err := networkParse.PrivateEndpointID(*val.Properties.PrivateEndpoint.Id)
		if err != nil {
			continue
		}
		if privateEndpointName == privateEndpointId.Name {
			privateEndpointConnectionId = val.Id
			break
		}
	}
	if privateEndpointConnectionId == nil {
		return fmt.Errorf("finding private endpoint connection ID for %s: %+v", id, err)
	}

	parsedPrivateEndpointConnectionId, err := mariadb.ParsePrivateEndpointConnectionID(*privateEndpointConnectionId)
	if err != nil {
		return err
	}

	input := mariadb.PrivateEndpointConnection{
		Properties: &mariadb.PrivateEndpointConnectionProperties{
			PrivateLinkServiceConnectionState: &mariadb.PrivateLinkServiceConnectionStateProperty{
				Status:      string(mariadbServers.PrivateLinkServiceConnectionStateStatusApproved),
				Description: approvalDescription,
			},
		},
	}
	err = client.CreateOrUpdateThenPoll(ctx, *parsedPrivateEndpointConnectionId, input)
	if err != nil {
		return fmt.Errorf("approving %s: %+v", id, err)
	}
	return nil
}

func approveManagedPrivateEndpointForMySQL(ctx context.Context, meta interface{}, id parse.ManagedPrivateEndpointId, resourceId, privateEndpointName string) error {
	client := meta.(*clients.Client).MySQL.PrivateEndpointConnectionClient

	targetResourceId, err := mysqlParse.ServerID(resourceId)
	if err != nil {
		return err
	}

	result, err := client.ListByServer(ctx, targetResourceId.ResourceGroup, targetResourceId.Name)
	if err != nil {
		return err
	}

	var privateEndpointConnectionName *string
	for _, val := range result.Values() {
		if val.PrivateEndpoint == nil || val.PrivateEndpoint.ID == nil {
			continue
		}
		privateEndpointId, err := networkParse.PrivateEndpointID(*val.PrivateEndpoint.ID)
		if err != nil {
			continue
		}
		if privateEndpointName == privateEndpointId.Name {
			privateEndpointConnectionName = val.Name
			break
		}
	}
	if privateEndpointConnectionName == nil {
		return fmt.Errorf("finding private endpoint connection name for %s: %+v", id, err)
	}

	parameters := mysql.PrivateEndpointConnection{
		PrivateEndpointConnectionProperties: &mysql.PrivateEndpointConnectionProperties{
			PrivateLinkServiceConnectionState: &mysql.PrivateLinkServiceConnectionStateProperty{
				Status:      utils.String(string(mysql.Approved)),
				Description: utils.String(approvalDescription),
			},
		},
	}
	future, err := client.CreateOrUpdate(ctx, targetResourceId.ResourceGroup, targetResourceId.Name, *privateEndpointConnectionName, parameters)
	if err != nil {
		return fmt.Errorf("approving %s: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for approval of %s: %+v", id, err)
	}

	return nil
}

func approveManagedPrivateEndpointForPostgreSQL(ctx context.Context, meta interface{}, id parse.ManagedPrivateEndpointId, resourceId, privateEndpointName string) error {
	client := meta.(*clients.Client).Postgres.PrivateEndpointConnectionClient

	targetResourceId, err := postgresql.ParseServerID(resourceId)
	if err != nil {
		return err
	}

	result, err := client.ListByServer(ctx, *targetResourceId)
	if err != nil {
		return err
	}

	var privateEndpointConnectionId *string
	for _, val := range *result.Model {
		if val.Properties == nil || val.Properties.PrivateEndpoint == nil || val.Properties.PrivateEndpoint.Id == nil {
			continue
		}
		privateEndpointId, err := networkParse.PrivateEndpointID(*val.Properties.PrivateEndpoint.Id)
		if err != nil {
			continue
		}
		if privateEndpointName == privateEndpointId.Name {
			privateEndpointConnectionId = val.Id
			break
		}
	}
	if privateEndpointConnectionId == nil {
		return fmt.Errorf("finding private endpoint connection ID for %s: %+v", id, err)
	}

	parsedPrivateEndpointConnectionId, err := postgresql.ParsePrivateEndpointConnectionID(*privateEndpointConnectionId)
	if err != nil {
		return err
	}

	input := postgresql.PrivateEndpointConnection{
		Properties: &postgresql.PrivateEndpointConnectionProperties{
			PrivateLinkServiceConnectionState: &postgresql.PrivateLinkServiceConnectionStateProperty{
				Status:      string(postgresqlServers.PrivateLinkServiceConnectionStateStatusApproved),
				Description: approvalDescription,
			},
		},
	}
	err = client.CreateOrUpdateThenPoll(ctx, *parsedPrivateEndpointConnectionId, input)
	if err != nil {
		return fmt.Errorf("approving %s: %+v", id, err)
	}

	return nil
}

func approveManagedPrivateEndpointForSearch(ctx context.Context, meta interface{}, id parse.ManagedPrivateEndpointId, resourceId, privateEndpointName string) error {
	client := meta.(*clients.Client).Search.PrivateEndpointConnectionsClient

	targetResourceId, err := search.ParseSearchServiceID(resourceId)
	if err != nil {
		return err
	}

	result, err := client.ListByService(ctx, *targetResourceId, search.DefaultListByServiceOperationOptions())
	if err != nil {
		return err
	}

	var privateEndpointConnectionId *string
	for _, val := range *result.Model {
		if val.Properties == nil || val.Properties.PrivateEndpoint == nil || val.Properties.PrivateEndpoint.Id == nil {
			continue
		}
		privateEndpointId, err := networkParse.PrivateEndpointID(*val.Properties.PrivateEndpoint.Id)
		if err != nil {
			continue
		}
		if privateEndpointName == privateEndpointId.Name {
			privateEndpointConnectionId = val.Id
			break
		}
	}
	if privateEndpointConnectionId == nil {
		return fmt.Errorf("finding private endpoint connection ID for %s: %+v", id, err)
	}

	parsedPrivateEndpointConnectionId, err := search.ParsePrivateEndpointConnectionID(*privateEndpointConnectionId)
	if err != nil {
		return err
	}

	approvedStatus := search.PrivateLinkServiceConnectionStatusApproved
	privateEndpointConnection := search.PrivateEndpointConnection{
		Properties: &search.PrivateEndpointConnectionProperties{
			PrivateLinkServiceConnectionState: &search.PrivateEndpointConnectionPropertiesPrivateLinkServiceConnectionState{
				Status:      &approvedStatus,
				Description: utils.String(approvalDescription),
			},
		},
	}
	_, err = client.Update(ctx, *parsedPrivateEndpointConnectionId, privateEndpointConnection, search.DefaultUpdateOperationOptions())
	if err != nil {
		return fmt.Errorf("approving %s: %+v", id, err)
	}

	return nil
}

func approveManagedPrivateEndpointForSQL(ctx context.Context, meta interface{}, id parse.ManagedPrivateEndpointId, resourceId, privateEndpointName string) error {
	client := meta.(*clients.Client).Sql.PrivateEndpointConnectionsClient

	targetResourceId, err := sqlParse.ServerID(resourceId)
	if err != nil {
		return err
	}

	result, err := client.ListByServer(ctx, targetResourceId.ResourceGroup, targetResourceId.Name)
	if err != nil {
		return err
	}

	var privateEndpointConnectionName *string
	for _, val := range result.Values() {
		if val.PrivateEndpoint == nil || val.PrivateEndpoint.ID == nil {
			continue
		}
		privateEndpointId, err := networkParse.PrivateEndpointID(*val.PrivateEndpoint.ID)
		if err != nil {
			continue
		}
		if privateEndpointName == privateEndpointId.Name {
			privateEndpointConnectionName = val.Name
			break
		}
	}
	if privateEndpointConnectionName == nil {
		return fmt.Errorf("finding private endpoint connection name for %s: %+v", id, err)
	}

	parameters := sql.PrivateEndpointConnection{
		PrivateEndpointConnectionProperties: &sql.PrivateEndpointConnectionProperties{
			PrivateLinkServiceConnectionState: &sql.PrivateLinkServiceConnectionStateProperty{
				Status:      utils.String("Approved"),
				Description: utils.String(approvalDescription),
			},
		},
	}
	future, err := client.CreateOrUpdate(ctx, targetResourceId.ResourceGroup, targetResourceId.Name, *privateEndpointConnectionName, parameters)
	if err != nil {
		return fmt.Errorf("approving %s: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for approval of %s: %+v", id, err)
	}

	return nil
}

func approveManagedPrivateEndpointForStorage(ctx context.Context, meta interface{}, id parse.ManagedPrivateEndpointId, resourceId, privateEndpointName string) error {
	client := meta.(*clients.Client).Storage.PrivateEndpointConnectionClient

	targetResourceId, err := storageParse.StorageAccountID(resourceId)
	if err != nil {
		return err
	}

	result, err := client.List(ctx, targetResourceId.ResourceGroup, targetResourceId.Name)
	if err != nil {
		return err
	}

	var privateEndpointConnectionName *string
	for _, val := range *result.Value {
		if val.PrivateEndpoint == nil || val.PrivateEndpoint.ID == nil {
			continue
		}
		privateEndpointId, err := networkParse.PrivateEndpointID(*val.PrivateEndpoint.ID)
		if err != nil {
			continue
		}
		if privateEndpointName == privateEndpointId.Name {
			privateEndpointConnectionName = val.Name
			break
		}
	}
	if privateEndpointConnectionName == nil {
		return fmt.Errorf("finding private endpoint connection name for %s: %+v", id, err)
	}

	properties := storage.PrivateEndpointConnection{
		PrivateEndpointConnectionProperties: &storage.PrivateEndpointConnectionProperties{
			PrivateLinkServiceConnectionState: &storage.PrivateLinkServiceConnectionState{
				Status:      storage.PrivateEndpointServiceConnectionStatusApproved,
				Description: utils.String(approvalDescription),
			},
		},
	}
	_, err = client.Put(ctx, targetResourceId.ResourceGroup, targetResourceId.Name, *privateEndpointConnectionName, properties)
	if err != nil {
		return fmt.Errorf("approving %s: %+v", id, err)
	}

	return nil
}

func approveManagedPrivateEndpointForSynapse(ctx context.Context, meta interface{}, id parse.ManagedPrivateEndpointId, resourceId, privateEndpointName string) error {
	client := meta.(*clients.Client).Synapse.PrivateEndpointConnectionClient

	targetResourceId, err := parse.WorkspaceID(resourceId)
	if err != nil {
		return err
	}

	result, err := client.List(ctx, targetResourceId.ResourceGroup, targetResourceId.Name)
	if err != nil {
		return err
	}

	var privateEndpointConnectionName *string
	for _, val := range result.Values() {
		if val.PrivateEndpoint == nil || val.PrivateEndpoint.ID == nil {
			continue
		}
		privateEndpointId, err := networkParse.PrivateEndpointID(*val.PrivateEndpoint.ID)
		if err != nil {
			continue
		}
		if privateEndpointName == privateEndpointId.Name {
			privateEndpointConnectionName = val.Name
			break
		}
	}
	if privateEndpointConnectionName == nil {
		return fmt.Errorf("finding private endpoint connection name for %s: %+v", id, err)
	}

	request := synapse.PrivateEndpointConnection{
		PrivateEndpointConnectionProperties: &synapse.PrivateEndpointConnectionProperties{
			PrivateLinkServiceConnectionState: &synapse.PrivateLinkServiceConnectionState{
				Status:      utils.String("Approved"),
				Description: utils.String(approvalDescription),
			},
		},
	}
	future, err := client.Create(ctx, request, targetResourceId.ResourceGroup, targetResourceId.Name, *privateEndpointConnectionName)
	if err != nil {
		return fmt.Errorf("approving %s: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for approval of %s: %+v", id, err)
	}

	return nil
}
