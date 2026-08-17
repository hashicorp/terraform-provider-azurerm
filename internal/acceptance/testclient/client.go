// Copyright IBM Corp. 2014, 2025
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
	_clients   = make(map[string]*clients.Client)
	clientLock = &sync.Mutex{}
)

func Build() (*clients.Client, error) {
	return BuildWithTestName(CurrentTestName())
}

func BuildWithTestName(testName string) (*clients.Client, error) {
	clientLock.Lock()
	defer clientLock.Unlock()

	c, ok := _clients[testName]
	if !ok {
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
			if env, err = environments.FromEndpoint(ctx, fmt.Sprintf("https://%s", metadataHost)); err != nil {
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
			AuthConfig:        &authConfig,
			TerraformVersion:  os.Getenv("TERRAFORM_CORE_VERSION"),
			Features:          features.Default(),
			StorageUseAzureAD: false,
			SubscriptionID:    os.Getenv("ARM_SUBSCRIPTION_ID"),
			TestName:          testName,
		}

		client, err := clients.Build(ctx, clientBuilder)
		if err != nil {
			return nil, fmt.Errorf("building test client: %+v", err)
		}

		c = client
		_clients[testName] = c
	}

	return c, nil
}
