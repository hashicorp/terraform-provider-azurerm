// Package passes provides analysis passes for the azurerm-linter.
//
// This package contains multiple analyzers that enforce coding standards and best practices
// for the Azure Resource Manager Terraform Provider (azurerm).
//
// # Available Analyzers
//
// # AZBP001 - String Validation Check
//
// Reports when String type schema fields (Required or Optional) do not have a ValidateFunc.
//
// Reference: https://github.com/hashicorp/terraform-provider-azurerm/blob/main/contributing/topics/guide-new-fields-to-resource.md#schema
//
// Flagged:
//
//	"name": {
//	    Type:     pluginsdk.TypeString,
//	    Required: true,
//	    // Missing ValidateFunc!
//	}
//
// Correct:
//
//	"name": {
//	    Type:         pluginsdk.TypeString,
//	    Required:     true,
//	    ValidateFunc: validation.StringIsNotEmpty,
//	}
//
// AZBP002 - Optional+Computed Documentation Check
//
// Reports when schema properties are marked as both Optional and Computed without
// proper documentation explaining why this pattern is necessary.
//
// Reference: https://github.com/hashicorp/terraform-provider-azurerm/blob/main/contributing/topics/best-practices.md#setting-properties-to-optional--computed
//
// Flagged:
//
//	"etag": {
//	    Type:     pluginsdk.TypeString,
//	    Optional: true,
//	    Computed: true,  // Missing NOTE: O+C comment
//	}
//
// Correct:
//
//	"etag": {
//	    Type:     pluginsdk.TypeString,
//	    Optional: true,
//	    // NOTE: O+C Azure generates a new value every time this resource is updated
//	    Computed: true,
//	}
//
// AZSD001 - MaxItems:1 Flattening Check
//
// Reports when blocks with MaxItems: 1 contain only a single nested property without
// proper justification. These should typically be flattened for better user experience.
//
// Flagged:
//
//	"config": {
//	    Type:     pluginsdk.TypeList,
//	    MaxItems: 1,
//	    Elem: &pluginsdk.Resource{
//	        Schema: map[string]*pluginsdk.Schema{
//	            "value": {...},  // Only one property - should be flattened
//	        },
//	    },
//	}
//
// Correct (flattened):
//
//	"config_value": {...}
//
// Correct (with explanation):
//
//	"config": {
//	    Type:     pluginsdk.TypeList,
//	    MaxItems: 1,
//	    // Additional properties will be added per service team confirmation
//	    Elem: &pluginsdk.Resource{
//	        Schema: map[string]*pluginsdk.Schema{
//	            "value": {...},
//	        },
//	    },
//	}
//
// # AZSD002 - AtLeastOneOf Validation for TypeList Fields
//
// Reports when a TypeList block contains multiple optional nested fields but none
// of them have AtLeastOneOf validation set. When all nested fields are optional,
// at least one should use AtLeastOneOf to ensure users specify at least one option.
//
// Flagged:
//
//	"setting": {
//	    Type:     pluginsdk.TypeList,
//	    Optional: true,
//	    MaxItems: 1,
//	    Elem: &pluginsdk.Resource{
//	        Schema: map[string]*pluginsdk.Schema{
//	            "linux": {
//	                Type:     pluginsdk.TypeList,
//	                Optional: true,
//	                // Missing AtLeastOneOf!
//	            },
//	            "windows": {
//	                Type:     pluginsdk.TypeList,
//	                Optional: true,
//	                // Missing AtLeastOneOf!
//	            },
//	        },
//	    },
//	}
//
// Correct:
//
//	"setting": {
//	    Type:     pluginsdk.TypeList,
//	    Optional: true,
//	    MaxItems: 1,
//	    Elem: &pluginsdk.Resource{
//	        Schema: map[string]*pluginsdk.Schema{
//	            "linux": {
//	                Type:     pluginsdk.TypeList,
//	                Optional: true,
//	                AtLeastOneOf: []string{"setting.0.linux", "setting.0.windows"},
//	            },
//	            "windows": {
//	                Type:     pluginsdk.TypeList,
//	                Optional: true,
//	                AtLeastOneOf: []string{"setting.0.linux", "setting.0.windows"},
//	            },
//	        },
//	    },
//	}
//
// # AZBP003 - Enum Conversion Convention
//
// Reports when go-azure-sdk enum types are converted using pointer.To() with
// explicit type conversion instead of the generic pointer.ToEnum[T]() function.
// Using ToEnum is safer as it's specifically designed for enum types.
//
// Flagged:
//
//	return &managedclusters.ManagedClusterBootstrapProfile{
//	    ArtifactSource: pointer.To(managedclusters.ArtifactSource(config["artifact_source"].(string))),
//	}
//
// Correct:
//
//	return &managedclusters.ManagedClusterBootstrapProfile{
//	    ArtifactSource: pointer.ToEnum[managedclusters.ArtifactSource](config["artifact_source"].(string)),
//	}
//
// # AZBP004 - Pointer Dereferencing Convention
//
// Reports when a variable is initialized to its zero value and then conditionally
// assigned via pointer dereferencing. This pattern can be simplified using pointer.From(),
// which is more concise and safely handles nil cases by returning the zero value.
//
// Flagged:
//
//	name := ""
//	if input.Name != nil {
//	    name = *input.Name
//	}
//
// Correct:
//
//	name := pointer.From(input.Name)
//
// # AZRN001 - Percentage Suffix Convention
//
// Reports when percentage properties use '_in_percent' suffix instead of the
// standardized '_percentage' suffix.
//
// Flagged:
//
//	"cpu_threshold_in_percent": {...}
//
// Correct:
//
//	"cpu_threshold_percentage": {...}
//
// # AZRE001 - Error Creation Convention
//
// Reports when fixed error strings (without format placeholders) use fmt.Errorf()
// instead of errors.New().
//
// Flagged:
//
//	return fmt.Errorf("something went wrong")
//
// Correct:
//
//	return errors.New("something went wrong")
//	return fmt.Errorf("value %s is invalid", value)  // with placeholder, OK
//
// # AZNR001 - Schema Field Ordering
//
// Reports when schema fields are not ordered according to the provider's conventions.
// For top-level schemas, checks that name, resource_group_name, and location appear
// in the correct relative order. For nested schemas, checks that required, optional,
// and computed fields are properly grouped and alphabetically sorted.
//
// Reference: https://github.com/hashicorp/terraform-provider-azurerm/blob/main/contributing/topics/guide-new-resource.md
//
// When git filter is enabled, it only validates on newly created files
//
// Required order:
//  1. Special ID fields (name, resource_group_name in order)
//  2. Location field
//  3. Required fields (sorted alphabetically for nested schemas)
//  4. Optional fields (sorted alphabetically)
//  5. Computed fields (sorted alphabetically)
package passes
