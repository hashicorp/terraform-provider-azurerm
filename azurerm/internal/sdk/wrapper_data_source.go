package sdk

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
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
		logger:     ConsoleLogger{},
	}
}

// DataSource returns the Terraform Plugin SDK type for this DataSource implementation
func (rw *DataSourceWrapper) DataSource() (*schema.Resource, error) {
	resourceSchema, err := combineSchema(rw.dataSource.Arguments(), rw.dataSource.Attributes())
	if err != nil {
		return nil, fmt.Errorf("building Schema: %+v", err)
	}

	modelObj := rw.dataSource.ModelObject()
	if err := ValidateModelObject(&modelObj); err != nil {
		return nil, fmt.Errorf("validating model for %q: %+v", rw.dataSource.ResourceType(), err)
	}

	var d = func(duration time.Duration) *time.Duration {
		return &duration
	}

	resource := schema.Resource{
		Schema: *resourceSchema,
		Read: func(d *schema.ResourceData, meta interface{}) error {
			ctx, metaData := runArgs(d, meta, rw.logger)
			wrappedCtx, cancel := timeouts.ForRead(ctx, d)
			defer cancel()
			return rw.dataSource.Read().Func(wrappedCtx, metaData)
		},
		Timeouts: &schema.ResourceTimeout{
			Read: d(rw.dataSource.Read().Timeout),
		},
	}

	return &resource, nil
}
