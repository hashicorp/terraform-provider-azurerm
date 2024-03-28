// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apidiagnostic"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apimanagementservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apioperation"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apioperationpolicy"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apioperationtag"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apipolicy"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apirelease"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apischema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apitag"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apitagdescription"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apiversionset"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apiversionsets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/authorizationserver"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/backend"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/cache"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/certificate"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/delegationsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/deletedservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/diagnostic"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/emailtemplates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/gateway"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/gatewayapi"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/gatewaycertificateauthority"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/gatewayhostnameconfiguration"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/group"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/groupuser"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/identityprovider"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/logger"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/namedvalue"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/notificationrecipientemail"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/notificationrecipientuser"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/openidconnectprovider"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/policy"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/policyfragment"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/product"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/productapi"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/productgroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/productpolicy"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/producttag"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/schema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/signinsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/signupsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/subscription"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/tag"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/tenantaccess"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/user"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ApiClient                          *api.ApiClient
	ApiDiagnosticClient                *apidiagnostic.ApiDiagnosticClient
	ApiOperationPoliciesClient         *apioperationpolicy.ApiOperationPolicyClient
	ApiOperationsClient                *apioperation.ApiOperationClient
	ApiOperationTagClient              *apioperationtag.ApiOperationTagClient
	ApiPoliciesClient                  *apipolicy.ApiPolicyClient
	ApiReleasesClient                  *apirelease.ApiReleaseClient
	ApiSchemasClient                   *apischema.ApiSchemaClient
	ApiTagClient                       *apitag.ApiTagClient
	ApiTagDescriptionClient            *apitagdescription.ApiTagDescriptionClient
	ApiVersionSetClient                *apiversionset.ApiVersionSetClient
	ApiVersionSetsClient               *apiversionsets.ApiVersionSetsClient
	AuthorizationServersClient         *authorizationserver.AuthorizationServerClient
	BackendClient                      *backend.BackendClient
	CacheClient                        *cache.CacheClient
	CertificatesClient                 *certificate.CertificateClient
	DelegationSettingsClient           *delegationsettings.DelegationSettingsClient
	DeletedServicesClient              *deletedservice.DeletedServiceClient
	DiagnosticClient                   *diagnostic.DiagnosticClient
	EmailTemplatesClient               *emailtemplates.EmailTemplatesClient
	GatewayApisClient                  *gatewayapi.GatewayApiClient
	GatewayCertificateAuthorityClient  *gatewaycertificateauthority.GatewayCertificateAuthorityClient
	GatewayClient                      *gateway.GatewayClient
	GatewayHostNameConfigurationClient *gatewayhostnameconfiguration.GatewayHostnameConfigurationClient
	GlobalSchemaClient                 *schema.SchemaClient
	GroupClient                        *group.GroupClient
	GroupUsersClient                   *groupuser.GroupUserClient
	IdentityProviderClient             *identityprovider.IdentityProviderClient
	LoggerClient                       *logger.LoggerClient
	NamedValueClient                   *namedvalue.NamedValueClient
	NotificationRecipientEmailClient   *notificationrecipientemail.NotificationRecipientEmailClient
	NotificationRecipientUserClient    *notificationrecipientuser.NotificationRecipientUserClient
	OpenIdConnectClient                *openidconnectprovider.OpenidConnectProviderClient
	PolicyClient                       *policy.PolicyClient
	PolicyFragmentClient               *policyfragment.PolicyFragmentClient
	ProductApisClient                  *productapi.ProductApiClient
	ProductGroupsClient                *productgroup.ProductGroupClient
	ProductPoliciesClient              *productpolicy.ProductPolicyClient
	ProductsClient                     *product.ProductClient
	ProductTagClient                   *producttag.ProductTagClient
	ServiceClient                      *apimanagementservice.ApiManagementServiceClient
	SignInClient                       *signinsettings.SignInSettingsClient
	SignUpClient                       *signupsettings.SignUpSettingsClient
	SubscriptionsClient                *subscription.SubscriptionClient
	TagClient                          *tag.TagClient
	TenantAccessClient                 *tenantaccess.TenantAccessClient
	UsersClient                        *user.UserClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	apiClient, err := api.NewApiClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Api client: %+v", err)
	}
	o.Configure(apiClient.Client, o.Authorizers.ResourceManager)

	apiDiagnosticClient, err := apidiagnostic.NewApiDiagnosticClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Api Diagnostic client: %+v", err)
	}
	o.Configure(apiDiagnosticClient.Client, o.Authorizers.ResourceManager)

	apiPoliciesClient, err := apipolicy.NewApiPolicyClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Api Policies client: %+v", err)
	}
	o.Configure(apiPoliciesClient.Client, o.Authorizers.ResourceManager)

	apiOperationsClient, err := apioperation.NewApiOperationClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Api Operations client: %+v", err)
	}
	o.Configure(apiOperationsClient.Client, o.Authorizers.ResourceManager)

	apiOperationTagClient, err := apioperationtag.NewApiOperationTagClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Api Operation Tag client: %+v", err)
	}
	o.Configure(apiOperationTagClient.Client, o.Authorizers.ResourceManager)

	apiOperationPoliciesClient, err := apioperationpolicy.NewApiOperationPolicyClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Api Operation Policies client: %+v", err)
	}
	o.Configure(apiOperationPoliciesClient.Client, o.Authorizers.ResourceManager)

	apiReleasesClient, err := apirelease.NewApiReleaseClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Api Releases client: %+v", err)
	}
	o.Configure(apiReleasesClient.Client, o.Authorizers.ResourceManager)

	apiSchemasClient, err := apischema.NewApiSchemaClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Api Schemas client: %+v", err)
	}
	o.Configure(apiSchemasClient.Client, o.Authorizers.ResourceManager)

	apiTagClient, err := apitag.NewApiTagClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Api Tag client: %+v", err)
	}
	o.Configure(apiTagClient.Client, o.Authorizers.ResourceManager)

	apiVersionSetClient, err := apiversionset.NewApiVersionSetClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Api Version Set client: %+v", err)
	}
	o.Configure(apiVersionSetClient.Client, o.Authorizers.ResourceManager)

	apiVersionSetsClient, err := apiversionsets.NewApiVersionSetsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Api Version Sets client: %+v", err)
	}
	o.Configure(apiVersionSetsClient.Client, o.Authorizers.ResourceManager)

	apiTagDescriptionClient, err := apitagdescription.NewApiTagDescriptionClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Api Tag Description client: %+v", err)
	}
	o.Configure(apiTagDescriptionClient.Client, o.Authorizers.ResourceManager)

	authorizationServersClient, err := authorizationserver.NewAuthorizationServerClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Authorization Servers client: %+v", err)
	}
	o.Configure(authorizationServersClient.Client, o.Authorizers.ResourceManager)

	backendClient, err := backend.NewBackendClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building backend client: %+v", err)
	}
	o.Configure(backendClient.Client, o.Authorizers.ResourceManager)

	cacheClient, err := cache.NewCacheClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building cache client: %+v", err)
	}
	o.Configure(cacheClient.Client, o.Authorizers.ResourceManager)

	certificatesClient, err := certificate.NewCertificateClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building certificates client: %+v", err)
	}
	o.Configure(certificatesClient.Client, o.Authorizers.ResourceManager)

	diagnosticClient, err := diagnostic.NewDiagnosticClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building diagnostic client: %+v", err)
	}
	o.Configure(diagnosticClient.Client, o.Authorizers.ResourceManager)

	delegationSettingsClient, err := delegationsettings.NewDelegationSettingsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Delegation Settings client: %+v", err)
	}
	o.Configure(delegationSettingsClient.Client, o.Authorizers.ResourceManager)

	deletedServicesClient, err := deletedservice.NewDeletedServiceClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Deleted Services client: %+v", err)
	}
	o.Configure(deletedServicesClient.Client, o.Authorizers.ResourceManager)

	emailTemplatesClient, err := emailtemplates.NewEmailTemplatesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Email Templates client: %+v", err)
	}
	o.Configure(emailTemplatesClient.Client, o.Authorizers.ResourceManager)

	gatewayClient, err := gateway.NewGatewayClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Gateway client: %+v", err)
	}
	o.Configure(gatewayClient.Client, o.Authorizers.ResourceManager)

	gatewayCertificateAuthorityClient, err := gatewaycertificateauthority.NewGatewayCertificateAuthorityClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Gateway Certificate Authority client: %+v", err)
	}
	o.Configure(gatewayCertificateAuthorityClient.Client, o.Authorizers.ResourceManager)

	gatewayApisClient, err := gatewayapi.NewGatewayApiClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Gateway Apis client: %+v", err)
	}
	o.Configure(gatewayApisClient.Client, o.Authorizers.ResourceManager)

	gatewayHostnameConfigurationClient, err := gatewayhostnameconfiguration.NewGatewayHostnameConfigurationClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Gateway Hostname Configuration client: %+v", err)
	}
	o.Configure(gatewayHostnameConfigurationClient.Client, o.Authorizers.ResourceManager)

	globalSchemaClient, err := schema.NewSchemaClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Schema client: %+v", err)
	}
	o.Configure(globalSchemaClient.Client, o.Authorizers.ResourceManager)

	groupClient, err := group.NewGroupClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Group client: %+v", err)
	}
	o.Configure(groupClient.Client, o.Authorizers.ResourceManager)

	groupUsersClient, err := groupuser.NewGroupUserClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Group Users client: %+v", err)
	}
	o.Configure(groupUsersClient.Client, o.Authorizers.ResourceManager)

	identityProviderClient, err := identityprovider.NewIdentityProviderClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Identity Provider client: %+v", err)
	}
	o.Configure(identityProviderClient.Client, o.Authorizers.ResourceManager)

	namedValueClient, err := namedvalue.NewNamedValueClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Named Value client: %+v", err)
	}
	o.Configure(namedValueClient.Client, o.Authorizers.ResourceManager)

	notificationRecipientEmailClient, err := notificationrecipientemail.NewNotificationRecipientEmailClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Notification Recipient Email client: %+v", err)
	}
	o.Configure(notificationRecipientEmailClient.Client, o.Authorizers.ResourceManager)

	notificationRecipientUserClient, err := notificationrecipientuser.NewNotificationRecipientUserClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Notification Recipient User client: %+v", err)
	}
	o.Configure(notificationRecipientUserClient.Client, o.Authorizers.ResourceManager)

	loggerClient, err := logger.NewLoggerClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Logger client: %+v", err)
	}
	o.Configure(loggerClient.Client, o.Authorizers.ResourceManager)

	openIdConnectClient, err := openidconnectprovider.NewOpenidConnectProviderClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building OpenId Connect client: %+v", err)
	}
	o.Configure(openIdConnectClient.Client, o.Authorizers.ResourceManager)

	policyClient, err := policy.NewPolicyClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Policy client: %+v", err)
	}
	o.Configure(policyClient.Client, o.Authorizers.ResourceManager)

	policyFragmentClient, err := policyfragment.NewPolicyFragmentClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Policy Fragment client: %+v", err)
	}
	o.Configure(policyFragmentClient.Client, o.Authorizers.ResourceManager)

	productsClient, err := product.NewProductClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Products client: %+v", err)
	}
	o.Configure(productsClient.Client, o.Authorizers.ResourceManager)

	productTagClient, err := producttag.NewProductTagClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Product Tag client: %+v", err)
	}
	o.Configure(productTagClient.Client, o.Authorizers.ResourceManager)

	productApisClient, err := productapi.NewProductApiClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Product Apis client: %+v", err)
	}
	o.Configure(productApisClient.Client, o.Authorizers.ResourceManager)

	productGroupsClient, err := productgroup.NewProductGroupClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Product Groups client: %+v", err)
	}
	o.Configure(productGroupsClient.Client, o.Authorizers.ResourceManager)

	productPoliciesClient, err := productpolicy.NewProductPolicyClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Product Policies client: %+v", err)
	}
	o.Configure(productPoliciesClient.Client, o.Authorizers.ResourceManager)

	serviceClient, err := apimanagementservice.NewApiManagementServiceClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Service client: %+v", err)
	}
	o.Configure(serviceClient.Client, o.Authorizers.ResourceManager)

	signInClient, err := signinsettings.NewSignInSettingsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building SignIn client: %+v", err)
	}
	o.Configure(signInClient.Client, o.Authorizers.ResourceManager)

	signUpClient, err := signupsettings.NewSignUpSettingsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building SignUp client: %+v", err)
	}
	o.Configure(signUpClient.Client, o.Authorizers.ResourceManager)

	subscriptionsClient, err := subscription.NewSubscriptionClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Subscriptions client: %+v", err)
	}
	o.Configure(subscriptionsClient.Client, o.Authorizers.ResourceManager)

	tagClient, err := tag.NewTagClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building tag client: %+v", err)
	}
	o.Configure(tagClient.Client, o.Authorizers.ResourceManager)

	tenantAccessClient, err := tenantaccess.NewTenantAccessClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Tenant Access client: %+v", err)
	}
	o.Configure(tenantAccessClient.Client, o.Authorizers.ResourceManager)

	usersClient, err := user.NewUserClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building users client: %+v", err)
	}
	o.Configure(usersClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ApiClient:                          apiClient,
		ApiDiagnosticClient:                apiDiagnosticClient,
		ApiOperationPoliciesClient:         apiOperationPoliciesClient,
		ApiOperationsClient:                apiOperationsClient,
		ApiOperationTagClient:              apiOperationTagClient,
		ApiPoliciesClient:                  apiPoliciesClient,
		ApiReleasesClient:                  apiReleasesClient,
		ApiSchemasClient:                   apiSchemasClient,
		ApiTagClient:                       apiTagClient,
		ApiTagDescriptionClient:            apiTagDescriptionClient,
		ApiVersionSetClient:                apiVersionSetClient,
		ApiVersionSetsClient:               apiVersionSetsClient,
		AuthorizationServersClient:         authorizationServersClient,
		BackendClient:                      backendClient,
		CacheClient:                        cacheClient,
		CertificatesClient:                 certificatesClient,
		DelegationSettingsClient:           delegationSettingsClient,
		DeletedServicesClient:              deletedServicesClient,
		DiagnosticClient:                   diagnosticClient,
		EmailTemplatesClient:               emailTemplatesClient,
		GatewayApisClient:                  gatewayApisClient,
		GatewayCertificateAuthorityClient:  gatewayCertificateAuthorityClient,
		GatewayClient:                      gatewayClient,
		GatewayHostNameConfigurationClient: gatewayHostnameConfigurationClient,
		GlobalSchemaClient:                 globalSchemaClient,
		GroupClient:                        groupClient,
		GroupUsersClient:                   groupUsersClient,
		IdentityProviderClient:             identityProviderClient,
		LoggerClient:                       loggerClient,
		NamedValueClient:                   namedValueClient,
		NotificationRecipientEmailClient:   notificationRecipientEmailClient,
		NotificationRecipientUserClient:    notificationRecipientUserClient,
		OpenIdConnectClient:                openIdConnectClient,
		PolicyClient:                       policyClient,
		PolicyFragmentClient:               policyFragmentClient,
		ProductApisClient:                  productApisClient,
		ProductGroupsClient:                productGroupsClient,
		ProductPoliciesClient:              productPoliciesClient,
		ProductsClient:                     productsClient,
		ProductTagClient:                   productTagClient,
		ServiceClient:                      serviceClient,
		SignInClient:                       signInClient,
		SignUpClient:                       signUpClient,
		SubscriptionsClient:                subscriptionsClient,
		TagClient:                          tagClient,
		TenantAccessClient:                 tenantAccessClient,
		UsersClient:                        usersClient,
	}, nil
}
