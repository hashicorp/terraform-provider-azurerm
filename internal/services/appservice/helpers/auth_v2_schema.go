package helpers

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"strings"
)

type AuthV2Settings struct {
	// Platform
	AuthEnabled    bool   `tfschema:"auth_enabled"`
	RuntimeVersion string `tfschema:"runtime_version"`
	ConfigFilePath string `tfschema:"config_file_path"`
	// Global
	RequireAuth           bool     `tfschema:"require_authentication"`
	UnauthenticatedAction string   `tfschema:"unauthenticated_action"`
	DefaultAuthProvider   string   `tfschema:"default_provider"`
	ExcludedPaths         []string `tfschema:"excluded_paths"`
	// IdentityProviders
	AppleAuth                []AppleAuthV2Settings        `tfschema:"apple"`
	AzureActiveDirectoryAuth []AadAuthV2Settings          `tfschema:"active_directory"`
	AzureStaticWebAuth       []StaticWebAppAuthV2Settings `tfschema:"azure_static_web_app"`
	CustomOIDCAuth           []CustomOIDCAuthV2Settings   `tfschema:"custom_oidc"`
	FacebookAuth             []FacebookAuthV2Settings     `tfschema:"facebook"`
	GithubAuth               []GithubAuthV2Settings       `tfschema:"github"`
	GoogleAuth               []GoogleAuthV2Settings       `tfschema:"google"`
	MicrosoftAuth            []MicrosoftAuthV2Settings    `tfschema:"microsoft"`
	TwitterAuth              []TwitterAuthV2Settings      `tfschema:"twitter"`
	// Login
	Login []AuthV2Login `tfschema:"login"`
	// HTTPSettings
	RequireHTTPS                       bool   `tfschema:"require_https"`
	HttpRoutesAPIPrefix                string `tfschema:"http_route_api_prefix"`
	ForwardProxyConvention             string `tfschema:"forward_proxy_convention"`
	ForwardProxyCustomHostHeaderName   string `tfschema:"forward_proxy_custom_host_header_name"`
	ForwardProxyCustomSchemeHeaderName string `tfschema:"forward_proxy_custom_scheme_header_name"`
}

func AuthV2SettingsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"auth_enabled": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Description: "Should the AuthV2 Settings be enabled. Defaults to `false`",
				},

				"runtime_version": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      `~1`,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The Runtime Version of the Authentication and Authorisation feature of this App. Defaults to `~1`",
				},

				"config_file_path": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The path to the App Auth settings. **Note:** Relative Paths are evaluated from the Site Root directory.",
				},

				"require_authentication": {
					Type:        pluginsdk.TypeBool,
					Default:     true,
					Optional:    true,
					Description: "Should the authentication flow be used for all requests. Defaults to `true`",
				},

				"unauthenticated_action": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(web.UnauthenticatedClientActionV2RedirectToLoginPage),
					ValidateFunc: validation.StringInSlice([]string{
						string(web.UnauthenticatedClientActionV2RedirectToLoginPage),
						string(web.UnauthenticatedClientActionV2AllowAnonymous),
						string(web.UnauthenticatedClientActionV2Return401),
						string(web.UnauthenticatedClientActionV2Return403),
					}, false),
				},

				"default_provider": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					// ValidateFunc: validation.StringInSlice([]string{}, false), // TODO - find the correct strings for the Auth names
					Description: "The Default Authentication Provider to use when the `unauthenticated_action` is set to `RedirectToLoginPage`.",
				},

				"excluded_paths": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "The paths which should be excluded from the `unauthenticated_action` when it is set to `RedirectToLoginPage`.",
				},

				"apple": AppleAuthV2SettingsSchema(),

				"active_directory": AadAuthV2SettingsSchema(),

				"azure_static_web_app": StaticWebAppAuthV2SettingsSchema(),

				"custom_oidc": CustomOIDCAuthV2SettingsSchema(),

				"facebook": FacebookAuthV2SettingsSchema(),

				"github": GithubAuthV2SettingsSchema(),

				"google": GoogleAuthV2SettingsSchema(),

				"microsoft": MicrosoftAuthV2SettingsSchema(),

				"twitter": TwitterAuthV2SettingsSchema(),

				"login": authV2LoginSchema(),

				"require_https": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     true,
					Description: "Should HTTPS be required on connections? Defaults to true.",
				},

				"http_route_api_prefix": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      "/.auth",
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The prefix that should precede all the authentication and authorisation paths. Defaults to `/.auth`",
				},

				"forward_proxy_convention": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(web.ForwardProxyConventionNoProxy),
					ValidateFunc: validation.StringInSlice([]string{
						string(web.ForwardProxyConventionNoProxy),
						string(web.ForwardProxyConventionCustom),
						string(web.ForwardProxyConventionStandard),
					}, false),
					Description: "The convention used to determine the url of the request made. Possible values include 'ForwardProxyConventionNoProxy', 'ForwardProxyConventionStandard', 'ForwardProxyConventionCustom'. Defaults to `ForwardProxyConventionNoProxy`.",
				},

				"forward_proxy_custom_host_header_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The name of the header containing the host of the request.",
				},

				"forward_proxy_custom_scheme_header_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The name of the header containing the scheme of the request.",
				},
			},
		},
	}
}

type AuthV2Login struct {
	LogoutEndpoint                string   `tfschema:"logout_endpoint"`
	TokenStoreEnabled             bool     `tfschema:"token_store_enabled"`
	TokenRefreshExtension         float64  `tfschema:"token_refresh_extension_time"`
	TokenFilesystemPath           string   `tfschema:"token_store_path"`
	TokenBlobStorageSAS           string   `tfschema:"token_store_sas_setting_name"`
	PreserveURLFragmentsForLogins bool     `tfschema:"preserve_url_fragments_for_logins"`
	AllowedExternalRedirectURLs   []string `tfschema:"allowed_external_redirect_urls"`
	CookieExpirationConvention    string   `tfschema:"cookie_expiration_convention"`
	CookieExpirationTime          string   `tfschema:"cookie_expiration_time"`
	ValidateNonce                 bool     `tfschema:"validate_nonce"`
	NonceExpirationTime           string   `tfschema:"nonce_expiration_time"`
}

func authV2LoginSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		MaxItems: 1,
		Required: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"logout_endpoint": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Description: "The endpoint to which logout requests should be made.",
				},

				"token_store_enabled": {
					Type:        pluginsdk.TypeBool,
					Computed:    true,
					Description: "Is the Token Store configuration Enabled. True if one of `token_store_path` or `token_store_sas_setting_name` are set to a non-empty value, otherwise `false`.",
				},

				"token_refresh_extension_time": {
					Type:         pluginsdk.TypeFloat,
					Optional:     true,
					Default:      72,
					ValidateFunc: validation.FloatAtLeast(1),
					Description:  "The number of hours after session token expiration that a session token can be used to call the token refresh API. Defaults to `72` hours.",
				},

				"token_store_path": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					ConflictsWith: []string{
						"auth_v2_settings.0.login.0.token_store_sas_setting_name",
					},
					Description: "The directory path in the App Filesystem in which the tokens will be stored.",
				},

				"token_store_sas_setting_name": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ConflictsWith: []string{
						"auth_v2_settings.0.login.0.token_store_path",
					},
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The name of the app setting which contains the SAS URL of the blob storage containing the tokens.",
				},

				"preserve_url_fragments_for_logins": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Should the fragments from the request be preserved after the login request is made. Defaults to `false`",
				},

				"allowed_external_redirect_urls": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					Description: "External URLs that can be redirected to as part of logging in or logging out of the app. This is an advanced setting typically only needed by Windows Store application backends. **Note:** URLs within the current domain are always implicitly allowed.",
				},

				"cookie_expiration_convention": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(web.CookieExpirationConventionFixedTime),
					ValidateFunc: validation.StringInSlice([]string{
						string(web.CookieExpirationConventionIdentityProviderDerived),
						string(web.CookieExpirationConventionFixedTime),
					}, false),
				},

				"cookie_expiration_time": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  "08:00:00",
					// ValidateFunc: // TODO - Find the allowed strings / values for validation
					Description: "The time after the request is made when the session cookie should expire.",
				},

				"validate_nonce": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     true,
					Description: "Should the nonce be validated while completing the login flow. Defaults to `true`",
				},

				"nonce_expiration_time": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  "00:05:00",
					// ValidateFunc: TODO
					Description: "The time after the request is made when the nonce should expire.",
				},
			},
		},
	}
}

func expandAuthV2LoginSettings(input []AuthV2Login) *web.Login {
	if len(input) == 0 {
		return nil
	}
	login := input[0]
	result := &web.Login{
		TokenStore: &web.TokenStore{
			Enabled: utils.Bool(false),
		},
		PreserveURLFragmentsForLogins: utils.Bool(login.PreserveURLFragmentsForLogins),
		Nonce: &web.Nonce{
			ValidateNonce:           nil,
			NonceExpirationInterval: nil,
		},
	}

	if login.TokenFilesystemPath != "" || login.TokenBlobStorageSAS != "" {
		result.TokenStore.Enabled = utils.Bool(true)
		if login.TokenFilesystemPath != "" {
			result.TokenStore.FileSystem = &web.FileSystemTokenStore{
				Directory: utils.String(login.TokenFilesystemPath),
			}
		}
		if login.TokenBlobStorageSAS != "" {
			result.TokenStore.AzureBlobStorage = &web.BlobStorageTokenStore{
				BlobStorageTokenStoreProperties: &web.BlobStorageTokenStoreProperties{
					SasURLSettingName: utils.String(login.TokenBlobStorageSAS),
				},
			}
		}
	}

	if login.LogoutEndpoint != "" {
		result.Routes = &web.LoginRoutes{
			LogoutEndpoint: utils.String(login.LogoutEndpoint),
		}
	}
	result.TokenStore.TokenRefreshExtensionHours = utils.Float(login.TokenRefreshExtension)
	if login.TokenFilesystemPath != "" {
		result.TokenStore.FileSystem = &web.FileSystemTokenStore{
			Directory: utils.String(login.TokenFilesystemPath),
		}
	}
	if login.TokenBlobStorageSAS != "" {
		result.TokenStore.AzureBlobStorage = &web.BlobStorageTokenStore{
			BlobStorageTokenStoreProperties: &web.BlobStorageTokenStoreProperties{
				SasURLSettingName: utils.String(login.TokenBlobStorageSAS),
			},
		}
	}
	if len(login.AllowedExternalRedirectURLs) > 0 {
		result.AllowedExternalRedirectUrls = &login.AllowedExternalRedirectURLs
	}
	if login.CookieExpirationTime != "" || login.CookieExpirationConvention != "" {
		result.CookieExpiration = &web.CookieExpiration{}
		if login.CookieExpirationTime != "" {
			result.CookieExpiration.TimeToExpiration = utils.String(login.CookieExpirationTime)
		}
		if login.CookieExpirationConvention != "" {
			result.CookieExpiration.Convention = web.CookieExpirationConvention(login.CookieExpirationConvention)
		}
	}

	result.Nonce = &web.Nonce{
		ValidateNonce: utils.Bool(login.ValidateNonce),
	}

	if login.NonceExpirationTime != "" {
		if login.NonceExpirationTime != "" {
			result.Nonce.NonceExpirationInterval = utils.String(login.NonceExpirationTime)
		}
	}

	return result
}

func flattenAuthV2LoginSettings(input *web.Login) []AuthV2Login {
	if input == nil {
		return []AuthV2Login{{}}
	}
	result := AuthV2Login{
		PreserveURLFragmentsForLogins: utils.NormaliseNilableBool(input.PreserveURLFragmentsForLogins),
		AllowedExternalRedirectURLs:   nil,
	}
	if routes := input.Routes; routes != nil {
		result.LogoutEndpoint = utils.NormalizeNilableString(routes.LogoutEndpoint)
	}
	if token := input.TokenStore; token != nil {
		result.TokenRefreshExtension = utils.NormalizeNilableFloat(token.TokenRefreshExtensionHours)
		if fs := token.FileSystem; fs != nil {
			result.TokenFilesystemPath = utils.NormalizeNilableString(fs.Directory)
		}
		if bs := token.AzureBlobStorage; bs != nil && bs.BlobStorageTokenStoreProperties != nil {
			result.TokenBlobStorageSAS = utils.NormalizeNilableString(bs.SasURLSettingName)
		}
	}

	if nonce := input.Nonce; nonce != nil {
		result.NonceExpirationTime = utils.NormalizeNilableString(nonce.NonceExpirationInterval)
		result.ValidateNonce = utils.NormaliseNilableBool(nonce.ValidateNonce)
	}

	if cookie := input.CookieExpiration; cookie != nil {
		result.CookieExpirationConvention = string(cookie.Convention)
		result.CookieExpirationTime = utils.NormalizeNilableString(cookie.TimeToExpiration)
	}

	if input.AllowedExternalRedirectUrls != nil {
		result.AllowedExternalRedirectURLs = *input.AllowedExternalRedirectUrls
	}

	return []AuthV2Login{result}
}

type AppleAuthV2Settings struct {
	ClientId                string   `tfschema:"client_id"`
	ClientSecretSettingName string   `tfschema:"client_secret_setting_name"`
	LoginScopes             []string `tfschema:"login_scopes"`
}

func AppleAuthV2SettingsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		AtLeastOneOf: []string{
			"auth_v2_settings.0.apple",
			"auth_v2_settings.0.active_directory",
			"auth_v2_settings.0.azure_static_web_app",
			"auth_v2_settings.0.custom_oidc",
			"auth_v2_settings.0.facebook",
			"auth_v2_settings.0.github",
			"auth_v2_settings.0.google",
			"auth_v2_settings.0.microsoft",
			"auth_v2_settings.0.twitter",
		},
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"client_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The OpenID Connect Client ID for the Apple web application.",
				},

				"client_secret_setting_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The app setting name that contains the `client_secret` value used for Apple Login.",
				},

				"login_scopes": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "Specifies a list of OAuth 2.0 scopes that will be requested as part of Apple Sign-In authentication. If not specified, \"openid\", \"profile\", and \"email\" are used as default scopes.",
				},
			},
		},
	}
}

func AppleAuthV2SettingsSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"client_id": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The OpenID Connect Client ID for the Apple web application.",
				},

				"client_secret_setting_name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The app setting name that contains the `client_secret` value used for Apple Login.",
				},

				"login_scopes": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "Specifies a list of OAuth 2.0 scopes that will be requested as part of Apple Sign-In authentication. If not specified, \"openid\", \"profile\", and \"email\" are used as default scopes.",
				},
			},
		},
	}
}

func expandAppleAuthV2Settings(input []AppleAuthV2Settings) *web.Apple {
	if len(input) == 1 {
		apple := input[0]
		result := &web.Apple{
			AppleProperties: &web.AppleProperties{
				Enabled: utils.Bool(true),
				Registration: &web.AppleRegistration{
					ClientID:                utils.String(apple.ClientId),
					ClientSecretSettingName: utils.String(apple.ClientSecretSettingName),
				},
				Login: &web.LoginScopes{},
			},
		}
		if len(apple.LoginScopes) > 0 {
			result.AppleProperties.Login.Scopes = &apple.LoginScopes
		}

		return result
	}

	return &web.Apple{
		AppleProperties: &web.AppleProperties{
			Enabled: utils.Bool(true),
		},
	}
}

func flattenAppleAuthV2Settings(input *web.Apple) []AppleAuthV2Settings {
	if input == nil || input.AppleProperties == nil {
		return nil
	}
	result := AppleAuthV2Settings{
		ClientId:                "",
		ClientSecretSettingName: "",
		LoginScopes:             nil,
	}

	props := *input.AppleProperties
	if reg := props.Registration; reg != nil {
		result.ClientId = utils.NormalizeNilableString(reg.ClientID)
		result.ClientSecretSettingName = utils.NormalizeNilableString(reg.ClientSecretSettingName)
	}
	if loginScopes := props.Login; loginScopes != nil && loginScopes.Scopes != nil {
		result.LoginScopes = *loginScopes.Scopes
	}

	return []AppleAuthV2Settings{result}
}

type AadAuthV2Settings struct {
	TenantAuthURI                     string            `tfschema:"tenant_auth_endpoint"` // Maps to OpenIDIssuer, takes the form `https://login.microsoftonline.com/v2.0/{tenant-guid}/`
	ClientId                          string            `tfschema:"client_id"`
	ClientSecretSettingName           string            `tfschema:"client_secret_setting_name"`
	ClientSecretCertificateThumbprint string            `tfschema:"client_secret_certificate_thumbprint"`
	LoginParameters                   map[string]string `tfschema:"login_parameters"`
	DisableWWWAuth                    bool              `tfschema:"disable_www_authentication"`
	JWTAllowedGroups                  []string          `tfschema:"jwt_allowed_groups"`
	JWTAllowedClientApps              []string          `tfschema:"jwt_allowed_client_applications"`
	AllowedApplications               []string          `tfschema:"allowed_applications"`
	AllowedAudiences                  []string          `tfschema:"allowed_audiences"`
	AllowedIdentities                 []string          `tfschema:"allowed_identities"`
	AllowedGroups                     []string          `tfschema:"allowed_groups"`
}

func AadAuthV2SettingsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		AtLeastOneOf: []string{
			"auth_v2_settings.0.apple",
			"auth_v2_settings.0.active_directory",
			"auth_v2_settings.0.azure_static_web_app",
			"auth_v2_settings.0.custom_oidc",
			"auth_v2_settings.0.facebook",
			"auth_v2_settings.0.github",
			"auth_v2_settings.0.google",
			"auth_v2_settings.0.microsoft",
			"auth_v2_settings.0.twitter",
		},
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"client_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The ID of the Client to use to authenticate with Azure Active Directory.",
				},

				"tenant_auth_endpoint": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					Description:  "The Azure Tenant Endpoint for the Authenticating Tenant. e.g. `https://login.microsoftonline.com/v2.0/{tenant-guid}/`",
				},

				"client_secret_setting_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					ExactlyOneOf: []string{
						"auth_v2_settings.0.active_directory.0.client_secret_setting_name",
						"auth_v2_settings.0.active_directory.0.client_secret_certificate_thumbprint",
					},
					Description: "The App Setting name that contains the client secret of the Client.",
				},

				"client_secret_certificate_thumbprint": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					ExactlyOneOf: []string{
						"auth_v2_settings.0.active_directory.0.client_secret_setting_name",
						"auth_v2_settings.0.active_directory.0.client_secret_certificate_thumbprint",
					},
					Description: "The thumbprint of the certificate used for signing purposes.",
				},

				"jwt_allowed_groups": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					Description: "A list of Allowed Groups in the JWT Claim.",
				},

				"jwt_allowed_client_applications": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					Description: "A list of Allowed Client Applications in the JWT Claim.",
				},

				"disable_www_authentication": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Description: "Should the www-authenticate provider should be omitted from the request? Defaults to `false`",
				},

				"allowed_groups": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					Description: "The list of allowed Group Names for the Default Authorisation Policy.",
				},

				"allowed_identities": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty, // TODO - Can we use identity Validation here?
					},
					Description: "The list of allowed Identities for the Default Authorisation Policy.",
				},

				"allowed_applications": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty, // TODO - Can we use identity Validation here?
					},
					Description: "The list of allowed Applications for the Default Authorisation Policy.",
				},

				"login_parameters": {
					Type:     pluginsdk.TypeMap,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "A map of key-value pairs to send to the Authorisation Endpoint when a user logs in.",
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

func AadAuthV2SettingsSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"client_id": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The ID of the Client to use to authenticate with Azure Active Directory.",
				},

				"tenant_id": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The Azure Tenant URI for the Authenticating Tenant. e.g. `https://login.microsoftonline.com/v2.0/{tenant-guid}/`",
				},

				"client_secret_setting_name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The App Setting name that contains the client secret of the Client.",
				},

				"client_secret_certificate_thumbprint": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The thumbprint of the certificate used for signing purposes.",
				},

				"jwt_allowed_groups": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "A list of Allowed Groups in the JWT Claim.",
				},

				"jwt_allowed_client_applications": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "A list of Allowed Client Applications in the JWT Claim.",
				},

				"allowed_groups": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"allowed_identities": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"allowed_applications": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"login_parameters": {
					Type:     pluginsdk.TypeMap,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"allowed_audiences": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "Specifies a list of Allowed audience values to consider when validating JWTs issued by Azure Active Directory.",
				},
			},
		},
	}
}

func expandAadAuthV2Settings(input []AadAuthV2Settings) *web.AzureActiveDirectory {
	if len(input) == 1 {
		aad := input[0]
		result := &web.AzureActiveDirectory{
			Enabled: utils.Bool(true),
			Registration: &web.AzureActiveDirectoryRegistration{
				AzureActiveDirectoryRegistrationProperties: &web.AzureActiveDirectoryRegistrationProperties{
					OpenIDIssuer: utils.String(aad.TenantAuthURI),
					ClientID:     utils.String(aad.ClientId),
				},
			},
			Login: &web.AzureActiveDirectoryLogin{
				AzureActiveDirectoryLoginProperties: &web.AzureActiveDirectoryLoginProperties{
					DisableWWWAuthenticate: utils.Bool(aad.DisableWWWAuth),
				},
			},
			Validation: &web.AzureActiveDirectoryValidation{
				AzureActiveDirectoryValidationProperties: &web.AzureActiveDirectoryValidationProperties{
					JwtClaimChecks: &web.JwtClaimChecks{},
					DefaultAuthorizationPolicy: &web.DefaultAuthorizationPolicy{
						AllowedPrincipals: &web.AllowedPrincipals{
							AllowedPrincipalsProperties: &web.AllowedPrincipalsProperties{},
						},
					},
				},
			},
		}

		if aad.ClientSecretCertificateThumbprint != "" {
			result.Registration.ClientSecretCertificateThumbprint = utils.String(aad.ClientSecretCertificateThumbprint)
		} else {
			result.Registration.ClientSecretSettingName = utils.String(aad.ClientSecretSettingName)
		}

		if len(aad.LoginParameters) > 0 {
			params := make([]string, 0)
			for k, v := range aad.LoginParameters {
				params = append(params, fmt.Sprintf("%s=%s", k, v))
			}
			result.Login.LoginParameters = &params
		}

		if len(aad.JWTAllowedGroups) > 0 {
			result.Validation.AzureActiveDirectoryValidationProperties.JwtClaimChecks.AllowedGroups = &aad.JWTAllowedGroups
		}
		if len(aad.JWTAllowedClientApps) > 0 {
			result.Validation.AzureActiveDirectoryValidationProperties.JwtClaimChecks.AllowedClientApplications = &aad.JWTAllowedClientApps
		}
		if len(aad.AllowedAudiences) > 0 {
			result.Validation.AllowedAudiences = &aad.AllowedAudiences
		}
		if len(aad.AllowedGroups) > 0 {
			result.Validation.AzureActiveDirectoryValidationProperties.DefaultAuthorizationPolicy.AllowedPrincipals.AllowedPrincipalsProperties.Groups = &aad.AllowedGroups
		}
		if len(aad.AllowedIdentities) > 0 {
			result.Validation.AzureActiveDirectoryValidationProperties.DefaultAuthorizationPolicy.AllowedPrincipals.AllowedPrincipalsProperties.Identities = &aad.AllowedIdentities
		}
		if len(aad.AllowedApplications) > 0 {
			result.Validation.AzureActiveDirectoryValidationProperties.DefaultAuthorizationPolicy.AllowedApplications = &aad.AllowedApplications
		}

		return result
	}

	return &web.AzureActiveDirectory{
		Enabled: utils.Bool(true),
	}
}

func flattenAadAuthV2Settings(input *web.AzureActiveDirectory) []AadAuthV2Settings {
	if input == nil || input.Registration == nil {
		return nil
	}

	result := AadAuthV2Settings{}

	if reg := input.Registration; reg != nil && reg.AzureActiveDirectoryRegistrationProperties != nil {
		result.TenantAuthURI = utils.NormalizeNilableString(reg.OpenIDIssuer)
		result.ClientId = utils.NormalizeNilableString(reg.ClientID)
		result.ClientSecretSettingName = utils.NormalizeNilableString(reg.ClientSecretSettingName)
		result.ClientSecretCertificateThumbprint = utils.NormalizeNilableString(reg.ClientSecretCertificateThumbprint)
	}

	if login := input.Login; login != nil && login.AzureActiveDirectoryLoginProperties != nil {
		result.DisableWWWAuth = utils.NormaliseNilableBool(login.DisableWWWAuthenticate)
		if loginParamsRaw := login.LoginParameters; loginParamsRaw != nil {
			loginParams := make(map[string]string)
			for _, v := range *loginParamsRaw {
				parts := strings.Split(v, "=")
				if len(parts) == 2 && parts[0] != "" {
					loginParams[parts[0]] = parts[1]
				}
			}
			result.LoginParameters = loginParams
		}

	}

	if validationRaw := input.Validation; validationRaw != nil && validationRaw.AzureActiveDirectoryValidationProperties != nil {
		if validationRaw.AllowedAudiences != nil {
			result.AllowedAudiences = *validationRaw.AllowedAudiences
		}
		if jwt := validationRaw.JwtClaimChecks; jwt != nil {
			if jwt.AllowedGroups != nil {
				result.JWTAllowedGroups = *jwt.AllowedGroups
			}
			if jwt.AllowedClientApplications != nil {
				result.JWTAllowedClientApps = *jwt.AllowedClientApplications
			}
		}
		if defaultPolicy := validationRaw.DefaultAuthorizationPolicy; defaultPolicy != nil {
			if defaultPolicy.AllowedApplications != nil {
				result.AllowedApplications = *defaultPolicy.AllowedApplications
			}
			if defaultPolicy.AllowedPrincipals != nil && defaultPolicy.AllowedPrincipals.AllowedPrincipalsProperties != nil {
				if defaultPolicy.AllowedPrincipals.AllowedPrincipalsProperties.Groups != nil {
					result.AllowedGroups = *defaultPolicy.AllowedPrincipals.AllowedPrincipalsProperties.Groups
				}
				if defaultPolicy.AllowedPrincipals.AllowedPrincipalsProperties.Identities != nil {
					result.AllowedIdentities = *defaultPolicy.AllowedPrincipals.AllowedPrincipalsProperties.Identities
				}
			}
		}
	}

	return []AadAuthV2Settings{result}
}

type StaticWebAppAuthV2Settings struct {
	ClientId string `tfschema:"client_id"`
}

func StaticWebAppAuthV2SettingsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		AtLeastOneOf: []string{
			"auth_v2_settings.0.apple",
			"auth_v2_settings.0.active_directory",
			"auth_v2_settings.0.azure_static_web_app",
			"auth_v2_settings.0.custom_oidc",
			"auth_v2_settings.0.facebook",
			"auth_v2_settings.0.github",
			"auth_v2_settings.0.google",
			"auth_v2_settings.0.microsoft",
			"auth_v2_settings.0.twitter",
		},
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"client_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The ID of the Client to use to authenticate with Azure Static Web App Authentication.",
				},
			},
		},
	}
}

func StaticWebAppAuthV2SettingsSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"client_id": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The ID of the Client used to authenticate with Azure Static Web App Authentication.",
				},
			},
		},
	}
}

func expandStaticWebAppAuthV2Settings(input []StaticWebAppAuthV2Settings) *web.AzureStaticWebApps {
	if len(input) == 1 {
		swa := input[0]
		return &web.AzureStaticWebApps{
			AzureStaticWebAppsProperties: &web.AzureStaticWebAppsProperties{
				Enabled: utils.Bool(true),
				Registration: &web.AzureStaticWebAppsRegistration{
					ClientID: utils.String(swa.ClientId),
				},
			},
		}
	}

	return nil
}

func flattenStaticWebAppAuthV2Settings(input *web.AzureStaticWebApps) []StaticWebAppAuthV2Settings {
	if input == nil {
		return nil
	}

	result := StaticWebAppAuthV2Settings{}

	if props := input.AzureStaticWebAppsProperties; props != nil && utils.NormaliseNilableBool(props.Enabled) {
		if props.Registration != nil {
			result.ClientId = utils.NormalizeNilableString(props.Registration.ClientID)
		}
	}

	return []StaticWebAppAuthV2Settings{result}
}

type CustomOIDCAuthV2Settings struct {
	Name                        string   `tfschema:"name"`
	ClientId                    string   `tfschema:"client_id"`
	ClientCredentialMethod      string   `tfschema:"client_credential_method"`
	ClientSecretSettingName     string   `tfschema:"client_secret_setting_name"`
	AuthorizationEndpoint       string   `tfschema:"authorisation_endpoint"`
	TokenEndpoint               string   `tfschema:"token_endpoint"`
	IssuerEndpoint              string   `tfschema:"issuer_endpoint"`
	CertificationURI            string   `tfschema:"certification_uri"`
	OpenIDConfigurationEndpoint string   `tfschema:"openid_configuration_endpoint"`
	NameClaimType               string   `tfschema:"name_claim_type"`
	Scopes                      []string `tfschema:"scopes"`
}

func CustomOIDCAuthV2SettingsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		AtLeastOneOf: []string{
			"auth_v2_settings.0.apple",
			"auth_v2_settings.0.active_directory",
			"auth_v2_settings.0.azure_static_web_app",
			"auth_v2_settings.0.custom_oidc",
			"auth_v2_settings.0.facebook",
			"auth_v2_settings.0.github",
			"auth_v2_settings.0.google",
			"auth_v2_settings.0.microsoft",
			"auth_v2_settings.0.twitter",
		},
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The name of the Custom OIDC Authentication Provider.",
				},

				"client_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The ID of the Client to use to authenticate with this Custom OIDC.",
				},

				"openid_configuration_endpoint": {
					Type:        pluginsdk.TypeString,
					Required:    true,
					Description: "The endpoint that contains all the configuration endpoints for this Custom OIDC provider.",
				},

				"name_claim_type": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The name of the claim that contains the users name.",
				},

				"scopes": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					Description: "The list of the scopes that should be requested while authenticating.",
				},

				"client_credential_method": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The Client Credential Method used. Currently the only supported value is `ClientSecretPost`",
				},

				"client_secret_setting_name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The App Setting name that contains the secret for this Custom OIDC Client.",
				},

				"authorisation_endpoint": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The endpoint to make the Authorisation Request.",
				},

				"token_endpoint": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The endpoint used to request a Token.",
				},

				"issuer_endpoint": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The endpoint that issued the Token.",
				},

				"certification_uri": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The endpoint that provides the keys necessary to validate the token.",
				},
			},
		},
	}
}

func CustomOIDCAuthV2SettingsSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The name of the Custom OIDC Authentication Provider.",
				},

				"client_id": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The ID of the Client used to authenticate with this Custom OIDC.",
				},

				"name_claim_type": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The name of the claim that contains the users name.",
				},

				"openid_configuration_endpoint": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The endpoint that contains all the configuration endpoints for this Custom OIDC provider.",
				},

				"scopes": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "The list of the scopes that should be requested while authenticating.",
				},

				"client_credential_method": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The Client Credential Method used. Currently the only supported value is `ClientSecretPost`",
				},

				"client_secret_setting_name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The App Setting name that contains the secret for this Custom OIDC Client.",
				},

				"authorisation_endpoint": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The endpoint to make the Authorisation Request.",
				},

				"token_endpoint": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The endpoint used to request a Token.",
				},

				"issuer_endpoint": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The endpoint that issued the Token.",
				},

				"certification_uri": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The endpoint that provides the keys necessary to validate the token.",
				},
			},
		},
	}
}

func expandCustomOIDCAuthV2Settings(input []CustomOIDCAuthV2Settings) map[string]*web.CustomOpenIDConnectProvider {
	if len(input) == 0 {
		return nil
	}
	result := make(map[string]*web.CustomOpenIDConnectProvider)
	for _, v := range input {
		if v.Name == "" {
			continue
		}
		provider := &web.CustomOpenIDConnectProvider{
			CustomOpenIDConnectProviderProperties: &web.CustomOpenIDConnectProviderProperties{
				Enabled: utils.Bool(true),
				Registration: &web.OpenIDConnectRegistration{
					ClientID:         utils.String(v.ClientId),
					ClientCredential: nil,
					OpenIDConnectConfiguration: &web.OpenIDConnectConfig{
						WellKnownOpenIDConfiguration: utils.String(v.OpenIDConfigurationEndpoint),
					},
				},
				Login: &web.OpenIDConnectLogin{
					NameClaimType: nil,
					Scopes:        nil,
				},
			},
		}

		if v.NameClaimType != "" {
			provider.Login.NameClaimType = utils.String(v.NameClaimType)

		}
		if len(v.Scopes) > 0 {
			provider.Login.Scopes = &v.Scopes
		}

		result[v.Name] = provider
	}

	return result
}

func flattenCustomOIDCAuthV2Settings(input map[string]*web.CustomOpenIDConnectProvider) []CustomOIDCAuthV2Settings {
	if len(input) == 0 {
		return nil
	}

	result := make([]CustomOIDCAuthV2Settings, 0)
	for k, v := range input {
		if props := v.CustomOpenIDConnectProviderProperties; props == nil || props.Enabled == nil || !*props.Enabled {
			continue
		} else {
			provider := CustomOIDCAuthV2Settings{
				Name: k,
			}
			if reg := props.Registration; reg != nil {
				provider.ClientId = utils.NormalizeNilableString(reg.ClientID)
				if reg.ClientCredential != nil {
					provider.ClientSecretSettingName = utils.NormalizeNilableString(reg.ClientCredential.ClientSecretSettingName)
					provider.ClientCredentialMethod = string(reg.ClientCredential.Method)
				}
				if config := reg.OpenIDConnectConfiguration; config != nil {
					provider.OpenIDConfigurationEndpoint = utils.NormalizeNilableString(config.WellKnownOpenIDConfiguration)
					provider.AuthorizationEndpoint = utils.NormalizeNilableString(config.AuthorizationEndpoint)
					provider.TokenEndpoint = utils.NormalizeNilableString(config.TokenEndpoint)
					provider.IssuerEndpoint = utils.NormalizeNilableString(config.Issuer)
					provider.CertificationURI = utils.NormalizeNilableString(config.CertificationURI)
				}
			}
			if login := props.Login; login != nil {
				if login.Scopes != nil {
					provider.Scopes = *login.Scopes
				}
				provider.NameClaimType = utils.NormalizeNilableString(login.NameClaimType)
			}
			result = append(result, provider)
		}
	}

	return result
}

type FacebookAuthV2Settings struct {
	AppId                string   `tfschema:"app_id"`
	AppSecretSettingName string   `tfschema:"app_secret_setting_name"`
	LoginScopes          []string `tfschema:"login_scopes"`
	GraphAPIVersion      string   `tfschema:"graph_api_version"`
}

func FacebookAuthV2SettingsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		AtLeastOneOf: []string{
			"auth_v2_settings.0.apple",
			"auth_v2_settings.0.active_directory",
			"auth_v2_settings.0.azure_static_web_app",
			"auth_v2_settings.0.custom_oidc",
			"auth_v2_settings.0.facebook",
			"auth_v2_settings.0.github",
			"auth_v2_settings.0.google",
			"auth_v2_settings.0.microsoft",
			"auth_v2_settings.0.twitter",
		},
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"app_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The App ID of the Facebook app used for login.",
				},

				"app_secret_setting_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The app setting name that contains the `app_secret` value used for Facebook Login.",
				},

				"graph_api_version": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The version of the Facebook API to be used while logging in.",
				},

				"login_scopes": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "Specifies a list of scopes to be requested as part of Facebook Login authentication.",
				},
			},
		},
	}
}

func FacebookAuthV2SettingsSchemaComputed() *pluginsdk.Schema {
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

				"app_secret_setting_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The app setting name that contains the `app_secret` value used for Facebook Login.",
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

func expandFacebookAuthV2Settings(input []FacebookAuthV2Settings) *web.Facebook {
	if len(input) == 1 {
		facebook := input[0]
		result := &web.Facebook{
			Enabled: utils.Bool(true),
			Registration: &web.AppRegistration{
				AppRegistrationProperties: &web.AppRegistrationProperties{
					AppID:                utils.String(facebook.AppId),
					AppSecretSettingName: utils.String(facebook.AppSecretSettingName),
				},
			},
		}

		if facebook.GraphAPIVersion != "" {
			result.GraphAPIVersion = utils.String(facebook.GraphAPIVersion)
		}
		if len(facebook.LoginScopes) > 0 {
			result.Login = &web.LoginScopes{
				Scopes: &facebook.LoginScopes,
			}
		}

		return result
	}

	return &web.Facebook{
		Enabled: utils.Bool(true),
	}
}

func flattenFacebookAuthV2Settings(input *web.Facebook) []FacebookAuthV2Settings {
	if input == nil {
		return nil
	}

	result := FacebookAuthV2Settings{
		GraphAPIVersion: utils.NormalizeNilableString(input.GraphAPIVersion),
	}

	if reg := input.Registration; reg != nil {
		if reg.AppRegistrationProperties != nil {
			result.AppId = utils.NormalizeNilableString(reg.AppID)
			result.AppSecretSettingName = utils.NormalizeNilableString(reg.AppSecretSettingName)
		}
	}
	if login := input.Login; login != nil && login.Scopes != nil && len(*login.Scopes) != 0 {
		result.LoginScopes = *login.Scopes
	}

	return []FacebookAuthV2Settings{result}
}

type GithubAuthV2Settings struct {
	ClientId                string   `tfschema:"client_id"`
	ClientSecretSettingName string   `tfschema:"client_secret_setting_name"`
	LoginScopes             []string `tfschema:"login_scopes"`
}

func GithubAuthV2SettingsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		AtLeastOneOf: []string{
			"auth_v2_settings.0.apple",
			"auth_v2_settings.0.active_directory",
			"auth_v2_settings.0.azure_static_web_app",
			"auth_v2_settings.0.custom_oidc",
			"auth_v2_settings.0.facebook",
			"auth_v2_settings.0.github",
			"auth_v2_settings.0.google",
			"auth_v2_settings.0.microsoft",
			"auth_v2_settings.0.twitter",
		},
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"client_id": {
					Type:        pluginsdk.TypeString,
					Required:    true,
					Description: "The ID of the GitHub app used for login.",
				},

				"client_secret_setting_name": {
					Type:        pluginsdk.TypeString,
					Required:    true,
					Description: "The app setting name that contains the `client_secret` value used for GitHub Login.",
				},

				"login_scopes": {
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

func GithubAuthV2SettingsSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"client_id": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The ID of the GitHub app used for login.",
				},

				"client_secret_setting_name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The app setting name that contains the `client_secret` value used for GitHub Login.",
				},

				"login_scopes": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "Specifies a list of OAuth 2.0 scopes that will be requested as part of GitHub Login authentication.",
				},
			},
		},
	}
}

func expandGitHubAuthV2Settings(input []GithubAuthV2Settings) *web.GitHub {
	if len(input) == 1 {
		github := input[0]
		result := &web.GitHub{
			GitHubProperties: &web.GitHubProperties{
				Enabled: utils.Bool(true),
				Registration: &web.ClientRegistration{
					ClientID:                utils.String(github.ClientId),
					ClientSecretSettingName: utils.String(github.ClientSecretSettingName),
				},
				Login: &web.LoginScopes{},
			},
		}
		if len(github.LoginScopes) > 0 {
			result.GitHubProperties.Login = &web.LoginScopes{Scopes: &github.LoginScopes}
		}

		return result
	}

	return &web.GitHub{
		GitHubProperties: &web.GitHubProperties{
			Enabled: utils.Bool(true),
		},
	}
}

func flattenGitHubAuthV2Settings(input *web.GitHub) []GithubAuthV2Settings {
	if input == nil || input.GitHubProperties == nil {
		return nil
	}

	props := *input.GitHubProperties
	if props.Enabled == nil || !*props.Enabled {
		return nil
	}

	result := GithubAuthV2Settings{}

	if reg := props.Registration; reg != nil {
		result.ClientId = utils.NormalizeNilableString(reg.ClientID)
		result.ClientSecretSettingName = utils.NormalizeNilableString(reg.ClientSecretSettingName)
	}
	if login := props.Login; login != nil && login.Scopes != nil {
		result.LoginScopes = *login.Scopes
	}

	return []GithubAuthV2Settings{result}
}

type GoogleAuthV2Settings struct {
	ClientId                string   `tfschema:"client_id"`
	ClientSecretSettingName string   `tfschema:"client_secret_setting_name"`
	AllowedAudiences        []string `tfschema:"allowed_audiences"`
	LoginScopes             []string `tfschema:"login_scopes"`
}

func GoogleAuthV2SettingsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		AtLeastOneOf: []string{
			"auth_v2_settings.0.apple",
			"auth_v2_settings.0.active_directory",
			"auth_v2_settings.0.azure_static_web_app",
			"auth_v2_settings.0.custom_oidc",
			"auth_v2_settings.0.facebook",
			"auth_v2_settings.0.github",
			"auth_v2_settings.0.google",
			"auth_v2_settings.0.microsoft",
			"auth_v2_settings.0.twitter",
		},
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"client_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The OpenID Connect Client ID for the Google web application.",
				},

				"client_secret_setting_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The app setting name that contains the `client_secret` value used for Google Login.",
				},

				"allowed_audiences": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					Description: "Specifies a list of Allowed Audiences that will be requested as part of Google Sign-In authentication.",
				},

				"login_scopes": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					Description: "Specifies a list of Login scopes that will be requested as part of Google Sign-In authentication.",
				},
			},
		},
	}
}

func GoogleAuthV2SettingsSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"client_id": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The OpenID Connect Client ID for the Google web application.",
				},

				"client_secret_setting_name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The app setting name that contains the `client_secret` value used for Google Login.",
				},

				"allowed_audiences": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "Specifies a list of Allowed Audiences that will be requested as part of Google Sign-In authentication. If not specified, \"openid\", \"profile\", and \"email\" are used as default scopes.",
				},

				"login_scopes": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "The list of OAuth 2.0 scopes requested as part of Google Sign-In authentication. If not specified, \"openid\", \"profile\", and \"email\" are used as default scopes.",
				},
			},
		},
	}
}

func expandGoogleAuthV2Settings(input []GoogleAuthV2Settings) *web.Google {
	if len(input) == 1 {
		google := input[0]
		result := &web.Google{
			GoogleProperties: &web.GoogleProperties{
				Enabled: utils.Bool(true),
				Registration: &web.ClientRegistration{
					ClientID:                utils.String(google.ClientId),
					ClientSecretSettingName: utils.String(google.ClientSecretSettingName),
				},
			},
		}

		if len(google.AllowedAudiences) != 0 {
			result.GoogleProperties.Validation = &web.AllowedAudiencesValidation{
				AllowedAudiences: &google.AllowedAudiences,
			}
		}
		if len(google.LoginScopes) != 0 {
			result.GoogleProperties.Login = &web.LoginScopes{
				Scopes: &google.LoginScopes,
			}
		}

		return result
	}

	return &web.Google{
		GoogleProperties: &web.GoogleProperties{
			Enabled: utils.Bool(true),
		},
	}
}

func flattenGoogleAuthV2Settings(input *web.Google) []GoogleAuthV2Settings {
	if input == nil || input.GoogleProperties == nil {
		return nil
	}

	props := *input.GoogleProperties
	if props.Enabled == nil || !*props.Enabled {
		return nil
	}

	result := GoogleAuthV2Settings{}

	if reg := props.Registration; reg != nil {
		result.ClientId = utils.NormalizeNilableString(reg.ClientID)
		result.ClientSecretSettingName = utils.NormalizeNilableString(reg.ClientSecretSettingName)
	}
	if login := input.Login; login != nil && login.Scopes != nil {
		result.LoginScopes = *login.Scopes
	}
	if val := input.Validation; val != nil && val.AllowedAudiences != nil {
		result.LoginScopes = *val.AllowedAudiences
	}

	return []GoogleAuthV2Settings{result}
}

type MicrosoftAuthV2Settings struct {
	ClientId                string   `tfschema:"client_id"`
	ClientSecretSettingName string   `tfschema:"client_secret_setting_name"`
	AllowedAudiences        []string `tfschema:"allowed_audiences"`
	LoginScopes             []string `tfschema:"login_scopes"`
}

func MicrosoftAuthV2SettingsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		AtLeastOneOf: []string{
			"auth_v2_settings.0.apple",
			"auth_v2_settings.0.active_directory",
			"auth_v2_settings.0.azure_static_web_app",
			"auth_v2_settings.0.custom_oidc",
			"auth_v2_settings.0.facebook",
			"auth_v2_settings.0.github",
			"auth_v2_settings.0.google",
			"auth_v2_settings.0.microsoft",
			"auth_v2_settings.0.twitter",
		},
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"client_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The OAuth 2.0 client ID that was created for the app used for authentication.",
				},

				"client_secret_setting_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The app setting name containing the OAuth 2.0 client secret that was created for the app used for authentication.",
				},

				"allowed_audiences": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					Description: "Specifies a list of Allowed Audiences that will be requested as part of Microsoft Sign-In authentication.",
				},

				"login_scopes": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "The list of Login scopes that will be requested as part of Microsoft Account authentication.",
				},
			},
		},
	}
}

func MicrosoftAuthV2SettingsSchemaComputed() *pluginsdk.Schema {
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

				"client_secret_setting_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The app setting name containing the OAuth 2.0 client secret that was created for the app used for authentication.",
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

func expandMicrosoftAuthV2Settings(input []MicrosoftAuthV2Settings) *web.LegacyMicrosoftAccount {
	if len(input) == 1 {
		msft := input[0]
		result := &web.LegacyMicrosoftAccount{
			LegacyMicrosoftAccountProperties: &web.LegacyMicrosoftAccountProperties{
				Enabled: utils.Bool(true),
				Registration: &web.ClientRegistration{
					ClientID:                utils.String(msft.ClientId),
					ClientSecretSettingName: utils.String(msft.ClientSecretSettingName),
				},
			},
		}
		if len(msft.AllowedAudiences) != 0 {
			result.LegacyMicrosoftAccountProperties.Validation = &web.AllowedAudiencesValidation{
				AllowedAudiences: &msft.AllowedAudiences,
			}
		}
		if len(msft.LoginScopes) != 0 {
			result.LegacyMicrosoftAccountProperties.Login = &web.LoginScopes{
				Scopes: &msft.LoginScopes,
			}
		}

		return result
	}

	return &web.LegacyMicrosoftAccount{
		LegacyMicrosoftAccountProperties: &web.LegacyMicrosoftAccountProperties{
			Enabled: utils.Bool(true),
		},
	}
}

func flattenMicrosoftAuthV2Settings(input *web.LegacyMicrosoftAccount) []MicrosoftAuthV2Settings {
	if input == nil || input.LegacyMicrosoftAccountProperties == nil {
		return nil
	}

	props := *input.LegacyMicrosoftAccountProperties
	if props.Enabled == nil || !*props.Enabled {
		return nil
	}

	result := MicrosoftAuthV2Settings{}

	if reg := props.Registration; reg != nil {
		result.ClientId = utils.NormalizeNilableString(reg.ClientID)
		result.ClientSecretSettingName = utils.NormalizeNilableString(reg.ClientSecretSettingName)
	}
	if login := input.Login; login != nil && login.Scopes != nil {
		result.LoginScopes = *login.Scopes
	}
	if val := input.Validation; val != nil && val.AllowedAudiences != nil {
		result.LoginScopes = *val.AllowedAudiences
	}

	return []MicrosoftAuthV2Settings{result}
}

type TwitterAuthV2Settings struct {
	ConsumerKey               string `tfschema:"consumer_key"`
	ConsumerSecretSettingName string `tfschema:"consumer_secret_setting_name"`
}

func TwitterAuthV2SettingsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		AtLeastOneOf: []string{
			"auth_v2_settings.0.apple",
			"auth_v2_settings.0.active_directory",
			"auth_v2_settings.0.azure_static_web_app",
			"auth_v2_settings.0.custom_oidc",
			"auth_v2_settings.0.facebook",
			"auth_v2_settings.0.github",
			"auth_v2_settings.0.google",
			"auth_v2_settings.0.microsoft",
			"auth_v2_settings.0.twitter",
		},
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"consumer_key": {
					Type:        pluginsdk.TypeString,
					Required:    true,
					Description: "The OAuth 1.0a consumer key of the Twitter application used for sign-in.",
				},

				"consumer_secret_setting_name": {
					Type:        pluginsdk.TypeString,
					Required:    true,
					Description: "The app setting name that contains the OAuth 1.0a consumer secret of the Twitter application used for sign-in.",
				},
			},
		},
	}
}

func TwitterAuthV2SettingsSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"consumer_key": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The OAuth 1.0a consumer key of the Twitter application used for sign-in.",
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

func expandTwitterAuthV2Settings(input []TwitterAuthV2Settings) *web.Twitter {
	if len(input) == 1 {
		twitter := input[0]
		result := &web.Twitter{
			TwitterProperties: &web.TwitterProperties{
				Enabled: utils.Bool(true),
				Registration: &web.TwitterRegistration{
					ConsumerKey:               utils.String(twitter.ConsumerKey),
					ConsumerSecretSettingName: utils.String(twitter.ConsumerSecretSettingName),
				},
			},
		}

		return result
	}

	return &web.Twitter{
		TwitterProperties: &web.TwitterProperties{
			Enabled: utils.Bool(true),
		},
	}
}

func flattenTwitterAuthV2Settings(input *web.Twitter) []TwitterAuthV2Settings {
	if input == nil || input.TwitterProperties == nil {
		return nil
	}

	props := *input.TwitterProperties
	if utils.NormaliseNilableBool(props.Enabled) {
		result := TwitterAuthV2Settings{}
		if reg := props.Registration; reg != nil {
			result.ConsumerKey = utils.NormalizeNilableString(reg.ConsumerKey)
			result.ConsumerSecretSettingName = utils.NormalizeNilableString(reg.ConsumerSecretSettingName)
		}
		return []TwitterAuthV2Settings{result}
	}

	return nil
}

func ExpandAuthV2Settings(input []AuthV2Settings) *web.SiteAuthSettingsV2 {
	result := &web.SiteAuthSettingsV2{}
	if len(input) != 1 {
		return result
	}

	settings := input[0]

	props := &web.SiteAuthSettingsV2Properties{
		Platform: &web.AuthPlatform{
			Enabled: utils.Bool(settings.AuthEnabled),
		},
		GlobalValidation: &web.GlobalValidation{
			RequireAuthentication: utils.Bool(settings.RequireAuth),
		},
		IdentityProviders: &web.IdentityProviders{
			AzureActiveDirectory:         expandAadAuthV2Settings(settings.AzureActiveDirectoryAuth),
			Facebook:                     expandFacebookAuthV2Settings(settings.FacebookAuth),
			GitHub:                       expandGitHubAuthV2Settings(settings.GithubAuth),
			Google:                       expandGoogleAuthV2Settings(settings.GoogleAuth),
			Twitter:                      expandTwitterAuthV2Settings(settings.TwitterAuth),
			CustomOpenIDConnectProviders: expandCustomOIDCAuthV2Settings(settings.CustomOIDCAuth),
			LegacyMicrosoftAccount:       expandMicrosoftAuthV2Settings(settings.MicrosoftAuth),
			Apple:                        expandAppleAuthV2Settings(settings.AppleAuth),
			AzureStaticWebApps:           expandStaticWebAppAuthV2Settings(settings.AzureStaticWebAuth),
		},
		Login: expandAuthV2LoginSettings(settings.Login),
		HTTPSettings: &web.HTTPSettings{
			RequireHTTPS: utils.Bool(settings.RequireHTTPS),
			ForwardProxy: &web.ForwardProxy{
				Convention: web.ForwardProxyConvention(settings.ForwardProxyConvention),
			},
		},
	}

	// Platform
	if settings.RuntimeVersion != "" {
		props.Platform.RuntimeVersion = utils.String(settings.RuntimeVersion)
	}
	if settings.ConfigFilePath != "" {
		props.Platform.ConfigFilePath = utils.String(settings.ConfigFilePath)
	}

	// Global
	if settings.UnauthenticatedAction != "" {
		props.GlobalValidation.UnauthenticatedClientAction = web.UnauthenticatedClientActionV2(settings.UnauthenticatedAction)
	}
	if settings.DefaultAuthProvider != "" {
		props.GlobalValidation.RedirectToProvider = utils.String(settings.DefaultAuthProvider)
	}
	if len(settings.ExcludedPaths) > 0 {
		props.GlobalValidation.ExcludedPaths = &settings.ExcludedPaths
	}

	// HTTP
	if settings.ForwardProxyCustomHostHeaderName != "" {
		props.HTTPSettings.ForwardProxy.CustomHostHeaderName = utils.String(settings.ForwardProxyCustomHostHeaderName)
	}
	if settings.ForwardProxyCustomSchemeHeaderName != "" {
		props.HTTPSettings.ForwardProxy.CustomProtoHeaderName = utils.String(settings.ForwardProxyCustomSchemeHeaderName)
	}
	if settings.HttpRoutesAPIPrefix != "" {
		props.HTTPSettings = &web.HTTPSettings{
			Routes: &web.HTTPSettingsRoutes{
				APIPrefix: utils.String(settings.HttpRoutesAPIPrefix),
			},
		}
	}

	result.SiteAuthSettingsV2Properties = props

	return result
}

func FlattenAuthV2Settings(input web.SiteAuthSettingsV2) []AuthV2Settings {
	if input.SiteAuthSettingsV2Properties == nil {
		return nil
	}

	settings := *input.SiteAuthSettingsV2Properties

	result := AuthV2Settings{}

	if platform := settings.Platform; platform != nil {
		result.AuthEnabled = utils.NormaliseNilableBool(platform.Enabled)
		result.RuntimeVersion = utils.NormalizeNilableString(platform.RuntimeVersion)
		result.ConfigFilePath = utils.NormalizeNilableString(platform.ConfigFilePath)
	}

	if global := settings.GlobalValidation; global != nil {
		result.RequireAuth = utils.NormaliseNilableBool(global.RequireAuthentication)
		result.UnauthenticatedAction = string(global.UnauthenticatedClientAction)
		result.DefaultAuthProvider = utils.NormalizeNilableString(global.RedirectToProvider)
		if global.ExcludedPaths != nil && len(*global.ExcludedPaths) > 0 {
			result.ExcludedPaths = *global.ExcludedPaths
		}
	}

	if http := settings.HTTPSettings; http != nil {
		result.RequireHTTPS = utils.NormaliseNilableBool(http.RequireHTTPS)
		if http.Routes != nil {
			result.HttpRoutesAPIPrefix = utils.NormalizeNilableString(http.Routes.APIPrefix)
		}
		if fp := http.ForwardProxy; fp != nil {
			result.ForwardProxyConvention = string(fp.Convention)
			result.ForwardProxyCustomHostHeaderName = utils.NormalizeNilableString(fp.CustomHostHeaderName)
			result.ForwardProxyCustomSchemeHeaderName = utils.NormalizeNilableString(fp.CustomProtoHeaderName)
		}
	}

	if login := settings.Login; login != nil {
		result.Login = flattenAuthV2LoginSettings(login)
	}

	if authProviders := settings.IdentityProviders; authProviders != nil {
		result.AppleAuth = flattenAppleAuthV2Settings(authProviders.Apple)
		result.AzureActiveDirectoryAuth = flattenAadAuthV2Settings(authProviders.AzureActiveDirectory)
		result.AzureStaticWebAuth = flattenStaticWebAppAuthV2Settings(authProviders.AzureStaticWebApps)
		result.CustomOIDCAuth = flattenCustomOIDCAuthV2Settings(authProviders.CustomOpenIDConnectProviders)
		result.FacebookAuth = flattenFacebookAuthV2Settings(authProviders.Facebook)
		result.GithubAuth = flattenGitHubAuthV2Settings(authProviders.GitHub)
		result.GoogleAuth = flattenGoogleAuthV2Settings(authProviders.Google)
		result.MicrosoftAuth = flattenMicrosoftAuthV2Settings(authProviders.LegacyMicrosoftAccount)
		result.TwitterAuth = flattenTwitterAuthV2Settings(authProviders.Twitter)
	}

	return []AuthV2Settings{result}
}
