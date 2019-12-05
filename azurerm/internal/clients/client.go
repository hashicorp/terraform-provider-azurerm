package clients

import (
	"context"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
	analysisServices "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/analysisservices/client"
	apiManagement "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/client"
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

type Client struct {
	// StopContext is used for propagating control from Terraform Core (e.g. Ctrl/Cmd+C)
	StopContext context.Context

	Account *ResourceManagerAccount

	AnalysisServices *analysisServices.Client
	ApiManagement    *apiManagement.Client
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

	// TODO: Phase 2
	DataBricks  *databricks.Client
	DataFactory *datafactory.Client
	Datalake    *datalake.Client
	DevSpace    *devspace.Client
	DevTestLabs *devtestlabs.Client
	Dns         *dns.Client
	EventGrid   *eventgrid.Client
	Eventhub    *eventhub.Client
	Frontdoor   *frontdoor.Client
	Graph       *graph.Client
	HDInsight   *hdinsight.Client
	Healthcare  *healthcare.Client

	// TODO: Phrase 3
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
	Msi              *msi.Client
	Mssql            *mssql.Client
	Mysql            *mysql.Client

	// TODO: Phase 4
	Netapp           *netapp.Client
	Network          *network.Client
	NotificationHubs *notificationhub.Client
	Policy           *policy.Client
	Portal           *portal.Client
	Postgres         *postgres.Client
	PrivateDns       *privatedns.Client
	RecoveryServices *recoveryservices.Client
	Redis            *redis.Client
	Relay            *relay.Client
	Resource         *resource.Client

	// TODO: Phase 5
	Scheduler       *scheduler.Client
	Search          *search.Client
	SecurityCenter  *securitycenter.Client
	ServiceBus      *servicebus.Client
	ServiceFabric   *servicefabric.Client
	SignalR         *signalr.Client
	Storage         *storage.Client
	StreamAnalytics *streamanalytics.Client
	Subscription    *subscription.Client
	Sql             *sql.Client
	TrafficManager  *trafficmanager.Client
	Web             *web.Client
}

func (client *Client) Build(o *common.ClientOptions) error {
	client.AnalysisServices = analysisServices.NewClient(o)
	client.ApiManagement = apiManagement.NewClient(o)
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
	client.Healthcare = healthcare.NewClient(o)
	client.IoTHub = iothub.NewClient(o)

	return nil
}
