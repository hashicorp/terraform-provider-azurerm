package ar

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/sender"
	"github.com/hashicorp/terraform/httpclient"
	"github.com/terraform-providers/terraform-provider-azurerm/version"
)

type ClientOptions struct {
	Authorizer                 autorest.Authorizer
	ProviderName               string
	PartnerId                  string
	PollingDuration            time.Duration
	SkipProviderReg            bool
	EnableCorrelationRequestID bool
}

func ConfigureClient(c *autorest.Client, o *ClientOptions) {
	if o.PartnerId != "" {
		setUserAgent(c, o.ProviderName, o.PartnerId)
	}

	c.Authorizer = o.Authorizer
	c.Sender = sender.BuildSender(o.ProviderName)
	c.PollingDuration = o.PollingDuration
	c.SkipResourceProviderRegistration = o.SkipProviderReg
	if o.EnableCorrelationRequestID {
		c.RequestInspector = WithCorrelationRequestID(CorrelationRequestID())
	}
}

func setUserAgent(client *autorest.Client, providerName, partnerID string) {
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

	log.Printf("[DEBUG] %s Client User Agent: %s\n", providerName, client.UserAgent)
}
