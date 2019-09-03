package common

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/sender"
	"github.com/hashicorp/terraform/httpclient"
	"github.com/terraform-providers/terraform-provider-azurerm/version"
)

type ClientOptions struct {
	SubscriptionId string
	TenantID       string
	PartnerId      string

	GraphAuthorizer           autorest.Authorizer
	GraphEndpoint             string
	KeyVaultAuthorizer        autorest.Authorizer
	ResourceManagerAuthorizer autorest.Authorizer
	ResourceManagerEndpoint   string
	StorageAuthorizer         autorest.Authorizer

	PollingDuration             time.Duration
	SkipProviderReg             bool
	DisableCorrelationRequestID bool
	Environment                 azure.Environment
}

func (o ClientOptions) ConfigureClient(c *autorest.Client, authorizer autorest.Authorizer) {
	if o.PartnerId != "" {
		setUserAgent(c, o.PartnerId)
	}

	c.Authorizer = authorizer
	c.Sender = sender.BuildSender("AzureRM")
	c.PollingDuration = o.PollingDuration
	c.SkipResourceProviderRegistration = o.SkipProviderReg
	if !o.DisableCorrelationRequestID {
		c.RequestInspector = WithCorrelationRequestID(CorrelationRequestID())
	}
}

func setUserAgent(client *autorest.Client, partnerID string) {
	tfUserAgent := httpclient.UserAgentString()

	providerUserAgent := fmt.Sprintf("%s terraform-provider-azurerm/%s", tfUserAgent, version.ProviderVersion)
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
