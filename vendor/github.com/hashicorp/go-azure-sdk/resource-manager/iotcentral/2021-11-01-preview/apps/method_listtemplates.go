package apps

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

type ListTemplatesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AppTemplate
}

type ListTemplatesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []AppTemplate
}

// ListTemplates ...
func (c AppsClient) ListTemplates(ctx context.Context, id commonids.SubscriptionId) (result ListTemplatesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/providers/Microsoft.IoTCentral/appTemplates", id.ID()),
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
		Values *[]AppTemplate `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListTemplatesComplete retrieves all the results into a single object
func (c AppsClient) ListTemplatesComplete(ctx context.Context, id commonids.SubscriptionId) (ListTemplatesCompleteResult, error) {
	return c.ListTemplatesCompleteMatchingPredicate(ctx, id, AppTemplateOperationPredicate{})
}

// ListTemplatesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppsClient) ListTemplatesCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate AppTemplateOperationPredicate) (result ListTemplatesCompleteResult, err error) {
	items := make([]AppTemplate, 0)

	resp, err := c.ListTemplates(ctx, id)
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

	result = ListTemplatesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
