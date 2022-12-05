package client

import (
	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2021-08-01/apimanagement"
	pandoraAPIMGlobalSchema "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ApiClient                          *apimanagement.APIClient
	ApiDiagnosticClient                *apimanagement.APIDiagnosticClient
	ApiOperationPoliciesClient         *apimanagement.APIOperationPolicyClient
	ApiOperationsClient                *apimanagement.APIOperationClient
	ApiPoliciesClient                  *apimanagement.APIPolicyClient
	ApiReleasesClient                  *apimanagement.APIReleaseClient
	ApiSchemasClient                   *apimanagement.APISchemaClient
	ApiTagDescriptionClient            *apimanagement.APITagDescriptionClient
	ApiVersionSetClient                *apimanagement.APIVersionSetClient
	AuthorizationServersClient         *apimanagement.AuthorizationServerClient
	BackendClient                      *apimanagement.BackendClient
	CacheClient                        *apimanagement.CacheClient
	CertificatesClient                 *apimanagement.CertificateClient
	DeletedServicesClient              *apimanagement.DeletedServicesClient
	DiagnosticClient                   *apimanagement.DiagnosticClient
	EmailTemplateClient                *apimanagement.EmailTemplateClient
	GatewayApisClient                  *apimanagement.GatewayAPIClient
	GatewayCertificateAuthorityClient  *apimanagement.GatewayCertificateAuthorityClient
	GatewayClient                      *apimanagement.GatewayClient
	GatewayHostNameConfigurationClient *apimanagement.GatewayHostnameConfigurationClient
	GlobalSchemaClient                 *pandoraAPIMGlobalSchema.SchemaClient
	GroupClient                        *apimanagement.GroupClient
	GroupUsersClient                   *apimanagement.GroupUserClient
	IdentityProviderClient             *apimanagement.IdentityProviderClient
	LoggerClient                       *apimanagement.LoggerClient
	NamedValueClient                   *apimanagement.NamedValueClient
	NotificationRecipientEmailClient   *apimanagement.NotificationRecipientEmailClient
	NotificationRecipientUserClient    *apimanagement.NotificationRecipientUserClient
	OpenIdConnectClient                *apimanagement.OpenIDConnectProviderClient
	PolicyClient                       *apimanagement.PolicyClient
	ProductApisClient                  *apimanagement.ProductAPIClient
	ProductGroupsClient                *apimanagement.ProductGroupClient
	ProductPoliciesClient              *apimanagement.ProductPolicyClient
	ProductsClient                     *apimanagement.ProductClient
	ServiceClient                      *apimanagement.ServiceClient
	SignInClient                       *apimanagement.SignInSettingsClient
	SignUpClient                       *apimanagement.SignUpSettingsClient
	SubscriptionsClient                *apimanagement.SubscriptionClient
	TagClient                          *apimanagement.TagClient
	TenantAccessClient                 *apimanagement.TenantAccessClient
	UsersClient                        *apimanagement.UserClient
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

	apiReleasesClient := apimanagement.NewAPIReleaseClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&apiReleasesClient.Client, o.ResourceManagerAuthorizer)

	apiSchemasClient := apimanagement.NewAPISchemaClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&apiSchemasClient.Client, o.ResourceManagerAuthorizer)

	apiVersionSetClient := apimanagement.NewAPIVersionSetClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&apiVersionSetClient.Client, o.ResourceManagerAuthorizer)

	apiTagDescriptionClient := apimanagement.NewAPITagDescriptionClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&apiTagDescriptionClient.Client, o.ResourceManagerAuthorizer)

	authorizationServersClient := apimanagement.NewAuthorizationServerClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&authorizationServersClient.Client, o.ResourceManagerAuthorizer)

	backendClient := apimanagement.NewBackendClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&backendClient.Client, o.ResourceManagerAuthorizer)

	cacheClient := apimanagement.NewCacheClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&cacheClient.Client, o.ResourceManagerAuthorizer)

	certificatesClient := apimanagement.NewCertificateClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&certificatesClient.Client, o.ResourceManagerAuthorizer)

	diagnosticClient := apimanagement.NewDiagnosticClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&diagnosticClient.Client, o.ResourceManagerAuthorizer)

	deletedServicesClient := apimanagement.NewDeletedServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&deletedServicesClient.Client, o.ResourceManagerAuthorizer)

	emailTemplateClient := apimanagement.NewEmailTemplateClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&emailTemplateClient.Client, o.ResourceManagerAuthorizer)

	gatewayClient := apimanagement.NewGatewayClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&gatewayClient.Client, o.ResourceManagerAuthorizer)

	gatewayCertificateAuthorityClient := apimanagement.NewGatewayCertificateAuthorityClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&gatewayCertificateAuthorityClient.Client, o.ResourceManagerAuthorizer)

	gatewayApisClient := apimanagement.NewGatewayAPIClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&gatewayApisClient.Client, o.ResourceManagerAuthorizer)

	gatewayHostnameConfigurationClient := apimanagement.NewGatewayHostnameConfigurationClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&gatewayHostnameConfigurationClient.Client, o.ResourceManagerAuthorizer)

	globalSchemaClient := pandoraAPIMGlobalSchema.NewSchemaClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&globalSchemaClient.Client, o.ResourceManagerAuthorizer)

	groupClient := apimanagement.NewGroupClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&groupClient.Client, o.ResourceManagerAuthorizer)

	groupUsersClient := apimanagement.NewGroupUserClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&groupUsersClient.Client, o.ResourceManagerAuthorizer)

	identityProviderClient := apimanagement.NewIdentityProviderClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&identityProviderClient.Client, o.ResourceManagerAuthorizer)

	namedValueClient := apimanagement.NewNamedValueClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&namedValueClient.Client, o.ResourceManagerAuthorizer)

	notificationRecipientEmailClient := apimanagement.NewNotificationRecipientEmailClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&notificationRecipientEmailClient.Client, o.ResourceManagerAuthorizer)

	notificationRecipientUserClient := apimanagement.NewNotificationRecipientUserClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&notificationRecipientUserClient.Client, o.ResourceManagerAuthorizer)

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

	tagClient := apimanagement.NewTagClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&tagClient.Client, o.ResourceManagerAuthorizer)

	tenantAccessClient := apimanagement.NewTenantAccessClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&tenantAccessClient.Client, o.ResourceManagerAuthorizer)

	usersClient := apimanagement.NewUserClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&usersClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ApiClient:                          &apiClient,
		ApiDiagnosticClient:                &apiDiagnosticClient,
		ApiOperationPoliciesClient:         &apiOperationPoliciesClient,
		ApiOperationsClient:                &apiOperationsClient,
		ApiPoliciesClient:                  &apiPoliciesClient,
		ApiReleasesClient:                  &apiReleasesClient,
		ApiSchemasClient:                   &apiSchemasClient,
		ApiTagDescriptionClient:            &apiTagDescriptionClient,
		ApiVersionSetClient:                &apiVersionSetClient,
		AuthorizationServersClient:         &authorizationServersClient,
		BackendClient:                      &backendClient,
		CacheClient:                        &cacheClient,
		CertificatesClient:                 &certificatesClient,
		DeletedServicesClient:              &deletedServicesClient,
		DiagnosticClient:                   &diagnosticClient,
		EmailTemplateClient:                &emailTemplateClient,
		GatewayApisClient:                  &gatewayApisClient,
		GatewayCertificateAuthorityClient:  &gatewayCertificateAuthorityClient,
		GatewayClient:                      &gatewayClient,
		GatewayHostNameConfigurationClient: &gatewayHostnameConfigurationClient,
		GlobalSchemaClient:                 &globalSchemaClient,
		GroupClient:                        &groupClient,
		GroupUsersClient:                   &groupUsersClient,
		IdentityProviderClient:             &identityProviderClient,
		LoggerClient:                       &loggerClient,
		NamedValueClient:                   &namedValueClient,
		NotificationRecipientEmailClient:   &notificationRecipientEmailClient,
		NotificationRecipientUserClient:    &notificationRecipientUserClient,
		OpenIdConnectClient:                &openIdConnectClient,
		PolicyClient:                       &policyClient,
		ProductApisClient:                  &productApisClient,
		ProductGroupsClient:                &productGroupsClient,
		ProductPoliciesClient:              &productPoliciesClient,
		ProductsClient:                     &productsClient,
		ServiceClient:                      &serviceClient,
		SignInClient:                       &signInClient,
		SignUpClient:                       &signUpClient,
		SubscriptionsClient:                &subscriptionsClient,
		TagClient:                          &tagClient,
		TenantAccessClient:                 &tenantAccessClient,
		UsersClient:                        &usersClient,
	}
}
