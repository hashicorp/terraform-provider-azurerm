package apimanagement

import (
	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2018-01-01/apimanagement"
	"github.com/Azure/go-autorest/autorest"
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

func BuildClient(endpoint, subscriptionId, partnerId string, auth autorest.Authorizer, skipProviderReg bool) *Client {
	c := Client{}
	c.ApiClient = apimanagement.NewAPIClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ApiClient.Client, auth, partnerId, skipProviderReg)

	c.ApiPoliciesClient = apimanagement.NewAPIPolicyClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ApiPoliciesClient.Client, auth, partnerId, skipProviderReg)

	c.ApiOperationsClient = apimanagement.NewAPIOperationClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ApiOperationsClient.Client, auth, partnerId, skipProviderReg)

	c.ApiOperationPoliciesClient = apimanagement.NewAPIOperationPolicyClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ApiOperationPoliciesClient.Client, auth, partnerId, skipProviderReg)

	c.ApiSchemasClient = apimanagement.NewAPISchemaClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ApiSchemasClient.Client, auth, partnerId, skipProviderReg)

	c.ApiVersionSetClient = apimanagement.NewAPIVersionSetClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ApiVersionSetClient.Client, auth, partnerId, skipProviderReg)

	c.AuthorizationServersClient = apimanagement.NewAuthorizationServerClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.AuthorizationServersClient.Client, auth, partnerId, skipProviderReg)

	c.CertificatesClient = apimanagement.NewCertificateClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.CertificatesClient.Client, auth, partnerId, skipProviderReg)

	c.GroupClient = apimanagement.NewGroupClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.GroupClient.Client, auth, partnerId, skipProviderReg)

	c.GroupUsersClient = apimanagement.NewGroupUserClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.GroupUsersClient.Client, auth, partnerId, skipProviderReg)

	c.LoggerClient = apimanagement.NewLoggerClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.LoggerClient.Client, auth, partnerId, skipProviderReg)

	c.PolicyClient = apimanagement.NewPolicyClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.PolicyClient.Client, auth, partnerId, skipProviderReg)

	c.ServiceClient = apimanagement.NewServiceClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ServiceClient.Client, auth, partnerId, skipProviderReg)

	c.SignInClient = apimanagement.NewSignInSettingsClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.SignInClient.Client, auth, partnerId, skipProviderReg)

	c.SignUpClient = apimanagement.NewSignUpSettingsClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.SignUpClient.Client, auth, partnerId, skipProviderReg)

	c.OpenIdConnectClient = apimanagement.NewOpenIDConnectProviderClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.OpenIdConnectClient.Client, auth, partnerId, skipProviderReg)

	c.ProductsClient = apimanagement.NewProductClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ProductsClient.Client, auth, partnerId, skipProviderReg)

	c.ProductApisClient = apimanagement.NewProductAPIClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ProductApisClient.Client, auth, partnerId, skipProviderReg)

	c.ProductGroupsClient = apimanagement.NewProductGroupClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ProductGroupsClient.Client, auth, partnerId, skipProviderReg)

	c.ProductPoliciesClient = apimanagement.NewProductPolicyClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ProductPoliciesClient.Client, auth, partnerId, skipProviderReg)

	c.PropertyClient = apimanagement.NewPropertyClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.PropertyClient.Client, auth, partnerId, skipProviderReg)

	c.SubscriptionsClient = apimanagement.NewSubscriptionClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.SubscriptionsClient.Client, auth, partnerId, skipProviderReg)

	c.UsersClient = apimanagement.NewUserClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.UsersClient.Client, auth, partnerId, skipProviderReg)

	return &c
}
