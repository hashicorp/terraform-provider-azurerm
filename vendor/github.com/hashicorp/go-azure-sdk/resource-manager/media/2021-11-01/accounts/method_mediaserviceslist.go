package accounts

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

type MediaservicesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]MediaService
}

type MediaservicesListCompleteResult struct {
	Items []MediaService
}

// MediaservicesList ...
func (c AccountsClient) MediaservicesList(ctx context.Context, id commonids.ResourceGroupId) (result MediaservicesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/providers/Microsoft.Media/mediaServices", id.ID()),
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
		Values *[]MediaService `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// MediaservicesListComplete retrieves all the results into a single object
func (c AccountsClient) MediaservicesListComplete(ctx context.Context, id commonids.ResourceGroupId) (MediaservicesListCompleteResult, error) {
	return c.MediaservicesListCompleteMatchingPredicate(ctx, id, MediaServiceOperationPredicate{})
}

// MediaservicesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AccountsClient) MediaservicesListCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate MediaServiceOperationPredicate) (result MediaservicesListCompleteResult, err error) {
	items := make([]MediaService, 0)

	resp, err := c.MediaservicesList(ctx, id)
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

	result = MediaservicesListCompleteResult{
		Items: items,
	}
	return
}
