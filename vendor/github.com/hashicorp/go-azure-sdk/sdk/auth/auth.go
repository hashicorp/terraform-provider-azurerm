// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MPL-2.0 License. See NOTICE.txt in the project root for license information.

package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// NewAuthorizerFromCredentials returns a suitable Authorizer depending on what is defined in the Credentials
// Authorizers are selected for authentication methods in the following preferential order:
// - Client certificate authentication
// - Client secret authentication
// - OIDC authentication
// - GitHub OIDC authentication
// - MSI authentication
// - Azure CLI authentication
//
// Whether one of these is returned depends on whether it is enabled in the Credentials, and whether sufficient
// configuration fields are set to enable that authentication method.
//
// For client certificate authentication, specify TenantID, ClientID and ClientCertificateData / ClientCertificatePath.
// For client secret authentication, specify TenantID, ClientID and ClientSecret.
// For OIDC authentication, specify TenantID, ClientID and OIDCAssertionToken.
// For GitHub OIDC authentication, specify TenantID, ClientID, GitHubOIDCTokenRequestURL and GitHubOIDCTokenRequestToken.
// MSI authentication (if enabled) using the Azure Metadata Service is then attempted
// Azure CLI authentication (if enabled) is attempted last
//
// It's recommended to only enable the mechanisms you have configured and are known to work in the execution
// environment. If any authentication mechanism fails due to misconfiguration or some other error, the function
// will return (nil, error) and later mechanisms will not be attempted.
func NewAuthorizerFromCredentials(ctx context.Context, c Credentials, api environments.Api) (Authorizer, error) {
	if c.EnableAuthenticatingUsingClientCertificate && strings.TrimSpace(c.TenantID) != "" && strings.TrimSpace(c.ClientID) != "" && (len(c.ClientCertificateData) > 0 || strings.TrimSpace(c.ClientCertificatePath) != "") {
		opts := ClientCertificateAuthorizerOptions{
			Environment:  c.Environment,
			Api:          api,
			TenantId:     c.TenantID,
			AuxTenantIds: c.AuxiliaryTenantIDs,
			ClientId:     c.ClientID,
			Pkcs12Data:   c.ClientCertificateData,
			Pkcs12Path:   c.ClientCertificatePath,
			Pkcs12Pass:   c.ClientCertificatePassword,
		}
		a, err := NewClientCertificateAuthorizer(ctx, opts)
		if err != nil {
			return nil, fmt.Errorf("could not configure ClientCertificate Authorizer: %s", err)
		}
		if a != nil {
			return a, nil
		}
	}

	if c.EnableAuthenticatingUsingClientSecret && strings.TrimSpace(c.TenantID) != "" && strings.TrimSpace(c.ClientID) != "" && strings.TrimSpace(c.ClientSecret) != "" {
		opts := ClientSecretAuthorizerOptions{
			Environment:  c.Environment,
			Api:          api,
			TenantId:     c.TenantID,
			AuxTenantIds: c.AuxiliaryTenantIDs,
			ClientId:     c.ClientID,
			ClientSecret: c.ClientSecret,
		}
		a, err := NewClientSecretAuthorizer(ctx, opts)
		if err != nil {
			return nil, fmt.Errorf("could not configure ClientSecret Authorizer: %s", err)
		}
		if a != nil {
			return a, nil
		}
	}

	if c.EnableAuthenticationUsingOIDC && strings.TrimSpace(c.TenantID) != "" && strings.TrimSpace(c.ClientID) != "" && strings.TrimSpace(c.OIDCAssertionToken) != "" {
		opts := OIDCAuthorizerOptions{
			Environment:        c.Environment,
			Api:                api,
			TenantId:           c.TenantID,
			AuxiliaryTenantIds: c.AuxiliaryTenantIDs,
			ClientId:           c.ClientID,
			FederatedAssertion: c.OIDCAssertionToken,
		}
		a, err := NewOIDCAuthorizer(ctx, opts)
		if err != nil {
			return nil, fmt.Errorf("could not configure OIDC Authorizer: %s", err)
		}
		if a != nil {
			return a, nil
		}
	}

	if c.EnableAuthenticationUsingGitHubOIDC && strings.TrimSpace(c.TenantID) != "" && strings.TrimSpace(c.ClientID) != "" && strings.TrimSpace(c.GitHubOIDCTokenRequestURL) != "" && strings.TrimSpace(c.GitHubOIDCTokenRequestToken) != "" {
		opts := GitHubOIDCAuthorizerOptions{
			Api:                 api,
			AuxiliaryTenantIds:  c.AuxiliaryTenantIDs,
			ClientId:            c.ClientID,
			Environment:         c.Environment,
			IdTokenRequestUrl:   c.GitHubOIDCTokenRequestURL,
			IdTokenRequestToken: c.GitHubOIDCTokenRequestToken,
			TenantId:            c.TenantID,
		}
		a, err := NewGitHubOIDCAuthorizer(context.Background(), opts)
		if err != nil {
			return nil, fmt.Errorf("could not configure GitHubOIDC Authorizer: %s", err)
		}
		if a != nil {
			return a, nil
		}
	}

	if c.EnableAuthenticatingUsingManagedIdentity {
		opts := ManagedIdentityAuthorizerOptions{
			Api:                           api,
			ClientId:                      c.ClientID,
			CustomManagedIdentityEndpoint: c.CustomManagedIdentityEndpoint,
		}
		a, err := NewManagedIdentityAuthorizer(ctx, opts)
		if err != nil {
			return nil, fmt.Errorf("could not configure MSI Authorizer: %s", err)
		}
		if a != nil {
			return a, nil
		}
	}

	if c.EnableAuthenticatingUsingAzureCLI {
		opts := AzureCliAuthorizerOptions{
			Api:                api,
			TenantId:           c.TenantID,
			AuxTenantIds:       c.AuxiliaryTenantIDs,
			SubscriptionIdHint: c.AzureCliSubscriptionIDHint,
		}
		a, err := NewAzureCliAuthorizer(ctx, opts)
		if err != nil {
			return nil, fmt.Errorf("could not configure AzureCli Authorizer: %s", err)
		}
		if a != nil {
			return a, nil
		}
	}

	return nil, fmt.Errorf("no Authorizer could be configured, please check your configuration")
}
