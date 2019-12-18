package clients

import (
	"context"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
	analysisServices "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/analysisservices/client"
	apiManagement "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/client"
	appConfiguration "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/appconfiguration/client"
	applicationInsights "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/applicationinsights/client"
	authorization "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/authorization/client"
	automation "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/automation/client"
	batch "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/batch/client"
	bot "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/bot/client"
	cdn "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn/client"
	cognitiveServices "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cognitive/client"
	compute "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/client"
	containerServices "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/containers/client"
	cosmosdb "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/client"
	databricks "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databricks/client"
	datafactory "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datafactory/client"
	datalake "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datalake/client"
	devspace "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/devspace/client"
	devtestlabs "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/devtestlabs/client"
	dns "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dns/client"
	eventgrid "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventgrid/client"
	eventhub "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/client"
	frontdoor "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/frontdoor/client"
	graph "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/graph/client"
	hdinsight "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hdinsight/client"
	healthcare "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/healthcare/client"
	iothub "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iothub/client"
	keyvault "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/client"
	kusto "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto/client"
	loganalytics "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/client"
	logic "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/logic/client"
	managementgroup "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managementgroup/client"
	maps "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/maps/client"
	mariadb "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mariadb/client"
	media "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/media/client"
	monitor "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/monitor/client"
	msi "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/client"
	mssql "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/client"
	mysql "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mysql/client"
	netapp "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/netapp/client"
	network "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/client"
	notificationhub "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/notificationhub/client"
	policy "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/client"
	portal "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/portal/client"
	postgres "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/client"
	privatedns "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/privatedns/client"
	recoveryServices "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/recoveryservices/client"
	redis "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/redis/client"
	relay "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/relay/client"
	resource "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/client"
	scheduler "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/scheduler/client"
	search "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/search/client"
	securityCenter "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter/client"
	serviceBus "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicebus/client"
	serviceFabric "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicefabric/client"
	signalr "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/signalr/client"
	sql "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sql/client"
	storage "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/client"
	streamAnalytics "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/streamanalytics/client"
	subscription "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/subscription/client"
	trafficManager "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/trafficmanager/client"
	web "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/client"
)

type Client struct {
	// StopContext is used for propagating control from Terraform Core (e.g. Ctrl/Cmd+C)
	StopContext context.Context

	Account *ResourceManagerAccount

	AnalysisServices *analysisServices.Client
	ApiManagement    *apiManagement.Client
	AppConfiguration *appConfiguration.Client
	AppInsights      *applicationInsights.Client
	Authorization    *authorization.Client
	Automation       *automation.Client
	Batch            *batch.Client
	Bot              *bot.Client
	Cdn              *cdn.Client
	Cognitive        *cognitiveServices.Client
	Compute          *compute.Client
	Containers       *containerServices.Client
	Cosmos           *cosmosdb.Client
	DataBricks       *databricks.Client
	DataFactory      *datafactory.Client
	Datalake         *datalake.Client
	DevSpace         *devspace.Client
	DevTestLabs      *devtestlabs.Client
	Dns              *dns.Client
	EventGrid        *eventgrid.Client
	Eventhub         *eventhub.Client
	Frontdoor        *frontdoor.Client
	Graph            *graph.Client
	HDInsight        *hdinsight.Client
	HealthCare       *healthcare.Client
	IoTHub           *iothub.Client
	KeyVault         *keyvault.Client
	Kusto            *kusto.Client
	LogAnalytics     *loganalytics.Client
	Logic            *logic.Client
	ManagementGroups *managementgroup.Client
	Maps             *maps.Client
	MariaDB          *mariadb.Client
	Media            *media.Client
	Monitor          *monitor.Client
	MSI              *msi.Client
	MSSQL            *mssql.Client
	MySQL            *mysql.Client
	NetApp           *netapp.Client
	Network          *network.Client
	NotificationHubs *notificationhub.Client
	Policy           *policy.Client
	Portal           *portal.Client
	Postgres         *postgres.Client
	PrivateDns       *privatedns.Client
	RecoveryServices *recoveryServices.Client
	Redis            *redis.Client
	Relay            *relay.Client
	Resource         *resource.Client
	Scheduler        *scheduler.Client
	Search           *search.Client
	SecurityCenter   *securityCenter.Client
	ServiceBus       *serviceBus.Client
	ServiceFabric    *serviceFabric.Client
	SignalR          *signalr.Client
	Storage          *storage.Client
	StreamAnalytics  *streamAnalytics.Client
	Subscription     *subscription.Client
	Sql              *sql.Client
	TrafficManager   *trafficManager.Client
	Web              *web.Client
}

func (client *Client) Build(o *common.ClientOptions) error {
	client.AnalysisServices = analysisServices.NewClient(o)
	client.ApiManagement = apiManagement.NewClient(o)
	client.AppConfiguration = appConfiguration.NewClient(o)
	client.AppInsights = applicationInsights.NewClient(o)
	client.Authorization = authorization.NewClient(o)
	client.Automation = automation.NewClient(o)
	client.Batch = batch.NewClient(o)
	client.Bot = bot.NewClient(o)
	client.Cdn = cdn.NewClient(o)
	client.Cognitive = cognitiveServices.NewClient(o)
	client.Compute = compute.NewClient(o)
	client.Containers = containerServices.NewClient(o)
	client.Cosmos = cosmosdb.NewClient(o)
	client.DataBricks = databricks.NewClient(o)
	client.DataFactory = datafactory.NewClient(o)
	client.Datalake = datalake.NewClient(o)
	client.DevSpace = devspace.NewClient(o)
	client.DevTestLabs = devtestlabs.NewClient(o)
	client.Dns = dns.NewClient(o)
	client.EventGrid = eventgrid.NewClient(o)
	client.Eventhub = eventhub.NewClient(o)
	client.Frontdoor = frontdoor.NewClient(o)
	client.Graph = graph.NewClient(o)
	client.HDInsight = hdinsight.NewClient(o)
	client.HealthCare = healthcare.NewClient(o)
	client.IoTHub = iothub.NewClient(o)
	client.KeyVault = keyvault.NewClient(o)
	client.Kusto = kusto.NewClient(o)
	client.LogAnalytics = loganalytics.NewClient(o)
	client.Logic = logic.NewClient(o)
	client.ManagementGroups = managementgroup.NewClient(o)
	client.Maps = maps.NewClient(o)
	client.MariaDB = mariadb.NewClient(o)
	client.Media = media.NewClient(o)
	client.Monitor = monitor.NewClient(o)
	client.MSI = msi.NewClient(o)
	client.MSSQL = mssql.NewClient(o)
	client.MySQL = mysql.NewClient(o)
	client.NetApp = netapp.NewClient(o)
	client.Network = network.NewClient(o)
	client.NotificationHubs = notificationhub.NewClient(o)
	client.Policy = policy.NewClient(o)
	client.Portal = portal.NewClient(o)
	client.Postgres = postgres.NewClient(o)
	client.PrivateDns = privatedns.NewClient(o)
	client.RecoveryServices = recoveryServices.NewClient(o)
	client.Redis = redis.NewClient(o)
	client.Relay = relay.NewClient(o)
	client.Resource = resource.NewClient(o)
	client.Scheduler = scheduler.NewClient(o)
	client.Search = search.NewClient(o)
	client.SecurityCenter = securityCenter.NewClient(o)
	client.ServiceBus = serviceBus.NewClient(o)
	client.ServiceFabric = serviceFabric.NewClient(o)
	client.SignalR = signalr.NewClient(o)
	client.Sql = sql.NewClient(o)
	client.Storage = storage.NewClient(o)
	client.StreamAnalytics = streamAnalytics.NewClient(o)
	client.Subscription = subscription.NewClient(o)
	client.TrafficManager = trafficManager.NewClient(o)
	client.Web = web.NewClient(o)

	return nil
}
