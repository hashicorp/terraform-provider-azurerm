package appplatform

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BuildServiceListSupportedBuildpacksOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SupportedBuildpackResource
}

type BuildServiceListSupportedBuildpacksCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SupportedBuildpackResource
}

type BuildServiceListSupportedBuildpacksCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *BuildServiceListSupportedBuildpacksCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// BuildServiceListSupportedBuildpacks ...
func (c AppPlatformClient) BuildServiceListSupportedBuildpacks(ctx context.Context, id BuildServiceId) (result BuildServiceListSupportedBuildpacksOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &BuildServiceListSupportedBuildpacksCustomPager{},
		Path:       fmt.Sprintf("%s/supportedBuildPacks", id.ID()),
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
		Values *[]SupportedBuildpackResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// BuildServiceListSupportedBuildpacksComplete retrieves all the results into a single object
func (c AppPlatformClient) BuildServiceListSupportedBuildpacksComplete(ctx context.Context, id BuildServiceId) (BuildServiceListSupportedBuildpacksCompleteResult, error) {
	return c.BuildServiceListSupportedBuildpacksCompleteMatchingPredicate(ctx, id, SupportedBuildpackResourceOperationPredicate{})
}

// BuildServiceListSupportedBuildpacksCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) BuildServiceListSupportedBuildpacksCompleteMatchingPredicate(ctx context.Context, id BuildServiceId, predicate SupportedBuildpackResourceOperationPredicate) (result BuildServiceListSupportedBuildpacksCompleteResult, err error) {
	items := make([]SupportedBuildpackResource, 0)

	resp, err := c.BuildServiceListSupportedBuildpacks(ctx, id)
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

	result = BuildServiceListSupportedBuildpacksCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
