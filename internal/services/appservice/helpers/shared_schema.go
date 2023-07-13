// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/web/2022-09-01/web"
)

type IpRestriction struct {
	IpAddress    string                 `tfschema:"ip_address"`
	ServiceTag   string                 `tfschema:"service_tag"`
	VnetSubnetId string                 `tfschema:"virtual_network_subnet_id"`
	Name         string                 `tfschema:"name"`
	Priority     int                    `tfschema:"priority"`
	Action       string                 `tfschema:"action"`
	Headers      []IpRestrictionHeaders `tfschema:"headers"`
}

type IpRestrictionHeaders struct {
	XForwardedHost []string `tfschema:"x_forwarded_host"`
	XForwardedFor  []string `tfschema:"x_forwarded_for"`
	XAzureFDID     []string `tfschema:"x_azure_fdid"`
	XFDHealthProbe []string `tfschema:"x_fd_health_probe"`
}

func (v IpRestriction) Validate() error {
	hasIpAddress := v.IpAddress != ""
	hasServiceTag := v.ServiceTag != ""
	hasVnetSubnetId := v.VnetSubnetId != ""

	if (hasIpAddress && hasServiceTag) || (hasIpAddress && hasVnetSubnetId) || (hasServiceTag && hasVnetSubnetId) {
		return fmt.Errorf("only one of `ip_address`, `service_tag`, or `virtual_network_subnet_id` can be specified")
	}

	if !hasIpAddress && !hasServiceTag && !hasVnetSubnetId {
		return fmt.Errorf("one of `ip_address`, `service_tag`, or `virtual_network_subnet_id` must be specified")
	}

	return nil
}

func IpRestrictionSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"ip_address": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validate.IsIpOrCIDRRangeList,
					Description:  "The CIDR notation of the IP or IP Range to match. For example: `10.0.0.0/24` or `192.168.10.1/32` or `fe80::/64` or `13.107.6.152/31,13.107.128.0/22`",
				},

				"service_tag": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The Service Tag used for this IP Restriction.",
				},

				"virtual_network_subnet_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: commonids.ValidateSubnetID,
					Description:  "The Virtual Network Subnet ID used for this IP Restriction.",
				},

				"name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The name which should be used for this `ip_restriction`.",
				},

				"priority": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      65000,
					ValidateFunc: validation.IntBetween(1, 2147483647),
					Description:  "The priority value of this `ip_restriction`.",
				},

				"action": {
					Type:     pluginsdk.TypeString,
					Default:  "Allow",
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"Allow",
						"Deny",
					}, false),
					Description: "The action to take. Possible values are `Allow` or `Deny`.",
				},

				"headers": IpRestrictionHeadersSchema(),
			},
		},
	}
}

func IpRestrictionSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"ip_address": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The CIDR notation of the IP or IP Range to match.",
				},

				"service_tag": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The Service Tag used for this IP Restriction.",
				},

				"virtual_network_subnet_id": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The Virtual Network Subnet ID used for this IP Restriction.",
				},

				"name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The name used for this `ip_restriction`.",
				},

				"priority": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "The priority value of this `ip_restriction`.",
				},

				"action": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The action to take.",
				},

				"headers": IpRestrictionHeadersSchemaComputed(),
			},
		},
	}
}

func IpRestrictionHeadersSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:       pluginsdk.TypeList,
		MaxItems:   1,
		Optional:   true,
		ConfigMode: pluginsdk.SchemaConfigModeAttr,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"x_forwarded_host": {
					Type:     pluginsdk.TypeList,
					MaxItems: 8,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "Specifies a list of Hosts for which matching should be applied.",
				},

				"x_forwarded_for": {
					Type:     pluginsdk.TypeList,
					MaxItems: 8,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.IsCIDR,
					},
					Description: "Specifies a list of addresses for which matching should be applied. Omitting this value means allow any.",
				},

				"x_azure_fdid": { // Front Door ID (UUID)
					Type:     pluginsdk.TypeList,
					MaxItems: 8,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.IsUUID,
					},
					Description: "Specifies a list of Azure Front Door IDs.",
				},

				"x_fd_health_probe": { // 1 or absent
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
						ValidateFunc: validation.StringInSlice([]string{
							"1",
						}, false),
					},
				},
			},
		},
	}
}

func IpRestrictionHeadersSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"x_forwarded_host": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "The list of Hosts for which matching will be applied.",
				},

				"x_forwarded_for": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "The list of addresses for which matching is applied.",
				},

				"x_azure_fdid": { // Front Door ID (UUID)
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "The list of Azure Front Door IDs.",
				},

				"x_fd_health_probe": { // 1 or absent
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "Specifies if a Front Door Health Probe is expected.",
				},
			},
		},
	}
}

type CorsSetting struct {
	AllowedOrigins     []string `tfschema:"allowed_origins"`
	SupportCredentials bool     `tfschema:"support_credentials"`
}

func CorsSettingsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"allowed_origins": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					MinItems: 1,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "Specifies a list of origins that should be allowed to make cross-origin calls.",
				},

				"support_credentials": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Are credentials allowed in CORS requests? Defaults to `false`.",
				},
			},
		},
	}
}

func CorsSettingsSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"allowed_origins": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "The list of origins that are allowed to make cross-origin calls.",
				},

				"support_credentials": {
					Type:        pluginsdk.TypeBool,
					Computed:    true,
					Description: "Are credentials allowed in CORS requests?",
				},
			},
		},
	}
}

func FlattenCorsSettings(input *web.CorsSettings) []CorsSetting {
	if input == nil {
		return []CorsSetting{}
	}

	cors := *input
	if len(pointer.From(cors.AllowedOrigins)) == 0 && !pointer.From(cors.SupportCredentials) {
		return []CorsSetting{}
	}

	return []CorsSetting{{
		SupportCredentials: pointer.From(cors.SupportCredentials),
		AllowedOrigins:     pointer.From(cors.AllowedOrigins),
	}}
}

func ExpandCorsSettings(input []CorsSetting) *web.CorsSettings {
	if len(input) != 1 {
		return &web.CorsSettings{}
	}
	cors := input[0]

	return &web.CorsSettings{
		AllowedOrigins:     pointer.To(cors.AllowedOrigins),
		SupportCredentials: pointer.To(cors.SupportCredentials),
	}
}

type SourceControl struct {
	RepoURL           string `tfschema:"repo_url"`
	Branch            string `tfschema:"branch"`
	ManualIntegration bool   `tfschema:"manual_integration"`
	UseMercurial      bool   `tfschema:"use_mercurial"`
	RollbackEnabled   bool   `tfschema:"rollback_enabled"`
}

type SiteCredential struct {
	Username string `tfschema:"name"`
	Password string `tfschema:"password"`
}

func SiteCredentialSchema() *pluginsdk.Schema { // TODO - This can apparently be disabled as a security option for the service?
	return &pluginsdk.Schema{
		Type:      pluginsdk.TypeList,
		Computed:  true,
		Sensitive: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Sensitive:   true,
					Description: "The Site Credentials Username used for publishing.",
				},

				"password": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Sensitive:   true,
					Description: "The Site Credentials Password used for publishing.",
				},
			},
		},
	}
}

type AuthSettings struct {
	Enabled                     bool                    `tfschema:"enabled"`
	AdditionalLoginParameters   map[string]string       `tfschema:"additional_login_parameters"`
	AllowedExternalRedirectUrls []string                `tfschema:"allowed_external_redirect_urls"`
	DefaultProvider             string                  `tfschema:"default_provider"`
	Issuer                      string                  `tfschema:"issuer"`
	RuntimeVersion              string                  `tfschema:"runtime_version"`
	TokenRefreshExtensionHours  float64                 `tfschema:"token_refresh_extension_hours"`
	TokenStoreEnabled           bool                    `tfschema:"token_store_enabled"`
	UnauthenticatedClientAction string                  `tfschema:"unauthenticated_client_action"`
	AzureActiveDirectoryAuth    []AadAuthSettings       `tfschema:"active_directory"`
	FacebookAuth                []FacebookAuthSettings  `tfschema:"facebook"`
	GithubAuth                  []GithubAuthSettings    `tfschema:"github"`
	GoogleAuth                  []GoogleAuthSettings    `tfschema:"google"`
	MicrosoftAuth               []MicrosoftAuthSettings `tfschema:"microsoft"`
	TwitterAuth                 []TwitterAuthSettings   `tfschema:"twitter"`
}

func AuthSettingsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"enabled": {
					Type:        pluginsdk.TypeBool,
					Required:    true,
					Description: "Should the Authentication / Authorization feature be enabled?",
				},

				"additional_login_parameters": {
					Type:     pluginsdk.TypeMap,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "Specifies a map of Login Parameters to send to the OpenID Connect authorization endpoint when a user logs in.",
				},

				"allowed_external_redirect_urls": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					Description: "Specifies a list of External URLs that can be redirected to as part of logging in or logging out of the Windows Web App.",
				},

				"default_provider": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true, // Once set, cannot be unset
					ValidateFunc: validation.StringInSlice([]string{
						string(web.BuiltInAuthenticationProviderAzureActiveDirectory),
						string(web.BuiltInAuthenticationProviderFacebook),
						string(web.BuiltInAuthenticationProviderGithub),
						string(web.BuiltInAuthenticationProviderGoogle),
						string(web.BuiltInAuthenticationProviderMicrosoftAccount),
						string(web.BuiltInAuthenticationProviderTwitter),
					}, false),
					Description: "The default authentication provider to use when multiple providers are configured. Possible values include: `AzureActiveDirectory`, `Facebook`, `Google`, `MicrosoftAccount`, `Twitter`, `Github`.",
				},

				"issuer": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
					Description:  "The OpenID Connect Issuer URI that represents the entity which issues access tokens.",
				},

				"runtime_version": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Computed:    true,
					Description: "The RuntimeVersion of the Authentication / Authorization feature in use.",
				},

				"token_refresh_extension_hours": {
					Type:        pluginsdk.TypeFloat,
					Optional:    true,
					Default:     72,
					Description: "The number of hours after session token expiration that a session token can be used to call the token refresh API. Defaults to `72` hours.",
				},

				"token_store_enabled": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Should the Windows Web App durably store platform-specific security tokens that are obtained during login flows? Defaults to `false`.",
				},

				"unauthenticated_client_action": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true, // Once set, cannot be removed
					ValidateFunc: validation.StringInSlice([]string{
						string(web.UnauthenticatedClientActionAllowAnonymous),
						string(web.UnauthenticatedClientActionRedirectToLoginPage),
					}, false),
					Description: "The action to take when an unauthenticated client attempts to access the app. Possible values include: `RedirectToLoginPage`, `AllowAnonymous`.",
				},

				"active_directory": AadAuthSettingsSchema(),

				"facebook": FacebookAuthSettingsSchema(),

				"github": GithubAuthSettingsSchema(),

				"google": GoogleAuthSettingsSchema(),

				"microsoft": MicrosoftAuthSettingsSchema(),

				"twitter": TwitterAuthSettingsSchema(),
			},
		},
	}
}

func AuthSettingsSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"enabled": {
					Type:        pluginsdk.TypeBool,
					Computed:    true,
					Description: "Is the Authentication / Authorization feature enabled?",
				},

				"additional_login_parameters": {
					Type:     pluginsdk.TypeMap,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "The map of Login Parameters sent to the OpenID Connect authorization endpoint when a user logs in.",
				},

				"allowed_external_redirect_urls": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "Specifies a list of External URLs that can be redirected to as part of logging in or logging out of the Windows Web App.",
				},

				"default_provider": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The default authentication provider used when multiple providers are configured. Possible values include: `AzureActiveDirectory`, `Facebook`, `Google`, `MicrosoftAccount`, `Twitter`, `Github`.",
				},

				"issuer": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The OpenID Connect Issuer URI that represents the entity which issues access tokens.",
				},

				"runtime_version": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The RuntimeVersion of the Authentication / Authorization feature in use.",
				},

				"token_refresh_extension_hours": {
					Type:        pluginsdk.TypeFloat,
					Computed:    true,
					Description: "The number of hours after session token expiration that a session token can be used to call the token refresh API.",
				},

				"token_store_enabled": {
					Type:        pluginsdk.TypeBool,
					Computed:    true,
					Description: "Are platform-specific security tokens that are obtained during login flows durably stored?",
				},

				"unauthenticated_client_action": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The action taken when an unauthenticated client attempts to access the app.",
				},

				"active_directory": AadAuthSettingsSchemaComputed(),

				"facebook": FacebookAuthSettingsSchemaComputed(),

				"github": GithubAuthSettingsSchemaComputed(),

				"google": GoogleAuthSettingsSchemaComputed(),

				"microsoft": MicrosoftAuthSettingsSchemaComputed(),

				"twitter": TwitterAuthSettingsSchemaComputed(),
			},
		},
	}
}

type AadAuthSettings struct {
	ClientId                string   `tfschema:"client_id"`
	ClientSecret            string   `tfschema:"client_secret"`
	ClientSecretSettingName string   `tfschema:"client_secret_setting_name"`
	AllowedAudiences        []string `tfschema:"allowed_audiences"`
}

func AadAuthSettingsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"client_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The ID of the Client to use to authenticate with Azure Active Directory.",
				},

				"client_secret": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
					ConflictsWith: []string{
						"auth_settings.0.active_directory.0.client_secret_setting_name",
					},
					Description: "The Client Secret for the Client ID. Cannot be used with `client_secret_setting_name`.",
				},

				"client_secret_setting_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					ConflictsWith: []string{
						"auth_settings.0.active_directory.0.client_secret",
					},
					Description: "The App Setting name that contains the client secret of the Client. Cannot be used with `client_secret`.",
				},

				"allowed_audiences": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "Specifies a list of Allowed audience values to consider when validating JWTs issued by Azure Active Directory.",
				},
			},
		},
	}
}

func AadAuthSettingsSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"client_id": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The ID of the Client to use to authenticate with Azure Active Directory.",
				},

				"client_secret": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Sensitive:   true,
					Description: "The Client Secret for the Client ID.",
				},

				"client_secret_setting_name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The App Setting name that contains the client secret of the Client.",
				},

				"allowed_audiences": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "The list of Allowed audience values considered when validating JWTs issued by Azure Active Directory.",
				},
			},
		},
	}
}

type FacebookAuthSettings struct {
	AppId                string   `tfschema:"app_id"`
	AppSecret            string   `tfschema:"app_secret"`
	AppSecretSettingName string   `tfschema:"app_secret_setting_name"`
	OauthScopes          []string `tfschema:"oauth_scopes"`
}

func FacebookAuthSettingsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"app_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The App ID of the Facebook app used for login.",
				},

				"app_secret": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
					ExactlyOneOf: []string{
						"auth_settings.0.facebook.0.app_secret",
						"auth_settings.0.facebook.0.app_secret_setting_name",
					},
					Description: "The App Secret of the Facebook app used for Facebook Login. Cannot be specified with `app_secret_setting_name`.",
				},

				"app_secret_setting_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					ExactlyOneOf: []string{
						"auth_settings.0.facebook.0.app_secret",
						"auth_settings.0.facebook.0.app_secret_setting_name",
					},
					Description: "The app setting name that contains the `app_secret` value used for Facebook Login. Cannot be specified with `app_secret`.",
				},

				"oauth_scopes": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "Specifies a list of OAuth 2.0 scopes to be requested as part of Facebook Login authentication.",
				},
			},
		},
	}
}

func FacebookAuthSettingsSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"app_id": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The App ID of the Facebook app used for login.",
				},

				"app_secret": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Sensitive:   true,
					Description: "The App Secret of the Facebook app used for Facebook Login.",
				},

				"app_secret_setting_name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The app setting name that contains the `app_secret` value used for Facebook Login.",
				},

				"oauth_scopes": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "The list of OAuth 2.0 scopes requested as part of Facebook Login authentication.",
				},
			},
		},
	}
}

type GoogleAuthSettings struct {
	ClientId                string   `tfschema:"client_id"`
	ClientSecret            string   `tfschema:"client_secret"`
	ClientSecretSettingName string   `tfschema:"client_secret_setting_name"`
	OauthScopes             []string `tfschema:"oauth_scopes"`
}

func GoogleAuthSettingsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"client_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The OpenID Connect Client ID for the Google web application.",
				},

				"client_secret": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
					ExactlyOneOf: []string{
						"auth_settings.0.google.0.client_secret",
						"auth_settings.0.google.0.client_secret_setting_name",
					},
					Description: "The client secret associated with the Google web application.  Cannot be specified with `client_secret_setting_name`.",
				},

				"client_secret_setting_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					ExactlyOneOf: []string{
						"auth_settings.0.google.0.client_secret",
						"auth_settings.0.google.0.client_secret_setting_name",
					},
					Description: "The app setting name that contains the `client_secret` value used for Google Login. Cannot be specified with `client_secret`.",
				},

				"oauth_scopes": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "Specifies a list of OAuth 2.0 scopes that will be requested as part of Google Sign-In authentication. If not specified, \"openid\", \"profile\", and \"email\" are used as default scopes.",
				},
			},
		},
	}
}

func GoogleAuthSettingsSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"client_id": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The OpenID Connect Client ID for the Google web application.",
				},

				"client_secret": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Sensitive:   true,
					Description: "The client secret associated with the Google web application.",
				},

				"client_secret_setting_name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The app setting name that contains the `client_secret` value used for Google Login.",
				},

				"oauth_scopes": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "The list of OAuth 2.0 scopes that requested as part of Google Sign-In authentication.",
				},
			},
		},
	}
}

type MicrosoftAuthSettings struct {
	ClientId                string   `tfschema:"client_id"`
	ClientSecret            string   `tfschema:"client_secret"`
	ClientSecretSettingName string   `tfschema:"client_secret_setting_name"`
	OauthScopes             []string `tfschema:"oauth_scopes"`
}

func MicrosoftAuthSettingsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"client_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The OAuth 2.0 client ID that was created for the app used for authentication.",
				},

				"client_secret": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
					ExactlyOneOf: []string{
						"auth_settings.0.microsoft.0.client_secret",
						"auth_settings.0.microsoft.0.client_secret_setting_name",
					},
					Description: "The OAuth 2.0 client secret that was created for the app used for authentication. Cannot be specified with `client_secret_setting_name`.",
				},

				"client_secret_setting_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					ExactlyOneOf: []string{
						"auth_settings.0.microsoft.0.client_secret",
						"auth_settings.0.microsoft.0.client_secret_setting_name",
					},
					Description: "The app setting name containing the OAuth 2.0 client secret that was created for the app used for authentication. Cannot be specified with `client_secret`.",
				},

				"oauth_scopes": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "The list of OAuth 2.0 scopes that will be requested as part of Microsoft Account authentication. If not specified, `wl.basic` is used as the default scope.",
				},
			},
		},
	}
}

func MicrosoftAuthSettingsSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"client_id": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The OAuth 2.0 client ID that was created for the app used for authentication.",
				},

				"client_secret": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Sensitive:   true,
					Description: "The OAuth 2.0 client secret that was created for the app used for authentication.",
				},

				"client_secret_setting_name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The app setting name containing the OAuth 2.0 client secret that was created for the app used for authentication.",
				},

				"oauth_scopes": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "The list of OAuth 2.0 scopes requested as part of Microsoft Account authentication.",
				},
			},
		},
	}
}

type TwitterAuthSettings struct {
	ConsumerKey               string `tfschema:"consumer_key"`
	ConsumerSecret            string `tfschema:"consumer_secret"`
	ConsumerSecretSettingName string `tfschema:"consumer_secret_setting_name"`
}

func TwitterAuthSettingsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"consumer_key": {
					Type:        pluginsdk.TypeString,
					Required:    true,
					Description: "The OAuth 1.0a consumer key of the Twitter application used for sign-in.",
				},

				"consumer_secret": {
					Type:      pluginsdk.TypeString,
					Optional:  true,
					Sensitive: true,
					ExactlyOneOf: []string{
						"auth_settings.0.twitter.0.consumer_secret",
						"auth_settings.0.twitter.0.consumer_secret_setting_name",
					},
					Description: "The OAuth 1.0a consumer secret of the Twitter application used for sign-in. Cannot be specified with `consumer_secret_setting_name`.",
				},

				"consumer_secret_setting_name": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Description: "The app setting name that contains the OAuth 1.0a consumer secret of the Twitter application used for sign-in. Cannot be specified with `consumer_secret`.",
				},
			},
		},
	}
}

func TwitterAuthSettingsSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"consumer_key": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The OAuth 1.0a consumer key of the Twitter application used for sign-in.",
				},

				"consumer_secret": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Sensitive:   true,
					Description: "The OAuth 1.0a consumer secret of the Twitter application used for sign-in.",
				},

				"consumer_secret_setting_name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The app setting name that contains the OAuth 1.0a consumer secret of the Twitter application used for sign-in.",
				},
			},
		},
	}
}

type GithubAuthSettings struct {
	ClientId                string   `tfschema:"client_id"`
	ClientSecret            string   `tfschema:"client_secret"`
	ClientSecretSettingName string   `tfschema:"client_secret_setting_name"`
	OAuthScopes             []string `tfschema:"oauth_scopes"`
}

func GithubAuthSettingsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"client_id": {
					Type:        pluginsdk.TypeString,
					Required:    true,
					Description: "The ID of the GitHub app used for login.",
				},

				"client_secret": {
					Type:      pluginsdk.TypeString,
					Optional:  true,
					Sensitive: true,
					ExactlyOneOf: []string{
						"auth_settings.0.github.0.client_secret",
						"auth_settings.0.github.0.client_secret_setting_name",
					},
					Description: "The Client Secret of the GitHub app used for GitHub Login. Cannot be specified with `client_secret_setting_name`.",
				},

				"client_secret_setting_name": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ExactlyOneOf: []string{
						"auth_settings.0.github.0.client_secret",
						"auth_settings.0.github.0.client_secret_setting_name",
					},
					Description: "The app setting name that contains the `client_secret` value used for GitHub Login. Cannot be specified with `client_secret`.",
				},

				"oauth_scopes": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "Specifies a list of OAuth 2.0 scopes that will be requested as part of GitHub Login authentication.",
				},
			},
		},
	}
}

func GithubAuthSettingsSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"client_id": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The ID of the GitHub app used for login.",
				},

				"client_secret": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Sensitive:   true,
					Description: "The Client Secret of the GitHub app used for GitHub Login.",
				},

				"client_secret_setting_name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The app setting name that contains the `client_secret` value used for GitHub Login.",
				},

				"oauth_scopes": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "The list of OAuth 2.0 scopes requested as part of GitHub Login authentication.",
				},
			},
		},
	}
}

func ExpandIpRestrictions(restrictions []IpRestriction) (*[]web.IPSecurityRestriction, error) {
	expanded := make([]web.IPSecurityRestriction, 0)
	if len(restrictions) == 0 {
		return &expanded, nil
	}

	for _, v := range restrictions {
		if err := v.Validate(); err != nil {
			return nil, err
		}

		var restriction web.IPSecurityRestriction
		if v.Name != "" {
			restriction.Name = utils.String(v.Name)
		}

		if v.IpAddress != "" {
			restriction.IPAddress = utils.String(v.IpAddress)
		}

		if v.ServiceTag != "" {
			restriction.IPAddress = utils.String(v.ServiceTag)
			restriction.Tag = web.IPFilterTagServiceTag
		}

		if v.VnetSubnetId != "" {
			restriction.VnetSubnetResourceID = utils.String(v.VnetSubnetId)
		}

		restriction.Priority = utils.Int32(int32(v.Priority))

		restriction.Action = utils.String(v.Action)

		restriction.Headers = expandIpRestrictionHeaders(v.Headers)

		expanded = append(expanded, restriction)
	}

	return &expanded, nil
}

func expandIpRestrictionHeaders(headers []IpRestrictionHeaders) map[string][]string {
	result := make(map[string][]string)
	if len(headers) == 0 {
		return result
	}

	for _, v := range headers {
		if len(v.XForwardedHost) > 0 {
			result["x-forwarded-host"] = v.XForwardedHost
		}
		if len(v.XForwardedFor) > 0 {
			result["x-forwarded-for"] = v.XForwardedFor
		}
		if len(v.XAzureFDID) > 0 {
			result["x-azure-fdid"] = v.XAzureFDID
		}
		if len(v.XFDHealthProbe) > 0 {
			result["x-fd-healthprobe"] = v.XFDHealthProbe
		}
	}

	return result
}

func ExpandAuthSettings(auth []AuthSettings) *web.SiteAuthSettings {
	result := &web.SiteAuthSettings{}
	if len(auth) == 0 {
		return result
	}

	props := &web.SiteAuthSettingsProperties{}

	v := auth[0]

	props.Enabled = utils.Bool(v.Enabled)

	additionalLoginParams := make([]string, 0)
	if len(v.AdditionalLoginParameters) > 0 {
		for k, s := range v.AdditionalLoginParameters {
			additionalLoginParams = append(additionalLoginParams, fmt.Sprintf("%s=%s", k, s))
		}
		props.AdditionalLoginParams = &additionalLoginParams
	}

	props.AllowedExternalRedirectUrls = &v.AllowedExternalRedirectUrls

	props.DefaultProvider = web.BuiltInAuthenticationProvider(v.DefaultProvider)

	props.Issuer = utils.String(v.Issuer)

	props.RuntimeVersion = utils.String(v.RuntimeVersion)

	props.TokenStoreEnabled = utils.Bool(v.TokenStoreEnabled)

	props.TokenRefreshExtensionHours = utils.Float(v.TokenRefreshExtensionHours)

	props.UnauthenticatedClientAction = web.UnauthenticatedClientAction(v.UnauthenticatedClientAction)

	a := AadAuthSettings{}
	if len(v.AzureActiveDirectoryAuth) > 0 {
		a = v.AzureActiveDirectoryAuth[0]
	}
	props.ClientID = utils.String(a.ClientId)

	if a.ClientSecret != "" {
		props.ClientSecret = utils.String(a.ClientSecret)
	}

	if a.ClientSecretSettingName != "" {
		props.ClientSecretSettingName = utils.String(a.ClientSecretSettingName)
	}

	props.AllowedAudiences = &a.AllowedAudiences

	f := FacebookAuthSettings{}
	if len(v.FacebookAuth) > 0 {
		f = v.FacebookAuth[0]
	}
	props.FacebookAppID = utils.String(f.AppId)
	props.FacebookAppSecret = utils.String(f.AppSecret)
	props.FacebookAppSecretSettingName = utils.String(f.AppSecretSettingName)
	props.FacebookOAuthScopes = &f.OauthScopes

	gh := GithubAuthSettings{}
	if len(v.GithubAuth) > 0 {
		gh = v.GithubAuth[0]
	}
	props.GitHubClientID = utils.String(gh.ClientId)
	props.GitHubClientSecret = utils.String(gh.ClientSecret)
	props.GitHubClientSecretSettingName = utils.String(gh.ClientSecretSettingName)
	props.GitHubOAuthScopes = &gh.OAuthScopes

	g := GoogleAuthSettings{}
	if len(v.GoogleAuth) > 0 {
		g = v.GoogleAuth[0]
	}

	props.GoogleClientID = utils.String(g.ClientId)
	props.GoogleClientSecret = utils.String(g.ClientSecret)
	props.GoogleClientSecretSettingName = utils.String(g.ClientSecretSettingName)
	props.GoogleOAuthScopes = &g.OauthScopes

	m := MicrosoftAuthSettings{}
	if len(v.MicrosoftAuth) > 0 {
		m = v.MicrosoftAuth[0]
	}
	props.MicrosoftAccountClientID = utils.String(m.ClientId)
	props.MicrosoftAccountClientSecret = utils.String(m.ClientSecret)
	props.MicrosoftAccountClientSecretSettingName = utils.String(m.ClientSecretSettingName)
	props.MicrosoftAccountOAuthScopes = &m.OauthScopes

	t := TwitterAuthSettings{}
	if len(v.TwitterAuth) > 0 {
		t = v.TwitterAuth[0]
	}
	props.TwitterConsumerKey = utils.String(t.ConsumerKey)
	props.TwitterConsumerSecret = utils.String(t.ConsumerSecret)
	props.TwitterConsumerSecretSettingName = utils.String(t.ConsumerSecretSettingName)

	result.SiteAuthSettingsProperties = props

	return result
}

func FlattenAuthSettings(auth web.SiteAuthSettings) []AuthSettings {
	if auth.SiteAuthSettingsProperties == nil || !pointer.From(auth.Enabled) || strings.ToLower(pointer.From(auth.ConfigVersion)) != "v1" {
		return []AuthSettings{}
	}

	props := *auth.SiteAuthSettingsProperties

	result := AuthSettings{
		DefaultProvider:             string(props.DefaultProvider),
		UnauthenticatedClientAction: string(props.UnauthenticatedClientAction),
	}

	if props.Enabled != nil {
		result.Enabled = *props.Enabled
	}

	if props.AdditionalLoginParams != nil {
		params := make(map[string]string)
		for _, v := range *props.AdditionalLoginParams {
			parts := strings.Split(v, "=")
			if len(parts) != 2 {
				continue
			}
			params[parts[0]] = parts[1]
		}
		result.AdditionalLoginParameters = params
	}

	var allowedRedirectUrls []string
	if props.AllowedExternalRedirectUrls != nil {
		allowedRedirectUrls = *props.AllowedExternalRedirectUrls
	}
	result.AllowedExternalRedirectUrls = allowedRedirectUrls

	if props.Issuer != nil {
		result.Issuer = *props.Issuer
	}

	if props.RuntimeVersion != nil {
		result.RuntimeVersion = *props.RuntimeVersion
	}

	if props.TokenRefreshExtensionHours != nil {
		result.TokenRefreshExtensionHours = *props.TokenRefreshExtensionHours
	}

	if props.TokenStoreEnabled != nil {
		result.TokenStoreEnabled = *props.TokenStoreEnabled
	}

	// AAD Auth
	if props.ClientID != nil {
		aadAuthSettings := AadAuthSettings{
			ClientId: *props.ClientID,
		}

		if props.ClientSecret != nil {
			aadAuthSettings.ClientSecret = *props.ClientSecret
		}

		if props.ClientSecretSettingName != nil {
			aadAuthSettings.ClientSecretSettingName = *props.ClientSecretSettingName
		}

		if props.AllowedAudiences != nil {
			aadAuthSettings.AllowedAudiences = *props.AllowedAudiences
		}

		result.AzureActiveDirectoryAuth = []AadAuthSettings{aadAuthSettings}
	}

	if props.FacebookAppID != nil {
		facebookAuthSettings := FacebookAuthSettings{
			AppId: *props.FacebookAppID,
		}

		if props.FacebookAppSecret != nil {
			facebookAuthSettings.AppSecret = *props.FacebookAppSecret
		}

		if props.FacebookAppSecretSettingName != nil {
			facebookAuthSettings.AppSecretSettingName = *props.FacebookAppSecretSettingName
		}

		if props.FacebookOAuthScopes != nil {
			facebookAuthSettings.OauthScopes = *props.FacebookOAuthScopes
		}

		result.FacebookAuth = []FacebookAuthSettings{facebookAuthSettings}
	}

	if props.GitHubClientID != nil {
		githubAuthSetting := GithubAuthSettings{
			ClientId: *props.GitHubClientID,
		}

		if props.GitHubClientSecret != nil {
			githubAuthSetting.ClientSecret = *props.GitHubClientSecret
		}

		if props.GitHubClientSecretSettingName != nil {
			githubAuthSetting.ClientSecretSettingName = *props.GitHubClientSecretSettingName
		}

		result.GithubAuth = []GithubAuthSettings{githubAuthSetting}
	}

	if props.GoogleClientID != nil {
		googleAuthSettings := GoogleAuthSettings{
			ClientId: *props.GoogleClientID,
		}

		if props.GoogleClientSecret != nil {
			googleAuthSettings.ClientSecret = *props.GoogleClientSecret
		}

		if props.GoogleClientSecretSettingName != nil {
			googleAuthSettings.ClientSecretSettingName = *props.GoogleClientSecretSettingName
		}

		if props.GoogleOAuthScopes != nil {
			googleAuthSettings.OauthScopes = *props.GoogleOAuthScopes
		}

		result.GoogleAuth = []GoogleAuthSettings{googleAuthSettings}
	}

	if props.MicrosoftAccountClientID != nil {
		microsoftAuthSettings := MicrosoftAuthSettings{
			ClientId: *props.MicrosoftAccountClientID,
		}

		if props.MicrosoftAccountClientSecret != nil {
			microsoftAuthSettings.ClientSecret = *props.MicrosoftAccountClientSecret
		}

		if props.MicrosoftAccountClientSecretSettingName != nil {
			microsoftAuthSettings.ClientSecretSettingName = *props.MicrosoftAccountClientSecretSettingName
		}

		if props.MicrosoftAccountOAuthScopes != nil {
			microsoftAuthSettings.OauthScopes = *props.MicrosoftAccountOAuthScopes
		}

		result.MicrosoftAuth = []MicrosoftAuthSettings{microsoftAuthSettings}
	}

	if props.TwitterConsumerKey != nil {
		twitterAuthSetting := TwitterAuthSettings{
			ConsumerKey: *props.TwitterConsumerKey,
		}
		if props.TwitterConsumerSecret != nil {
			twitterAuthSetting.ConsumerSecret = *props.TwitterConsumerSecret
		}
		if props.TwitterConsumerSecretSettingName != nil {
			twitterAuthSetting.ConsumerSecretSettingName = *props.TwitterConsumerSecretSettingName
		}

		result.TwitterAuth = []TwitterAuthSettings{twitterAuthSetting}
	}

	return []AuthSettings{result}
}

func FlattenIpRestrictions(ipRestrictionsList *[]web.IPSecurityRestriction) []IpRestriction {
	if ipRestrictionsList == nil {
		return []IpRestriction{}
	}

	var ipRestrictions []IpRestriction
	for _, v := range *ipRestrictionsList {
		ipRestriction := IpRestriction{}

		if v.Name != nil {
			ipRestriction.Name = *v.Name
		}

		if v.IPAddress != nil {
			if *v.IPAddress == "Any" {
				continue
			}

			if v.Tag == web.IPFilterTagServiceTag {
				ipRestriction.ServiceTag = *v.IPAddress
			} else {
				ipRestriction.IpAddress = *v.IPAddress
			}
		}

		if v.VnetSubnetResourceID != nil {
			ipRestriction.VnetSubnetId = *v.VnetSubnetResourceID
		}

		if v.Priority != nil {
			ipRestriction.Priority = int(*v.Priority)
		}

		if v.Action != nil {
			ipRestriction.Action = *v.Action
		}

		ipRestriction.Headers = flattenIpRestrictionHeaders(v.Headers)

		ipRestrictions = append(ipRestrictions, ipRestriction)
	}

	return ipRestrictions
}

func flattenIpRestrictionHeaders(headers map[string][]string) []IpRestrictionHeaders {
	if len(headers) == 0 {
		return []IpRestrictionHeaders{}
	}
	ipRestrictionHeader := IpRestrictionHeaders{}
	if xForwardFor, ok := headers["x-forwarded-for"]; ok {
		ipRestrictionHeader.XForwardedFor = xForwardFor
	}

	if xForwardedHost, ok := headers["x-forwarded-host"]; ok {
		ipRestrictionHeader.XForwardedHost = xForwardedHost
	}

	if xAzureFDID, ok := headers["x-azure-fdid"]; ok {
		ipRestrictionHeader.XAzureFDID = xAzureFDID
	}

	if xFDHealthProbe, ok := headers["x-fd-healthprobe"]; ok {
		ipRestrictionHeader.XFDHealthProbe = xFDHealthProbe
	}

	return []IpRestrictionHeaders{ipRestrictionHeader}
}

func FlattenWebStringDictionary(input web.StringDictionary) map[string]string {
	result := make(map[string]string)
	for k, v := range input.Properties {
		result[k] = utils.NormalizeNilableString(v)
	}

	return result
}

func FlattenSiteCredentials(input web.User) []SiteCredential {
	var result []SiteCredential
	if input.UserProperties == nil {
		return result
	}

	userProps := *input.UserProperties
	result = append(result, SiteCredential{
		Username: utils.NormalizeNilableString(userProps.PublishingUserName),
		Password: utils.NormalizeNilableString(userProps.PublishingPassword),
	})

	return result
}

type StickySettings struct {
	AppSettingNames       []string `tfschema:"app_setting_names"`
	ConnectionStringNames []string `tfschema:"connection_string_names"`
}

func StickySettingsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"app_setting_names": {
					Type:     pluginsdk.TypeList,
					MinItems: 1,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					AtLeastOneOf: []string{
						"sticky_settings.0.app_setting_names",
						"sticky_settings.0.connection_string_names",
					},
				},

				"connection_string_names": {
					Type:     pluginsdk.TypeList,
					MinItems: 1,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					AtLeastOneOf: []string{
						"sticky_settings.0.app_setting_names",
						"sticky_settings.0.connection_string_names",
					},
				},
			},
		},
	}
}

func StickySettingsComputedSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"app_setting_names": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"connection_string_names": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},
			},
		},
	}
}

func ExpandStickySettings(input []StickySettings) *web.SlotConfigNames {
	if len(input) == 0 {
		return nil
	}

	return &web.SlotConfigNames{
		AppSettingNames:       &input[0].AppSettingNames,
		ConnectionStringNames: &input[0].ConnectionStringNames,
	}
}

func FlattenStickySettings(input *web.SlotConfigNames) []StickySettings {
	result := StickySettings{}
	if input == nil || (input.AppSettingNames == nil && input.ConnectionStringNames == nil) || (len(*input.AppSettingNames) == 0 && len(*input.ConnectionStringNames) == 0) {
		return []StickySettings{}
	}

	if input.AppSettingNames != nil && len(*input.AppSettingNames) > 0 {
		result.AppSettingNames = *input.AppSettingNames
	}

	if input.ConnectionStringNames != nil && len(*input.ConnectionStringNames) > 0 {
		result.ConnectionStringNames = *input.ConnectionStringNames
	}

	return []StickySettings{result}
}
