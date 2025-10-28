---
applyTo: "internal/**/*.go"
description: API evolution and versioning patterns for the Terraform AzureRM provider including handling Azure API changes, backward compatibility, and migration strategies.
---

# API Evolution Patterns

API evolution and versioning patterns for the Terraform AzureRM provider including handling Azure API changes, backward compatibility, and migration strategies.

**Quick navigation:** [🔄 Version Management](#🔄-azure-api-version-management) | [⬆️ Backward Compatibility](#⬆️-backward-compatibility-patterns) | [🚀 Migration Strategies](#🚀-migration-strategies) | [📊 Deprecation](#📊-deprecation-management)

## 🔄 Azure API Version Management

### Handling Azure API Version Changes

```go
// Handling Azure API version changes
func expandResourceV2(input []interface{}, apiVersion string) interface{} {
    switch apiVersion {
    case "2023-01-01":
        return expandResourceV2_20230101(input)
    case "2024-01-01":
        return expandResourceV2_20240101(input)
    default:
        return expandResourceV2_Latest(input)
    }
}

// Version-specific expansion
func expandResourceV2_20230101(input []interface{}) *azuretype.ResourceV1 {
    if len(input) == 0 || input[0] == nil {
        return nil
    }

    raw := input[0].(map[string]interface{})

    return &azuretype.ResourceV1{
        // Legacy field mapping for older API version
        LegacyField: pointer.To(raw["legacy_field"].(string)),
        CommonField: pointer.To(raw["common_field"].(string)),
    }
}

func expandResourceV2_20240101(input []interface{}) *azuretype.ResourceV2 {
    if len(input) == 0 || input[0] == nil {
        return nil
    }

    raw := input[0].(map[string]interface{})

    return &azuretype.ResourceV2{
        // New field mapping for current API version
        NewField:    pointer.To(raw["new_field"].(string)),
        CommonField: pointer.To(raw["common_field"].(string)),
        // Enhanced properties not available in v1
        EnhancedProperties: expandEnhancedProperties(raw["enhanced_properties"].([]interface{})),
    }
}
```

### API Version Detection

```go
func detectAPIVersion(ctx context.Context, client *azuretype.Client) (string, error) {
    // Try latest API version first
    latestVersion := "2024-01-01"
    if err := testAPIVersion(ctx, client, latestVersion); err == nil {
        return latestVersion, nil
    }

    // Fallback to previous versions
    fallbackVersions := []string{"2023-06-01", "2023-01-01", "2022-12-01"}
    for _, version := range fallbackVersions {
        if err := testAPIVersion(ctx, client, version); err == nil {
            return version, nil
        }
    }

    return "", fmt.Errorf("no compatible API version found")
}

func testAPIVersion(ctx context.Context, client *azuretype.Client, apiVersion string) error {
    // Test API version by making a lightweight call
    _, err := client.GetWithAPIVersion(ctx, "test", apiVersion)
    return err
}
```

### Feature Flag Management

```go
type FeatureFlags struct {
    UseNewAPI       bool
    EnablePreview   bool
    BackwardCompat  bool
}

func (r ServiceResource) Create() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            flags := getFeatureFlags(metadata)

            if flags.UseNewAPI {
                return r.createWithNewAPI(ctx, metadata)
            }

            return r.createWithLegacyAPI(ctx, metadata)
        },
    }
}

func getFeatureFlags(metadata sdk.ResourceMetaData) FeatureFlags {
    // Extract feature flags from provider configuration
    return FeatureFlags{
        UseNewAPI:      metadata.Client.Features.UseLatestAPI,
        EnablePreview:  metadata.Client.Features.EnablePreviewFeatures,
        BackwardCompat: metadata.Client.Features.MaintainBackwardCompatibility,
    }
}
```

## ⬆️ Backward Compatibility Patterns

### Legacy Field Support

```go
// Backward compatibility patterns
func flattenResourceWithCompatibility(input *azuretype.Resource) []interface{} {
    result := flattenResourceStandard(input)

    // Handle deprecated fields for state compatibility
    if input.LegacyField != nil {
        // Maintain backward compatibility while migrating
        result[0].(map[string]interface{})["legacy_field"] = pointer.From(input.LegacyField)
    }

    return result
}

// Schema with deprecated field support
func resourceSchema() map[string]*pluginsdk.Schema {
    return map[string]*pluginsdk.Schema{
        "new_field": {
            Type:     pluginsdk.TypeString,
            Optional: true,
        },
        "legacy_field": {
            Type:       pluginsdk.TypeString,
            Optional:   true,
            Deprecated: "This field is deprecated. Use 'new_field' instead.",
        },
    }
}
```

### State Migration Handling

```go
func migrateStateV0ToV1(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
    // Handle field renames
    if legacyValue, exists := rawState["legacy_field"]; exists {
        rawState["new_field"] = legacyValue
        delete(rawState, "legacy_field")
    }

    // Handle structure changes
    if oldConfig, exists := rawState["old_config"]; exists {
        newConfig := migrateConfigStructure(oldConfig)
        rawState["new_config"] = newConfig
        delete(rawState, "old_config")
    }

    // Update schema version
    rawState["schema_version"] = 1

    return rawState, nil
}

func migrateConfigStructure(oldConfig interface{}) interface{} {
    // Convert old configuration structure to new format
    oldList := oldConfig.([]interface{})
    if len(oldList) == 0 {
        return []interface{}{}
    }

    oldMap := oldList[0].(map[string]interface{})

    // Transform to new structure
    newMap := map[string]interface{}{
        "setting":     oldMap["old_setting"],
        "new_setting": "default_value", // New field with default
    }

    return []interface{}{newMap}
}
```

### Gradual Migration Strategy

```go
func (r ServiceResource) Update() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            var model ServiceResourceModel
            if err := metadata.Decode(&model); err != nil {
                return fmt.Errorf("decoding: %+v", err)
            }

            // Check if resource is using legacy configuration
            if isLegacyConfiguration(model) {
                // Migrate to new configuration during update
                model = migrateLegacyConfiguration(model)

                // Log migration for tracking
                metadata.Logger.Infof("Migrated legacy configuration for %s", model.Id)
            }

            // Use new API for update
            return r.updateWithNewAPI(ctx, metadata, model)
        },
    }
}

func isLegacyConfiguration(model ServiceResourceModel) bool {
    // Check for presence of legacy fields or patterns
    return model.LegacyField != "" || len(model.OldConfiguration) > 0
}

func migrateLegacyConfiguration(model ServiceResourceModel) ServiceResourceModel {
    // Transform legacy configuration to new format
    if model.LegacyField != "" {
        model.NewField = model.LegacyField
        model.LegacyField = ""
    }

    // Migrate configuration blocks
    if len(model.OldConfiguration) > 0 {
        model.NewConfiguration = transformConfigurationStructure(model.OldConfiguration)
        model.OldConfiguration = nil
    }

    return model
}
```

## 🚀 Migration Strategies

### Phased Migration Approach

```go
type MigrationPhase int

const (
    PhaseDetection MigrationPhase = iota
    PhaseWarning
    PhaseMigration
    PhaseEnforcement
)

func (r ServiceResource) getCurrentMigrationPhase() MigrationPhase {
    // Determine current migration phase based on timeline
    now := time.Now()

    switch {
    case now.Before(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)):
        return PhaseDetection
    case now.Before(time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC)):
        return PhaseWarning
    case now.Before(time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)):
        return PhaseMigration
    default:
        return PhaseEnforcement
    }
}

func (r ServiceResource) handleMigrationPhase(ctx context.Context, metadata sdk.ResourceMetaData) error {
    phase := r.getCurrentMigrationPhase()

    switch phase {
    case PhaseDetection:
        return r.detectLegacyUsage(metadata)
    case PhaseWarning:
        return r.warnLegacyUsage(metadata)
    case PhaseMigration:
        return r.migrateLegacyUsage(ctx, metadata)
    case PhaseEnforcement:
        return r.enforceMigration(metadata)
    }

    return nil
}
```

### Automated Migration Tools

```go
func generateMigrationScript(oldConfig, newConfig interface{}) (string, error) {
    script := `# Terraform Configuration Migration Script
# Generated: ` + time.Now().Format(time.RFC3339) + `

# Replace the following configuration:
# OLD:
/*
` + formatConfiguration(oldConfig) + `
*/

# NEW:
/*
` + formatConfiguration(newConfig) + `
*/

# Migration commands:
terraform state replace-provider old/provider new/provider
terraform plan  # Review changes
terraform apply # Apply migration
`

    return script, nil
}

func analyzeConfigurationChanges(old, new interface{}) ConfigurationDiff {
    diff := ConfigurationDiff{
        FieldChanges: make(map[string]FieldChange),
        Warnings:     []string{},
        Errors:       []string{},
    }

    // Analyze field-by-field differences
    oldFields := extractFields(old)
    newFields := extractFields(new)

    for fieldName, oldValue := range oldFields {
        if newValue, exists := newFields[fieldName]; exists {
            if !reflect.DeepEqual(oldValue, newValue) {
                diff.FieldChanges[fieldName] = FieldChange{
                    Old: oldValue,
                    New: newValue,
                    Type: "modified",
                }
            }
        } else {
            diff.FieldChanges[fieldName] = FieldChange{
                Old: oldValue,
                Type: "removed",
            }
            diff.Warnings = append(diff.Warnings, fmt.Sprintf("Field '%s' will be removed", fieldName))
        }
    }

    return diff
}
```

## 📊 Deprecation Management

### Deprecation Timeline

```go
type DeprecationSchedule struct {
    FeatureName    string
    DeprecatedIn   string // Provider version
    WarningPhase   time.Time
    RemovalPhase   time.Time
    ReplacementAPI string
}

var deprecationSchedule = []DeprecationSchedule{
    {
        FeatureName:    "legacy_field",
        DeprecatedIn:   "3.80.0",
        WarningPhase:   time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
        RemovalPhase:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
        ReplacementAPI: "new_field",
    },
}

func checkDeprecationStatus(featureName string) (DeprecationSchedule, bool) {
    for _, schedule := range deprecationSchedule {
        if schedule.FeatureName == featureName {
            return schedule, true
        }
    }
    return DeprecationSchedule{}, false
}
```

### Deprecation Warnings

```go
func validateDeprecatedFeatures(metadata sdk.ResourceMetaData, model interface{}) []string {
    warnings := []string{}

    // Check for deprecated fields
    if hasDeprecatedField(model, "legacy_field") {
        schedule, _ := checkDeprecationStatus("legacy_field")
        warning := fmt.Sprintf(
            "Field 'legacy_field' is deprecated as of provider version %s. "+
            "It will be removed on %s. Use '%s' instead.",
            schedule.DeprecatedIn,
            schedule.RemovalPhase.Format("2006-01-02"),
            schedule.ReplacementAPI,
        )
        warnings = append(warnings, warning)
    }

    return warnings
}

func logDeprecationWarnings(warnings []string, logger interface{}) {
    for _, warning := range warnings {
        // Log with appropriate severity
        if logger != nil {
            if sdkLogger, ok := logger.(sdk.Logger); ok {
                sdkLogger.Warnf("DEPRECATION: %s", warning)
            }
        }
    }
}
```

### Breaking Change Communication

```go
func generateBreakingChangeNotice(changes []BreakingChange) string {
    notice := fmt.Sprintf(`
BREAKING CHANGES DETECTED
========================

The following breaking changes will affect your configuration:

`)

    for i, change := range changes {
        notice += fmt.Sprintf(`
%d. %s
   Impact: %s
   Action Required: %s
   Migration Deadline: %s

`, i+1, change.Description, change.Impact, change.ActionRequired, change.Deadline.Format("2006-01-02"))
    }

    notice += `
For detailed migration instructions, see:
https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/migration-guide
`

    return notice
}

type BreakingChange struct {
    Description    string
    Impact         string
    ActionRequired string
    Deadline       time.Time
}
```

---
[⬆️ Back to top](#api-evolution-patterns)
