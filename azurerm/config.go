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
	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2018-12-01/batch"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-06-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb"
	analyticsAccount "github.com/Azure/azure-sdk-for-go/services/datalake/analytics/mgmt/2016-11-01/account"
	"github.com/Azure/azure-sdk-for-go/services/datalake/store/2016-11-01/filesystem"
	storeAccount "github.com/Azure/azure-sdk-for-go/services/datalake/store/mgmt/2016-11-01/account"
	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	keyVault "github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2018-02-14/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/mariadb/mgmt/2018-06-01/mariadb"
	"github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2017-12-01/mysql"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
	"github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2017-12-01/postgresql"
	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2018-01-01-preview/authorization"
	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-03-01/insights"
	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2015-05-01-preview/sql"
	MsSql "github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-10-01-preview/sql"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2016-09-01/locks"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-06-01/subscriptions"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-04-01/storage"
	"github.com/Azure/azure-sdk-for-go/services/streamanalytics/mgmt/2016-03-01/streamanalytics"
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2018-02-01/web"
	mainStorage "github.com/Azure/azure-sdk-for-go/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/go-azure-helpers/sender"
	"github.com/hashicorp/terraform/httpclient"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/applicationinsights"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/automation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cognitive"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/containers"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databricks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datafactory"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/devspace"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/devtestlabs"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dns"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventgrid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hdinsight"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iothub"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/logic"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managementgroup"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/media"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/notificationhub"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/privatedns"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/recoveryservices"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/redis"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/relay"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/scheduler"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/search"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicebus"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicefabric"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/signalr"
	intStor "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/trafficmanager"
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
	environment              azure.Environment
	skipProviderRegistration bool

	StopContext context.Context

	// Services
	apiManagement    *apimanagement.Client
	appInsights      *applicationinsights.Client
	automation       *automation.Client
	cdn              *cdn.Client
	cognitive        *cognitive.Client
	containers       *containers.Client
	databricks       *databricks.Client
	dataFactory      *datafactory.Client
	devSpace         *devspace.Client
	devTestLabs      *devtestlabs.Client
	dns              *dns.Client
	privateDns       *privatedns.Client
	eventGrid        *eventgrid.Client
	eventhub         *eventhub.Client
	hdinsight        *hdinsight.Client
	iothub           *iothub.Client
	logAnalytics     *loganalytics.Client
	logic            *logic.Client
	managementGroups *managementgroup.Client
	media            *media.Client
	msi              *msi.Client
	notificationHubs *notificationhub.Client
	policy           *policy.Client
	recoveryServices *recoveryservices.Client
	redis            *redis.Client
	relay            *relay.Client
	scheduler        *scheduler.Client
	search           *search.Client
	securityCenter   *securitycenter.Client
	servicebus       *servicebus.Client
	serviceFabric    *servicefabric.Client
	signalr          *signalr.Client
	storage          *intStor.Client
	trafficManager   *trafficmanager.Client

	// TODO: refactor
	cosmosAccountsClient documentdb.DatabaseAccountsClient

	// Authentication
	roleAssignmentsClient   authorization.RoleAssignmentsClient
	roleDefinitionsClient   authorization.RoleDefinitionsClient
	applicationsClient      graphrbac.ApplicationsClient
	servicePrincipalsClient graphrbac.ServicePrincipalsClient

	// Autoscale Settings
	autoscaleSettingsClient insights.AutoscaleSettingsClient

	// Batch
	batchAccountClient     batch.AccountClient
	batchApplicationClient batch.ApplicationClient
	batchCertificateClient batch.CertificateClient
	batchPoolClient        batch.PoolClient

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

	// Databases
	mariadbDatabasesClient                   mariadb.DatabasesClient
	mariadbFirewallRulesClient               mariadb.FirewallRulesClient
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

	// Data Lake Store
	dataLakeStoreAccountClient       storeAccount.AccountsClient
	dataLakeStoreFirewallRulesClient storeAccount.FirewallRulesClient
	dataLakeStoreFilesClient         filesystem.Client

	// Data Lake Analytics
	dataLakeAnalyticsAccountClient       analyticsAccount.AccountsClient
	dataLakeAnalyticsFirewallRulesClient analyticsAccount.FirewallRulesClient

	// KeyVault
	keyVaultClient           keyvault.VaultsClient
	keyVaultManagementClient keyVault.BaseClient

	// Monitor
	monitorActionGroupsClient               insights.ActionGroupsClient
	monitorActivityLogAlertsClient          insights.ActivityLogAlertsClient
	monitorAlertRulesClient                 insights.AlertRulesClient
	monitorDiagnosticSettingsClient         insights.DiagnosticSettingsClient
	monitorDiagnosticSettingsCategoryClient insights.DiagnosticSettingsCategoryClient
	monitorLogProfilesClient                insights.LogProfilesClient
	monitorMetricAlertsClient               insights.MetricAlertsClient

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
	netProfileClient                network.ProfilesClient
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

	// Resources
	managementLocksClient locks.ManagementLocksClient
	deploymentsClient     resources.DeploymentsClient
	providersClient       resourcesprofile.ProvidersClient
	resourcesClient       resources.Client
	resourceGroupsClient  resources.GroupsClient
	subscriptionsClient   subscriptions.Client

	// Storage
	storageServiceClient storage.AccountsClient
	storageUsageClient   storage.UsagesClient

	// Stream Analytics
	streamAnalyticsFunctionsClient       streamanalytics.FunctionsClient
	streamAnalyticsJobsClient            streamanalytics.StreamingJobsClient
	streamAnalyticsInputsClient          streamanalytics.InputsClient
	streamAnalyticsOutputsClient         streamanalytics.OutputsClient
	streamAnalyticsTransformationsClient streamanalytics.TransformationsClient

	// Web
	appServicePlansClient web.AppServicePlansClient
	appServicesClient     web.AppsClient
}

func (c *ArmClient) configureClient(client *autorest.Client, auth autorest.Authorizer) {
	setUserAgent(client, c.partnerId)
	client.Authorizer = auth
	client.RequestInspector = common.WithCorrelationRequestID(common.CorrelationRequestID())
	client.Sender = sender.BuildSender("AzureRM")
	client.SkipResourceProviderRegistration = c.skipProviderRegistration
	client.PollingDuration = 180 * time.Minute
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

	sender := sender.BuildSender("AzureRM")

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

	o := &common.ClientOptions{
		GraphAuthorizer:            graphAuth,
		GraphEndpoint:              graphEndpoint,
		KeyVaultAuthorizer:         keyVaultAuth,
		ResourceManagerAuthorizer:  auth,
		ResourceManagerEndpoint:    endpoint,
		SubscriptionId:             c.SubscriptionID,
		PartnerId:                  partnerId,
		PollingDuration:            60 * time.Minute,
		SkipProviderReg:            skipProviderRegistration,
		EnableCorrelationRequestID: true,
	}

	client.apiManagement = apimanagement.BuildClient(o)
	client.appInsights = applicationinsights.BuildClient(o)
	client.automation = automation.BuildClient(o)
	client.cdn = cdn.BuildClient(o)
	client.cognitive = cognitive.BuildClient(o)
	client.containers = containers.BuildClient(o)
	client.databricks = databricks.BuildClient(o)
	client.dataFactory = datafactory.BuildClient(o)
	client.devSpace = devspace.BuildClient(o)
	client.devTestLabs = devtestlabs.BuildClient(o)
	client.dns = dns.BuildClient(o)
	client.eventGrid = eventgrid.BuildClient(o)
	client.eventhub = eventhub.BuildClient(o)
	client.hdinsight = hdinsight.BuildClient(o)
	client.iothub = iothub.BuildClient(o)
	client.logic = logic.BuildClient(o)
	client.logAnalytics = loganalytics.BuildClient(o)
	client.media = media.BuildClient(o)
	client.msi = msi.BuildClient(o)
	client.managementGroups = managementgroup.BuildClient(o)
	client.notificationHubs = notificationhub.BuildClient(o)
	client.policy = policy.BuildClient(o)
	client.privateDns = privatedns.BuildClient(o)
	client.recoveryServices = recoveryservices.BuildClient(o)
	client.redis = redis.BuildClient(o)
	client.relay = relay.BuildClient(o)
	client.search = search.BuildClient(o)
	client.securityCenter = securitycenter.BuildClient(o)
	client.servicebus = servicebus.BuildClient(o)
	client.serviceFabric = servicefabric.BuildClient(o)
	client.scheduler = scheduler.BuildClient(o)
	client.signalr = signalr.BuildClient(o)
	client.trafficManager = trafficmanager.BuildClient(o)

	client.registerAuthentication(endpoint, graphEndpoint, c.SubscriptionID, c.TenantID, auth, graphAuth)
	client.registerBatchClients(endpoint, c.SubscriptionID, auth)
	client.registerComputeClients(endpoint, c.SubscriptionID, auth)
	client.registerCosmosAccountsClients(endpoint, c.SubscriptionID, auth)
	client.registerDatabases(endpoint, c.SubscriptionID, auth, sender)
	client.registerDataLakeStoreClients(endpoint, c.SubscriptionID, auth)
	client.registerKeyVaultClients(endpoint, c.SubscriptionID, auth, keyVaultAuth)
	client.registerMonitorClients(endpoint, c.SubscriptionID, auth)
	client.registerNetworkingClients(endpoint, c.SubscriptionID, auth)
	client.registerResourcesClients(endpoint, c.SubscriptionID, auth)
	client.registerStorageClients(endpoint, c.SubscriptionID, auth)
	client.registerStreamAnalyticsClients(endpoint, c.SubscriptionID, auth)
	client.registerWebClients(endpoint, c.SubscriptionID, auth)

	return &client, nil
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

	batchApplicationClient := batch.NewApplicationClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&batchApplicationClient.Client, auth)
	c.batchApplicationClient = batchApplicationClient

	batchCertificateClient := batch.NewCertificateClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&batchCertificateClient.Client, auth)
	c.batchCertificateClient = batchCertificateClient

	batchPool := batch.NewPoolClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&batchPool.Client, auth)
	c.batchPoolClient = batchPool
}

func (c *ArmClient) registerCosmosAccountsClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	ca := documentdb.NewDatabaseAccountsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&ca.Client, auth)
	c.cosmosAccountsClient = ca
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

func (c *ArmClient) registerDatabases(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	mariadbDBClient := mariadb.NewDatabasesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&mariadbDBClient.Client, auth)
	c.mariadbDatabasesClient = mariadbDBClient

	mariadbFWClient := mariadb.NewFirewallRulesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&mariadbFWClient.Client, auth)
	c.mariadbFirewallRulesClient = mariadbFWClient

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

func (c *ArmClient) registerKeyVaultClients(endpoint, subscriptionId string, auth autorest.Authorizer, keyVaultAuth autorest.Authorizer) {
	keyVaultClient := keyvault.NewVaultsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&keyVaultClient.Client, auth)
	c.keyVaultClient = keyVaultClient

	keyVaultManagementClient := keyVault.New()
	c.configureClient(&keyVaultManagementClient.Client, keyVaultAuth)
	c.keyVaultManagementClient = keyVaultManagementClient
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

	netProfileClient := network.NewProfilesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&netProfileClient.Client, auth)
	c.netProfileClient = netProfileClient

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

	watchersClient := network.NewWatchersClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&watchersClient.Client, auth)
	c.watcherClient = watchersClient
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

func (c *ArmClient) registerStorageClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	accountsClient := storage.NewAccountsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&accountsClient.Client, auth)
	c.storageServiceClient = accountsClient

	usageClient := storage.NewUsagesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&usageClient.Client, auth)
	c.storageUsageClient = usageClient

	c.storage = intStor.BuildClient(accountsClient)
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

func (c *ArmClient) registerWebClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	appServicePlansClient := web.NewAppServicePlansClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&appServicePlansClient.Client, auth)
	c.appServicePlansClient = appServicePlansClient

	appsClient := web.NewAppsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&appsClient.Client, auth)
	c.appServicesClient = appsClient
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
