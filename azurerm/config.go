package azurerm

import (
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/go-azure-helpers/sender"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/analysisservices"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/applicationinsights"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/authorization"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/automation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/batch"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/bot"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cognitive"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/containers"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databricks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datafactory"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datalake"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/devspace"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/devtestlabs"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dns"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventgrid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/frontdoor"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/graph"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hdinsight"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iothub"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/logic"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managementgroup"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/maps"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mariadb"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/media"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/monitor"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mysql"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/notificationhub"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/privatedns"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/recoveryservices"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/redis"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/relay"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/scheduler"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/search"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicebus"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicefabric"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/signalr"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sql"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/streamanalytics"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/subscription"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/trafficmanager"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web"
)

// ArmClient contains the handles to all the specific Azure Resource Manager
// resource classes' respective clients.
type ArmClient struct {
	// inherit the fields from the parent, so that we should be able to set/access these at either level
	clients.Client

	clientId                 string
	tenantId                 string
	subscriptionId           string
	partnerId                string
	usingServicePrincipal    bool
	environment              azure.Environment
	skipProviderRegistration bool

	// Services
	// NOTE: all new services should be Public as they're going to be relocated in the near-future
	analysisservices *analysisservices.Client
	apiManagement    *apimanagement.Client
	appInsights      *applicationinsights.Client
	automation       *automation.Client
	authorization    *authorization.Client
	batch            *batch.Client
	bot              *bot.Client
	cdn              *cdn.Client
	cognitive        *cognitive.Client
	compute          *clients.ComputeClient
	containers       *containers.Client
	cosmos           *cosmos.Client
	databricks       *databricks.Client
	dataFactory      *datafactory.Client
	datalake         *datalake.Client
	devSpace         *devspace.Client
	devTestLabs      *devtestlabs.Client
	dns              *dns.Client
	privateDns       *privatedns.Client
	eventGrid        *eventgrid.Client
	eventhub         *eventhub.Client
	frontdoor        *frontdoor.Client
	graph            *graph.Client
	hdinsight        *hdinsight.Client
	iothub           *iothub.Client
	keyvault         *keyvault.Client
	kusto            *kusto.Client
	logAnalytics     *loganalytics.Client
	logic            *logic.Client
	managementGroups *managementgroup.Client
	maps             *maps.Client
	mariadb          *mariadb.Client
	media            *media.Client
	monitor          *monitor.Client
	mysql            *mysql.Client
	msi              *msi.Client
	mssql            *mssql.Client
	network          *network.Client
	notificationHubs *notificationhub.Client
	policy           *policy.Client
	postgres         *postgres.Client
	recoveryServices *recoveryservices.Client
	redis            *redis.Client
	relay            *relay.Client
	resource         *resource.Client
	Scheduler        *scheduler.Client
	Search           *search.Client
	SecurityCenter   *securitycenter.Client
	ServiceBus       *servicebus.Client
	ServiceFabric    *servicefabric.Client
	SignalR          *signalr.Client
	Storage          *storage.Client
	StreamAnalytics  *streamanalytics.Client
	Subscription     *subscription.Client
	Sql              *sql.Client
	TrafficManager   *trafficmanager.Client
	web              *web.Client
}

// getArmClient is a helper method which returns a fully instantiated
// *ArmClient based on the Config's current settings.
func getArmClient(authConfig *authentication.Config, skipProviderRegistration bool, partnerId string, disableCorrelationRequestID bool) (*ArmClient, error) {
	env, err := authentication.DetermineEnvironment(authConfig.Environment)
	if err != nil {
		return nil, err
	}

	// client declarations:
	client := ArmClient{
		Client: clients.Client{},

		clientId:                 authConfig.ClientID,
		tenantId:                 authConfig.TenantID,
		subscriptionId:           authConfig.SubscriptionID,
		partnerId:                partnerId,
		environment:              *env,
		usingServicePrincipal:    authConfig.AuthenticatedAsAServicePrincipal,
		skipProviderRegistration: skipProviderRegistration,
	}

	oauthConfig, err := authConfig.BuildOAuthConfig(env.ActiveDirectoryEndpoint)
	if err != nil {
		return nil, err
	}

	// OAuthConfigForTenant returns a pointer, which can be nil.
	if oauthConfig == nil {
		return nil, fmt.Errorf("Unable to configure OAuthConfig for tenant %s", authConfig.TenantID)
	}

	sender := sender.BuildSender("AzureRM")

	// Resource Manager endpoints
	endpoint := env.ResourceManagerEndpoint
	auth, err := authConfig.GetAuthorizationToken(sender, oauthConfig, env.TokenAudience)
	if err != nil {
		return nil, err
	}

	// Graph Endpoints
	graphEndpoint := env.GraphEndpoint
	graphAuth, err := authConfig.GetAuthorizationToken(sender, oauthConfig, graphEndpoint)
	if err != nil {
		return nil, err
	}

	// Storage Endpoints
	storageAuth := authConfig.BearerAuthorizerCallback(sender, oauthConfig)

	// Key Vault Endpoints
	keyVaultAuth := authConfig.BearerAuthorizerCallback(sender, oauthConfig)

	o := &common.ClientOptions{
		SubscriptionId:              authConfig.SubscriptionID,
		TenantID:                    authConfig.TenantID,
		PartnerId:                   partnerId,
		GraphAuthorizer:             graphAuth,
		GraphEndpoint:               graphEndpoint,
		KeyVaultAuthorizer:          keyVaultAuth,
		ResourceManagerAuthorizer:   auth,
		ResourceManagerEndpoint:     endpoint,
		StorageAuthorizer:           storageAuth,
		PollingDuration:             180 * time.Minute,
		SkipProviderReg:             skipProviderRegistration,
		DisableCorrelationRequestID: disableCorrelationRequestID,
		Environment:                 *env,
	}

	client.analysisservices = analysisservices.BuildClient(o)
	client.apiManagement = apimanagement.BuildClient(o)
	client.appInsights = applicationinsights.BuildClient(o)
	client.automation = automation.BuildClient(o)
	client.authorization = authorization.BuildClient(o)
	client.batch = batch.BuildClient(o)
	client.bot = bot.BuildClient(o)
	client.cdn = cdn.BuildClient(o)
	client.cognitive = cognitive.BuildClient(o)
	client.compute = clients.NewComputeClient(o)
	client.containers = containers.BuildClient(o)
	client.cosmos = cosmos.BuildClient(o)
	client.databricks = databricks.BuildClient(o)
	client.dataFactory = datafactory.BuildClient(o)
	client.datalake = datalake.BuildClient(o)
	client.devSpace = devspace.BuildClient(o)
	client.devTestLabs = devtestlabs.BuildClient(o)
	client.dns = dns.BuildClient(o)
	client.eventGrid = eventgrid.BuildClient(o)
	client.eventhub = eventhub.BuildClient(o)
	client.frontdoor = frontdoor.BuildClient(o)
	client.graph = graph.BuildClient(o)
	client.hdinsight = hdinsight.BuildClient(o)
	client.iothub = iothub.BuildClient(o)
	client.keyvault = keyvault.BuildClient(o)
	client.kusto = kusto.BuildClient(o)
	client.logic = logic.BuildClient(o)
	client.logAnalytics = loganalytics.BuildClient(o)
	client.maps = maps.BuildClient(o)
	client.mariadb = mariadb.BuildClient(o)
	client.media = media.BuildClient(o)
	client.monitor = monitor.BuildClient(o)
	client.mssql = mssql.BuildClient(o)
	client.msi = msi.BuildClient(o)
	client.mysql = mysql.BuildClient(o)
	client.managementGroups = managementgroup.BuildClient(o)
	client.network = network.BuildClient(o)
	client.notificationHubs = notificationhub.BuildClient(o)
	client.policy = policy.BuildClient(o)
	client.postgres = postgres.BuildClient(o)
	client.privateDns = privatedns.BuildClient(o)
	client.recoveryServices = recoveryservices.BuildClient(o)
	client.redis = redis.BuildClient(o)
	client.relay = relay.BuildClient(o)
	client.resource = resource.BuildClient(o)
	client.Search = search.BuildClient(o)
	client.SecurityCenter = securitycenter.BuildClient(o)
	client.ServiceBus = servicebus.BuildClient(o)
	client.ServiceFabric = servicefabric.BuildClient(o)
	client.Scheduler = scheduler.BuildClient(o)
	client.SignalR = signalr.BuildClient(o)
	client.StreamAnalytics = streamanalytics.BuildClient(o)
	client.Storage = storage.BuildClient(o)
	client.Subscription = subscription.BuildClient(o)
	client.Sql = sql.BuildClient(o)
	client.TrafficManager = trafficmanager.BuildClient(o)
	client.web = web.BuildClient(o)

	return &client, nil
}
