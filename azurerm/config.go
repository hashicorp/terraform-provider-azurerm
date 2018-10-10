package azurerm

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"sync"
	"time"

	appinsights "github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2015-05-01/insights"
	"github.com/Azure/azure-sdk-for-go/services/automation/mgmt/2015-10-31/automation"
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2017-10-12/cdn"
	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/mgmt/2017-04-18/cognitiveservices"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-06-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2018-04-01/containerinstance"
	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2017-10-01/containerregistry"
	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2018-03-31/containerservice"
	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb"
	"github.com/Azure/azure-sdk-for-go/services/databricks/mgmt/2018-04-01/databricks"
	analyticsAccount "github.com/Azure/azure-sdk-for-go/services/datalake/analytics/mgmt/2016-11-01/account"
	"github.com/Azure/azure-sdk-for-go/services/datalake/store/2016-11-01/filesystem"
	storeAccount "github.com/Azure/azure-sdk-for-go/services/datalake/store/mgmt/2016-11-01/account"
	"github.com/Azure/azure-sdk-for-go/services/devtestlabs/mgmt/2016-05-15/dtl"
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/mgmt/2018-01-01/eventgrid"
	"github.com/Azure/azure-sdk-for-go/services/eventhub/mgmt/2017-04-01/eventhub"
	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/azure-sdk-for-go/services/iothub/mgmt/2018-04-01/devices"
	keyVault "github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2018-02-14/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/logic/mgmt/2016-06-01/logic"
	"github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2017-12-01/mysql"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-04-01/network"
	"github.com/Azure/azure-sdk-for-go/services/notificationhubs/mgmt/2017-04-01/notificationhubs"
	"github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2017-12-01/postgresql"
	"github.com/Azure/azure-sdk-for-go/services/preview/apimanagement/mgmt/2018-06-01-preview/apimanagement"
	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2018-01-01-preview/authorization"
	"github.com/Azure/azure-sdk-for-go/services/preview/dns/mgmt/2018-03-01-preview/dns"
	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-03-01/insights"
	"github.com/Azure/azure-sdk-for-go/services/preview/msi/mgmt/2015-08-31-preview/msi"
	"github.com/Azure/azure-sdk-for-go/services/preview/operationalinsights/mgmt/2015-11-01-preview/operationalinsights"
	"github.com/Azure/azure-sdk-for-go/services/preview/operationsmanagement/mgmt/2015-11-01-preview/operationsmanagement"
	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2018-03-01-preview/management"
	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2015-05-01-preview/sql"
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2016-06-01/recoveryservices"
	"github.com/Azure/azure-sdk-for-go/services/redis/mgmt/2018-03-01/redis"
	"github.com/Azure/azure-sdk-for-go/services/relay/mgmt/2017-04-01/relay"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2016-06-01/subscriptions"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2016-09-01/locks"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2016-12-01/policy"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"
	"github.com/Azure/azure-sdk-for-go/services/scheduler/mgmt/2016-03-01/scheduler"
	"github.com/Azure/azure-sdk-for-go/services/search/mgmt/2015-08-19/search"
	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/Azure/azure-sdk-for-go/services/servicefabric/mgmt/2018-02-01/servicefabric"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2017-10-01/storage"
	"github.com/Azure/azure-sdk-for-go/services/trafficmanager/mgmt/2017-05-01/trafficmanager"
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2018-02-01/web"

	mainStorage "github.com/Azure/azure-sdk-for-go/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/authentication"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/Azure/azure-sdk-for-go/services/preview/devspaces/mgmt/2018-06-01-preview/devspaces"
)

// ArmClient contains the handles to all the specific Azure Resource Manager
// resource classes' respective clients.
type ArmClient struct {
	clientId                 string
	tenantId                 string
	subscriptionId           string
	usingServicePrincipal    bool
	environment              azure.Environment
	skipProviderRegistration bool

	StopContext context.Context

	cosmosDBClient documentdb.DatabaseAccountsClient

	automationAccountClient      automation.AccountClient
	automationRunbookClient      automation.RunbookClient
	automationCredentialClient   automation.CredentialClient
	automationScheduleClient     automation.ScheduleClient
	automationRunbookDraftClient automation.RunbookDraftClient

	dnsClient   dns.RecordSetsClient
	zonesClient dns.ZonesClient

	containerRegistryClient  containerregistry.RegistriesClient
	containerServicesClient  containerservice.ContainerServicesClient
	kubernetesClustersClient containerservice.ManagedClustersClient
	containerGroupsClient    containerinstance.ContainerGroupsClient

	eventGridTopicsClient       eventgrid.TopicsClient
	eventHubClient              eventhub.EventHubsClient
	eventHubConsumerGroupClient eventhub.ConsumerGroupsClient
	eventHubNamespacesClient    eventhub.NamespacesClient

	workspacesClient operationalinsights.WorkspacesClient
	solutionsClient  operationsmanagement.SolutionsClient

	redisClient               redis.Client
	redisFirewallClient       redis.FirewallRulesClient
	redisPatchSchedulesClient redis.PatchSchedulesClient

	// API Management
	apiManagementServiceClient apimanagement.ServiceClient

	// Application Insights
	appInsightsClient appinsights.ComponentsClient

	// Authentication
	roleAssignmentsClient   authorization.RoleAssignmentsClient
	roleDefinitionsClient   authorization.RoleDefinitionsClient
	applicationsClient      graphrbac.ApplicationsClient
	servicePrincipalsClient graphrbac.ServicePrincipalsClient

	// Autoscale Settings
	autoscaleSettingsClient insights.AutoscaleSettingsClient

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
	devTestVirtualNetworksClient dtl.VirtualNetworksClient

	// DevSpaces
	devSpacesControllersClient devspaces.ControllersClient

	// Databases
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
	sqlFirewallRulesClient                   sql.FirewallRulesClient
	sqlServersClient                         sql.ServersClient
	sqlServerAzureADAdministratorsClient     sql.ServerAzureADAdministratorsClient
	sqlVirtualNetworkRulesClient             sql.VirtualNetworkRulesClient

	// Data Lake Store
	dataLakeStoreAccountClient       storeAccount.AccountsClient
	dataLakeStoreFirewallRulesClient storeAccount.FirewallRulesClient
	dataLakeStoreFilesClient         filesystem.Client

	// Data Lake Analytics
	dataLakeAnalyticsAccountClient       analyticsAccount.AccountsClient
	dataLakeAnalyticsFirewallRulesClient analyticsAccount.FirewallRulesClient

	// Databricks
	databricksWorkspacesClient databricks.WorkspacesClient

	// KeyVault
	keyVaultClient           keyvault.VaultsClient
	keyVaultManagementClient keyVault.BaseClient

	// Logic
	logicWorkflowsClient logic.WorkflowsClient

	// Management Groups
	managementGroupsClient             managementgroups.Client
	managementGroupsSubscriptionClient managementgroups.SubscriptionsClient

	// Monitor
	monitorActionGroupsClient insights.ActionGroupsClient
	monitorAlertRulesClient   insights.AlertRulesClient

	// MSI
	userAssignedIdentitiesClient msi.UserAssignedIdentitiesClient

	// Networking
	applicationGatewayClient        network.ApplicationGatewaysClient
	applicationSecurityGroupsClient network.ApplicationSecurityGroupsClient
	azureFirewallsClient            network.AzureFirewallsClient
	expressRouteAuthsClient         network.ExpressRouteCircuitAuthorizationsClient
	expressRouteCircuitClient       network.ExpressRouteCircuitsClient
	expressRoutePeeringsClient      network.ExpressRouteCircuitPeeringsClient
	ifaceClient                     network.InterfacesClient
	loadBalancerClient              network.LoadBalancersClient
	localNetConnClient              network.LocalNetworkGatewaysClient
	packetCapturesClient            network.PacketCapturesClient
	publicIPClient                  network.PublicIPAddressesClient
	routesClient                    network.RoutesClient
	routeTablesClient               network.RouteTablesClient
	secGroupClient                  network.SecurityGroupsClient
	secRuleClient                   network.SecurityRulesClient
	subnetClient                    network.SubnetsClient
	netUsageClient                  network.UsagesClient
	vnetGatewayConnectionsClient    network.VirtualNetworkGatewayConnectionsClient
	vnetGatewayClient               network.VirtualNetworkGatewaysClient
	vnetClient                      network.VirtualNetworksClient
	vnetPeeringsClient              network.VirtualNetworkPeeringsClient
	watcherClient                   network.WatchersClient

	// Notification Hubs
	notificationHubsClient       notificationhubs.Client
	notificationNamespacesClient notificationhubs.NamespacesClient

	// Recovery Services
	recoveryServicesVaultsClient recoveryservices.VaultsClient

	// Relay
	relayNamespacesClient relay.NamespacesClient

	// Resources
	managementLocksClient locks.ManagementLocksClient
	deploymentsClient     resources.DeploymentsClient
	providersClient       resources.ProvidersClient
	resourcesClient       resources.Client
	resourceGroupsClient  resources.GroupsClient
	subscriptionsClient   subscriptions.Client

	// Scheduler
	schedulerJobCollectionsClient scheduler.JobCollectionsClient
	schedulerJobsClient           scheduler.JobsClient

	// Search
	searchServicesClient search.ServicesClient

	// ServiceBus
	serviceBusQueuesClient            servicebus.QueuesClient
	serviceBusNamespacesClient        servicebus.NamespacesClient
	serviceBusTopicsClient            servicebus.TopicsClient
	serviceBusSubscriptionsClient     servicebus.SubscriptionsClient
	serviceBusSubscriptionRulesClient servicebus.RulesClient

	// Service Fabric
	serviceFabricClustersClient servicefabric.ClustersClient

	// Storage
	storageServiceClient storage.AccountsClient
	storageUsageClient   storage.UsageClient

	// Traffic Manager
	trafficManagerGeographialHierarchiesClient trafficmanager.GeographicHierarchiesClient
	trafficManagerProfilesClient               trafficmanager.ProfilesClient
	trafficManagerEndpointsClient              trafficmanager.EndpointsClient

	// Web
	appServicePlansClient web.AppServicePlansClient
	appServicesClient     web.AppsClient

	// Policy
	policyAssignmentsClient policy.AssignmentsClient
	policyDefinitionsClient policy.DefinitionsClient
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
	setUserAgent(client)
	client.Authorizer = auth
	//client.RequestInspector = azure.WithClientID(clientRequestID())
	client.Sender = autorest.CreateSender(withRequestLogging())
	client.SkipResourceProviderRegistration = c.skipProviderRegistration
	client.PollingDuration = 60 * time.Minute
}

func withRequestLogging() autorest.SendDecorator {
	return func(s autorest.Sender) autorest.Sender {
		return autorest.SenderFunc(func(r *http.Request) (*http.Response, error) {
			// dump request to wire format
			if dump, err := httputil.DumpRequestOut(r, true); err == nil {
				log.Printf("[DEBUG] AzureRM Request: \n%s\n", dump)
			} else {
				// fallback to basic message
				log.Printf("[DEBUG] AzureRM Request: %s to %s\n", r.Method, r.URL)
			}

			resp, err := s.Do(r)
			if resp != nil {
				// dump response to wire format
				if dump, err := httputil.DumpResponse(resp, true); err == nil {
					log.Printf("[DEBUG] AzureRM Response for %s: \n%s\n", r.URL, dump)
				} else {
					// fallback to basic message
					log.Printf("[DEBUG] AzureRM Response: %s for %s\n", resp.Status, r.URL)
				}
			} else {
				log.Printf("[DEBUG] Request to %s completed with no response", r.URL)
			}
			return resp, err
		})
	}
}

func setUserAgent(client *autorest.Client) {
	tfVersion := fmt.Sprintf("HashiCorp-Terraform-v%s", terraform.VersionString())

	// if the user agent already has a value append the Terraform user agent string
	if curUserAgent := client.UserAgent; curUserAgent != "" {
		client.UserAgent = fmt.Sprintf("%s;%s", curUserAgent, tfVersion)
	} else {
		client.UserAgent = tfVersion
	}

	// append the CloudShell version to the user agent if it exists
	if azureAgent := os.Getenv("AZURE_HTTP_USER_AGENT"); azureAgent != "" {
		client.UserAgent = fmt.Sprintf("%s;%s", client.UserAgent, azureAgent)
	}

	log.Printf("[DEBUG] AzureRM Client User Agent: %s\n", client.UserAgent)
}

func getAuthorizationToken(c *authentication.Config, oauthConfig *adal.OAuthConfig, endpoint string) (*autorest.BearerAuthorizer, error) {
	useServicePrincipal := c.ClientSecret != ""

	if useServicePrincipal {
		spt, err := adal.NewServicePrincipalToken(*oauthConfig, c.ClientID, c.ClientSecret, endpoint)
		if err != nil {
			return nil, err
		}

		auth := autorest.NewBearerAuthorizer(spt)
		return auth, nil
	}

	if c.UseMsi {
		spt, err := adal.NewServicePrincipalTokenFromMSI(c.MsiEndpoint, endpoint)
		if err != nil {
			return nil, err
		}
		auth := autorest.NewBearerAuthorizer(spt)
		return auth, nil
	}

	if c.IsCloudShell {
		// load the refreshed tokens from the Azure CLI
		err := c.LoadTokensFromAzureCLI()
		if err != nil {
			return nil, fmt.Errorf("Error loading the refreshed CloudShell tokens: %+v", err)
		}
	}

	spt, err := adal.NewServicePrincipalTokenFromManualToken(*oauthConfig, c.ClientID, endpoint, *c.AccessToken)
	if err != nil {
		return nil, err
	}

	err = spt.Refresh()

	if err != nil {
		return nil, fmt.Errorf("Error refreshing Service Principal Token: %+v", err)
	}

	auth := autorest.NewBearerAuthorizer(spt)
	return auth, nil
}

// getArmClient is a helper method which returns a fully instantiated
// *ArmClient based on the Config's current settings.
func getArmClient(c *authentication.Config) (*ArmClient, error) {
	// detect cloud from environment
	env, envErr := azure.EnvironmentFromName(c.Environment)
	if envErr != nil {
		// try again with wrapped value to support readable values like german instead of AZUREGERMANCLOUD
		wrapped := fmt.Sprintf("AZURE%sCLOUD", c.Environment)
		var innerErr error
		if env, innerErr = azure.EnvironmentFromName(wrapped); innerErr != nil {
			return nil, envErr
		}
	}

	// client declarations:
	client := ArmClient{
		clientId:                 c.ClientID,
		tenantId:                 c.TenantID,
		subscriptionId:           c.SubscriptionID,
		environment:              env,
		usingServicePrincipal:    c.ClientSecret != "",
		skipProviderRegistration: c.SkipProviderRegistration,
	}

	oauthConfig, err := adal.NewOAuthConfig(env.ActiveDirectoryEndpoint, c.TenantID)
	if err != nil {
		return nil, err
	}

	// OAuthConfigForTenant returns a pointer, which can be nil.
	if oauthConfig == nil {
		return nil, fmt.Errorf("Unable to configure OAuthConfig for tenant %s", c.TenantID)
	}

	sender := autorest.CreateSender(withRequestLogging())

	// Resource Manager endpoints
	endpoint := env.ResourceManagerEndpoint
	auth, err := getAuthorizationToken(c, oauthConfig, endpoint)
	if err != nil {
		return nil, err
	}

	// Graph Endpoints
	graphEndpoint := env.GraphEndpoint
	graphAuth, err := getAuthorizationToken(c, oauthConfig, graphEndpoint)
	if err != nil {
		return nil, err
	}

	// Key Vault Endpoints
	keyVaultAuth := autorest.NewBearerAuthorizerCallback(sender, func(tenantID, resource string) (*autorest.BearerAuthorizer, error) {
		keyVaultSpt, err := getAuthorizationToken(c, oauthConfig, resource)
		if err != nil {
			return nil, err
		}

		return keyVaultSpt, nil
	})

	client.registerApiManagementServiceClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerAppInsightsClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerAutomationClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerAuthentication(endpoint, graphEndpoint, c.SubscriptionID, c.TenantID, auth, graphAuth, sender)
	client.registerCDNClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerCognitiveServiceClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerComputeClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerContainerInstanceClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerContainerRegistryClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerContainerServicesClients(endpoint, c.SubscriptionID, auth)
	client.registerCosmosDBClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerDatabricksClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerDatabases(endpoint, c.SubscriptionID, auth, sender)
	client.registerDataLakeStoreClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerDeviceClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerDevSpacesClients(endpoint, c.SubscriptionID, auth)
	client.registerDevTestClients(endpoint, c.SubscriptionID, auth)
	client.registerDNSClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerEventGridClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerEventHubClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerKeyVaultClients(endpoint, c.SubscriptionID, auth, keyVaultAuth, sender)
	client.registerLogicClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerMonitorClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerNetworkingClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerNotificationHubsClient(endpoint, c.SubscriptionID, auth, sender)
	client.registerOperationalInsightsClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerRecoveryServiceClients(endpoint, c.SubscriptionID, auth)
	client.registerPolicyClients(endpoint, c.SubscriptionID, auth)
	client.registerManagementGroupClients(endpoint, auth)
	client.registerRedisClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerRelayClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerResourcesClients(endpoint, c.SubscriptionID, auth)
	client.registerSearchClients(endpoint, c.SubscriptionID, auth)
	client.registerServiceBusClients(endpoint, c.SubscriptionID, auth)
	client.registerServiceFabricClients(endpoint, c.SubscriptionID, auth)
	client.registerSchedulerClients(endpoint, c.SubscriptionID, auth)
	client.registerStorageClients(endpoint, c.SubscriptionID, auth)
	client.registerTrafficManagerClients(endpoint, c.SubscriptionID, auth)
	client.registerWebClients(endpoint, c.SubscriptionID, auth)

	return &client, nil
}

func (c *ArmClient) registerApiManagementServiceClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	ams := apimanagement.NewServiceClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&ams.Client, auth)
	c.apiManagementServiceClient = ams
}

func (c *ArmClient) registerAppInsightsClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	ai := appinsights.NewComponentsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&ai.Client, auth)
	c.appInsightsClient = ai
}

func (c *ArmClient) registerAutomationClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	accountClient := automation.NewAccountClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&accountClient.Client, auth)
	c.automationAccountClient = accountClient

	credentialClient := automation.NewCredentialClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&credentialClient.Client, auth)
	c.automationCredentialClient = credentialClient

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

func (c *ArmClient) registerAuthentication(endpoint, graphEndpoint, subscriptionId, tenantId string, auth, graphAuth autorest.Authorizer, sender autorest.Sender) {
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

func (c *ArmClient) registerCDNClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
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

func (c *ArmClient) registerCognitiveServiceClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	accountsClient := cognitiveservices.NewAccountsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&accountsClient.Client, auth)
	c.cognitiveAccountsClient = accountsClient
}

func (c *ArmClient) registerCosmosDBClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	cdb := documentdb.NewDatabaseAccountsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&cdb.Client, auth)
	c.cosmosDBClient = cdb
}

func (c *ArmClient) registerComputeClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
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

func (c *ArmClient) registerContainerInstanceClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	cgc := containerinstance.NewContainerGroupsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&cgc.Client, auth)
	c.containerGroupsClient = cgc
}

func (c *ArmClient) registerContainerRegistryClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	crc := containerregistry.NewRegistriesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&crc.Client, auth)
	c.containerRegistryClient = crc
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

func (c *ArmClient) registerDatabricksClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	databricksWorkspacesClient := databricks.NewWorkspacesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&databricksWorkspacesClient.Client, auth)
	c.databricksWorkspacesClient = databricksWorkspacesClient
}

func (c *ArmClient) registerDatabases(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
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
	setUserAgent(&sqlDTDPClient.Client)
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

func (c *ArmClient) registerDataLakeStoreClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
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

func (c *ArmClient) registerDeviceClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	iotClient := devices.NewIotHubResourceClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&iotClient.Client, auth)
	c.iothubResourceClient = iotClient
}

func (c *ArmClient) registerDevTestClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	labsClient := dtl.NewLabsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&labsClient.Client, auth)
	c.devTestLabsClient = labsClient

	devTestVirtualNetworksClient := dtl.NewVirtualNetworksClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&devTestVirtualNetworksClient.Client, auth)
	c.devTestVirtualNetworksClient = devTestVirtualNetworksClient
}

func (c *ArmClient) registerDevSpacesClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	controllersClient := devspaces.NewControllersClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&controllersClient.Client, auth)
	c.devSpacesControllersClient = controllersClient
}

func (c *ArmClient) registerDNSClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	dn := dns.NewRecordSetsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&dn.Client, auth)
	c.dnsClient = dn

	zo := dns.NewZonesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&zo.Client, auth)
	c.zonesClient = zo
}

func (c *ArmClient) registerEventGridClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	egtc := eventgrid.NewTopicsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&egtc.Client, auth)
	c.eventGridTopicsClient = egtc
}

func (c *ArmClient) registerEventHubClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
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

func (c *ArmClient) registerKeyVaultClients(endpoint, subscriptionId string, auth autorest.Authorizer, keyVaultAuth autorest.Authorizer, sender autorest.Sender) {
	keyVaultClient := keyvault.NewVaultsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&keyVaultClient.Client, auth)
	c.keyVaultClient = keyVaultClient

	keyVaultManagementClient := keyVault.New()
	c.configureClient(&keyVaultManagementClient.Client, keyVaultAuth)
	c.keyVaultManagementClient = keyVaultManagementClient
}

func (c *ArmClient) registerLogicClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	workflowsClient := logic.NewWorkflowsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&workflowsClient.Client, auth)
	c.logicWorkflowsClient = workflowsClient
}

func (c *ArmClient) registerMonitorClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	agc := insights.NewActionGroupsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&agc.Client, auth)
	c.monitorActionGroupsClient = agc

	arc := insights.NewAlertRulesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&arc.Client, auth)
	c.monitorAlertRulesClient = arc

	autoscaleSettingsClient := insights.NewAutoscaleSettingsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&autoscaleSettingsClient.Client, auth)
	c.autoscaleSettingsClient = autoscaleSettingsClient
}

func (c *ArmClient) registerNetworkingClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	applicationGatewaysClient := network.NewApplicationGatewaysClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&applicationGatewaysClient.Client, auth)
	c.applicationGatewayClient = applicationGatewaysClient

	appSecurityGroupsClient := network.NewApplicationSecurityGroupsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&appSecurityGroupsClient.Client, auth)
	c.applicationSecurityGroupsClient = appSecurityGroupsClient

	azureFirewallsClient := network.NewAzureFirewallsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&azureFirewallsClient.Client, auth)
	c.azureFirewallsClient = azureFirewallsClient

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

func (c *ArmClient) registerNotificationHubsClient(endpoint, subscriptionId string, auth *autorest.BearerAuthorizer, sender autorest.Sender) {
	namespacesClient := notificationhubs.NewNamespacesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&namespacesClient.Client, auth)
	c.notificationNamespacesClient = namespacesClient

	notificationHubsClient := notificationhubs.NewClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&notificationHubsClient.Client, auth)
	c.notificationHubsClient = notificationHubsClient
}

func (c *ArmClient) registerOperationalInsightsClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	opwc := operationalinsights.NewWorkspacesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&opwc.Client, auth)
	c.workspacesClient = opwc

	solutionsClient := operationsmanagement.NewSolutionsClientWithBaseURI(endpoint, subscriptionId, "Microsoft.OperationsManagement", "solutions", "testing")
	c.configureClient(&solutionsClient.Client, auth)
	c.solutionsClient = solutionsClient
}

func (c *ArmClient) registerRecoveryServiceClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	vaultsClient := recoveryservices.NewVaultsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&vaultsClient.Client, auth)
	c.recoveryServicesVaultsClient = vaultsClient
}

func (c *ArmClient) registerRedisClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
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

func (c *ArmClient) registerRelayClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
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

	providersClient := resources.NewProvidersClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&providersClient.Client, auth)
	c.providersClient = providersClient

	subscriptionsClient := subscriptions.NewClientWithBaseURI(endpoint)
	c.configureClient(&subscriptionsClient.Client, auth)
	c.subscriptionsClient = subscriptionsClient
}

func (c *ArmClient) registerSchedulerClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	jobCollectionsClient := scheduler.NewJobCollectionsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&jobCollectionsClient.Client, auth)
	c.schedulerJobCollectionsClient = jobCollectionsClient

	jobsClient := scheduler.NewJobsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&jobsClient.Client, auth)
	c.schedulerJobsClient = jobsClient
}

func (c *ArmClient) registerSearchClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	searchClient := search.NewServicesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&searchClient.Client, auth)
	c.searchServicesClient = searchClient
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

func (c *ArmClient) registerStorageClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	accountsClient := storage.NewAccountsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&accountsClient.Client, auth)
	c.storageServiceClient = accountsClient

	usageClient := storage.NewUsageClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&usageClient.Client, auth)
	c.storageUsageClient = usageClient
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

func (armClient *ArmClient) getKeyForStorageAccount(ctx context.Context, resourceGroupName, storageAccountName string) (string, bool, error) {
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
		accountKeys, err := armClient.storageServiceClient.ListKeys(ctx, resourceGroupName, storageAccountName)
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

func (armClient *ArmClient) getBlobStorageClientForStorageAccount(ctx context.Context, resourceGroupName, storageAccountName string) (*mainStorage.BlobStorageClient, bool, error) {
	key, accountExists, err := armClient.getKeyForStorageAccount(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return nil, accountExists, err
	}
	if accountExists == false {
		return nil, false, nil
	}

	storageClient, err := mainStorage.NewClient(storageAccountName, key, armClient.environment.StorageEndpointSuffix,
		mainStorage.DefaultAPIVersion, true)
	if err != nil {
		return nil, true, fmt.Errorf("Error creating storage client for storage storeAccount %q: %s", storageAccountName, err)
	}

	blobClient := storageClient.GetBlobService()
	return &blobClient, true, nil
}

func (armClient *ArmClient) getFileServiceClientForStorageAccount(ctx context.Context, resourceGroupName, storageAccountName string) (*mainStorage.FileServiceClient, bool, error) {
	key, accountExists, err := armClient.getKeyForStorageAccount(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return nil, accountExists, err
	}
	if accountExists == false {
		return nil, false, nil
	}

	storageClient, err := mainStorage.NewClient(storageAccountName, key, armClient.environment.StorageEndpointSuffix,
		mainStorage.DefaultAPIVersion, true)
	if err != nil {
		return nil, true, fmt.Errorf("Error creating storage client for storage storeAccount %q: %s", storageAccountName, err)
	}

	fileClient := storageClient.GetFileService()
	return &fileClient, true, nil
}

func (armClient *ArmClient) getTableServiceClientForStorageAccount(ctx context.Context, resourceGroupName, storageAccountName string) (*mainStorage.TableServiceClient, bool, error) {
	key, accountExists, err := armClient.getKeyForStorageAccount(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return nil, accountExists, err
	}
	if accountExists == false {
		return nil, false, nil
	}

	storageClient, err := mainStorage.NewClient(storageAccountName, key, armClient.environment.StorageEndpointSuffix,
		mainStorage.DefaultAPIVersion, true)
	if err != nil {
		return nil, true, fmt.Errorf("Error creating storage client for storage storeAccount %q: %s", storageAccountName, err)
	}

	tableClient := storageClient.GetTableService()
	return &tableClient, true, nil
}

func (armClient *ArmClient) getQueueServiceClientForStorageAccount(ctx context.Context, resourceGroupName, storageAccountName string) (*mainStorage.QueueServiceClient, bool, error) {
	key, accountExists, err := armClient.getKeyForStorageAccount(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return nil, accountExists, err
	}
	if accountExists == false {
		return nil, false, nil
	}

	storageClient, err := mainStorage.NewClient(storageAccountName, key, armClient.environment.StorageEndpointSuffix,
		mainStorage.DefaultAPIVersion, true)
	if err != nil {
		return nil, true, fmt.Errorf("Error creating storage client for storage storeAccount %q: %s", storageAccountName, err)
	}

	queueClient := storageClient.GetQueueService()
	return &queueClient, true, nil
}
