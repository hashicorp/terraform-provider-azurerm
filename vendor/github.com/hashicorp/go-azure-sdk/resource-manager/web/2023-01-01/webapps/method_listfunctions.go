package webapps

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

type ListFunctionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]FunctionEnvelope
}

type ListFunctionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []FunctionEnvelope
}

type ListFunctionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListFunctionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListFunctions ...
func (c WebAppsClient) ListFunctions(ctx context.Context, id commonids.AppServiceId) (result ListFunctionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListFunctionsCustomPager{},
		Path:       fmt.Sprintf("%s/functions", id.ID()),
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
		Values *[]FunctionEnvelope `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListFunctionsComplete retrieves all the results into a single object
func (c WebAppsClient) ListFunctionsComplete(ctx context.Context, id commonids.AppServiceId) (ListFunctionsCompleteResult, error) {
	return c.ListFunctionsCompleteMatchingPredicate(ctx, id, FunctionEnvelopeOperationPredicate{})
}

// ListFunctionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListFunctionsCompleteMatchingPredicate(ctx context.Context, id commonids.AppServiceId, predicate FunctionEnvelopeOperationPredicate) (result ListFunctionsCompleteResult, err error) {
	items := make([]FunctionEnvelope, 0)

	resp, err := c.ListFunctions(ctx, id)
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

	result = ListFunctionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
