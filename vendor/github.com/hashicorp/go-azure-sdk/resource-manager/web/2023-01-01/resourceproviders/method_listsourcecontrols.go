package resourceproviders

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListSourceControlsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SourceControl
}

type ListSourceControlsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SourceControl
}

type ListSourceControlsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListSourceControlsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListSourceControls ...
func (c ResourceProvidersClient) ListSourceControls(ctx context.Context) (result ListSourceControlsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListSourceControlsCustomPager{},
		Path:       "/providers/Microsoft.Web/sourceControls",
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
		Values *[]SourceControl `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListSourceControlsComplete retrieves all the results into a single object
func (c ResourceProvidersClient) ListSourceControlsComplete(ctx context.Context) (ListSourceControlsCompleteResult, error) {
	return c.ListSourceControlsCompleteMatchingPredicate(ctx, SourceControlOperationPredicate{})
}

// ListSourceControlsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ResourceProvidersClient) ListSourceControlsCompleteMatchingPredicate(ctx context.Context, predicate SourceControlOperationPredicate) (result ListSourceControlsCompleteResult, err error) {
	items := make([]SourceControl, 0)

	resp, err := c.ListSourceControls(ctx)
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

	result = ListSourceControlsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
