package v2024_05_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/allpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/api"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apidiagnostic"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apigateway"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apigatewayconfigconnection"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apiissue"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apiissueattachment"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apiissuecomment"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementgatewayskus"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementserviceskus"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementworkspacelinks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apioperation"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apioperationpolicy"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apioperationsbytag"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apioperationtag"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apipolicy"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apiproduct"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apirelease"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apirevision"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apisbytag"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apischema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apitag"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apitagdescription"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apiversionset"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apiversionsets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apiwiki"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/authorization"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/authorizationaccesspolicy"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/authorizationconfirmconsentcode"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/authorizationloginlinks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/authorizationprovider"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/authorizations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/authorizationserver"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/backend"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/backendreconnect"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/cache"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/certificate"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/contenttype"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/contenttypecontentitem"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/delegationsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/deletedservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/diagnostic"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/documentationresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/emailtemplate"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/emailtemplates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/gateway"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/gatewayapi"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/gatewaycertificateauthority"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/gatewaygeneratetoken"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/gatewayhostnameconfiguration"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/gatewayinvalidatedebugcredentials"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/gatewaylistdebugcredentials"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/gatewaylistkeys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/gatewaylisttrace"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/gatewayregeneratekey"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/graphqlapiresolver"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/graphqlapiresolverpolicy"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/group"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/groupuser"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/identityprovider"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/issue"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/logger"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/namedvalue"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/networkstatus"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/notification"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/notificationrecipientemail"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/notificationrecipientuser"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/openidconnectprovider"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/outboundnetworkdependenciesendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/performconnectivitycheck"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/policy"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/policydescription"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/policyfragment"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/policyrestriction"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/policyrestrictions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/policyrestrictionsvalidations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/portalconfig"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/portalrevision"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/portalsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/privateendpointconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/product"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/productapi"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/productapilink"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/productgroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/productgrouplink"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/productpolicy"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/productsbytag"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/productsubscription"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/producttag"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/productwiki"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/quotabycounterkeys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/quotabyperiodkeys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/region"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/reports"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/schema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/signinsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/signupsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/skus"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/subscription"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/tag"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/tagapilink"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/tagoperationlink"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/tagproductlink"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/tagresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/tenantaccess"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/tenantaccessgit"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/tenantconfiguration"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/tenantconfigurationsyncstate"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/tenantsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/user"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/userconfirmationpasswordsend"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/usergroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/useridentity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/users"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/usersubscription"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/usertoken"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspace"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspacepolicy"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	AllPolicies                          *allpolicies.AllPoliciesClient
	Api                                  *api.ApiClient
	ApiDiagnostic                        *apidiagnostic.ApiDiagnosticClient
	ApiGateway                           *apigateway.ApiGatewayClient
	ApiGatewayConfigConnection           *apigatewayconfigconnection.ApiGatewayConfigConnectionClient
	ApiIssue                             *apiissue.ApiIssueClient
	ApiIssueAttachment                   *apiissueattachment.ApiIssueAttachmentClient
	ApiIssueComment                      *apiissuecomment.ApiIssueCommentClient
	ApiManagementGatewaySkus             *apimanagementgatewayskus.ApiManagementGatewaySkusClient
	ApiManagementService                 *apimanagementservice.ApiManagementServiceClient
	ApiManagementServiceSkus             *apimanagementserviceskus.ApiManagementServiceSkusClient
	ApiManagementWorkspaceLinks          *apimanagementworkspacelinks.ApiManagementWorkspaceLinksClient
	ApiOperation                         *apioperation.ApiOperationClient
	ApiOperationPolicy                   *apioperationpolicy.ApiOperationPolicyClient
	ApiOperationTag                      *apioperationtag.ApiOperationTagClient
	ApiOperationsByTag                   *apioperationsbytag.ApiOperationsByTagClient
	ApiPolicy                            *apipolicy.ApiPolicyClient
	ApiProduct                           *apiproduct.ApiProductClient
	ApiRelease                           *apirelease.ApiReleaseClient
	ApiRevision                          *apirevision.ApiRevisionClient
	ApiSchema                            *apischema.ApiSchemaClient
	ApiTag                               *apitag.ApiTagClient
	ApiTagDescription                    *apitagdescription.ApiTagDescriptionClient
	ApiVersionSet                        *apiversionset.ApiVersionSetClient
	ApiVersionSets                       *apiversionsets.ApiVersionSetsClient
	ApiWiki                              *apiwiki.ApiWikiClient
	ApisByTag                            *apisbytag.ApisByTagClient
	Authorization                        *authorization.AuthorizationClient
	AuthorizationAccessPolicy            *authorizationaccesspolicy.AuthorizationAccessPolicyClient
	AuthorizationConfirmConsentCode      *authorizationconfirmconsentcode.AuthorizationConfirmConsentCodeClient
	AuthorizationLoginLinks              *authorizationloginlinks.AuthorizationLoginLinksClient
	AuthorizationProvider                *authorizationprovider.AuthorizationProviderClient
	AuthorizationServer                  *authorizationserver.AuthorizationServerClient
	Authorizations                       *authorizations.AuthorizationsClient
	Backend                              *backend.BackendClient
	BackendReconnect                     *backendreconnect.BackendReconnectClient
	Cache                                *cache.CacheClient
	Certificate                          *certificate.CertificateClient
	ContentType                          *contenttype.ContentTypeClient
	ContentTypeContentItem               *contenttypecontentitem.ContentTypeContentItemClient
	DelegationSettings                   *delegationsettings.DelegationSettingsClient
	DeletedService                       *deletedservice.DeletedServiceClient
	Diagnostic                           *diagnostic.DiagnosticClient
	DocumentationResource                *documentationresource.DocumentationResourceClient
	EmailTemplate                        *emailtemplate.EmailTemplateClient
	EmailTemplates                       *emailtemplates.EmailTemplatesClient
	Gateway                              *gateway.GatewayClient
	GatewayApi                           *gatewayapi.GatewayApiClient
	GatewayCertificateAuthority          *gatewaycertificateauthority.GatewayCertificateAuthorityClient
	GatewayGenerateToken                 *gatewaygeneratetoken.GatewayGenerateTokenClient
	GatewayHostnameConfiguration         *gatewayhostnameconfiguration.GatewayHostnameConfigurationClient
	GatewayInvalidateDebugCredentials    *gatewayinvalidatedebugcredentials.GatewayInvalidateDebugCredentialsClient
	GatewayListDebugCredentials          *gatewaylistdebugcredentials.GatewayListDebugCredentialsClient
	GatewayListKeys                      *gatewaylistkeys.GatewayListKeysClient
	GatewayListTrace                     *gatewaylisttrace.GatewayListTraceClient
	GatewayRegenerateKey                 *gatewayregeneratekey.GatewayRegenerateKeyClient
	GraphQLApiResolver                   *graphqlapiresolver.GraphQLApiResolverClient
	GraphQLApiResolverPolicy             *graphqlapiresolverpolicy.GraphQLApiResolverPolicyClient
	Group                                *group.GroupClient
	GroupUser                            *groupuser.GroupUserClient
	IdentityProvider                     *identityprovider.IdentityProviderClient
	Issue                                *issue.IssueClient
	Logger                               *logger.LoggerClient
	NamedValue                           *namedvalue.NamedValueClient
	NetworkStatus                        *networkstatus.NetworkStatusClient
	Notification                         *notification.NotificationClient
	NotificationRecipientEmail           *notificationrecipientemail.NotificationRecipientEmailClient
	NotificationRecipientUser            *notificationrecipientuser.NotificationRecipientUserClient
	OpenidConnectProvider                *openidconnectprovider.OpenidConnectProviderClient
	OutboundNetworkDependenciesEndpoints *outboundnetworkdependenciesendpoints.OutboundNetworkDependenciesEndpointsClient
	PerformConnectivityCheck             *performconnectivitycheck.PerformConnectivityCheckClient
	Policy                               *policy.PolicyClient
	PolicyDescription                    *policydescription.PolicyDescriptionClient
	PolicyFragment                       *policyfragment.PolicyFragmentClient
	PolicyRestriction                    *policyrestriction.PolicyRestrictionClient
	PolicyRestrictions                   *policyrestrictions.PolicyRestrictionsClient
	PolicyRestrictionsValidations        *policyrestrictionsvalidations.PolicyRestrictionsValidationsClient
	PortalConfig                         *portalconfig.PortalConfigClient
	PortalRevision                       *portalrevision.PortalRevisionClient
	PortalSettings                       *portalsettings.PortalSettingsClient
	PrivateEndpointConnections           *privateendpointconnections.PrivateEndpointConnectionsClient
	Product                              *product.ProductClient
	ProductApi                           *productapi.ProductApiClient
	ProductApiLink                       *productapilink.ProductApiLinkClient
	ProductGroup                         *productgroup.ProductGroupClient
	ProductGroupLink                     *productgrouplink.ProductGroupLinkClient
	ProductPolicy                        *productpolicy.ProductPolicyClient
	ProductSubscription                  *productsubscription.ProductSubscriptionClient
	ProductTag                           *producttag.ProductTagClient
	ProductWiki                          *productwiki.ProductWikiClient
	ProductsByTag                        *productsbytag.ProductsByTagClient
	QuotaByCounterKeys                   *quotabycounterkeys.QuotaByCounterKeysClient
	QuotaByPeriodKeys                    *quotabyperiodkeys.QuotaByPeriodKeysClient
	Region                               *region.RegionClient
	Reports                              *reports.ReportsClient
	Schema                               *schema.SchemaClient
	SignInSettings                       *signinsettings.SignInSettingsClient
	SignUpSettings                       *signupsettings.SignUpSettingsClient
	Skus                                 *skus.SkusClient
	Subscription                         *subscription.SubscriptionClient
	Tag                                  *tag.TagClient
	TagApiLink                           *tagapilink.TagApiLinkClient
	TagOperationLink                     *tagoperationlink.TagOperationLinkClient
	TagProductLink                       *tagproductlink.TagProductLinkClient
	TagResource                          *tagresource.TagResourceClient
	TenantAccess                         *tenantaccess.TenantAccessClient
	TenantAccessGit                      *tenantaccessgit.TenantAccessGitClient
	TenantConfiguration                  *tenantconfiguration.TenantConfigurationClient
	TenantConfigurationSyncState         *tenantconfigurationsyncstate.TenantConfigurationSyncStateClient
	TenantSettings                       *tenantsettings.TenantSettingsClient
	User                                 *user.UserClient
	UserConfirmationPasswordSend         *userconfirmationpasswordsend.UserConfirmationPasswordSendClient
	UserGroup                            *usergroup.UserGroupClient
	UserIdentity                         *useridentity.UserIdentityClient
	UserSubscription                     *usersubscription.UserSubscriptionClient
	UserToken                            *usertoken.UserTokenClient
	Users                                *users.UsersClient
	Workspace                            *workspace.WorkspaceClient
	WorkspacePolicy                      *workspacepolicy.WorkspacePolicyClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	allPoliciesClient, err := allpolicies.NewAllPoliciesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AllPolicies client: %+v", err)
	}
	configureFunc(allPoliciesClient.Client)

	apiClient, err := api.NewApiClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Api client: %+v", err)
	}
	configureFunc(apiClient.Client)

	apiDiagnosticClient, err := apidiagnostic.NewApiDiagnosticClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ApiDiagnostic client: %+v", err)
	}
	configureFunc(apiDiagnosticClient.Client)

	apiGatewayClient, err := apigateway.NewApiGatewayClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ApiGateway client: %+v", err)
	}
	configureFunc(apiGatewayClient.Client)

	apiGatewayConfigConnectionClient, err := apigatewayconfigconnection.NewApiGatewayConfigConnectionClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ApiGatewayConfigConnection client: %+v", err)
	}
	configureFunc(apiGatewayConfigConnectionClient.Client)

	apiIssueAttachmentClient, err := apiissueattachment.NewApiIssueAttachmentClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ApiIssueAttachment client: %+v", err)
	}
	configureFunc(apiIssueAttachmentClient.Client)

	apiIssueClient, err := apiissue.NewApiIssueClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ApiIssue client: %+v", err)
	}
	configureFunc(apiIssueClient.Client)

	apiIssueCommentClient, err := apiissuecomment.NewApiIssueCommentClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ApiIssueComment client: %+v", err)
	}
	configureFunc(apiIssueCommentClient.Client)

	apiManagementGatewaySkusClient, err := apimanagementgatewayskus.NewApiManagementGatewaySkusClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ApiManagementGatewaySkus client: %+v", err)
	}
	configureFunc(apiManagementGatewaySkusClient.Client)

	apiManagementServiceClient, err := apimanagementservice.NewApiManagementServiceClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ApiManagementService client: %+v", err)
	}
	configureFunc(apiManagementServiceClient.Client)

	apiManagementServiceSkusClient, err := apimanagementserviceskus.NewApiManagementServiceSkusClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ApiManagementServiceSkus client: %+v", err)
	}
	configureFunc(apiManagementServiceSkusClient.Client)

	apiManagementWorkspaceLinksClient, err := apimanagementworkspacelinks.NewApiManagementWorkspaceLinksClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ApiManagementWorkspaceLinks client: %+v", err)
	}
	configureFunc(apiManagementWorkspaceLinksClient.Client)

	apiOperationClient, err := apioperation.NewApiOperationClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ApiOperation client: %+v", err)
	}
	configureFunc(apiOperationClient.Client)

	apiOperationPolicyClient, err := apioperationpolicy.NewApiOperationPolicyClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ApiOperationPolicy client: %+v", err)
	}
	configureFunc(apiOperationPolicyClient.Client)

	apiOperationTagClient, err := apioperationtag.NewApiOperationTagClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ApiOperationTag client: %+v", err)
	}
	configureFunc(apiOperationTagClient.Client)

	apiOperationsByTagClient, err := apioperationsbytag.NewApiOperationsByTagClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ApiOperationsByTag client: %+v", err)
	}
	configureFunc(apiOperationsByTagClient.Client)

	apiPolicyClient, err := apipolicy.NewApiPolicyClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ApiPolicy client: %+v", err)
	}
	configureFunc(apiPolicyClient.Client)

	apiProductClient, err := apiproduct.NewApiProductClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ApiProduct client: %+v", err)
	}
	configureFunc(apiProductClient.Client)

	apiReleaseClient, err := apirelease.NewApiReleaseClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ApiRelease client: %+v", err)
	}
	configureFunc(apiReleaseClient.Client)

	apiRevisionClient, err := apirevision.NewApiRevisionClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ApiRevision client: %+v", err)
	}
	configureFunc(apiRevisionClient.Client)

	apiSchemaClient, err := apischema.NewApiSchemaClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ApiSchema client: %+v", err)
	}
	configureFunc(apiSchemaClient.Client)

	apiTagClient, err := apitag.NewApiTagClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ApiTag client: %+v", err)
	}
	configureFunc(apiTagClient.Client)

	apiTagDescriptionClient, err := apitagdescription.NewApiTagDescriptionClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ApiTagDescription client: %+v", err)
	}
	configureFunc(apiTagDescriptionClient.Client)

	apiVersionSetClient, err := apiversionset.NewApiVersionSetClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ApiVersionSet client: %+v", err)
	}
	configureFunc(apiVersionSetClient.Client)

	apiVersionSetsClient, err := apiversionsets.NewApiVersionSetsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ApiVersionSets client: %+v", err)
	}
	configureFunc(apiVersionSetsClient.Client)

	apiWikiClient, err := apiwiki.NewApiWikiClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ApiWiki client: %+v", err)
	}
	configureFunc(apiWikiClient.Client)

	apisByTagClient, err := apisbytag.NewApisByTagClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ApisByTag client: %+v", err)
	}
	configureFunc(apisByTagClient.Client)

	authorizationAccessPolicyClient, err := authorizationaccesspolicy.NewAuthorizationAccessPolicyClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AuthorizationAccessPolicy client: %+v", err)
	}
	configureFunc(authorizationAccessPolicyClient.Client)

	authorizationClient, err := authorization.NewAuthorizationClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Authorization client: %+v", err)
	}
	configureFunc(authorizationClient.Client)

	authorizationConfirmConsentCodeClient, err := authorizationconfirmconsentcode.NewAuthorizationConfirmConsentCodeClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AuthorizationConfirmConsentCode client: %+v", err)
	}
	configureFunc(authorizationConfirmConsentCodeClient.Client)

	authorizationLoginLinksClient, err := authorizationloginlinks.NewAuthorizationLoginLinksClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AuthorizationLoginLinks client: %+v", err)
	}
	configureFunc(authorizationLoginLinksClient.Client)

	authorizationProviderClient, err := authorizationprovider.NewAuthorizationProviderClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AuthorizationProvider client: %+v", err)
	}
	configureFunc(authorizationProviderClient.Client)

	authorizationServerClient, err := authorizationserver.NewAuthorizationServerClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AuthorizationServer client: %+v", err)
	}
	configureFunc(authorizationServerClient.Client)

	authorizationsClient, err := authorizations.NewAuthorizationsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Authorizations client: %+v", err)
	}
	configureFunc(authorizationsClient.Client)

	backendClient, err := backend.NewBackendClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Backend client: %+v", err)
	}
	configureFunc(backendClient.Client)

	backendReconnectClient, err := backendreconnect.NewBackendReconnectClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building BackendReconnect client: %+v", err)
	}
	configureFunc(backendReconnectClient.Client)

	cacheClient, err := cache.NewCacheClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Cache client: %+v", err)
	}
	configureFunc(cacheClient.Client)

	certificateClient, err := certificate.NewCertificateClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Certificate client: %+v", err)
	}
	configureFunc(certificateClient.Client)

	contentTypeClient, err := contenttype.NewContentTypeClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ContentType client: %+v", err)
	}
	configureFunc(contentTypeClient.Client)

	contentTypeContentItemClient, err := contenttypecontentitem.NewContentTypeContentItemClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ContentTypeContentItem client: %+v", err)
	}
	configureFunc(contentTypeContentItemClient.Client)

	delegationSettingsClient, err := delegationsettings.NewDelegationSettingsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building DelegationSettings client: %+v", err)
	}
	configureFunc(delegationSettingsClient.Client)

	deletedServiceClient, err := deletedservice.NewDeletedServiceClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building DeletedService client: %+v", err)
	}
	configureFunc(deletedServiceClient.Client)

	diagnosticClient, err := diagnostic.NewDiagnosticClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Diagnostic client: %+v", err)
	}
	configureFunc(diagnosticClient.Client)

	documentationResourceClient, err := documentationresource.NewDocumentationResourceClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building DocumentationResource client: %+v", err)
	}
	configureFunc(documentationResourceClient.Client)

	emailTemplateClient, err := emailtemplate.NewEmailTemplateClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building EmailTemplate client: %+v", err)
	}
	configureFunc(emailTemplateClient.Client)

	emailTemplatesClient, err := emailtemplates.NewEmailTemplatesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building EmailTemplates client: %+v", err)
	}
	configureFunc(emailTemplatesClient.Client)

	gatewayApiClient, err := gatewayapi.NewGatewayApiClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building GatewayApi client: %+v", err)
	}
	configureFunc(gatewayApiClient.Client)

	gatewayCertificateAuthorityClient, err := gatewaycertificateauthority.NewGatewayCertificateAuthorityClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building GatewayCertificateAuthority client: %+v", err)
	}
	configureFunc(gatewayCertificateAuthorityClient.Client)

	gatewayClient, err := gateway.NewGatewayClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Gateway client: %+v", err)
	}
	configureFunc(gatewayClient.Client)

	gatewayGenerateTokenClient, err := gatewaygeneratetoken.NewGatewayGenerateTokenClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building GatewayGenerateToken client: %+v", err)
	}
	configureFunc(gatewayGenerateTokenClient.Client)

	gatewayHostnameConfigurationClient, err := gatewayhostnameconfiguration.NewGatewayHostnameConfigurationClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building GatewayHostnameConfiguration client: %+v", err)
	}
	configureFunc(gatewayHostnameConfigurationClient.Client)

	gatewayInvalidateDebugCredentialsClient, err := gatewayinvalidatedebugcredentials.NewGatewayInvalidateDebugCredentialsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building GatewayInvalidateDebugCredentials client: %+v", err)
	}
	configureFunc(gatewayInvalidateDebugCredentialsClient.Client)

	gatewayListDebugCredentialsClient, err := gatewaylistdebugcredentials.NewGatewayListDebugCredentialsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building GatewayListDebugCredentials client: %+v", err)
	}
	configureFunc(gatewayListDebugCredentialsClient.Client)

	gatewayListKeysClient, err := gatewaylistkeys.NewGatewayListKeysClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building GatewayListKeys client: %+v", err)
	}
	configureFunc(gatewayListKeysClient.Client)

	gatewayListTraceClient, err := gatewaylisttrace.NewGatewayListTraceClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building GatewayListTrace client: %+v", err)
	}
	configureFunc(gatewayListTraceClient.Client)

	gatewayRegenerateKeyClient, err := gatewayregeneratekey.NewGatewayRegenerateKeyClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building GatewayRegenerateKey client: %+v", err)
	}
	configureFunc(gatewayRegenerateKeyClient.Client)

	graphQLApiResolverClient, err := graphqlapiresolver.NewGraphQLApiResolverClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building GraphQLApiResolver client: %+v", err)
	}
	configureFunc(graphQLApiResolverClient.Client)

	graphQLApiResolverPolicyClient, err := graphqlapiresolverpolicy.NewGraphQLApiResolverPolicyClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building GraphQLApiResolverPolicy client: %+v", err)
	}
	configureFunc(graphQLApiResolverPolicyClient.Client)

	groupClient, err := group.NewGroupClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Group client: %+v", err)
	}
	configureFunc(groupClient.Client)

	groupUserClient, err := groupuser.NewGroupUserClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building GroupUser client: %+v", err)
	}
	configureFunc(groupUserClient.Client)

	identityProviderClient, err := identityprovider.NewIdentityProviderClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building IdentityProvider client: %+v", err)
	}
	configureFunc(identityProviderClient.Client)

	issueClient, err := issue.NewIssueClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Issue client: %+v", err)
	}
	configureFunc(issueClient.Client)

	loggerClient, err := logger.NewLoggerClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Logger client: %+v", err)
	}
	configureFunc(loggerClient.Client)

	namedValueClient, err := namedvalue.NewNamedValueClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building NamedValue client: %+v", err)
	}
	configureFunc(namedValueClient.Client)

	networkStatusClient, err := networkstatus.NewNetworkStatusClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building NetworkStatus client: %+v", err)
	}
	configureFunc(networkStatusClient.Client)

	notificationClient, err := notification.NewNotificationClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Notification client: %+v", err)
	}
	configureFunc(notificationClient.Client)

	notificationRecipientEmailClient, err := notificationrecipientemail.NewNotificationRecipientEmailClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building NotificationRecipientEmail client: %+v", err)
	}
	configureFunc(notificationRecipientEmailClient.Client)

	notificationRecipientUserClient, err := notificationrecipientuser.NewNotificationRecipientUserClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building NotificationRecipientUser client: %+v", err)
	}
	configureFunc(notificationRecipientUserClient.Client)

	openidConnectProviderClient, err := openidconnectprovider.NewOpenidConnectProviderClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building OpenidConnectProvider client: %+v", err)
	}
	configureFunc(openidConnectProviderClient.Client)

	outboundNetworkDependenciesEndpointsClient, err := outboundnetworkdependenciesendpoints.NewOutboundNetworkDependenciesEndpointsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building OutboundNetworkDependenciesEndpoints client: %+v", err)
	}
	configureFunc(outboundNetworkDependenciesEndpointsClient.Client)

	performConnectivityCheckClient, err := performconnectivitycheck.NewPerformConnectivityCheckClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PerformConnectivityCheck client: %+v", err)
	}
	configureFunc(performConnectivityCheckClient.Client)

	policyClient, err := policy.NewPolicyClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Policy client: %+v", err)
	}
	configureFunc(policyClient.Client)

	policyDescriptionClient, err := policydescription.NewPolicyDescriptionClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PolicyDescription client: %+v", err)
	}
	configureFunc(policyDescriptionClient.Client)

	policyFragmentClient, err := policyfragment.NewPolicyFragmentClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PolicyFragment client: %+v", err)
	}
	configureFunc(policyFragmentClient.Client)

	policyRestrictionClient, err := policyrestriction.NewPolicyRestrictionClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PolicyRestriction client: %+v", err)
	}
	configureFunc(policyRestrictionClient.Client)

	policyRestrictionsClient, err := policyrestrictions.NewPolicyRestrictionsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PolicyRestrictions client: %+v", err)
	}
	configureFunc(policyRestrictionsClient.Client)

	policyRestrictionsValidationsClient, err := policyrestrictionsvalidations.NewPolicyRestrictionsValidationsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PolicyRestrictionsValidations client: %+v", err)
	}
	configureFunc(policyRestrictionsValidationsClient.Client)

	portalConfigClient, err := portalconfig.NewPortalConfigClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PortalConfig client: %+v", err)
	}
	configureFunc(portalConfigClient.Client)

	portalRevisionClient, err := portalrevision.NewPortalRevisionClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PortalRevision client: %+v", err)
	}
	configureFunc(portalRevisionClient.Client)

	portalSettingsClient, err := portalsettings.NewPortalSettingsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PortalSettings client: %+v", err)
	}
	configureFunc(portalSettingsClient.Client)

	privateEndpointConnectionsClient, err := privateendpointconnections.NewPrivateEndpointConnectionsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PrivateEndpointConnections client: %+v", err)
	}
	configureFunc(privateEndpointConnectionsClient.Client)

	productApiClient, err := productapi.NewProductApiClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ProductApi client: %+v", err)
	}
	configureFunc(productApiClient.Client)

	productApiLinkClient, err := productapilink.NewProductApiLinkClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ProductApiLink client: %+v", err)
	}
	configureFunc(productApiLinkClient.Client)

	productClient, err := product.NewProductClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Product client: %+v", err)
	}
	configureFunc(productClient.Client)

	productGroupClient, err := productgroup.NewProductGroupClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ProductGroup client: %+v", err)
	}
	configureFunc(productGroupClient.Client)

	productGroupLinkClient, err := productgrouplink.NewProductGroupLinkClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ProductGroupLink client: %+v", err)
	}
	configureFunc(productGroupLinkClient.Client)

	productPolicyClient, err := productpolicy.NewProductPolicyClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ProductPolicy client: %+v", err)
	}
	configureFunc(productPolicyClient.Client)

	productSubscriptionClient, err := productsubscription.NewProductSubscriptionClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ProductSubscription client: %+v", err)
	}
	configureFunc(productSubscriptionClient.Client)

	productTagClient, err := producttag.NewProductTagClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ProductTag client: %+v", err)
	}
	configureFunc(productTagClient.Client)

	productWikiClient, err := productwiki.NewProductWikiClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ProductWiki client: %+v", err)
	}
	configureFunc(productWikiClient.Client)

	productsByTagClient, err := productsbytag.NewProductsByTagClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ProductsByTag client: %+v", err)
	}
	configureFunc(productsByTagClient.Client)

	quotaByCounterKeysClient, err := quotabycounterkeys.NewQuotaByCounterKeysClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building QuotaByCounterKeys client: %+v", err)
	}
	configureFunc(quotaByCounterKeysClient.Client)

	quotaByPeriodKeysClient, err := quotabyperiodkeys.NewQuotaByPeriodKeysClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building QuotaByPeriodKeys client: %+v", err)
	}
	configureFunc(quotaByPeriodKeysClient.Client)

	regionClient, err := region.NewRegionClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Region client: %+v", err)
	}
	configureFunc(regionClient.Client)

	reportsClient, err := reports.NewReportsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Reports client: %+v", err)
	}
	configureFunc(reportsClient.Client)

	schemaClient, err := schema.NewSchemaClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Schema client: %+v", err)
	}
	configureFunc(schemaClient.Client)

	signInSettingsClient, err := signinsettings.NewSignInSettingsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SignInSettings client: %+v", err)
	}
	configureFunc(signInSettingsClient.Client)

	signUpSettingsClient, err := signupsettings.NewSignUpSettingsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SignUpSettings client: %+v", err)
	}
	configureFunc(signUpSettingsClient.Client)

	skusClient, err := skus.NewSkusClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Skus client: %+v", err)
	}
	configureFunc(skusClient.Client)

	subscriptionClient, err := subscription.NewSubscriptionClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Subscription client: %+v", err)
	}
	configureFunc(subscriptionClient.Client)

	tagApiLinkClient, err := tagapilink.NewTagApiLinkClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building TagApiLink client: %+v", err)
	}
	configureFunc(tagApiLinkClient.Client)

	tagClient, err := tag.NewTagClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Tag client: %+v", err)
	}
	configureFunc(tagClient.Client)

	tagOperationLinkClient, err := tagoperationlink.NewTagOperationLinkClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building TagOperationLink client: %+v", err)
	}
	configureFunc(tagOperationLinkClient.Client)

	tagProductLinkClient, err := tagproductlink.NewTagProductLinkClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building TagProductLink client: %+v", err)
	}
	configureFunc(tagProductLinkClient.Client)

	tagResourceClient, err := tagresource.NewTagResourceClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building TagResource client: %+v", err)
	}
	configureFunc(tagResourceClient.Client)

	tenantAccessClient, err := tenantaccess.NewTenantAccessClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building TenantAccess client: %+v", err)
	}
	configureFunc(tenantAccessClient.Client)

	tenantAccessGitClient, err := tenantaccessgit.NewTenantAccessGitClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building TenantAccessGit client: %+v", err)
	}
	configureFunc(tenantAccessGitClient.Client)

	tenantConfigurationClient, err := tenantconfiguration.NewTenantConfigurationClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building TenantConfiguration client: %+v", err)
	}
	configureFunc(tenantConfigurationClient.Client)

	tenantConfigurationSyncStateClient, err := tenantconfigurationsyncstate.NewTenantConfigurationSyncStateClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building TenantConfigurationSyncState client: %+v", err)
	}
	configureFunc(tenantConfigurationSyncStateClient.Client)

	tenantSettingsClient, err := tenantsettings.NewTenantSettingsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building TenantSettings client: %+v", err)
	}
	configureFunc(tenantSettingsClient.Client)

	userClient, err := user.NewUserClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building User client: %+v", err)
	}
	configureFunc(userClient.Client)

	userConfirmationPasswordSendClient, err := userconfirmationpasswordsend.NewUserConfirmationPasswordSendClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building UserConfirmationPasswordSend client: %+v", err)
	}
	configureFunc(userConfirmationPasswordSendClient.Client)

	userGroupClient, err := usergroup.NewUserGroupClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building UserGroup client: %+v", err)
	}
	configureFunc(userGroupClient.Client)

	userIdentityClient, err := useridentity.NewUserIdentityClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building UserIdentity client: %+v", err)
	}
	configureFunc(userIdentityClient.Client)

	userSubscriptionClient, err := usersubscription.NewUserSubscriptionClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building UserSubscription client: %+v", err)
	}
	configureFunc(userSubscriptionClient.Client)

	userTokenClient, err := usertoken.NewUserTokenClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building UserToken client: %+v", err)
	}
	configureFunc(userTokenClient.Client)

	usersClient, err := users.NewUsersClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Users client: %+v", err)
	}
	configureFunc(usersClient.Client)

	workspaceClient, err := workspace.NewWorkspaceClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Workspace client: %+v", err)
	}
	configureFunc(workspaceClient.Client)

	workspacePolicyClient, err := workspacepolicy.NewWorkspacePolicyClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building WorkspacePolicy client: %+v", err)
	}
	configureFunc(workspacePolicyClient.Client)

	return &Client{
		AllPolicies:                          allPoliciesClient,
		Api:                                  apiClient,
		ApiDiagnostic:                        apiDiagnosticClient,
		ApiGateway:                           apiGatewayClient,
		ApiGatewayConfigConnection:           apiGatewayConfigConnectionClient,
		ApiIssue:                             apiIssueClient,
		ApiIssueAttachment:                   apiIssueAttachmentClient,
		ApiIssueComment:                      apiIssueCommentClient,
		ApiManagementGatewaySkus:             apiManagementGatewaySkusClient,
		ApiManagementService:                 apiManagementServiceClient,
		ApiManagementServiceSkus:             apiManagementServiceSkusClient,
		ApiManagementWorkspaceLinks:          apiManagementWorkspaceLinksClient,
		ApiOperation:                         apiOperationClient,
		ApiOperationPolicy:                   apiOperationPolicyClient,
		ApiOperationTag:                      apiOperationTagClient,
		ApiOperationsByTag:                   apiOperationsByTagClient,
		ApiPolicy:                            apiPolicyClient,
		ApiProduct:                           apiProductClient,
		ApiRelease:                           apiReleaseClient,
		ApiRevision:                          apiRevisionClient,
		ApiSchema:                            apiSchemaClient,
		ApiTag:                               apiTagClient,
		ApiTagDescription:                    apiTagDescriptionClient,
		ApiVersionSet:                        apiVersionSetClient,
		ApiVersionSets:                       apiVersionSetsClient,
		ApiWiki:                              apiWikiClient,
		ApisByTag:                            apisByTagClient,
		Authorization:                        authorizationClient,
		AuthorizationAccessPolicy:            authorizationAccessPolicyClient,
		AuthorizationConfirmConsentCode:      authorizationConfirmConsentCodeClient,
		AuthorizationLoginLinks:              authorizationLoginLinksClient,
		AuthorizationProvider:                authorizationProviderClient,
		AuthorizationServer:                  authorizationServerClient,
		Authorizations:                       authorizationsClient,
		Backend:                              backendClient,
		BackendReconnect:                     backendReconnectClient,
		Cache:                                cacheClient,
		Certificate:                          certificateClient,
		ContentType:                          contentTypeClient,
		ContentTypeContentItem:               contentTypeContentItemClient,
		DelegationSettings:                   delegationSettingsClient,
		DeletedService:                       deletedServiceClient,
		Diagnostic:                           diagnosticClient,
		DocumentationResource:                documentationResourceClient,
		EmailTemplate:                        emailTemplateClient,
		EmailTemplates:                       emailTemplatesClient,
		Gateway:                              gatewayClient,
		GatewayApi:                           gatewayApiClient,
		GatewayCertificateAuthority:          gatewayCertificateAuthorityClient,
		GatewayGenerateToken:                 gatewayGenerateTokenClient,
		GatewayHostnameConfiguration:         gatewayHostnameConfigurationClient,
		GatewayInvalidateDebugCredentials:    gatewayInvalidateDebugCredentialsClient,
		GatewayListDebugCredentials:          gatewayListDebugCredentialsClient,
		GatewayListKeys:                      gatewayListKeysClient,
		GatewayListTrace:                     gatewayListTraceClient,
		GatewayRegenerateKey:                 gatewayRegenerateKeyClient,
		GraphQLApiResolver:                   graphQLApiResolverClient,
		GraphQLApiResolverPolicy:             graphQLApiResolverPolicyClient,
		Group:                                groupClient,
		GroupUser:                            groupUserClient,
		IdentityProvider:                     identityProviderClient,
		Issue:                                issueClient,
		Logger:                               loggerClient,
		NamedValue:                           namedValueClient,
		NetworkStatus:                        networkStatusClient,
		Notification:                         notificationClient,
		NotificationRecipientEmail:           notificationRecipientEmailClient,
		NotificationRecipientUser:            notificationRecipientUserClient,
		OpenidConnectProvider:                openidConnectProviderClient,
		OutboundNetworkDependenciesEndpoints: outboundNetworkDependenciesEndpointsClient,
		PerformConnectivityCheck:             performConnectivityCheckClient,
		Policy:                               policyClient,
		PolicyDescription:                    policyDescriptionClient,
		PolicyFragment:                       policyFragmentClient,
		PolicyRestriction:                    policyRestrictionClient,
		PolicyRestrictions:                   policyRestrictionsClient,
		PolicyRestrictionsValidations:        policyRestrictionsValidationsClient,
		PortalConfig:                         portalConfigClient,
		PortalRevision:                       portalRevisionClient,
		PortalSettings:                       portalSettingsClient,
		PrivateEndpointConnections:           privateEndpointConnectionsClient,
		Product:                              productClient,
		ProductApi:                           productApiClient,
		ProductApiLink:                       productApiLinkClient,
		ProductGroup:                         productGroupClient,
		ProductGroupLink:                     productGroupLinkClient,
		ProductPolicy:                        productPolicyClient,
		ProductSubscription:                  productSubscriptionClient,
		ProductTag:                           productTagClient,
		ProductWiki:                          productWikiClient,
		ProductsByTag:                        productsByTagClient,
		QuotaByCounterKeys:                   quotaByCounterKeysClient,
		QuotaByPeriodKeys:                    quotaByPeriodKeysClient,
		Region:                               regionClient,
		Reports:                              reportsClient,
		Schema:                               schemaClient,
		SignInSettings:                       signInSettingsClient,
		SignUpSettings:                       signUpSettingsClient,
		Skus:                                 skusClient,
		Subscription:                         subscriptionClient,
		Tag:                                  tagClient,
		TagApiLink:                           tagApiLinkClient,
		TagOperationLink:                     tagOperationLinkClient,
		TagProductLink:                       tagProductLinkClient,
		TagResource:                          tagResourceClient,
		TenantAccess:                         tenantAccessClient,
		TenantAccessGit:                      tenantAccessGitClient,
		TenantConfiguration:                  tenantConfigurationClient,
		TenantConfigurationSyncState:         tenantConfigurationSyncStateClient,
		TenantSettings:                       tenantSettingsClient,
		User:                                 userClient,
		UserConfirmationPasswordSend:         userConfirmationPasswordSendClient,
		UserGroup:                            userGroupClient,
		UserIdentity:                         userIdentityClient,
		UserSubscription:                     userSubscriptionClient,
		UserToken:                            userTokenClient,
		Users:                                usersClient,
		Workspace:                            workspaceClient,
		WorkspacePolicy:                      workspacePolicyClient,
	}, nil
}
