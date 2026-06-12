package jobs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListDetectorsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Diagnostics
}

type ListDetectorsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Diagnostics
}

type ListDetectorsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListDetectorsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListDetectors ...
func (c JobsClient) ListDetectors(ctx context.Context, id JobId) (result ListDetectorsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListDetectorsCustomPager{},
		Path:       fmt.Sprintf("%s/detectors", id.ID()),
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
		Values *[]Diagnostics `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListDetectorsComplete retrieves all the results into a single object
func (c JobsClient) ListDetectorsComplete(ctx context.Context, id JobId) (ListDetectorsCompleteResult, error) {
	return c.ListDetectorsCompleteMatchingPredicate(ctx, id, DiagnosticsOperationPredicate{})
}

// ListDetectorsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c JobsClient) ListDetectorsCompleteMatchingPredicate(ctx context.Context, id JobId, predicate DiagnosticsOperationPredicate) (result ListDetectorsCompleteResult, err error) {
	items := make([]Diagnostics, 0)

	resp, err := c.ListDetectors(ctx, id)
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

	result = ListDetectorsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
