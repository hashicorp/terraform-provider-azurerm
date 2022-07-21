package clients

import (
	"context"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	aadb2c "github.com/hashicorp/terraform-provider-azurerm/internal/services/aadb2c/client"
	advisor "github.com/hashicorp/terraform-provider-azurerm/internal/services/advisor/client"
	analysisServices "github.com/hashicorp/terraform-provider-azurerm/internal/services/analysisservices/client"
	apiManagement "github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/client"
	appConfiguration "github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/client"
	applicationInsights "github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/client"
	appService "github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/client"
	attestation "github.com/hashicorp/terraform-provider-azurerm/internal/services/attestation/client"
	authorization "github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/client"
	automation "github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/client"
	azureStackHCI "github.com/hashicorp/terraform-provider-azurerm/internal/services/azurestackhci/client"
	batch "github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/client"
	blueprints "github.com/hashicorp/terraform-provider-azurerm/internal/services/blueprints/client"
	bot "github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/client"
	cdn "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/client"
	cognitiveServices "github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/client"
	communication "github.com/hashicorp/terraform-provider-azurerm/internal/services/communication/client"
	compute "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/client"
	confidentialledger "github.com/hashicorp/terraform-provider-azurerm/internal/services/confidentialledger/client"
	connections "github.com/hashicorp/terraform-provider-azurerm/internal/services/connections/client"
	consumption "github.com/hashicorp/terraform-provider-azurerm/internal/services/consumption/client"
	containerServices "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/client"
	cosmosdb "github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/client"
	costmanagement "github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/client"
	customproviders "github.com/hashicorp/terraform-provider-azurerm/internal/services/customproviders/client"
	datamigration "github.com/hashicorp/terraform-provider-azurerm/internal/services/databasemigration/client"
	databoxedge "github.com/hashicorp/terraform-provider-azurerm/internal/services/databoxedge/client"
	databricks "github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/client"
	datafactory "github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/client"
	dataprotection "github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/client"
	datashare "github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/client"
	desktopvirtualization "github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/client"
	devtestlabs "github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/client"
	digitaltwins "github.com/hashicorp/terraform-provider-azurerm/internal/services/digitaltwins/client"
	disks "github.com/hashicorp/terraform-provider-azurerm/internal/services/disks/client"
	dns "github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/client"
	domainservices "github.com/hashicorp/terraform-provider-azurerm/internal/services/domainservices/client"
	elastic "github.com/hashicorp/terraform-provider-azurerm/internal/services/elastic/client"
	eventgrid "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventgrid/client"
	eventhub "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/client"
	firewall "github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/client"
	fluidrelay "github.com/hashicorp/terraform-provider-azurerm/internal/services/fluidrelay/client"
	frontdoor "github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/client"
	hdinsight "github.com/hashicorp/terraform-provider-azurerm/internal/services/hdinsight/client"
	healthcare "github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/client"
	hpccache "github.com/hashicorp/terraform-provider-azurerm/internal/services/hpccache/client"
	hsm "github.com/hashicorp/terraform-provider-azurerm/internal/services/hsm/client"
	iotcentral "github.com/hashicorp/terraform-provider-azurerm/internal/services/iotcentral/client"
	iothub "github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/client"
	timeseriesinsights "github.com/hashicorp/terraform-provider-azurerm/internal/services/iottimeseriesinsights/client"
	keyvault "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/client"
	kusto "github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/client"
	legacy "github.com/hashicorp/terraform-provider-azurerm/internal/services/legacy/client"
	lighthouse "github.com/hashicorp/terraform-provider-azurerm/internal/services/lighthouse/client"
	loadbalancers "github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/client"
	loadtest "github.com/hashicorp/terraform-provider-azurerm/internal/services/loadtest/client"
	loganalytics "github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/client"
	logic "github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/client"
	logz "github.com/hashicorp/terraform-provider-azurerm/internal/services/logz/client"
	machinelearning "github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/client"
	maintenance "github.com/hashicorp/terraform-provider-azurerm/internal/services/maintenance/client"
	managedapplication "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedapplications/client"
	managementgroup "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/client"
	maps "github.com/hashicorp/terraform-provider-azurerm/internal/services/maps/client"
	mariadb "github.com/hashicorp/terraform-provider-azurerm/internal/services/mariadb/client"
	media "github.com/hashicorp/terraform-provider-azurerm/internal/services/media/client"
	mixedreality "github.com/hashicorp/terraform-provider-azurerm/internal/services/mixedreality/client"
	monitor "github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/client"
	msi "github.com/hashicorp/terraform-provider-azurerm/internal/services/msi/client"
	mssql "github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/client"
	mysql "github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/client"
	netapp "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/client"
	network "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/client"
	notificationhub "github.com/hashicorp/terraform-provider-azurerm/internal/services/notificationhub/client"
	policy "github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/client"
	portal "github.com/hashicorp/terraform-provider-azurerm/internal/services/portal/client"
	postgres "github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/client"
	powerBI "github.com/hashicorp/terraform-provider-azurerm/internal/services/powerbi/client"
	privatedns "github.com/hashicorp/terraform-provider-azurerm/internal/services/privatedns/client"
	purview "github.com/hashicorp/terraform-provider-azurerm/internal/services/purview/client"
	recoveryServices "github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/client"
	redis "github.com/hashicorp/terraform-provider-azurerm/internal/services/redis/client"
	redisenterprise "github.com/hashicorp/terraform-provider-azurerm/internal/services/redisenterprise/client"
	relay "github.com/hashicorp/terraform-provider-azurerm/internal/services/relay/client"
	resource "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/client"
	search "github.com/hashicorp/terraform-provider-azurerm/internal/services/search/client"
	securityCenter "github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/client"
	sentinel "github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/client"
	serviceBus "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/client"
	serviceFabric "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicefabric/client"
	serviceFabricManaged "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicefabricmanaged/client"
	signalr "github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/client"
	appPlatform "github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/client"
	sql "github.com/hashicorp/terraform-provider-azurerm/internal/services/sql/client"
	storage "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	streamAnalytics "github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/client"
	subscription "github.com/hashicorp/terraform-provider-azurerm/internal/services/subscription/client"
	synapse "github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/client"
	trafficManager "github.com/hashicorp/terraform-provider-azurerm/internal/services/trafficmanager/client"
	videoAnalyzer "github.com/hashicorp/terraform-provider-azurerm/internal/services/videoanalyzer/client"
	vmware "github.com/hashicorp/terraform-provider-azurerm/internal/services/vmware/client"
	web "github.com/hashicorp/terraform-provider-azurerm/internal/services/web/client"
)

type Client struct {
	// StopContext is used for propagating control from Terraform Core (e.g. Ctrl/Cmd+C)
	StopContext context.Context

	Account  *ResourceManagerAccount
	Features features.UserFeatures

	AadB2c                *aadb2c.Client
	Advisor               *advisor.Client
	AnalysisServices      *analysisServices.Client
	ApiManagement         *apiManagement.Client
	AppConfiguration      *appConfiguration.Client
	AppInsights           *applicationInsights.Client
	AppPlatform           *appPlatform.Client
	AppService            *appService.Client
	Attestation           *attestation.Client
	Authorization         *authorization.Client
	Automation            *automation.Client
	AzureStackHCI         *azureStackHCI.Client
	Batch                 *batch.Client
	Blueprints            *blueprints.Client
	Bot                   *bot.Client
	Cdn                   *cdn.Client
	Cognitive             *cognitiveServices.Client
	Communication         *communication.Client
	Compute               *compute.Client
	ConfidentialLedger    *confidentialledger.Client
	Connections           *connections.Client
	Consumption           *consumption.Client
	Containers            *containerServices.Client
	Cosmos                *cosmosdb.Client
	CostManagement        *costmanagement.Client
	CustomProviders       *customproviders.Client
	DatabaseMigration     *datamigration.Client
	DataBricks            *databricks.Client
	DataboxEdge           *databoxedge.Client
	DataFactory           *datafactory.Client
	DataProtection        *dataprotection.Client
	DataShare             *datashare.Client
	DesktopVirtualization *desktopvirtualization.Client
	DevTestLabs           *devtestlabs.Client
	DigitalTwins          *digitaltwins.Client
	Disks                 *disks.Client
	Dns                   *dns.Client
	DomainServices        *domainservices.Client
	Elastic               *elastic.Client
	EventGrid             *eventgrid.Client
	Eventhub              *eventhub.Client
	Firewall              *firewall.Client
	FluidRelay            *fluidrelay.Client
	Frontdoor             *frontdoor.Client
	HPCCache              *hpccache.Client
	HSM                   *hsm.Client
	HDInsight             *hdinsight.Client
	HealthCare            *healthcare.Client
	IoTCentral            *iotcentral.Client
	IoTHub                *iothub.Client
	IoTTimeSeriesInsights *timeseriesinsights.Client
	KeyVault              *keyvault.Client
	Kusto                 *kusto.Client
	Legacy                *legacy.Client
	Lighthouse            *lighthouse.Client
	LoadBalancers         *loadbalancers.Client
	LoadTest              *loadtest.Client
	LogAnalytics          *loganalytics.Client
	Logic                 *logic.Client
	Logz                  *logz.Client
	MachineLearning       *machinelearning.Client
	Maintenance           *maintenance.Client
	ManagedApplication    *managedapplication.Client
	ManagementGroups      *managementgroup.Client
	Maps                  *maps.Client
	MariaDB               *mariadb.Client
	Media                 *media.Client
	MixedReality          *mixedreality.Client
	Monitor               *monitor.Client
	MSI                   *msi.Client
	MSSQL                 *mssql.Client
	MySQL                 *mysql.Client
	NetApp                *netapp.Client
	Network               *network.Client
	NotificationHubs      *notificationhub.Client
	Policy                *policy.Client
	Portal                *portal.Client
	Postgres              *postgres.Client
	PowerBI               *powerBI.Client
	PrivateDns            *privatedns.Client
	Purview               *purview.Client
	RecoveryServices      *recoveryServices.Client
	Redis                 *redis.Client
	RedisEnterprise       *redisenterprise.Client
	Relay                 *relay.Client
	Resource              *resource.Client
	Search                *search.Client
	SecurityCenter        *securityCenter.Client
	Sentinel              *sentinel.Client
	ServiceBus            *serviceBus.Client
	ServiceFabric         *serviceFabric.Client
	ServiceFabricManaged  *serviceFabricManaged.Client
	SignalR               *signalr.Client
	Storage               *storage.Client
	StreamAnalytics       *streamAnalytics.Client
	Subscription          *subscription.Client
	Sql                   *sql.Client
	Synapse               *synapse.Client
	TrafficManager        *trafficManager.Client
	VideoAnalyzer         *videoAnalyzer.Client
	Vmware                *vmware.Client
	Web                   *web.Client
}

// NOTE: it should be possible for this method to become Private once the top level Client's removed

func (client *Client) Build(ctx context.Context, o *common.ClientOptions) error {
	autorest.Count429AsRetry = false
	// Disable the Azure SDK for Go's validation since it's unhelpful for our use-case
	validation.Disabled = true

	client.Features = o.Features
	client.StopContext = ctx

	client.AadB2c = aadb2c.NewClient(o)
	client.Advisor = advisor.NewClient(o)
	client.AnalysisServices = analysisServices.NewClient(o)
	client.ApiManagement = apiManagement.NewClient(o)
	client.AppConfiguration = appConfiguration.NewClient(o)
	client.AppInsights = applicationInsights.NewClient(o)
	client.AppPlatform = appPlatform.NewClient(o)
	client.AppService = appService.NewClient(o)
	client.Attestation = attestation.NewClient(o)
	client.Authorization = authorization.NewClient(o)
	client.Automation = automation.NewClient(o)
	client.AzureStackHCI = azureStackHCI.NewClient(o)
	client.Batch = batch.NewClient(o)
	client.Blueprints = blueprints.NewClient(o)
	client.Bot = bot.NewClient(o)
	client.Cdn = cdn.NewClient(o)
	client.Cognitive = cognitiveServices.NewClient(o)
	client.Communication = communication.NewClient(o)
	client.Compute = compute.NewClient(o)
	client.ConfidentialLedger = confidentialledger.NewClient(o)
	client.Connections = connections.NewClient(o)
	client.Consumption = consumption.NewClient(o)
	client.Containers = containerServices.NewClient(o)
	client.Cosmos = cosmosdb.NewClient(o)
	client.CostManagement = costmanagement.NewClient(o)
	client.CustomProviders = customproviders.NewClient(o)
	client.DatabaseMigration = datamigration.NewClient(o)
	client.DataBricks = databricks.NewClient(o)
	client.DataboxEdge = databoxedge.NewClient(o)
	client.DataFactory = datafactory.NewClient(o)
	client.DataProtection = dataprotection.NewClient(o)
	client.DataShare = datashare.NewClient(o)
	client.DesktopVirtualization = desktopvirtualization.NewClient(o)
	client.DevTestLabs = devtestlabs.NewClient(o)
	client.DigitalTwins = digitaltwins.NewClient(o)
	client.Disks = disks.NewClient(o)
	client.Dns = dns.NewClient(o)
	client.DomainServices = domainservices.NewClient(o)
	client.Elastic = elastic.NewClient(o)
	client.EventGrid = eventgrid.NewClient(o)
	client.Eventhub = eventhub.NewClient(o)
	client.Firewall = firewall.NewClient(o)
	client.FluidRelay = fluidrelay.NewClient(o)
	client.Frontdoor = frontdoor.NewClient(o)
	client.HPCCache = hpccache.NewClient(o)
	client.HSM = hsm.NewClient(o)
	client.HDInsight = hdinsight.NewClient(o)
	client.HealthCare = healthcare.NewClient(o)
	client.IoTCentral = iotcentral.NewClient(o)
	client.IoTHub = iothub.NewClient(o)
	client.IoTTimeSeriesInsights = timeseriesinsights.NewClient(o)
	client.KeyVault = keyvault.NewClient(o)
	client.Kusto = kusto.NewClient(o)
	client.Legacy = legacy.NewClient(o)
	client.Lighthouse = lighthouse.NewClient(o)
	client.LogAnalytics = loganalytics.NewClient(o)
	client.LoadBalancers = loadbalancers.NewClient(o)
	client.LoadTest = loadtest.NewClient(o)
	client.Logic = logic.NewClient(o)
	client.Logz = logz.NewClient(o)
	client.MachineLearning = machinelearning.NewClient(o)
	client.Maintenance = maintenance.NewClient(o)
	client.ManagedApplication = managedapplication.NewClient(o)
	client.ManagementGroups = managementgroup.NewClient(o)
	client.Maps = maps.NewClient(o)
	client.MariaDB = mariadb.NewClient(o)
	client.Media = media.NewClient(o)
	client.MixedReality = mixedreality.NewClient(o)
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
	client.PowerBI = powerBI.NewClient(o)
	client.PrivateDns = privatedns.NewClient(o)
	client.Purview = purview.NewClient(o)
	client.RecoveryServices = recoveryServices.NewClient(o)
	client.Redis = redis.NewClient(o)
	client.RedisEnterprise = redisenterprise.NewClient(o)
	client.Relay = relay.NewClient(o)
	client.Resource = resource.NewClient(o)
	client.Search = search.NewClient(o)
	client.SecurityCenter = securityCenter.NewClient(o)
	client.Sentinel = sentinel.NewClient(o)
	client.ServiceBus = serviceBus.NewClient(o)
	client.ServiceFabric = serviceFabric.NewClient(o)
	client.ServiceFabricManaged = serviceFabricManaged.NewClient(o)
	client.SignalR = signalr.NewClient(o)
	client.Sql = sql.NewClient(o)
	client.Storage = storage.NewClient(o)
	client.StreamAnalytics = streamAnalytics.NewClient(o)
	client.Subscription = subscription.NewClient(o)
	client.Synapse = synapse.NewClient(o)
	client.TrafficManager = trafficManager.NewClient(o)
	client.VideoAnalyzer = videoAnalyzer.NewClient(o)
	client.Vmware = vmware.NewClient(o)
	client.Web = web.NewClient(o)

	return nil
}
