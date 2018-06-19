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
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2017-04-02/cdn"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2017-12-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2018-04-01/containerinstance"
	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2017-10-01/containerregistry"
	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2017-09-30/containerservice"
	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb"
	"github.com/Azure/azure-sdk-for-go/services/datalake/store/mgmt/2016-11-01/account"
	"github.com/Azure/azure-sdk-for-go/services/dns/mgmt/2016-04-01/dns"
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/mgmt/2018-01-01/eventgrid"
	"github.com/Azure/azure-sdk-for-go/services/eventhub/mgmt/2017-04-01/eventhub"
	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/azure-sdk-for-go/services/iothub/mgmt/2017-07-01/devices"
	keyVault "github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2016-10-01/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/monitor/mgmt/2018-03-01/insights"
	"github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2017-12-01/mysql"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2017-12-01/postgresql"
	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2018-01-01-preview/authorization"
	"github.com/Azure/azure-sdk-for-go/services/preview/operationalinsights/mgmt/2015-11-01-preview/operationalinsights"
	"github.com/Azure/azure-sdk-for-go/services/preview/operationsmanagement/mgmt/2015-11-01-preview/operationsmanagement"
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
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2017-10-01/storage"
	"github.com/Azure/azure-sdk-for-go/services/trafficmanager/mgmt/2017-05-01/trafficmanager"
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2016-09-01/web"
	mainStorage "github.com/Azure/azure-sdk-for-go/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/authentication"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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

	automationAccountClient    automation.AccountClient
	automationRunbookClient    automation.RunbookClient
	automationCredentialClient automation.CredentialClient
	automationScheduleClient   automation.ScheduleClient

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

	// Application Insights
	appInsightsClient appinsights.ComponentsClient

	// Authentication
	roleAssignmentsClient   authorization.RoleAssignmentsClient
	roleDefinitionsClient   authorization.RoleDefinitionsClient
	servicePrincipalsClient graphrbac.ServicePrincipalsClient

	// CDN
	cdnCustomDomainsClient cdn.CustomDomainsClient
	cdnEndpointsClient     cdn.EndpointsClient
	cdnProfilesClient      cdn.ProfilesClient

	// Compute
	availSetClient         compute.AvailabilitySetsClient
	diskClient             compute.DisksClient
	imageClient            compute.ImagesClient
	snapshotsClient        compute.SnapshotsClient
	usageOpsClient         compute.UsageClient
	vmExtensionImageClient compute.VirtualMachineExtensionImagesClient
	vmExtensionClient      compute.VirtualMachineExtensionsClient
	vmScaleSetClient       compute.VirtualMachineScaleSetsClient
	vmImageClient          compute.VirtualMachineImagesClient
	vmClient               compute.VirtualMachinesClient

	// Devices
	iothubResourceClient devices.IotHubResourceClient

	// Databases
	mysqlConfigurationsClient            mysql.ConfigurationsClient
	mysqlDatabasesClient                 mysql.DatabasesClient
	mysqlFirewallRulesClient             mysql.FirewallRulesClient
	mysqlServersClient                   mysql.ServersClient
	postgresqlConfigurationsClient       postgresql.ConfigurationsClient
	postgresqlDatabasesClient            postgresql.DatabasesClient
	postgresqlFirewallRulesClient        postgresql.FirewallRulesClient
	postgresqlServersClient              postgresql.ServersClient
	sqlDatabasesClient                   sql.DatabasesClient
	sqlElasticPoolsClient                sql.ElasticPoolsClient
	sqlFirewallRulesClient               sql.FirewallRulesClient
	sqlServersClient                     sql.ServersClient
	sqlServerAzureADAdministratorsClient sql.ServerAzureADAdministratorsClient
	sqlVirtualNetworkRulesClient         sql.VirtualNetworkRulesClient

	// Data Lake Store
	dataLakeStoreAccountClient account.AccountsClient

	// KeyVault
	keyVaultClient           keyvault.VaultsClient
	keyVaultManagementClient keyVault.BaseClient

	// Monitor
	monitorAlertRulesClient insights.AlertRulesClient

	// Networking
	applicationGatewayClient        network.ApplicationGatewaysClient
	applicationSecurityGroupsClient network.ApplicationSecurityGroupsClient
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

	// Search
	searchServicesClient search.ServicesClient

	// ServiceBus
	serviceBusQueuesClient            servicebus.QueuesClient
	serviceBusNamespacesClient        servicebus.NamespacesClient
	serviceBusTopicsClient            servicebus.TopicsClient
	serviceBusSubscriptionsClient     servicebus.SubscriptionsClient
	serviceBusSubscriptionRulesClient servicebus.RulesClient

	//Scheduler
	schedulerJobCollectionsClient scheduler.JobCollectionsClient

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

func (c *ArmClient) configureClient(client *autorest.Client, auth autorest.Authorizer) {
	setUserAgent(client)
	client.Authorizer = auth
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

	client.registerAppInsightsClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerAutomationClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerAuthentication(endpoint, graphEndpoint, c.SubscriptionID, c.TenantID, auth, graphAuth, sender)
	client.registerCDNClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerComputeClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerContainerInstanceClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerContainerRegistryClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerContainerServicesClients(endpoint, c.SubscriptionID, auth)
	client.registerCosmosDBClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerDatabases(endpoint, c.SubscriptionID, auth, sender)
	client.registerDataLakeStoreAccountClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerDeviceClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerDNSClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerEventGridClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerEventHubClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerKeyVaultClients(endpoint, c.SubscriptionID, auth, keyVaultAuth, sender)
	client.registerMonitorClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerNetworkingClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerOperationalInsightsClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerRecoveryServiceClients(endpoint, c.SubscriptionID, auth)
	client.registerRedisClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerRelayClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerResourcesClients(endpoint, c.SubscriptionID, auth)
	client.registerSearchClients(endpoint, c.SubscriptionID, auth)
	client.registerServiceBusClients(endpoint, c.SubscriptionID, auth)
	client.registerSchedulerClients(endpoint, c.SubscriptionID, auth)
	client.registerStorageClients(endpoint, c.SubscriptionID, auth)
	client.registerTrafficManagerClients(endpoint, c.SubscriptionID, auth)
	client.registerWebClients(endpoint, c.SubscriptionID, auth)
	client.registerPolicyClients(endpoint, c.SubscriptionID, auth)

	return &client, nil
}

func (c *ArmClient) registerAppInsightsClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	ai := appinsights.NewComponentsClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&ai.Client)
	ai.Authorizer = auth
	ai.Sender = sender
	ai.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.appInsightsClient = ai
}

func (c *ArmClient) registerAutomationClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	accountClient := automation.NewAccountClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&accountClient.Client)
	accountClient.Authorizer = auth
	accountClient.Sender = sender
	accountClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.automationAccountClient = accountClient

	credentialClient := automation.NewCredentialClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&credentialClient.Client)
	credentialClient.Authorizer = auth
	credentialClient.Sender = sender
	credentialClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.automationCredentialClient = credentialClient

	runbookClient := automation.NewRunbookClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&runbookClient.Client)
	runbookClient.Authorizer = auth
	runbookClient.Sender = sender
	runbookClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.automationRunbookClient = runbookClient

	scheduleClient := automation.NewScheduleClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&scheduleClient.Client)
	scheduleClient.Authorizer = auth
	scheduleClient.Sender = sender
	scheduleClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.automationScheduleClient = scheduleClient
}

func (c *ArmClient) registerAuthentication(endpoint, graphEndpoint, subscriptionId, tenantId string, auth, graphAuth autorest.Authorizer, sender autorest.Sender) {
	assignmentsClient := authorization.NewRoleAssignmentsClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&assignmentsClient.Client)
	assignmentsClient.Authorizer = auth
	assignmentsClient.Sender = sender
	assignmentsClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.roleAssignmentsClient = assignmentsClient

	definitionsClient := authorization.NewRoleDefinitionsClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&definitionsClient.Client)
	definitionsClient.Authorizer = auth
	definitionsClient.Sender = sender
	definitionsClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.roleDefinitionsClient = definitionsClient

	servicePrincipalsClient := graphrbac.NewServicePrincipalsClientWithBaseURI(graphEndpoint, tenantId)
	setUserAgent(&servicePrincipalsClient.Client)
	servicePrincipalsClient.Authorizer = graphAuth
	servicePrincipalsClient.Sender = sender
	servicePrincipalsClient.SkipResourceProviderRegistration = c.skipProviderRegistration
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

func (c *ArmClient) registerDatabases(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	// MySQL
	mysqlConfigClient := mysql.NewConfigurationsClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&mysqlConfigClient.Client)
	mysqlConfigClient.Authorizer = auth
	mysqlConfigClient.Sender = sender
	mysqlConfigClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.mysqlConfigurationsClient = mysqlConfigClient

	mysqlDBClient := mysql.NewDatabasesClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&mysqlDBClient.Client)
	mysqlDBClient.Authorizer = auth
	mysqlDBClient.Sender = sender
	mysqlDBClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.mysqlDatabasesClient = mysqlDBClient

	mysqlFWClient := mysql.NewFirewallRulesClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&mysqlFWClient.Client)
	mysqlFWClient.Authorizer = auth
	mysqlFWClient.Sender = sender
	mysqlFWClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.mysqlFirewallRulesClient = mysqlFWClient

	mysqlServersClient := mysql.NewServersClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&mysqlServersClient.Client)
	mysqlServersClient.Authorizer = auth
	mysqlServersClient.Sender = sender
	mysqlServersClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.mysqlServersClient = mysqlServersClient

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

	// SQL Azure
	sqlDBClient := sql.NewDatabasesClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&sqlDBClient.Client)
	sqlDBClient.Authorizer = auth
	sqlDBClient.Sender = sender
	sqlDBClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.sqlDatabasesClient = sqlDBClient

	sqlFWClient := sql.NewFirewallRulesClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&sqlFWClient.Client)
	sqlFWClient.Authorizer = auth
	sqlFWClient.Sender = sender
	sqlFWClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.sqlFirewallRulesClient = sqlFWClient

	sqlEPClient := sql.NewElasticPoolsClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&sqlEPClient.Client)
	sqlEPClient.Authorizer = auth
	sqlEPClient.Sender = sender
	sqlEPClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.sqlElasticPoolsClient = sqlEPClient

	sqlSrvClient := sql.NewServersClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&sqlSrvClient.Client)
	sqlSrvClient.Authorizer = auth
	sqlSrvClient.Sender = sender
	sqlSrvClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.sqlServersClient = sqlSrvClient

	sqlADClient := sql.NewServerAzureADAdministratorsClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&sqlADClient.Client)
	sqlADClient.Authorizer = auth
	sqlADClient.Sender = sender
	sqlADClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.sqlServerAzureADAdministratorsClient = sqlADClient

	sqlVNRClient := sql.NewVirtualNetworkRulesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&sqlVNRClient.Client, auth)
	c.sqlVirtualNetworkRulesClient = sqlVNRClient
}

func (c *ArmClient) registerDataLakeStoreAccountClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	dataLakeStoreAccountClient := account.NewAccountsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&dataLakeStoreAccountClient.Client, auth)
	c.dataLakeStoreAccountClient = dataLakeStoreAccountClient
}

func (c *ArmClient) registerDeviceClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	iotClient := devices.NewIotHubResourceClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&iotClient.Client, auth)
	c.iothubResourceClient = iotClient
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
	setUserAgent(&egtc.Client)
	egtc.Authorizer = auth
	egtc.Sender = sender
	egtc.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.eventGridTopicsClient = egtc
}

func (c *ArmClient) registerEventHubClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	ehc := eventhub.NewEventHubsClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&ehc.Client)
	ehc.Authorizer = auth
	ehc.Sender = sender
	ehc.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.eventHubClient = ehc

	chcgc := eventhub.NewConsumerGroupsClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&chcgc.Client)
	chcgc.Authorizer = auth
	chcgc.Sender = sender
	chcgc.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.eventHubConsumerGroupClient = chcgc

	ehnc := eventhub.NewNamespacesClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&ehnc.Client)
	ehnc.Authorizer = auth
	ehnc.Sender = sender
	ehnc.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.eventHubNamespacesClient = ehnc
}

func (c *ArmClient) registerKeyVaultClients(endpoint, subscriptionId string, auth autorest.Authorizer, keyVaultAuth autorest.Authorizer, sender autorest.Sender) {
	keyVaultClient := keyvault.NewVaultsClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&keyVaultClient.Client)
	keyVaultClient.Authorizer = auth
	keyVaultClient.Sender = sender
	keyVaultClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.keyVaultClient = keyVaultClient

	keyVaultManagementClient := keyVault.New()
	setUserAgent(&keyVaultManagementClient.Client)
	keyVaultManagementClient.Authorizer = keyVaultAuth
	keyVaultManagementClient.Sender = sender
	keyVaultManagementClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.keyVaultManagementClient = keyVaultManagementClient
}

func (c *ArmClient) registerMonitorClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	arc := insights.NewAlertRulesClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&arc.Client)
	arc.Authorizer = auth
	arc.Sender = autorest.CreateSender(withRequestLogging())
	c.monitorAlertRulesClient = arc
}

func (c *ArmClient) registerNetworkingClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	applicationGatewaysClient := network.NewApplicationGatewaysClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&applicationGatewaysClient.Client, auth)
	c.applicationGatewayClient = applicationGatewaysClient

	appSecurityGroupsClient := network.NewApplicationSecurityGroupsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&appSecurityGroupsClient.Client, auth)
	c.applicationSecurityGroupsClient = appSecurityGroupsClient

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

	watchersClient := network.NewWatchersClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&watchersClient.Client, auth)
	c.watcherClient = watchersClient
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

func (c *ArmClient) registerSchedulerClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	jobsClient := scheduler.NewJobCollectionsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&jobsClient.Client, auth)
	c.schedulerJobCollectionsClient = jobsClient
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
			// account still exists.
			return "", true, fmt.Errorf("Error retrieving keys for storage account %q: %s", storageAccountName, err)
		}

		if accountKeys.Keys == nil {
			return "", false, fmt.Errorf("Nil key returned for storage account %q", storageAccountName)
		}

		keys := *accountKeys.Keys
		if len(keys) <= 0 {
			return "", false, fmt.Errorf("No keys returned for storage account %q", storageAccountName)
		}

		keyPtr := keys[0].Value
		if keyPtr == nil {
			return "", false, fmt.Errorf("The first key returned is nil for storage account %q", storageAccountName)
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
		return nil, true, fmt.Errorf("Error creating storage client for storage account %q: %s", storageAccountName, err)
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
		return nil, true, fmt.Errorf("Error creating storage client for storage account %q: %s", storageAccountName, err)
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
		return nil, true, fmt.Errorf("Error creating storage client for storage account %q: %s", storageAccountName, err)
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
		return nil, true, fmt.Errorf("Error creating storage client for storage account %q: %s", storageAccountName, err)
	}

	queueClient := storageClient.GetQueueService()
	return &queueClient, true, nil
}
