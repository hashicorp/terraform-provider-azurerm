package sdk

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/framework/commonschema"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type FrameworkDataSourceWrapper struct {
	ResourceMetadata

	FrameworkWrappedDataSource

	Model interface{}
}

var _ datasource.DataSourceWithConfigure = &FrameworkDataSourceWrapper{}

func (d *FrameworkDataSourceWrapper) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	d.FrameworkWrappedDataSource.Schema(ctx, request, response)
	response.Schema.Attributes["id"] = commonschema.IDAttributeDataSource()

	if response.Schema.Blocks == nil {
		response.Schema.Blocks = map[string]schema.Block{}
	}

	response.Schema.Blocks["timeouts"] = timeouts.Block(ctx)
}

func (d *FrameworkDataSourceWrapper) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	d.DefaultsDataSource(request, response)

	d.Model = d.ModelObject()
	if _, ok := d.FrameworkWrappedDataSource.(FrameworkWrappedDataSourceWithConfigure); ok {
		d.FrameworkWrappedDataSource.(FrameworkWrappedDataSourceWithConfigure).Configure(ctx, request, response, d.ResourceMetadata)
	}
}

func (d *FrameworkDataSourceWrapper) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = d.FrameworkWrappedDataSource.ResourceType()
}

func (d *FrameworkDataSourceWrapper) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	ctx, cancel := context.WithTimeout(ctx, d.ResourceMetadata.TimeoutRead)
	defer cancel()

	config := d.FrameworkWrappedDataSource.ModelObject()

	if ok := d.ResourceMetadata.DecodeDataSourceRead(ctx, request, response, config); !ok {
		return
	}

	d.FrameworkWrappedDataSource.Read(ctx, request, response, d.ResourceMetadata, config)

	d.ResourceMetadata.EncodeDataSourceRead(ctx, response, config)
}

func (d *FrameworkDataSourceWrapper) DataSource() func() datasource.DataSource {
	return func() datasource.DataSource {
		return d
	}
}
