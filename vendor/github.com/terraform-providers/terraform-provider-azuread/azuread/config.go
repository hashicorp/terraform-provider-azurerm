package azuread

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/go-azure-helpers/sender"
	"github.com/hashicorp/terraform/httpclient"
	"github.com/terraform-providers/terraform-provider-azuread/version"
)

// ArmClient contains the handles to all the specific Azure ADger resource classes' respective clients.
type ArmClient struct {
	subscriptionID string
	clientID       string
	tenantID       string
	environment    azure.Environment

	StopContext context.Context

	// azure AD clients
	applicationsClient      graphrbac.ApplicationsClient
	domainsClient           graphrbac.DomainsClient
	groupsClient            graphrbac.GroupsClient
	servicePrincipalsClient graphrbac.ServicePrincipalsClient
	usersClient             graphrbac.UsersClient
}

// getArmClient is a helper method which returns a fully instantiated *ArmClient based on the auth Config's current settings.
func getArmClient(authCfg *authentication.Config) (*ArmClient, error) {
	env, err := authentication.DetermineEnvironment(authCfg.Environment)
	if err != nil {
		return nil, err
	}

	// client declarations:
	client := ArmClient{
		subscriptionID: authCfg.SubscriptionID,
		clientID:       authCfg.ClientID,
		tenantID:       authCfg.TenantID,
		environment:    *env,
	}

	sender := sender.BuildSender("AzureAD")

	oauth, err := authCfg.BuildOAuthConfig(env.ActiveDirectoryEndpoint)
	if err != nil {
		return nil, err
	}

	// Graph Endpoints
	graphEndpoint := env.GraphEndpoint
	graphAuthorizer, err := authCfg.GetAuthorizationToken(sender, oauth, graphEndpoint)
	if err != nil {
		return nil, err
	}

	client.registerGraphRBACClients(graphEndpoint, authCfg.TenantID, graphAuthorizer)

	return &client, nil
}

func (c *ArmClient) registerGraphRBACClients(endpoint, tenantID string, authorizer autorest.Authorizer) {
	c.applicationsClient = graphrbac.NewApplicationsClientWithBaseURI(endpoint, tenantID)
	configureClient(&c.applicationsClient.Client, authorizer)

	c.domainsClient = graphrbac.NewDomainsClientWithBaseURI(endpoint, tenantID)
	configureClient(&c.domainsClient.Client, authorizer)

	c.groupsClient = graphrbac.NewGroupsClientWithBaseURI(endpoint, tenantID)
	configureClient(&c.groupsClient.Client, authorizer)

	c.servicePrincipalsClient = graphrbac.NewServicePrincipalsClientWithBaseURI(endpoint, tenantID)
	configureClient(&c.servicePrincipalsClient.Client, authorizer)

	c.usersClient = graphrbac.NewUsersClientWithBaseURI(endpoint, tenantID)
	configureClient(&c.usersClient.Client, authorizer)
}

func configureClient(client *autorest.Client, auth autorest.Authorizer) {
	setUserAgent(client)
	client.Authorizer = auth
	client.Sender = sender.BuildSender("AzureAD")
	client.SkipResourceProviderRegistration = false
	client.PollingDuration = 60 * time.Minute
}

//could be moved to helpers
func setUserAgent(client *autorest.Client) {
	// TODO: This is the SDK version not the CLI version, once we are on 0.12, should revisit
	tfUserAgent := httpclient.UserAgentString()

	pv := version.ProviderVersion
	providerUserAgent := fmt.Sprintf("%s terraform-provider-azuread/%s", tfUserAgent, pv)
	client.UserAgent = strings.TrimSpace(fmt.Sprintf("%s %s", client.UserAgent, providerUserAgent))

	// append the CloudShell version to the user agent if it exists
	if azureAgent := os.Getenv("AZURE_HTTP_USER_AGENT"); azureAgent != "" {
		client.UserAgent = fmt.Sprintf("%s %s", client.UserAgent, azureAgent)
	}

	log.Printf("[DEBUG] AzureAD Client User Agent: %s\n", client.UserAgent)
}
