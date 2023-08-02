package monitors

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListUserRolesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]UserRoleResponse
}

type ListUserRolesCompleteResult struct {
	Items []UserRoleResponse
}

// ListUserRoles ...
func (c MonitorsClient) ListUserRoles(ctx context.Context, id MonitorId, input UserRoleRequest) (result ListUserRolesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/listUserRoles", id.ID()),
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
		Values *[]UserRoleResponse `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListUserRolesComplete retrieves all the results into a single object
func (c MonitorsClient) ListUserRolesComplete(ctx context.Context, id MonitorId, input UserRoleRequest) (ListUserRolesCompleteResult, error) {
	return c.ListUserRolesCompleteMatchingPredicate(ctx, id, input, UserRoleResponseOperationPredicate{})
}

// ListUserRolesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c MonitorsClient) ListUserRolesCompleteMatchingPredicate(ctx context.Context, id MonitorId, input UserRoleRequest, predicate UserRoleResponseOperationPredicate) (result ListUserRolesCompleteResult, err error) {
	items := make([]UserRoleResponse, 0)

	resp, err := c.ListUserRoles(ctx, id, input)
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

	result = ListUserRolesCompleteResult{
		Items: items,
	}
	return
}
