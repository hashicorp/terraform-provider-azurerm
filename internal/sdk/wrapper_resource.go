// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// ResourceWrapper is a wrapper for converting a Resource implementation
// into the object used by the Terraform Plugin SDK
type ResourceWrapper struct {
	logger   Logger
	resource Resource
}

// NewResourceWrapper returns a ResourceWrapper for this Resource implementation
func NewResourceWrapper(resource Resource) ResourceWrapper {
	return ResourceWrapper{
		logger:   &DiagnosticsLogger{},
		resource: resource,
	}
}

// Resource returns the Terraform Plugin SDK type for this Resource implementation
func (rw *ResourceWrapper) Resource() (*schema.Resource, error) {
	resourceSchema, err := combineSchema(rw.resource.Arguments(), rw.resource.Attributes())
	if err != nil {
		return nil, fmt.Errorf("building Schema: %+v", err)
	}

	modelObj := rw.resource.ModelObject()
	if modelObj != nil {
		if err := ValidateModelObject(modelObj); err != nil {
			return nil, fmt.Errorf("validating model for %q: %+v", rw.resource.ResourceType(), err)
		}
	}

	d := func(duration time.Duration) *time.Duration {
		return &duration
	}

	resource := schema.Resource{
		Schema: *resourceSchema,

		CreateContext: rw.diagnosticsWrapper(func(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
			metaData := runArgs(d, meta, rw.logger)
			err := rw.resource.Create().Func(ctx, metaData)
			if err != nil {
				return err
			}
			// NOTE: whilst this may look like we should use the Read
			// functions timeout here, we're still /technically/ in the
			// Create function so reusing that timeout should be sufficient
			return rw.resource.Read().Func(ctx, metaData)
		}),

		// looks like these could be reused, easiest if they're not
		ReadContext: rw.diagnosticsWrapper(func(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
			metaData := runArgs(d, meta, rw.logger)
			return rw.resource.Read().Func(ctx, metaData)
		}),
		DeleteContext: rw.diagnosticsWrapper(func(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
			metaData := runArgs(d, meta, rw.logger)
			return rw.resource.Delete().Func(ctx, metaData)
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: d(rw.resource.Create().Timeout),
			Read:   d(rw.resource.Read().Timeout),
			Delete: d(rw.resource.Delete().Timeout),
		},
		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			fn := rw.resource.IDValidationFunc()
			warnings, errs := fn(id, "id")
			if len(warnings) > 0 {
				for _, warning := range warnings {
					rw.logger.Warn(warning)
				}
			}
			if len(errs) > 0 {
				out := ""
				for _, err := range errs {
					out += err.Error()
				}
				return errors.New(out)
			}

			return nil
		}, func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) ([]*pluginsdk.ResourceData, error) {
			if v, ok := rw.resource.(ResourceWithCustomImporter); ok {
				metaData := runArgs(d, meta, rw.logger)

				ctx, cancel := context.WithTimeout(ctx, rw.resource.Read().Timeout)
				defer cancel()
				err := v.CustomImporter()(ctx, metaData)
				if err != nil {
					return nil, err
				}

				return []*pluginsdk.ResourceData{metaData.ResourceData}, nil
			}

			return schema.ImportStatePassthroughContext(ctx, d, meta)
		}),
	}

	// Not all resources support update - so this is an separate interface
	// implementations can opt to interface
	if v, ok := rw.resource.(ResourceWithUpdate); ok {
		resource.UpdateContext = rw.diagnosticsWrapper(func(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
			metaData := runArgs(d, meta, rw.logger)

			err := v.Update().Func(ctx, metaData)
			if err != nil {
				return err
			}
			// whilst this may look like we should use the Update timeout here
			// we're still "technically" in the update method, so reusing the
			// Update's timeout should be fine
			return rw.resource.Read().Func(ctx, metaData)
		})
		resource.Timeouts.Update = d(v.Update().Timeout)
	}

	if v, ok := rw.resource.(ResourceWithCustomizeDiff); ok {
		resource.CustomizeDiff = func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			client := meta.(*clients.Client)
			ctx, cancel := context.WithTimeout(ctx, v.CustomizeDiff().Timeout)
			defer cancel()
			metaData := ResourceMetaData{
				Client:                   client,
				Logger:                   rw.logger,
				ResourceDiff:             d,
				serializationDebugLogger: NullLogger{},
			}

			return v.CustomizeDiff().Func(ctx, metaData)
		}
	}

	if v, ok := rw.resource.(ResourceWithDeprecationAndNoReplacement); ok {
		resource.DeprecationMessage = v.DeprecationMessage()
	}
	if v, ok := rw.resource.(ResourceWithDeprecationReplacedBy); ok {
		if resource.DeprecationMessage != "" {
			return nil, fmt.Errorf("Resource %q can implement either ResourceWithDeprecationAndNoReplacement or ResourceWithDeprecationReplacedBy but not both", rw.resource.ResourceType())
		}

		replacementResourceType := v.DeprecatedInFavourOfResource()
		if replacementResourceType == "" {
			return nil, fmt.Errorf("Resource %q must return a non-empty DeprecatedInFavourOfResource if implementing ResourceWithDeprecationReplacedBy", rw.resource.ResourceType())
		}

		resource.DeprecationMessage = fmt.Sprintf(`The %[1]q resource has been deprecated and replaced by the %[2]q resource.

The existing %[1]q resource will remain available until the next
major version of the Azure Provider however the existing resource is feature-frozen
and we recommend using the %[2]q resource instead.
`, rw.resource.ResourceType(), replacementResourceType)
	}

	if v, ok := rw.resource.(ResourceWithStateMigration); ok {
		stateUpgradeData := v.StateUpgraders()
		resource.SchemaVersion = stateUpgradeData.SchemaVersion
		resource.StateUpgraders = pluginsdk.StateUpgrades(stateUpgradeData.Upgraders)
	}
	// TODO: State Migrations

	return &resource, nil
}

func (rw *ResourceWrapper) diagnosticsWrapper(in func(ctx context.Context, d *schema.ResourceData, meta interface{}) error) func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diagnosticsWrapper(in, rw.logger)
}

func diagnosticsWrapper(in func(ctx context.Context, d *schema.ResourceData, meta interface{}) error, logger Logger) func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		out := make([]diag.Diagnostic, 0)
		if err := in(ctx, d, meta); err != nil {
			out = append(out, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       err.Error(),
				Detail:        err.Error(),
				AttributePath: nil,
			})
		}

		if diagsLogger, ok := logger.(*DiagnosticsLogger); ok {
			out = append(out, diagsLogger.diagnostics...)
		}

		return out
	}
}
