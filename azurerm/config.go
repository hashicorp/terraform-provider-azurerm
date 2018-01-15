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

	"github.com/Azure/azure-sdk-for-go/arm/appinsights"
	"github.com/Azure/azure-sdk-for-go/arm/authorization"
	"github.com/Azure/azure-sdk-for-go/arm/automation"
	"github.com/Azure/azure-sdk-for-go/arm/cdn"
	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/Azure/azure-sdk-for-go/arm/containerinstance"
	"github.com/Azure/azure-sdk-for-go/arm/containerregistry"
	"github.com/Azure/azure-sdk-for-go/arm/containerservice"
	"github.com/Azure/azure-sdk-for-go/arm/cosmos-db"
	"github.com/Azure/azure-sdk-for-go/arm/disk"
	"github.com/Azure/azure-sdk-for-go/arm/dns"
	"github.com/Azure/azure-sdk-for-go/arm/eventgrid"
	"github.com/Azure/azure-sdk-for-go/arm/eventhub"
	"github.com/Azure/azure-sdk-for-go/arm/graphrbac"
	"github.com/Azure/azure-sdk-for-go/arm/keyvault"
	"github.com/Azure/azure-sdk-for-go/arm/mysql"
	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/Azure/azure-sdk-for-go/arm/operationalinsights"
	"github.com/Azure/azure-sdk-for-go/arm/postgresql"
	"github.com/Azure/azure-sdk-for-go/arm/redis"
	"github.com/Azure/azure-sdk-for-go/arm/resources/locks"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/azure-sdk-for-go/arm/resources/subscriptions"
	keyVault "github.com/Azure/azure-sdk-for-go/dataplane/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/scheduler/mgmt/2016-03-01/scheduler"
	"github.com/Azure/azure-sdk-for-go/services/search/mgmt/2015-08-19/search"
	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/Azure/azure-sdk-for-go/services/sql/mgmt/2015-05-01-preview/sql"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2017-06-01/storage"
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

	availSetClient         compute.AvailabilitySetsClient
	usageOpsClient         compute.UsageClient
	vmExtensionImageClient compute.VirtualMachineExtensionImagesClient
	vmExtensionClient      compute.VirtualMachineExtensionsClient
	vmScaleSetClient       compute.VirtualMachineScaleSetsClient
	vmImageClient          compute.VirtualMachineImagesClient
	vmClient               compute.VirtualMachinesClient
	imageClient            compute.ImagesClient

	diskClient                 disk.DisksClient
	snapshotsClient            disk.SnapshotsClient
	cosmosDBClient             cosmosdb.DatabaseAccountsClient
	automationAccountClient    automation.AccountClient
	automationRunbookClient    automation.RunbookClient
	automationCredentialClient automation.CredentialClient
	automationScheduleClient   automation.ScheduleClient

	applicationGatewayClient     network.ApplicationGatewaysClient
	ifaceClient                  network.InterfacesClient
	expressRouteCircuitClient    network.ExpressRouteCircuitsClient
	loadBalancerClient           network.LoadBalancersClient
	localNetConnClient           network.LocalNetworkGatewaysClient
	publicIPClient               network.PublicIPAddressesClient
	secGroupClient               network.SecurityGroupsClient
	secRuleClient                network.SecurityRulesClient
	subnetClient                 network.SubnetsClient
	netUsageClient               network.UsagesClient
	vnetGatewayConnectionsClient network.VirtualNetworkGatewayConnectionsClient
	vnetGatewayClient            network.VirtualNetworkGatewaysClient
	vnetClient                   network.VirtualNetworksClient
	vnetPeeringsClient           network.VirtualNetworkPeeringsClient
	routeTablesClient            network.RouteTablesClient
	routesClient                 network.RoutesClient
	dnsClient                    dns.RecordSetsClient
	zonesClient                  dns.ZonesClient

	cdnProfilesClient  cdn.ProfilesClient
	cdnEndpointsClient cdn.EndpointsClient

	containerRegistryClient containerregistry.RegistriesClient
	containerServicesClient containerservice.ContainerServicesClient
	containerGroupsClient   containerinstance.ContainerGroupsClient

	eventGridTopicsClient       eventgrid.TopicsClient
	eventHubClient              eventhub.EventHubsClient
	eventHubConsumerGroupClient eventhub.ConsumerGroupsClient
	eventHubNamespacesClient    eventhub.NamespacesClient

	workspacesClient operationalinsights.WorkspacesClient

	providers           resources.ProvidersClient
	resourceGroupClient resources.GroupsClient
	tagsClient          resources.TagsClient
	resourceFindClient  resources.GroupClient

	subscriptionsGroupClient subscriptions.GroupClient

	deploymentsClient resources.DeploymentsClient

	redisClient               redis.GroupClient
	redisFirewallClient       redis.FirewallRuleClient
	redisPatchSchedulesClient redis.PatchSchedulesClient

	keyVaultClient           keyvault.VaultsClient
	keyVaultManagementClient keyVault.ManagementClient

	appInsightsClient appinsights.ComponentsClient

	// Authentication
	roleAssignmentsClient   authorization.RoleAssignmentsClient
	roleDefinitionsClient   authorization.RoleDefinitionsClient
	servicePrincipalsClient graphrbac.ServicePrincipalsClient

	// Databases
	mysqlConfigurationsClient      mysql.ConfigurationsClient
	mysqlDatabasesClient           mysql.DatabasesClient
	mysqlFirewallRulesClient       mysql.FirewallRulesClient
	mysqlServersClient             mysql.ServersClient
	postgresqlConfigurationsClient postgresql.ConfigurationsClient
	postgresqlDatabasesClient      postgresql.DatabasesClient
	postgresqlFirewallRulesClient  postgresql.FirewallRulesClient
	postgresqlServersClient        postgresql.ServersClient
	sqlDatabasesClient             sql.DatabasesClient
	sqlElasticPoolsClient          sql.ElasticPoolsClient
	sqlFirewallRulesClient         sql.FirewallRulesClient
	sqlServersClient               sql.ServersClient

	// Networking
	watcherClient network.WatchersClient

	// Resources
	managementLocksClient locks.ManagementLocksClient

	// Scheduler
	jobsClient            scheduler.JobsClient
	jobsCollectionsClient scheduler.JobCollectionsClient

	// Search
	searchServicesClient search.ServicesClient

	// ServiceBus
	serviceBusQueuesClient        servicebus.QueuesClient
	serviceBusNamespacesClient    servicebus.NamespacesClient
	serviceBusTopicsClient        servicebus.TopicsClient
	serviceBusSubscriptionsClient servicebus.SubscriptionsClient

	// Storage
	storageServiceClient storage.AccountsClient
	storageUsageClient   storage.UsageClient

	// Traffic Manager
	trafficManagerProfilesClient  trafficmanager.ProfilesClient
	trafficManagerEndpointsClient trafficmanager.EndpointsClient

	// Web
	appServicePlansClient web.AppServicePlansClient
	appServicesClient     web.AppsClient
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

	// NOTE: these declarations should be left separate for clarity should the
	// clients be wished to be configured with custom Responders/PollingModes etc...
	asc := compute.NewAvailabilitySetsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&asc.Client)
	asc.Authorizer = auth
	asc.Sender = sender
	asc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.availSetClient = asc

	uoc := compute.NewUsageClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&uoc.Client)
	uoc.Authorizer = auth
	uoc.Sender = sender
	uoc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.usageOpsClient = uoc

	vmeic := compute.NewVirtualMachineExtensionImagesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&vmeic.Client)
	vmeic.Authorizer = auth
	vmeic.Sender = sender
	vmeic.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.vmExtensionImageClient = vmeic

	vmec := compute.NewVirtualMachineExtensionsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&vmec.Client)
	vmec.Authorizer = auth
	vmec.Sender = sender
	vmec.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.vmExtensionClient = vmec

	vmic := compute.NewVirtualMachineImagesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&vmic.Client)
	vmic.Authorizer = auth
	vmic.Sender = sender
	vmic.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.vmImageClient = vmic

	vmssc := compute.NewVirtualMachineScaleSetsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&vmssc.Client)
	vmssc.Authorizer = auth
	vmssc.Sender = sender
	vmssc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.vmScaleSetClient = vmssc

	vmc := compute.NewVirtualMachinesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&vmc.Client)
	vmc.Authorizer = auth
	vmc.Sender = sender
	vmc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.vmClient = vmc

	agc := network.NewApplicationGatewaysClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&agc.Client)
	agc.Authorizer = auth
	agc.Sender = sender
	agc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.applicationGatewayClient = agc

	crc := containerregistry.NewRegistriesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&crc.Client)
	crc.Authorizer = auth
	crc.Sender = sender
	crc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.containerRegistryClient = crc

	csc := containerservice.NewContainerServicesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&csc.Client)
	csc.Authorizer = auth
	csc.Sender = sender
	csc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.containerServicesClient = csc

	cgc := containerinstance.NewContainerGroupsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&cgc.Client)
	cgc.Authorizer = auth
	cgc.Sender = autorest.CreateSender(withRequestLogging())
	cgc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.containerGroupsClient = cgc

	cdb := cosmosdb.NewDatabaseAccountsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&cdb.Client)
	cdb.Authorizer = auth
	cdb.Sender = sender
	cdb.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.cosmosDBClient = cdb

	img := compute.NewImagesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&img.Client)
	img.Authorizer = auth
	img.Sender = sender
	img.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.imageClient = img

	egtc := eventgrid.NewTopicsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&egtc.Client)
	egtc.Authorizer = auth
	egtc.Sender = sender
	egtc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.eventGridTopicsClient = egtc

	ehc := eventhub.NewEventHubsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&ehc.Client)
	ehc.Authorizer = auth
	ehc.Sender = sender
	ehc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.eventHubClient = ehc

	chcgc := eventhub.NewConsumerGroupsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&chcgc.Client)
	chcgc.Authorizer = auth
	chcgc.Sender = sender
	chcgc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.eventHubConsumerGroupClient = chcgc

	ehnc := eventhub.NewNamespacesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&ehnc.Client)
	ehnc.Authorizer = auth
	ehnc.Sender = sender
	ehnc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.eventHubNamespacesClient = ehnc

	ifc := network.NewInterfacesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&ifc.Client)
	ifc.Authorizer = auth
	ifc.Sender = sender
	ifc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.ifaceClient = ifc

	erc := network.NewExpressRouteCircuitsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&erc.Client)
	erc.Authorizer = auth
	erc.Sender = sender
	erc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.expressRouteCircuitClient = erc

	lbc := network.NewLoadBalancersClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&lbc.Client)
	lbc.Authorizer = auth
	lbc.Sender = sender
	lbc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.loadBalancerClient = lbc

	lgc := network.NewLocalNetworkGatewaysClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&lgc.Client)
	lgc.Authorizer = auth
	lgc.Sender = sender
	lgc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.localNetConnClient = lgc

	opwc := operationalinsights.NewWorkspacesClient(c.SubscriptionID)
	setUserAgent(&opwc.Client)
	opwc.Authorizer = auth
	opwc.Sender = autorest.CreateSender(withRequestLogging())
	opwc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.workspacesClient = opwc

	pipc := network.NewPublicIPAddressesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&pipc.Client)
	pipc.Authorizer = auth
	pipc.Sender = sender
	pipc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.publicIPClient = pipc

	sgc := network.NewSecurityGroupsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&sgc.Client)
	sgc.Authorizer = auth
	sgc.Sender = sender
	sgc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.secGroupClient = sgc

	src := network.NewSecurityRulesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&src.Client)
	src.Authorizer = auth
	src.Sender = sender
	src.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.secRuleClient = src

	snc := network.NewSubnetsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&snc.Client)
	snc.Authorizer = auth
	snc.Sender = sender
	snc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.subnetClient = snc

	vgcc := network.NewVirtualNetworkGatewayConnectionsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&vgcc.Client)
	vgcc.Authorizer = auth
	vgcc.Sender = sender
	vgcc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.vnetGatewayConnectionsClient = vgcc

	vgc := network.NewVirtualNetworkGatewaysClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&vgc.Client)
	vgc.Authorizer = auth
	vgc.Sender = sender
	vgc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.vnetGatewayClient = vgc

	vnc := network.NewVirtualNetworksClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&vnc.Client)
	vnc.Authorizer = auth
	vnc.Sender = sender
	vnc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.vnetClient = vnc

	vnpc := network.NewVirtualNetworkPeeringsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&vnpc.Client)
	vnpc.Authorizer = auth
	vnpc.Sender = sender
	vnpc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.vnetPeeringsClient = vnpc

	rtc := network.NewRouteTablesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&rtc.Client)
	rtc.Authorizer = auth
	rtc.Sender = sender
	rtc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.routeTablesClient = rtc

	rc := network.NewRoutesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&rc.Client)
	rc.Authorizer = auth
	rc.Sender = sender
	rc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.routesClient = rc

	dn := dns.NewRecordSetsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&dn.Client)
	dn.Authorizer = auth
	dn.Sender = sender
	dn.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.dnsClient = dn

	zo := dns.NewZonesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&zo.Client)
	zo.Authorizer = auth
	zo.Sender = sender
	zo.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.zonesClient = zo

	rgc := resources.NewGroupsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&rgc.Client)
	rgc.Authorizer = auth
	rgc.Sender = sender
	rgc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.resourceGroupClient = rgc

	pc := resources.NewProvidersClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&pc.Client)
	pc.Authorizer = auth
	pc.Sender = sender
	pc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.providers = pc

	tc := resources.NewTagsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&tc.Client)
	tc.Authorizer = auth
	tc.Sender = sender
	tc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.tagsClient = tc

	rf := resources.NewGroupClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&rf.Client)
	rf.Authorizer = auth
	rf.Sender = sender
	rf.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.resourceFindClient = rf

	subgc := subscriptions.NewGroupClientWithBaseURI(endpoint)
	setUserAgent(&subgc.Client)
	subgc.Authorizer = auth
	subgc.Sender = sender
	subgc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.subscriptionsGroupClient = subgc

	dc := resources.NewDeploymentsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&dc.Client)
	dc.Authorizer = auth
	dc.Sender = sender
	dc.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.deploymentsClient = dc

	ai := appinsights.NewComponentsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&ai.Client)
	ai.Authorizer = auth
	ai.Sender = sender
	ai.SkipResourceProviderRegistration = c.SkipProviderRegistration
	client.appInsightsClient = ai

	client.registerAutomationClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerAuthentication(endpoint, graphEndpoint, c.SubscriptionID, c.TenantID, auth, graphAuth, sender)
	client.registerCDNClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerDatabases(endpoint, c.SubscriptionID, auth, sender)
	client.registerDisks(endpoint, c.SubscriptionID, auth, sender)
	client.registerKeyVaultClients(endpoint, c.SubscriptionID, auth, keyVaultAuth, sender)
	client.registerNetworkingClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerRedisClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerResourcesClients(endpoint, c.SubscriptionID, auth, sender)
	client.registerSchedulerClients(endpoint, c.SubscriptionID, auth)
	client.registerSearchClients(endpoint, c.SubscriptionID, auth)
	client.registerServiceBusClients(endpoint, c.SubscriptionID, auth)
	client.registerStorageClients(endpoint, c.SubscriptionID, auth)
	client.registerTrafficManagerClients(endpoint, c.SubscriptionID, auth)
	client.registerWebClients(endpoint, c.SubscriptionID, auth)

	return &client, nil
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
	endpointsClient := cdn.NewEndpointsClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&endpointsClient.Client)
	endpointsClient.Authorizer = auth
	endpointsClient.Sender = sender
	endpointsClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.cdnEndpointsClient = endpointsClient

	profilesClient := cdn.NewProfilesClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&profilesClient.Client)
	profilesClient.Authorizer = auth
	profilesClient.Sender = sender
	profilesClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.cdnProfilesClient = profilesClient
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
	setUserAgent(&postgresqlConfigClient.Client)
	postgresqlConfigClient.Authorizer = auth
	postgresqlConfigClient.Sender = autorest.CreateSender(withRequestLogging())
	postgresqlConfigClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.postgresqlConfigurationsClient = postgresqlConfigClient

	postgresqlDBClient := postgresql.NewDatabasesClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&postgresqlDBClient.Client)
	postgresqlDBClient.Authorizer = auth
	postgresqlDBClient.Sender = autorest.CreateSender(withRequestLogging())
	postgresqlDBClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.postgresqlDatabasesClient = postgresqlDBClient

	postgresqlFWClient := postgresql.NewFirewallRulesClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&postgresqlFWClient.Client)
	postgresqlFWClient.Authorizer = auth
	postgresqlFWClient.Sender = autorest.CreateSender(withRequestLogging())
	postgresqlFWClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.postgresqlFirewallRulesClient = postgresqlFWClient

	postgresqlSrvClient := postgresql.NewServersClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&postgresqlSrvClient.Client)
	postgresqlSrvClient.Authorizer = auth
	postgresqlSrvClient.Sender = autorest.CreateSender(withRequestLogging())
	postgresqlSrvClient.SkipResourceProviderRegistration = c.skipProviderRegistration
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
}

func (c *ArmClient) registerDisks(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	diskClient := disk.NewDisksClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&diskClient.Client)
	diskClient.Authorizer = auth
	diskClient.Sender = sender
	diskClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.diskClient = diskClient

	snapshotsClient := disk.NewSnapshotsClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&snapshotsClient.Client)
	snapshotsClient.Authorizer = auth
	snapshotsClient.Sender = sender
	snapshotsClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.snapshotsClient = snapshotsClient
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

func (c *ArmClient) registerNetworkingClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	// TODO: move the other networking stuff in here, gradually
	watchersClient := network.NewWatchersClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&watchersClient.Client)
	watchersClient.Authorizer = auth
	watchersClient.Sender = sender
	watchersClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.watcherClient = watchersClient
}

func (c *ArmClient) registerRedisClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	groupsClient := redis.NewGroupClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&groupsClient.Client)
	groupsClient.Authorizer = auth
	groupsClient.Sender = sender
	groupsClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.redisClient = groupsClient

	firewallRuleClient := redis.NewFirewallRuleClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&firewallRuleClient.Client)
	firewallRuleClient.Authorizer = auth
	firewallRuleClient.Sender = sender
	firewallRuleClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.redisFirewallClient = firewallRuleClient

	patchSchedulesClient := redis.NewPatchSchedulesClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&patchSchedulesClient.Client)
	patchSchedulesClient.Authorizer = auth
	patchSchedulesClient.Sender = sender
	patchSchedulesClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.redisPatchSchedulesClient = patchSchedulesClient
}

func (c *ArmClient) registerResourcesClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	locksClient := locks.NewManagementLocksClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&locksClient.Client)
	locksClient.Authorizer = auth
	locksClient.Sender = sender
	locksClient.SkipResourceProviderRegistration = c.skipProviderRegistration
	c.managementLocksClient = locksClient
}

func (c *ArmClient) registerSchedulerClients(endpoint, subscriptionId string, auth autorest.Authorizer) {
	jobsClient := scheduler.NewJobsClientWithBaseURI(endpoint, c.subscriptionId)
	c.configureClient(&jobsClient.Client, auth)
	c.jobsClient = jobsClient

	collectionsClient := scheduler.NewJobCollectionsClientWithBaseURI(endpoint, subscriptionId)
	c.configureClient(&collectionsClient.Client, auth)
	c.jobsCollectionsClient = collectionsClient
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

var (
	storageKeyCacheMu sync.RWMutex
	storageKeyCache   = make(map[string]string)
)

func (armClient *ArmClient) getKeyForStorageAccount(resourceGroupName, storageAccountName string) (string, bool, error) {
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
		ctx := armClient.StopContext
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

func (armClient *ArmClient) getBlobStorageClientForStorageAccount(resourceGroupName, storageAccountName string) (*mainStorage.BlobStorageClient, bool, error) {
	key, accountExists, err := armClient.getKeyForStorageAccount(resourceGroupName, storageAccountName)
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

func (armClient *ArmClient) getFileServiceClientForStorageAccount(resourceGroupName, storageAccountName string) (*mainStorage.FileServiceClient, bool, error) {
	key, accountExists, err := armClient.getKeyForStorageAccount(resourceGroupName, storageAccountName)
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

func (armClient *ArmClient) getTableServiceClientForStorageAccount(resourceGroupName, storageAccountName string) (*mainStorage.TableServiceClient, bool, error) {
	key, accountExists, err := armClient.getKeyForStorageAccount(resourceGroupName, storageAccountName)
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

func (armClient *ArmClient) getQueueServiceClientForStorageAccount(resourceGroupName, storageAccountName string) (*mainStorage.QueueServiceClient, bool, error) {
	key, accountExists, err := armClient.getKeyForStorageAccount(resourceGroupName, storageAccountName)
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
