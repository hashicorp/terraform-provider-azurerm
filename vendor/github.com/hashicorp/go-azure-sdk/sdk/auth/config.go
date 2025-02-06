// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MPL-2.0 License. See NOTICE.txt in the project root for license information.

package auth

import (
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Credentials sets up NewAuthorizer to return an Authorizer based on the provided credentails.
type Credentials struct {
	// Specifies the national cloud environment to use
	Environment environments.Environment

	// AuxiliaryTenantIDs specifies the Auxiliary Tenant IDs for which to obtain tokens in a multi-tenant scenario.
	AuxiliaryTenantIDs []string
	// ClientID specifies the Client ID for the application used to authenticate the connection
	ClientID string
	// TenantID specifies the Azure Active Directory Tenant to connect to, which must be a valid UUID.
	TenantID string

	// EnableAuthenticatingUsingAzureCLI specifies whether Azure CLI authentication should be checked.
	EnableAuthenticatingUsingAzureCLI bool
	// AzureCliSubscriptionIDHint is the subscription to target when selecting an account with which to obtain an access token
	// Used to hint to Azure CLI which of its signed-in accounts it should select, based on apparent access to the subscription.
	AzureCliSubscriptionIDHint string

	// EnableAuthenticatingUsingClientCertificate specifies whether Client Certificate authentication should be checked.
	EnableAuthenticatingUsingClientCertificate bool
	// ClientCertificateData specifies the contents of a Client Certificate PKCS#12 bundle.
	ClientCertificateData []byte
	// ClientCertificatePath specifies the path to a Client Certificate PKCS#12 bundle (.pfx file)
	ClientCertificatePath string
	// ClientCertificatePassword specifies the encryption password to unlock a Client Certificate.
	ClientCertificatePassword string

	// EnableAuthenticatingUsingClientSecret specifies whether Client Secret authentication should be used.
	EnableAuthenticatingUsingClientSecret bool
	// ClientSecret specifies the Secret used authenticate using Client Secret authentication.
	ClientSecret string

	// EnableAuthenticatingUsingManagedIdentity specifies whether Managed Identity authentication should be checked.
	EnableAuthenticatingUsingManagedIdentity bool
	// CustomManagedIdentityEndpoint specifies a custom endpoint which should be used for Managed Identity.
	CustomManagedIdentityEndpoint string

	// Enables OIDC authentication (federated client credentials).
	EnableAuthenticationUsingOIDC bool
	// OIDCAssertionToken specifies the OIDC Assertion Token to authenticate using Client Credentials.
	OIDCAssertionToken string

	// EnableAuthenticationUsingGitHubOIDC specifies whether GitHub OIDC
	EnableAuthenticationUsingGitHubOIDC bool
	// GitHubOIDCTokenRequestURL specifies the URL for GitHub's OIDC provider
	GitHubOIDCTokenRequestURL string
	// GitHubOIDCTokenRequestToken specifies the bearer token for the request to GitHub's OIDC provider
	GitHubOIDCTokenRequestToken string
}
