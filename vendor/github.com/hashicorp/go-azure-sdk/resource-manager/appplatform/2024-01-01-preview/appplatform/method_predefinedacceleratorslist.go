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

type PredefinedAcceleratorsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PredefinedAcceleratorResource
}

type PredefinedAcceleratorsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PredefinedAcceleratorResource
}

type PredefinedAcceleratorsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *PredefinedAcceleratorsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// PredefinedAcceleratorsList ...
func (c AppPlatformClient) PredefinedAcceleratorsList(ctx context.Context, id ApplicationAcceleratorId) (result PredefinedAcceleratorsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &PredefinedAcceleratorsListCustomPager{},
		Path:       fmt.Sprintf("%s/predefinedAccelerators", id.ID()),
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
		Values *[]PredefinedAcceleratorResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// PredefinedAcceleratorsListComplete retrieves all the results into a single object
func (c AppPlatformClient) PredefinedAcceleratorsListComplete(ctx context.Context, id ApplicationAcceleratorId) (PredefinedAcceleratorsListCompleteResult, error) {
	return c.PredefinedAcceleratorsListCompleteMatchingPredicate(ctx, id, PredefinedAcceleratorResourceOperationPredicate{})
}

// PredefinedAcceleratorsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) PredefinedAcceleratorsListCompleteMatchingPredicate(ctx context.Context, id ApplicationAcceleratorId, predicate PredefinedAcceleratorResourceOperationPredicate) (result PredefinedAcceleratorsListCompleteResult, err error) {
	items := make([]PredefinedAcceleratorResource, 0)

	resp, err := c.PredefinedAcceleratorsList(ctx, id)
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

	result = PredefinedAcceleratorsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
