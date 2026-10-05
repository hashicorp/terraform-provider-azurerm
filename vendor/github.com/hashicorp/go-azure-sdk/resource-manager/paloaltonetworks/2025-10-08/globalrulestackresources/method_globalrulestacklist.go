package globalrulestackresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GlobalRulestackListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]GlobalRulestackResource
}

type GlobalRulestackListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []GlobalRulestackResource
}

type GlobalRulestackListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GlobalRulestackListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GlobalRulestackList ...
func (c GlobalRulestackResourcesClient) GlobalRulestackList(ctx context.Context) (result GlobalRulestackListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &GlobalRulestackListCustomPager{},
		Path:       "/providers/PaloAltoNetworks.Cloudngfw/globalRulestacks",
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
		Values *[]GlobalRulestackResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GlobalRulestackListComplete retrieves all the results into a single object
func (c GlobalRulestackResourcesClient) GlobalRulestackListComplete(ctx context.Context) (GlobalRulestackListCompleteResult, error) {
	return c.GlobalRulestackListCompleteMatchingPredicate(ctx, GlobalRulestackResourceOperationPredicate{})
}

// GlobalRulestackListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c GlobalRulestackResourcesClient) GlobalRulestackListCompleteMatchingPredicate(ctx context.Context, predicate GlobalRulestackResourceOperationPredicate) (result GlobalRulestackListCompleteResult, err error) {
	items := make([]GlobalRulestackResource, 0)

	resp, err := c.GlobalRulestackList(ctx)
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

	result = GlobalRulestackListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
