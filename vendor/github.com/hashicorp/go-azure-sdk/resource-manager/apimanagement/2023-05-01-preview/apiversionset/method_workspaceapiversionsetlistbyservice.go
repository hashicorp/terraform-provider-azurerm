package apiversionset

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceApiVersionSetListByServiceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ApiVersionSetContract
}

type WorkspaceApiVersionSetListByServiceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ApiVersionSetContract
}

type WorkspaceApiVersionSetListByServiceOperationOptions struct {
	Filter *string
	Skip   *int64
	Top    *int64
}

func DefaultWorkspaceApiVersionSetListByServiceOperationOptions() WorkspaceApiVersionSetListByServiceOperationOptions {
	return WorkspaceApiVersionSetListByServiceOperationOptions{}
}

func (o WorkspaceApiVersionSetListByServiceOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o WorkspaceApiVersionSetListByServiceOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o WorkspaceApiVersionSetListByServiceOperationOptions) ToQuery() *client.QueryParams {
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

type WorkspaceApiVersionSetListByServiceCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *WorkspaceApiVersionSetListByServiceCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// WorkspaceApiVersionSetListByService ...
func (c ApiVersionSetClient) WorkspaceApiVersionSetListByService(ctx context.Context, id WorkspaceId, options WorkspaceApiVersionSetListByServiceOperationOptions) (result WorkspaceApiVersionSetListByServiceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &WorkspaceApiVersionSetListByServiceCustomPager{},
		Path:          fmt.Sprintf("%s/apiVersionSets", id.ID()),
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
		Values *[]ApiVersionSetContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkspaceApiVersionSetListByServiceComplete retrieves all the results into a single object
func (c ApiVersionSetClient) WorkspaceApiVersionSetListByServiceComplete(ctx context.Context, id WorkspaceId, options WorkspaceApiVersionSetListByServiceOperationOptions) (WorkspaceApiVersionSetListByServiceCompleteResult, error) {
	return c.WorkspaceApiVersionSetListByServiceCompleteMatchingPredicate(ctx, id, options, ApiVersionSetContractOperationPredicate{})
}

// WorkspaceApiVersionSetListByServiceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ApiVersionSetClient) WorkspaceApiVersionSetListByServiceCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, options WorkspaceApiVersionSetListByServiceOperationOptions, predicate ApiVersionSetContractOperationPredicate) (result WorkspaceApiVersionSetListByServiceCompleteResult, err error) {
	items := make([]ApiVersionSetContract, 0)

	resp, err := c.WorkspaceApiVersionSetListByService(ctx, id, options)
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

	result = WorkspaceApiVersionSetListByServiceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
