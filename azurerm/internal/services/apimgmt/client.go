package apimgmt

import "github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2018-01-01/apimanagement"

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
