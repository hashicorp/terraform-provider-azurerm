# `internal/preflight`

This package wraps the Azure Preflight Validation API
(`POST /providers/Microsoft.Resources/validateResources`) to allow the AzureRM provider to
validate resource configurations at `terraform plan` time — before any changes are applied.

Preflight is gated behind the `features.enhanced_validation.preflight_enabled` feature flag and
is implemented via `CustomizeDiff` on resources that support it.

---

## Constraint: PUT payloads only

The Azure Preflight Validation API validates full ARM **PUT** payloads. **PATCH operations are
not supported.** The `properties` value passed to `NewValidationRequest` must represent the
complete resource body exactly as it will be sent to the ARM API. A partial or PATCH-style
payload will produce unreliable validation results — either false positives (blocked valid
configs) or false negatives (invalid configs that pass undetected).

---

## When to implement preflight validation

Add preflight to a resource when:

1. The resource has a stable `expandCreateForMyResource` (or equivalent) function that returns the full ARM body.
2. The resource uses `CustomizeDiff` (implemented via `sdk.ResourceWithCustomizeDiff`).
3. The resource type is supported by the Azure Preflight Validation API.

Preflight must always be wrapped in a feature check inside `CustomizeDiff`:

```go
if metadata.Client.Features.EnhancedValidation.PreflightEnabled {
    // validation logic
}
```

---

## Constructing the request

### `NewValidationRequest` — default

For most resources, the resource type is derived automatically from the ARM ID's path segments.
`NewValidationRequest` handles this:

```go
resId := mypackage.NewMyResourceID(metadata.Client.Account.SubscriptionId, model.ResourceGroupName, model.Name)
preflightValidate, err := preflight.NewValidationRequest(pointer.To(model.Location), pointer.To(resId), "2025-01-01", req)
```

The `resourceType` field in the request (used to select validation rules) is extracted from the
ID's static path segment — e.g. `"virtualNetworks"` from a `VirtualNetworkId`, or `"sites"` from
an `AppServiceId`.

### `NewValidationRequestWithTypeOverride` — when the type differs from the ID segment

Some Azure resource providers bundle multiple product offerings under a single namespace. In these
cases the preflight API uses a different type discriminator than the ARM ID segment. Use
`NewValidationRequestWithTypeOverride` to specify the correct value explicitly:

```go
resId := redisenterprise.NewRedisEnterpriseID(metadata.Client.Account.SubscriptionId, model.ResourceGroupName, model.Name)
// Microsoft.Cache hosts both "redis" (classic) and "redisEnterprise" products under the same
// provider. The preflight API uses "redis" as the type discriminator for this resource.
preflightValidate, err := preflight.NewValidationRequestWithTypeOverride(pointer.To(model.Location), pointer.To(resId), "redis", "2025-07-01", req)
```

---

## Naming convention

Expand functions used with preflight must be named for their resource type to avoid redeclaration
conflicts in packages that contain multiple resources:

- `expandCreateForMyResource`
- `expandUpdateForMyResource`

The `ForPreflight` suffix is not used — the resource name provides the necessary disambiguation.

---

## Data completeness in `CustomizeDiff`

`ResourceDiff` always contains the **complete planned state**, not just changed values. When
`metadata.DecodeDiff(&model)` is called, the SDK resolves each field through a priority chain of
`state → config → diff → newDiff`. Unchanged fields on an update are satisfied from the prior
`state` layer and are never empty.

This means `expandCreateForMyResource(model)` has all the data it needs during an update, in the
same way it does during a create. **Pattern 3 exists for structural differences in the request
body, not because data is missing.**

---

## Implementation patterns

Three patterns are available. Choose based on how the resource handles create vs. update.

### Pattern 1 — Reuse the create expand function for all operations

**Use when:** The resource uses idempotent PUT semantics for both create and update, and the
same full body is sent in both cases.

```go
func (r MyResource) CustomizeDiff() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 5 * time.Minute,
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            if metadata.ResourceDiff == nil {
                return nil
            }

            var model MyResourceModel
            if err := metadata.DecodeDiff(&model); err != nil {
                return err
            }

            if metadata.Client.Features.EnhancedValidation.PreflightEnabled {
                // Only validate when there are actual changes or the resource is new.
                if len(metadata.ResourceDiff.GetChangedKeysPrefix("")) > 0 || metadata.ResourceDiff.Id() == "" {
                    req, err := expandCreateForMyResource(model)
                    if err != nil {
                        return err
                    }

                    resId := mypackage.NewMyResourceID(metadata.Client.Account.SubscriptionId, model.ResourceGroupName, model.Name)
                    preflightValidate, err := preflight.NewValidationRequest(pointer.To(model.Location), pointer.To(resId), "2025-01-01", req)
                    if err != nil {
                        return fmt.Errorf("constructing preflight validation request: %w", err)
                    }

                    if err = preflightValidate.ValidateResource(ctx, metadata); err != nil {
                        return err
                    }
                }
            }

            // ... remaining CustomizeDiff logic
            return nil
        },
    }
}
```

---

### Pattern 2 — Validate on create and ForceNew replacement only

**Use when:** You want to skip preflight for in-place updates, e.g. to avoid false positives
from immutable fields that appear in the full PUT body, or when the update path uses PATCH.

`ResourceDiff` does not expose a `RequiresNew()` method — ForceNew replacement is detected
by checking whether any changed key has `ForceNew: true` in the resource schema:

```go
if metadata.Client.Features.EnhancedValidation.PreflightEnabled {
    isNewResource := metadata.ResourceDiff.Id() == ""
    isForceNewReplacement := false

    if !isNewResource {
        for _, key := range metadata.ResourceDiff.GetChangedKeysPrefix("") {
            if s, ok := r.Arguments()[key]; ok && s.ForceNew {
                isForceNewReplacement = true
                break
            }
        }
    }

    if isNewResource || isForceNewReplacement {
        req, err := expandCreateForMyResource(model)
        if err != nil {
            return err
        }

        resId := mypackage.NewMyResourceID(metadata.Client.Account.SubscriptionId, model.ResourceGroupName, model.Name)
        preflightValidate, err := preflight.NewValidationRequest(pointer.To(model.Location), pointer.To(resId), "2025-01-01", req)
        if err != nil {
            return fmt.Errorf("constructing preflight validation request: %w", err)
        }

        if err = preflightValidate.ValidateResource(ctx, metadata); err != nil {
            return err
        }
    }
}
```

> ForceNew replacement results in a destroy + create, so `expandCreateForMyResource` is always
> the correct payload — the old resource is being destroyed and a new one is being created in
> its place.

---

### Pattern 3 — Different expand functions for create vs. update

**Use when:** The create and update PUT bodies are **structurally different** — for example,
a field that is required at create time is rejected if re-sent on update (common for immutable
`kind` or `sku.tier` fields on some ARM resource types).

This is **not** about data availability — `ResourceDiff` always contains the full planned state.
See [Data completeness in CustomizeDiff](#data-completeness-in-customizediff) above.

```go
if metadata.Client.Features.EnhancedValidation.PreflightEnabled {
    if len(metadata.ResourceDiff.GetChangedKeysPrefix("")) > 0 || metadata.ResourceDiff.Id() == "" {
        var req any
        var err error

        if metadata.ResourceDiff.Id() == "" {
            req, err = expandCreateForMyResource(model)
        } else {
            req, err = expandUpdateForMyResource(model)
        }
        if err != nil {
            return err
        }

        resId := mypackage.NewMyResourceID(metadata.Client.Account.SubscriptionId, model.ResourceGroupName, model.Name)
        preflightValidate, err := preflight.NewValidationRequest(pointer.To(model.Location), pointer.To(resId), "2025-01-01", req)
        if err != nil {
            return fmt.Errorf("constructing preflight validation request: %w", err)
        }

        if err = preflightValidate.ValidateResource(ctx, metadata); err != nil {
            return err
        }
    }
}
```

> If the resource's existing `expandUpdate` returns a PATCH body, do **not** pass it to
> preflight. Instead, create a dedicated `expandUpdateForMyResource` that constructs the
> complete PUT body from the model. The preflight API does not support partial payloads.

---

## Decision guide

| Question | Answer | Pattern |
|---|---|---|
| Does `update()` send a full PUT body with the same shape as `create()`? | Yes | **1** |
| Do you want to skip plan-time API calls for in-place updates? | Yes | **2** |
| Does `update()` send a structurally different (but still full) PUT body than `create()`? | Yes | **3** |
| Does `update()` send a PATCH (partial) body? | Yes | **3** with a dedicated `expandUpdateForMyResource` |

---

## Plan-time behaviour

`CustomizeDiff` runs during Terraform's `PlanResourceChange` phase. This means:

- Multiple resources' `CustomizeDiff` functions run **concurrently** for independent resources.
- The preflight API validates the configuration payload against ARM schema and Azure Policy.
  It does **not** check whether the resource or its dependencies currently exist in Azure.
- Making HTTP calls from `CustomizeDiff` adds latency to `terraform plan`. The change-detection
  guard (`len(GetChangedKeysPrefix("")) > 0 || Id() == ""`) avoids redundant API calls when
  nothing has changed between plan invocations.

---

## Common pitfalls

- **Forgetting the nil guard:** `CustomizeDiff` is called in contexts where `ResourceDiff`
  may be nil (e.g. during import). Always guard with `if metadata.ResourceDiff == nil { return nil }`.
- **Partial payloads:** Passing a PATCH-style body will silently under-validate. Ensure the
  expand function returns every field that will appear in the actual ARM PUT call.
- **Missing the change guard:** Without `len(GetChangedKeysPrefix("")) > 0 || Id() == ""`,
  preflight will make an API call on every plan even when nothing has changed.

---

## Known limitation: cross-resource computed values

`ResourceDiff` contains the complete planned state for all **known** values. When a config
attribute references a computed output from another resource that does not yet exist, Terraform
marks that value as `(known after apply)`. At the `ResourceDiff` level this surfaces as the
field's zero value (`""`, `0`, `[]`, etc.) rather than the real value.

This means the preflight payload will contain an empty value for that field, while the actual
`Create` payload — called after the dependency has been applied — will contain the real value.

**Example:**

```hcl
resource "azurerm_user_assigned_identity" "identity" { ... }

resource "azurerm_managed_redis" "redis" {
  identity {
    identity_ids = [azurerm_user_assigned_identity.identity.id]  # (known after apply) at plan time
  }
}
```

During `CustomizeDiff`, `identity.identity_ids` will be `[]` in the preflight payload. During
`Create` (after Apply has resolved the dependency), it will contain the real UAI resource ID.

### Impact by field type

| Scenario | Effect on preflight | Risk |
|---|---|---|
| `Optional` field with cross-resource ref | Field absent from payload | **False negative** — validation gap, not an error |
| `Optional+Computed` field with cross-resource ref | Zero value sent | **False negative** — ARM typically accepts or ignores |
| `Required` field with cross-resource ref (e.g. `location = rg.location`) | Empty value sent | **False positive** — ARM may reject the empty required field |

### Practical guidance

For most resources the risk is limited to **false negatives** (reduced coverage), not
**false positives** (spurious plan failures). The `Required` fields in most AzureRM resources
(`name`, `location`, `resource_group_name`) are almost always literal values in user config,
not computed cross-resource references.

If a resource commonly has a `Required` field populated via a cross-resource computed reference,
consider guarding the preflight call to avoid false positives:

```go
// Skip preflight if a required field is not yet known at plan time.
if model.Location == "" {
    return nil
}
```

> **Note:** `terraform-plugin-framework` exposes the full plan including unknown value markers.
> When Framework support is added, it will be possible to detect interpolated values explicitly
> and skip or adjust the preflight payload accordingly.