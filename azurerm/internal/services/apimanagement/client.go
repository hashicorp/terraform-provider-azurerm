package apimanagement

import (
	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2018-01-01/apimanagement"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/ar"
)

type Client struct {
	ApiClient                  apimanagement.APIClient
	ApiPoliciesClient          apimanagement.APIPolicyClient
	ApiOperationsClient        apimanagement.APIOperationClient
	ApiOperationPoliciesClient apimanagement.APIOperationPolicyClient
	ApiSchemasClient           apimanagement.APISchemaClient
	ApiVersionSetClient        apimanagement.APIVersionSetClient
	AuthorizationServersClient apimanagement.AuthorizationServerClient
	CertificatesClient         apimanagement.CertificateClient
	GroupClient                apimanagement.GroupClient
	GroupUsersClient           apimanagement.GroupUserClient
	LoggerClient               apimanagement.LoggerClient
	OpenIdConnectClient        apimanagement.OpenIDConnectProviderClient
	PolicyClient               apimanagement.PolicyClient
	ProductsClient             apimanagement.ProductClient
	ProductApisClient          apimanagement.ProductAPIClient
	ProductGroupsClient        apimanagement.ProductGroupClient
	ProductPoliciesClient      apimanagement.ProductPolicyClient
	PropertyClient             apimanagement.PropertyClient
	ServiceClient              apimanagement.ServiceClient
	SignInClient               apimanagement.SignInSettingsClient
	SignUpClient               apimanagement.SignUpSettingsClient
	SubscriptionsClient        apimanagement.SubscriptionClient
	UsersClient                apimanagement.UserClient
}

func BuildClient(endpoint, subscriptionId string, o *ar.ClientOptions) *Client {
	c := Client{}

	c.ApiClient = apimanagement.NewAPIClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ApiClient.Client, o)

	c.ApiPoliciesClient = apimanagement.NewAPIPolicyClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ApiPoliciesClient.Client, o)

	c.ApiOperationsClient = apimanagement.NewAPIOperationClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ApiOperationsClient.Client, o)

	c.ApiOperationPoliciesClient = apimanagement.NewAPIOperationPolicyClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ApiOperationPoliciesClient.Client, o)

	c.ApiSchemasClient = apimanagement.NewAPISchemaClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ApiSchemasClient.Client, o)

	c.ApiVersionSetClient = apimanagement.NewAPIVersionSetClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ApiVersionSetClient.Client, o)

	c.AuthorizationServersClient = apimanagement.NewAuthorizationServerClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.AuthorizationServersClient.Client, o)

	c.CertificatesClient = apimanagement.NewCertificateClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.CertificatesClient.Client, o)

	c.GroupClient = apimanagement.NewGroupClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.GroupClient.Client, o)

	c.GroupUsersClient = apimanagement.NewGroupUserClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.GroupUsersClient.Client, o)

	c.LoggerClient = apimanagement.NewLoggerClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.LoggerClient.Client, o)

	c.PolicyClient = apimanagement.NewPolicyClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.PolicyClient.Client, o)

	c.ServiceClient = apimanagement.NewServiceClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ServiceClient.Client, o)

	c.SignInClient = apimanagement.NewSignInSettingsClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.SignInClient.Client, o)

	c.SignUpClient = apimanagement.NewSignUpSettingsClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.SignUpClient.Client, o)

	c.OpenIdConnectClient = apimanagement.NewOpenIDConnectProviderClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.OpenIdConnectClient.Client, o)

	c.ProductsClient = apimanagement.NewProductClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ProductsClient.Client, o)

	c.ProductApisClient = apimanagement.NewProductAPIClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ProductApisClient.Client, o)

	c.ProductGroupsClient = apimanagement.NewProductGroupClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ProductGroupsClient.Client, o)

	c.ProductPoliciesClient = apimanagement.NewProductPolicyClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ProductPoliciesClient.Client, o)

	c.PropertyClient = apimanagement.NewPropertyClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.PropertyClient.Client, o)

	c.SubscriptionsClient = apimanagement.NewSubscriptionClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.SubscriptionsClient.Client, o)

	c.UsersClient = apimanagement.NewUserClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.UsersClient.Client, o)

	return &c
}
