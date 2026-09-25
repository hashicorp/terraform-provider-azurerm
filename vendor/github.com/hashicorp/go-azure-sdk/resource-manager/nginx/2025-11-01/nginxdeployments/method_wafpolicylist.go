package nginxdeployments

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WafPolicyListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]NginxDeploymentWafPolicyMetadata
}

type WafPolicyListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []NginxDeploymentWafPolicyMetadata
}

type WafPolicyListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *WafPolicyListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// WafPolicyList ...
func (c NginxDeploymentsClient) WafPolicyList(ctx context.Context, id NginxDeploymentId) (result WafPolicyListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &WafPolicyListCustomPager{},
		Path:       fmt.Sprintf("%s/wafPolicies", id.ID()),
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
		Values *[]NginxDeploymentWafPolicyMetadata `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WafPolicyListComplete retrieves all the results into a single object
func (c NginxDeploymentsClient) WafPolicyListComplete(ctx context.Context, id NginxDeploymentId) (WafPolicyListCompleteResult, error) {
	return c.WafPolicyListCompleteMatchingPredicate(ctx, id, NginxDeploymentWafPolicyMetadataOperationPredicate{})
}

// WafPolicyListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c NginxDeploymentsClient) WafPolicyListCompleteMatchingPredicate(ctx context.Context, id NginxDeploymentId, predicate NginxDeploymentWafPolicyMetadataOperationPredicate) (result WafPolicyListCompleteResult, err error) {
	items := make([]NginxDeploymentWafPolicyMetadata, 0)

	resp, err := c.WafPolicyList(ctx, id)
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

	result = WafPolicyListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
