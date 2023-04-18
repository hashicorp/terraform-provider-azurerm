package common

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/sender"
	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/version"
)

type Authorizers struct {
	BatchManagement auth.Authorizer
	KeyVault        auth.Authorizer
	ResourceManager auth.Authorizer
	Storage         auth.Authorizer
	Synapse         auth.Authorizer

	// Some data-plane APIs require a token scoped for a specific endpoint
	AuthorizerFunc ApiAuthorizerFunc
}

type ApiAuthorizerFunc func(api environments.Api) (auth.Authorizer, error)

type ClientOptions struct {
	Authorizers *Authorizers
	Environment environments.Environment
	Features    features.UserFeatures

	SubscriptionId   string
	TenantId         string
	PartnerId        string
	TerraformVersion string

	CustomCorrelationRequestID  string
	DisableCorrelationRequestID bool

	DisableTerraformPartnerID bool
	SkipProviderReg           bool
	StorageUseAzureAD         bool

	// Keep these around for convenience with Autorest based clients, remove when we are no longer using autorest
	AzureEnvironment        azure.Environment
	ResourceManagerEndpoint string

	// Legacy authorizers for go-autorest
	BatchManagementAuthorizer autorest.Authorizer
	KeyVaultAuthorizer        autorest.Authorizer
	ResourceManagerAuthorizer autorest.Authorizer
	StorageAuthorizer         autorest.Authorizer
	SynapseAuthorizer         autorest.Authorizer
}

func (o ClientOptions) Clone() ClientOptions {
	return ClientOptions{
		Authorizers:                 o.Authorizers,
		Environment:                 o.Environment,
		Features:                    o.Features,
		SubscriptionId:              o.SubscriptionId,
		TenantId:                    o.TenantId,
		PartnerId:                   o.PartnerId,
		TerraformVersion:            o.TerraformVersion,
		CustomCorrelationRequestID:  o.CustomCorrelationRequestID,
		DisableCorrelationRequestID: o.DisableCorrelationRequestID,
		DisableTerraformPartnerID:   o.DisableTerraformPartnerID,
		SkipProviderReg:             o.SkipProviderReg,
		StorageUseAzureAD:           o.StorageUseAzureAD,
		AzureEnvironment:            o.AzureEnvironment,
		ResourceManagerEndpoint:     o.ResourceManagerEndpoint,
		BatchManagementAuthorizer:   o.BatchManagementAuthorizer,
		KeyVaultAuthorizer:          o.KeyVaultAuthorizer,
		ResourceManagerAuthorizer:   o.ResourceManagerAuthorizer,
		StorageAuthorizer:           o.StorageAuthorizer,
		SynapseAuthorizer:           o.SynapseAuthorizer,
	}
}

// Configure set up a resourcemanager.Client using an auth.Authorizer from hashicorp/go-azure-sdk
func (o ClientOptions) Configure(c *resourcemanager.Client, authorizer auth.Authorizer) {
	c.Authorizer = authorizer
	c.UserAgent = userAgent(c.UserAgent, o.TerraformVersion, o.PartnerId, o.DisableTerraformPartnerID)

	requestMiddlewares := make([]client.RequestMiddleware, 0)
	if !o.DisableCorrelationRequestID {
		id := o.CustomCorrelationRequestID
		if id == "" {
			id = correlationRequestID()
		}
		requestMiddlewares = append(requestMiddlewares, correlationRequestIDMiddleware(id))
	}
	requestMiddlewares = append(requestMiddlewares, requestLoggerMiddleware("AzureRM"))
	c.RequestMiddlewares = &requestMiddlewares

	c.ResponseMiddlewares = &[]client.ResponseMiddleware{
		responseLoggerMiddleware("AzureRM"),
	}
}

// ConfigureClient sets up an autorest.Client using an autorest.Authorizer
func (o ClientOptions) ConfigureClient(c *autorest.Client, authorizer autorest.Authorizer) {
	// The default retry duration is 30 second, combined with the default retry attempts (default to 3) and the backoff logic,
	// makes the Read function always exceeds context deadline (since the provider uses 5min as the timeout value for read).
	// This brings users useless error message.
	// Changing it to 5 second will makes the total retries to be 75sec (= 5sec + 10sec + 20sec + 40sec), which is less than 5min,
	// so that users can always get the real error message instead.
	c.RetryDuration = 5 * time.Second
	c.UserAgent = userAgent(c.UserAgent, o.TerraformVersion, o.PartnerId, o.DisableTerraformPartnerID)

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

func userAgent(userAgent, tfVersion, partnerID string, disableTerraformPartnerID bool) string {
	tfUserAgent := fmt.Sprintf("HashiCorp Terraform/%s (+https://www.terraform.io) Terraform Plugin SDK/%s", tfVersion, meta.SDKVersionString())

	providerUserAgent := fmt.Sprintf("%s terraform-provider-azurerm/%s", tfUserAgent, version.ProviderVersion)
	if features.FourPointOhBeta() {
		providerUserAgent = fmt.Sprintf("%s terraform-provider-azurerm/%s+4.0-beta", tfUserAgent, version.ProviderVersion)
	}
	userAgent = strings.TrimSpace(fmt.Sprintf("%s %s", userAgent, providerUserAgent))

	// append the CloudShell version to the user agent if it exists
	if azureAgent := os.Getenv("AZURE_HTTP_USER_AGENT"); azureAgent != "" {
		userAgent = fmt.Sprintf("%s %s", userAgent, azureAgent)
	}

	// only one pid can be interpreted currently
	// hence, send partner ID if present, otherwise send Terraform GUID
	// unless users have opted out
	if partnerID == "" && !disableTerraformPartnerID {
		// Microsoft’s Terraform Partner ID is this specific GUID
		partnerID = "222c6c49-1b0a-5959-a213-6608f9eb8820"
	}

	if partnerID != "" {
		// Tolerate partnerID UUIDs without the "pid-" prefix
		userAgent = fmt.Sprintf("%s pid-%s", userAgent, strings.TrimPrefix(partnerID, "pid-"))
	}

	return userAgent
}
