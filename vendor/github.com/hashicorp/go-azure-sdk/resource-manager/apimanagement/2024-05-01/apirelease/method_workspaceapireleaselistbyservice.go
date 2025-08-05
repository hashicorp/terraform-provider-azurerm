package apirelease

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceApiReleaseListByServiceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ApiReleaseContract
}

type WorkspaceApiReleaseListByServiceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ApiReleaseContract
}

type WorkspaceApiReleaseListByServiceOperationOptions struct {
	Filter *string
	Skip   *int64
	Top    *int64
}

func DefaultWorkspaceApiReleaseListByServiceOperationOptions() WorkspaceApiReleaseListByServiceOperationOptions {
	return WorkspaceApiReleaseListByServiceOperationOptions{}
}

func (o WorkspaceApiReleaseListByServiceOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o WorkspaceApiReleaseListByServiceOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o WorkspaceApiReleaseListByServiceOperationOptions) ToQuery() *client.QueryParams {
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

type WorkspaceApiReleaseListByServiceCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *WorkspaceApiReleaseListByServiceCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// WorkspaceApiReleaseListByService ...
func (c ApiReleaseClient) WorkspaceApiReleaseListByService(ctx context.Context, id WorkspaceApiId, options WorkspaceApiReleaseListByServiceOperationOptions) (result WorkspaceApiReleaseListByServiceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &WorkspaceApiReleaseListByServiceCustomPager{},
		Path:          fmt.Sprintf("%s/releases", id.ID()),
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
		Values *[]ApiReleaseContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkspaceApiReleaseListByServiceComplete retrieves all the results into a single object
func (c ApiReleaseClient) WorkspaceApiReleaseListByServiceComplete(ctx context.Context, id WorkspaceApiId, options WorkspaceApiReleaseListByServiceOperationOptions) (WorkspaceApiReleaseListByServiceCompleteResult, error) {
	return c.WorkspaceApiReleaseListByServiceCompleteMatchingPredicate(ctx, id, options, ApiReleaseContractOperationPredicate{})
}

// WorkspaceApiReleaseListByServiceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ApiReleaseClient) WorkspaceApiReleaseListByServiceCompleteMatchingPredicate(ctx context.Context, id WorkspaceApiId, options WorkspaceApiReleaseListByServiceOperationOptions, predicate ApiReleaseContractOperationPredicate) (result WorkspaceApiReleaseListByServiceCompleteResult, err error) {
	items := make([]ApiReleaseContract, 0)

	resp, err := c.WorkspaceApiReleaseListByService(ctx, id, options)
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

	result = WorkspaceApiReleaseListByServiceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
