package digitaltwinsinstance

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

type DigitalTwinsListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DigitalTwinsDescription
}

type DigitalTwinsListByResourceGroupCompleteResult struct {
	Items []DigitalTwinsDescription
}

// DigitalTwinsListByResourceGroup ...
func (c DigitalTwinsInstanceClient) DigitalTwinsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result DigitalTwinsListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/providers/Microsoft.DigitalTwins/digitalTwinsInstances", id.ID()),
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
		Values *[]DigitalTwinsDescription `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// DigitalTwinsListByResourceGroupComplete retrieves all the results into a single object
func (c DigitalTwinsInstanceClient) DigitalTwinsListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (DigitalTwinsListByResourceGroupCompleteResult, error) {
	return c.DigitalTwinsListByResourceGroupCompleteMatchingPredicate(ctx, id, DigitalTwinsDescriptionOperationPredicate{})
}

// DigitalTwinsListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DigitalTwinsInstanceClient) DigitalTwinsListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate DigitalTwinsDescriptionOperationPredicate) (result DigitalTwinsListByResourceGroupCompleteResult, err error) {
	items := make([]DigitalTwinsDescription, 0)

	resp, err := c.DigitalTwinsListByResourceGroup(ctx, id)
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

	result = DigitalTwinsListByResourceGroupCompleteResult{
		Items: items,
	}
	return
}
