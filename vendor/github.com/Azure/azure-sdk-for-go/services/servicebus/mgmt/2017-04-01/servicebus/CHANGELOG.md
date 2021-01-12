Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewSBAuthorizationRuleListResultPage` parameter(s) have been changed from `(func(context.Context, SBAuthorizationRuleListResult) (SBAuthorizationRuleListResult, error))` to `(SBAuthorizationRuleListResult, func(context.Context, SBAuthorizationRuleListResult) (SBAuthorizationRuleListResult, error))`
- Function `NewSBTopicListResultPage` parameter(s) have been changed from `(func(context.Context, SBTopicListResult) (SBTopicListResult, error))` to `(SBTopicListResult, func(context.Context, SBTopicListResult) (SBTopicListResult, error))`
- Function `NewRuleListResultPage` parameter(s) have been changed from `(func(context.Context, RuleListResult) (RuleListResult, error))` to `(RuleListResult, func(context.Context, RuleListResult) (RuleListResult, error))`
- Function `NewEventHubListResultPage` parameter(s) have been changed from `(func(context.Context, EventHubListResult) (EventHubListResult, error))` to `(EventHubListResult, func(context.Context, EventHubListResult) (EventHubListResult, error))`
- Function `NewPremiumMessagingRegionsListResultPage` parameter(s) have been changed from `(func(context.Context, PremiumMessagingRegionsListResult) (PremiumMessagingRegionsListResult, error))` to `(PremiumMessagingRegionsListResult, func(context.Context, PremiumMessagingRegionsListResult) (PremiumMessagingRegionsListResult, error))`
- Function `NewArmDisasterRecoveryListResultPage` parameter(s) have been changed from `(func(context.Context, ArmDisasterRecoveryListResult) (ArmDisasterRecoveryListResult, error))` to `(ArmDisasterRecoveryListResult, func(context.Context, ArmDisasterRecoveryListResult) (ArmDisasterRecoveryListResult, error))`
- Function `NewSBNamespaceListResultPage` parameter(s) have been changed from `(func(context.Context, SBNamespaceListResult) (SBNamespaceListResult, error))` to `(SBNamespaceListResult, func(context.Context, SBNamespaceListResult) (SBNamespaceListResult, error))`
- Function `NewSBSubscriptionListResultPage` parameter(s) have been changed from `(func(context.Context, SBSubscriptionListResult) (SBSubscriptionListResult, error))` to `(SBSubscriptionListResult, func(context.Context, SBSubscriptionListResult) (SBSubscriptionListResult, error))`
- Function `NewNetworkRuleSetListResultPage` parameter(s) have been changed from `(func(context.Context, NetworkRuleSetListResult) (NetworkRuleSetListResult, error))` to `(NetworkRuleSetListResult, func(context.Context, NetworkRuleSetListResult) (NetworkRuleSetListResult, error))`
- Function `NewOperationListResultPage` parameter(s) have been changed from `(func(context.Context, OperationListResult) (OperationListResult, error))` to `(OperationListResult, func(context.Context, OperationListResult) (OperationListResult, error))`
- Function `NewMigrationConfigListResultPage` parameter(s) have been changed from `(func(context.Context, MigrationConfigListResult) (MigrationConfigListResult, error))` to `(MigrationConfigListResult, func(context.Context, MigrationConfigListResult) (MigrationConfigListResult, error))`
- Function `NewSBQueueListResultPage` parameter(s) have been changed from `(func(context.Context, SBQueueListResult) (SBQueueListResult, error))` to `(SBQueueListResult, func(context.Context, SBQueueListResult) (SBQueueListResult, error))`
- Struct `AuthorizationRuleProperties` has been removed
- Field `Code` of struct `ErrorResponse` has been removed
- Field `Message` of struct `ErrorResponse` has been removed

## New Content

- New struct `ErrorAdditionalInfo`
- New struct `ErrorResponseError`
- New field `Error` in struct `ErrorResponse`
