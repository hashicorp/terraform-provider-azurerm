// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package typed_upgrade

import "fmt"

// Plan is a read-only assessment of what a typed upgrade would do.
type Plan struct {
	// Supported reports whether the file can be upgraded.
	Supported bool
	// Reason is set when Supported is false.
	Reason string

	// Metadata derived from analysis.
	TerraformType string
	StructName    string
	ModelName     string
	HasUpdate     bool

	// Warnings is a list of migration notes the caller should present to the user.
	Warnings []string
}

// PlanUpgrade assesses the resource and reports what a typed upgrade would
// change, without mutating anything.
func PlanUpgrade(info *Info) Plan {
	p := Plan{
		Supported:     true,
		TerraformType: info.TerraformType,
		StructName:    info.StructName,
		ModelName:     info.ModelName,
		HasUpdate:     info.HasUpdate,
	}

	if info.SDKPackage == "" {
		p.Warnings = append(p.Warnings, "could not determine the go-azure-sdk package; IDValidationFunc may need manual adjustment")
	}
	if info.IDValidateFunc == "" {
		p.Warnings = append(p.Warnings, "could not determine the ID validation function; set IDValidationFunc manually after generation")
	}
	if len(info.Fields) == 0 {
		p.Warnings = append(p.Warnings, "no schema fields extracted; model fields will need to be added manually")
	}
	if info.HasCustomizeDiff {
		p.Warnings = append(p.Warnings, "CustomizeDiff detected — the generated stub panics; migrate the diff logic manually")
	}
	if info.HasSchemaVersion && info.SchemaVersion > 0 {
		p.Warnings = append(p.Warnings,
			fmt.Sprintf("SchemaVersion %d detected — carry over the state upgraders and keep SchemaVersion=%d in the typed resource", info.SchemaVersion, info.SchemaVersion))
	}
	if info.DeprecationMsg != "" {
		p.Warnings = append(p.Warnings, "resource is deprecated — include the DeprecationMessage in the typed resource documentation")
	}

	// Block-related warnings.
	p.Warnings = append(p.Warnings, PlanBlockWarnings(info)...)

	return p
}

// UpgradeOptions configures what Upgrade produces.
type UpgradeOptions struct {
	// Write causes Upgrade to write the generated file and rename the original.
	// When false (dry run), only the generated source is returned.
	Write bool

	// Overwrite allows Upgrade to overwrite the generated file when it already
	// exists.
	Overwrite bool

	// UpdateRegistration causes Upgrade to also update the service package's
	// registration.go: adding the new struct to Resources() and removing it
	// from SupportedResources().
	UpdateRegistration bool

	// RegistrationPath is the path to the registration.go file. When empty and
	// UpdateRegistration is true, Upgrade derives it from the resource file path.
	RegistrationPath string

	// TerraformType overrides the derived terraform resource type string
	// (e.g. "azurerm_maps_creator"). Useful when the derived value is wrong.
	TerraformType string

	// Pandora resolution options — all optional. When ARMType (or Service +
	// Resource) is provided and the Pandora server is reachable, the IR is
	// resolved and used to generate typed expand/flatten helpers for nested
	// block fields.
	PandoraURL string
	ARMType    string
	Service    string
	Resource   string
	APIVersion string
}

// Result holds the output of a typed upgrade.
type Result struct {
	// GeneratedPath is the path the new typed resource was written to (or would
	// be written to in a dry run).
	GeneratedPath string

	// RenamedPath is the path the original resource file was renamed to (or
	// would be renamed).
	RenamedPath string

	// GeneratedSrc is the raw source of the generated typed resource.
	GeneratedSrc string

	// Plan is the assessment that was used.
	Plan Plan
}

// Upgrade generates a typed resource from an analyzed untyped resource.
// When opts.Write is false, it is a dry run: the generated source is returned
// in Result.GeneratedSrc but no files are modified.
func Upgrade(info *Info, opts UpgradeOptions) (*Result, error) {
	if opts.TerraformType != "" {
		info.TerraformType = opts.TerraformType
	}

	// Optionally resolve the Pandora IR for nested block expand/flatten generation.
	if opts.ARMType != "" || (opts.Service != "" && opts.Resource != "") {
		if err := ResolveBlocks(info, BlockResolveOptions{
			PandoraURL: opts.PandoraURL,
			ARMType:    opts.ARMType,
			Service:    opts.Service,
			Resource:   opts.Resource,
			APIVersion: opts.APIVersion,
		}); err != nil {
			// Non-fatal: degrade gracefully without block helpers.
			_ = err
		}
		// Re-run model transforms now that IR-derived block names are available.
		if info.PandoraIR != nil {
			info.applyModelTransforms()
		}
	}

	plan := PlanUpgrade(info)
	if !plan.Supported {
		return nil, fmt.Errorf("%s", plan.Reason)
	}

	src, err := Generate(info)
	if err != nil {
		return nil, fmt.Errorf("generating typed resource: %w", err)
	}

	// Derive file paths.
	generatedPath, renamedPath := derivePaths(info.Path)

	res := &Result{
		GeneratedPath: generatedPath,
		RenamedPath:   renamedPath,
		GeneratedSrc:  src,
		Plan:          plan,
	}

	if !opts.Write {
		return res, nil
	}

	if err := writeFiles(info.Path, generatedPath, renamedPath, []byte(src), opts.Overwrite); err != nil {
		return res, err
	}

	if opts.UpdateRegistration {
		regPath := opts.RegistrationPath
		if regPath == "" {
			regPath = deriveRegistrationPath(info.Path)
		}
		if regPath != "" {
			newSrc, changed, err := RegisterTypedResource(regPath, info.StructName, info.TerraformType)
			if err != nil {
				return res, fmt.Errorf("updating registration: %w", err)
			}
			if changed {
				if err := writeFile(regPath, newSrc, true); err != nil {
					return res, fmt.Errorf("writing registration: %w", err)
				}
			}
		}
	}

	return res, nil
}
