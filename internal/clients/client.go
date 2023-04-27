package clients

import (
	"context"
	"fmt"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/validation"
	aadb2c_v2021_04_01_preview "github.com/hashicorp/go-azure-sdk/resource-manager/aadb2c/2021-04-01-preview"
	analysisservices_v2017_08_01 "github.com/hashicorp/go-azure-sdk/resource-manager/analysisservices/2017-08-01"
	azurestackhci_v2022_12_01 "github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2022-12-01"
	datadog_v2021_03_01 "github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01"
	dns_v2018_05_01 "github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01"
	fluidrelay_2022_05_26 "github.com/hashicorp/go-azure-sdk/resource-manager/fluidrelay/2022-05-26"
	nginx2 "github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2022-08-01"
	redis_v2022_06_01 "github.com/hashicorp/go-azure-sdk/resource-manager/redis/2022-06-01"
	timeseriesinsights_v2020_05_15 "github.com/hashicorp/go-azure-sdk/resource-manager/timeseriesinsights/2020-05-15"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	aadb2c "github.com/hashicorp/terraform-provider-azurerm/internal/services/aadb2c/client"
	advisor "github.com/hashicorp/terraform-provider-azurerm/internal/services/advisor/client"
	analysisServices "github.com/hashicorp/terraform-provider-azurerm/internal/services/analysisservices/client"
	apiManagement "github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/client"
	appConfiguration "github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/client"
	applicationInsights "github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/client"
	appService "github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/client"
	arckubernetes "github.com/hashicorp/terraform-provider-azurerm/internal/services/arckubernetes/client"
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
	containerapps "github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/client"
	containerServices "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/client"
	cosmosdb "github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/client"
	costmanagement "github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/client"
	customproviders "github.com/hashicorp/terraform-provider-azurerm/internal/services/customproviders/client"
	dashboard "github.com/hashicorp/terraform-provider-azurerm/internal/services/dashboard/client"
	datamigration "github.com/hashicorp/terraform-provider-azurerm/internal/services/databasemigration/client"
	databoxedge "github.com/hashicorp/terraform-provider-azurerm/internal/services/databoxedge/client"
	databricks "github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/client"
	datadog "github.com/hashicorp/terraform-provider-azurerm/internal/services/datadog/client"
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
	hybridcompute "github.com/hashicorp/terraform-provider-azurerm/internal/services/hybridcompute/client"
	iotcentral "github.com/hashicorp/terraform-provider-azurerm/internal/services/iotcentral/client"
	iothub "github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/client"
	timeseriesinsights "github.com/hashicorp/terraform-provider-azurerm/internal/services/iottimeseriesinsights/client"
	keyvault "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/client"
	kusto "github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/client"
	labservice "github.com/hashicorp/terraform-provider-azurerm/internal/services/labservice/client"
	legacy "github.com/hashicorp/terraform-provider-azurerm/internal/services/legacy/client"
	lighthouse "github.com/hashicorp/terraform-provider-azurerm/internal/services/lighthouse/client"
	loadbalancers "github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/client"
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
	mobilenetwork "github.com/hashicorp/terraform-provider-azurerm/internal/services/mobilenetwork/client"
	monitor "github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/client"
	mssql "github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/client"
	mssqlmanagedinstance "github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/client"
	mysql "github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/client"
	netapp "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/client"
	network "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/client"
	nginx "github.com/hashicorp/terraform-provider-azurerm/internal/services/nginx/client"
	notificationhub "github.com/hashicorp/terraform-provider-azurerm/internal/services/notificationhub/client"
	orbital "github.com/hashicorp/terraform-provider-azurerm/internal/services/orbital/client"
	policy "github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/client"
	portal "github.com/hashicorp/terraform-provider-azurerm/internal/services/portal/client"
	postgres "github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/client"
	powerBI "github.com/hashicorp/terraform-provider-azurerm/internal/services/powerbi/client"
	privatedns "github.com/hashicorp/terraform-provider-azurerm/internal/services/privatedns/client"
	dnsresolver "github.com/hashicorp/terraform-provider-azurerm/internal/services/privatednsresolver/client"
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
	serviceConnector "github.com/hashicorp/terraform-provider-azurerm/internal/services/serviceconnector/client"
	serviceFabric "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicefabric/client"
	serviceFabricManaged "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicefabricmanaged/client"
	signalr "github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/client"
	appPlatform "github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/client"
	sql "github.com/hashicorp/terraform-provider-azurerm/internal/services/sql/client"
	storage "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	storageMover "github.com/hashicorp/terraform-provider-azurerm/internal/services/storagemover/client"
	streamAnalytics "github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/client"
	subscription "github.com/hashicorp/terraform-provider-azurerm/internal/services/subscription/client"
	synapse "github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/client"
	trafficManager "github.com/hashicorp/terraform-provider-azurerm/internal/services/trafficmanager/client"
	videoAnalyzer "github.com/hashicorp/terraform-provider-azurerm/internal/services/videoanalyzer/client"
	vmware "github.com/hashicorp/terraform-provider-azurerm/internal/services/vmware/client"
	voiceServices "github.com/hashicorp/terraform-provider-azurerm/internal/services/voiceservices/client"
	web "github.com/hashicorp/terraform-provider-azurerm/internal/services/web/client"
)

type Client struct {
	autoClient

	// StopContext is used for propagating control from Terraform Core (e.g. Ctrl/Cmd+C)
	StopContext context.Context

	Account  *ResourceManagerAccount
	Features features.UserFeatures

	AadB2c                *aadb2c_v2021_04_01_preview.Client
	Advisor               *advisor.Client
	AnalysisServices      *analysisservices_v2017_08_01.Client
	ApiManagement         *apiManagement.Client
	AppConfiguration      *appConfiguration.Client
	AppInsights           *applicationInsights.Client
	AppPlatform           *appPlatform.Client
	AppService            *appService.Client
	ArcKubernetes         *arckubernetes.Client
	Attestation           *attestation.Client
	Authorization         *authorization.Client
	Automation            *automation.Client
	AzureStackHCI         *azurestackhci_v2022_12_01.Client
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
	ContainerApps         *containerapps.Client
	Containers            *containerServices.Client
	Cosmos                *cosmosdb.Client
	CostManagement        *costmanagement.Client
	CustomProviders       *customproviders.Client
	Dashboard             *dashboard.Client
	DatabaseMigration     *datamigration.Client
	DataBricks            *databricks.Client
	DataboxEdge           *databoxedge.Client
	Datadog               *datadog_v2021_03_01.Client
	DataFactory           *datafactory.Client
	DataProtection        *dataprotection.Client
	DataShare             *datashare.Client
	DesktopVirtualization *desktopvirtualization.Client
	DevTestLabs           *devtestlabs.Client
	DigitalTwins          *digitaltwins.Client
	Disks                 *disks.Client
	Dns                   *dns_v2018_05_01.Client
	DomainServices        *domainservices.Client
	Elastic               *elastic.Client
	EventGrid             *eventgrid.Client
	Eventhub              *eventhub.Client
	Firewall              *firewall.Client
	FluidRelay            *fluidrelay_2022_05_26.Client
	Frontdoor             *frontdoor.Client
	HPCCache              *hpccache.Client
	HSM                   *hsm.Client
	HDInsight             *hdinsight.Client
	HybridCompute         *hybridcompute.Client
	HealthCare            *healthcare.Client
	IoTCentral            *iotcentral.Client
	IoTHub                *iothub.Client
	IoTTimeSeriesInsights *timeseriesinsights_v2020_05_15.Client
	KeyVault              *keyvault.Client
	Kusto                 *kusto.Client
	LabService            *labservice.Client
	Legacy                *legacy.Client
	Lighthouse            *lighthouse.Client
	LoadBalancers         *loadbalancers.Client
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
	MobileNetwork         *mobilenetwork.Client
	MSSQL                 *mssql.Client
	MSSQLManagedInstance  *mssqlmanagedinstance.Client
	MySQL                 *mysql.Client
	NetApp                *netapp.Client
	Network               *network.Client
	Nginx                 *nginx2.Client
	NotificationHubs      *notificationhub.Client
	Orbital               *orbital.Client
	Policy                *policy.Client
	Portal                *portal.Client
	Postgres              *postgres.Client
	PowerBI               *powerBI.Client
	PrivateDns            *privatedns.Client
	PrivateDnsResolver    *dnsresolver.Client
	Purview               *purview.Client
	RecoveryServices      *recoveryServices.Client
	Redis                 *redis_v2022_06_01.Client
	RedisEnterprise       *redisenterprise.Client
	Relay                 *relay.Client
	Resource              *resource.Client
	Search                *search.Client
	SecurityCenter        *securityCenter.Client
	Sentinel              *sentinel.Client
	ServiceBus            *serviceBus.Client
	ServiceConnector      *serviceConnector.Client
	ServiceFabric         *serviceFabric.Client
	ServiceFabricManaged  *serviceFabricManaged.Client
	SignalR               *signalr.Client
	Storage               *storage.Client
	StorageMover          *storageMover.Client
	StreamAnalytics       *streamAnalytics.Client
	Subscription          *subscription.Client
	Sql                   *sql.Client
	Synapse               *synapse.Client
	TrafficManager        *trafficManager.Client
	VideoAnalyzer         *videoAnalyzer.Client
	Vmware                *vmware.Client
	VoiceServices         *voiceServices.Client
	Web                   *web.Client
}

// NOTE: it should be possible for this method to become Private once the top level Client's removed

func (client *Client) Build(ctx context.Context, o *common.ClientOptions) error {
	autorest.Count429AsRetry = false
	// Disable the Azure SDK for Go's validation since it's unhelpful for our use-case
	validation.Disabled = true

	if err := buildAutoClients(&client.autoClient, o); err != nil {
		return fmt.Errorf("building auto-clients: %+v", err)
	}

	client.Features = o.Features
	client.StopContext = ctx

	var err error

	if client.AadB2c, err = aadb2c.NewClient(o); err != nil {
		return fmt.Errorf("building clients for AadB2c: %+v", err)
	}
	client.Advisor = advisor.NewClient(o)
	client.AnalysisServices = analysisServices.NewClient(o)
	client.ApiManagement = apiManagement.NewClient(o)
	client.AppConfiguration = appConfiguration.NewClient(o)
	client.AppInsights = applicationInsights.NewClient(o)
	client.AppPlatform = appPlatform.NewClient(o)
	client.AppService = appService.NewClient(o)
	client.ArcKubernetes = arckubernetes.NewClient(o)
	if client.Attestation, err = attestation.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Attestation: %+v", err)
	}
	client.Authorization = authorization.NewClient(o)
	client.Automation = automation.NewClient(o)
	client.AzureStackHCI = azureStackHCI.NewClient(o)
	client.Batch = batch.NewClient(o)
	client.Blueprints = blueprints.NewClient(o)
	client.Bot = bot.NewClient(o)
	client.Cdn = cdn.NewClient(o)
	client.Cognitive = cognitiveServices.NewClient(o)
	if client.Communication, err = communication.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Communication: %+v", err)
	}
	client.Compute = compute.NewClient(o)
	if client.ConfidentialLedger, err = confidentialledger.NewClient(o); err != nil {
		return fmt.Errorf("building clients for ConfidentialLedger: %+v", err)
	}
	client.Connections = connections.NewClient(o)
	if client.Consumption, err = consumption.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Consumption: %+v", err)
	}
	if client.Containers, err = containerServices.NewContainersClient(o); err != nil {
		return fmt.Errorf("building clients for Containers: %+v", err)
	}
	client.ContainerApps = containerapps.NewClient(o)
	client.Cosmos = cosmosdb.NewClient(o)
	if client.CostManagement, err = costmanagement.NewClient(o); err != nil {
		return fmt.Errorf("building clients for CostManagement: %+v", err)
	}
	if client.CustomProviders, err = customproviders.NewClient(o); err != nil {
		return fmt.Errorf("building clients for CustomProviders: %+v", err)
	}
	if client.Dashboard, err = dashboard.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Dashboard: %+v", err)
	}
	client.DatabaseMigration = datamigration.NewClient(o)
	if client.DataBricks, err = databricks.NewClient(o); err != nil {
		return fmt.Errorf("building clients for DataBricks: %+v", err)
	}
	if client.DataboxEdge, err = databoxedge.NewClient(o); err != nil {
		return fmt.Errorf("building clients for DataboxEdge: %+v", err)
	}
	if client.Datadog, err = datadog.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Datadog: %+v", err)
	}
	client.DataFactory = datafactory.NewClient(o)
	if client.DataProtection, err = dataprotection.NewClient(o); err != nil {
		return fmt.Errorf("building clients for DataProtection: %+v", err)
	}
	if client.DataShare, err = datashare.NewClient(o); err != nil {
		return fmt.Errorf("building clients for DataShare: %+v", err)
	}
	if client.DesktopVirtualization, err = desktopvirtualization.NewClient(o); err != nil {
		return fmt.Errorf("building clients for DesktopVirtualization: %+v", err)
	}
	client.DevTestLabs = devtestlabs.NewClient(o)
	if client.DigitalTwins, err = digitaltwins.NewClient(o); err != nil {
		return fmt.Errorf("building clients for DigitalTwins: %+v", err)
	}
	client.Disks = disks.NewClient(o)
	if client.Dns, err = dns.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Dns: %+v", err)
	}
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
	client.HybridCompute = hybridcompute.NewClient(o)
	client.IoTCentral = iotcentral.NewClient(o)
	client.IoTHub = iothub.NewClient(o)
	client.IoTTimeSeriesInsights = timeseriesinsights.NewClient(o)
	client.KeyVault = keyvault.NewClient(o)
	client.Kusto = kusto.NewClient(o)
	client.LabService = labservice.NewClient(o)
	client.Legacy = legacy.NewClient(o)
	client.Lighthouse = lighthouse.NewClient(o)
	client.LogAnalytics = loganalytics.NewClient(o)
	client.LoadBalancers = loadbalancers.NewClient(o)
	client.Logic = logic.NewClient(o)
	client.Logz = logz.NewClient(o)
	client.MachineLearning = machinelearning.NewClient(o)
	client.Maintenance = maintenance.NewClient(o)
	client.ManagedApplication = managedapplication.NewClient(o)
	client.ManagementGroups = managementgroup.NewClient(o)
	if client.Maps, err = maps.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Maps: %+v", err)
	}
	client.MariaDB = mariadb.NewClient(o)
	if client.Media, err = media.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Media: %+v", err)
	}
	client.MixedReality = mixedreality.NewClient(o)
	client.Monitor = monitor.NewClient(o)
	client.MobileNetwork = mobilenetwork.NewClient(o)
	client.MSSQL = mssql.NewClient(o)
	client.MySQL = mysql.NewClient(o)
	client.NetApp = netapp.NewClient(o)
	client.Network = network.NewClient(o)
	client.Nginx = nginx.NewClient(o)
	client.NotificationHubs = notificationhub.NewClient(o)
	client.Orbital = orbital.NewClient(o)
	client.Policy = policy.NewClient(o)
	client.Portal = portal.NewClient(o)
	client.Postgres = postgres.NewClient(o)
	client.PowerBI = powerBI.NewClient(o)
	client.PrivateDns = privatedns.NewClient(o)
	client.PrivateDnsResolver = dnsresolver.NewClient(o)
	client.Purview = purview.NewClient(o)
	client.RecoveryServices = recoveryServices.NewClient(o)
	client.Redis = redis.NewClient(o)
	client.RedisEnterprise = redisenterprise.NewClient(o)
	client.Relay = relay.NewClient(o)
	client.Resource = resource.NewClient(o)
	client.Search = search.NewClient(o)
	client.SecurityCenter = securityCenter.NewClient(o)
	client.Sentinel = sentinel.NewClient(o)
	if client.ServiceBus, err = serviceBus.NewClient(o); err != nil {
		return fmt.Errorf("building clients for ServiceBus: %+v", err)
	}
	client.ServiceConnector = serviceConnector.NewClient(o)
	client.ServiceFabric = serviceFabric.NewClient(o)
	client.ServiceFabricManaged = serviceFabricManaged.NewClient(o)
	if client.SignalR, err = signalr.NewClient(o); err != nil {
		return fmt.Errorf("building clients for SignalR: %+v", err)
	}
	client.Sql = sql.NewClient(o)
	client.Storage = storage.NewClient(o)
	if client.StorageMover, err = storageMover.NewClient(o); err != nil {
		return fmt.Errorf("building clients for StorageMover: %+v", err)
	}
	client.StreamAnalytics = streamAnalytics.NewClient(o)
	client.Subscription = subscription.NewClient(o)
	client.Synapse = synapse.NewClient(o)
	client.TrafficManager = trafficManager.NewClient(o)
	client.VideoAnalyzer = videoAnalyzer.NewClient(o)
	client.Vmware = vmware.NewClient(o)
	client.VoiceServices = voiceServices.NewClient(o)
	client.Web = web.NewClient(o)

	return nil
}
