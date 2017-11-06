package azurerm

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

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
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/azure-sdk-for-go/arm/resources/subscriptions"
	"github.com/Azure/azure-sdk-for-go/arm/scheduler"
	"github.com/Azure/azure-sdk-for-go/arm/search"
	"github.com/Azure/azure-sdk-for-go/arm/servicebus"
	"github.com/Azure/azure-sdk-for-go/arm/sql"
	"github.com/Azure/azure-sdk-for-go/arm/storage"
	"github.com/Azure/azure-sdk-for-go/arm/trafficmanager"
	"github.com/Azure/azure-sdk-for-go/arm/web"
	keyVault "github.com/Azure/azure-sdk-for-go/dataplane/keyvault"
	mainStorage "github.com/Azure/azure-sdk-for-go/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/terraform/terraform"
)

// ArmClient contains the handles to all the specific Azure Resource Manager
// resource classes' respective clients.
type ArmClient struct {
	clientId              string
	tenantId              string
	subscriptionId        string
	usingServicePrincipal bool
	environment           azure.Environment

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

	jobsClient            scheduler.JobsClient
	jobsCollectionsClient scheduler.JobCollectionsClient

	storageServiceClient storage.AccountsClient
	storageUsageClient   storage.UsageClient

	deploymentsClient resources.DeploymentsClient

	redisClient         redis.GroupClient
	redisFirewallClient redis.FirewallRuleClient

	trafficManagerProfilesClient  trafficmanager.ProfilesClient
	trafficManagerEndpointsClient trafficmanager.EndpointsClient

	searchServicesClient          search.ServicesClient
	serviceBusNamespacesClient    servicebus.NamespacesClient
	serviceBusQueuesClient        servicebus.QueuesClient
	serviceBusTopicsClient        servicebus.TopicsClient
	serviceBusSubscriptionsClient servicebus.SubscriptionsClient

	keyVaultClient           keyvault.VaultsClient
	keyVaultManagementClient keyVault.ManagementClient

	appServicePlansClient web.AppServicePlansClient
	appServicesClient     web.AppsClient

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
	version := terraform.VersionString()
	client.UserAgent = fmt.Sprintf("HashiCorp-Terraform-v%s", version)

	// append the CloudShell version to the user agent
	if azureAgent := os.Getenv("AZURE_HTTP_USER_AGENT"); azureAgent != "" {
		client.UserAgent = fmt.Sprintf("%s;%s", client.UserAgent, azureAgent)
	}

	log.Printf("[DEBUG] User Agent: %s", client.UserAgent)
}

func (c *Config) getAuthorizationToken(oauthConfig *adal.OAuthConfig, endpoint string) (*autorest.BearerAuthorizer, error) {
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
func (c *Config) getArmClient() (*ArmClient, error) {
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
		clientId:              c.ClientID,
		tenantId:              c.TenantID,
		subscriptionId:        c.SubscriptionID,
		environment:           env,
		usingServicePrincipal: c.ClientSecret != "",
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
	auth, err := c.getAuthorizationToken(oauthConfig, endpoint)
	if err != nil {
		return nil, err
	}

	// Graph Endpoints
	graphEndpoint := env.GraphEndpoint
	graphAuth, err := c.getAuthorizationToken(oauthConfig, graphEndpoint)
	if err != nil {
		return nil, err
	}

	// Key Vault Endpoints
	keyVaultAuth := autorest.NewBearerAuthorizerCallback(sender, func(tenantID, resource string) (*autorest.BearerAuthorizer, error) {
		keyVaultSpt, err := c.getAuthorizationToken(oauthConfig, resource)
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
	client.availSetClient = asc

	uoc := compute.NewUsageClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&uoc.Client)
	uoc.Authorizer = auth
	uoc.Sender = sender
	client.usageOpsClient = uoc

	vmeic := compute.NewVirtualMachineExtensionImagesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&vmeic.Client)
	vmeic.Authorizer = auth
	vmeic.Sender = sender
	client.vmExtensionImageClient = vmeic

	vmec := compute.NewVirtualMachineExtensionsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&vmec.Client)
	vmec.Authorizer = auth
	vmec.Sender = sender
	client.vmExtensionClient = vmec

	vmic := compute.NewVirtualMachineImagesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&vmic.Client)
	vmic.Authorizer = auth
	vmic.Sender = sender
	client.vmImageClient = vmic

	vmssc := compute.NewVirtualMachineScaleSetsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&vmssc.Client)
	vmssc.Authorizer = auth
	vmssc.Sender = sender
	client.vmScaleSetClient = vmssc

	vmc := compute.NewVirtualMachinesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&vmc.Client)
	vmc.Authorizer = auth
	vmc.Sender = sender
	client.vmClient = vmc

	agc := network.NewApplicationGatewaysClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&agc.Client)
	agc.Authorizer = auth
	agc.Sender = sender
	client.applicationGatewayClient = agc

	crc := containerregistry.NewRegistriesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&crc.Client)
	crc.Authorizer = auth
	crc.Sender = sender
	client.containerRegistryClient = crc

	csc := containerservice.NewContainerServicesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&csc.Client)
	csc.Authorizer = auth
	csc.Sender = sender
	client.containerServicesClient = csc

	cgc := containerinstance.NewContainerGroupsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&cgc.Client)
	cgc.Authorizer = auth
	cgc.Sender = autorest.CreateSender(withRequestLogging())
	client.containerGroupsClient = cgc

	cdb := cosmosdb.NewDatabaseAccountsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&cdb.Client)
	cdb.Authorizer = auth
	cdb.Sender = sender
	client.cosmosDBClient = cdb

	img := compute.NewImagesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&img.Client)
	img.Authorizer = auth
	img.Sender = sender
	client.imageClient = img

	egtc := eventgrid.NewTopicsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&egtc.Client)
	egtc.Authorizer = auth
	egtc.Sender = sender
	client.eventGridTopicsClient = egtc

	ehc := eventhub.NewEventHubsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&ehc.Client)
	ehc.Authorizer = auth
	ehc.Sender = sender
	client.eventHubClient = ehc

	chcgc := eventhub.NewConsumerGroupsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&chcgc.Client)
	chcgc.Authorizer = auth
	chcgc.Sender = sender
	client.eventHubConsumerGroupClient = chcgc

	ehnc := eventhub.NewNamespacesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&ehnc.Client)
	ehnc.Authorizer = auth
	ehnc.Sender = sender
	client.eventHubNamespacesClient = ehnc

	ifc := network.NewInterfacesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&ifc.Client)
	ifc.Authorizer = auth
	ifc.Sender = sender
	client.ifaceClient = ifc

	erc := network.NewExpressRouteCircuitsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&erc.Client)
	erc.Authorizer = auth
	erc.Sender = sender
	client.expressRouteCircuitClient = erc

	lbc := network.NewLoadBalancersClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&lbc.Client)
	lbc.Authorizer = auth
	lbc.Sender = sender
	client.loadBalancerClient = lbc

	lgc := network.NewLocalNetworkGatewaysClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&lgc.Client)
	lgc.Authorizer = auth
	lgc.Sender = sender
	client.localNetConnClient = lgc

	opwc := operationalinsights.NewWorkspacesClient(c.SubscriptionID)
	setUserAgent(&opwc.Client)
	opwc.Authorizer = auth
	opwc.Sender = autorest.CreateSender(withRequestLogging())
	client.workspacesClient = opwc

	pipc := network.NewPublicIPAddressesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&pipc.Client)
	pipc.Authorizer = auth
	pipc.Sender = sender
	client.publicIPClient = pipc

	sgc := network.NewSecurityGroupsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&sgc.Client)
	sgc.Authorizer = auth
	sgc.Sender = sender
	client.secGroupClient = sgc

	src := network.NewSecurityRulesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&src.Client)
	src.Authorizer = auth
	src.Sender = sender
	client.secRuleClient = src

	snc := network.NewSubnetsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&snc.Client)
	snc.Authorizer = auth
	snc.Sender = sender
	client.subnetClient = snc

	vgcc := network.NewVirtualNetworkGatewayConnectionsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&vgcc.Client)
	vgcc.Authorizer = auth
	vgcc.Sender = sender
	client.vnetGatewayConnectionsClient = vgcc

	vgc := network.NewVirtualNetworkGatewaysClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&vgc.Client)
	vgc.Authorizer = auth
	vgc.Sender = sender
	client.vnetGatewayClient = vgc

	vnc := network.NewVirtualNetworksClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&vnc.Client)
	vnc.Authorizer = auth
	vnc.Sender = sender
	client.vnetClient = vnc

	vnpc := network.NewVirtualNetworkPeeringsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&vnpc.Client)
	vnpc.Authorizer = auth
	vnpc.Sender = sender
	client.vnetPeeringsClient = vnpc

	rtc := network.NewRouteTablesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&rtc.Client)
	rtc.Authorizer = auth
	rtc.Sender = sender
	client.routeTablesClient = rtc

	rc := network.NewRoutesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&rc.Client)
	rc.Authorizer = auth
	rc.Sender = sender
	client.routesClient = rc

	dn := dns.NewRecordSetsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&dn.Client)
	dn.Authorizer = auth
	dn.Sender = sender
	client.dnsClient = dn

	zo := dns.NewZonesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&zo.Client)
	zo.Authorizer = auth
	zo.Sender = sender
	client.zonesClient = zo

	rgc := resources.NewGroupsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&rgc.Client)
	rgc.Authorizer = auth
	rgc.Sender = sender
	client.resourceGroupClient = rgc

	pc := resources.NewProvidersClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&pc.Client)
	pc.Authorizer = auth
	pc.Sender = sender
	client.providers = pc

	tc := resources.NewTagsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&tc.Client)
	tc.Authorizer = auth
	tc.Sender = sender
	client.tagsClient = tc

	rf := resources.NewGroupClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&rf.Client)
	rf.Authorizer = auth
	rf.Sender = sender
	client.resourceFindClient = rf

	subgc := subscriptions.NewGroupClientWithBaseURI(endpoint)
	setUserAgent(&subgc.Client)
	subgc.Authorizer = auth
	subgc.Sender = sender
	client.subscriptionsGroupClient = subgc

	jc := scheduler.NewJobsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&jc.Client)
	jc.Authorizer = auth
	jc.Sender = sender
	client.jobsClient = jc

	jcc := scheduler.NewJobCollectionsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&jcc.Client)
	jcc.Authorizer = auth
	jcc.Sender = sender
	client.jobsCollectionsClient = jcc

	ssc := storage.NewAccountsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&ssc.Client)
	ssc.Authorizer = auth
	ssc.Sender = sender
	client.storageServiceClient = ssc

	suc := storage.NewUsageClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&suc.Client)
	suc.Authorizer = auth
	suc.Sender = sender
	client.storageUsageClient = suc

	cpc := cdn.NewProfilesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&cpc.Client)
	cpc.Authorizer = auth
	cpc.Sender = sender
	client.cdnProfilesClient = cpc

	cec := cdn.NewEndpointsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&cec.Client)
	cec.Authorizer = auth
	cec.Sender = sender
	client.cdnEndpointsClient = cec

	dc := resources.NewDeploymentsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&dc.Client)
	dc.Authorizer = auth
	dc.Sender = sender
	client.deploymentsClient = dc

	tmpc := trafficmanager.NewProfilesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&tmpc.Client)
	tmpc.Authorizer = auth
	tmpc.Sender = sender
	client.trafficManagerProfilesClient = tmpc

	tmec := trafficmanager.NewEndpointsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&tmec.Client)
	tmec.Authorizer = auth
	tmec.Sender = sender
	client.trafficManagerEndpointsClient = tmec

	sesc := search.NewServicesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&sesc.Client)
	sesc.Authorizer = auth
	sesc.Sender = sender
	client.searchServicesClient = sesc

	sbnc := servicebus.NewNamespacesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&sbnc.Client)
	sbnc.Authorizer = auth
	sbnc.Sender = sender
	client.serviceBusNamespacesClient = sbnc

	sbqc := servicebus.NewQueuesClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&sbqc.Client)
	sbqc.Authorizer = auth
	sbqc.Sender = sender
	client.serviceBusQueuesClient = sbqc

	sbtc := servicebus.NewTopicsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&sbtc.Client)
	sbtc.Authorizer = auth
	sbtc.Sender = sender
	client.serviceBusTopicsClient = sbtc

	sbsc := servicebus.NewSubscriptionsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&sbsc.Client)
	sbsc.Authorizer = auth
	sbsc.Sender = sender
	client.serviceBusSubscriptionsClient = sbsc

	aspc := web.NewAppServicePlansClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&aspc.Client)
	aspc.Authorizer = auth
	aspc.Sender = sender
	client.appServicePlansClient = aspc

	ac := web.NewAppsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&ac.Client)
	ac.Authorizer = auth
	ac.Sender = autorest.CreateSender(withRequestLogging())
	client.appServicesClient = ac

	ai := appinsights.NewComponentsClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&ai.Client)
	ai.Authorizer = auth
	ai.Sender = sender
	client.appInsightsClient = ai

	aadb := automation.NewAccountClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&aadb.Client)
	aadb.Authorizer = auth
	aadb.Sender = sender
	client.automationAccountClient = aadb

	arc := automation.NewRunbookClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&arc.Client)
	arc.Authorizer = auth
	arc.Sender = sender
	client.automationRunbookClient = arc

	acc := automation.NewCredentialClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&acc.Client)
	acc.Authorizer = auth
	acc.Sender = sender
	client.automationCredentialClient = acc

	aschc := automation.NewScheduleClientWithBaseURI(endpoint, c.SubscriptionID)
	setUserAgent(&aschc.Client)
	aschc.Authorizer = auth
	aschc.Sender = sender
	client.automationScheduleClient = aschc

	client.registerAuthentication(endpoint, graphEndpoint, c.SubscriptionID, c.TenantID, auth, graphAuth, sender)
	client.registerDatabases(endpoint, c.SubscriptionID, auth, sender)
	client.registerDisks(endpoint, c.SubscriptionID, auth, sender)
	client.registerKeyVaultClients(endpoint, c.SubscriptionID, auth, keyVaultAuth, sender)
	client.registerRedisClients(endpoint, c.SubscriptionID, auth, sender)

	return &client, nil
}

func (c *ArmClient) registerAuthentication(endpoint, graphEndpoint, subscriptionId, tenantId string, auth, graphAuth autorest.Authorizer, sender autorest.Sender) {
	spc := graphrbac.NewServicePrincipalsClientWithBaseURI(graphEndpoint, tenantId)
	setUserAgent(&spc.Client)
	spc.Authorizer = graphAuth
	spc.Sender = sender
	c.servicePrincipalsClient = spc

	rac := authorization.NewRoleAssignmentsClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&rac.Client)
	rac.Authorizer = auth
	rac.Sender = sender
	c.roleAssignmentsClient = rac

	rdc := authorization.NewRoleDefinitionsClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&rdc.Client)
	rdc.Authorizer = auth
	rdc.Sender = sender
	c.roleDefinitionsClient = rdc
}

func (c *ArmClient) registerDatabases(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	// MySQL
	mysqlConfigClient := mysql.NewConfigurationsClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&mysqlConfigClient.Client)
	mysqlConfigClient.Authorizer = auth
	mysqlConfigClient.Sender = sender
	c.mysqlConfigurationsClient = mysqlConfigClient

	mysqlDBClient := mysql.NewDatabasesClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&mysqlDBClient.Client)
	mysqlDBClient.Authorizer = auth
	mysqlDBClient.Sender = sender
	c.mysqlDatabasesClient = mysqlDBClient

	mysqlFWClient := mysql.NewFirewallRulesClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&mysqlFWClient.Client)
	mysqlFWClient.Authorizer = auth
	mysqlFWClient.Sender = sender
	c.mysqlFirewallRulesClient = mysqlFWClient

	mysqlServersClient := mysql.NewServersClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&mysqlServersClient.Client)
	mysqlServersClient.Authorizer = auth
	mysqlServersClient.Sender = sender
	c.mysqlServersClient = mysqlServersClient

	// PostgreSQL
	postgresqlConfigClient := postgresql.NewConfigurationsClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&postgresqlConfigClient.Client)
	postgresqlConfigClient.Authorizer = auth
	postgresqlConfigClient.Sender = autorest.CreateSender(withRequestLogging())
	c.postgresqlConfigurationsClient = postgresqlConfigClient

	postgresqlDBClient := postgresql.NewDatabasesClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&postgresqlDBClient.Client)
	postgresqlDBClient.Authorizer = auth
	postgresqlDBClient.Sender = autorest.CreateSender(withRequestLogging())
	c.postgresqlDatabasesClient = postgresqlDBClient

	postgresqlFWClient := postgresql.NewFirewallRulesClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&postgresqlFWClient.Client)
	postgresqlFWClient.Authorizer = auth
	postgresqlFWClient.Sender = autorest.CreateSender(withRequestLogging())
	c.postgresqlFirewallRulesClient = postgresqlFWClient

	postgresqlSrvClient := postgresql.NewServersClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&postgresqlSrvClient.Client)
	postgresqlSrvClient.Authorizer = auth
	postgresqlSrvClient.Sender = autorest.CreateSender(withRequestLogging())
	c.postgresqlServersClient = postgresqlSrvClient

	// SQL Azure
	sqlDBClient := sql.NewDatabasesClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&sqlDBClient.Client)
	sqlDBClient.Authorizer = auth
	sqlDBClient.Sender = sender
	c.sqlDatabasesClient = sqlDBClient

	sqlFWClient := sql.NewFirewallRulesClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&sqlFWClient.Client)
	sqlFWClient.Authorizer = auth
	sqlFWClient.Sender = sender
	c.sqlFirewallRulesClient = sqlFWClient

	sqlEPClient := sql.NewElasticPoolsClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&sqlEPClient.Client)
	sqlEPClient.Authorizer = auth
	sqlEPClient.Sender = sender
	c.sqlElasticPoolsClient = sqlEPClient

	sqlSrvClient := sql.NewServersClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&sqlSrvClient.Client)
	sqlSrvClient.Authorizer = auth
	sqlSrvClient.Sender = sender
	c.sqlServersClient = sqlSrvClient
}

func (c *ArmClient) registerDisks(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	diskClient := disk.NewDisksClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&diskClient.Client)
	diskClient.Authorizer = auth
	diskClient.Sender = sender
	c.diskClient = diskClient

	snapshotsClient := disk.NewSnapshotsClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&snapshotsClient.Client)
	snapshotsClient.Authorizer = auth
	snapshotsClient.Sender = sender
	c.snapshotsClient = snapshotsClient
}

func (c *ArmClient) registerKeyVaultClients(endpoint, subscriptionId string, auth autorest.Authorizer, keyVaultAuth autorest.Authorizer, sender autorest.Sender) {
	keyVaultClient := keyvault.NewVaultsClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&keyVaultClient.Client)
	keyVaultClient.Authorizer = auth
	keyVaultClient.Sender = sender
	c.keyVaultClient = keyVaultClient

	keyVaultManagementClient := keyVault.New()
	setUserAgent(&keyVaultManagementClient.Client)
	keyVaultManagementClient.Authorizer = keyVaultAuth
	keyVaultManagementClient.Sender = sender
	c.keyVaultManagementClient = keyVaultManagementClient
}

func (c *ArmClient) registerRedisClients(endpoint, subscriptionId string, auth autorest.Authorizer, sender autorest.Sender) {
	rdc := redis.NewGroupClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&rdc.Client)
	rdc.Authorizer = auth
	rdc.Sender = sender
	c.redisClient = rdc

	rdfc := redis.NewFirewallRuleClientWithBaseURI(endpoint, subscriptionId)
	setUserAgent(&rdfc.Client)
	rdfc.Authorizer = auth
	rdfc.Sender = sender
	c.redisFirewallClient = rdfc
}

func (armClient *ArmClient) getKeyForStorageAccount(resourceGroupName, storageAccountName string) (string, bool, error) {
	accountKeys, err := armClient.storageServiceClient.ListKeys(resourceGroupName, storageAccountName)
	if accountKeys.StatusCode == http.StatusNotFound {
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
	return *keys[0].Value, true, nil
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
