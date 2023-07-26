package workspaces

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceFeaturesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AmlUserFeature
}

type WorkspaceFeaturesListCompleteResult struct {
	Items []AmlUserFeature
}

// WorkspaceFeaturesList ...
func (c WorkspacesClient) WorkspaceFeaturesList(ctx context.Context, id WorkspaceId) (result WorkspaceFeaturesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/features", id.ID()),
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
		Values *[]AmlUserFeature `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkspaceFeaturesListComplete retrieves all the results into a single object
func (c WorkspacesClient) WorkspaceFeaturesListComplete(ctx context.Context, id WorkspaceId) (WorkspaceFeaturesListCompleteResult, error) {
	return c.WorkspaceFeaturesListCompleteMatchingPredicate(ctx, id, AmlUserFeatureOperationPredicate{})
}

// WorkspaceFeaturesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WorkspacesClient) WorkspaceFeaturesListCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, predicate AmlUserFeatureOperationPredicate) (result WorkspaceFeaturesListCompleteResult, err error) {
	items := make([]AmlUserFeature, 0)

	resp, err := c.WorkspaceFeaturesList(ctx, id)
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

	result = WorkspaceFeaturesListCompleteResult{
		Items: items,
	}
	return
}
