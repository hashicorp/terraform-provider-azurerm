package testhelpers

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane/storage"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	Environment          environments.Environment
	ResourceGroupsClient *resourcegroups.ResourceGroupsClient
	StorageAccountClient *storageaccounts.StorageAccountsClient
	SubscriptionId       string

	resourceManagerAuth auth.Authorizer
	storageAuth         auth.Authorizer
}

type TestResources struct {
	ResourceGroup      string
	StorageAccountName string
	StorageAccountKey  string
}

func (c Client) BuildTestResources(ctx context.Context, resourceGroup, name string, kind storageaccounts.Kind) (*TestResources, error) {
	return c.buildTestResources(ctx, resourceGroup, name, kind, false, "")
}
func (c Client) BuildTestResourcesWithHns(ctx context.Context, resourceGroup, name string, kind storageaccounts.Kind) (*TestResources, error) {
	return c.buildTestResources(ctx, resourceGroup, name, kind, true, "")
}
func (c Client) BuildTestResourcesWithSku(ctx context.Context, resourceGroup, name string, kind storageaccounts.Kind, sku storageaccounts.SkuName) (*TestResources, error) {
	return c.buildTestResources(ctx, resourceGroup, name, kind, false, sku)
}
func (c Client) buildTestResources(ctx context.Context, resourceGroup, name string, kind storageaccounts.Kind, enableHns bool, sku storageaccounts.SkuName) (*TestResources, error) {
	location := os.Getenv("ARM_TEST_LOCATION")
	resourceGroupId := commonids.NewResourceGroupID(c.SubscriptionId, resourceGroup)
	resourceGroupPayload := resourcegroups.ResourceGroup{
		Location: location,
	}
	if _, err := c.ResourceGroupsClient.CreateOrUpdate(ctx, resourceGroupId, resourceGroupPayload); err != nil {
		return nil, fmt.Errorf("error creating Resource Group %q: %s", resourceGroup, err)
	}

	props := storageaccounts.StorageAccountPropertiesCreateParameters{
		AllowBlobPublicAccess: pointer.To(true),
		PublicNetworkAccess:   pointer.To(storageaccounts.PublicNetworkAccessEnabled),
	}
	if kind == storageaccounts.KindBlobStorage {
		props.AccessTier = pointer.To(storageaccounts.AccessTierHot)
	}
	if enableHns {
		props.IsHnsEnabled = &enableHns
	}
	if sku == "" {
		sku = storageaccounts.SkuNameStandardLRS
	}

	payload := storageaccounts.StorageAccountCreateParameters{
		Location: location,
		Sku: storageaccounts.Sku{
			Name: sku,
		},
		Kind:       kind,
		Properties: &props,
	}
	storageAccountId := commonids.NewStorageAccountID(c.SubscriptionId, resourceGroup, name)
	if err := c.StorageAccountClient.CreateThenPoll(ctx, storageAccountId, payload); err != nil {
		return nil, fmt.Errorf("error creating %s: %+v", storageAccountId, err)
	}

	var options storageaccounts.ListKeysOperationOptions
	keys, err := c.StorageAccountClient.ListKeys(ctx, storageAccountId, options)
	if err != nil {
		return nil, fmt.Errorf("error listing keys for %s: %+v", storageAccountId, err)
	}

	// sure we could poll to get around the inconsistency, but where's the fun in that
	time.Sleep(5 * time.Second)

	accountKeys := *keys.Model.Keys
	return &TestResources{
		ResourceGroup:      resourceGroup,
		StorageAccountName: name,
		StorageAccountKey:  *(accountKeys[0]).Value,
	}, nil
}

func (c Client) DestroyTestResources(ctx context.Context, resourceGroup, name string) error {
	storageAccountId := commonids.NewStorageAccountID(c.SubscriptionId, resourceGroup, name)
	if _, err := c.StorageAccountClient.Delete(ctx, storageAccountId); err != nil {
		return fmt.Errorf("error deleting %s: %+v", storageAccountId, err)
	}

	resourceGroupId := commonids.NewResourceGroupID(c.SubscriptionId, resourceGroup)
	if err := c.ResourceGroupsClient.DeleteThenPoll(ctx, resourceGroupId, resourcegroups.DefaultDeleteOperationOptions()); err != nil {
		return fmt.Errorf("error deleting Resource Group %q: %s", resourceGroup, err)
	}

	return nil
}

func Build(ctx context.Context, t *testing.T) (*Client, error) {
	if os.Getenv("ACCTEST") == "" {
		t.Skip("Skipping as `ACCTEST` hasn't been set")
	}

	environmentName := os.Getenv("ARM_ENVIRONMENT")
	env, err := environments.FromName(environmentName)
	if err != nil {
		return nil, fmt.Errorf("determining environment %q: %+v", environmentName, err)
	}
	if env == nil {
		return nil, fmt.Errorf("environment was nil: %s", err)
	}

	authConfig := auth.Credentials{
		Environment:  *env,
		ClientID:     os.Getenv("ARM_CLIENT_ID"),
		TenantID:     os.Getenv("ARM_TENANT_ID"),
		ClientSecret: os.Getenv("ARM_CLIENT_SECRET"),

		EnableAuthenticatingUsingClientCertificate: true,
		EnableAuthenticatingUsingClientSecret:      true,
		EnableAuthenticatingUsingAzureCLI:          false,
		EnableAuthenticatingUsingManagedIdentity:   false,
		EnableAuthenticationUsingOIDC:              false,
		EnableAuthenticationUsingGitHubOIDC:        false,
	}

	resourceManagerAuth, err := auth.NewAuthorizerFromCredentials(ctx, authConfig, authConfig.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("unable to build authorizer for Resource Manager API: %+v", err)
	}

	storageAuthorizer, err := auth.NewAuthorizerFromCredentials(ctx, authConfig, authConfig.Environment.Storage)
	if err != nil {
		return nil, fmt.Errorf("unable to build authorizer for Storage API: %+v", err)
	}

	client := Client{
		Environment:    *env,
		SubscriptionId: os.Getenv("ARM_SUBSCRIPTION_ID"),

		// internal
		resourceManagerAuth: resourceManagerAuth,
		storageAuth:         storageAuthorizer,
	}

	resourceGroupsClient, err := resourcegroups.NewResourceGroupsClientWithBaseURI(env.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Resource Groups client: %+v", err)
	}
	client.Configure(resourceGroupsClient.Client.Client, client.resourceManagerAuth)
	client.ResourceGroupsClient = resourceGroupsClient

	storageClient, err := storageaccounts.NewStorageAccountsClientWithBaseURI(env.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Storage Accounts client: %+v", err)
	}
	client.Configure(storageClient.Client.Client, client.resourceManagerAuth)
	client.StorageAccountClient = storageClient

	return &client, nil
}

func (c Client) Configure(client *client.Client, authorizer auth.Authorizer) {
	client.Authorizer = authorizer
	// TODO: add logging
}

func (c Client) PrepareWithResourceManagerAuth(input *storage.BaseClient) {
	input.WithAuthorizer(c.storageAuth)
}

func (c Client) PrepareWithSharedKeyAuth(input *storage.BaseClient, data *TestResources, keyType auth.SharedKeyType) error {
	auth, err := auth.NewSharedKeyAuthorizer(data.StorageAccountName, data.StorageAccountKey, keyType)
	if err != nil {
		return fmt.Errorf("building SharedKey authorizer: %+v", err)
	}
	input.WithAuthorizer(auth)
	return nil
}
