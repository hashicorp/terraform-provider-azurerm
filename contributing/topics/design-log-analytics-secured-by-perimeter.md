# Design: Support `SecuredByPerimeter` for `azurerm_log_analytics_workspace`

**Status:** RFC — awaiting team review  
**Date:** 2026-05-19 (updated 2026-06-02 with reviewer feedback)  
**Issue:** Customer request — unable to manage Log Analytics Workspaces that use Azure Network Security Perimeter (NSP) via Terraform  
**API Ref:** [`PublicNetworkAccessType` enum (2025-07-01)](https://learn.microsoft.com/en-us/rest/api/loganalytics/workspaces/create-or-update?view=rest-loganalytics-2025-07-01&tabs=HTTP#publicnetworkaccesstype)

---

## Problem

`internet_ingestion_enabled` and `internet_query_enabled` are currently `TypeBool`. The Azure API's `PublicNetworkAccessType` enum has three values — `Enabled`, `Disabled`, and `SecuredByPerimeter` — but a bool can only express two of them.

When a workspace is configured with `SecuredByPerimeter` (e.g. via the portal or Azure Policy), the provider's read path evaluates it as `false` and plans an unwanted write to `Disabled` on the next apply. Users are forced to use `ignore_changes` as a workaround.

**Key research finding:** `SecuredByPerimeter` was introduced in the `2025-07-01` stable API release. It is **not present** in the currently vendored `2023-09-01` SDK package. Patching vendor files directly is prohibited by project convention (`go mod vendor` must be the sole source of truth). An API version upgrade to `2025-07-01` is therefore a prerequisite — not an option.

---

## Proposed Change

Per the provider's [breaking changes guide](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/contributing/topics/guide-breaking-changes.md#breaking-schema-changes-and-deprecations), changing the type of an existing property is a breaking change and must follow the deprecation pattern — **not** an in-place type change. The correct approach is to introduce a new string property alongside the existing bool, let both coexist as Optional+Computed (O+C) in v4.x, then remove the deprecated bool in v5.0.

No state migration is required because the bool type of `internet_ingestion_enabled` never changes in v4.x. Both attributes are written to state simultaneously on each read.

### v4.x — Add new string attributes, deprecate the bools

**Schema** (gated under `if !features.FivePointOh()`):

```go
// New string attributes — the permanent replacement
"internet_ingestion_access_type": {
    Type:     pluginsdk.TypeString,
    Optional: true,
    Computed: true,
    ValidateFunc: validation.StringInSlice([]string{
        "Enabled", "Disabled", "SecuredByPerimeter",
    }, false),
    ConflictsWith: []string{"internet_ingestion_enabled"},
},
"internet_query_access_type": {
    Type:     pluginsdk.TypeString,
    Optional: true,
    Computed: true,
    ValidateFunc: validation.StringInSlice([]string{
        "Enabled", "Disabled", "SecuredByPerimeter",
    }, false),
    ConflictsWith: []string{"internet_query_enabled"},
},

// Existing bool attributes — deprecated, kept as O+C until v5.0
"internet_ingestion_enabled": {
    Type:          pluginsdk.TypeBool,
    Optional:      true,
    Computed:      true,
    Deprecated:    "`internet_ingestion_enabled` has been deprecated in favour of `internet_ingestion_access_type` and will be removed in v5.0 of the AzureRM Provider",
    ConflictsWith: []string{"internet_ingestion_access_type"},
},
"internet_query_enabled": {
    Type:          pluginsdk.TypeBool,
    Optional:      true,
    Computed:      true,
    Deprecated:    "`internet_query_enabled` has been deprecated in favour of `internet_query_access_type` and will be removed in v5.0 of the AzureRM Provider",
    ConflictsWith: []string{"internet_query_access_type"},
},
```

Outside the feature flag block, the v5.0 schema defines only the string attributes with no `ConflictsWith` and no `Computed`:

```go
// In v5.0 (default schema, outside the !FivePointOh block)
"internet_ingestion_access_type": {
    Type:     pluginsdk.TypeString,
    Optional: true,
    Default:  "Enabled",
    ValidateFunc: validation.StringInSlice([]string{
        "Enabled", "Disabled", "SecuredByPerimeter",
    }, false),
},
```

**Create/Update logic** (prefer the new string attribute; fall back to the deprecated bool):

```go
if !features.FivePointOh() {
    if v, ok := d.GetOk("internet_ingestion_access_type"); ok {
        internetIngestionEnabled = workspaces.PublicNetworkAccessType(v.(string))
    } else if d.Get("internet_ingestion_enabled").(bool) {
        internetIngestionEnabled = workspaces.PublicNetworkAccessTypeEnabled
    } else {
        internetIngestionEnabled = workspaces.PublicNetworkAccessTypeDisabled
    }
} else {
    internetIngestionEnabled = workspaces.PublicNetworkAccessType(d.Get("internet_ingestion_access_type").(string))
}
```

**Read logic** — set both attributes on every read:

```go
if v := props.PublicNetworkAccessForIngestion; v != nil {
    d.Set("internet_ingestion_access_type", string(*v))

    if !features.FivePointOh() {
        d.Set("internet_ingestion_enabled", *v == workspaces.PublicNetworkAccessTypeEnabled)
    }
}
```

> When the API returns `SecuredByPerimeter`, `internet_ingestion_enabled` (bool) is set to `false` in state. This is an accepted approximation during the deprecation period — users on the deprecated bool will see `false` but no unwanted plan diff, because `Computed: true` suppresses it. They should migrate to `internet_ingestion_access_type`.

### v5.0 — Remove the deprecated bools

Delete the `if !features.FivePointOh() { ... }` block entirely. The `internet_ingestion_enabled` and `internet_query_enabled` attributes are removed from the schema, CRUD logic, and documentation. Add entries to `website/docs/5.0-upgrade-guide.markdown`.

### SDK prerequisite

`SecuredByPerimeter` is not present in the currently vendored `2023-09-01` SDK. The constant must be available for the new `ValidateFunc` and CRUD logic to compile. Two options:

- **Preferred:** Bump the vendored `go-azure-sdk` to a release that includes `operationalinsights/2025-07-01` (which natively defines the constant) and update the client import. Run `go mod tidy && go mod vendor` — do not edit vendor files by hand.
- **Acceptable short-term:** If a qualifying SDK release is not yet available, the constant can be defined locally in the provider package as `const publicNetworkAccessTypeSecuredByPerimeter = "SecuredByPerimeter"` and the SDK import left at `2023-09-01`. The `parsePublicNetworkAccessType` unmarshal function in the SDK already falls back to a best-effort pass-through for unknown values, so the Read path will work correctly even without the constant in the SDK.

---

## PoC Test Results (2026-05-26)

A proof-of-concept was built locally implementing the **in-place type change** (bool → string directly) to understand what breaking changes that approach would produce. The PoC was not the correct final approach (see reviewer comment below), but its findings directly informed the design.

`terraform plan` was run against three configs — `"Enabled"`, `"Disabled"`, and `"SecuredByPerimeter"` — and passed cleanly, confirming the schema logic and SDK constant changes work correctly.

**Reviewer comment (2026-06-02):** An in-place type change is not necessary and introduces avoidable breakage. The correct pattern per [`guide-breaking-changes.md`](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/contributing/topics/guide-breaking-changes.md#breaking-schema-changes-and-deprecations) is to introduce a new property and deprecate the old one, with both marked O+C in v4.x. This avoids a state migration entirely because the bool type of the existing attribute never changes.

---

## Breaking Changes

| # | Impact | Detail |
|---|---|---|
| 1 | **v4.x: no state migration** | The bool attributes stay as `TypeBool` throughout v4.x. No `SchemaVersion` bump is needed. Existing state is unaffected. |
| 2 | **v5.0: HCL update required** | Users must replace `internet_ingestion_enabled = true/false` with `internet_ingestion_access_type = "Enabled"/"Disabled"` when upgrading to v5.0. The deprecation warning in v4.x gives them advance notice. |
| 3 | **`SecuredByPerimeter` on read in v4.x** | When the API returns `SecuredByPerimeter`, the deprecated bool is set to `false` in state (best-effort approximation). No plan diff is produced because both attributes are `Computed`. Users managing `SecuredByPerimeter` workspaces should adopt `internet_ingestion_access_type` immediately. |

---

## Files to Change

| File | Change |
|---|---|
| `go.mod` / vendor | SDK bump to obtain `SecuredByPerimeter` constant (or define locally — see SDK prerequisite above) |
| `internal/services/loganalytics/log_analytics_workspace_resource.go` | Add new string attributes; deprecate old bools; update CRUD logic; gate all of it behind `!features.FivePointOh()` |
| `internal/services/loganalytics/log_analytics_workspace_data_source.go` | Add `internet_ingestion_access_type` and `internet_query_access_type` as `Computed` string attributes |
| `internal/services/loganalytics/log_analytics_workspace_resource_test.go` | Add acceptance tests for `SecuredByPerimeter`; one test config keeps the deprecated bool to guard against regression |
| `website/docs/r/log_analytics_workspace.html.markdown` | Document new string attributes; add deprecation notice to old bool docs |
| `website/docs/5.0-upgrade-guide.markdown` | Add removal notice for `internet_ingestion_enabled` and `internet_query_enabled` |

---

## Open Questions

1. **SDK:** Does the current pinned `go-azure-sdk` version already bundle `operationalinsights/2025-07-01`? If not, should we bump the SDK now or use a local constant as a temporary measure?
2. **NSP acceptance test:** A `SecuredByPerimeter` acceptance test requires a subscription with NSP configured. Does the CI environment support this, or does this test need a separate environment label?
