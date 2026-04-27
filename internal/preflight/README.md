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

Preflight is always wrapped in:

```go
if metadata.Client.Features.EnhancedValidation.PreflightEnabled {
    // validation logic
}
```

---

## Implementation Patterns

Three patterns are available. Choose based on how the resource handles create vs. update.

> **Naming convention:** Expand functions used with preflight must be named for their resource
> type to avoid redeclaration conflicts in packages that contain multiple resources:
> `expandCreateForMyResource`, `expandUpdateForMyResource`. The `ForPreflight` suffix is
> not used — the resource name provides the necessary disambiguation.

> **Data completeness:** `ResourceDiff` always contains the complete planned state, not just
> changed values. When `metadata.DecodeDiff(&model)` is called, the SDK resolves each field
> through a priority chain of `state → config → diff → newDiff`. Unchanged fields on an
> update are satisfied from the prior `state` layer — they are never empty. This means
> `expandCreateForMyResource(model)` has all the data it needs during an update in the same
> way it does during a create. **Pattern 3 is about structural differences in the request
> body, not data availability.**

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
                // expandCreate is reused for updates: this resource uses full PUT semantics
                // for both create and update operations.
                if len(metadata.ResourceDiff.GetChangedKeysPrefix("")) > 0 || metadata.ResourceDiff.Id() == "" {
                    req, err := expandCreateForMyResource(model)
                    if err != nil {
                        return err
                    }

                    resId := mypackage.NewMyResourceID(metadata.Client.Account.SubscriptionId, model.ResourceGroupName, model.Name)
                    preflightValidate, err := preflight.NewValidationRequest(pointer.To(model.Location), pointer.To(resId), "myResourceType", "2025-01-01", req)
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
from a resource with immutable fields that are present in the full PUT body.

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
        preflightValidate, err := preflight.NewValidationRequest(pointer.To(model.Location), pointer.To(resId), "myResourceType", "2025-01-01", req)
        if err != nil {
            return fmt.Errorf("constructing preflight validation request: %w", err)
        }

        if err = preflightValidate.ValidateResource(ctx, metadata); err != nil {
            return err
        }
    }
}
```

> ForceNew replacement results in a destroy + create, so `expandCreateForMyResource` is always the
> correct payload — the old resource is being destroyed and a new one created.

---

### Pattern 3 — Different expand functions for create vs. update

**Use when:** The create and update PUT bodies are **structurally different** — for example,
a field that is required at create time is rejected if re-sent on update (common for immutable
`kind` or `sku.tier` fields on some ARM resource types), or the update body must include fields
that reference computed values that only exist after the resource has been created. This is
**not** about data availability — `ResourceDiff` always contains the full planned state.

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
        preflightValidate, err := preflight.NewValidationRequest(pointer.To(model.Location), pointer.To(resId), "myResourceType", "2025-01-01", req)
        if err != nil {
            return fmt.Errorf("constructing preflight validation request: %w", err)
        }

        if err = preflightValidate.ValidateResource(ctx, metadata); err != nil {
            return err
        }
    }
}
```

> **Important:** If the resource's existing `expandUpdate` returns a PATCH body, do **not**
> pass it to preflight. Instead, create a dedicated `expandUpdateForMyResource` that constructs
> the complete PUT body from the model. The preflight API does not support partial payloads.

---

## Decision guide

| Question | Answer | Pattern |
|---|---|---|
| Does `update()` send a full PUT body with the same shape as `create()`? | Yes | **1** |
| Do you want to skip plan-time API calls for in-place updates? | Yes | **2** |
| Does `update()` send a different (but still full) PUT body than `create()`? | Yes | **3** |
| Does `update()` send a PATCH (partial) body? | Yes | **3** with a dedicated `expandUpdateForMyResource` |

---

## DAG / lifecycle context

`CustomizeDiff` runs during Terraform's `PlanResourceChange` phase. This means:

- Multiple resources' `CustomizeDiff` functions run **concurrently** for independent resources.
- The preflight API validates the configuration payload against ARM schema and Azure Policy.
  It does **not** check whether the resource or its dependencies currently exist in Azure.
- Making HTTP calls from `CustomizeDiff` is supported but adds latency to `terraform plan`.
  The change-detection guard (`len(GetChangedKeysPrefix("")) > 0 || Id() == ""`) avoids
  redundant API calls when nothing has changed between plan invocations.

---

## Common pitfalls

- **Forgetting the nil guard:** `CustomizeDiff` is called in contexts where `ResourceDiff`
  may be nil (e.g. during import). Always guard with `if metadata.ResourceDiff == nil { return nil }`.
- **Partial payloads:** Passing a PATCH-style body will silently under-validate. Ensure the
  expand function returns every field that will be in the actual ARM PUT call.
- **Missing the change guard:** Without `len(GetChangedKeysPrefix("")) > 0 || Id() == ""`,
  preflight will make an API call on every plan even when nothing has changed.
