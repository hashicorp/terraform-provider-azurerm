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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/version"
)

func ConfigureClient(client *autorest.Client, auth autorest.Authorizer, partnerId string, skipProviderReg bool) {
	setUserAgent(client, partnerId)
	client.Authorizer = auth
	client.Sender = sender.BuildSender("AzureRM")
	client.PollingDuration = 180 * time.Minute
	client.SkipResourceProviderRegistration = skipProviderReg
	client.RequestInspector = azure.WithCorrelationRequestID(azure.CorrelationRequestID())

}

func setUserAgent(client *autorest.Client, partnerID string) {
	// TODO: This is the SDK version not the CLI version, once we are on 0.12, should revisit
	tfUserAgent := httpclient.UserAgentString()

	pv := version.ProviderVersion
	providerUserAgent := fmt.Sprintf("%s terraform-provider-azurerm/%s", tfUserAgent, pv)
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
