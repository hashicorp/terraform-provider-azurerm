---
applyTo: "internal/**/*.go"
description: Migration patterns and upgrade procedures for the Terraform AzureRM provider including implementation approach transitions, breaking changes, and version compatibility.
---

# Migration Guide

**Quick navigation:** [🔄 Implementation Migration](#🔄-implementation-approach-migration) | [💔 Breaking Changes](#💔-breaking-change-patterns) | [📦 Version Compatibility](#📦-version-compatibility) | [🚧 Upgrade Procedures](#🚧-upgrade-procedures)

## 🔄 Implementation Approach Migration

### Migration Decision Matrix

| Scenario | Action | Approach |
|----------|--------|----------|
| New Resource | Always use Typed Resource Implementation | Start with typed from day one |
| Bug Fix (< 5 lines) | Maintain Untyped Implementation | Quick fix in existing pattern |
| Feature Addition (< 50 lines) | Consider migration if touching >30% of resource | Evaluate cost/benefit |
| Major Refactor (> 50 lines) | Migrate to Typed Implementation | Plan migration with comprehensive testing |
| EOL/Deprecation Planning | Plan Typed Migration | Include in deprecation timeline |

### Service Registration During Migration

**Dual Registration Pattern:**
Services often need to be registered in both lists temporarily during migration:

```go
// In internal/provider/services.go

func SupportedTypedServices() []sdk.TypedServiceRegistration {
    services := []sdk.TypedServiceRegistration{
        // Add service here when it has any typed resources
        cdn.Registration{},
        // ...other services
    }
    return services
}

func SupportedUntypedServices() []sdk.UntypedServiceRegistration {
    return func() []sdk.UntypedServiceRegistration {
        out := []sdk.UntypedServiceRegistration{
            // Keep service here until all resources are migrated
            cdn.Registration{},
            // ...other services
        }
        return out
    }()
}
```

**Registration Implementation Requirements:**
```go
var (
	_ sdk.TypedServiceRegistration                   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/connections"
}

// REQUIRED: Always implement both typed and untyped functions
func (r Registration) DataSources() []sdk.DataSource {
    return []sdk.DataSource{
        // Typed data sources here, or empty slice if none exist
        ApiConnectionDataSource{},
    }
}

func (r Registration) Resources() []sdk.Resource {
    return []sdk.Resource{
        // Typed resources here, or empty slice if none exist
        // Add typed resources here when implemented
    }
}

func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
    return map[string]*pluginsdk.Resource{
        // Untyped data sources here, or empty map if none exist
        "azurerm_managed_api": dataSourceManagedApi(),
    }
}

func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
    return map[string]*pluginsdk.Resource{
        // Untyped resources here, or empty map if none exist
        "azurerm_api_connection": resourceApiConnection(),
    }
}
```

### Step-by-Step Migration Process

**Phase 1: Assessment and Planning**
```go
// 1. Create backup branch
git checkout -b migration/resource-name-to-typed

// 2. Analyze existing untyped implementation
// Study these files:
// - internal/services/servicename/resource_name_resource.go
// - internal/services/servicename/resource_name_resource_test.go
// - website/docs/r/service_name_resource.html.markdown

// 3. Document current behavior
// - All schema fields and their types
// - CRUD operation behaviors
// - Error handling patterns
// - CustomizeDiff validations
// - Import functionality
```

**Phase 2: Model Structure Creation**
```go
// Create the typed model structure
type ServiceNameResourceModel struct {
    // Required fields first (alphabetical within category)
    Name              string            `tfschema:"name"`
    ResourceGroup     string            `tfschema:"resource_group_name"`
    Location          string            `tfschema:"location"`

    // Optional fields (alphabetical)
    Enabled           bool              `tfschema:"enabled"`
    Settings          map[string]string `tfschema:"settings"`
    Tags              map[string]string `tfschema:"tags"`

    // Complex nested structures
    Configuration     []ConfigModel     `tfschema:"configuration"`

    // Computed attributes last
    Id                string            `tfschema:"id"`
    Endpoint          string            `tfschema:"endpoint"`
}

// Nested model structures
type ConfigModel struct {
    Setting1 string `tfschema:"setting1"`
    Setting2 string `tfschema:"setting2"`
}
```

**Phase 3: Resource Structure Implementation**
```go
// Implement the resource structure
type ServiceNameResource struct{}

// Required interface implementations
var (
    _ sdk.Resource           = ServiceNameResource{}
    _ sdk.ResourceWithUpdate = ServiceNameResource{} // Only if resource supports updates
)

func (r ServiceNameResource) ResourceType() string {
    return "azurerm_service_name_resource"
}

func (r ServiceNameResource) ModelObject() interface{} {
    return &ServiceNameResourceModel{}
}

func (r ServiceNameResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
    // Reuse existing ID validation function
    return parse.ValidateServiceNameResourceID
}
```

**Phase 4: CRUD Migration**
```go
// Migrate Create operation
func (r ServiceNameResource) Create() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 30 * time.Minute, // Copy timeout from untyped resource
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            // 1. Client access (new pattern)
            client := metadata.Client.ServiceName.ResourceClient
            subscriptionId := metadata.Client.Account.SubscriptionId

            // 2. Model decoding (new pattern)
            var model ServiceNameResourceModel
            if err := metadata.Decode(&model); err != nil {
                return fmt.Errorf("decoding: %+v", err)
            }

            // 3. Resource ID creation (same as untyped)
            id := parse.NewServiceNameResourceID(subscriptionId, model.ResourceGroup, model.Name)

            // 4. Import check (new pattern)
            metadata.Logger.Infof("Import check for %s", id)
            existing, err := client.Get(ctx, id)
            if err != nil && !response.WasNotFound(existing.HttpResponse) {
                return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
            }

            if !response.WasNotFound(existing.HttpResponse) {
                return metadata.ResourceRequiresImport(r.ResourceType(), id)
            }

            // 5. API call (similar to untyped but with model data)
            metadata.Logger.Infof("Creating %s", id)

            properties := servicenametype.Resource{
                Location: model.Location,
                Properties: &servicenametype.ResourceProperties{
                    Enabled: pointer.To(model.Enabled),
                    // Map other model fields to API structure
                },
                Tags: &model.Tags,
            }

            if err := client.CreateOrUpdateThenPoll(ctx, id, properties); err != nil {
                return fmt.Errorf("creating %s: %+v", id, err)
            }

            // 6. ID setting (new pattern)
            metadata.SetID(id)
            return nil
        },
    }
}
```

**Phase 5: Testing Migration**
```go
// Update test patterns
func TestAccServiceNameResource_basic(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_service_name_resource", "test")
    r := ServiceNameResource{} // New pattern: use struct instead of function

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config: r.basic(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctest-%d", data.RandomInteger)),
                check.That(data.ResourceName).Key("resource_group_name").HasValue(fmt.Sprintf("acctestRG-%d", data.RandomInteger)),
            ),
        },
        data.ImportStep(), // Keep existing import tests
    })
}

// Update Exists function for typed resource
func (r ServiceNameResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
    id, err := parse.ServiceNameResourceID(state.ID)
    if err != nil {
        return nil, err
    }

    resp, err := clients.ServiceName.ResourceClient.Get(ctx, *id)
    if err != nil {
        return nil, fmt.Errorf("reading %s: %+v", *id, err)
    }

    return utils.Bool(resp.Model != nil), nil
}
```

**Phase 6: Service Registration Update**
```go
// Update service registration to include typed resource
func (r Registration) Resources() []sdk.Resource {
    return []sdk.Resource{
        ServiceNameResource{}, // Add migrated typed resource
        // Keep other typed resources
    }
}

func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
    return map[string]*pluginsdk.Resource{
        // Remove migrated resource from here
        // "azurerm_service_name_resource": resourceServiceNameResource(), // REMOVE

        // Keep other untyped resources
        "azurerm_other_resource": resourceOtherResource(),
    }
}
```

### Migration Validation Checklist

**Functionality Verification:**
- [ ] All CRUD operations work correctly
- [ ] Resource import functionality preserved
- [ ] CustomizeDiff validations migrated correctly
- [ ] Error handling maintains same user experience
- [ ] Timeout configurations preserved
- [ ] All schema fields and validation rules maintained

**Testing Verification:**
- [ ] All existing acceptance tests pass
- [ ] Import tests continue to work
- [ ] Error scenario tests function correctly
- [ ] Test configurations require no changes
- [ ] Performance characteristics remain similar

**State Compatibility Verification:**
- [ ] Existing Terraform state works without manual intervention
- [ ] Resource attributes remain accessible
- [ ] Computed values continue to be populated correctly
- [ ] No unexpected ForceNew behaviors introduced

**Documentation Verification:**
- [ ] Resource documentation requires no changes
- [ ] Import syntax remains the same
- [ ] Examples continue to work
- [ ] Attribute descriptions remain accurate

---
[⬆️ Back to top](#migration-guide)

## 💔 Breaking Change Patterns

### Field Rename Migration

**When Implementing Field Renames:**
```go
// BEFORE - Generic field name
"scrubbing_rule": {
    Type:     pluginsdk.TypeSet,
    Optional: true,
    Elem: &pluginsdk.Resource{
        Schema: map[string]*pluginsdk.Schema{
            "match_variable": {
                Type:     pluginsdk.TypeString,
                Required: true,
            },
        },
    },
},

// AFTER - Descriptive field name
"log_scrubbing_rule": {
    Type:     pluginsdk.TypeSet,
    Optional: true,
    Elem: &pluginsdk.Resource{
        Schema: map[string]*pluginsdk.Schema{
            "match_variable": {
                Type:     pluginsdk.TypeString,
                Required: true,
            },
        },
    },
},
```

**Field Rename Testing Requirements:**
- [ ] Resource implementation updated with new field name
- [ ] Data source implementation updated for consistency
- [ ] All test configurations updated
- [ ] Documentation updated with new field name
- [ ] Import functionality verified
- [ ] State compatibility ensured

### Schema Flattening Breaking Changes

**Before Flattening (v3.x):**
```hcl
resource "azurerm_cdn_frontdoor_profile" "example" {
  name = "example"

  log_scrubbing {
    enabled = true

    scrubbing_rule {
      match_variable = "QueryStringArgNames"
    }
  }
}
```

**After Flattening (v4.x):**
```hcl
resource "azurerm_cdn_frontdoor_profile" "example" {
  name = "example"

  scrubbing_rule {
    match_variable = "QueryStringArgNames"
  }
}
```

### Version-Specific Breaking Changes

**v3.x to v4.x Migration Patterns:**
- **Field Renames**: Generic names → Descriptive names
- **Schema Flattening**: Remove unnecessary wrapper structures
- **"None" Pattern Adoption**: Remove explicit "None" values from schema
- **Optional+Computed Evolution**: Simplify Azure-managed defaults

**v4.x to v5.x Planned Changes:**
- **Typed Resource Migration**: Complete migration from untyped to typed resources
- **Legacy Pattern Removal**: Remove deprecated patterns and anti-patterns
- **SDK Updates**: Migration to newer Azure SDK patterns

---
[⬆️ Back to top](#migration-guide)

## 📦 Version Compatibility

### Terraform Plugin SDK Compatibility

| Feature | Plugin SDK v2.0+ | Plugin SDK v2.10+ | Plugin SDK v2.20+ | Notes |
|---------|------------------|-------------------|-------------------|--------|
| Basic Typed Resources | ❌ | ✅ | ✅ | Minimum version for typed resource framework |
| `metadata.Decode()` | ❌ | ✅ | ✅ | State decoding for typed resources |
| `metadata.Encode()` | ❌ | ✅ | ✅ | State encoding for typed resources |
| `metadata.ResourceRequiresImport()` | ❌ | ✅ | ✅ | Import conflict detection |
| `metadata.MarkAsGone()` | ❌ | ✅ | ✅ | Resource deletion handling |
| `metadata.Logger` | ❌ | ✅ | ✅ | Structured logging |
| `sdk.ResourceWithUpdate` | ❌ | ✅ | ✅ | Update operation interface |
| Enhanced `CustomizeDiff` | ✅ | ✅ | ✅ | Available in all v2.x versions |

### AzureRM Provider Framework Evolution

| AzureRM Version | Typed Resources | Migration Support | Dual Registration | Recommendation |
|-----------------|----------------|-------------------|-------------------|----------------|
| v3.0 - v3.50 | ❌ | ❌ | ❌ | Use untyped resources only |
| v3.51 - v3.80 | ⚠️ | ⚠️ | ❌ | Early typed resource support (experimental) |
| v3.81+ | ✅ | ✅ | ✅ | Full typed resource support with migration capabilities |
| v4.0+ (planned) | ✅ | ✅ | ✅ | Preferred typed resource implementation |

### Azure SDK for Go Compatibility

| Azure SDK Version | Pointer Helpers | Response Helpers | Polling Support | Migration Impact |
|-------------------|----------------|------------------|-----------------|------------------|
| HashiCorp Go Azure SDK v0.20+ | ✅ | ✅ | ✅ | Full migration support |
| Azure SDK for Go v68+ | ⚠️ | ✅ | ✅ | Limited pointer helper support |
| Legacy Azure SDK | ❌ | ⚠️ | ⚠️ | Migration not recommended |

### Migration Timeline Recommendations

```go
// Version-specific migration approach
func planMigration(providerVersion string) MigrationStrategy {
    switch {
    case providerVersion >= "v3.81":
        return MigrationStrategy{
            TypedResources:    true,
            DualRegistration: true,
            Recommendation:   "Full migration support available",
        }
    case providerVersion >= "v3.51":
        return MigrationStrategy{
            TypedResources:    true,
            DualRegistration: false,
            Recommendation:   "Experimental - wait for v3.81+ for production use",
        }
    default:
        return MigrationStrategy{
            TypedResources:    false,
            DualRegistration: false,
            Recommendation:   "Upgrade provider version before migration",
        }
    }
}
```

---
[⬆️ Back to top](#migration-guide)

## 🚧 Upgrade Procedures

### Pre-Migration Compatibility Checks

```go
// Version compatibility validation
func validateMigrationCompatibility() error {
    // Check Plugin SDK version
    if !hasPluginSDKVersion("v2.10.0") {
        return fmt.Errorf("migration requires Terraform Plugin SDK v2.10.0 or later")
    }

    // Check AzureRM provider framework
    if !hasTypedResourceSupport() {
        return fmt.Errorf("migration requires AzureRM provider v3.81 or later")
    }

    // Check Azure SDK compatibility
    if !hasPointerHelpers() {
        return fmt.Errorf("migration requires HashiCorp Go Azure SDK v0.20 or later")
    }

    return nil
}
```

### State Migration Patterns

**Schema Version Migration:**
```go
func resourceServiceNameV0() *pluginsdk.Resource {
    return &pluginsdk.Resource{
        Schema: map[string]*pluginsdk.Schema{
            // Old schema definition
            "scrubbing_rule": {
                Type:     pluginsdk.TypeSet,
                Optional: true,
                Elem: &pluginsdk.Resource{
                    Schema: oldScrubbingRuleSchema(),
                },
            },
        },
    }
}

func resourceServiceName() *pluginsdk.Resource {
    return &pluginsdk.Resource{
        SchemaVersion: 1,
        StateUpgraders: []pluginsdk.StateUpgrader{
            {
                Type:    resourceServiceNameV0().CoreConfigSchema().ImpliedType(),
                Upgrade: resourceServiceNameStateUpgradeV0ToV1,
                Version: 0,
            },
        },
        Schema: map[string]*pluginsdk.Schema{
            // New schema definition
            "log_scrubbing_rule": {
                Type:     pluginsdk.TypeSet,
                Optional: true,
                Elem: &pluginsdk.Resource{
                    Schema: newScrubbingRuleSchema(),
                },
            },
        },
    }
}

func resourceServiceNameStateUpgradeV0ToV1(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
    // Migrate field names
    if scrubbingRules, ok := rawState["scrubbing_rule"]; ok {
        rawState["log_scrubbing_rule"] = scrubbingRules
        delete(rawState, "scrubbing_rule")
    }

    return rawState, nil
}
```

### Breaking Change Communication

**Changelog Entry Pattern:**
```markdown
## 4.0.0 (Unreleased)

BREAKING CHANGES:

* **Field Rename**: `azurerm_cdn_frontdoor_profile` - the `scrubbing_rule` field has been renamed to `log_scrubbing_rule` for better clarity ([#12345](https://github.com/hashicorp/terraform-provider-azurerm/pull/12345))

* **Schema Flattening**: `azurerm_cdn_frontdoor_profile` - the `log_scrubbing` wrapper block has been removed, `scrubbing_rule` blocks are now configured directly on the profile ([#12346](https://github.com/hashicorp/terraform-provider-azurerm/pull/12346))

NOTES:

* This release contains significant breaking changes. Please see the [4.0 Upgrade Guide](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/4.0-upgrade-guide) for details on how to upgrade your configurations.
```

**Migration Guide Entry:**
```markdown
# Upgrading to 4.0

## CDN Front Door Profile Changes

### Field Rename: `scrubbing_rule` → `log_scrubbing_rule`

**Before (v3.x):**
```hcl
resource "azurerm_cdn_frontdoor_profile" "example" {
  scrubbing_rule {
    match_variable = "QueryStringArgNames"
  }
}
```

**After (v4.x):**
```hcl
resource "azurerm_cdn_frontdoor_profile" "example" {
  log_scrubbing_rule {
    match_variable = "QueryStringArgNames"
  }
}
```

**Migration Steps:**
1. Update your configuration files to use `log_scrubbing_rule`
2. Run `terraform plan` to verify changes
3. Apply the configuration
```

### Common Migration Pitfalls

**1. Schema Mismatch Issues**
```go
// Problem: tfschema tag doesn't match schema key
type BadModel struct {
    Name string `tfschema:"resource_name"` // Wrong: should be "name"
}

// Solution: Ensure tfschema tags exactly match schema keys
type GoodModel struct {
    Name string `tfschema:"name"` // Correct: matches schema key
}
```

**2. Import Conflict Detection**
```go
// Problem: Using old import conflict pattern
if existing.StatusCode != http.StatusNotFound {
    return tf.ImportAsExistsError("azurerm_resource", id.ID())
}

// Solution: Use metadata pattern
if !response.WasNotFound(existing.HttpResponse) {
    return metadata.ResourceRequiresImport(r.ResourceType(), id)
}
```

**3. State Management**
```go
// Problem: Trying to use d.Set() in typed resource
d.Set("name", model.Name) // Wrong pattern

// Solution: Use metadata.Encode()
return metadata.Encode(&state) // Correct pattern
```

### Post-Migration Verification

**Migration Success Checklist:**
- [ ] All tests pass with new implementation
- [ ] Import functionality preserved
- [ ] State compatibility maintained
- [ ] Documentation updated
- [ ] Performance benchmarks within acceptable range
- [ ] No breaking changes to user configurations (unless intentional)

**Rollback Plan:**
```go
// Maintain ability to revert if critical issues arise
// Keep backup branch: migration/resource-name-to-typed-backup
// Document rollback procedure in pull request description
// Test rollback path before merging
```

---
[⬆️ Back to top](#migration-guide)

---

## Quick Reference Links

- 🏗️ **Main Implementation Guide**: [implementation-guide.instructions.md](./implementation-guide.instructions.md)
- ⚡ **Azure Patterns**: [azure-patterns.instructions.md](./azure-patterns.instructions.md)
- 🚨 **Error Patterns**: [error-patterns.instructions.md](./error-patterns.instructions.md)
- 📋 **Schema Patterns**: [schema-patterns.instructions.md](./schema-patterns.instructions.md)
- 🧪 **Testing Guide**: [testing-guidelines.instructions.md](./testing-guidelines.instructions.md)
- 📝 **Documentation Guide**: [documentation-guidelines.instructions.md](./documentation-guidelines.instructions.md)
