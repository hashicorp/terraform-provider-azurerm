package packageresource

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByRuntimeEnvironmentOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Package
}

type ListByRuntimeEnvironmentCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Package
}

type ListByRuntimeEnvironmentCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByRuntimeEnvironmentCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByRuntimeEnvironment ...
func (c PackageResourceClient) ListByRuntimeEnvironment(ctx context.Context, id RuntimeEnvironmentId) (result ListByRuntimeEnvironmentOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByRuntimeEnvironmentCustomPager{},
		Path:       fmt.Sprintf("%s/packages", id.ID()),
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
		Values *[]Package `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByRuntimeEnvironmentComplete retrieves all the results into a single object
func (c PackageResourceClient) ListByRuntimeEnvironmentComplete(ctx context.Context, id RuntimeEnvironmentId) (ListByRuntimeEnvironmentCompleteResult, error) {
	return c.ListByRuntimeEnvironmentCompleteMatchingPredicate(ctx, id, PackageOperationPredicate{})
}

// ListByRuntimeEnvironmentCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c PackageResourceClient) ListByRuntimeEnvironmentCompleteMatchingPredicate(ctx context.Context, id RuntimeEnvironmentId, predicate PackageOperationPredicate) (result ListByRuntimeEnvironmentCompleteResult, err error) {
	items := make([]Package, 0)

	resp, err := c.ListByRuntimeEnvironment(ctx, id)
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

	result = ListByRuntimeEnvironmentCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
