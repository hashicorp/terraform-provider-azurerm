package sdk

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
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
		if err := ValidateModelObject(&modelObj); err != nil {
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
			warnings, errors := fn(id, "id")
			if len(warnings) > 0 {
				for _, warning := range warnings {
					rw.logger.Warn(warning)
				}
			}
			if len(errors) > 0 {
				out := ""
				for _, error := range errors {
					out += error.Error()
				}
				return fmt.Errorf(out)
			}

			return nil
		}, func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) ([]*pluginsdk.ResourceData, error) {
			if v, ok := rw.resource.(ResourceWithCustomImporter); ok {
				metaData := runArgs(d, meta, rw.logger)

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

	if v, ok := rw.resource.(ResourceWithDeprecation); ok {
		message := v.DeprecationMessage()
		if message == "" {
			return nil, fmt.Errorf("Resource %q must return a non-empty DeprecationMessage if implementing ResourceWithDeprecation", rw.resource.ResourceType())
		}

		resource.DeprecationMessage = message
	}

	// TODO: CustomizeDiff
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
