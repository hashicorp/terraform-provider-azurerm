package sdk

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceWrapper is a wrapper for converting a DataSource implementation
// into the object used by the Terraform Plugin SDK
type DataSourceWrapper struct {
	dataSource DataSource
	logger     Logger
}

// NewDataSourceWrapper returns a DataSourceWrapper for this Data Source implementation
func NewDataSourceWrapper(dataSource DataSource) DataSourceWrapper {
	return DataSourceWrapper{
		dataSource: dataSource,
		logger:     &DiagnosticsLogger{},
	}
}

// DataSource returns the Terraform Plugin SDK type for this DataSource implementation
func (dw *DataSourceWrapper) DataSource() (*schema.Resource, error) {
	resourceSchema, err := combineSchema(dw.dataSource.Arguments(), dw.dataSource.Attributes())
	if err != nil {
		return nil, fmt.Errorf("building Schema: %+v", err)
	}

	modelObj := dw.dataSource.ModelObject()
	if modelObj != nil {
		if err := ValidateModelObject(&modelObj); err != nil {
			return nil, fmt.Errorf("validating model for %q: %+v", dw.dataSource.ResourceType(), err)
		}
	}

	d := func(duration time.Duration) *time.Duration {
		return &duration
	}

	resource := schema.Resource{
		Schema: *resourceSchema,
		ReadContext: dw.diagnosticsWrapper(func(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
			metaData := runArgs(d, meta, dw.logger)
			return dw.dataSource.Read().Func(ctx, metaData)
		}),
		Timeouts: &schema.ResourceTimeout{
			Read: d(dw.dataSource.Read().Timeout),
		},
	}

	return &resource, nil
}

func (dw *DataSourceWrapper) diagnosticsWrapper(in func(ctx context.Context, d *schema.ResourceData, meta interface{}) error) schema.ReadContextFunc {
	return diagnosticsWrapper(in, dw.logger)
}
