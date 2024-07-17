package managedclusters

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListGuardrailsVersionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]GuardrailsAvailableVersion
}

type ListGuardrailsVersionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []GuardrailsAvailableVersion
}

type ListGuardrailsVersionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListGuardrailsVersionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListGuardrailsVersions ...
func (c ManagedClustersClient) ListGuardrailsVersions(ctx context.Context, id LocationId) (result ListGuardrailsVersionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListGuardrailsVersionsCustomPager{},
		Path:       fmt.Sprintf("%s/guardrailsVersions", id.ID()),
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
		Values *[]GuardrailsAvailableVersion `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListGuardrailsVersionsComplete retrieves all the results into a single object
func (c ManagedClustersClient) ListGuardrailsVersionsComplete(ctx context.Context, id LocationId) (ListGuardrailsVersionsCompleteResult, error) {
	return c.ListGuardrailsVersionsCompleteMatchingPredicate(ctx, id, GuardrailsAvailableVersionOperationPredicate{})
}

// ListGuardrailsVersionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ManagedClustersClient) ListGuardrailsVersionsCompleteMatchingPredicate(ctx context.Context, id LocationId, predicate GuardrailsAvailableVersionOperationPredicate) (result ListGuardrailsVersionsCompleteResult, err error) {
	items := make([]GuardrailsAvailableVersion, 0)

	resp, err := c.ListGuardrailsVersions(ctx, id)
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

	result = ListGuardrailsVersionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
