package profiles

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AFDProfilesListResourceUsageOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Usage
}

type AFDProfilesListResourceUsageCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Usage
}

type AFDProfilesListResourceUsageCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *AFDProfilesListResourceUsageCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// AFDProfilesListResourceUsage ...
func (c ProfilesClient) AFDProfilesListResourceUsage(ctx context.Context, id ProfileId) (result AFDProfilesListResourceUsageOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &AFDProfilesListResourceUsageCustomPager{},
		Path:       fmt.Sprintf("%s/usages", id.ID()),
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
		Values *[]Usage `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// AFDProfilesListResourceUsageComplete retrieves all the results into a single object
func (c ProfilesClient) AFDProfilesListResourceUsageComplete(ctx context.Context, id ProfileId) (AFDProfilesListResourceUsageCompleteResult, error) {
	return c.AFDProfilesListResourceUsageCompleteMatchingPredicate(ctx, id, UsageOperationPredicate{})
}

// AFDProfilesListResourceUsageCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ProfilesClient) AFDProfilesListResourceUsageCompleteMatchingPredicate(ctx context.Context, id ProfileId, predicate UsageOperationPredicate) (result AFDProfilesListResourceUsageCompleteResult, err error) {
	items := make([]Usage, 0)

	resp, err := c.AFDProfilesListResourceUsage(ctx, id)
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

	result = AFDProfilesListResourceUsageCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
