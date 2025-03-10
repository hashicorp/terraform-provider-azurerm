package apirevision

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceApiRevisionListByServiceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ApiRevisionContract
}

type WorkspaceApiRevisionListByServiceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ApiRevisionContract
}

type WorkspaceApiRevisionListByServiceOperationOptions struct {
	Filter *string
	Skip   *int64
	Top    *int64
}

func DefaultWorkspaceApiRevisionListByServiceOperationOptions() WorkspaceApiRevisionListByServiceOperationOptions {
	return WorkspaceApiRevisionListByServiceOperationOptions{}
}

func (o WorkspaceApiRevisionListByServiceOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o WorkspaceApiRevisionListByServiceOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o WorkspaceApiRevisionListByServiceOperationOptions) ToQuery() *client.QueryParams {
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

type WorkspaceApiRevisionListByServiceCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *WorkspaceApiRevisionListByServiceCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// WorkspaceApiRevisionListByService ...
func (c ApiRevisionClient) WorkspaceApiRevisionListByService(ctx context.Context, id WorkspaceApiId, options WorkspaceApiRevisionListByServiceOperationOptions) (result WorkspaceApiRevisionListByServiceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &WorkspaceApiRevisionListByServiceCustomPager{},
		Path:          fmt.Sprintf("%s/revisions", id.ID()),
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
		Values *[]ApiRevisionContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkspaceApiRevisionListByServiceComplete retrieves all the results into a single object
func (c ApiRevisionClient) WorkspaceApiRevisionListByServiceComplete(ctx context.Context, id WorkspaceApiId, options WorkspaceApiRevisionListByServiceOperationOptions) (WorkspaceApiRevisionListByServiceCompleteResult, error) {
	return c.WorkspaceApiRevisionListByServiceCompleteMatchingPredicate(ctx, id, options, ApiRevisionContractOperationPredicate{})
}

// WorkspaceApiRevisionListByServiceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ApiRevisionClient) WorkspaceApiRevisionListByServiceCompleteMatchingPredicate(ctx context.Context, id WorkspaceApiId, options WorkspaceApiRevisionListByServiceOperationOptions, predicate ApiRevisionContractOperationPredicate) (result WorkspaceApiRevisionListByServiceCompleteResult, err error) {
	items := make([]ApiRevisionContract, 0)

	resp, err := c.WorkspaceApiRevisionListByService(ctx, id, options)
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

	result = WorkspaceApiRevisionListByServiceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
