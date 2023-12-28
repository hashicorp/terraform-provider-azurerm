package elasticsanskus

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SkusListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *SkuInformationList
}

type SkusListOperationOptions struct {
	Filter *string
}

func DefaultSkusListOperationOptions() SkusListOperationOptions {
	return SkusListOperationOptions{}
}

func (o SkusListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o SkusListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o SkusListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

// SkusList ...
func (c ElasticSanSkusClient) SkusList(ctx context.Context, id commonids.SubscriptionId, options SkusListOperationOptions) (result SkusListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/providers/Microsoft.ElasticSan/skus", id.ID()),
		OptionsObject: options,
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	if err = resp.Unmarshal(&result.Model); err != nil {
		return
	}

	return
}
