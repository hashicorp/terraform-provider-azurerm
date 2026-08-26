package jobexecutions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByAgentOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]JobExecution
}

type ListByAgentCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []JobExecution
}

type ListByAgentOperationOptions struct {
	CreateTimeMax *string
	CreateTimeMin *string
	EndTimeMax    *string
	EndTimeMin    *string
	IsActive      *bool
	Skip          *int64
	Top           *int64
}

func DefaultListByAgentOperationOptions() ListByAgentOperationOptions {
	return ListByAgentOperationOptions{}
}

func (o ListByAgentOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByAgentOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListByAgentOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.CreateTimeMax != nil {
		out.Append("createTimeMax", fmt.Sprintf("%v", *o.CreateTimeMax))
	}
	if o.CreateTimeMin != nil {
		out.Append("createTimeMin", fmt.Sprintf("%v", *o.CreateTimeMin))
	}
	if o.EndTimeMax != nil {
		out.Append("endTimeMax", fmt.Sprintf("%v", *o.EndTimeMax))
	}
	if o.EndTimeMin != nil {
		out.Append("endTimeMin", fmt.Sprintf("%v", *o.EndTimeMin))
	}
	if o.IsActive != nil {
		out.Append("isActive", fmt.Sprintf("%v", *o.IsActive))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ListByAgentCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByAgentCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByAgent ...
func (c JobExecutionsClient) ListByAgent(ctx context.Context, id JobAgentId, options ListByAgentOperationOptions) (result ListByAgentOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListByAgentCustomPager{},
		Path:          fmt.Sprintf("%s/executions", id.ID()),
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
		Values *[]JobExecution `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByAgentComplete retrieves all the results into a single object
func (c JobExecutionsClient) ListByAgentComplete(ctx context.Context, id JobAgentId, options ListByAgentOperationOptions) (ListByAgentCompleteResult, error) {
	return c.ListByAgentCompleteMatchingPredicate(ctx, id, options, JobExecutionOperationPredicate{})
}

// ListByAgentCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c JobExecutionsClient) ListByAgentCompleteMatchingPredicate(ctx context.Context, id JobAgentId, options ListByAgentOperationOptions, predicate JobExecutionOperationPredicate) (result ListByAgentCompleteResult, err error) {
	items := make([]JobExecution, 0)

	resp, err := c.ListByAgent(ctx, id, options)
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

	result = ListByAgentCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
