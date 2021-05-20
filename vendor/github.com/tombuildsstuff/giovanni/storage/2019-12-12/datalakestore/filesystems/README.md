## Data Lake Storage Gen2 File Systems SDK for API version 2019-12-12

This package allows you to interact with the Data Lake Storage Gen2 File Systems API

### Supported Authorizers

* Azure Active Directory (for the Resource Endpoint `https://storage.azure.com`)

### Example Usage

```go
package main

import (
	"context"
	"fmt"
	"os"
    "time"
	
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
    "github.com/Azure/go-autorest/autorest/azure"
    "github.com/hashicorp/go-azure-helpers/authentication"
    "github.com/hashicorp/go-azure-helpers/sender"
    "github.com/tombuildsstuff/giovanni/storage/2019-12-12/datalakestore/filesystems"
)

func Example() error {
	accountName := "storageaccount1"
    fileSystemName := "filesystem1"

    builder := &authentication.Builder{
        SubscriptionID: os.Getenv("ARM_SUBSCRIPTION_ID"),
        ClientID:       os.Getenv("ARM_CLIENT_ID"),
        ClientSecret:   os.Getenv("ARM_CLIENT_SECRET"),
        TenantID:       os.Getenv("ARM_TENANT_ID"),
        Environment:    os.Getenv("ARM_ENVIRONMENT"),

        // Feature Toggles
        SupportsClientSecretAuth: true,
    }

    c, err := builder.Build()
    if err != nil {
        return fmt.Errorf("Error building AzureRM Client: %s", err)
    }

    env, err := authentication.DetermineEnvironment(c.Environment)
    if err != nil {
        return err
    }

    oauthConfig, err := adal.NewOAuthConfig(env.ActiveDirectoryEndpoint, c.TenantID)
	if err != nil {
		return err
	}

	// OAuthConfigForTenant returns a pointer, which can be nil.
	if oauthConfig == nil {
		return fmt.Errorf("Unable to configure OAuthConfig for tenant %s", c.TenantID)
	}

    sender := sender.BuildSender("AzureRM")
    ctx := context.Background()

    storageAuth, err := config.GetAuthorizationToken(sender, oauthConfig, "https://storage.azure.com/")
	if err != nil {
		return fmt.Errorf("Error retrieving Authorization Token")
	}

   
    fileSystemsClient := filesystems.NewWithEnvironment(env)
	fileSystemsClient.Client.Authorizer = storageAuth

	input := filesystems.CreateInput{
		Properties: map[string]string{},
	}
	if _, err = fileSystemsClient.Create(ctx, accountName, fileSystemName, input); err != nil {
		return fmt.Errorf("Error creating: %s", err)
	}
	
    return nil 
}
```