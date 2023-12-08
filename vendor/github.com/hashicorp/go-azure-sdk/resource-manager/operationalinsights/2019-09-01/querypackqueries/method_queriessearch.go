package querypackqueries

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueriesSearchOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]LogAnalyticsQueryPackQuery
}

type QueriesSearchCompleteResult struct {
	Items []LogAnalyticsQueryPackQuery
}

type QueriesSearchOperationOptions struct {
	IncludeBody *bool
	Top         *int64
}

func DefaultQueriesSearchOperationOptions() QueriesSearchOperationOptions {
	return QueriesSearchOperationOptions{}
}

func (o QueriesSearchOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o QueriesSearchOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o QueriesSearchOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.IncludeBody != nil {
		out.Append("includeBody", fmt.Sprintf("%v", *o.IncludeBody))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// QueriesSearch ...
func (c QueryPackQueriesClient) QueriesSearch(ctx context.Context, id QueryPackId, input LogAnalyticsQueryPackQuerySearchProperties, options QueriesSearchOperationOptions) (result QueriesSearchOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		Path:          fmt.Sprintf("%s/queries/search", id.ID()),
		OptionsObject: options,
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
		Values *[]LogAnalyticsQueryPackQuery `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// QueriesSearchComplete retrieves all the results into a single object
func (c QueryPackQueriesClient) QueriesSearchComplete(ctx context.Context, id QueryPackId, input LogAnalyticsQueryPackQuerySearchProperties, options QueriesSearchOperationOptions) (QueriesSearchCompleteResult, error) {
	return c.QueriesSearchCompleteMatchingPredicate(ctx, id, input, options, LogAnalyticsQueryPackQueryOperationPredicate{})
}

// QueriesSearchCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c QueryPackQueriesClient) QueriesSearchCompleteMatchingPredicate(ctx context.Context, id QueryPackId, input LogAnalyticsQueryPackQuerySearchProperties, options QueriesSearchOperationOptions, predicate LogAnalyticsQueryPackQueryOperationPredicate) (result QueriesSearchCompleteResult, err error) {
	items := make([]LogAnalyticsQueryPackQuery, 0)

	resp, err := c.QueriesSearch(ctx, id, input, options)
	if err != nil {
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

	result = QueriesSearchCompleteResult{
		Items: items,
	}
	return
}
