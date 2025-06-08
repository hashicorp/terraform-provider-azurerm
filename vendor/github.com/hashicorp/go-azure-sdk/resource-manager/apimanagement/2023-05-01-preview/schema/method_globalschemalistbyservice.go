package schema

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GlobalSchemaListByServiceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]GlobalSchemaContract
}

type GlobalSchemaListByServiceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []GlobalSchemaContract
}

type GlobalSchemaListByServiceOperationOptions struct {
	Filter *string
	Skip   *int64
	Top    *int64
}

func DefaultGlobalSchemaListByServiceOperationOptions() GlobalSchemaListByServiceOperationOptions {
	return GlobalSchemaListByServiceOperationOptions{}
}

func (o GlobalSchemaListByServiceOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GlobalSchemaListByServiceOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GlobalSchemaListByServiceOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type GlobalSchemaListByServiceCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GlobalSchemaListByServiceCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GlobalSchemaListByService ...
func (c SchemaClient) GlobalSchemaListByService(ctx context.Context, id ServiceId, options GlobalSchemaListByServiceOperationOptions) (result GlobalSchemaListByServiceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &GlobalSchemaListByServiceCustomPager{},
		Path:          fmt.Sprintf("%s/schemas", id.ID()),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.ExecutePaged(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var values struct {
		Values *[]GlobalSchemaContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GlobalSchemaListByServiceComplete retrieves all the results into a single object
func (c SchemaClient) GlobalSchemaListByServiceComplete(ctx context.Context, id ServiceId, options GlobalSchemaListByServiceOperationOptions) (GlobalSchemaListByServiceCompleteResult, error) {
	return c.GlobalSchemaListByServiceCompleteMatchingPredicate(ctx, id, options, GlobalSchemaContractOperationPredicate{})
}

// GlobalSchemaListByServiceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SchemaClient) GlobalSchemaListByServiceCompleteMatchingPredicate(ctx context.Context, id ServiceId, options GlobalSchemaListByServiceOperationOptions, predicate GlobalSchemaContractOperationPredicate) (result GlobalSchemaListByServiceCompleteResult, err error) {
	items := make([]GlobalSchemaContract, 0)

	resp, err := c.GlobalSchemaListByService(ctx, id, options)
	if err != nil {
		result.LatestHttpResponse = resp.HttpResponse
		err = fmt.Errorf("loading results: %+v", err)
		return
	}
	if resp.Model != nil {
		for _, v := range *resp.Model {
			if predicate.Matches(v) {
				items = append(items, v)
			}
		}
	}

	result = GlobalSchemaListByServiceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
