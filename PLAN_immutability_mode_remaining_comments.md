# Remaining work for feature/mssql_database-immutability_mode

This file tracks the outstanding PR comments to resolve after the SQL API version
bump PR (feature/mssql_sql_api_version_bump) merges and this branch rebases on top.

---

## Comment 2 — Remove `immutable_backups_enabled`, expose only `immutability_mode`
**File:** `internal/services/mssql/helper/sql_retention_policies.go` lines 63–75

### What to do
- Remove the `immutable_backups_enabled` schema field entirely.
- Keep only `immutability_mode` (Optional, Computed, ValidateFunc for Locked/Unlocked).
- In `ExpandLongTermRetentionPolicy`: infer `TimeBasedImmutability` from whether
  `immutability_mode` is set:
  - mode is non-empty  → `TimeBasedImmutability = Enabled`, send `TimeBasedImmutabilityMode`
  - mode is empty/unset → `TimeBasedImmutability = Disabled`, omit mode
- In `FlattenLongTermRetentionPolicy`: remove `immutable_backups_enabled` from the
  returned map; derive it from `TimeBasedImmutabilityMode` being non-empty.
- Remove `RequiredWith` constraint (no longer needed since `immutable_backups_enabled`
  is gone).
- Update `mssql_database_resource_test.go`: replace all references to
  `immutable_backups_enabled` with just `immutability_mode`.
- Update `website/docs/r/mssql_database.html.markdown`: remove the
  `immutable_backups_enabled` argument doc line, update `immutability_mode` description
  to reflect it also controls the on/off state.

---

## Comment 3 — Use `pointer.From` in `FlattenLongTermRetentionPolicy`
**File:** `internal/services/mssql/helper/sql_retention_policies.go` lines 178–186

### What to do
Replace the manual nil-check + dereference blocks with `pointer.From`:

```go
// BEFORE
immutableBackupsEnabled := false
if input.Properties.TimeBasedImmutability != nil {
    immutableBackupsEnabled = *input.Properties.TimeBasedImmutability == longtermretentionpolicies.TimeBasedImmutabilityEnabled
}

immutabilityMode := ""
if input.Properties.TimeBasedImmutabilityMode != nil {
    immutabilityMode = string(*input.Properties.TimeBasedImmutabilityMode)
}

// AFTER (once immutable_backups_enabled is removed per comment 2,
//         only immutabilityMode remains)
immutabilityMode := string(pointer.From(input.Properties.TimeBasedImmutabilityMode))
```

Note: also apply `pointer.From` to the existing fields that use the same pattern
(`yearlyRetention` at line 175 still uses `*input.Properties.YearlyRetention` directly
instead of `pointer.From`).

---

## Comment 4 — `CustomizeDiff` to block `Locked → Unlocked`
**File:** `internal/services/mssql/mssql_database_resource.go` (CustomizeDiff block, ~line 77)
**Docs:** `website/docs/r/mssql_database.html.markdown` lines 290–291

### What to do
Add a `pluginsdk.ForceNewIfChange` (or a custom error func) inside `CustomizeDiff`
for `long_term_retention_policy.0.immutability_mode` that:
- If old == "Locked" and new != "Locked" (i.e. trying to go back to Unlocked or unset)
  → return an error at plan time: "once `immutability_mode` is set to `Locked` it
    cannot be changed or removed as the immutability is permanent"
- If old == "Unlocked" and new == "Locked"
  → allow the change (this direction is valid)

Pattern to follow: `pluginsdk.ForceNewIfChange("enclave_type", ...)` at line 80.

Example:
```go
pluginsdk.ForceNewIfChange("long_term_retention_policy.0.immutability_mode",
    func(ctx context.Context, old, new, _ interface{}) bool {
        return old.(string) == string(longtermretentionpolicies.TimeBasedImmutabilityModeLocked) &&
               new.(string) != string(longtermretentionpolicies.TimeBasedImmutabilityModeLocked)
    }),
```

Or use a custom error func if ForceNew is not appropriate (since recreating the
database may not be desirable — confirm with reviewer).

Update the docs to reflect the ForceNew/error behaviour.

---

## After all comments resolved
1. Run: `go build ./internal/services/mssql/...`
2. Run: `go fmt ./internal/services/mssql/... && go vet ./internal/services/mssql/...`
3. Run lint: `~/go/bin/golangci-lint run ./internal/services/mssql/...`
4. Run relevant tests:
   ```
   TF_ACC=1 go test ./internal/services/mssql/ -v \
     -run "TestAccMsSqlDatabase_withLongTermRetentionPolicy|TestAccMsSqlDatabase_threatDetectionPolicy" \
     -timeout 30m
   ```
5. Push and update PR.
