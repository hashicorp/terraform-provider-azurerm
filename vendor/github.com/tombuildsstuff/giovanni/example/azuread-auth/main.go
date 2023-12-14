package main

// TODO: update & re-enable this (see #75)
// The last stable example can be found here: https://github.com/tombuildsstuff/giovanni/tree/v0.20.0/example/azuread-auth

//import (
//	"context"
//	"fmt"
//	"log"
//	"net/http"
//	"os"
//
//	"github.com/Azure/go-autorest/autorest"
//	"github.com/hashicorp/go-azure-helpers/authentication"
//	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/blob/containers"
//)
//
//func main() {
//	log.Printf("[DEBUG] Started..")
//
//	// NOTE: fill this in
//	storageAccountName := "example"
//
//	log.Printf("[DEBUG] Building Client..")
//	client, err := buildClient()
//	if err != nil {
//		panic(fmt.Errorf("Error building client: %s", err))
//	}
//
//	ctx := context.TODO()
//	containerName := "armauth"
//	input := containers.CreateInput{
//		AccessLevel: containers.Private,
//		MetaData: map[string]string{
//			"hello": "world",
//		},
//	}
//	log.Printf("[DEBUG] Creating Container..")
//	if _, err := client.ContainersClient.Create(ctx, storageAccountName, containerName, input); err != nil {
//		panic(fmt.Errorf("Error creating container: %s", err))
//	}
//
//	log.Printf("[DEBUG] Retrieving Container..")
//	container, err := client.ContainersClient.GetProperties(ctx, storageAccountName, containerName)
//	if err != nil {
//		panic(fmt.Errorf("Error reading properties for container: %s", err))
//	}
//
//	log.Printf("[DEBUG] MetaData: %+v", container.MetaData)
//}
//
//type Client struct {
//	ContainersClient containers.Client
//}
//
//func buildClient() (*Client, error) {
//	// we're using github.com/hashicorp/go-azure-helpers since it makes this simpler
//	// but you can use an Authorizer from github.com/Azure/go-autorest directly too
//	builder := &authentication.Builder{
//		SubscriptionID: os.Getenv("ARM_SUBSCRIPTION_ID"),
//		ClientID:       os.Getenv("ARM_CLIENT_ID"),
//		ClientSecret:   os.Getenv("ARM_CLIENT_SECRET"),
//		TenantID:       os.Getenv("ARM_TENANT_ID"),
//		Environment:    os.Getenv("ARM_ENVIRONMENT"),
//
//		// Feature Toggles
//		SupportsClientSecretAuth: true,
//		SupportsAzureCliToken:    true,
//	}
//
//	config, err := builder.Build()
//	if err != nil {
//		return nil, fmt.Errorf("Error building AzureRM Client: %s", err)
//	}
//
//	env, err := authentication.DetermineEnvironment(config.Environment)
//	if err != nil {
//		return nil, err
//	}
//
//	oauthConfig, err := config.BuildOAuthConfig(env.ActiveDirectoryEndpoint)
//	if err != nil {
//		return nil, err
//	}
//
//	// OAuthConfigForTenant returns a pointer, which can be nil.
//	if oauthConfig == nil {
//		return nil, fmt.Errorf("Unable to configure OAuthConfig for tenant %s", config.TenantID)
//	}
//
//	// support for HTTP Proxies
//	sender := autorest.DecorateSender(&http.Client{
//		Transport: &http.Transport{
//			Proxy: http.ProxyFromEnvironment,
//		},
//	})
//
//	storageAuth, err := config.GetAuthorizationToken(sender, oauthConfig, "https://storage.azure.com/")
//	if err != nil {
//		return nil, err
//	}
//
//	containersClient := containers.New()
//	containersClient.Client.Authorizer = storageAuth
//
//	result := &Client{
//		ContainersClient: containersClient,
//	}
//
//	return result, nil
//}
