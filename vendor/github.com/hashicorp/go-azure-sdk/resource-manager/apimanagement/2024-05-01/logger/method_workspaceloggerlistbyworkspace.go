package logger

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceLoggerListByWorkspaceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]LoggerContract
}

type WorkspaceLoggerListByWorkspaceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []LoggerContract
}

type WorkspaceLoggerListByWorkspaceOperationOptions struct {
	Filter *string
	Skip   *int64
	Top    *int64
}

func DefaultWorkspaceLoggerListByWorkspaceOperationOptions() WorkspaceLoggerListByWorkspaceOperationOptions {
	return WorkspaceLoggerListByWorkspaceOperationOptions{}
}

func (o WorkspaceLoggerListByWorkspaceOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o WorkspaceLoggerListByWorkspaceOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o WorkspaceLoggerListByWorkspaceOperationOptions) ToQuery() *client.QueryParams {
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

type WorkspaceLoggerListByWorkspaceCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *WorkspaceLoggerListByWorkspaceCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// WorkspaceLoggerListByWorkspace ...
func (c LoggerClient) WorkspaceLoggerListByWorkspace(ctx context.Context, id WorkspaceId, options WorkspaceLoggerListByWorkspaceOperationOptions) (result WorkspaceLoggerListByWorkspaceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &WorkspaceLoggerListByWorkspaceCustomPager{},
		Path:          fmt.Sprintf("%s/loggers", id.ID()),
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
		Values *[]LoggerContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkspaceLoggerListByWorkspaceComplete retrieves all the results into a single object
func (c LoggerClient) WorkspaceLoggerListByWorkspaceComplete(ctx context.Context, id WorkspaceId, options WorkspaceLoggerListByWorkspaceOperationOptions) (WorkspaceLoggerListByWorkspaceCompleteResult, error) {
	return c.WorkspaceLoggerListByWorkspaceCompleteMatchingPredicate(ctx, id, options, LoggerContractOperationPredicate{})
}

// WorkspaceLoggerListByWorkspaceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c LoggerClient) WorkspaceLoggerListByWorkspaceCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, options WorkspaceLoggerListByWorkspaceOperationOptions, predicate LoggerContractOperationPredicate) (result WorkspaceLoggerListByWorkspaceCompleteResult, err error) {
	items := make([]LoggerContract, 0)

	resp, err := c.WorkspaceLoggerListByWorkspace(ctx, id, options)
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

	result = WorkspaceLoggerListByWorkspaceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
