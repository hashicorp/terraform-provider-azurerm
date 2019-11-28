package azurerm

import (
	"context"
	"fmt"
	"time"

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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/healthcare"
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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/netapp"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/notificationhub"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/portal"
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
}

type armClientBuilder struct {
	authConfig                  *authentication.Config
	skipProviderRegistration    bool
	terraformVersion            string
	partnerId                   string
	disableCorrelationRequestID bool
	disableTerraformPartnerID   bool
}

// getArmClient is a helper method which returns a fully instantiated
// *ArmClient based on the Config's current settings.
func getArmClient(ctx context.Context, builder armClientBuilder) (*ArmClient, error) {
	env, err := authentication.DetermineEnvironment(builder.authConfig.Environment)
	if err != nil {
		return nil, err
	}

	// client declarations:
	account, err := clients.NewResourceManagerAccount(ctx, *builder.authConfig, *env)
	if err != nil {
		return nil, fmt.Errorf("Error building account: %+v", err)
	}

	client := ArmClient{
		Client: clients.Client{
			Account: account,
		},
	}

	oauthConfig, err := builder.authConfig.BuildOAuthConfig(env.ActiveDirectoryEndpoint)
	if err != nil {
		return nil, err
	}

	// OAuthConfigForTenant returns a pointer, which can be nil.
	if oauthConfig == nil {
		return nil, fmt.Errorf("Unable to configure OAuthConfig for tenant %s", builder.authConfig.TenantID)
	}

	sender := sender.BuildSender("AzureRM")

	// Resource Manager endpoints
	endpoint := env.ResourceManagerEndpoint
	auth, err := builder.authConfig.GetAuthorizationToken(sender, oauthConfig, env.TokenAudience)
	if err != nil {
		return nil, err
	}

	// Graph Endpoints
	graphEndpoint := env.GraphEndpoint
	graphAuth, err := builder.authConfig.GetAuthorizationToken(sender, oauthConfig, graphEndpoint)
	if err != nil {
		return nil, err
	}

	// Storage Endpoints
	storageAuth, err := builder.authConfig.GetAuthorizationToken(sender, oauthConfig, env.ResourceIdentifiers.Storage)
	if err != nil {
		return nil, err
	}

	// Key Vault Endpoints
	keyVaultAuth := builder.authConfig.BearerAuthorizerCallback(sender, oauthConfig)

	o := &common.ClientOptions{
		SubscriptionId:              builder.authConfig.SubscriptionID,
		TenantID:                    builder.authConfig.TenantID,
		PartnerId:                   builder.partnerId,
		TerraformVersion:            builder.terraformVersion,
		GraphAuthorizer:             graphAuth,
		GraphEndpoint:               graphEndpoint,
		KeyVaultAuthorizer:          keyVaultAuth,
		ResourceManagerAuthorizer:   auth,
		ResourceManagerEndpoint:     endpoint,
		StorageAuthorizer:           storageAuth,
		PollingDuration:             180 * time.Minute,
		SkipProviderReg:             builder.skipProviderRegistration,
		DisableCorrelationRequestID: builder.disableCorrelationRequestID,
		DisableTerraformPartnerID:   builder.disableTerraformPartnerID,
		Environment:                 *env,
	}

	client.AnalysisServices = analysisservices.BuildClient(o)
	client.ApiManagement = apimanagement.BuildClient(o)
	client.AppInsights = applicationinsights.BuildClient(o)
	client.Automation = automation.BuildClient(o)
	client.Authorization = authorization.BuildClient(o)
	client.Batch = batch.BuildClient(o)
	client.Bot = bot.BuildClient(o)
	client.Cdn = cdn.BuildClient(o)
	client.Cognitive = cognitive.BuildClient(o)
	client.Compute = clients.NewComputeClient(o)
	client.Containers = containers.BuildClient(o)
	client.Cosmos = cosmos.BuildClient(o)
	client.DataBricks = databricks.BuildClient(o)
	client.DataFactory = datafactory.BuildClient(o)
	client.Datalake = datalake.BuildClient(o)
	client.DevSpace = devspace.BuildClient(o)
	client.DevTestLabs = devtestlabs.BuildClient(o)
	client.Dns = dns.BuildClient(o)
	client.EventGrid = eventgrid.BuildClient(o)
	client.Eventhub = eventhub.BuildClient(o)
	client.Frontdoor = frontdoor.BuildClient(o)
	client.Graph = graph.BuildClient(o)
	client.HDInsight = hdinsight.BuildClient(o)
	client.Healthcare = healthcare.BuildClient(o)
	client.IoTHub = iothub.BuildClient(o)
	client.KeyVault = keyvault.BuildClient(o)
	client.Kusto = kusto.BuildClient(o)
	client.Logic = logic.BuildClient(o)
	client.LogAnalytics = loganalytics.BuildClient(o)
	client.Maps = maps.BuildClient(o)
	client.MariaDB = mariadb.BuildClient(o)
	client.Media = media.BuildClient(o)
	client.Monitor = monitor.BuildClient(o)
	client.Mssql = mssql.BuildClient(o)
	client.Msi = msi.BuildClient(o)
	client.Mysql = mysql.BuildClient(o)
	client.ManagementGroups = managementgroup.BuildClient(o)
	client.Netapp = netapp.BuildClient(o)
	client.Network = network.BuildClient(o)
	client.NotificationHubs = notificationhub.BuildClient(o)
	client.Policy = policy.BuildClient(o)
	client.Portal = portal.BuildClient(o)
	client.Postgres = postgres.BuildClient(o)
	client.PrivateDns = privatedns.BuildClient(o)
	client.RecoveryServices = recoveryservices.BuildClient(o)
	client.Redis = redis.BuildClient(o)
	client.Relay = relay.BuildClient(o)
	client.Resource = resource.BuildClient(o)
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
	client.Web = web.BuildClient(o)

	return &client, nil
}
