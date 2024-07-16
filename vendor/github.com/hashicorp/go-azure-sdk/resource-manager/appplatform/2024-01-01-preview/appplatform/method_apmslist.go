package appplatform

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

type ApmsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ApmResource
}

type ApmsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ApmResource
}

type ApmsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ApmsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ApmsList ...
func (c AppPlatformClient) ApmsList(ctx context.Context, id commonids.SpringCloudServiceId) (result ApmsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ApmsListCustomPager{},
		Path:       fmt.Sprintf("%s/apms", id.ID()),
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
		Values *[]ApmResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ApmsListComplete retrieves all the results into a single object
func (c AppPlatformClient) ApmsListComplete(ctx context.Context, id commonids.SpringCloudServiceId) (ApmsListCompleteResult, error) {
	return c.ApmsListCompleteMatchingPredicate(ctx, id, ApmResourceOperationPredicate{})
}

// ApmsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) ApmsListCompleteMatchingPredicate(ctx context.Context, id commonids.SpringCloudServiceId, predicate ApmResourceOperationPredicate) (result ApmsListCompleteResult, err error) {
	items := make([]ApmResource, 0)

	resp, err := c.ApmsList(ctx, id)
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

	result = ApmsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
