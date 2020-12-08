package client

import (
	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2019-12-01/apimanagement"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ApiClient                  *apimanagement.APIClient
	ApiDiagnosticClient        *apimanagement.APIDiagnosticClient
	ApiPoliciesClient          *apimanagement.APIPolicyClient
	ApiOperationsClient        *apimanagement.APIOperationClient
	ApiOperationPoliciesClient *apimanagement.APIOperationPolicyClient
	ApiSchemasClient           *apimanagement.APISchemaClient
	ApiVersionSetClient        *apimanagement.APIVersionSetClient
	AuthorizationServersClient *apimanagement.AuthorizationServerClient
	BackendClient              *apimanagement.BackendClient
	CertificatesClient         *apimanagement.CertificateClient
	DiagnosticClient           *apimanagement.DiagnosticClient
	GroupClient                *apimanagement.GroupClient
	GroupUsersClient           *apimanagement.GroupUserClient
	IdentityProviderClient     *apimanagement.IdentityProviderClient
	LoggerClient               *apimanagement.LoggerClient
	NamedValueClient           *apimanagement.NamedValueClient
	OpenIdConnectClient        *apimanagement.OpenIDConnectProviderClient
	PolicyClient               *apimanagement.PolicyClient
	ProductsClient             *apimanagement.ProductClient
	ProductApisClient          *apimanagement.ProductAPIClient
	ProductGroupsClient        *apimanagement.ProductGroupClient
	ProductPoliciesClient      *apimanagement.ProductPolicyClient
	ServiceClient              *apimanagement.ServiceClient
	SignInClient               *apimanagement.SignInSettingsClient
	SignUpClient               *apimanagement.SignUpSettingsClient
	SubscriptionsClient        *apimanagement.SubscriptionClient
	UsersClient                *apimanagement.UserClient
}

func NewClient(o *common.ClientOptions) *Client {
	apiClient := apimanagement.NewAPIClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&apiClient.Client, o.ResourceManagerAuthorizer)

	apiDiagnosticClient := apimanagement.NewAPIDiagnosticClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&apiDiagnosticClient.Client, o.ResourceManagerAuthorizer)

	apiPoliciesClient := apimanagement.NewAPIPolicyClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&apiPoliciesClient.Client, o.ResourceManagerAuthorizer)

	apiOperationsClient := apimanagement.NewAPIOperationClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&apiOperationsClient.Client, o.ResourceManagerAuthorizer)

	apiOperationPoliciesClient := apimanagement.NewAPIOperationPolicyClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&apiOperationPoliciesClient.Client, o.ResourceManagerAuthorizer)

	apiSchemasClient := apimanagement.NewAPISchemaClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&apiSchemasClient.Client, o.ResourceManagerAuthorizer)

	apiVersionSetClient := apimanagement.NewAPIVersionSetClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&apiVersionSetClient.Client, o.ResourceManagerAuthorizer)

	authorizationServersClient := apimanagement.NewAuthorizationServerClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&authorizationServersClient.Client, o.ResourceManagerAuthorizer)

	backendClient := apimanagement.NewBackendClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&backendClient.Client, o.ResourceManagerAuthorizer)

	certificatesClient := apimanagement.NewCertificateClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&certificatesClient.Client, o.ResourceManagerAuthorizer)

	diagnosticClient := apimanagement.NewDiagnosticClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&diagnosticClient.Client, o.ResourceManagerAuthorizer)

	groupClient := apimanagement.NewGroupClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&groupClient.Client, o.ResourceManagerAuthorizer)

	groupUsersClient := apimanagement.NewGroupUserClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&groupUsersClient.Client, o.ResourceManagerAuthorizer)

	identityProviderClient := apimanagement.NewIdentityProviderClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&identityProviderClient.Client, o.ResourceManagerAuthorizer)

	namedValueClient := apimanagement.NewNamedValueClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&namedValueClient.Client, o.ResourceManagerAuthorizer)

	loggerClient := apimanagement.NewLoggerClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&loggerClient.Client, o.ResourceManagerAuthorizer)

	openIdConnectClient := apimanagement.NewOpenIDConnectProviderClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&openIdConnectClient.Client, o.ResourceManagerAuthorizer)

	policyClient := apimanagement.NewPolicyClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&policyClient.Client, o.ResourceManagerAuthorizer)

	productsClient := apimanagement.NewProductClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&productsClient.Client, o.ResourceManagerAuthorizer)

	productApisClient := apimanagement.NewProductAPIClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&productApisClient.Client, o.ResourceManagerAuthorizer)

	productGroupsClient := apimanagement.NewProductGroupClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&productGroupsClient.Client, o.ResourceManagerAuthorizer)

	productPoliciesClient := apimanagement.NewProductPolicyClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&productPoliciesClient.Client, o.ResourceManagerAuthorizer)

	serviceClient := apimanagement.NewServiceClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serviceClient.Client, o.ResourceManagerAuthorizer)

	signInClient := apimanagement.NewSignInSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&signInClient.Client, o.ResourceManagerAuthorizer)

	signUpClient := apimanagement.NewSignUpSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&signUpClient.Client, o.ResourceManagerAuthorizer)

	subscriptionsClient := apimanagement.NewSubscriptionClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&subscriptionsClient.Client, o.ResourceManagerAuthorizer)

	usersClient := apimanagement.NewUserClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&usersClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ApiClient:                  &apiClient,
		ApiDiagnosticClient:        &apiDiagnosticClient,
		ApiPoliciesClient:          &apiPoliciesClient,
		ApiOperationsClient:        &apiOperationsClient,
		ApiOperationPoliciesClient: &apiOperationPoliciesClient,
		ApiSchemasClient:           &apiSchemasClient,
		ApiVersionSetClient:        &apiVersionSetClient,
		AuthorizationServersClient: &authorizationServersClient,
		BackendClient:              &backendClient,
		CertificatesClient:         &certificatesClient,
		DiagnosticClient:           &diagnosticClient,
		GroupClient:                &groupClient,
		GroupUsersClient:           &groupUsersClient,
		IdentityProviderClient:     &identityProviderClient,
		LoggerClient:               &loggerClient,
		NamedValueClient:           &namedValueClient,
		OpenIdConnectClient:        &openIdConnectClient,
		PolicyClient:               &policyClient,
		ProductsClient:             &productsClient,
		ProductApisClient:          &productApisClient,
		ProductGroupsClient:        &productGroupsClient,
		ProductPoliciesClient:      &productPoliciesClient,
		ServiceClient:              &serviceClient,
		SignInClient:               &signInClient,
		SignUpClient:               &signUpClient,
		SubscriptionsClient:        &subscriptionsClient,
		UsersClient:                &usersClient,
	}
}
