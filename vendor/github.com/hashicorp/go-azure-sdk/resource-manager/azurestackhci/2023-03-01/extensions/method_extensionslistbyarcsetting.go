package extensions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExtensionsListByArcSettingOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Extension
}

type ExtensionsListByArcSettingCompleteResult struct {
	Items []Extension
}

// ExtensionsListByArcSetting ...
func (c ExtensionsClient) ExtensionsListByArcSetting(ctx context.Context, id ArcSettingId) (result ExtensionsListByArcSettingOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/extensions", id.ID()),
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
		Values *[]Extension `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ExtensionsListByArcSettingComplete retrieves all the results into a single object
func (c ExtensionsClient) ExtensionsListByArcSettingComplete(ctx context.Context, id ArcSettingId) (ExtensionsListByArcSettingCompleteResult, error) {
	return c.ExtensionsListByArcSettingCompleteMatchingPredicate(ctx, id, ExtensionOperationPredicate{})
}

// ExtensionsListByArcSettingCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ExtensionsClient) ExtensionsListByArcSettingCompleteMatchingPredicate(ctx context.Context, id ArcSettingId, predicate ExtensionOperationPredicate) (result ExtensionsListByArcSettingCompleteResult, err error) {
	items := make([]Extension, 0)

	resp, err := c.ExtensionsListByArcSetting(ctx, id)
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

	result = ExtensionsListByArcSettingCompleteResult{
		Items: items,
	}
	return
}
