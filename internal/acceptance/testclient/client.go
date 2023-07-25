// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package testclient

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
)

var (
	_client    *clients.Client
	clientLock = &sync.Mutex{}
)

func Build() (*clients.Client, error) {
	clientLock.Lock()
	defer clientLock.Unlock()

	if _client == nil {
		var (
			ctx = context.TODO()

			env *environments.Environment
			err error

			metadataHost = os.Getenv("ARM_METADATA_HOSTNAME")
		)

		envName, exists := os.LookupEnv("ARM_ENVIRONMENT")
		if !exists {
			envName = "public"
		}

		if metadataHost != "" {
			if env, err = environments.FromEndpoint(ctx, fmt.Sprintf("https://%s", metadataHost), envName); err != nil {
				return nil, fmt.Errorf("building test client: %+v", err)
			}
		} else if env, err = environments.FromName(envName); err != nil {
			return nil, fmt.Errorf("building test client: %+v", err)
		}

		authConfig := auth.Credentials{
			Environment: *env,
			ClientID:    os.Getenv("ARM_CLIENT_ID"),
			TenantID:    os.Getenv("ARM_TENANT_ID"),

			ClientCertificatePath:     os.Getenv("ARM_CLIENT_CERTIFICATE_PATH"),
			ClientCertificatePassword: os.Getenv("ARM_CLIENT_CERTIFICATE_PASSWORD"),
			ClientSecret:              os.Getenv("ARM_CLIENT_SECRET"),

			EnableAuthenticatingUsingClientCertificate: true,
			EnableAuthenticatingUsingClientSecret:      true,
			EnableAuthenticatingUsingAzureCLI:          false,
			EnableAuthenticatingUsingManagedIdentity:   false,
			EnableAuthenticationUsingOIDC:              false,
			EnableAuthenticationUsingGitHubOIDC:        false,
		}

		clientBuilder := clients.ClientBuilder{
			AuthConfig:               &authConfig,
			SkipProviderRegistration: true,
			TerraformVersion:         os.Getenv("TERRAFORM_CORE_VERSION"),
			Features:                 features.Default(),
			StorageUseAzureAD:        false,
			SubscriptionID:           os.Getenv("ARM_SUBSCRIPTION_ID"),
		}

		client, err := clients.Build(ctx, clientBuilder)
		if err != nil {
			return nil, fmt.Errorf("building test client: %+v", err)
		}

		_client = client
	}

	return _client, nil
}
