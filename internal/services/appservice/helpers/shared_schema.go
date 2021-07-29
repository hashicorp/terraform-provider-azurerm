package helpers

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-01-15/web"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	msiParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/parse"
	msiValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/validate"
	networkValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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
		Type:       pluginsdk.TypeList,
		Optional:   true,
		Computed:   true,
		ConfigMode: pluginsdk.SchemaConfigModeAttr,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"ip_address": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.Any(
						validate.IPv4Address,
						validate.CIDR,
					),
				},

				"service_tag": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"virtual_network_subnet_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: networkValidate.SubnetID,
				},

				"name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"priority": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      65000,
					ValidateFunc: validation.IntBetween(1, 2147483647),
				},

				"action": {
					Type:     pluginsdk.TypeString,
					Default:  "Allow",
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"Allow",
						"Deny",
					}, false),
				},

				"headers": IpRestrictionHeadersSchema(),
			},
		},
	}
}

func IpRestrictionSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:       pluginsdk.TypeList,
		Optional:   true,
		Computed:   true,
		ConfigMode: pluginsdk.SchemaConfigModeAttr,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"ip_address": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"service_tag": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"virtual_network_subnet_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"priority": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},

				"action": {
					Type:     pluginsdk.TypeString,
					Computed: true,
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
		Computed:   true,
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
				},

				"x_forwarded_for": {
					Type:     pluginsdk.TypeList,
					MaxItems: 8,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.IsCIDR,
					},
				},

				"x_azure_fdid": { // Front Door ID (UUID)
					Type:     pluginsdk.TypeList,
					MaxItems: 8,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.IsUUID,
					},
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
		Type:       pluginsdk.TypeList,
		Computed:   true,
		ConfigMode: pluginsdk.SchemaConfigModeAttr,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"x_forwarded_host": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"x_forwarded_for": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"x_azure_fdid": { // Front Door ID (UUID)
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"x_fd_health_probe": { // 1 or absent
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

type Identity struct {
	IdentityIds []string `tfschema:"identity_ids"`
	Type        string   `tfschema:"type"`
	PrincipalId string   `tfschema:"principal_id"`
	TenantId    string   `tfschema:"tenant_id"`
}

func IdentitySchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"identity_ids": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MinItems: 1,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: msiValidate.UserAssignedIdentityID,
					},
				},

				"type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.ManagedServiceIdentityTypeNone),
						string(web.ManagedServiceIdentityTypeSystemAssigned),
						string(web.ManagedServiceIdentityTypeSystemAssignedUserAssigned),
						string(web.ManagedServiceIdentityTypeUserAssigned),
					}, true),
				},

				"principal_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"tenant_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func IdentitySchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"identity_ids": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"type": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"principal_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"tenant_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
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
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"allowed_origins": {
					Type:     pluginsdk.TypeSet,
					Required: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"support_credentials": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
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
					Type:     pluginsdk.TypeSet,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"support_credentials": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},
			},
		},
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

func SiteCredentialSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"password": {
					Type:      pluginsdk.TypeString,
					Computed:  true,
					Sensitive: true,
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
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"enabled": {
					Type:     pluginsdk.TypeBool,
					Required: true,
				},

				"additional_login_parameters": {
					Type:     pluginsdk.TypeMap,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"allowed_external_redirect_urls": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},

				"default_provider": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.BuiltInAuthenticationProviderAzureActiveDirectory),
						string(web.BuiltInAuthenticationProviderFacebook),
						string(web.BuiltInAuthenticationProviderGithub),
						string(web.BuiltInAuthenticationProviderGoogle),
						string(web.BuiltInAuthenticationProviderMicrosoftAccount),
						string(web.BuiltInAuthenticationProviderTwitter),
					}, false),
				},

				"issuer": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
				},

				"runtime_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
				},

				"token_refresh_extension_hours": {
					Type:     pluginsdk.TypeFloat,
					Optional: true,
					Default:  72,
				},

				"token_store_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"unauthenticated_client_action": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.UnauthenticatedClientActionAllowAnonymous),
						string(web.UnauthenticatedClientActionRedirectToLoginPage),
					}, false),
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
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"additional_login_parameters": {
					Type:     pluginsdk.TypeMap,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"allowed_external_redirect_urls": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"default_provider": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"issuer": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"runtime_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"token_refresh_extension_hours": {
					Type:     pluginsdk.TypeFloat,
					Computed: true,
				},

				"token_store_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"unauthenticated_client_action": {
					Type:     pluginsdk.TypeString,
					Computed: true,
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
				},

				"client_secret": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
					ExactlyOneOf: []string{
						"auth_settings.0.active_directory.0.client_secret",
						"auth_settings.0.active_directory.0.client_secret_setting_name",
					},
				},

				"client_secret_setting_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					ExactlyOneOf: []string{
						"auth_settings.0.active_directory.0.client_secret",
						"auth_settings.0.active_directory.0.client_secret_setting_name",
					},
				},

				"allowed_audiences": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},
			},
		},
	}
}

func AadAuthSettingsSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"client_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"client_secret": {
					Type:      pluginsdk.TypeString,
					Computed:  true,
					Sensitive: true,
				},

				"client_secret_setting_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"allowed_audiences": {
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
				},

				"app_secret_setting_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					ExactlyOneOf: []string{
						"auth_settings.0.facebook.0.app_secret",
						"auth_settings.0.facebook.0.app_secret_setting_name",
					},
				},

				"oauth_scopes": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
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
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"app_secret": {
					Type:      pluginsdk.TypeString,
					Computed:  true,
					Sensitive: true,
				},

				"app_secret_setting_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"oauth_scopes": {
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

type GoogleAuthSettings struct {
	ClientId                string   `tfschema:"client_id"`
	ClientSecret            string   `tfschema:"client_schema"`
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
				},

				"client_secret_setting_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					ExactlyOneOf: []string{
						"auth_settings.0.google.0.client_secret",
						"auth_settings.0.google.0.client_secret_setting_name",
					},
				},

				"oauth_scopes": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
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
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"client_secret": {
					Type:      pluginsdk.TypeString,
					Computed:  true,
					Sensitive: true,
				},

				"client_secret_setting_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"oauth_scopes": {
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

type MicrosoftAuthSettings struct {
	ClientId                string   `tfschema:"client_id"`
	ClientSecret            string   `tfschema:"client_schema"`
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
				},

				"client_secret_setting_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					ExactlyOneOf: []string{
						"auth_settings.0.microsoft.0.client_secret",
						"auth_settings.0.microsoft.0.client_secret_setting_name",
					},
				},

				"oauth_scopes": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
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
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"client_secret": {
					Type:      pluginsdk.TypeString,
					Computed:  true,
					Sensitive: true,
				},

				"client_secret_setting_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"oauth_scopes": {
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
					Type:     pluginsdk.TypeString,
					Required: true,
				},

				"consumer_secret": {
					Type:      pluginsdk.TypeString,
					Optional:  true,
					Sensitive: true,
					ExactlyOneOf: []string{
						"auth_settings.0.twitter.0.consumer_secret",
						"auth_settings.0.twitter.0.consumer_secret_setting_name",
					},
				},

				"consumer_secret_setting_name": {
					Type:     pluginsdk.TypeString,
					Optional: true,
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
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"consumer_secret": {
					Type:      pluginsdk.TypeString,
					Computed:  true,
					Sensitive: true,
				},

				"consumer_secret_setting_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
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
					Type:     pluginsdk.TypeString,
					Required: true,
				},

				"client_secret": {
					Type:      pluginsdk.TypeString,
					Optional:  true,
					Sensitive: true,
					ExactlyOneOf: []string{
						"auth_settings.0.github.0.client_secret",
						"auth_settings.0.github.0.client_secret_setting_name",
					},
				},

				"client_secret_setting_name": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ExactlyOneOf: []string{
						"auth_settings.0.github.0.client_secret",
						"auth_settings.0.github.0.client_secret_setting_name",
					},
				},

				"oauth_scopes": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
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
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"client_secret": {
					Type:      pluginsdk.TypeString,
					Computed:  true,
					Sensitive: true,
				},

				"client_secret_setting_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"oauth_scopes": {
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

func ExpandIpRestrictions(restrictions []IpRestriction) (*[]web.IPSecurityRestriction, error) {
	var expanded []web.IPSecurityRestriction
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
			result["x-azure-fd-id"] = v.XAzureFDID
		}
		if len(v.XFDHealthProbe) > 0 {
			result["x-fd-healthprobe"] = v.XFDHealthProbe
		}
	}
	return result
}

func ExpandCorsSettings(input []CorsSetting) *web.CorsSettings {
	if len(input) == 0 {
		return nil
	}
	var result web.CorsSettings
	for _, v := range input {
		if v.SupportCredentials {
			result.SupportCredentials = utils.Bool(v.SupportCredentials)
		}

		result.AllowedOrigins = &v.AllowedOrigins
	}
	return &result
}

func ExpandIdentity(identities []Identity) *web.ManagedServiceIdentity {
	if len(identities) == 0 {
		return nil
	}
	var result web.ManagedServiceIdentity
	for _, v := range identities {
		result.Type = web.ManagedServiceIdentityType(v.Type)
		if result.Type == web.ManagedServiceIdentityTypeUserAssigned || result.Type == web.ManagedServiceIdentityTypeSystemAssignedUserAssigned {
			identityIds := make(map[string]*web.UserAssignedIdentity)
			for _, i := range v.IdentityIds {
				identityIds[i] = &web.UserAssignedIdentity{}
			}
			result.UserAssignedIdentities = identityIds
		}
	}
	return &result
}

func ExpandAuthSettings(auth []AuthSettings) *web.SiteAuthSettings {
	if len(auth) == 0 {
		return nil
	}

	props := &web.SiteAuthSettingsProperties{}

	for _, v := range auth {
		if v.Enabled {
			props.Enabled = utils.Bool(v.Enabled)
		}
		if len(v.AdditionalLoginParameters) > 0 {
			var additionalLoginParams []string
			for k, s := range v.AdditionalLoginParameters {
				additionalLoginParams = append(additionalLoginParams, fmt.Sprintf("%s=%s", k, s))
			}
			props.AdditionalLoginParams = &additionalLoginParams
		}

		if len(v.AllowedExternalRedirectUrls) != 0 && v.AllowedExternalRedirectUrls != nil {
			props.AllowedExternalRedirectUrls = &v.AllowedExternalRedirectUrls
		}

		if v.DefaultProvider != "" {
			props.DefaultProvider = web.BuiltInAuthenticationProvider(v.DefaultProvider)
		}

		if v.Issuer != "" {
			props.Issuer = utils.String(v.Issuer)
		}

		if v.RuntimeVersion != "" {
			props.RuntimeVersion = utils.String(v.RuntimeVersion)
		}

		if v.TokenRefreshExtensionHours != 0 {
			props.TokenRefreshExtensionHours = utils.Float(v.TokenRefreshExtensionHours)
		}

		if v.UnauthenticatedClientAction != "" {
			props.UnauthenticatedClientAction = web.UnauthenticatedClientAction(v.UnauthenticatedClientAction)
		}

		if len(v.AzureActiveDirectoryAuth) == 1 {
			a := v.AzureActiveDirectoryAuth[0]
			props.ClientID = utils.String(a.ClientId)

			if a.ClientSecret != "" {
				props.ClientSecret = utils.String(a.ClientSecret)
			}

			if a.ClientSecretSettingName != "" {
				props.ClientSecretSettingName = utils.String(a.ClientSecretSettingName)
			}

			props.AllowedAudiences = &a.AllowedAudiences
		}

		if len(v.FacebookAuth) == 1 {
			f := v.FacebookAuth[0]
			props.FacebookAppID = utils.String(f.AppId)

			if f.AppSecret != "" {
				props.FacebookAppSecret = utils.String(f.AppSecret)
			}

			if f.AppSecretSettingName != "" {
				props.FacebookAppSecretSettingName = utils.String(f.AppSecretSettingName)
			}

			props.FacebookOAuthScopes = &f.OauthScopes
		}

		if len(v.GithubAuth) == 1 {
			g := v.GithubAuth[0]
			props.GitHubClientID = utils.String(g.ClientId)
			if g.ClientSecret != "" {
				props.GitHubClientID = utils.String(g.ClientId)
			}

			if g.ClientSecretSettingName != "" {
				props.GitHubClientSecretSettingName = utils.String(g.ClientSecretSettingName)
			}

			props.GitHubOAuthScopes = &g.OAuthScopes
		}

		if len(v.GoogleAuth) == 1 {
			g := v.GoogleAuth[0]
			props.GoogleClientID = utils.String(g.ClientId)

			if g.ClientSecret != "" {
				props.GoogleClientSecret = utils.String(g.ClientSecret)
			}

			if g.ClientSecretSettingName != "" {
				props.GoogleClientSecretSettingName = utils.String(g.ClientSecretSettingName)
			}

			props.GoogleOAuthScopes = &g.OauthScopes
		}

		if len(v.MicrosoftAuth) == 1 {
			m := v.MicrosoftAuth[0]
			props.MicrosoftAccountClientID = utils.String(m.ClientId)

			if m.ClientSecret != "" {
				props.MicrosoftAccountClientSecret = utils.String(m.ClientSecret)
			}

			if m.ClientSecretSettingName != "" {
				props.MicrosoftAccountClientSecretSettingName = utils.String(m.ClientSecretSettingName)
			}

			props.MicrosoftAccountOAuthScopes = &m.OauthScopes
		}

		if len(v.TwitterAuth) == 1 {
			t := v.TwitterAuth[0]
			props.TwitterConsumerKey = utils.String(t.ConsumerKey)

			if t.ConsumerSecret != "" {
				props.TwitterConsumerSecret = utils.String(t.ConsumerSecret)
			}

			if t.ConsumerSecretSettingName != "" {
				props.TwitterConsumerSecretSettingName = utils.String(t.ConsumerSecretSettingName)
			}
		}
	}

	return &web.SiteAuthSettings{
		SiteAuthSettingsProperties: props,
	}
}

func FlattenAuthSettings(auth web.SiteAuthSettings) []AuthSettings {
	if auth.SiteAuthSettingsProperties == nil {
		return nil
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

	if props.AllowedExternalRedirectUrls != nil {
		result.AllowedExternalRedirectUrls = *props.AllowedExternalRedirectUrls
	}

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

func FlattenIdentity(appIdentity *web.ManagedServiceIdentity) []Identity {
	if appIdentity == nil {
		return nil
	}
	identity := Identity{
		Type: string(appIdentity.Type),
	}

	if len(appIdentity.UserAssignedIdentities) != 0 {
		var identityIds []string
		for k := range appIdentity.UserAssignedIdentities {
			// Service can return broken case IDs, so we normalise here and discard invalid entries
			id, err := msiParse.UserAssignedIdentityID(k)
			if err == nil {
				identityIds = append(identityIds, id.ID())
			}
		}
		identity.IdentityIds = identityIds
	}

	if appIdentity.PrincipalID != nil {
		identity.PrincipalId = *appIdentity.PrincipalID
	}

	if appIdentity.TenantID != nil {
		identity.TenantId = *appIdentity.TenantID
	}

	return []Identity{identity}
}

func FlattenIpRestrictions(ipRestrictionsList *[]web.IPSecurityRestriction) []IpRestriction {
	if ipRestrictionsList == nil {
		return nil
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
			ipRestriction.IpAddress = *v.IPAddress
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
		return nil
	}
	ipRestrictionHeader := IpRestrictionHeaders{}
	if xForwardFor, ok := headers["x-forwarded-for"]; ok {
		ipRestrictionHeader.XForwardedFor = xForwardFor
	}

	if xForwardedHost, ok := headers["x-forwarded-host"]; ok {
		ipRestrictionHeader.XForwardedHost = xForwardedHost
	}

	if xAzureFDID, ok := headers["x-azure-fd-id"]; ok {
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
