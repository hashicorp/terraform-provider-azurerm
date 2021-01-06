Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewEventDataCollectionPage` parameter(s) have been changed from `(func(context.Context, EventDataCollection) (EventDataCollection, error))` to `(EventDataCollection, func(context.Context, EventDataCollection) (EventDataCollection, error))`
- Function `NewAutoscaleSettingResourceCollectionPage` parameter(s) have been changed from `(func(context.Context, AutoscaleSettingResourceCollection) (AutoscaleSettingResourceCollection, error))` to `(AutoscaleSettingResourceCollection, func(context.Context, AutoscaleSettingResourceCollection) (AutoscaleSettingResourceCollection, error))`

## New Content

- New field `SkipMetricValidation` in struct `MultiMetricCriteria`
- New field `SkipMetricValidation` in struct `MetricCriteria`
- New field `SkipMetricValidation` in struct `DynamicMetricCriteria`
