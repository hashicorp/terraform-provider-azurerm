# Change History

## Breaking Changes

### Removed Constants

1. AccessRights.Listen
1. AccessRights.Manage
1. AccessRights.SendEnumValue
1. NamespaceType.Messaging
1. NamespaceType.NotificationHub
1. SkuName.Basic
1. SkuName.Free
1. SkuName.Standard

### Signature Changes

#### Funcs

1. NamespacesClient.ListKeys
	- Returns
		- From: SharedAccessAuthorizationRuleListResult, error
		- To: ResourceListKeys, error
1. NamespacesClient.ListKeysResponder
	- Returns
		- From: SharedAccessAuthorizationRuleListResult, error
		- To: ResourceListKeys, error

## Additive Changes

### New Constants

1. AccessRights.AccessRightsListen
1. AccessRights.AccessRightsManage
1. AccessRights.AccessRightsSend
1. NamespaceType.NamespaceTypeMessaging
1. NamespaceType.NamespaceTypeNotificationHub
1. SkuName.SkuNameBasic
1. SkuName.SkuNameFree
1. SkuName.SkuNameStandard
