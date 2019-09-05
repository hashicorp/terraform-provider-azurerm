package apimanagement

import (
	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2018-01-01/apimanagement"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ApiClient                  *apimanagement.APIClient
	ApiPoliciesClient          *apimanagement.APIPolicyClient
	ApiOperationsClient        *apimanagement.APIOperationClient
	ApiOperationPoliciesClient *apimanagement.APIOperationPolicyClient
	ApiSchemasClient           *apimanagement.APISchemaClient
	ApiVersionSetClient        *apimanagement.APIVersionSetClient
	AuthorizationServersClient *apimanagement.AuthorizationServerClient
	BackendClient              *apimanagement.BackendClient
	CertificatesClient         *apimanagement.CertificateClient
	GroupClient                *apimanagement.GroupClient
	GroupUsersClient           *apimanagement.GroupUserClient
	LoggerClient               *apimanagement.LoggerClient
	OpenIdConnectClient        *apimanagement.OpenIDConnectProviderClient
	PolicyClient               *apimanagement.PolicyClient
	ProductsClient             *apimanagement.ProductClient
	ProductApisClient          *apimanagement.ProductAPIClient
	ProductGroupsClient        *apimanagement.ProductGroupClient
	ProductPoliciesClient      *apimanagement.ProductPolicyClient
	PropertyClient             *apimanagement.PropertyClient
	ServiceClient              *apimanagement.ServiceClient
	SignInClient               *apimanagement.SignInSettingsClient
	SignUpClient               *apimanagement.SignUpSettingsClient
	SubscriptionsClient        *apimanagement.SubscriptionClient
	UsersClient                *apimanagement.UserClient
}

func BuildClient(o *common.ClientOptions) *Client {

	ApiClient := apimanagement.NewAPIClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ApiClient.Client, o.ResourceManagerAuthorizer)

	ApiPoliciesClient := apimanagement.NewAPIPolicyClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ApiPoliciesClient.Client, o.ResourceManagerAuthorizer)

	ApiOperationsClient := apimanagement.NewAPIOperationClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ApiOperationsClient.Client, o.ResourceManagerAuthorizer)

	ApiOperationPoliciesClient := apimanagement.NewAPIOperationPolicyClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ApiOperationPoliciesClient.Client, o.ResourceManagerAuthorizer)

	ApiSchemasClient := apimanagement.NewAPISchemaClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ApiSchemasClient.Client, o.ResourceManagerAuthorizer)

	ApiVersionSetClient := apimanagement.NewAPIVersionSetClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ApiVersionSetClient.Client, o.ResourceManagerAuthorizer)

	AuthorizationServersClient := apimanagement.NewAuthorizationServerClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&AuthorizationServersClient.Client, o.ResourceManagerAuthorizer)

	BackendClient := apimanagement.NewBackendClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&BackendClient.Client, o.ResourceManagerAuthorizer)

	CertificatesClient := apimanagement.NewCertificateClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&CertificatesClient.Client, o.ResourceManagerAuthorizer)

	GroupClient := apimanagement.NewGroupClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&GroupClient.Client, o.ResourceManagerAuthorizer)

	GroupUsersClient := apimanagement.NewGroupUserClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&GroupUsersClient.Client, o.ResourceManagerAuthorizer)

	LoggerClient := apimanagement.NewLoggerClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&LoggerClient.Client, o.ResourceManagerAuthorizer)

	PolicyClient := apimanagement.NewPolicyClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&PolicyClient.Client, o.ResourceManagerAuthorizer)

	ServiceClient := apimanagement.NewServiceClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ServiceClient.Client, o.ResourceManagerAuthorizer)

	SignInClient := apimanagement.NewSignInSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&SignInClient.Client, o.ResourceManagerAuthorizer)

	SignUpClient := apimanagement.NewSignUpSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&SignUpClient.Client, o.ResourceManagerAuthorizer)

	OpenIdConnectClient := apimanagement.NewOpenIDConnectProviderClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&OpenIdConnectClient.Client, o.ResourceManagerAuthorizer)

	ProductsClient := apimanagement.NewProductClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ProductsClient.Client, o.ResourceManagerAuthorizer)

	ProductApisClient := apimanagement.NewProductAPIClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ProductApisClient.Client, o.ResourceManagerAuthorizer)

	ProductGroupsClient := apimanagement.NewProductGroupClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ProductGroupsClient.Client, o.ResourceManagerAuthorizer)

	ProductPoliciesClient := apimanagement.NewProductPolicyClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ProductPoliciesClient.Client, o.ResourceManagerAuthorizer)

	PropertyClient := apimanagement.NewPropertyClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&PropertyClient.Client, o.ResourceManagerAuthorizer)

	SubscriptionsClient := apimanagement.NewSubscriptionClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&SubscriptionsClient.Client, o.ResourceManagerAuthorizer)

	UsersClient := apimanagement.NewUserClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&UsersClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ApiClient:                  &ApiClient,
		ApiPoliciesClient:          &ApiPoliciesClient,
		ApiOperationsClient:        &ApiOperationsClient,
		ApiOperationPoliciesClient: &ApiOperationPoliciesClient,
		ApiSchemasClient:           &ApiSchemasClient,
		ApiVersionSetClient:        &ApiVersionSetClient,
		AuthorizationServersClient: &AuthorizationServersClient,
		BackendClient:              &BackendClient,
		CertificatesClient:         &CertificatesClient,
		GroupClient:                &GroupClient,
		GroupUsersClient:           &GroupUsersClient,
		LoggerClient:               &LoggerClient,
		OpenIdConnectClient:        &OpenIdConnectClient,
		PolicyClient:               &PolicyClient,
		ProductsClient:             &ProductsClient,
		ProductApisClient:          &ProductApisClient,
		ProductGroupsClient:        &ProductGroupsClient,
		ProductPoliciesClient:      &ProductPoliciesClient,
		PropertyClient:             &PropertyClient,
		ServiceClient:              &ServiceClient,
		SignInClient:               &SignInClient,
		SignUpClient:               &SignUpClient,
		SubscriptionsClient:        &SubscriptionsClient,
		UsersClient:                &UsersClient,
	}
}
