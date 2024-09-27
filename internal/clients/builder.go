// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package clients

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	authWrapper "github.com/hashicorp/go-azure-sdk/sdk/auth/autorest"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceproviders"
)

type ClientBuilder struct {
	AuthConfig *auth.Credentials
	Features   features.UserFeatures

	CustomCorrelationRequestID  string
	DisableCorrelationRequestID bool
	DisableTerraformPartnerID   bool
	MetadataHost                string
	PartnerID                   string
	RegisteredResourceProviders resourceproviders.ResourceProviders
	StorageUseAzureAD           bool
	SubscriptionID              string
	TerraformVersion            string
}

const azureStackEnvironmentError = `
The AzureRM Provider supports the different Azure Public Clouds - including China, Public,
and US Government - however it does not support Azure Stack due to differences in API and
feature availability.

Terraform instead offers a separate "azurestack" provider which supports the functionality
and APIs available in Azure Stack via Azure Stack Profiles.
`

func Build(ctx context.Context, builder ClientBuilder) (*Client, error) {
	var err error

	// point folks towards the separate Azure Stack Provider when using Azure Stack
	if builder.AuthConfig.Environment.IsAzureStack() {
		return nil, fmt.Errorf(azureStackEnvironmentError)
	}

	var resourceManagerAuth, storageAuth, synapseAuth, batchManagementAuth, keyVaultAuth auth.Authorizer

	resourceManagerAuth, err = auth.NewAuthorizerFromCredentials(ctx, *builder.AuthConfig, builder.AuthConfig.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("unable to build authorizer for Resource Manager API: %+v", err)
	}

	storageAuth, err = auth.NewAuthorizerFromCredentials(ctx, *builder.AuthConfig, builder.AuthConfig.Environment.Storage)
	if err != nil {
		return nil, fmt.Errorf("unable to build authorizer for Storage API: %+v", err)
	}

	keyVaultAuth, err = auth.NewAuthorizerFromCredentials(ctx, *builder.AuthConfig, builder.AuthConfig.Environment.KeyVault)
	if err != nil {
		return nil, fmt.Errorf("unable to build authorizer for Key Vault API: %+v", err)
	}

	if builder.AuthConfig.Environment.Synapse.Available() {
		synapseAuth, err = auth.NewAuthorizerFromCredentials(ctx, *builder.AuthConfig, builder.AuthConfig.Environment.Synapse)
		if err != nil {
			return nil, fmt.Errorf("unable to build authorizer for Synapse API: %+v", err)
		}
	} else {
		log.Printf("[DEBUG] Skipping building the Synapse Authorizer since this is not supported in the current Azure Environment")
	}

	if builder.AuthConfig.Environment.Batch.Available() {
		batchManagementAuth, err = auth.NewAuthorizerFromCredentials(ctx, *builder.AuthConfig, builder.AuthConfig.Environment.Batch)
		if err != nil {
			return nil, fmt.Errorf("unable to build authorizer for Batch Management API: %+v", err)
		}
	} else {
		log.Printf("[DEBUG] Skipping building the Batch Management Authorizer since this is not supported in the current Azure Environment")
	}

	// Helper for obtaining endpoint-specific tokens
	authorizerFunc := common.ApiAuthorizerFunc(func(api environments.Api) (auth.Authorizer, error) {
		authorizer, err := auth.NewAuthorizerFromCredentials(ctx, *builder.AuthConfig, api)
		if err != nil {
			return nil, fmt.Errorf("building custom authorizer for API %q: %+v", api.Name(), err)
		}

		return authorizer, nil
	})

	account, err := NewResourceManagerAccount(ctx, *builder.AuthConfig, builder.SubscriptionID, builder.RegisteredResourceProviders)
	if err != nil {
		return nil, fmt.Errorf("building account: %+v", err)
	}

	var managedHSMAuth auth.Authorizer
	if builder.AuthConfig.Environment.ManagedHSM.Available() {
		managedHSMAuth, err = auth.NewAuthorizerFromCredentials(ctx, *builder.AuthConfig, builder.AuthConfig.Environment.ManagedHSM)
		if err != nil {
			return nil, fmt.Errorf("unable to build authorizer for Managed HSM API: %+v", err)
		}
	} else {
		log.Printf("[DEBUG] Skipping building the Managed HSM Authorizer since this is not supported in the current Azure Environment")
	}

	resourceManagerEndpoint, ok := builder.AuthConfig.Environment.ResourceManager.Endpoint()
	if !ok {
		return nil, fmt.Errorf("unable to determine resource manager endpoint for the current environment")
	}

	client := Client{
		Account: account,
	}

	o := &common.ClientOptions{
		Authorizers: &common.Authorizers{
			BatchManagement: batchManagementAuth,
			KeyVault:        keyVaultAuth,
			ManagedHSM:      managedHSMAuth,
			ResourceManager: resourceManagerAuth,
			Storage:         storageAuth,
			Synapse:         synapseAuth,
			AuthorizerFunc:  authorizerFunc,
		},

		AuthConfig:  builder.AuthConfig,
		Environment: builder.AuthConfig.Environment,
		Features:    builder.Features,

		SubscriptionId:   account.SubscriptionId,
		TenantId:         account.TenantId,
		PartnerId:        builder.PartnerID,
		TerraformVersion: builder.TerraformVersion,

		BatchManagementAuthorizer: authWrapper.AutorestAuthorizer(batchManagementAuth),
		KeyVaultAuthorizer:        authWrapper.AutorestAuthorizer(keyVaultAuth).BearerAuthorizerCallback(),
		ManagedHSMAuthorizer:      authWrapper.AutorestAuthorizer(managedHSMAuth).BearerAuthorizerCallback(),
		ResourceManagerAuthorizer: authWrapper.AutorestAuthorizer(resourceManagerAuth),
		SynapseAuthorizer:         authWrapper.AutorestAuthorizer(synapseAuth),

		CustomCorrelationRequestID:  builder.CustomCorrelationRequestID,
		DisableCorrelationRequestID: builder.DisableCorrelationRequestID,
		DisableTerraformPartnerID:   builder.DisableTerraformPartnerID,
		SkipProviderReg:             len(builder.RegisteredResourceProviders) == 0,
		StorageUseAzureAD:           builder.StorageUseAzureAD,

		ResourceManagerEndpoint: *resourceManagerEndpoint,
	}

	if err := client.Build(ctx, o); err != nil {
		return nil, fmt.Errorf("building Client: %+v", err)
	}

	if features.EnhancedValidationEnabled() {
		subscriptionId := commonids.NewSubscriptionID(client.Account.SubscriptionId)

		ctx2, cancel := context.WithTimeout(ctx, 10*time.Minute)
		defer cancel()

		location.CacheSupportedLocations(ctx2, *resourceManagerEndpoint)
		if err := resourceproviders.CacheSupportedProviders(ctx2, client.Resource.ResourceProvidersClient, subscriptionId); err != nil {
			log.Printf("[DEBUG] error retrieving providers: %s. Enhanced validation will be unavailable", err)
		}
	}

	return &client, nil
}
