package projectresource

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProjectsListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Project
}

type ProjectsListByResourceGroupCompleteResult struct {
	Items []Project
}

// ProjectsListByResourceGroup ...
func (c ProjectResourceClient) ProjectsListByResourceGroup(ctx context.Context, id ServiceId) (result ProjectsListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/projects", id.ID()),
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
		Values *[]Project `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ProjectsListByResourceGroupComplete retrieves all the results into a single object
func (c ProjectResourceClient) ProjectsListByResourceGroupComplete(ctx context.Context, id ServiceId) (ProjectsListByResourceGroupCompleteResult, error) {
	return c.ProjectsListByResourceGroupCompleteMatchingPredicate(ctx, id, ProjectOperationPredicate{})
}

// ProjectsListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ProjectResourceClient) ProjectsListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id ServiceId, predicate ProjectOperationPredicate) (result ProjectsListByResourceGroupCompleteResult, err error) {
	items := make([]Project, 0)

	resp, err := c.ProjectsListByResourceGroup(ctx, id)
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

	result = ProjectsListByResourceGroupCompleteResult{
		Items: items,
	}
	return
}
