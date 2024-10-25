package storagetaskassignments

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

type InstancesReportListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]StorageTaskReportInstance
}

type InstancesReportListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []StorageTaskReportInstance
}

type InstancesReportListOperationOptions struct {
	Filter      *string
	Maxpagesize *int64
}

func DefaultInstancesReportListOperationOptions() InstancesReportListOperationOptions {
	return InstancesReportListOperationOptions{}
}

func (o InstancesReportListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o InstancesReportListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o InstancesReportListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Maxpagesize != nil {
		out.Append("$maxpagesize", fmt.Sprintf("%v", *o.Maxpagesize))
	}
	return &out
}

type InstancesReportListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *InstancesReportListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// InstancesReportList ...
func (c StorageTaskAssignmentsClient) InstancesReportList(ctx context.Context, id commonids.StorageAccountId, options InstancesReportListOperationOptions) (result InstancesReportListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &InstancesReportListCustomPager{},
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

// InstancesReportListComplete retrieves all the results into a single object
func (c StorageTaskAssignmentsClient) InstancesReportListComplete(ctx context.Context, id commonids.StorageAccountId, options InstancesReportListOperationOptions) (InstancesReportListCompleteResult, error) {
	return c.InstancesReportListCompleteMatchingPredicate(ctx, id, options, StorageTaskReportInstanceOperationPredicate{})
}

// InstancesReportListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c StorageTaskAssignmentsClient) InstancesReportListCompleteMatchingPredicate(ctx context.Context, id commonids.StorageAccountId, options InstancesReportListOperationOptions, predicate StorageTaskReportInstanceOperationPredicate) (result InstancesReportListCompleteResult, err error) {
	items := make([]StorageTaskReportInstance, 0)

	resp, err := c.InstancesReportList(ctx, id, options)
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

	result = InstancesReportListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
