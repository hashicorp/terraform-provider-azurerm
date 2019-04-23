package azurerm

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	resourcesprofile "github.com/Azure/azure-sdk-for-go/profiles/2017-03-09/resources/mgmt/resources"
	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2018-01-01/apimanagement"
	appinsights "github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2015-05-01/insights"
	"github.com/Azure/azure-sdk-for-go/services/automation/mgmt/2015-10-31/automation"
	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2018-12-01/batch"
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2017-10-12/cdn"
	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/mgmt/2017-04-18/cognitiveservices"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-06-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2018-10-01/containerinstance"
	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2017-10-01/containerregistry"
	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2018-03-31/containerservice"
	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb"
	"github.com/Azure/azure-sdk-for-go/services/databricks/mgmt/2018-04-01/databricks"
	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	analyticsAccount "github.com/Azure/azure-sdk-for-go/services/datalake/analytics/mgmt/2016-11-01/account"
	"github.com/Azure/azure-sdk-for-go/services/datalake/store/2016-11-01/filesystem"
	storeAccount "github.com/Azure/azure-sdk-for-go/services/datalake/store/mgmt/2016-11-01/account"
	"github.com/Azure/azure-sdk-for-go/services/devtestlabs/mgmt/2016-05-15/dtl"
	"github.com/Azure/azure-sdk-for-go/services/eventhub/mgmt/2017-04-01/eventhub"
	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	keyVault "github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2018-02-14/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/logic/mgmt/2016-06-01/logic"
	"github.com/Azure/azure-sdk-for-go/services/mediaservices/mgmt/2018-07-01/media"
	"github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2017-12-01/mysql"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
	"github.com/Azure/azure-sdk-for-go/services/notificationhubs/mgmt/2017-04-01/notificationhubs"
	"github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2017-12-01/postgresql"
	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2018-01-01-preview/authorization"
	"github.com/Azure/azure-sdk-for-go/services/preview/devspaces/mgmt/2018-06-01-preview/devspaces"
	"github.com/Azure/azure-sdk-for-go/services/preview/dns/mgmt/2018-03-01-preview/dns"
	"github.com/Azure/azure-sdk-for-go/services/preview/eventgrid/mgmt/2018-09-15-preview/eventgrid"
	"github.com/Azure/azure-sdk-for-go/services/preview/hdinsight/mgmt/2018-06-01-preview/hdinsight"
	"github.com/Azure/azure-sdk-for-go/services/preview/iothub/mgmt/2018-12-01-preview/devices"
	"github.com/Azure/azure-sdk-for-go/services/preview/mariadb/mgmt/2018-06-01-preview/mariadb"
	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-03-01/insights"
	"github.com/Azure/azure-sdk-for-go/services/preview/msi/mgmt/2015-08-31-preview/msi"
	"github.com/Azure/azure-sdk-for-go/services/preview/operationalinsights/mgmt/2015-11-01-preview/operationalinsights"
	"github.com/Azure/azure-sdk-for-go/services/preview/operationsmanagement/mgmt/2015-11-01-preview/operationsmanagement"
	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2018-03-01-preview/managementgroups"
	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v1.0/security"
	"github.com/Azure/azure-sdk-for-go/services/preview/signalr/mgmt/2018-03-01-preview/signalr"
	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2015-05-01-preview/sql"
	MsSql "github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-10-01-preview/sql"
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2016-06-01/recoveryservices"
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2017-07-01/backup"
	"github.com/Azure/azure-sdk-for-go/services/redis/mgmt/2018-03-01/redis"
	"github.com/Azure/azure-sdk-for-go/services/relay/mgmt/2017-04-01/relay"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2016-06-01/subscriptions"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2016-09-01/locks"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/policy"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/Azure/azure-sdk-for-go/services/scheduler/mgmt/2016-03-01/scheduler"
	"github.com/Azure/azure-sdk-for-go/services/search/mgmt/2015-08-19/search"
	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/Azure/azure-sdk-for-go/services/servicefabric/mgmt/2018-02-01/servicefabric"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2018-02-01/storage"
	"github.com/Azure/azure-sdk-for-go/services/streamanalytics/mgmt/2016-03-01/streamanalytics"
	"github.com/Azure/azure-sdk-for-go/services/trafficmanager/mgmt/2018-04-01/trafficmanager"
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2018-02-01/web"

	mainStorage "github.com/Azure/azure-sdk-for-go/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	az "github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform/httpclient"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/terraform-providers/terraform-provider-azurerm/version"
)

// ArmClient contains the handles to all the specific Azure Resource Manager
// resource classes' respective clients.
type ArmClient struct {
	clientId                 string
	tenantId                 string
	subscriptionId           string
	partnerId                string
	usingServicePrincipal    bool
	environment              az.Environment
	skipProviderRegistration bool

	StopContext context.Context

	cosmosDBClient documentdb.DatabaseAccountsClient

	automationAccountClient               automation.AccountClient
	automationAgentRegistrationInfoClient automation.AgentRegistrationInformationClient
	automationCredentialClient            automation.CredentialClient
	automationDscConfigurationClient      automation.DscConfigurationClient
	automationDscNodeConfigurationClient  automation.DscNodeConfigurationClient
	automationModuleClient                automation.ModuleClient
	automationRunbookClient               automation.RunbookClient
	automationRunbookDraftClient          automation.RunbookDraftClient
	automationScheduleClient              automation.ScheduleClient

	dnsClient   dns.RecordSetsClient
	zonesClient dns.ZonesClient

	containerRegistryClient             containerregistry.RegistriesClient
	containerRegistryReplicationsClient containerregistry.ReplicationsClient
	containerServicesClient             containerservice.ContainerServicesClient
	kubernetesClustersClient            containerservice.ManagedClustersClient
	containerGroupsClient               containerinstance.ContainerGroupsClient

	eventGridDomainsClient            eventgrid.DomainsClient
	eventGridEventSubscriptionsClient eventgrid.EventSubscriptionsClient
	eventGridTopicsClient             eventgrid.TopicsClient
	eventHubClient                    eventhub.EventHubsClient
	eventHubConsumerGroupClient       eventhub.ConsumerGroupsClient
	eventHubNamespacesClient          eventhub.NamespacesClient

	solutionsClient operationsmanagement.SolutionsClient

	redisClient               redis.Client
	redisFirewallClient       redis.FirewallRulesClient
	redisPatchSchedulesClient redis.PatchSchedulesClient

	// API Management
	apiManagementApiClient                  apimanagement.APIClient
	apiManagementApiOperationsClient        apimanagement.APIOperationClient
	apiManagementApiVersionSetClient        apimanagement.APIVersionSetClient
	apiManagementAuthorizationServersClient apimanagement.AuthorizationServerClient
	apiManagementCertificatesClient         apimanagement.CertificateClient
	apiManagementGroupClient                apimanagement.GroupClient
	apiManagementGroupUsersClient           apimanagement.GroupUserClient
	apiManagementLoggerClient               apimanagement.LoggerClient
	apiManagementOpenIdConnectClient        apimanagement.OpenIDConnectProviderClient
	apiManagementPolicyClient               apimanagement.PolicyClient
	apiManagementProductsClient             apimanagement.ProductClient
	apiManagementProductApisClient          apimanagement.ProductAPIClient
	apiManagementProductGroupsClient        apimanagement.ProductGroupClient
	apiManagementPropertyClient             apimanagement.PropertyClient
	apiManagementServiceClient              apimanagement.ServiceClient
	apiManagementSignInClient               apimanagement.SignInSettingsClient
	apiManagementSignUpClient               apimanagement.SignUpSettingsClient
	apiManagementSubscriptionsClient        apimanagement.SubscriptionClient
	apiManagementUsersClient                apimanagement.UserClient

	// Application Insights
	appInsightsClient       appinsights.ComponentsClient
	appInsightsAPIKeyClient appinsights.APIKeysClient

	// Authentication
	roleAssignmentsClient   authorization.RoleAssignmentsClient
	roleDefinitionsClient   authorization.RoleDefinitionsClient
	applicationsClient      graphrbac.ApplicationsClient
	servicePrincipalsClient graphrbac.ServicePrincipalsClient

	// Autoscale Settings
	autoscaleSettingsClient insights.AutoscaleSettingsClient

	// Batch
	batchAccountClient     batch.AccountClient
	batchCertificateClient batch.CertificateClient
	batchPoolClient        batch.PoolClient

	// CDN
	cdnCustomDomainsClient cdn.CustomDomainsClient
	cdnEndpointsClient     cdn.EndpointsClient
	cdnProfilesClient      cdn.ProfilesClient

	// Cognitive Services
	cognitiveAccountsClient cognitiveservices.AccountsClient

	// Compute
	availSetClient             compute.AvailabilitySetsClient
	diskClient                 compute.DisksClient
	imageClient                compute.ImagesClient
	galleriesClient            compute.GalleriesClient
	galleryImagesClient        compute.GalleryImagesClient
	galleryImageVersionsClient compute.GalleryImageVersionsClient
	snapshotsClient            compute.SnapshotsClient
	usageOpsClient             compute.UsageClient
	vmExtensionImageClient     compute.VirtualMachineExtensionImagesClient
	vmExtensionClient          compute.VirtualMachineExtensionsClient
	vmScaleSetClient           compute.VirtualMachineScaleSetsClient
	vmImageClient              compute.VirtualMachineImagesClient
	vmClient                   compute.VirtualMachinesClient

	// Devices
	iothubResourceClient devices.IotHubResourceClient

	// DevTestLabs
	devTestLabsClient            dtl.LabsClient
	devTestPoliciesClient        dtl.PoliciesClient
	devTestVirtualMachinesClient dtl.VirtualMachinesClient
	devTestVirtualNetworksClient dtl.VirtualNetworksClient

	// DevSpace
	devSpaceControllerClient devspaces.ControllersClient

	// Databases
	mariadbDatabasesClient                   mariadb.DatabasesClient
	mariadbServersClient                     mariadb.ServersClient
	mysqlConfigurationsClient                mysql.ConfigurationsClient
	mysqlDatabasesClient                     mysql.DatabasesClient
	mysqlFirewallRulesClient                 mysql.FirewallRulesClient
	mysqlServersClient                       mysql.ServersClient
	mysqlVirtualNetworkRulesClient           mysql.VirtualNetworkRulesClient
	postgresqlConfigurationsClient           postgresql.ConfigurationsClient
	postgresqlDatabasesClient                postgresql.DatabasesClient
	postgresqlFirewallRulesClient            postgresql.FirewallRulesClient
	postgresqlServersClient                  postgresql.ServersClient
	postgresqlVirtualNetworkRulesClient      postgresql.VirtualNetworkRulesClient
	sqlDatabasesClient                       sql.DatabasesClient
	sqlDatabaseThreatDetectionPoliciesClient sql.DatabaseThreatDetectionPoliciesClient
	sqlElasticPoolsClient                    sql.ElasticPoolsClient
	// Client for the new 2017-10-01-preview SQL API which implements vCore, DTU, and Azure data standards
	msSqlElasticPoolsClient              MsSql.ElasticPoolsClient
	sqlFirewallRulesClient               sql.FirewallRulesClient
	sqlServersClient                     sql.ServersClient
	sqlServerAzureADAdministratorsClient sql.ServerAzureADAdministratorsClient
	sqlVirtualNetworkRulesClient         sql.VirtualNetworkRulesClient

	// Data Factory
	dataFactoryPipelineClient      datafactory.PipelinesClient
	dataFactoryClient              datafactory.FactoriesClient
	dataFactoryDatasetClient       datafactory.DatasetsClient
	dataFactoryLinkedServiceClient datafactory.LinkedServicesClient

	// Data Lake Store
	dataLakeStoreAccountClient       storeAccount.AccountsClient
	dataLakeStoreFirewallRulesClient storeAccount.FirewallRulesClient
	dataLakeStoreFilesClient         filesystem.Client

	// Data Lake Analytics
	dataLakeAnalyticsAccountClient       analyticsAccount.AccountsClient
	dataLakeAnalyticsFirewallRulesClient analyticsAccount.FirewallRulesClient

	// Databricks
	databricksWorkspacesClient databricks.WorkspacesClient

	// HDInsight
	hdinsightApplicationsClient   hdinsight.ApplicationsClient
	hdinsightClustersClient       hdinsight.ClustersClient
	hdinsightConfigurationsClient hdinsight.ConfigurationsClient

	// KeyVault
	keyVaultClient           keyvault.VaultsClient
	keyVaultManagementClient keyVault.BaseClient

	// Log Analytics
	linkedServicesClient operationalinsights.LinkedServicesClient
	workspacesClient     operationalinsights.WorkspacesClient

	// Logic
	logicWorkflowsClient logic.WorkflowsClient

	// Management Groups
	managementGroupsClient             managementgroups.Client
	managementGroupsSubscriptionClient managementgroups.SubscriptionsClient

	// Media
	mediaServicesClient media.MediaservicesClient

	// Monitor
	monitorActionGroupsClient               insights.ActionGroupsClient
	monitorActivityLogAlertsClient          insights.ActivityLogAlertsClient
	monitorAlertRulesClient                 insights.AlertRulesClient
	monitorDiagnosticSettingsClient         insights.DiagnosticSettingsClient
	monitorDiagnosticSettingsCategoryClient insights.DiagnosticSettingsCategoryClient
	monitorLogProfilesClient                insights.LogProfilesClient
	monitorMetricAlertsClient               insights.MetricAlertsClient

	// MSI
	userAssignedIdentitiesClient msi.UserAssignedIdentitiesClient

	// Networking
	applicationGatewayClient        network.ApplicationGatewaysClient
	applicationSecurityGroupsClient network.ApplicationSecurityGroupsClient
	azureFirewallsClient            network.AzureFirewallsClient
	connectionMonitorsClient        network.ConnectionMonitorsClient
	ddosProtectionPlanClient        network.DdosProtectionPlansClient
	expressRouteAuthsClient         network.ExpressRouteCircuitAuthorizationsClient
	expressRouteCircuitClient       network.ExpressRouteCircuitsClient
	expressRoutePeeringsClient      network.ExpressRouteCircuitPeeringsClient
	ifaceClient                     network.InterfacesClient
	loadBalancerClient              network.LoadBalancersClient
	localNetConnClient              network.LocalNetworkGatewaysClient
	packetCapturesClient            network.PacketCapturesClient
	publicIPClient                  network.PublicIPAddressesClient
	publicIPPrefixClient            network.PublicIPPrefixesClient
	routesClient                    network.RoutesClient
	routeTablesClient               network.RouteTablesClient
	secGroupClient                  network.SecurityGroupsClient
	secRuleClient                   network.SecurityRulesClient
	subnetClient                    network.SubnetsClient
	vnetGatewayConnectionsClient    network.VirtualNetworkGatewayConnectionsClient
	vnetGatewayClient               network.VirtualNetworkGatewaysClient
	vnetClient                      network.VirtualNetworksClient
	vnetPeeringsClient              network.VirtualNetworkPeeringsClient
	watcherClient                   network.WatchersClient

	// Notification Hubs
	notificationHubsClient       notificationhubs.Client
	notificationNamespacesClient notificationhubs.NamespacesClient

	// Recovery Services
	recoveryServicesVaultsClient             recoveryservices.VaultsClient
	recoveryServicesProtectedItemsClient     backup.ProtectedItemsGroupClient
	recoveryServicesProtectionPoliciesClient backup.ProtectionPoliciesClient

	// Relay
	relayNamespacesClient relay.NamespacesClient

	// Resources
	managementLocksClient locks.ManagementLocksClient
	deploymentsClient     resources.DeploymentsClient
	providersClient       resourcesprofile.ProvidersClient
	resourcesClient       resources.Client
	resourceGroupsClient  resources.GroupsClient
	subscriptionsClient   subscriptions.Client

	// Scheduler
	schedulerJobCollectionsClient scheduler.JobCollectionsClient //nolint: megacheck
	schedulerJobsClient           scheduler.JobsClient           //nolint: megacheck

	// Search
	searchServicesClient  search.ServicesClient
	searchAdminKeysClient search.AdminKeysClient

	// Security Centre
	securityCenterPricingClient   security.PricingsClient
	securityCenterContactsClient  security.ContactsClient
	securityCenterWorkspaceClient security.WorkspaceSettingsClient

	// ServiceBus
	serviceBusQueuesClient            servicebus.QueuesClient
	serviceBusNamespacesClient        servicebus.NamespacesClient
	serviceBusTopicsClient            servicebus.TopicsClient
	serviceBusSubscriptionsClient     servicebus.SubscriptionsClient
	serviceBusSubscriptionRulesClient servicebus.RulesClient

	// Service Fabric
	serviceFabricClustersClient servicefabric.ClustersClient

	// SignalR
	signalRClient signalr.Client

	// Storage
	storageServiceClient storage.AccountsClient
	storageUsageClient   storage.UsageClient

	// Stream Analytics
	streamAnalyticsFunctionsClient       streamanalytics.FunctionsClient
	streamAnalyticsJobsClient            streamanalytics.StreamingJobsClient
	streamAnalyticsInputsClient          streamanalytics.InputsClient
	streamAnalyticsOutputsClient         streamanalytics.OutputsClient
	streamAnalyticsTransformationsClient streamanalytics.TransformationsClient

	// Traffic Manager
	trafficManagerGeographialHierarchiesClient trafficmanager.GeographicHierarchiesClient
	trafficManagerProfilesClient               trafficmanager.ProfilesClient
	trafficManagerEndpointsClient              trafficmanager.EndpointsClient

	// Web
	appServicePlansClient web.AppServicePlansClient
	appServicesClient     web.AppsClient

	// Policy
	policyAssignmentsClient    policy.AssignmentsClient
	policyDefinitionsClient    policy.DefinitionsClient
	policySetDefinitionsClient policy.SetDefinitionsClient
}

var (
	msClientRequestIDOnce sync.Once
	msClientRequestID     string
)

// clientRequestID generates a UUID to pass through `x-ms-client-request-id` header.
func clientRequestID() string {
	msClientRequestIDOnce.Do(func() {
		var err error
		msClientRequestID, err = uuid.GenerateUUID()

		if err != nil {
			log.Printf("[WARN] Fail to generate uuid for msClientRequestID: %+v", err)
		}
	})

	log.Printf("[DEBUG] AzureRM Client Request Id: %s", msClientRequestID)
	return msClientRequestID
}

func (c *ArmClient) configureClient(client *autorest.Client, auth autorest.Authorizer) {
	setUserAgent(client, c.partnerId)
	client.Authorizer = auth
	//client.RequestInspector = azure.WithClientID(clientRequestID())
	client.Sender = azure.BuildSender()
	client.SkipResourceProviderRegistration = c.skipProviderRegistration
	client.PollingDuration = 60 * time.Minute
}

func setUserAgent(client *autorest.Client, partnerID string) {
	// TODO: This is the SDK version not the CLI version, once we are on 0.12, should revisit
	tfUserAgent := httpclient.UserAgentString()

	pv := version.ProviderVersion
	providerUserAgent := fmt.Sprintf("%s terraform-provider-azurerm/%s", tfUserAgent, pv)
	client.UserAgent = strings.TrimSpace(fmt.Sprintf("%s %s", client.UserAgent, providerUserAgent))

	// append the CloudShell version to the user agent if it exists
	if azureAgent := os.Getenv("AZURE_HTTP_USER_AGENT"); azureAgent != "" {
		client.UserAgent = fmt.Sprintf("%s %s", client.UserAgent, azureAgent)
	}

	if partnerID != "" {
		client.UserAgent = fmt.Sprintf("%s pid-%s", client.UserAgent, partnerID)
	}

	log.Printf("[DEBUG] AzureRM Client User Agent: %s\n", client.UserAgent)
}

// getArmClient is a helper method which returns a fully instantiated
// *ArmClient based on the Config's current settings.
func getArmClient(c *authentication.Config, skipProviderRegistration bool, partnerId string) (*ArmClient, error) {
	env, err := authentication.DetermineEnvironment(c.Environment)
	if err != nil {
		return nil, err
	}

	// client declarations:
	client := ArmClient{
		clientId:                 c.ClientID,
		tenantId:                 c.TenantID,
		subscriptionId:           c.SubscriptionID,
		partnerId:                partnerId,
		environment:              *env,
		usingServicePrincipal:    c.AuthenticatedAsAServicePrincipal,
		skipProviderRegistration: skipProviderRegistration,
	}

	oauthConfig, err := adal.NewOAuthConfig(env.ActiveDirectoryEndpoint, c.TenantID)
	if err != nil {
		return nil, err
	}

	// OAuthConfigForTenant returns a pointer, which can be nil.
	if oauthConfig == nil {
		return nil, fmt.Errorf("Unable to configure OAuthConfig for tenant %s", c.TenantID)
	}

	sender := azure.BuildSender()

	// Resource Manager endpoints
	endpoint := env.ResourceManagerEndpoint
	auth, err := c.GetAuthorizationToken(sender, oauthConfig, env.TokenAudience)
	if err != nil {
		return nil, err
	}

	// Graph Endpoints
	graphEndpoint := env.GraphEndpoint
	graphAuth, err := c.GetAuthorizationToken(sender, oauthConfig, graphEndpoint)
	if err != nil {
		return nil, err
	}

	// Key Vault Endpoints
	keyVaultAuth := autorest.NewBearerAuthorizerCallback(sender, func(tenantID, resource string) (*autorest.BearerAuthorizer, error) {
		keyVaultSpt, err := c.GetAuthorizationToken(sender, oauthConfig, resource)
		if err != nil {
			return nil, err
		}

		return keyVaultSpt, nil
	})

	client.registerApiManagementServiceClients(endpoint, c.SubscriptionID, auth)
	client.registerAppInsightsClients(endpoint, c.SubscriptionID, auth)
	client.registerAutomationClients(endpoint, c.SubscriptionID, auth)
	client.registerAuthentication(endpoint, graphEndpoint, c.SubscriptionID, c.TenantID, auth, graphAuth)
	client.registerBatchClients(endpoint, c.SubscriptionID, auth)
	client.registerCDNClients(endpoint, c.SubscriptionID, auth)
	client.registerCognitiveServiceClients(endpoint, c.SubscriptionID, auth)
	client.registerComputeClients(endpoint, c.SubscriptionID, auth)
	client.registerContainerInstanceClients(endpoint, c.SubscriptionID, auth)
	client.registerContainerRegistryClients(endpoint, c.SubscriptionID, auth)
	client.registerContainerServicesClients(endpoint, c.SubscriptionID, auth)
	client.registerCosmosDBClients(endpoint, c.SubscriptionID, auth)
	client.registerDatabricksClients(endpoint, c.SubscriptionID, auth)
	client.registerDatabases(endpoint, c.SubscriptionID, auth, sender)
	client.registerDataFactoryClients(endpoint, c.SubscriptionID, auth)
	client.registerDataLakeStoreClients(endpoint, c.SubscriptionID, auth)
	client.registerDeviceClients(endpoint, c.SubscriptionID, auth)
	client.registerDevSpaceClients(endpoint, c.SubscriptionID, auth)
	client.registerDevTestClients(endpoint, c.SubscriptionID, auth)
	client.registerDNSClients(endpoint, c.SubscriptionID, auth)
	client.registerEventGridClients(endpoint, c.SubscriptionID, auth)
	client.registerEventHubClients(endpoint, c.SubscriptionID, auth)
	client.registerHDInsightsClients(endpoint, c.SubscriptionID, auth)
	client.registerKeyVaultClients(endpoint, c.SubscriptionID, auth, keyVaultAuth)
	client.registerLogicClients(endpoint, c.SubscriptionID, auth)
	client.registerMediaServiceClients(endpoint, c.SubscriptionID, auth)
	client.registerMonitorClients(endpoint, c.SubscriptionID, auth)
	client.registerNetworkingClients(endpoint, c.SubscriptionID, auth)
	client.registerNotificationHubsClient(endpoint, c.SubscriptionID, auth)
	client.registerOperationalInsightsClients(endpoint, c.SubscriptionID, auth)
	client.registerRecoveryServiceClients(endpoint, c.SubscriptionID, auth)
	client.registerPolicyClients(endpoint, c.SubscriptionID, auth)
	client.registerManagementGroupClients(endpoint, auth)
	client.registerRedisClients(endpoint, c.SubscriptionID, auth)
	client.registerRelayClients(endpoint, c.SubscriptionID, auth)
	client.registerResourcesClients(endpoint, c.SubscriptionID, auth)
	client.registerSearchClients(endpoint, c.SubscriptionID, auth)
	client.registerSecurityCenterClients(endpoint, c.SubscriptionID, auth)
	client.registerServiceBusClients(endpoint, c.SubscriptionID, auth)
	client.registerServiceFabricClients(endpoint, c.SubscriptionID, auth)
	client.registerSchedulerClients(endpoint, c.SubscriptionID, auth)
	client.registerSignalRClients(endpoint, c.SubscriptionID, auth)
	client.registerStorageClients(endpoint, c.SubscriptionID, auth)
	client.registerStreamAnalyticsClients(endpoint, c.SubscriptionID, auth)
	client.registerTrafficManagerClients(endpoint, c.SubscriptionID, auth)
	client.registerWebClients(endpoint, c.SubscriptionID, auth)

	return &client, nil
}

func (c *ArmClient) registerApiManagementServiceClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	apisClient := apimanagement.NewAPIClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&apisClient.Client, auth)
	c.apiManagementApiClient = apisClient

	apiOperationsClient := apimanagement.NewAPIOperationClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&apiOperationsClient.Client, auth)
	c.apiManagementApiOperationsClient = apiOperationsClient

	apiVersionSetClient := apimanagement.NewAPIVersionSetClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&apiVersionSetClient.Client, auth)
	c.apiManagementApiVersionSetClient = apiVersionSetClient

	authorizationServersClient := apimanagement.NewAuthorizationServerClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&authorizationServersClient.Client, auth)
	c.apiManagementAuthorizationServersClient = authorizationServersClient

	certificatesClient := apimanagement.NewCertificateClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&certificatesClient.Client, auth)
	c.apiManagementCertificatesClient = certificatesClient

	groupsClient := apimanagement.NewGroupClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&groupsClient.Client, auth)
	c.apiManagementGroupClient = groupsClient

	groupUsersClient := apimanagement.NewGroupUserClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&groupUsersClient.Client, auth)
	c.apiManagementGroupUsersClient = groupUsersClient

	loggerClient := apimanagement.NewLoggerClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&loggerClient.Client, auth)
	c.apiManagementLoggerClient = loggerClient

	policyClient := apimanagement.NewPolicyClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&policyClient.Client, auth)
	c.apiManagementPolicyClient = policyClient

	serviceClient := apimanagement.NewServiceClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&serviceClient.Client, auth)
	c.apiManagementServiceClient = serviceClient

	signInClient := apimanagement.NewSignInSettingsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&signInClient.Client, auth)
	c.apiManagementSignInClient = signInClient

	signUpClient := apimanagement.NewSignUpSettingsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&signUpClient.Client, auth)
	c.apiManagementSignUpClient = signUpClient

	openIdConnectClient := apimanagement.NewOpenIDConnectProviderClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&openIdConnectClient.Client, auth)
	c.apiManagementOpenIdConnectClient = openIdConnectClient

	productsClient := apimanagement.NewProductClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&productsClient.Client, auth)
	c.apiManagementProductsClient = productsClient

	productApisClient := apimanagement.NewProductAPIClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&productApisClient.Client, auth)
	c.apiManagementProductApisClient = productApisClient

	productGroupsClient := apimanagement.NewProductGroupClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&productGroupsClient.Client, auth)
	c.apiManagementProductGroupsClient = productGroupsClient

	propertiesClient := apimanagement.NewPropertyClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&propertiesClient.Client, auth)
	c.apiManagementPropertyClient = propertiesClient

	subscriptionsClient := apimanagement.NewSubscriptionClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&subscriptionsClient.Client, auth)
	c.apiManagementSubscriptionsClient = subscriptionsClient

	usersClient := apimanagement.NewUserClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&usersClient.Client, auth)
	c.apiManagementUsersClient = usersClient
}

func (c *ArmClient) registerAppInsightsClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	ai := appinsights.NewComponentsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&ai.Client, auth)
	c.appInsightsClient = ai

	aiak := appinsights.NewAPIKeysClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&aiak.Client, auth)
	c.appInsightsAPIKeyClient = aiak
}

func (c *ArmClient) registerAutomationClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	accountClient := automation.NewAccountClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&accountClient.Client, auth)
	c.automationAccountClient = accountClient

	agentRegistrationInfoClient := automation.NewAgentRegistrationInformationClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&agentRegistrationInfoClient.Client, auth)
	c.automationAgentRegistrationInfoClient = agentRegistrationInfoClient

	credentialClient := automation.NewCredentialClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&credentialClient.Client, auth)
	c.automationCredentialClient = credentialClient

	dscConfigurationClient := automation.NewDscConfigurationClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&dscConfigurationClient.Client, auth)
	c.automationDscConfigurationClient = dscConfigurationClient

	dscNodeConfigurationClient := automation.NewDscNodeConfigurationClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&dscNodeConfigurationClient.Client, auth)
	c.automationDscNodeConfigurationClient = dscNodeConfigurationClient

	moduleClient := automation.NewModuleClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&moduleClient.Client, auth)
	c.automationModuleClient = moduleClient

	runbookClient := automation.NewRunbookClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&runbookClient.Client, auth)
	c.automationRunbookClient = runbookClient

	scheduleClient := automation.NewScheduleClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&scheduleClient.Client, auth)
	c.automationScheduleClient = scheduleClient

	runbookDraftClient := automation.NewRunbookDraftClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&runbookDraftClient.Client, auth)
	c.automationRunbookDraftClient = runbookDraftClient
}

func (c *ArmClient) registerAuthentication(endpoint, graphEndpoint, subscriptionId, tenantId string, auth, graphAuth autorest.Authorizer) {
	assignmentsClient := authorization.NewRoleAssignmentsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&assignmentsClient.Client, auth)
	c.roleAssignmentsClient = assignmentsClient

	definitionsClient := authorization.NewRoleDefinitionsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&definitionsClient.Client, auth)
	c.roleDefinitionsClient = definitionsClient

	applicationsClient := graphrbac.NewApplicationsClientWithBaseURI(graphEndpoint, tenantId)
	c.configureClient(&applicationsClient.Client, graphAuth)
	c.applicationsClient = applicationsClient

	servicePrincipalsClient := graphrbac.NewServicePrincipalsClientWithBaseURI(graphEndpoint, tenantId)
	c.configureClient(&servicePrincipalsClient.Client, graphAuth)
	c.servicePrincipalsClient = servicePrincipalsClient
}

func (c *ArmClient) registerBatchClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	batchAccount := batch.NewAccountClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&batchAccount.Client, auth)
	c.batchAccountClient = batchAccount

	batchCertificateClient := batch.NewCertificateClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&batchCertificateClient.Client, auth)
	c.batchCertificateClient = batchCertificateClient

	batchPool := batch.NewPoolClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&batchPool.Client, auth)
	c.batchPoolClient = batchPool
}

func (c *ArmClient) registerCDNClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	customDomainsClient := cdn.NewCustomDomainsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&customDomainsClient.Client, auth)
	c.cdnCustomDomainsClient = customDomainsClient

	endpointsClient := cdn.NewEndpointsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&endpointsClient.Client, auth)
	c.cdnEndpointsClient = endpointsClient

	profilesClient := cdn.NewProfilesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&profilesClient.Client, auth)
	c.cdnProfilesClient = profilesClient
}

func (c *ArmClient) registerCognitiveServiceClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	accountsClient := cognitiveservices.NewAccountsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&accountsClient.Client, auth)
	c.cognitiveAccountsClient = accountsClient
}

func (c *ArmClient) registerCosmosDBClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	cdb := documentdb.NewDatabaseAccountsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&cdb.Client, auth)
	c.cosmosDBClient = cdb
}

func (c *ArmClient) registerMediaServiceClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	mediaServicesClient := media.NewMediaservicesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&mediaServicesClient.Client, auth)
	c.mediaServicesClient = mediaServicesClient
}

func (c *ArmClient) registerComputeClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	availabilitySetsClient := compute.NewAvailabilitySetsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&availabilitySetsClient.Client, auth)
	c.availSetClient = availabilitySetsClient

	diskClient := compute.NewDisksClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&diskClient.Client, auth)
	c.diskClient = diskClient

	imagesClient := compute.NewImagesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&imagesClient.Client, auth)
	c.imageClient = imagesClient

	snapshotsClient := compute.NewSnapshotsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&snapshotsClient.Client, auth)
	c.snapshotsClient = snapshotsClient

	usageClient := compute.NewUsageClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&usageClient.Client, auth)
	c.usageOpsClient = usageClient

	extensionImagesClient := compute.NewVirtualMachineExtensionImagesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&extensionImagesClient.Client, auth)
	c.vmExtensionImageClient = extensionImagesClient

	extensionsClient := compute.NewVirtualMachineExtensionsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&extensionsClient.Client, auth)
	c.vmExtensionClient = extensionsClient

	virtualMachineImagesClient := compute.NewVirtualMachineImagesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&virtualMachineImagesClient.Client, auth)
	c.vmImageClient = virtualMachineImagesClient

	scaleSetsClient := compute.NewVirtualMachineScaleSetsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&scaleSetsClient.Client, auth)
	c.vmScaleSetClient = scaleSetsClient

	virtualMachinesClient := compute.NewVirtualMachinesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&virtualMachinesClient.Client, auth)
	c.vmClient = virtualMachinesClient

	galleriesClient := compute.NewGalleriesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&galleriesClient.Client, auth)
	c.galleriesClient = galleriesClient

	galleryImagesClient := compute.NewGalleryImagesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&galleryImagesClient.Client, auth)
	c.galleryImagesClient = galleryImagesClient

	galleryImageVersionsClient := compute.NewGalleryImageVersionsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&galleryImageVersionsClient.Client, auth)
	c.galleryImageVersionsClient = galleryImageVersionsClient
}

func (c *ArmClient) registerContainerInstanceClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	cgc := containerinstance.NewContainerGroupsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&cgc.Client, auth)
	c.containerGroupsClient = cgc
}

func (c *ArmClient) registerContainerRegistryClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	crc := containerregistry.NewRegistriesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&crc.Client, auth)
	c.containerRegistryClient = crc

	// container registry replicalication client
	crrc := containerregistry.NewReplicationsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&crrc.Client, auth)
	c.containerRegistryReplicationsClient = crrc
}

func (c *ArmClient) registerContainerServicesClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	// ACS
	containerServicesClient := containerservice.NewContainerServicesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&containerServicesClient.Client, auth)
	c.containerServicesClient = containerServicesClient

	// AKS
	kubernetesClustersClient := containerservice.NewManagedClustersClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&kubernetesClustersClient.Client, auth)
	c.kubernetesClustersClient = kubernetesClustersClient
}

func (c *ArmClient) registerDatabricksClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	databricksWorkspacesClient := databricks.NewWorkspacesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&databricksWorkspacesClient.Client, auth)
	c.databricksWorkspacesClient = databricksWorkspacesClient
}

func (c *ArmClient) registerDatabases(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	mariadbDBClient := mariadb.NewDatabasesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&mariadbDBClient.Client, auth)
	c.mariadbDatabasesClient = mariadbDBClient

	mariadbServersClient := mariadb.NewServersClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&mariadbServersClient.Client, auth)
	c.mariadbServersClient = mariadbServersClient

	// MySQL
	mysqlConfigClient := mysql.NewConfigurationsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&mysqlConfigClient.Client, auth)
	c.mysqlConfigurationsClient = mysqlConfigClient

	mysqlDBClient := mysql.NewDatabasesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&mysqlDBClient.Client, auth)
	c.mysqlDatabasesClient = mysqlDBClient

	mysqlFWClient := mysql.NewFirewallRulesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&mysqlFWClient.Client, auth)
	c.mysqlFirewallRulesClient = mysqlFWClient

	mysqlServersClient := mysql.NewServersClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&mysqlServersClient.Client, auth)
	c.mysqlServersClient = mysqlServersClient

	mysqlVirtualNetworkRulesClient := mysql.NewVirtualNetworkRulesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&mysqlVirtualNetworkRulesClient.Client, auth)
	c.mysqlVirtualNetworkRulesClient = mysqlVirtualNetworkRulesClient

	// PostgreSQL
	postgresqlConfigClient := postgresql.NewConfigurationsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&postgresqlConfigClient.Client, auth)
	c.postgresqlConfigurationsClient = postgresqlConfigClient

	postgresqlDBClient := postgresql.NewDatabasesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&postgresqlDBClient.Client, auth)
	c.postgresqlDatabasesClient = postgresqlDBClient

	postgresqlFWClient := postgresql.NewFirewallRulesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&postgresqlFWClient.Client, auth)
	c.postgresqlFirewallRulesClient = postgresqlFWClient

	postgresqlSrvClient := postgresql.NewServersClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&postgresqlSrvClient.Client, auth)
	c.postgresqlServersClient = postgresqlSrvClient

	postgresqlVNRClient := postgresql.NewVirtualNetworkRulesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&postgresqlVNRClient.Client, auth)
	c.postgresqlVirtualNetworkRulesClient = postgresqlVNRClient

	// SQL Azure
	sqlDBClient := sql.NewDatabasesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&sqlDBClient.Client, auth)
	c.sqlDatabasesClient = sqlDBClient

	sqlDTDPClient := sql.NewDatabaseThreatDetectionPoliciesClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&sqlDTDPClient.Client, "")
	sqlDTDPClient.Authorizer = auth
	sqlDTDPClient.Sender = sender
	sqlDTDPClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.sqlDatabaseThreatDetectionPoliciesClient = sqlDTDPClient

	sqlFWClient := sql.NewFirewallRulesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&sqlFWClient.Client, auth)
	c.sqlFirewallRulesClient = sqlFWClient

	sqlEPClient := sql.NewElasticPoolsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&sqlEPClient.Client, auth)
	c.sqlElasticPoolsClient = sqlEPClient

	MsSqlEPClient := MsSql.NewElasticPoolsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&MsSqlEPClient.Client, auth)
	c.msSqlElasticPoolsClient = MsSqlEPClient

	sqlSrvClient := sql.NewServersClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&sqlSrvClient.Client, auth)
	c.sqlServersClient = sqlSrvClient

	sqlADClient := sql.NewServerAzureADAdministratorsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&sqlADClient.Client, auth)
	c.sqlServerAzureADAdministratorsClient = sqlADClient

	sqlVNRClient := sql.NewVirtualNetworkRulesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&sqlVNRClient.Client, auth)
	c.sqlVirtualNetworkRulesClient = sqlVNRClient
}

func (c *ArmClient) registerDataFactoryClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	dataFactoryClient := datafactory.NewFactoriesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&dataFactoryClient.Client, auth)
	c.dataFactoryClient = dataFactoryClient

	dataFactoryDatasetClient := datafactory.NewDatasetsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&dataFactoryDatasetClient.Client, auth)
	c.dataFactoryDatasetClient = dataFactoryDatasetClient

	dataFactoryLinkedServiceClient := datafactory.NewLinkedServicesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&dataFactoryLinkedServiceClient.Client, auth)
	c.dataFactoryLinkedServiceClient = dataFactoryLinkedServiceClient

	dataFactoryPipelineClient := datafactory.NewPipelinesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&dataFactoryPipelineClient.Client, auth)
	c.dataFactoryPipelineClient = dataFactoryPipelineClient
}

func (c *ArmClient) registerDataLakeStoreClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	storeAccountClient := storeAccount.NewAccountsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&storeAccountClient.Client, auth)
	c.dataLakeStoreAccountClient = storeAccountClient

	storeFirewallRulesClient := storeAccount.NewFirewallRulesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&storeFirewallRulesClient.Client, auth)
	c.dataLakeStoreFirewallRulesClient = storeFirewallRulesClient

	analyticsAccountClient := analyticsAccount.NewAccountsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&analyticsAccountClient.Client, auth)
	c.dataLakeAnalyticsAccountClient = analyticsAccountClient

	filesClient := filesystem.NewClient()
	c.configureClient(&filesClient.Client, auth)
	c.dataLakeStoreFilesClient = filesClient

	analyticsFirewallRulesClient := analyticsAccount.NewFirewallRulesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&analyticsFirewallRulesClient.Client, auth)
	c.dataLakeAnalyticsFirewallRulesClient = analyticsFirewallRulesClient
}

func (c *ArmClient) registerDeviceClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	iotClient := devices.NewIotHubResourceClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&iotClient.Client, auth)
	c.iothubResourceClient = iotClient
}

func (c *ArmClient) registerDevTestClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	labsClient := dtl.NewLabsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&labsClient.Client, auth)
	c.devTestLabsClient = labsClient

	devTestPoliciesClient := dtl.NewPoliciesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&devTestPoliciesClient.Client, auth)
	c.devTestPoliciesClient = devTestPoliciesClient

	devTestVirtualMachinesClient := dtl.NewVirtualMachinesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&devTestVirtualMachinesClient.Client, auth)
	c.devTestVirtualMachinesClient = devTestVirtualMachinesClient

	devTestVirtualNetworksClient := dtl.NewVirtualNetworksClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&devTestVirtualNetworksClient.Client, auth)
	c.devTestVirtualNetworksClient = devTestVirtualNetworksClient
}

func (c *ArmClient) registerDevSpaceClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	controllersClient := devspaces.NewControllersClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&controllersClient.Client, auth)
	c.devSpaceControllerClient = controllersClient
}

func (c *ArmClient) registerDNSClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	dn := dns.NewRecordSetsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&dn.Client, auth)
	c.dnsClient = dn

	zo := dns.NewZonesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&zo.Client, auth)
	c.zonesClient = zo
}

func (c *ArmClient) registerEventGridClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	egtc := eventgrid.NewTopicsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&egtc.Client, auth)
	c.eventGridTopicsClient = egtc

	egdc := eventgrid.NewDomainsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&egdc.Client, auth)
	c.eventGridDomainsClient = egdc

	egesc := eventgrid.NewEventSubscriptionsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&egesc.Client, auth)
	c.eventGridEventSubscriptionsClient = egesc
}

func (c *ArmClient) registerEventHubClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	ehc := eventhub.NewEventHubsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&ehc.Client, auth)
	c.eventHubClient = ehc

	chcgc := eventhub.NewConsumerGroupsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&chcgc.Client, auth)
	c.eventHubConsumerGroupClient = chcgc

	ehnc := eventhub.NewNamespacesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&ehnc.Client, auth)
	c.eventHubNamespacesClient = ehnc
}

func (c *ArmClient) registerHDInsightsClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	applicationsClient := hdinsight.NewApplicationsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&applicationsClient.Client, auth)
	c.hdinsightApplicationsClient = applicationsClient

	clustersClient := hdinsight.NewClustersClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&clustersClient.Client, auth)
	c.hdinsightClustersClient = clustersClient

	configurationsClient := hdinsight.NewConfigurationsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&configurationsClient.Client, auth)
	c.hdinsightConfigurationsClient = configurationsClient
}

func (c *ArmClient) registerKeyVaultClients(endpoint, subscriptionId string, auth autorest.Authorizer, keyVaultAuth autorest.Authorizer) {
	keyVaultClient := keyvault.NewVaultsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&keyVaultClient.Client, auth)
	c.keyVaultClient = keyVaultClient

	keyVaultManagementClient := keyVault.New()
	c.configureClient(&keyVaultManagementClient.Client, keyVaultAuth)
	c.keyVaultManagementClient = keyVaultManagementClient
}

func (c *ArmClient) registerLogicClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	workflowsClient := logic.NewWorkflowsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&workflowsClient.Client, auth)
	c.logicWorkflowsClient = workflowsClient
}

func (c *ArmClient) registerMonitorClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	agc := insights.NewActionGroupsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&agc.Client, auth)
	c.monitorActionGroupsClient = agc

	alac := insights.NewActivityLogAlertsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&alac.Client, auth)
	c.monitorActivityLogAlertsClient = alac

	arc := insights.NewAlertRulesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&arc.Client, auth)
	c.monitorAlertRulesClient = arc

	monitorLogProfilesClient := insights.NewLogProfilesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&monitorLogProfilesClient.Client, auth)
	c.monitorLogProfilesClient = monitorLogProfilesClient

	mac := insights.NewMetricAlertsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&mac.Client, auth)
	c.monitorMetricAlertsClient = mac

	autoscaleSettingsClient := insights.NewAutoscaleSettingsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&autoscaleSettingsClient.Client, auth)
	c.autoscaleSettingsClient = autoscaleSettingsClient

	monitoringInsightsClient := insights.NewDiagnosticSettingsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&monitoringInsightsClient.Client, auth)
	c.monitorDiagnosticSettingsClient = monitoringInsightsClient

	monitoringCategorySettingsClient := insights.NewDiagnosticSettingsCategoryClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&monitoringCategorySettingsClient.Client, auth)
	c.monitorDiagnosticSettingsCategoryClient = monitoringCategorySettingsClient
}

func (c *ArmClient) registerNetworkingClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	applicationGatewaysClient := network.NewApplicationGatewaysClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&applicationGatewaysClient.Client, auth)
	c.applicationGatewayClient = applicationGatewaysClient

	appSecurityGroupsClient := network.NewApplicationSecurityGroupsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&appSecurityGroupsClient.Client, auth)
	c.applicationSecurityGroupsClient = appSecurityGroupsClient

	azureFirewallsClient := network.NewAzureFirewallsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&azureFirewallsClient.Client, auth)
	c.azureFirewallsClient = azureFirewallsClient

	connectionMonitorsClient := network.NewConnectionMonitorsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&connectionMonitorsClient.Client, auth)
	c.connectionMonitorsClient = connectionMonitorsClient

	ddosProtectionPlanClient := network.NewDdosProtectionPlansClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&ddosProtectionPlanClient.Client, auth)
	c.ddosProtectionPlanClient = ddosProtectionPlanClient

	expressRouteAuthsClient := network.NewExpressRouteCircuitAuthorizationsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&expressRouteAuthsClient.Client, auth)
	c.expressRouteAuthsClient = expressRouteAuthsClient

	expressRouteCircuitsClient := network.NewExpressRouteCircuitsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&expressRouteCircuitsClient.Client, auth)
	c.expressRouteCircuitClient = expressRouteCircuitsClient

	expressRoutePeeringsClient := network.NewExpressRouteCircuitPeeringsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&expressRoutePeeringsClient.Client, auth)
	c.expressRoutePeeringsClient = expressRoutePeeringsClient

	interfacesClient := network.NewInterfacesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&interfacesClient.Client, auth)
	c.ifaceClient = interfacesClient

	loadBalancersClient := network.NewLoadBalancersClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&loadBalancersClient.Client, auth)
	c.loadBalancerClient = loadBalancersClient

	localNetworkGatewaysClient := network.NewLocalNetworkGatewaysClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&localNetworkGatewaysClient.Client, auth)
	c.localNetConnClient = localNetworkGatewaysClient

	gatewaysClient := network.NewVirtualNetworkGatewaysClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&gatewaysClient.Client, auth)
	c.vnetGatewayClient = gatewaysClient

	gatewayConnectionsClient := network.NewVirtualNetworkGatewayConnectionsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&gatewayConnectionsClient.Client, auth)
	c.vnetGatewayConnectionsClient = gatewayConnectionsClient

	networksClient := network.NewVirtualNetworksClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&networksClient.Client, auth)
	c.vnetClient = networksClient

	packetCapturesClient := network.NewPacketCapturesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&packetCapturesClient.Client, auth)
	c.packetCapturesClient = packetCapturesClient

	peeringsClient := network.NewVirtualNetworkPeeringsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&peeringsClient.Client, auth)
	c.vnetPeeringsClient = peeringsClient

	publicIPAddressesClient := network.NewPublicIPAddressesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&publicIPAddressesClient.Client, auth)
	c.publicIPClient = publicIPAddressesClient

	publicIPPrefixesClient := network.NewPublicIPPrefixesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&publicIPPrefixesClient.Client, auth)
	c.publicIPPrefixClient = publicIPPrefixesClient

	routesClient := network.NewRoutesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&routesClient.Client, auth)
	c.routesClient = routesClient

	routeTablesClient := network.NewRouteTablesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&routeTablesClient.Client, auth)
	c.routeTablesClient = routeTablesClient

	securityGroupsClient := network.NewSecurityGroupsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&securityGroupsClient.Client, auth)
	c.secGroupClient = securityGroupsClient

	securityRulesClient := network.NewSecurityRulesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&securityRulesClient.Client, auth)
	c.secRuleClient = securityRulesClient

	subnetsClient := network.NewSubnetsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&subnetsClient.Client, auth)
	c.subnetClient = subnetsClient

	userAssignedIdentitiesClient := msi.NewUserAssignedIdentitiesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&userAssignedIdentitiesClient.Client, auth)
	c.userAssignedIdentitiesClient = userAssignedIdentitiesClient

	watchersClient := network.NewWatchersClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&watchersClient.Client, auth)
	c.watcherClient = watchersClient
}

func (c *ArmClient) registerNotificationHubsClient(endpoint, subscriptionId string, auth autorest.Authorizer) {
	namespacesClient := notificationhubs.NewNamespacesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&namespacesClient.Client, auth)
	c.notificationNamespacesClient = namespacesClient

	notificationHubsClient := notificationhubs.NewClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&notificationHubsClient.Client, auth)
	c.notificationHubsClient = notificationHubsClient
}

func (c *ArmClient) registerOperationalInsightsClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	opwc := operationalinsights.NewWorkspacesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&opwc.Client, auth)
	c.workspacesClient = opwc

	solutionsClient := operationsmanagement.NewSolutionsClientWithBaseURI(endpoint, subscriptionId, "Microsoft.OperationsManagement", "solutions", "testing")
	c.configureClient(&solutionsClient.Client, auth)
	c.solutionsClient = solutionsClient

	lsClient := operationalinsights.NewLinkedServicesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&lsClient.Client, auth)
	c.linkedServicesClient = lsClient
}

func (c *ArmClient) registerRecoveryServiceClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	vaultsClient := recoveryservices.NewVaultsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&vaultsClient.Client, auth)
	c.recoveryServicesVaultsClient = vaultsClient

	protectedItemsClient := backup.NewProtectedItemsGroupClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&protectedItemsClient.Client, auth)
	c.recoveryServicesProtectedItemsClient = protectedItemsClient

	protectionPoliciesClient := backup.NewProtectionPoliciesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&protectionPoliciesClient.Client, auth)
	c.recoveryServicesProtectionPoliciesClient = protectionPoliciesClient
}

func (c *ArmClient) registerRedisClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	redisClient := redis.NewClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&redisClient.Client, auth)
	c.redisClient = redisClient

	firewallRuleClient := redis.NewFirewallRulesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&firewallRuleClient.Client, auth)
	c.redisFirewallClient = firewallRuleClient

	patchSchedulesClient := redis.NewPatchSchedulesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&patchSchedulesClient.Client, auth)
	c.redisPatchSchedulesClient = patchSchedulesClient
}

func (c *ArmClient) registerRelayClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	relayNamespacesClient := relay.NewNamespacesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&relayNamespacesClient.Client, auth)
	c.relayNamespacesClient = relayNamespacesClient
}

func (c *ArmClient) registerResourcesClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	locksClient := locks.NewManagementLocksClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&locksClient.Client, auth)
	c.managementLocksClient = locksClient

	deploymentsClient := resources.NewDeploymentsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&deploymentsClient.Client, auth)
	c.deploymentsClient = deploymentsClient

	resourcesClient := resources.NewClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&resourcesClient.Client, auth)
	c.resourcesClient = resourcesClient

	resourceGroupsClient := resources.NewGroupsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&resourceGroupsClient.Client, auth)
	c.resourceGroupsClient = resourceGroupsClient

	subscriptionsClient := subscriptions.NewClientWithBaseURI(endpoint)
	c.configureClient(&subscriptionsClient.Client, auth)
	c.subscriptionsClient = subscriptionsClient

	// this has to come from the Profile since this is shared with Stack
	providersClient := resourcesprofile.NewProvidersClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&providersClient.Client, auth)
	c.providersClient = providersClient
}

func (c *ArmClient) registerSchedulerClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	jobCollectionsClient := scheduler.NewJobCollectionsClientWithBaseURI(endpoint, subscriptionId) //nolint: megacheck
	c.configureClient(&jobCollectionsClient.Client, auth)
	c.schedulerJobCollectionsClient = jobCollectionsClient

	jobsClient := scheduler.NewJobsClientWithBaseURI(endpoint, subscriptionId) //nolint: megacheck
	c.configureClient(&jobsClient.Client, auth)
	c.schedulerJobsClient = jobsClient
}

func (c *ArmClient) registerSearchClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	searchClient := search.NewServicesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&searchClient.Client, auth)
	c.searchServicesClient = searchClient

	searchAdminKeysClient := search.NewAdminKeysClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&searchAdminKeysClient.Client, auth)
	c.searchAdminKeysClient = searchAdminKeysClient
}

func (c *ArmClient) registerSecurityCenterClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	ascLocation := "Global"

	securityCenterPricingClient := security.NewPricingsClientWithBaseURI(endpoint, subscriptionId, ascLocation)
	c.configureClient(&securityCenterPricingClient.Client, auth)
	c.securityCenterPricingClient = securityCenterPricingClient

	securityCenterContactsClient := security.NewContactsClientWithBaseURI(endpoint, subscriptionId, ascLocation)
	c.configureClient(&securityCenterContactsClient.Client, auth)
	c.securityCenterContactsClient = securityCenterContactsClient

	securityCenterWorkspaceClient := security.NewWorkspaceSettingsClientWithBaseURI(endpoint, subscriptionId, ascLocation)
	c.configureClient(&securityCenterWorkspaceClient.Client, auth)
	c.securityCenterWorkspaceClient = securityCenterWorkspaceClient
}

func (c *ArmClient) registerServiceBusClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	queuesClient := servicebus.NewQueuesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&queuesClient.Client, auth)
	c.serviceBusQueuesClient = queuesClient

	namespacesClient := servicebus.NewNamespacesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&namespacesClient.Client, auth)
	c.serviceBusNamespacesClient = namespacesClient

	topicsClient := servicebus.NewTopicsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&topicsClient.Client, auth)
	c.serviceBusTopicsClient = topicsClient

	subscriptionsClient := servicebus.NewSubscriptionsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&subscriptionsClient.Client, auth)
	c.serviceBusSubscriptionsClient = subscriptionsClient

	subscriptionRulesClient := servicebus.NewRulesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&subscriptionRulesClient.Client, auth)
	c.serviceBusSubscriptionRulesClient = subscriptionRulesClient
}

func (c *ArmClient) registerServiceFabricClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	clustersClient := servicefabric.NewClustersClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&clustersClient.Client, auth)
	c.serviceFabricClustersClient = clustersClient
}

func (c *ArmClient) registerSignalRClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	sc := signalr.NewClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&sc.Client, auth)
	c.signalRClient = sc
}

func (c *ArmClient) registerStorageClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	accountsClient := storage.NewAccountsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&accountsClient.Client, auth)
	c.storageServiceClient = accountsClient

	usageClient := storage.NewUsageClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&usageClient.Client, auth)
	c.storageUsageClient = usageClient
}

func (c *ArmClient) registerStreamAnalyticsClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	functionsClient := streamanalytics.NewFunctionsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&functionsClient.Client, auth)
	c.streamAnalyticsFunctionsClient = functionsClient

	jobsClient := streamanalytics.NewStreamingJobsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&jobsClient.Client, auth)
	c.streamAnalyticsJobsClient = jobsClient

	inputsClient := streamanalytics.NewInputsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&inputsClient.Client, auth)
	c.streamAnalyticsInputsClient = inputsClient

	outputsClient := streamanalytics.NewOutputsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&outputsClient.Client, auth)
	c.streamAnalyticsOutputsClient = outputsClient

	transformationsClient := streamanalytics.NewTransformationsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&transformationsClient.Client, auth)
	c.streamAnalyticsTransformationsClient = transformationsClient
}

func (c *ArmClient) registerTrafficManagerClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	endpointsClient := trafficmanager.NewEndpointsClientWithBaseURI(endpoint, c.subscriptionId)
	c.configureClient(&endpointsClient.Client, auth)
	c.trafficManagerEndpointsClient = endpointsClient

	geographicalHierarchiesClient := trafficmanager.NewGeographicHierarchiesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&geographicalHierarchiesClient.Client, auth)
	c.trafficManagerGeographialHierarchiesClient = geographicalHierarchiesClient

	profilesClient := trafficmanager.NewProfilesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&profilesClient.Client, auth)
	c.trafficManagerProfilesClient = profilesClient
}

func (c *ArmClient) registerWebClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	appServicePlansClient := web.NewAppServicePlansClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&appServicePlansClient.Client, auth)
	c.appServicePlansClient = appServicePlansClient

	appsClient := web.NewAppsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&appsClient.Client, auth)
	c.appServicesClient = appsClient
}

func (c *ArmClient) registerPolicyClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	policyAssignmentsClient := policy.NewAssignmentsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&policyAssignmentsClient.Client, auth)
	c.policyAssignmentsClient = policyAssignmentsClient

	policyDefinitionsClient := policy.NewDefinitionsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&policyDefinitionsClient.Client, auth)
	c.policyDefinitionsClient = policyDefinitionsClient

	policySetDefinitionsClient := policy.NewSetDefinitionsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&policySetDefinitionsClient.Client, auth)
	c.policySetDefinitionsClient = policySetDefinitionsClient
}

func (c *ArmClient) registerManagementGroupClients(endpoint string, auth autorest.Authorizer) {
	managementGroupsClient := managementgroups.NewClientWithBaseURI(endpoint)
	c.configureClient(&managementGroupsClient.Client, auth)
	c.managementGroupsClient = managementGroupsClient

	managementGroupsSubscriptionClient := managementgroups.NewSubscriptionsClientWithBaseURI(endpoint)
	c.configureClient(&managementGroupsSubscriptionClient.Client, auth)
	c.managementGroupsSubscriptionClient = managementGroupsSubscriptionClient
}

var (
	storageKeyCacheMu sync.RWMutex
	storageKeyCache   = make(map[string]string)
)

func (c *ArmClient) getKeyForStorageAccount(ctx context.Context, resourceGroupName, storageAccountName string) (string, bool, error) {
	cacheIndex := resourceGroupName + "/" + storageAccountName
	storageKeyCacheMu.RLock()
	key, ok := storageKeyCache[cacheIndex]
	storageKeyCacheMu.RUnlock()

	if ok {
		return key, true, nil
	}

	storageKeyCacheMu.Lock()
	defer storageKeyCacheMu.Unlock()
	key, ok = storageKeyCache[cacheIndex]
	if !ok {
		accountKeys, err := c.storageServiceClient.ListKeys(ctx, resourceGroupName, storageAccountName)
		if utils.ResponseWasNotFound(accountKeys.Response) {
			return "", false, nil
		}
		if err != nil {
			// We assume this is a transient error rather than a 404 (which is caught above),  so assume the
			// storeAccount still exists.
			return "", true, fmt.Errorf("Error retrieving keys for storage storeAccount %q: %s", storageAccountName, err)
		}

		if accountKeys.Keys == nil {
			return "", false, fmt.Errorf("Nil key returned for storage storeAccount %q", storageAccountName)
		}

		keys := *accountKeys.Keys
		if len(keys) <= 0 {
			return "", false, fmt.Errorf("No keys returned for storage storeAccount %q", storageAccountName)
		}

		keyPtr := keys[0].Value
		if keyPtr == nil {
			return "", false, fmt.Errorf("The first key returned is nil for storage storeAccount %q", storageAccountName)
		}

		key = *keyPtr
		storageKeyCache[cacheIndex] = key
	}

	return key, true, nil
}

func (c *ArmClient) getBlobStorageClientForStorageAccount(ctx context.Context, resourceGroupName, storageAccountName string) (*mainStorage.BlobStorageClient, bool, error) {
	key, accountExists, err := c.getKeyForStorageAccount(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return nil, accountExists, err
	}
	if !accountExists {
		return nil, false, nil
	}

	storageClient, err := mainStorage.NewClient(storageAccountName, key, c.environment.StorageEndpointSuffix,
		mainStorage.DefaultAPIVersion, true)
	if err != nil {
		return nil, true, fmt.Errorf("Error creating storage client for storage storeAccount %q: %s", storageAccountName, err)
	}

	blobClient := storageClient.GetBlobService()
	return &blobClient, true, nil
}

func (c *ArmClient) getFileServiceClientForStorageAccount(ctx context.Context, resourceGroupName, storageAccountName string) (*mainStorage.FileServiceClient, bool, error) {
	key, accountExists, err := c.getKeyForStorageAccount(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return nil, accountExists, err
	}
	if !accountExists {
		return nil, false, nil
	}

	storageClient, err := mainStorage.NewClient(storageAccountName, key, c.environment.StorageEndpointSuffix,
		mainStorage.DefaultAPIVersion, true)
	if err != nil {
		return nil, true, fmt.Errorf("Error creating storage client for storage storeAccount %q: %s", storageAccountName, err)
	}

	fileClient := storageClient.GetFileService()
	return &fileClient, true, nil
}

func (c *ArmClient) getTableServiceClientForStorageAccount(ctx context.Context, resourceGroupName, storageAccountName string) (*mainStorage.TableServiceClient, bool, error) {
	key, accountExists, err := c.getKeyForStorageAccount(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return nil, accountExists, err
	}
	if !accountExists {
		return nil, false, nil
	}

	storageClient, err := mainStorage.NewClient(storageAccountName, key, c.environment.StorageEndpointSuffix,
		mainStorage.DefaultAPIVersion, true)
	if err != nil {
		return nil, true, fmt.Errorf("Error creating storage client for storage storeAccount %q: %s", storageAccountName, err)
	}

	tableClient := storageClient.GetTableService()
	return &tableClient, true, nil
}

func (c *ArmClient) getQueueServiceClientForStorageAccount(ctx context.Context, resourceGroupName, storageAccountName string) (*mainStorage.QueueServiceClient, bool, error) {
	key, accountExists, err := c.getKeyForStorageAccount(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return nil, accountExists, err
	}
	if !accountExists {
		return nil, false, nil
	}

	storageClient, err := mainStorage.NewClient(storageAccountName, key, c.environment.StorageEndpointSuffix,
		mainStorage.DefaultAPIVersion, true)
	if err != nil {
		return nil, true, fmt.Errorf("Error creating storage client for storage storeAccount %q: %s", storageAccountName, err)
	}

	queueClient := storageClient.GetQueueService()
	return &queueClient, true, nil
}
