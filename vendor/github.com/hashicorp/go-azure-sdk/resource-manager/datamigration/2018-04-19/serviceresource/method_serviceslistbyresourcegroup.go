package serviceresource

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServicesListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DataMigrationService
}

type ServicesListByResourceGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DataMigrationService
}

// ServicesListByResourceGroup ...
func (c ServiceResourceClient) ServicesListByResourceGroup(ctx context.Context, id ResourceGroupId) (result ServicesListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/providers/Microsoft.DataMigration/services", id.ID()),
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
		Values *[]DataMigrationService `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ServicesListByResourceGroupComplete retrieves all the results into a single object
func (c ServiceResourceClient) ServicesListByResourceGroupComplete(ctx context.Context, id ResourceGroupId) (ServicesListByResourceGroupCompleteResult, error) {
	return c.ServicesListByResourceGroupCompleteMatchingPredicate(ctx, id, DataMigrationServiceOperationPredicate{})
}

// ServicesListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ServiceResourceClient) ServicesListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id ResourceGroupId, predicate DataMigrationServiceOperationPredicate) (result ServicesListByResourceGroupCompleteResult, err error) {
	items := make([]DataMigrationService, 0)

	resp, err := c.ServicesListByResourceGroup(ctx, id)
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

	result = ServicesListByResourceGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
