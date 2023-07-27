package applicationgatewaywafdynamicmanifests

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayWafDynamicManifestsGetOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ApplicationGatewayWafDynamicManifestResult
}

type ApplicationGatewayWafDynamicManifestsGetCompleteResult struct {
	Items []ApplicationGatewayWafDynamicManifestResult
}

// ApplicationGatewayWafDynamicManifestsGet ...
func (c ApplicationGatewayWafDynamicManifestsClient) ApplicationGatewayWafDynamicManifestsGet(ctx context.Context, id LocationId) (result ApplicationGatewayWafDynamicManifestsGetOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/applicationGatewayWafDynamicManifests", id.ID()),
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
		Values *[]ApplicationGatewayWafDynamicManifestResult `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ApplicationGatewayWafDynamicManifestsGetComplete retrieves all the results into a single object
func (c ApplicationGatewayWafDynamicManifestsClient) ApplicationGatewayWafDynamicManifestsGetComplete(ctx context.Context, id LocationId) (ApplicationGatewayWafDynamicManifestsGetCompleteResult, error) {
	return c.ApplicationGatewayWafDynamicManifestsGetCompleteMatchingPredicate(ctx, id, ApplicationGatewayWafDynamicManifestResultOperationPredicate{})
}

// ApplicationGatewayWafDynamicManifestsGetCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ApplicationGatewayWafDynamicManifestsClient) ApplicationGatewayWafDynamicManifestsGetCompleteMatchingPredicate(ctx context.Context, id LocationId, predicate ApplicationGatewayWafDynamicManifestResultOperationPredicate) (result ApplicationGatewayWafDynamicManifestsGetCompleteResult, err error) {
	items := make([]ApplicationGatewayWafDynamicManifestResult, 0)

	resp, err := c.ApplicationGatewayWafDynamicManifestsGet(ctx, id)
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

	result = ApplicationGatewayWafDynamicManifestsGetCompleteResult{
		Items: items,
	}
	return
}
