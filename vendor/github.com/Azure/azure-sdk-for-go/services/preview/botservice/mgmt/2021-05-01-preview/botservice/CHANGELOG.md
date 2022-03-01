# Change History

## Breaking Changes

### Struct Changes

#### Removed Struct Fields

1. BotProperties.IsIsolated
1. WebChatSite.EnablePreview

### Signature Changes

#### Funcs

1. ChannelsClient.ListWithKeys
	- Returns
		- From: BotChannel, error
		- To: ListChannelWithKeysResponse, error
1. ChannelsClient.ListWithKeysResponder
	- Returns
		- From: BotChannel, error
		- To: ListChannelWithKeysResponse, error

## Additive Changes

### New Constants

1. PublicNetworkAccess.PublicNetworkAccessDisabled
1. PublicNetworkAccess.PublicNetworkAccessEnabled

### New Funcs

1. *ListChannelWithKeysResponse.UnmarshalJSON([]byte) error
1. ListChannelWithKeysResponse.MarshalJSON() ([]byte, error)
1. PossiblePublicNetworkAccessValues() []PublicNetworkAccess
1. Site.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. ChannelSettings
1. ListChannelWithKeysResponse
1. ServiceProviderParameterMetadata
1. ServiceProviderParameterMetadataConstraints
1. Site

#### New Struct Fields

1. AlexaChannel.Etag
1. AlexaChannel.Location
1. AlexaChannel.ProvisioningState
1. Bot.Zones
1. BotChannel.Zones
1. BotProperties.AllSettings
1. BotProperties.CmekEncryptionStatus
1. BotProperties.IsDeveloperAppInsightsAPIKeySet
1. BotProperties.IsStreamingSupported
1. BotProperties.ManifestURL
1. BotProperties.MigrationToken
1. BotProperties.Parameters
1. BotProperties.ProvisioningState
1. BotProperties.PublicNetworkAccess
1. BotProperties.PublishingCredentials
1. BotProperties.StorageResourceID
1. Channel.Etag
1. Channel.Location
1. Channel.ProvisioningState
1. ConnectionSetting.Zones
1. ConnectionSettingProperties.ID
1. ConnectionSettingProperties.Name
1. ConnectionSettingProperties.ProvisioningState
1. DirectLineChannel.Etag
1. DirectLineChannel.Location
1. DirectLineChannel.ProvisioningState
1. DirectLineChannelProperties.DirectLineEmbedCode
1. DirectLineSite.IsBlockUserUploadEnabled
1. DirectLineSpeechChannel.Etag
1. DirectLineSpeechChannel.Location
1. DirectLineSpeechChannel.ProvisioningState
1. EmailChannel.Etag
1. EmailChannel.Location
1. EmailChannel.ProvisioningState
1. FacebookChannel.Etag
1. FacebookChannel.Location
1. FacebookChannel.ProvisioningState
1. KikChannel.Etag
1. KikChannel.Location
1. KikChannel.ProvisioningState
1. LineChannel.Etag
1. LineChannel.Location
1. LineChannel.ProvisioningState
1. MsTeamsChannel.Etag
1. MsTeamsChannel.Location
1. MsTeamsChannel.ProvisioningState
1. MsTeamsChannelProperties.AcceptedTerms
1. MsTeamsChannelProperties.DeploymentEnvironment
1. MsTeamsChannelProperties.IncomingCallRoute
1. Resource.Zones
1. ServiceProviderParameter.Metadata
1. SkypeChannel.Etag
1. SkypeChannel.Location
1. SkypeChannel.ProvisioningState
1. SkypeChannelProperties.IncomingCallRoute
1. SlackChannel.Etag
1. SlackChannel.Location
1. SlackChannel.ProvisioningState
1. SmsChannel.Etag
1. SmsChannel.Location
1. SmsChannel.ProvisioningState
1. TelegramChannel.Etag
1. TelegramChannel.Location
1. TelegramChannel.ProvisioningState
1. WebChatChannel.Etag
1. WebChatChannel.Location
1. WebChatChannel.ProvisioningState
1. WebChatSite.IsWebchatPreviewEnabled
