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

type StorageTaskAssignmentListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]StorageTaskAssignment
}

type StorageTaskAssignmentListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []StorageTaskAssignment
}

type StorageTaskAssignmentListOperationOptions struct {
	Maxpagesize *int64
}

func DefaultStorageTaskAssignmentListOperationOptions() StorageTaskAssignmentListOperationOptions {
	return StorageTaskAssignmentListOperationOptions{}
}

func (o StorageTaskAssignmentListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o StorageTaskAssignmentListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o StorageTaskAssignmentListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Maxpagesize != nil {
		out.Append("$maxpagesize", fmt.Sprintf("%v", *o.Maxpagesize))
	}
	return &out
}

type StorageTaskAssignmentListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *StorageTaskAssignmentListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// StorageTaskAssignmentList ...
func (c StorageTasksClient) StorageTaskAssignmentList(ctx context.Context, id StorageTaskId, options StorageTaskAssignmentListOperationOptions) (result StorageTaskAssignmentListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &StorageTaskAssignmentListCustomPager{},
		Path:          fmt.Sprintf("%s/storageTaskAssignments", id.ID()),
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
		Values *[]StorageTaskAssignment `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// StorageTaskAssignmentListComplete retrieves all the results into a single object
func (c StorageTasksClient) StorageTaskAssignmentListComplete(ctx context.Context, id StorageTaskId, options StorageTaskAssignmentListOperationOptions) (StorageTaskAssignmentListCompleteResult, error) {
	return c.StorageTaskAssignmentListCompleteMatchingPredicate(ctx, id, options, StorageTaskAssignmentOperationPredicate{})
}

// StorageTaskAssignmentListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c StorageTasksClient) StorageTaskAssignmentListCompleteMatchingPredicate(ctx context.Context, id StorageTaskId, options StorageTaskAssignmentListOperationOptions, predicate StorageTaskAssignmentOperationPredicate) (result StorageTaskAssignmentListCompleteResult, err error) {
	items := make([]StorageTaskAssignment, 0)

	resp, err := c.StorageTaskAssignmentList(ctx, id, options)
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

	result = StorageTaskAssignmentListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
