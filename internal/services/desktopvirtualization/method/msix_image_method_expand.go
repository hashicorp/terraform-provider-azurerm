// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package method

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/msiximage"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func expandMsixImage(ctx context.Context, metadata sdk.ResourceMetaData, id msiximage.HostPoolId, input msiximage.MSIXImageURI) (result msiximage.ExpandOperationResponse, err error) {
	msixImageClient := metadata.Client.DesktopVirtualization.MsixImagesClient
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &msiximage.ExpandCustomPager{},
		Path:       fmt.Sprintf("%s/expandMsixImage", id.ID()),
	}

	req, err := msixImageClient.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	// `Expand` method in imported `msiximage` package is lack of the codes below, causing `POST` request body to be empty
	if err = req.Marshal(input); err != nil {
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
		Values *[]msiximage.ExpandMsixImage `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

func ExpandCompleteMsixImage(ctx context.Context, metadata sdk.ResourceMetaData, id msiximage.HostPoolId, input msiximage.MSIXImageURI) (msiximage.ExpandCompleteResult, error) {
	return expandCompleteMatchingPredicateMsixImage(ctx, metadata, id, input, msiximage.ExpandMsixImageOperationPredicate{})
}

func expandCompleteMatchingPredicateMsixImage(ctx context.Context, metadata sdk.ResourceMetaData, id msiximage.HostPoolId, input msiximage.MSIXImageURI, predicate msiximage.ExpandMsixImageOperationPredicate) (result msiximage.ExpandCompleteResult, err error) {
	items := make([]msiximage.ExpandMsixImage, 0)

	resp, err := expandMsixImage(ctx, metadata, id, input)
	if err != nil {
		result.LatestHttpResponse = resp.HttpResponse
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

	result = msiximage.ExpandCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
