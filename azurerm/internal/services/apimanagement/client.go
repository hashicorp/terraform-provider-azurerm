package apimanagement

import (
	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2018-01-01/apimanagement"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ApiClient                  apimanagement.APIClient
	ApiPoliciesClient          apimanagement.APIPolicyClient
	ApiOperationsClient        apimanagement.APIOperationClient
	ApiOperationPoliciesClient apimanagement.APIOperationPolicyClient
	ApiSchemasClient           apimanagement.APISchemaClient
	ApiVersionSetClient        apimanagement.APIVersionSetClient
	AuthorizationServersClient apimanagement.AuthorizationServerClient
	BackendClient              apimanagement.BackendClient
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

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.ApiClient = apimanagement.NewAPIClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ApiClient.Client, o.ResourceManagerAuthorizer)

	c.ApiPoliciesClient = apimanagement.NewAPIPolicyClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ApiPoliciesClient.Client, o.ResourceManagerAuthorizer)

	c.ApiOperationsClient = apimanagement.NewAPIOperationClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ApiOperationsClient.Client, o.ResourceManagerAuthorizer)

	c.ApiOperationPoliciesClient = apimanagement.NewAPIOperationPolicyClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ApiOperationPoliciesClient.Client, o.ResourceManagerAuthorizer)

	c.ApiSchemasClient = apimanagement.NewAPISchemaClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ApiSchemasClient.Client, o.ResourceManagerAuthorizer)

	c.ApiVersionSetClient = apimanagement.NewAPIVersionSetClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ApiVersionSetClient.Client, o.ResourceManagerAuthorizer)

	c.AuthorizationServersClient = apimanagement.NewAuthorizationServerClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.AuthorizationServersClient.Client, o.ResourceManagerAuthorizer)

	c.BackendClient = apimanagement.NewBackendClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.BackendClient.Client, o.ResourceManagerAuthorizer)

	c.CertificatesClient = apimanagement.NewCertificateClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.CertificatesClient.Client, o.ResourceManagerAuthorizer)

	c.GroupClient = apimanagement.NewGroupClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.GroupClient.Client, o.ResourceManagerAuthorizer)

	c.GroupUsersClient = apimanagement.NewGroupUserClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.GroupUsersClient.Client, o.ResourceManagerAuthorizer)

	c.LoggerClient = apimanagement.NewLoggerClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.LoggerClient.Client, o.ResourceManagerAuthorizer)

	c.PolicyClient = apimanagement.NewPolicyClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.PolicyClient.Client, o.ResourceManagerAuthorizer)

	c.ServiceClient = apimanagement.NewServiceClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ServiceClient.Client, o.ResourceManagerAuthorizer)

	c.SignInClient = apimanagement.NewSignInSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.SignInClient.Client, o.ResourceManagerAuthorizer)

	c.SignUpClient = apimanagement.NewSignUpSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.SignUpClient.Client, o.ResourceManagerAuthorizer)

	c.OpenIdConnectClient = apimanagement.NewOpenIDConnectProviderClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.OpenIdConnectClient.Client, o.ResourceManagerAuthorizer)

	c.ProductsClient = apimanagement.NewProductClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ProductsClient.Client, o.ResourceManagerAuthorizer)

	c.ProductApisClient = apimanagement.NewProductAPIClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ProductApisClient.Client, o.ResourceManagerAuthorizer)

	c.ProductGroupsClient = apimanagement.NewProductGroupClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ProductGroupsClient.Client, o.ResourceManagerAuthorizer)

	c.ProductPoliciesClient = apimanagement.NewProductPolicyClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ProductPoliciesClient.Client, o.ResourceManagerAuthorizer)

	c.PropertyClient = apimanagement.NewPropertyClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.PropertyClient.Client, o.ResourceManagerAuthorizer)

	c.SubscriptionsClient = apimanagement.NewSubscriptionClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.SubscriptionsClient.Client, o.ResourceManagerAuthorizer)

	c.UsersClient = apimanagement.NewUserClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.UsersClient.Client, o.ResourceManagerAuthorizer)

	return &c
}
