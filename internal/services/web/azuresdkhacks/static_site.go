package azuresdkhacks

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web"
	"github.com/Azure/go-autorest/autorest"
)

// Copied from github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web/models.go, except the "StaticSite" is using the local type defined in this package.
// - the "StaticSite" is using the local type defined in this package.
// - remove some read only properties
type StaticSiteARMResource struct {
	autorest.Response `json:"-"`
	*StaticSite       `json:"properties,omitempty"`
	Sku               *web.SkuDescription         `json:"sku,omitempty"`
	Identity          *web.ManagedServiceIdentity `json:"identity,omitempty"`
	Location          *string                     `json:"location,omitempty"`
	Tags              map[string]*string          `json:"tags"`
	Kind              *string                     `json:"kind,omitempty"`
	// ID                *string                     `json:"id,omitempty"`
	// Name              *string                     `json:"name,omitempty"`
	//Type              *string                     `json:"type,omitempty"`
}

// Copied from github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web/models.go
func (ssar StaticSiteARMResource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if ssar.StaticSite != nil {
		objectMap["properties"] = ssar.StaticSite
	}
	if ssar.Sku != nil {
		objectMap["sku"] = ssar.Sku
	}
	if ssar.Identity != nil {
		objectMap["identity"] = ssar.Identity
	}
	if ssar.Kind != nil {
		objectMap["kind"] = ssar.Kind
	}
	if ssar.Location != nil {
		objectMap["location"] = ssar.Location
	}
	if ssar.Tags != nil {
		objectMap["tags"] = ssar.Tags
	}
	return json.Marshal(objectMap)
}

// Copied from github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web/models.go
type StaticSite struct {
	DefaultHostname             *string                                                       `json:"defaultHostname,omitempty"`
	RepositoryURL               *string                                                       `json:"repositoryUrl,omitempty"`
	Branch                      *string                                                       `json:"branch,omitempty"`
	CustomDomains               *[]string                                                     `json:"customDomains,omitempty"`
	RepositoryToken             *string                                                       `json:"repositoryToken,omitempty"`
	BuildProperties             *web.StaticSiteBuildProperties                                `json:"buildProperties,omitempty"`
	PrivateEndpointConnections  *[]web.ResponseMessageEnvelopeRemotePrivateEndpointConnection `json:"privateEndpointConnections,omitempty"`
	StagingEnvironmentPolicy    web.StagingEnvironmentPolicy                                  `json:"stagingEnvironmentPolicy,omitempty"`
	AllowConfigFileUpdates      *bool                                                         `json:"allowConfigFileUpdates,omitempty"`
	TemplateProperties          *web.StaticSiteTemplateOptions                                `json:"templateProperties,omitempty"`
	ContentDistributionEndpoint *string                                                       `json:"contentDistributionEndpoint,omitempty"`
	KeyVaultReferenceIdentity   *string                                                       `json:"keyVaultReferenceIdentity,omitempty"`
	UserProvidedFunctionApps    *[]web.StaticSiteUserProvidedFunctionApp                      `json:"userProvidedFunctionApps,omitempty"`
	Provider                    *string                                                       `json:"provider,omitempty"`
}

// Copied from github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web/models.go, except including the read-only properties.
func (ss StaticSite) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if ss.RepositoryURL != nil {
		objectMap["repositoryUrl"] = ss.RepositoryURL
	}
	if ss.Branch != nil {
		objectMap["branch"] = ss.Branch
	}
	if ss.RepositoryToken != nil {
		objectMap["repositoryToken"] = ss.RepositoryToken
	}
	if ss.BuildProperties != nil {
		objectMap["buildProperties"] = ss.BuildProperties
	}
	if ss.StagingEnvironmentPolicy != "" {
		objectMap["stagingEnvironmentPolicy"] = ss.StagingEnvironmentPolicy
	}
	if ss.AllowConfigFileUpdates != nil {
		objectMap["allowConfigFileUpdates"] = ss.AllowConfigFileUpdates
	}
	if ss.TemplateProperties != nil {
		objectMap["templateProperties"] = ss.TemplateProperties
	}

	// Changes from here
	if ss.DefaultHostname != nil {
		objectMap["defaultHostname"] = ss.DefaultHostname
	}
	if ss.CustomDomains != nil {
		objectMap["customDomains"] = ss.CustomDomains
	}
	if ss.PrivateEndpointConnections != nil {
		objectMap["privateEndpointConnections"] = ss.PrivateEndpointConnections
	}
	if ss.ContentDistributionEndpoint != nil {
		objectMap["contentDistributionEndpoint"] = ss.ContentDistributionEndpoint
	}
	if ss.KeyVaultReferenceIdentity != nil {
		objectMap["keyVaultReferenceIdentity"] = ss.KeyVaultReferenceIdentity
	}
	if ss.UserProvidedFunctionApps != nil {
		objectMap["userProvidedFunctionApps"] = ss.UserProvidedFunctionApps
	}
	if ss.Provider != nil {
		objectMap["provider"] = ss.Provider
	}
	return json.Marshal(objectMap)
}

func CreateOrUpdateStaticSite(ctx context.Context, client *web.StaticSitesClient, resourceGroupName string, name string, staticSiteEnvelope web.StaticSiteARMResource) (result web.StaticSitesCreateOrUpdateStaticSiteFuture, err error) {
	req, err := CreateOrUpdateStaticSitePreparer(ctx, client, resourceGroupName, name, staticSiteEnvelope)
	if err != nil {
		err = autorest.NewErrorWithError(err, "web.StaticSitesClient", "CreateOrUpdateStaticSite", nil, "Failure preparing request")
		return
	}

	result, err = client.CreateOrUpdateStaticSiteSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "web.StaticSitesClient", "CreateOrUpdateStaticSite", result.Response(), "Failure sending request")
		return
	}

	return
}

func CreateOrUpdateStaticSitePreparer(ctx context.Context, client *web.StaticSitesClient, resourceGroupName string, name string, staticSiteEnvelope web.StaticSiteARMResource) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"name":              autorest.Encode("path", name),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2021-02-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	// Transform the body parameter from the SDK type to the locally defined type
	params := StaticSiteARMResource{
		Sku:      staticSiteEnvelope.Sku,
		Identity: staticSiteEnvelope.Identity,
		Location: staticSiteEnvelope.Location,
		Tags:     staticSiteEnvelope.Tags,
		Kind:     staticSiteEnvelope.Kind,
	}
	if props := staticSiteEnvelope.StaticSite; props != nil {
		params.StaticSite = &StaticSite{
			DefaultHostname:             props.DefaultHostname,
			RepositoryURL:               props.RepositoryURL,
			Branch:                      props.Branch,
			CustomDomains:               props.CustomDomains,
			RepositoryToken:             props.RepositoryToken,
			BuildProperties:             props.BuildProperties,
			PrivateEndpointConnections:  props.PrivateEndpointConnections,
			StagingEnvironmentPolicy:    props.StagingEnvironmentPolicy,
			AllowConfigFileUpdates:      props.AllowConfigFileUpdates,
			TemplateProperties:          props.TemplateProperties,
			ContentDistributionEndpoint: props.ContentDistributionEndpoint,
			KeyVaultReferenceIdentity:   props.KeyVaultReferenceIdentity,
			UserProvidedFunctionApps:    props.UserProvidedFunctionApps,
			Provider:                    props.Provider,
		}
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Web/staticSites/{name}", pathParameters),
		autorest.WithJSON(params),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}
