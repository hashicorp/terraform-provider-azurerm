// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
		if err := ValidateModelObject(modelObj); err != nil {
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

	if v, ok := dw.dataSource.(DataSourceWithDeprecationReplacedBy); ok {
		replacementDataSourceType := v.DeprecatedInFavourOfDataSource()
		if replacementDataSourceType == "" {
			return nil, fmt.Errorf("datasource %q must return a non-empty DeprecatedInFavourOfDataSource if implementing DataSourceWithDeprecationReplacedBy", dw.dataSource.ResourceType())
		}
		resource.DeprecationMessage = fmt.Sprintf(`The %[1]q datasource has been deprecated and replaced by the %[2]q datasource.

The existing %[1]q datasource will remain available until the next
major version of the Azure Provider however the existing datasource is feature-frozen
and we recommend using the %[2]q datasource instead.
`, dw.dataSource.ResourceType(), replacementDataSourceType)
	}

	return &resource, nil
}

func (dw *DataSourceWrapper) diagnosticsWrapper(in func(ctx context.Context, d *schema.ResourceData, meta interface{}) error) schema.ReadContextFunc {
	return diagnosticsWrapper(in, dw.logger)
}
