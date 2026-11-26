package storageaccounts

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

type StorageTaskAssignmentsInstancesReportListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]StorageTaskReportInstance
}

type StorageTaskAssignmentsInstancesReportListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []StorageTaskReportInstance
}

type StorageTaskAssignmentsInstancesReportListOperationOptions struct {
	Filter      *string
	Maxpagesize *int64
}

func DefaultStorageTaskAssignmentsInstancesReportListOperationOptions() StorageTaskAssignmentsInstancesReportListOperationOptions {
	return StorageTaskAssignmentsInstancesReportListOperationOptions{}
}

func (o StorageTaskAssignmentsInstancesReportListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o StorageTaskAssignmentsInstancesReportListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o StorageTaskAssignmentsInstancesReportListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Maxpagesize != nil {
		out.Append("$maxpagesize", fmt.Sprintf("%v", *o.Maxpagesize))
	}
	return &out
}

type StorageTaskAssignmentsInstancesReportListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *StorageTaskAssignmentsInstancesReportListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// StorageTaskAssignmentsInstancesReportList ...
func (c StorageAccountsClient) StorageTaskAssignmentsInstancesReportList(ctx context.Context, id commonids.StorageAccountId, options StorageTaskAssignmentsInstancesReportListOperationOptions) (result StorageTaskAssignmentsInstancesReportListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &StorageTaskAssignmentsInstancesReportListCustomPager{},
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

// StorageTaskAssignmentsInstancesReportListComplete retrieves all the results into a single object
func (c StorageAccountsClient) StorageTaskAssignmentsInstancesReportListComplete(ctx context.Context, id commonids.StorageAccountId, options StorageTaskAssignmentsInstancesReportListOperationOptions) (StorageTaskAssignmentsInstancesReportListCompleteResult, error) {
	return c.StorageTaskAssignmentsInstancesReportListCompleteMatchingPredicate(ctx, id, options, StorageTaskReportInstanceOperationPredicate{})
}

// StorageTaskAssignmentsInstancesReportListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c StorageAccountsClient) StorageTaskAssignmentsInstancesReportListCompleteMatchingPredicate(ctx context.Context, id commonids.StorageAccountId, options StorageTaskAssignmentsInstancesReportListOperationOptions, predicate StorageTaskReportInstanceOperationPredicate) (result StorageTaskAssignmentsInstancesReportListCompleteResult, err error) {
	items := make([]StorageTaskReportInstance, 0)

	resp, err := c.StorageTaskAssignmentsInstancesReportList(ctx, id, options)
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

	result = StorageTaskAssignmentsInstancesReportListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
