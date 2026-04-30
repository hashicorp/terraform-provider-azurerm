package nginxdeploymentapikeyresponses

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiKeysListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]NginxDeploymentApiKeyResponse
}

type ApiKeysListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []NginxDeploymentApiKeyResponse
}

type ApiKeysListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ApiKeysListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ApiKeysList ...
func (c NginxDeploymentApiKeyResponsesClient) ApiKeysList(ctx context.Context, id NginxDeploymentId) (result ApiKeysListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ApiKeysListCustomPager{},
		Path:       fmt.Sprintf("%s/apiKeys", id.ID()),
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
		Values *[]NginxDeploymentApiKeyResponse `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ApiKeysListComplete retrieves all the results into a single object
func (c NginxDeploymentApiKeyResponsesClient) ApiKeysListComplete(ctx context.Context, id NginxDeploymentId) (ApiKeysListCompleteResult, error) {
	return c.ApiKeysListCompleteMatchingPredicate(ctx, id, NginxDeploymentApiKeyResponseOperationPredicate{})
}

// ApiKeysListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c NginxDeploymentApiKeyResponsesClient) ApiKeysListCompleteMatchingPredicate(ctx context.Context, id NginxDeploymentId, predicate NginxDeploymentApiKeyResponseOperationPredicate) (result ApiKeysListCompleteResult, err error) {
	items := make([]NginxDeploymentApiKeyResponse, 0)

	resp, err := c.ApiKeysList(ctx, id)
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

	result = ApiKeysListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
