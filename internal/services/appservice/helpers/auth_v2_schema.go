// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/tombuildsstuff/kermit/sdk/web/2022-09-01/web"
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
	AppleAuth                []AppleAuthV2Settings        `tfschema:"apple_v2"`
	AzureActiveDirectoryAuth []AadAuthV2Settings          `tfschema:"active_directory_v2"`
	AzureStaticWebAuth       []StaticWebAppAuthV2Settings `tfschema:"azure_static_web_app_v2"`
	CustomOIDCAuth           []CustomOIDCAuthV2Settings   `tfschema:"custom_oidc_v2"`
	FacebookAuth             []FacebookAuthV2Settings     `tfschema:"facebook_v2"`
	GithubAuth               []GithubAuthV2Settings       `tfschema:"github_v2"`
	GoogleAuth               []GoogleAuthV2Settings       `tfschema:"google_v2"`
	MicrosoftAuth            []MicrosoftAuthV2Settings    `tfschema:"microsoft_v2"`
	TwitterAuth              []TwitterAuthV2Settings      `tfschema:"twitter_v2"`
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
					Optional:    true,
					Description: "Should the authentication flow be used for all requests.",
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
					Description: "The action to take for requests made without authentication. Possible values include `RedirectToLoginPage`, `AllowAnonymous`, `Return401`, and `Return403`. Defaults to `RedirectToLoginPage`.",
				},

				"default_provider": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					// ValidateFunc: validation.StringInSlice([]string{}, false), // TODO - find the correct strings for the Auth names
					Description: "The Default Authentication Provider to use when the `unauthenticated_action` is set to `RedirectToLoginPage`. Possible values include: `apple`, `azureactivedirectory`, `facebook`, `github`, `google`, `twitter` and the `name` of your `custom_oidc_v2` provider.",
				},

				"excluded_paths": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "The paths which should be excluded from the `unauthenticated_action` when it is set to `RedirectToLoginPage`.",
				},

				"apple_v2": AppleAuthV2SettingsSchema(),

				"active_directory_v2": AadAuthV2SettingsSchema(),

				"azure_static_web_app_v2": StaticWebAppAuthV2SettingsSchema(),

				"custom_oidc_v2": CustomOIDCAuthV2SettingsSchema(),

				"facebook_v2": FacebookAuthV2SettingsSchema(),

				"github_v2": GithubAuthV2SettingsSchema(),

				"google_v2": GoogleAuthV2SettingsSchema(),

				"microsoft_v2": MicrosoftAuthV2SettingsSchema(),

				"twitter_v2": TwitterAuthV2SettingsSchema(),

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
					Description: "The convention used to determine the url of the request made. Possible values include `ForwardProxyConventionNoProxy`, `ForwardProxyConventionStandard`, `ForwardProxyConventionCustom`. Defaults to `ForwardProxyConventionNoProxy`",
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
func AuthV2SettingsComputedSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"auth_enabled": {
					Type:        pluginsdk.TypeBool,
					Computed:    true,
					Description: "Is the AuthV2 Settings be enabled.",
				},

				"runtime_version": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The Runtime Version of the Authentication and Authorisation feature of this App.",
				},

				"config_file_path": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The path to the App Auth settings.",
				},

				"require_authentication": {
					Type:        pluginsdk.TypeBool,
					Computed:    true,
					Description: "Is the authentication flow used for all requests.",
				},

				"unauthenticated_action": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"default_provider": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The Default Authentication Provider used when the `unauthenticated_action` is set to `RedirectToLoginPage`.",
				},

				"excluded_paths": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "The paths which are excluded from the `unauthenticated_action` when it is set to `RedirectToLoginPage`.",
				},

				"apple_v2": AppleAuthV2SettingsSchemaComputed(),

				"active_directory_v2": AadAuthV2SettingsSchemaComputed(),

				"azure_static_web_app_v2": StaticWebAppAuthV2SettingsSchemaComputed(),

				"custom_oidc_v2": CustomOIDCAuthV2SettingsSchemaComputed(),

				"facebook_v2": FacebookAuthV2SettingsSchemaComputed(),

				"github_v2": GithubAuthV2SettingsSchemaComputed(),

				"google_v2": GoogleAuthV2SettingsSchemaComputed(),

				"microsoft_v2": MicrosoftAuthV2SettingsSchemaComputed(),

				"twitter_v2": TwitterAuthV2SettingsSchemaComputed(),

				"login": authV2LoginSchemaComputed(),

				"require_https": {
					Type:        pluginsdk.TypeBool,
					Computed:    true,
					Description: "Is HTTPS required on connections?",
				},

				"http_route_api_prefix": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The prefix that precedes all the authentication and authorisation paths.",
				},

				"forward_proxy_convention": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The convention used to determine the url of the request made.",
				},

				"forward_proxy_custom_host_header_name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The name of the header containing the host of the request.",
				},

				"forward_proxy_custom_scheme_header_name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The name of the header containing the scheme of the request.",
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
					Optional:    true,
					Default:     false,
					Description: "Should the Token Store configuration Enabled. Defaults to `false`",
				},

				"token_refresh_extension_time": {
					Type:         pluginsdk.TypeFloat,
					Optional:     true,
					Default:      72,
					ValidateFunc: validation.FloatAtLeast(0),
					Description:  "The number of hours after session token expiration that a session token can be used to call the token refresh API. Defaults to `72` hours.",
				},

				"token_store_path": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					ConflictsWith: []string{
						"auth_settings_v2.0.login.0.token_store_sas_setting_name",
					},
					Description: "The directory path in the App Filesystem in which the tokens will be stored.",
				},

				"token_store_sas_setting_name": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ConflictsWith: []string{
						"auth_settings_v2.0.login.0.token_store_path",
					},
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The name of the app setting which contains the SAS URL of the blob storage containing the tokens.",
				},

				"preserve_url_fragments_for_logins": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Should the fragments from the request be preserved after the login request is made. Defaults to `false`.",
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
					Description: "The method by which cookies expire. Possible values include: `FixedTime`, and `IdentityProviderDerived`. Defaults to `FixedTime`.",
				},

				"cookie_expiration_time": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      "08:00:00",
					ValidateFunc: validate.TimeInterval,
					Description:  "The time after the request is made when the session cookie should expire. Defaults to `08:00:00`.",
				},

				"validate_nonce": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     true,
					Description: "Should the nonce be validated while completing the login flow. Defaults to `true`.",
				},

				"nonce_expiration_time": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      "00:05:00",
					ValidateFunc: validate.TimeInterval,
					Description:  "The time after the request is made when the nonce should expire. Defaults to `00:05:00`.",
				},
			},
		},
	}
}

func authV2LoginSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"logout_endpoint": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The endpoint to which logout requests are made.",
				},

				"token_store_enabled": {
					Type:        pluginsdk.TypeBool,
					Computed:    true,
					Description: "Is the Token Store configuration Enabled.",
				},

				"token_refresh_extension_time": {
					Type:        pluginsdk.TypeFloat,
					Computed:    true,
					Description: "The number of hours after session token expiration that a session token can be used to call the token refresh API.",
				},

				"token_store_path": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The directory path in the App Filesystem in which the tokens are stored.",
				},

				"token_store_sas_setting_name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The name of the app setting which contains the SAS URL of the blob storage containing the tokens.",
				},

				"preserve_url_fragments_for_logins": {
					Type:        pluginsdk.TypeBool,
					Computed:    true,
					Description: "Are the fragments from the request be preserved after the login request is made.",
				},

				"allowed_external_redirect_urls": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "External URLs that can be redirected to as part of logging in or logging out of the app. This is an advanced setting typically only needed by Windows Store application backends. **Note:** URLs within the current domain are always implicitly allowed.",
				},

				"cookie_expiration_convention": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"cookie_expiration_time": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The time after the request is made when the session cookie will expire.",
				},

				"validate_nonce": {
					Type:        pluginsdk.TypeBool,
					Computed:    true,
					Description: "Is the nonce be validated while completing the login flow.",
				},

				"nonce_expiration_time": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The time after the request is made when the nonce will expire.",
				},
			},
		},
	}
}

func expandAuthV2LoginSettings(input []AuthV2Login) *webapps.Login {
	if len(input) == 0 {
		return nil
	}
	login := input[0]
	result := &webapps.Login{
		Routes: &webapps.LoginRoutes{},
		TokenStore: &webapps.TokenStore{
			Enabled:          pointer.To(login.TokenStoreEnabled),
			FileSystem:       &webapps.FileSystemTokenStore{},
			AzureBlobStorage: &webapps.BlobStorageTokenStore{},
		},
		PreserveURLFragmentsForLogins: pointer.To(login.PreserveURLFragmentsForLogins),
		Nonce: &webapps.Nonce{
			ValidateNonce:           pointer.To(login.ValidateNonce),
			NonceExpirationInterval: pointer.To(login.NonceExpirationTime),
		},
		CookieExpiration: &webapps.CookieExpiration{
			Convention:       pointer.To(webapps.CookieExpirationConvention(login.CookieExpirationConvention)),
			TimeToExpiration: pointer.To(login.CookieExpirationTime),
		},
	}

	if login.TokenFilesystemPath != "" || login.TokenBlobStorageSAS != "" {
		result.TokenStore.Enabled = pointer.To(true)
		if login.TokenFilesystemPath != "" {
			result.TokenStore.FileSystem = &webapps.FileSystemTokenStore{
				Directory: pointer.To(login.TokenFilesystemPath),
			}
		}
		if login.TokenBlobStorageSAS != "" {
			result.TokenStore.AzureBlobStorage = &webapps.BlobStorageTokenStore{
				SasURLSettingName: pointer.To(login.TokenBlobStorageSAS),
			}
		}
	}

	if login.LogoutEndpoint != "" {
		result.Routes = &webapps.LoginRoutes{
			LogoutEndpoint: pointer.To(login.LogoutEndpoint),
		}
	}
	result.TokenStore.TokenRefreshExtensionHours = pointer.To(login.TokenRefreshExtension)
	if login.TokenFilesystemPath != "" {
		result.TokenStore.FileSystem = &webapps.FileSystemTokenStore{
			Directory: pointer.To(login.TokenFilesystemPath),
		}
	}
	if login.TokenBlobStorageSAS != "" {
		result.TokenStore.AzureBlobStorage = &webapps.BlobStorageTokenStore{
			SasURLSettingName: pointer.To(login.TokenBlobStorageSAS),
		}
	}
	result.AllowedExternalRedirectURLs = pointer.To(login.AllowedExternalRedirectURLs)

	return result
}

func flattenAuthV2LoginSettings(input *webapps.Login) []AuthV2Login {
	if input == nil {
		return []AuthV2Login{}
	}
	result := AuthV2Login{
		PreserveURLFragmentsForLogins: pointer.From(input.PreserveURLFragmentsForLogins),
		AllowedExternalRedirectURLs:   pointer.From(input.AllowedExternalRedirectURLs),
	}
	if routes := input.Routes; routes != nil {
		result.LogoutEndpoint = pointer.From(routes.LogoutEndpoint)
	}
	if token := input.TokenStore; token != nil {
		result.TokenStoreEnabled = pointer.From(token.Enabled)
		result.TokenRefreshExtension = pointer.From(token.TokenRefreshExtensionHours)
		if fs := token.FileSystem; fs != nil {
			result.TokenFilesystemPath = pointer.From(fs.Directory)
		}
		if bs := token.AzureBlobStorage; bs != nil {
			result.TokenBlobStorageSAS = pointer.From(bs.SasURLSettingName)
		}
	}

	if nonce := input.Nonce; nonce != nil {
		result.NonceExpirationTime = pointer.From(nonce.NonceExpirationInterval)
		result.ValidateNonce = pointer.From(nonce.ValidateNonce)
	}

	if cookie := input.CookieExpiration; cookie != nil {
		result.CookieExpirationConvention = string(pointer.From(cookie.Convention))
		result.CookieExpirationTime = pointer.From(cookie.TimeToExpiration)
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
			"auth_settings_v2.0.apple_v2",
			"auth_settings_v2.0.active_directory_v2",
			"auth_settings_v2.0.azure_static_web_app_v2",
			"auth_settings_v2.0.custom_oidc_v2",
			"auth_settings_v2.0.facebook_v2",
			"auth_settings_v2.0.github_v2",
			"auth_settings_v2.0.google_v2",
			"auth_settings_v2.0.microsoft_v2",
			"auth_settings_v2.0.twitter_v2",
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
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},
			},
		},
	}
}

func AppleAuthV2SettingsSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
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

func expandAppleAuthV2Settings(input []AppleAuthV2Settings) *webapps.Apple {
	if len(input) == 1 {
		apple := input[0]
		return &webapps.Apple{
			Enabled: pointer.To(true),
			Registration: &webapps.AppleRegistration{
				ClientId:                pointer.To(apple.ClientId),
				ClientSecretSettingName: pointer.To(apple.ClientSecretSettingName),
			},
			Login: &webapps.LoginScopes{
				Scopes: pointer.To(apple.LoginScopes),
			},
		}
	}

	return &webapps.Apple{
		Enabled: pointer.To(false),
	}
}

func flattenAppleAuthV2Settings(input *webapps.Apple) []AppleAuthV2Settings {
	if input == nil || !pointer.From(input.Enabled) {
		return []AppleAuthV2Settings{}
	}
	result := AppleAuthV2Settings{}

	props := *input
	if reg := props.Registration; reg != nil {
		result.ClientId = pointer.From(reg.ClientId)
		result.ClientSecretSettingName = pointer.From(reg.ClientSecretSettingName)
	}
	if loginScopes := props.Login; loginScopes != nil {
		result.LoginScopes = pointer.From(loginScopes.Scopes)
	}

	return []AppleAuthV2Settings{result}
}

type AadAuthV2Settings struct {
	TenantAuthURI                     string            `tfschema:"tenant_auth_endpoint"` // Maps to OpenIDIssuer, takes the form `https://login.microsoftonline.com/v2.0/{tenant-guid}/`
	ClientId                          string            `tfschema:"client_id"`
	ClientSecretSettingName           string            `tfschema:"client_secret_setting_name"`
	ClientSecretCertificateThumbprint string            `tfschema:"client_secret_certificate_thumbprint"`
	LoginParameters                   map[string]string `tfschema:"login_parameters"`
	DisableWWWAuth                    bool              `tfschema:"www_authentication_disabled"`
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
			"auth_settings_v2.0.apple_v2",
			"auth_settings_v2.0.active_directory_v2",
			"auth_settings_v2.0.azure_static_web_app_v2",
			"auth_settings_v2.0.custom_oidc_v2",
			"auth_settings_v2.0.facebook_v2",
			"auth_settings_v2.0.github_v2",
			"auth_settings_v2.0.google_v2",
			"auth_settings_v2.0.microsoft_v2",
			"auth_settings_v2.0.twitter_v2",
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
					Description:  "The Azure Tenant Endpoint for the Authenticating Tenant. e.g. `https://login.microsoftonline.com/v2.0/{tenant-guid}/`.",
				},

				"client_secret_setting_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					ConflictsWith: []string{
						"auth_settings_v2.0.active_directory_v2.0.client_secret_certificate_thumbprint",
					},
					Description: "The App Setting name that contains the client secret of the Client.",
				},

				"client_secret_certificate_thumbprint": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					ConflictsWith: []string{
						"auth_settings_v2.0.active_directory_v2.0.client_secret_setting_name",
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

				"www_authentication_disabled": {
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
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"client_id": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The ID of the Client to use to authenticate with Azure Active Directory.",
				},

				"tenant_auth_endpoint": {
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

				"www_authentication_disabled": {
					Type:        pluginsdk.TypeBool,
					Computed:    true,
					Description: "Is the www-authenticate provider omitted from the request?",
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

func expandAadAuthV2Settings(input []AadAuthV2Settings) *webapps.AzureActiveDirectory {
	result := &webapps.AzureActiveDirectory{
		Enabled: pointer.To(false),
	}

	if len(input) == 1 {
		aad := input[0]
		result = &webapps.AzureActiveDirectory{
			Enabled: pointer.To(true),
			Registration: &webapps.AzureActiveDirectoryRegistration{
				OpenIdIssuer: pointer.To(aad.TenantAuthURI),
				ClientId:     pointer.To(aad.ClientId),
			},
			Login: &webapps.AzureActiveDirectoryLogin{
				DisableWWWAuthenticate: pointer.To(aad.DisableWWWAuth),
			},
		}

		if aad.ClientSecretSettingName != "" {
			result.Registration.ClientSecretSettingName = pointer.To(aad.ClientSecretSettingName)
		}

		if aad.ClientSecretCertificateThumbprint != "" {
			result.Registration.ClientSecretCertificateThumbprint = pointer.To(aad.ClientSecretCertificateThumbprint)
		}

		if len(aad.LoginParameters) > 0 {
			params := make([]string, 0)
			for k, v := range aad.LoginParameters {
				params = append(params, fmt.Sprintf("%s=%s", k, v))
			}
			result.Login.LoginParameters = &params
		}

		if len(aad.JWTAllowedGroups) != 0 || len(aad.JWTAllowedClientApps) != 0 {
			if result.Validation == nil {
				result.Validation = &webapps.AzureActiveDirectoryValidation{}
			}
			result.Validation.JwtClaimChecks = &webapps.JwtClaimChecks{}
			if len(aad.JWTAllowedGroups) != 0 {
				result.Validation.JwtClaimChecks.AllowedGroups = pointer.To(aad.JWTAllowedGroups)
			}
			if len(aad.JWTAllowedClientApps) != 0 {
				result.Validation.JwtClaimChecks.AllowedClientApplications = pointer.To(aad.JWTAllowedClientApps)
			}
		}

		if len(aad.AllowedGroups) > 0 || len(aad.AllowedIdentities) > 0 {
			if result.Validation == nil {
				result.Validation = &webapps.AzureActiveDirectoryValidation{}
			}
			result.Validation.DefaultAuthorizationPolicy = &webapps.DefaultAuthorizationPolicy{
				AllowedPrincipals: &webapps.AllowedPrincipals{},
			}
			if len(aad.AllowedGroups) > 0 {
				result.Validation.DefaultAuthorizationPolicy.AllowedPrincipals.Groups = pointer.To(aad.AllowedGroups)
			}
			if len(aad.AllowedIdentities) > 0 {
				result.Validation.DefaultAuthorizationPolicy.AllowedPrincipals.Identities = pointer.To(aad.AllowedIdentities)
			}
		}
		if len(aad.AllowedAudiences) > 0 {
			if result.Validation == nil {
				result.Validation = &webapps.AzureActiveDirectoryValidation{}
			}
			result.Validation.AllowedAudiences = pointer.To(aad.AllowedAudiences)
		}

		if len(aad.AllowedApplications) > 0 {
			if result.Validation == nil {
				result.Validation = &webapps.AzureActiveDirectoryValidation{}
			}
			if result.Validation.DefaultAuthorizationPolicy == nil {
				result.Validation.DefaultAuthorizationPolicy = &webapps.DefaultAuthorizationPolicy{}
			}
			result.Validation.DefaultAuthorizationPolicy.AllowedApplications = pointer.To(aad.AllowedApplications)
		}
	}

	return result
}

func flattenAadAuthV2Settings(input *webapps.AzureActiveDirectory) []AadAuthV2Settings {
	if input == nil || !pointer.From(input.Enabled) {
		return []AadAuthV2Settings{}
	}

	result := AadAuthV2Settings{}

	if reg := input.Registration; reg != nil {
		result.TenantAuthURI = pointer.From(reg.OpenIdIssuer)
		result.ClientId = pointer.From(reg.ClientId)
		result.ClientSecretSettingName = pointer.From(reg.ClientSecretSettingName)
		result.ClientSecretCertificateThumbprint = pointer.From(reg.ClientSecretCertificateThumbprint)
	}

	if login := input.Login; login != nil {
		result.DisableWWWAuth = pointer.From(login.DisableWWWAuthenticate)
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

	if validation := input.Validation; validation != nil {
		if validation.AllowedAudiences != nil {
			result.AllowedAudiences = *validation.AllowedAudiences
		}
		if jwt := validation.JwtClaimChecks; jwt != nil {
			result.JWTAllowedGroups = pointer.From(jwt.AllowedGroups)
			result.JWTAllowedClientApps = pointer.From(jwt.AllowedClientApplications)
		}
		if defaultPolicy := validation.DefaultAuthorizationPolicy; defaultPolicy != nil {
			result.AllowedApplications = pointer.From(defaultPolicy.AllowedApplications)
			if defaultPolicy.AllowedPrincipals != nil {
				result.AllowedGroups = pointer.From(defaultPolicy.AllowedPrincipals.Groups)
				result.AllowedIdentities = pointer.From(defaultPolicy.AllowedPrincipals.Identities)
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
			"auth_settings_v2.0.apple_v2",
			"auth_settings_v2.0.active_directory_v2",
			"auth_settings_v2.0.azure_static_web_app_v2",
			"auth_settings_v2.0.custom_oidc_v2",
			"auth_settings_v2.0.facebook_v2",
			"auth_settings_v2.0.github_v2",
			"auth_settings_v2.0.google_v2",
			"auth_settings_v2.0.microsoft_v2",
			"auth_settings_v2.0.twitter_v2",
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
		Computed: true,
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

func expandStaticWebAppAuthV2Settings(input []StaticWebAppAuthV2Settings) *webapps.AzureStaticWebApps {
	if len(input) == 1 {
		swa := input[0]
		return &webapps.AzureStaticWebApps{
			Enabled: pointer.To(true),
			Registration: &webapps.AzureStaticWebAppsRegistration{
				ClientId: pointer.To(swa.ClientId),
			},
		}
	}

	return &webapps.AzureStaticWebApps{
		Enabled: pointer.To(false),
	}
}

func flattenStaticWebAppAuthV2Settings(input *webapps.AzureStaticWebApps) []StaticWebAppAuthV2Settings {
	if input == nil || (input.Enabled != nil && !*input.Enabled) {
		return []StaticWebAppAuthV2Settings{}
	}

	result := StaticWebAppAuthV2Settings{}

	if props := input; props != nil && pointer.From(props.Enabled) {
		if props.Registration != nil {
			result.ClientId = pointer.From(props.Registration.ClientId)
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
			"auth_settings_v2.0.apple_v2",
			"auth_settings_v2.0.active_directory_v2",
			"auth_settings_v2.0.azure_static_web_app_v2",
			"auth_settings_v2.0.custom_oidc_v2",
			"auth_settings_v2.0.facebook_v2",
			"auth_settings_v2.0.github_v2",
			"auth_settings_v2.0.google_v2",
			"auth_settings_v2.0.microsoft_v2",
			"auth_settings_v2.0.twitter_v2",
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
					Description: "The Client Credential Method used. Currently the only supported value is `ClientSecretPost`.",
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
		Computed: true,
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

func expandCustomOIDCAuthV2Settings(input []CustomOIDCAuthV2Settings) map[string]webapps.CustomOpenIdConnectProvider {
	if len(input) == 0 {
		return nil
	}
	result := make(map[string]webapps.CustomOpenIdConnectProvider)
	for _, v := range input {
		if v.Name == "" {
			continue
		}
		provider := webapps.CustomOpenIdConnectProvider{
			Enabled: pointer.To(true),
			Registration: &webapps.OpenIdConnectRegistration{
				ClientId: pointer.To(v.ClientId),
				ClientCredential: &webapps.OpenIdConnectClientCredential{
					Method:                  pointer.To(webapps.ClientCredentialMethodClientSecretPost),
					ClientSecretSettingName: pointer.To(fmt.Sprintf("%s_PROVIDER_AUTHENTICATION_SECRET", strings.ToUpper(v.Name))),
				},
				OpenIdConnectConfiguration: &webapps.OpenIdConnectConfig{
					WellKnownOpenIdConfiguration: pointer.To(v.OpenIDConfigurationEndpoint),
				},
			},
			Login: &webapps.OpenIdConnectLogin{
				Scopes: pointer.To(v.Scopes),
			},
		}

		if v.NameClaimType != "" {
			provider.Login.NameClaimType = pointer.To(v.NameClaimType)
		}

		result[v.Name] = provider
	}

	return result
}

func flattenCustomOIDCAuthV2Settings(input *map[string]webapps.CustomOpenIdConnectProvider) []CustomOIDCAuthV2Settings {
	if input == nil || len(*input) == 0 {
		return []CustomOIDCAuthV2Settings{}
	}

	result := make([]CustomOIDCAuthV2Settings, 0)
	for k, v := range *input {
		if !pointer.From(v.Enabled) {
			continue
		} else {
			provider := CustomOIDCAuthV2Settings{
				Name: k,
			}
			if reg := v.Registration; reg != nil {
				provider.ClientId = pointer.From(reg.ClientId)
				if reg.ClientCredential != nil {
					provider.ClientSecretSettingName = pointer.From(reg.ClientCredential.ClientSecretSettingName)
					provider.ClientCredentialMethod = string(pointer.From(reg.ClientCredential.Method))
				}
				if config := reg.OpenIdConnectConfiguration; config != nil {
					provider.OpenIDConfigurationEndpoint = pointer.From(config.WellKnownOpenIdConfiguration)
					provider.AuthorizationEndpoint = pointer.From(config.AuthorizationEndpoint)
					provider.TokenEndpoint = pointer.From(config.TokenEndpoint)
					provider.IssuerEndpoint = pointer.From(config.Issuer)
					provider.CertificationURI = pointer.From(config.CertificationUri)
				}
			}
			if login := v.Login; login != nil {
				if login.Scopes != nil {
					provider.Scopes = *login.Scopes
				}
				provider.NameClaimType = pointer.From(login.NameClaimType)
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
			"auth_settings_v2.0.apple_v2",
			"auth_settings_v2.0.active_directory_v2",
			"auth_settings_v2.0.azure_static_web_app_v2",
			"auth_settings_v2.0.custom_oidc_v2",
			"auth_settings_v2.0.facebook_v2",
			"auth_settings_v2.0.github_v2",
			"auth_settings_v2.0.google_v2",
			"auth_settings_v2.0.microsoft_v2",
			"auth_settings_v2.0.twitter_v2",
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
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"app_id": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The App ID of the Facebook app used for login.",
				},

				"app_secret_setting_name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The app setting name that contains the `app_secret` value used for Facebook Login.",
				},

				"graph_api_version": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The version of the Facebook API to be used while logging in.",
				},

				"login_scopes": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "Specifies a list of OAuth 2.0 scopes to be requested as part of Facebook Login authentication.",
				},
			},
		},
	}
}

func expandFacebookAuthV2Settings(input []FacebookAuthV2Settings) *webapps.Facebook {
	if len(input) == 1 {
		facebook := input[0]
		result := &webapps.Facebook{
			Enabled: pointer.To(true),
			Registration: &webapps.AppRegistration{
				AppId:                pointer.To(facebook.AppId),
				AppSecretSettingName: pointer.To(facebook.AppSecretSettingName),
			},
		}

		result.GraphApiVersion = pointer.To(facebook.GraphAPIVersion)
		result.Login = &webapps.LoginScopes{
			Scopes: pointer.To(facebook.LoginScopes),
		}

		return result
	}

	return &webapps.Facebook{
		Enabled: pointer.To(false),
	}
}

func flattenFacebookAuthV2Settings(input *webapps.Facebook) []FacebookAuthV2Settings {
	if input == nil || !pointer.From(input.Enabled) {
		return []FacebookAuthV2Settings{}
	}

	result := FacebookAuthV2Settings{
		GraphAPIVersion: pointer.From(input.GraphApiVersion),
	}

	if reg := input.Registration; reg != nil {
		result.AppId = pointer.From(reg.AppId)
		result.AppSecretSettingName = pointer.From(reg.AppSecretSettingName)
	}
	if login := input.Login; login != nil {
		result.LoginScopes = pointer.From(login.Scopes)
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
			"auth_settings_v2.0.apple_v2",
			"auth_settings_v2.0.active_directory_v2",
			"auth_settings_v2.0.azure_static_web_app_v2",
			"auth_settings_v2.0.custom_oidc_v2",
			"auth_settings_v2.0.facebook_v2",
			"auth_settings_v2.0.github_v2",
			"auth_settings_v2.0.google_v2",
			"auth_settings_v2.0.microsoft_v2",
			"auth_settings_v2.0.twitter_v2",
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
		Computed: true,
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

func expandGitHubAuthV2Settings(input []GithubAuthV2Settings) *webapps.GitHub {
	if len(input) == 1 {
		github := input[0]
		result := &webapps.GitHub{
			Enabled: pointer.To(true),
			Registration: &webapps.ClientRegistration{
				ClientId:                pointer.To(github.ClientId),
				ClientSecretSettingName: pointer.To(github.ClientSecretSettingName),
			},
			Login: &webapps.LoginScopes{
				Scopes: pointer.To(github.LoginScopes),
			},
		}

		return result
	}

	return &webapps.GitHub{
		Enabled: pointer.To(false),
	}
}

func flattenGitHubAuthV2Settings(input *webapps.GitHub) []GithubAuthV2Settings {
	if input == nil || !pointer.From(input.Enabled) {
		return []GithubAuthV2Settings{}
	}

	result := GithubAuthV2Settings{}

	if reg := input.Registration; reg != nil {
		result.ClientId = pointer.From(reg.ClientId)
		result.ClientSecretSettingName = pointer.From(reg.ClientSecretSettingName)
	}
	if login := input.Login; login != nil && login.Scopes != nil {
		result.LoginScopes = pointer.From(login.Scopes)
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
			"auth_settings_v2.0.apple_v2",
			"auth_settings_v2.0.active_directory_v2",
			"auth_settings_v2.0.azure_static_web_app_v2",
			"auth_settings_v2.0.custom_oidc_v2",
			"auth_settings_v2.0.facebook_v2",
			"auth_settings_v2.0.github_v2",
			"auth_settings_v2.0.google_v2",
			"auth_settings_v2.0.microsoft_v2",
			"auth_settings_v2.0.twitter_v2",
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
		Computed: true,
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

func expandGoogleAuthV2Settings(input []GoogleAuthV2Settings) *webapps.Google {
	if len(input) == 1 {
		google := input[0]
		return &webapps.Google{
			Enabled: pointer.To(true),
			Registration: &webapps.ClientRegistration{
				ClientId:                pointer.To(google.ClientId),
				ClientSecretSettingName: pointer.To(google.ClientSecretSettingName),
			},
			Validation: &webapps.AllowedAudiencesValidation{
				AllowedAudiences: pointer.To(google.AllowedAudiences),
			},
			Login: &webapps.LoginScopes{
				Scopes: pointer.To(google.LoginScopes),
			},
		}

	}

	return &webapps.Google{
		Enabled: pointer.To(false),
	}
}

func flattenGoogleAuthV2Settings(input *webapps.Google) []GoogleAuthV2Settings {
	if input == nil || !pointer.From(input.Enabled) {
		return []GoogleAuthV2Settings{}
	}

	result := GoogleAuthV2Settings{}

	if reg := input.Registration; reg != nil {
		result.ClientId = pointer.From(reg.ClientId)
		result.ClientSecretSettingName = pointer.From(reg.ClientSecretSettingName)
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
			"auth_settings_v2.0.apple_v2",
			"auth_settings_v2.0.active_directory_v2",
			"auth_settings_v2.0.azure_static_web_app_v2",
			"auth_settings_v2.0.custom_oidc_v2",
			"auth_settings_v2.0.facebook_v2",
			"auth_settings_v2.0.github_v2",
			"auth_settings_v2.0.google_v2",
			"auth_settings_v2.0.microsoft_v2",
			"auth_settings_v2.0.twitter_v2",
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
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"client_id": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The OAuth 2.0 client ID that was created for the app used for authentication.",
				},

				"client_secret_setting_name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The app setting name containing the OAuth 2.0 client secret that was created for the app used for authentication.",
				},

				"allowed_audiences": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "Specifies a list of Allowed Audiences that will be requested as part of Microsoft Sign-In authentication.",
				},

				"login_scopes": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "The list of Login scopes that will be requested as part of Microsoft Account authentication.",
				},
			},
		},
	}
}

func expandMicrosoftAuthV2Settings(input []MicrosoftAuthV2Settings) *webapps.LegacyMicrosoftAccount {
	if len(input) == 1 {
		msft := input[0]
		return &webapps.LegacyMicrosoftAccount{
			Enabled: pointer.To(true),
			Registration: &webapps.ClientRegistration{
				ClientId:                pointer.To(msft.ClientId),
				ClientSecretSettingName: pointer.To(msft.ClientSecretSettingName),
			},
			Validation: &webapps.AllowedAudiencesValidation{
				AllowedAudiences: pointer.To(msft.AllowedAudiences),
			},
			Login: &webapps.LoginScopes{
				Scopes: pointer.To(msft.LoginScopes),
			},
		}
	}

	return &webapps.LegacyMicrosoftAccount{
		Enabled: pointer.To(false),
	}
}

func flattenMicrosoftAuthV2Settings(input *webapps.LegacyMicrosoftAccount) []MicrosoftAuthV2Settings {
	if input == nil || !pointer.From(input.Enabled) {
		return []MicrosoftAuthV2Settings{}
	}

	result := MicrosoftAuthV2Settings{}

	if reg := input.Registration; reg != nil {
		result.ClientId = pointer.From(reg.ClientId)
		result.ClientSecretSettingName = pointer.From(reg.ClientSecretSettingName)
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
			"auth_settings_v2.0.apple_v2",
			"auth_settings_v2.0.active_directory_v2",
			"auth_settings_v2.0.azure_static_web_app_v2",
			"auth_settings_v2.0.custom_oidc_v2",
			"auth_settings_v2.0.facebook_v2",
			"auth_settings_v2.0.github_v2",
			"auth_settings_v2.0.google_v2",
			"auth_settings_v2.0.microsoft_v2",
			"auth_settings_v2.0.twitter_v2",
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
		Computed: true,
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

func expandTwitterAuthV2Settings(input []TwitterAuthV2Settings) *webapps.Twitter {
	if len(input) == 1 {
		twitter := input[0]
		result := &webapps.Twitter{
			Enabled: pointer.To(true),
			Registration: &webapps.TwitterRegistration{
				ConsumerKey:               pointer.To(twitter.ConsumerKey),
				ConsumerSecretSettingName: pointer.To(twitter.ConsumerSecretSettingName),
			},
		}

		return result
	}

	return &webapps.Twitter{
		Enabled: pointer.To(false),
	}
}

func flattenTwitterAuthV2Settings(input *webapps.Twitter) []TwitterAuthV2Settings {
	if input == nil || !pointer.From(input.Enabled) {
		return []TwitterAuthV2Settings{}
	}

	if pointer.From(input.Enabled) {
		result := TwitterAuthV2Settings{}
		if reg := input.Registration; reg != nil {
			result.ConsumerKey = pointer.From(reg.ConsumerKey)
			result.ConsumerSecretSettingName = pointer.From(reg.ConsumerSecretSettingName)
		}
		return []TwitterAuthV2Settings{result}
	}

	return nil
}

func ExpandAuthV2Settings(input []AuthV2Settings) *webapps.SiteAuthSettingsV2 {
	result := &webapps.SiteAuthSettingsV2{}
	if len(input) != 1 {
		return result
	}

	settings := input[0]

	props := &webapps.SiteAuthSettingsV2Properties{
		Platform: &webapps.AuthPlatform{
			Enabled:        pointer.To(settings.AuthEnabled),
			RuntimeVersion: pointer.To(settings.RuntimeVersion),
		},
		GlobalValidation: &webapps.GlobalValidation{
			RequireAuthentication:       pointer.To(settings.RequireAuth),
			UnauthenticatedClientAction: pointer.To(webapps.UnauthenticatedClientActionV2(settings.UnauthenticatedAction)),
			ExcludedPaths:               pointer.To(settings.ExcludedPaths),
		},
		IdentityProviders: &webapps.IdentityProviders{
			AzureActiveDirectory:         expandAadAuthV2Settings(settings.AzureActiveDirectoryAuth),
			Facebook:                     expandFacebookAuthV2Settings(settings.FacebookAuth),
			GitHub:                       expandGitHubAuthV2Settings(settings.GithubAuth),
			Google:                       expandGoogleAuthV2Settings(settings.GoogleAuth),
			Twitter:                      expandTwitterAuthV2Settings(settings.TwitterAuth),
			CustomOpenIdConnectProviders: pointer.To(expandCustomOIDCAuthV2Settings(settings.CustomOIDCAuth)),
			LegacyMicrosoftAccount:       expandMicrosoftAuthV2Settings(settings.MicrosoftAuth),
			Apple:                        expandAppleAuthV2Settings(settings.AppleAuth),
			AzureStaticWebApps:           expandStaticWebAppAuthV2Settings(settings.AzureStaticWebAuth),
		},
		Login: expandAuthV2LoginSettings(settings.Login),
		HTTPSettings: &webapps.HTTPSettings{
			RequireHTTPS: pointer.To(settings.RequireHTTPS),
			Routes: &webapps.HTTPSettingsRoutes{
				ApiPrefix: pointer.To(settings.HttpRoutesAPIPrefix),
			},
			ForwardProxy: &webapps.ForwardProxy{
				Convention: pointer.To(webapps.ForwardProxyConvention(settings.ForwardProxyConvention)),
			},
		},
	}

	// Platform
	if settings.ConfigFilePath != "" {
		props.Platform.ConfigFilePath = pointer.To(settings.ConfigFilePath)
	}

	// Global
	if settings.DefaultAuthProvider != "" {
		props.GlobalValidation.RedirectToProvider = pointer.To(settings.DefaultAuthProvider)
	}

	// HTTP
	if settings.ForwardProxyCustomHostHeaderName != "" {
		props.HTTPSettings.ForwardProxy.CustomHostHeaderName = pointer.To(settings.ForwardProxyCustomHostHeaderName)
	}
	if settings.ForwardProxyCustomSchemeHeaderName != "" {
		props.HTTPSettings.ForwardProxy.CustomProtoHeaderName = pointer.To(settings.ForwardProxyCustomSchemeHeaderName)
	}

	result.Properties = props

	return result
}

func FlattenAuthV2Settings(input webapps.SiteAuthSettingsV2) []AuthV2Settings {
	if input.Properties == nil {
		return []AuthV2Settings{}
	}

	settings := *input.Properties

	result := AuthV2Settings{}

	if platform := settings.Platform; platform != nil {
		result.AuthEnabled = pointer.From(platform.Enabled)
		result.RuntimeVersion = pointer.From(platform.RuntimeVersion)
		result.ConfigFilePath = pointer.From(platform.ConfigFilePath)
	}

	if global := settings.GlobalValidation; global != nil {
		result.RequireAuth = pointer.From(global.RequireAuthentication)
		result.UnauthenticatedAction = string(pointer.From(global.UnauthenticatedClientAction))
		result.DefaultAuthProvider = pointer.From(global.RedirectToProvider)
		result.ExcludedPaths = pointer.From(global.ExcludedPaths)
	}

	if http := settings.HTTPSettings; http != nil {
		result.RequireHTTPS = pointer.From(http.RequireHTTPS)
		if http.Routes != nil {
			result.HttpRoutesAPIPrefix = pointer.From(http.Routes.ApiPrefix)
		}
		if fp := http.ForwardProxy; fp != nil {
			result.ForwardProxyConvention = string(pointer.From(fp.Convention))
			result.ForwardProxyCustomHostHeaderName = pointer.From(fp.CustomHostHeaderName)
			result.ForwardProxyCustomSchemeHeaderName = pointer.From(fp.CustomProtoHeaderName)
		}
	}

	if login := settings.Login; login != nil {
		result.Login = flattenAuthV2LoginSettings(login)
	}

	if authProviders := settings.IdentityProviders; authProviders != nil {
		result.AppleAuth = flattenAppleAuthV2Settings(authProviders.Apple)
		result.AzureActiveDirectoryAuth = flattenAadAuthV2Settings(authProviders.AzureActiveDirectory)
		result.AzureStaticWebAuth = flattenStaticWebAppAuthV2Settings(authProviders.AzureStaticWebApps)
		result.CustomOIDCAuth = flattenCustomOIDCAuthV2Settings(authProviders.CustomOpenIdConnectProviders)
		result.FacebookAuth = flattenFacebookAuthV2Settings(authProviders.Facebook)
		result.GithubAuth = flattenGitHubAuthV2Settings(authProviders.GitHub)
		result.GoogleAuth = flattenGoogleAuthV2Settings(authProviders.Google)
		result.MicrosoftAuth = flattenMicrosoftAuthV2Settings(authProviders.LegacyMicrosoftAccount)
		result.TwitterAuth = flattenTwitterAuthV2Settings(authProviders.Twitter)
	}

	return []AuthV2Settings{result}
}
