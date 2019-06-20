package apimanagement

import (
	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2018-01-01/apimanagement"
	"github.com/Azure/go-autorest/autorest"
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

func BuildClient(endpoint string, authorizer autorest.Authorizer, o *common.ClientOptions) *Client {
	c := Client{}

	c.ApiClient = apimanagement.NewAPIClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ApiClient.Client, authorizer)

	c.ApiPoliciesClient = apimanagement.NewAPIPolicyClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ApiPoliciesClient.Client, authorizer)

	c.ApiOperationsClient = apimanagement.NewAPIOperationClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ApiOperationsClient.Client, authorizer)

	c.ApiOperationPoliciesClient = apimanagement.NewAPIOperationPolicyClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ApiOperationPoliciesClient.Client, authorizer)

	c.ApiSchemasClient = apimanagement.NewAPISchemaClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ApiSchemasClient.Client, authorizer)

	c.ApiVersionSetClient = apimanagement.NewAPIVersionSetClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ApiVersionSetClient.Client, authorizer)

	c.AuthorizationServersClient = apimanagement.NewAuthorizationServerClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.AuthorizationServersClient.Client, authorizer)

	c.CertificatesClient = apimanagement.NewCertificateClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.CertificatesClient.Client, authorizer)

	c.GroupClient = apimanagement.NewGroupClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.GroupClient.Client, authorizer)

	c.GroupUsersClient = apimanagement.NewGroupUserClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.GroupUsersClient.Client, authorizer)

	c.LoggerClient = apimanagement.NewLoggerClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.LoggerClient.Client, authorizer)

	c.PolicyClient = apimanagement.NewPolicyClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.PolicyClient.Client, authorizer)

	c.ServiceClient = apimanagement.NewServiceClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ServiceClient.Client, authorizer)

	c.SignInClient = apimanagement.NewSignInSettingsClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.SignInClient.Client, authorizer)

	c.SignUpClient = apimanagement.NewSignUpSettingsClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.SignUpClient.Client, authorizer)

	c.OpenIdConnectClient = apimanagement.NewOpenIDConnectProviderClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.OpenIdConnectClient.Client, authorizer)

	c.ProductsClient = apimanagement.NewProductClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ProductsClient.Client, authorizer)

	c.ProductApisClient = apimanagement.NewProductAPIClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ProductApisClient.Client, authorizer)

	c.ProductGroupsClient = apimanagement.NewProductGroupClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ProductGroupsClient.Client, authorizer)

	c.ProductPoliciesClient = apimanagement.NewProductPolicyClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ProductPoliciesClient.Client, authorizer)

	c.PropertyClient = apimanagement.NewPropertyClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.PropertyClient.Client, authorizer)

	c.SubscriptionsClient = apimanagement.NewSubscriptionClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.SubscriptionsClient.Client, authorizer)

	c.UsersClient = apimanagement.NewUserClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.UsersClient.Client, authorizer)

	return &c
}
