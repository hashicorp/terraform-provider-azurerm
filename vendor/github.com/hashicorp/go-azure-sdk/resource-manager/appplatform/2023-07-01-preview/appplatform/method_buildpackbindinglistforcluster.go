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

type BuildpackBindingListForClusterOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]BuildpackBindingResource
}

type BuildpackBindingListForClusterCompleteResult struct {
	Items []BuildpackBindingResource
}

// BuildpackBindingListForCluster ...
func (c AppPlatformClient) BuildpackBindingListForCluster(ctx context.Context, id SpringId) (result BuildpackBindingListForClusterOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/buildPackBindings", id.ID()),
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
		Values *[]BuildpackBindingResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// BuildpackBindingListForClusterComplete retrieves all the results into a single object
func (c AppPlatformClient) BuildpackBindingListForClusterComplete(ctx context.Context, id SpringId) (BuildpackBindingListForClusterCompleteResult, error) {
	return c.BuildpackBindingListForClusterCompleteMatchingPredicate(ctx, id, BuildpackBindingResourceOperationPredicate{})
}

// BuildpackBindingListForClusterCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) BuildpackBindingListForClusterCompleteMatchingPredicate(ctx context.Context, id SpringId, predicate BuildpackBindingResourceOperationPredicate) (result BuildpackBindingListForClusterCompleteResult, err error) {
	items := make([]BuildpackBindingResource, 0)

	resp, err := c.BuildpackBindingListForCluster(ctx, id)
	if err != nil {
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

	result = BuildpackBindingListForClusterCompleteResult{
		Items: items,
	}
	return
}
