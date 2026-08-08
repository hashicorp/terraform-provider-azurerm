package autoimportjobs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByAmlFilesystemOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AutoImportJob
}

type ListByAmlFilesystemCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []AutoImportJob
}

type ListByAmlFilesystemCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByAmlFilesystemCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByAmlFilesystem ...
func (c AutoImportJobsClient) ListByAmlFilesystem(ctx context.Context, id AmlFilesystemId) (result ListByAmlFilesystemOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByAmlFilesystemCustomPager{},
		Path:       fmt.Sprintf("%s/autoImportJobs", id.ID()),
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
		Values *[]AutoImportJob `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByAmlFilesystemComplete retrieves all the results into a single object
func (c AutoImportJobsClient) ListByAmlFilesystemComplete(ctx context.Context, id AmlFilesystemId) (ListByAmlFilesystemCompleteResult, error) {
	return c.ListByAmlFilesystemCompleteMatchingPredicate(ctx, id, AutoImportJobOperationPredicate{})
}

// ListByAmlFilesystemCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AutoImportJobsClient) ListByAmlFilesystemCompleteMatchingPredicate(ctx context.Context, id AmlFilesystemId, predicate AutoImportJobOperationPredicate) (result ListByAmlFilesystemCompleteResult, err error) {
	items := make([]AutoImportJob, 0)

	resp, err := c.ListByAmlFilesystem(ctx, id)
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

	result = ListByAmlFilesystemCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
