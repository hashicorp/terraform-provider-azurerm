package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SiteAuthSettingsProperties struct {
	AadClaimsAuthorization                  *string                        `json:"aadClaimsAuthorization,omitempty"`
	AdditionalLoginParams                   *[]string                      `json:"additionalLoginParams,omitempty"`
	AllowedAudiences                        *[]string                      `json:"allowedAudiences,omitempty"`
	AllowedExternalRedirectURLs             *[]string                      `json:"allowedExternalRedirectUrls,omitempty"`
	AuthFilePath                            *string                        `json:"authFilePath,omitempty"`
	ClientId                                *string                        `json:"clientId,omitempty"`
	ClientSecret                            *string                        `json:"clientSecret,omitempty"`
	ClientSecretCertificateThumbprint       *string                        `json:"clientSecretCertificateThumbprint,omitempty"`
	ClientSecretSettingName                 *string                        `json:"clientSecretSettingName,omitempty"`
	ConfigVersion                           *string                        `json:"configVersion,omitempty"`
	DefaultProvider                         *BuiltInAuthenticationProvider `json:"defaultProvider,omitempty"`
	Enabled                                 *bool                          `json:"enabled,omitempty"`
	FacebookAppId                           *string                        `json:"facebookAppId,omitempty"`
	FacebookAppSecret                       *string                        `json:"facebookAppSecret,omitempty"`
	FacebookAppSecretSettingName            *string                        `json:"facebookAppSecretSettingName,omitempty"`
	FacebookOAuthScopes                     *[]string                      `json:"facebookOAuthScopes,omitempty"`
	GitHubClientId                          *string                        `json:"gitHubClientId,omitempty"`
	GitHubClientSecret                      *string                        `json:"gitHubClientSecret,omitempty"`
	GitHubClientSecretSettingName           *string                        `json:"gitHubClientSecretSettingName,omitempty"`
	GitHubOAuthScopes                       *[]string                      `json:"gitHubOAuthScopes,omitempty"`
	GoogleClientId                          *string                        `json:"googleClientId,omitempty"`
	GoogleClientSecret                      *string                        `json:"googleClientSecret,omitempty"`
	GoogleClientSecretSettingName           *string                        `json:"googleClientSecretSettingName,omitempty"`
	GoogleOAuthScopes                       *[]string                      `json:"googleOAuthScopes,omitempty"`
	IsAuthFromFile                          *string                        `json:"isAuthFromFile,omitempty"`
	Issuer                                  *string                        `json:"issuer,omitempty"`
	MicrosoftAccountClientId                *string                        `json:"microsoftAccountClientId,omitempty"`
	MicrosoftAccountClientSecret            *string                        `json:"microsoftAccountClientSecret,omitempty"`
	MicrosoftAccountClientSecretSettingName *string                        `json:"microsoftAccountClientSecretSettingName,omitempty"`
	MicrosoftAccountOAuthScopes             *[]string                      `json:"microsoftAccountOAuthScopes,omitempty"`
	RuntimeVersion                          *string                        `json:"runtimeVersion,omitempty"`
	TokenRefreshExtensionHours              *float64                       `json:"tokenRefreshExtensionHours,omitempty"`
	TokenStoreEnabled                       *bool                          `json:"tokenStoreEnabled,omitempty"`
	TwitterConsumerKey                      *string                        `json:"twitterConsumerKey,omitempty"`
	TwitterConsumerSecret                   *string                        `json:"twitterConsumerSecret,omitempty"`
	TwitterConsumerSecretSettingName        *string                        `json:"twitterConsumerSecretSettingName,omitempty"`
	UnauthenticatedClientAction             *UnauthenticatedClientAction   `json:"unauthenticatedClientAction,omitempty"`
	ValidateIssuer                          *bool                          `json:"validateIssuer,omitempty"`
}
