package grafanaresource

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

type GrafanaListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ManagedGrafana
}

type GrafanaListByResourceGroupCompleteResult struct {
	Items []ManagedGrafana
}

// GrafanaListByResourceGroup ...
func (c GrafanaResourceClient) GrafanaListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result GrafanaListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/providers/Microsoft.Dashboard/grafana", id.ID()),
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
		Values *[]ManagedGrafana `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GrafanaListByResourceGroupComplete retrieves all the results into a single object
func (c GrafanaResourceClient) GrafanaListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (GrafanaListByResourceGroupCompleteResult, error) {
	return c.GrafanaListByResourceGroupCompleteMatchingPredicate(ctx, id, ManagedGrafanaOperationPredicate{})
}

// GrafanaListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c GrafanaResourceClient) GrafanaListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate ManagedGrafanaOperationPredicate) (result GrafanaListByResourceGroupCompleteResult, err error) {
	items := make([]ManagedGrafana, 0)

	resp, err := c.GrafanaListByResourceGroup(ctx, id)
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

	result = GrafanaListByResourceGroupCompleteResult{
		Items: items,
	}
	return
}
