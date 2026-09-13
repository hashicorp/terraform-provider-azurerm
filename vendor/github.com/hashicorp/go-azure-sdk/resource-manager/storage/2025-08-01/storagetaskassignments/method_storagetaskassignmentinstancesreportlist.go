package storagetaskassignments

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageTaskAssignmentInstancesReportListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]StorageTaskReportInstance
}

type StorageTaskAssignmentInstancesReportListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []StorageTaskReportInstance
}

type StorageTaskAssignmentInstancesReportListOperationOptions struct {
	Filter      *string
	Maxpagesize *int64
}

func DefaultStorageTaskAssignmentInstancesReportListOperationOptions() StorageTaskAssignmentInstancesReportListOperationOptions {
	return StorageTaskAssignmentInstancesReportListOperationOptions{}
}

func (o StorageTaskAssignmentInstancesReportListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o StorageTaskAssignmentInstancesReportListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o StorageTaskAssignmentInstancesReportListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Maxpagesize != nil {
		out.Append("$maxpagesize", fmt.Sprintf("%v", *o.Maxpagesize))
	}
	return &out
}

type StorageTaskAssignmentInstancesReportListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *StorageTaskAssignmentInstancesReportListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// StorageTaskAssignmentInstancesReportList ...
func (c StorageTaskAssignmentsClient) StorageTaskAssignmentInstancesReportList(ctx context.Context, id StorageTaskAssignmentId, options StorageTaskAssignmentInstancesReportListOperationOptions) (result StorageTaskAssignmentInstancesReportListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &StorageTaskAssignmentInstancesReportListCustomPager{},
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

// StorageTaskAssignmentInstancesReportListComplete retrieves all the results into a single object
func (c StorageTaskAssignmentsClient) StorageTaskAssignmentInstancesReportListComplete(ctx context.Context, id StorageTaskAssignmentId, options StorageTaskAssignmentInstancesReportListOperationOptions) (StorageTaskAssignmentInstancesReportListCompleteResult, error) {
	return c.StorageTaskAssignmentInstancesReportListCompleteMatchingPredicate(ctx, id, options, StorageTaskReportInstanceOperationPredicate{})
}

// StorageTaskAssignmentInstancesReportListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c StorageTaskAssignmentsClient) StorageTaskAssignmentInstancesReportListCompleteMatchingPredicate(ctx context.Context, id StorageTaskAssignmentId, options StorageTaskAssignmentInstancesReportListOperationOptions, predicate StorageTaskReportInstanceOperationPredicate) (result StorageTaskAssignmentInstancesReportListCompleteResult, err error) {
	items := make([]StorageTaskReportInstance, 0)

	resp, err := c.StorageTaskAssignmentInstancesReportList(ctx, id, options)
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

	result = StorageTaskAssignmentInstancesReportListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
