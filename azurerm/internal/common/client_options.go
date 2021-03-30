package common

import (
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/policy"
	"github.com/hashicorp/go-azure-helpers/sender"
	"github.com/hashicorp/terraform-plugin-sdk/meta"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/version"
)

type ClientOptions struct {
	SubscriptionId   string
	TenantID         string
	PartnerId        string
	TerraformVersion string

	GraphAuthorizer           autorest.Authorizer
	GraphEndpoint             string
	KeyVaultAuthorizer        autorest.Authorizer
	ResourceManagerAuthorizer autorest.Authorizer
	ResourceManagerEndpoint   string
	StorageAuthorizer         autorest.Authorizer
	SynapseAuthorizer         autorest.Authorizer

	SkipProviderReg             bool
	CustomCorrelationRequestID  string
	DisableCorrelationRequestID bool
	DisableTerraformPartnerID   bool
	Environment                 azure.Environment
	Features                    features.UserFeatures
	StorageUseAzureAD           bool

	// Track 2
	ResourceManagerConnection *armcore.Connection
}

func (o ClientOptions) ConfigureClient(c *autorest.Client, authorizer autorest.Authorizer) {
	setUserAgent(c, o.TerraformVersion, o.PartnerId, o.DisableTerraformPartnerID)

	c.Authorizer = authorizer
	c.Sender = sender.BuildSender("AzureRM")
	c.SkipResourceProviderRegistration = o.SkipProviderReg
	if !o.DisableCorrelationRequestID {
		id := o.CustomCorrelationRequestID
		if id == "" {
			id = correlationRequestID()
		}
		c.RequestInspector = withCorrelationRequestID(id)
	}
}

func (o ClientOptions) BuildResourceManagerConnection(cred azcore.TokenCredential) {
	// get the policies
	p := azcore.NewPipeline(nil,
		azcore.NewTelemetryPolicy(&azcore.TelemetryOptions{
			Value:         TelemetryValue(o.TerraformVersion, o.PartnerId, o.DisableTerraformPartnerID),
			ApplicationID: "",
			Disabled:      false,
		}),
		policy.NewRequestLoggingPolicy(),
	)

	o.ResourceManagerConnection = armcore.NewConnectionWithPipeline(o.ResourceManagerEndpoint, p)
}

func setUserAgent(client *autorest.Client, tfVersion, partnerID string, disableTerraformPartnerID bool) {
	tfUserAgent := fmt.Sprintf("HashiCorp Terraform/%s (+https://www.terraform.io) Terraform Plugin SDK/%s", tfVersion, meta.SDKVersionString())

	providerUserAgent := fmt.Sprintf("%s terraform-provider-azurerm/%s", tfUserAgent, version.ProviderVersion)
	client.UserAgent = strings.TrimSpace(fmt.Sprintf("%s %s", client.UserAgent, providerUserAgent))

	// append the CloudShell version to the user agent if it exists
	if azureAgent := os.Getenv("AZURE_HTTP_USER_AGENT"); azureAgent != "" {
		client.UserAgent = fmt.Sprintf("%s %s", client.UserAgent, azureAgent)
	}

	// only one pid can be interpreted currently
	// hence, send partner ID if present, otherwise send Terraform GUID
	// unless users have opted out
	if partnerID == "" && !disableTerraformPartnerID {
		// Microsoft’s Terraform Partner ID is this specific GUID
		partnerID = "222c6c49-1b0a-5959-a213-6608f9eb8820"
	}

	if partnerID != "" {
		client.UserAgent = fmt.Sprintf("%s pid-%s", client.UserAgent, partnerID)
	}

	// TODO -- refactor this to reuse the TelemetryValue function
}

func TelemetryValue(tfVersion, partnerID string, disableTerraformPartnerID bool) string {
	tfUserAgent := fmt.Sprintf("HashiCorp Terraform/%s (+https://www.terraform.io) Terraform Plugin SDK/%s", tfVersion, meta.SDKVersionString())

	providerUserAgent := fmt.Sprintf("%s terraform-provider-azurerm/%s", tfUserAgent, version.ProviderVersion)

	// append the CloudShell version to the user agent if it exists
	if azureAgent := os.Getenv("AZURE_HTTP_USER_AGENT"); azureAgent != "" {
		providerUserAgent = fmt.Sprintf("%s %s", providerUserAgent, azureAgent)
	}

	// only one pid can be interpreted currently
	// hence, send partner ID if present, otherwise send Terraform GUID
	// unless users have opted out
	if partnerID == "" && !disableTerraformPartnerID {
		// Microsoft’s Terraform Partner ID is this specific GUID
		partnerID = "222c6c49-1b0a-5959-a213-6608f9eb8820"
	}

	if partnerID != "" {
		providerUserAgent = fmt.Sprintf("%s pid-%s", providerUserAgent, partnerID)
	}

	return providerUserAgent
}
