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

type ApplicationAcceleratorsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ApplicationAcceleratorResource
}

type ApplicationAcceleratorsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ApplicationAcceleratorResource
}

type ApplicationAcceleratorsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ApplicationAcceleratorsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ApplicationAcceleratorsList ...
func (c AppPlatformClient) ApplicationAcceleratorsList(ctx context.Context, id commonids.SpringCloudServiceId) (result ApplicationAcceleratorsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ApplicationAcceleratorsListCustomPager{},
		Path:       fmt.Sprintf("%s/applicationAccelerators", id.ID()),
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
		Values *[]ApplicationAcceleratorResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ApplicationAcceleratorsListComplete retrieves all the results into a single object
func (c AppPlatformClient) ApplicationAcceleratorsListComplete(ctx context.Context, id commonids.SpringCloudServiceId) (ApplicationAcceleratorsListCompleteResult, error) {
	return c.ApplicationAcceleratorsListCompleteMatchingPredicate(ctx, id, ApplicationAcceleratorResourceOperationPredicate{})
}

// ApplicationAcceleratorsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) ApplicationAcceleratorsListCompleteMatchingPredicate(ctx context.Context, id commonids.SpringCloudServiceId, predicate ApplicationAcceleratorResourceOperationPredicate) (result ApplicationAcceleratorsListCompleteResult, err error) {
	items := make([]ApplicationAcceleratorResource, 0)

	resp, err := c.ApplicationAcceleratorsList(ctx, id)
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

	result = ApplicationAcceleratorsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
