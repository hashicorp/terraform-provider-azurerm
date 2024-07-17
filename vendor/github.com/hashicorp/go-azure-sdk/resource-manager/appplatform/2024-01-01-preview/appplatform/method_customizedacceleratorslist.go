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

type CustomizedAcceleratorsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CustomizedAcceleratorResource
}

type CustomizedAcceleratorsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CustomizedAcceleratorResource
}

type CustomizedAcceleratorsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *CustomizedAcceleratorsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// CustomizedAcceleratorsList ...
func (c AppPlatformClient) CustomizedAcceleratorsList(ctx context.Context, id ApplicationAcceleratorId) (result CustomizedAcceleratorsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &CustomizedAcceleratorsListCustomPager{},
		Path:       fmt.Sprintf("%s/customizedAccelerators", id.ID()),
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
		Values *[]CustomizedAcceleratorResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// CustomizedAcceleratorsListComplete retrieves all the results into a single object
func (c AppPlatformClient) CustomizedAcceleratorsListComplete(ctx context.Context, id ApplicationAcceleratorId) (CustomizedAcceleratorsListCompleteResult, error) {
	return c.CustomizedAcceleratorsListCompleteMatchingPredicate(ctx, id, CustomizedAcceleratorResourceOperationPredicate{})
}

// CustomizedAcceleratorsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) CustomizedAcceleratorsListCompleteMatchingPredicate(ctx context.Context, id ApplicationAcceleratorId, predicate CustomizedAcceleratorResourceOperationPredicate) (result CustomizedAcceleratorsListCompleteResult, err error) {
	items := make([]CustomizedAcceleratorResource, 0)

	resp, err := c.CustomizedAcceleratorsList(ctx, id)
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

	result = CustomizedAcceleratorsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
