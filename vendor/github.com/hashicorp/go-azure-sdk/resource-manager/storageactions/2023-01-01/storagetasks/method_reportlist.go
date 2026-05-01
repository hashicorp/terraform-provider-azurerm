package storagetasks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReportListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]StorageTaskReportInstance
}

type ReportListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []StorageTaskReportInstance
}

type ReportListOperationOptions struct {
	Filter      *string
	Maxpagesize *int64
}

func DefaultReportListOperationOptions() ReportListOperationOptions {
	return ReportListOperationOptions{}
}

func (o ReportListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ReportListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ReportListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Maxpagesize != nil {
		out.Append("$maxpagesize", fmt.Sprintf("%v", *o.Maxpagesize))
	}
	return &out
}

type ReportListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ReportListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ReportList ...
func (c StorageTasksClient) ReportList(ctx context.Context, id StorageTaskId, options ReportListOperationOptions) (result ReportListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ReportListCustomPager{},
		Path:          fmt.Sprintf("%s/reports", id.ID()),
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
		Values *[]StorageTaskReportInstance `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ReportListComplete retrieves all the results into a single object
func (c StorageTasksClient) ReportListComplete(ctx context.Context, id StorageTaskId, options ReportListOperationOptions) (ReportListCompleteResult, error) {
	return c.ReportListCompleteMatchingPredicate(ctx, id, options, StorageTaskReportInstanceOperationPredicate{})
}

// ReportListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c StorageTasksClient) ReportListCompleteMatchingPredicate(ctx context.Context, id StorageTaskId, options ReportListOperationOptions, predicate StorageTaskReportInstanceOperationPredicate) (result ReportListCompleteResult, err error) {
	items := make([]StorageTaskReportInstance, 0)

	resp, err := c.ReportList(ctx, id, options)
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

	result = ReportListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
