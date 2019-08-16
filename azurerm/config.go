package azurerm

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-04-01/storage"
	mainStorage "github.com/Azure/azure-sdk-for-go/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/go-azure-helpers/sender"
	"github.com/hashicorp/terraform/httpclient"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/analysisservices"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/applicationinsights"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/authorization"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/automation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/batch"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cognitive"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute"
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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/graph"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hdinsight"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iothub"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault"
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
	intStor "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/streamanalytics"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/subscription"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/trafficmanager"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web"
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
	analysisservices *analysisservices.Client
	apiManagement    *apimanagement.Client
	appInsights      *applicationinsights.Client
	automation       *automation.Client
	authorization    *authorization.Client
	batch            *batch.Client
	cdn              *cdn.Client
	cognitive        *cognitive.Client
	compute          *compute.Client
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
	graph            *graph.Client
	hdinsight        *hdinsight.Client
	iothub           *iothub.Client
	keyvault         *keyvault.Client
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
	scheduler        *scheduler.Client
	search           *search.Client
	securityCenter   *securitycenter.Client
	servicebus       *servicebus.Client
	serviceFabric    *servicefabric.Client
	signalr          *signalr.Client
	storage          *intStor.Client
	streamanalytics  *streamanalytics.Client
	subscription     *subscription.Client
	sql              *sql.Client
	trafficManager   *trafficmanager.Client
	web              *web.Client

	// Storage
	storageServiceClient storage.AccountsClient
	storageUsageClient   storage.UsagesClient
}

// getArmClient is a helper method which returns a fully instantiated
// *ArmClient based on the Config's current settings.
func getArmClient(c *authentication.Config, skipProviderRegistration bool, partnerId string, disableCorrelationRequestID bool) (*ArmClient, error) {
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

	// Storage Endpoints
	storageEndpoint := env.ResourceIdentifiers.Storage
	storageAuth, err := c.GetAuthorizationToken(sender, oauthConfig, storageEndpoint)
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
		SubscriptionId:              c.SubscriptionID,
		TenantID:                    c.TenantID,
		PartnerId:                   partnerId,
		GraphAuthorizer:             graphAuth,
		GraphEndpoint:               graphEndpoint,
		KeyVaultAuthorizer:          keyVaultAuth,
		ResourceManagerAuthorizer:   auth,
		ResourceManagerEndpoint:     endpoint,
		StorageAuthorizer:           storageAuth,
		PollingDuration:             60 * time.Minute,
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
	client.cdn = cdn.BuildClient(o)
	client.cognitive = cognitive.BuildClient(o)
	client.compute = compute.BuildClient(o)
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
	client.graph = graph.BuildClient(o)
	client.hdinsight = hdinsight.BuildClient(o)
	client.iothub = iothub.BuildClient(o)
	client.keyvault = keyvault.BuildClient(o)
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
	client.search = search.BuildClient(o)
	client.securityCenter = securitycenter.BuildClient(o)
	client.servicebus = servicebus.BuildClient(o)
	client.serviceFabric = servicefabric.BuildClient(o)
	client.scheduler = scheduler.BuildClient(o)
	client.signalr = signalr.BuildClient(o)
	client.streamanalytics = streamanalytics.BuildClient(o)
	client.subscription = subscription.BuildClient(o)
	client.sql = sql.BuildClient(o)
	client.trafficManager = trafficmanager.BuildClient(o)
	client.web = web.BuildClient(o)

	client.registerStorageClients(endpoint, c.SubscriptionID, auth, o)

	return &client, nil
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

func (c *ArmClient) registerStorageClients(endpoint, subscriptionId string, auth autorest.Authorizer, options *common.ClientOptions) {
	accountsClient := storage.NewAccountsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&accountsClient.Client, auth)
	c.storageServiceClient = accountsClient

	usageClient := storage.NewUsagesClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&usageClient.Client, auth)
	c.storageUsageClient = usageClient

	c.storage = intStor.BuildClient(accountsClient, options)
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
