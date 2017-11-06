package azurerm

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure/cli"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	var p *schema.Provider
	p = &schema.Provider{
		Schema: map[string]*schema.Schema{
			"subscription_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_SUBSCRIPTION_ID", ""),
			},

			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_ID", ""),
			},

			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_SECRET", ""),
			},

			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_TENANT_ID", ""),
			},

			"environment": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_ENVIRONMENT", "public"),
			},

			"skip_credentials_validation": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_SKIP_CREDENTIALS_VALIDATION", false),
			},

			"skip_provider_registration": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_SKIP_PROVIDER_REGISTRATION", false),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"azurerm_builtin_role_definition": dataSourceArmBuiltInRoleDefinition(),
			"azurerm_client_config":           dataSourceArmClientConfig(),
			"azurerm_image":                   dataSourceArmImage(),
			"azurerm_key_vault_access_policy": dataSourceArmKeyVaultAccessPolicy(),
			"azurerm_managed_disk":            dataSourceArmManagedDisk(),
			"azurerm_platform_image":          dataSourceArmPlatformImage(),
			"azurerm_public_ip":               dataSourceArmPublicIP(),
			"azurerm_resource_group":          dataSourceArmResourceGroup(),
			"azurerm_role_definition":         dataSourceArmRoleDefinition(),
			"azurerm_snapshot":                dataSourceArmSnapshot(),
			"azurerm_subnet":                  dataSourceArmSubnet(),
			"azurerm_subscription":            dataSourceArmSubscription(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"azurerm_application_gateway":         resourceArmApplicationGateway(),
			"azurerm_application_insights":        resourceArmApplicationInsights(),
			"azurerm_app_service":                 resourceArmAppService(),
			"azurerm_app_service_plan":            resourceArmAppServicePlan(),
			"azurerm_automation_account":          resourceArmAutomationAccount(),
			"azurerm_automation_credential":       resourceArmAutomationCredential(),
			"azurerm_automation_runbook":          resourceArmAutomationRunbook(),
			"azurerm_automation_schedule":         resourceArmAutomationSchedule(),
			"azurerm_availability_set":            resourceArmAvailabilitySet(),
			"azurerm_cdn_endpoint":                resourceArmCdnEndpoint(),
			"azurerm_cdn_profile":                 resourceArmCdnProfile(),
			"azurerm_container_registry":          resourceArmContainerRegistry(),
			"azurerm_container_service":           resourceArmContainerService(),
			"azurerm_container_group":             resourceArmContainerGroup(),
			"azurerm_cosmosdb_account":            resourceArmCosmosDBAccount(),
			"azurerm_dns_a_record":                resourceArmDnsARecord(),
			"azurerm_dns_aaaa_record":             resourceArmDnsAAAARecord(),
			"azurerm_dns_cname_record":            resourceArmDnsCNameRecord(),
			"azurerm_dns_mx_record":               resourceArmDnsMxRecord(),
			"azurerm_dns_ns_record":               resourceArmDnsNsRecord(),
			"azurerm_dns_ptr_record":              resourceArmDnsPtrRecord(),
			"azurerm_dns_srv_record":              resourceArmDnsSrvRecord(),
			"azurerm_dns_txt_record":              resourceArmDnsTxtRecord(),
			"azurerm_dns_zone":                    resourceArmDnsZone(),
			"azurerm_eventgrid_topic":             resourceArmEventGridTopic(),
			"azurerm_eventhub":                    resourceArmEventHub(),
			"azurerm_eventhub_authorization_rule": resourceArmEventHubAuthorizationRule(),
			"azurerm_eventhub_consumer_group":     resourceArmEventHubConsumerGroup(),
			"azurerm_eventhub_namespace":          resourceArmEventHubNamespace(),
			"azurerm_express_route_circuit":       resourceArmExpressRouteCircuit(),
			"azurerm_image":                       resourceArmImage(),
			"azurerm_key_vault":                   resourceArmKeyVault(),
			"azurerm_key_vault_certificate":       resourceArmKeyVaultCertificate(),
			"azurerm_key_vault_key":               resourceArmKeyVaultKey(),
			"azurerm_key_vault_secret":            resourceArmKeyVaultSecret(),
			"azurerm_lb":                          resourceArmLoadBalancer(),
			"azurerm_lb_backend_address_pool":     resourceArmLoadBalancerBackendAddressPool(),
			"azurerm_lb_nat_rule":                 resourceArmLoadBalancerNatRule(),
			"azurerm_lb_nat_pool":                 resourceArmLoadBalancerNatPool(),
			"azurerm_lb_probe":                    resourceArmLoadBalancerProbe(),
			"azurerm_lb_rule":                     resourceArmLoadBalancerRule(),
			"azurerm_local_network_gateway":       resourceArmLocalNetworkGateway(),
			"azurerm_log_analytics_workspace":     resourceArmLogAnalyticsWorkspace(),
			"azurerm_managed_disk":                resourceArmManagedDisk(),
			"azurerm_mysql_configuration":         resourceArmMySQLConfiguration(),
			"azurerm_mysql_database":              resourceArmMySqlDatabase(),
			"azurerm_mysql_firewall_rule":         resourceArmMySqlFirewallRule(),
			"azurerm_mysql_server":                resourceArmMySqlServer(),
			"azurerm_network_interface":           resourceArmNetworkInterface(),
			"azurerm_network_security_group":      resourceArmNetworkSecurityGroup(),
			"azurerm_network_security_rule":       resourceArmNetworkSecurityRule(),
			"azurerm_postgresql_configuration":    resourceArmPostgreSQLConfiguration(),
			"azurerm_postgresql_database":         resourceArmPostgreSQLDatabase(),
			"azurerm_postgresql_firewall_rule":    resourceArmPostgreSQLFirewallRule(),
			"azurerm_postgresql_server":           resourceArmPostgreSQLServer(),
			"azurerm_public_ip":                   resourceArmPublicIp(),
			"azurerm_redis_cache":                 resourceArmRedisCache(),
			"azurerm_redis_firewall_rule":         resourceArmRedisFirewallRule(),
			"azurerm_resource_group":              resourceArmResourceGroup(),
			"azurerm_role_assignment":             resourceArmRoleAssignment(),
			"azurerm_role_definition":             resourceArmRoleDefinition(),
			"azurerm_route":                       resourceArmRoute(),
			"azurerm_route_table":                 resourceArmRouteTable(),
			"azurerm_search_service":              resourceArmSearchService(),
			"azurerm_servicebus_namespace":        resourceArmServiceBusNamespace(),
			"azurerm_servicebus_queue":            resourceArmServiceBusQueue(),
			"azurerm_servicebus_subscription":     resourceArmServiceBusSubscription(),
			"azurerm_servicebus_topic":            resourceArmServiceBusTopic(),
			"azurerm_snapshot":                    resourceArmSnapshot(),
			"azurerm_sql_database":                resourceArmSqlDatabase(),
			"azurerm_sql_elasticpool":             resourceArmSqlElasticPool(),
			"azurerm_sql_firewall_rule":           resourceArmSqlFirewallRule(),
			"azurerm_sql_server":                  resourceArmSqlServer(),
			"azurerm_storage_account":             resourceArmStorageAccount(),
			"azurerm_storage_blob":                resourceArmStorageBlob(),
			"azurerm_storage_container":           resourceArmStorageContainer(),
			"azurerm_storage_share":               resourceArmStorageShare(),
			"azurerm_storage_queue":               resourceArmStorageQueue(),
			"azurerm_storage_table":               resourceArmStorageTable(),
			"azurerm_subnet":                      resourceArmSubnet(),
			"azurerm_template_deployment":         resourceArmTemplateDeployment(),
			"azurerm_traffic_manager_endpoint":    resourceArmTrafficManagerEndpoint(),
			"azurerm_traffic_manager_profile":     resourceArmTrafficManagerProfile(),
			"azurerm_virtual_machine_extension":   resourceArmVirtualMachineExtensions(),
			"azurerm_virtual_machine":             resourceArmVirtualMachine(),
			"azurerm_virtual_machine_scale_set":   resourceArmVirtualMachineScaleSet(),
			"azurerm_virtual_network":             resourceArmVirtualNetwork(),
			"azurerm_virtual_network_peering":     resourceArmVirtualNetworkPeering(),
		},
	}

	p.ConfigureFunc = providerConfigure(p)

	return p
}

// Config is the configuration structure used to instantiate a
// new Azure management client.
type Config struct {
	ManagementURL string

	// Core
	ClientID                  string
	SubscriptionID            string
	TenantID                  string
	Environment               string
	SkipCredentialsValidation bool
	SkipProviderRegistration  bool

	// Service Principal Auth
	ClientSecret string

	// Bearer Auth
	AccessToken  *adal.Token
	IsCloudShell bool
}

func (c *Config) validateServicePrincipal() error {
	var err *multierror.Error

	if c.SubscriptionID == "" {
		err = multierror.Append(err, fmt.Errorf("Subscription ID must be configured for the AzureRM provider"))
	}
	if c.ClientID == "" {
		err = multierror.Append(err, fmt.Errorf("Client ID must be configured for the AzureRM provider"))
	}
	if c.ClientSecret == "" {
		err = multierror.Append(err, fmt.Errorf("Client Secret must be configured for the AzureRM provider"))
	}
	if c.TenantID == "" {
		err = multierror.Append(err, fmt.Errorf("Tenant ID must be configured for the AzureRM provider"))
	}
	if c.Environment == "" {
		err = multierror.Append(err, fmt.Errorf("Environment must be configured for the AzureRM provider"))
	}

	return err.ErrorOrNil()
}

func (c *Config) validateBearerAuth() error {
	var err *multierror.Error

	if c.AccessToken == nil {
		err = multierror.Append(err, fmt.Errorf("Access Token was not found in your Azure CLI Credentials.\n\nPlease login to the Azure CLI again via `az login`"))
	}

	if c.ClientID == "" {
		err = multierror.Append(err, fmt.Errorf("Client ID was not found in your Azure CLI Credentials.\n\nPlease login to the Azure CLI again via `az login`"))
	}

	if c.SubscriptionID == "" {
		err = multierror.Append(err, fmt.Errorf("Subscription ID was not found in your Azure CLI Credentials.\n\nPlease login to the Azure CLI again via `az login`"))
	}

	if c.TenantID == "" {
		err = multierror.Append(err, fmt.Errorf("Tenant ID was not found in your Azure CLI Credentials.\n\nPlease login to the Azure CLI again via `az login`"))
	}

	return err.ErrorOrNil()
}

func (c *Config) LoadTokensFromAzureCLI() error {
	profilePath, err := cli.ProfilePath()
	if err != nil {
		return fmt.Errorf("Error loading the Profile Path from the Azure CLI: %+v", err)
	}

	profile, err := cli.LoadProfile(profilePath)
	if err != nil {
		return fmt.Errorf("Azure CLI Authorization Profile was not found. Please ensure the Azure CLI is installed and then log-in with `az login`.")
	}

	// pull out the TenantID and Subscription ID from the Azure Profile
	for _, subscription := range profile.Subscriptions {
		if subscription.IsDefault {
			c.SubscriptionID = subscription.ID
			c.TenantID = subscription.TenantID
			c.Environment = normalizeEnvironmentName(subscription.EnvironmentName)
			break
		}
	}

	foundToken := false
	if c.TenantID != "" {
		// pull out the ClientID and the AccessToken from the Azure Access Token
		tokensPath, err := cli.AccessTokensPath()
		if err != nil {
			return fmt.Errorf("Error loading the Tokens Path from the Azure CLI: %+v", err)
		}

		tokens, err := cli.LoadTokens(tokensPath)
		if err != nil {
			return fmt.Errorf("Azure CLI Authorization Tokens were not found. Please ensure the Azure CLI is installed and then log-in with `az login`.")
		}

		for _, accessToken := range tokens {
			token, err := accessToken.ToADALToken()
			if err != nil {
				return fmt.Errorf("[DEBUG] Error converting access token to token: %+v", err)
			}

			expirationDate, err := cli.ParseExpirationDate(accessToken.ExpiresOn)
			if err != nil {
				return fmt.Errorf("Error parsing expiration date: %q", accessToken.ExpiresOn)
			}

			if expirationDate.UTC().Before(time.Now().UTC()) {
				log.Printf("[DEBUG] Token '%s' has expired", token.AccessToken)
				continue
			}

			if !strings.Contains(accessToken.Resource, "management") {
				log.Printf("[DEBUG] Resource '%s' isn't a management domain", accessToken.Resource)
				continue
			}

			if !strings.HasSuffix(accessToken.Authority, c.TenantID) {
				log.Printf("[DEBUG] Resource '%s' isn't for the correct Tenant", accessToken.Resource)
				continue
			}

			c.ClientID = accessToken.ClientID
			c.AccessToken = &token
			c.IsCloudShell = accessToken.RefreshToken == ""
			foundToken = true
			break
		}
	}

	if !foundToken {
		return fmt.Errorf("No valid (unexpired) Azure CLI Auth Tokens found. Please run `az login`.")
	}

	return nil
}

func normalizeEnvironmentName(input string) string {
	// Environment is stored as `Azure{Environment}Cloud`
	output := strings.ToLower(input)
	output = strings.TrimPrefix(output, "azure")
	output = strings.TrimSuffix(output, "cloud")

	// however Azure Public is `AzureCloud` in the CLI Profile and not `AzurePublicCloud`.
	if output == "" {
		return "public"
	}
	return output
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		config := &Config{
			SubscriptionID:            d.Get("subscription_id").(string),
			ClientID:                  d.Get("client_id").(string),
			ClientSecret:              d.Get("client_secret").(string),
			TenantID:                  d.Get("tenant_id").(string),
			Environment:               d.Get("environment").(string),
			SkipCredentialsValidation: d.Get("skip_credentials_validation").(bool),
			SkipProviderRegistration:  d.Get("skip_provider_registration").(bool),
		}

		if config.ClientSecret != "" {
			log.Printf("[DEBUG] Client Secret specified - using Service Principal for Authentication")
			if err := config.validateServicePrincipal(); err != nil {
				return nil, err
			}
		} else {
			log.Printf("[DEBUG] No Client Secret specified - loading credentials from Azure CLI")
			if err := config.LoadTokensFromAzureCLI(); err != nil {
				return nil, err
			}

			if err := config.validateBearerAuth(); err != nil {
				return nil, fmt.Errorf("Please specify either a Service Principal, or log in with the Azure CLI (using `az login`)")
			}
		}

		client, err := config.getArmClient()
		if err != nil {
			return nil, err
		}

		client.StopContext = p.StopContext()

		// replaces the context between tests
		p.MetaReset = func() error {
			client.StopContext = p.StopContext()
			return nil
		}

		if !config.SkipCredentialsValidation {
			// List all the available providers and their registration state to avoid unnecessary
			// requests. This also lets us check if the provider credentials are correct.
			providerList, err := client.providers.List(nil, "")
			if err != nil {
				return nil, fmt.Errorf("Unable to list provider registration status, it is possible that this is due to invalid "+
					"credentials or the service principal does not have permission to use the Resource Manager API, Azure "+
					"error: %s", err)
			}

			if !config.SkipProviderRegistration {
				err = registerAzureResourceProvidersWithSubscription(*providerList.Value, client.providers)
				if err != nil {
					return nil, err
				}
			}
		}

		return client, nil
	}
}

func registerProviderWithSubscription(providerName string, client resources.ProvidersClient) error {
	_, err := client.Register(providerName)
	if err != nil {
		return fmt.Errorf("Cannot register provider %s with Azure Resource Manager: %s.", providerName, err)
	}

	return nil
}

func determineAzureResourceProvidersToRegister(providerList []resources.Provider) map[string]struct{} {
	providers := map[string]struct{}{
		"Microsoft.Authorization":       {},
		"Microsoft.Automation":          {},
		"Microsoft.Cache":               {},
		"Microsoft.Cdn":                 {},
		"Microsoft.Compute":             {},
		"Microsoft.ContainerInstance":   {},
		"Microsoft.ContainerRegistry":   {},
		"Microsoft.ContainerService":    {},
		"Microsoft.DBforMySQL":          {},
		"Microsoft.DBforPostgreSQL":     {},
		"Microsoft.DocumentDB":          {},
		"Microsoft.EventGrid":           {},
		"Microsoft.EventHub":            {},
		"Microsoft.KeyVault":            {},
		"microsoft.insights":            {},
		"Microsoft.Network":             {},
		"Microsoft.OperationalInsights": {},
		"Microsoft.Resources":           {},
		"Microsoft.Search":              {},
		"Microsoft.ServiceBus":          {},
		"Microsoft.Sql":                 {},
		"Microsoft.Storage":             {},
	}

	// filter out any providers already registered
	for _, p := range providerList {
		if _, ok := providers[*p.Namespace]; !ok {
			continue
		}

		if strings.ToLower(*p.RegistrationState) == "registered" {
			log.Printf("[DEBUG] Skipping provider registration for namespace %s\n", *p.Namespace)
			delete(providers, *p.Namespace)
		}
	}

	return providers
}

// registerAzureResourceProvidersWithSubscription uses the providers client to register
// all Azure resource providers which the Terraform provider may require (regardless of
// whether they are actually used by the configuration or not). It was confirmed by Microsoft
// that this is the approach their own internal tools also take.
func registerAzureResourceProvidersWithSubscription(providerList []resources.Provider, client resources.ProvidersClient) error {
	providers := determineAzureResourceProvidersToRegister(providerList)

	var err error
	var wg sync.WaitGroup
	wg.Add(len(providers))

	for providerName := range providers {
		go func(p string) {
			defer wg.Done()
			log.Printf("[DEBUG] Registering provider with namespace %s\n", p)
			if innerErr := registerProviderWithSubscription(p, client); err != nil {
				err = innerErr
			}
		}(providerName)
	}

	wg.Wait()

	return err
}

// armMutexKV is the instance of MutexKV for ARM resources
var armMutexKV = mutexkv.NewMutexKV()

// Resource group names can be capitalised, but we store them in lowercase.
// Use a custom diff function to avoid creation of new resources.
func resourceAzurermResourceGroupNameDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	return strings.ToLower(old) == strings.ToLower(new)
}

// ignoreCaseDiffSuppressFunc is a DiffSuppressFunc from helper/schema that is
// used to ignore any case-changes in a return value.
func ignoreCaseDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	return strings.ToLower(old) == strings.ToLower(new)
}

// ignoreCaseStateFunc is a StateFunc from helper/schema that converts the
// supplied value to lower before saving to state for consistency.
func ignoreCaseStateFunc(val interface{}) string {
	return strings.ToLower(val.(string))
}

func userDataStateFunc(v interface{}) string {
	switch s := v.(type) {
	case string:
		s = base64Encode(s)
		hash := sha1.Sum([]byte(s))
		return hex.EncodeToString(hash[:])
	default:
		return ""
	}
}

// base64Encode encodes data if the input isn't already encoded using
// base64.StdEncoding.EncodeToString. If the input is already base64 encoded,
// return the original input unchanged.
func base64Encode(data string) string {
	// Check whether the data is already Base64 encoded; don't double-encode
	if isBase64Encoded(data) {
		return data
	}
	// data has not been encoded encode and return
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func isBase64Encoded(data string) bool {
	_, err := base64.StdEncoding.DecodeString(data)
	return err == nil
}
