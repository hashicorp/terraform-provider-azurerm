Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewRecordSetListResultPage` parameter(s) have been changed from `(func(context.Context, RecordSetListResult) (RecordSetListResult, error))` to `(RecordSetListResult, func(context.Context, RecordSetListResult) (RecordSetListResult, error))`
- Function `NewZoneListResultPage` parameter(s) have been changed from `(func(context.Context, ZoneListResult) (ZoneListResult, error))` to `(ZoneListResult, func(context.Context, ZoneListResult) (ZoneListResult, error))`

## New Content

- New field `MaxNumberOfRecordsPerRecordSet` in struct `ZoneProperties`
