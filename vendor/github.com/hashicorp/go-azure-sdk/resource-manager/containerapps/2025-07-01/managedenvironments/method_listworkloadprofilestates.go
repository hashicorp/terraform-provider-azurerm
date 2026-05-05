package managedenvironments

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListWorkloadProfileStatesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]WorkloadProfileStates
}

type ListWorkloadProfileStatesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []WorkloadProfileStates
}

type ListWorkloadProfileStatesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListWorkloadProfileStatesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListWorkloadProfileStates ...
func (c ManagedEnvironmentsClient) ListWorkloadProfileStates(ctx context.Context, id ManagedEnvironmentId) (result ListWorkloadProfileStatesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListWorkloadProfileStatesCustomPager{},
		Path:       fmt.Sprintf("%s/workloadProfileStates", id.ID()),
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
		Values *[]WorkloadProfileStates `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListWorkloadProfileStatesComplete retrieves all the results into a single object
func (c ManagedEnvironmentsClient) ListWorkloadProfileStatesComplete(ctx context.Context, id ManagedEnvironmentId) (ListWorkloadProfileStatesCompleteResult, error) {
	return c.ListWorkloadProfileStatesCompleteMatchingPredicate(ctx, id, WorkloadProfileStatesOperationPredicate{})
}

// ListWorkloadProfileStatesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ManagedEnvironmentsClient) ListWorkloadProfileStatesCompleteMatchingPredicate(ctx context.Context, id ManagedEnvironmentId, predicate WorkloadProfileStatesOperationPredicate) (result ListWorkloadProfileStatesCompleteResult, err error) {
	items := make([]WorkloadProfileStates, 0)

	resp, err := c.ListWorkloadProfileStates(ctx, id)
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

	result = ListWorkloadProfileStatesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
