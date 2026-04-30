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

type DefaultWafPolicyListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]NginxDeploymentDefaultWafPolicyProperties
}

type DefaultWafPolicyListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []NginxDeploymentDefaultWafPolicyProperties
}

type DefaultWafPolicyListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *DefaultWafPolicyListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// DefaultWafPolicyList ...
func (c NginxDeploymentsClient) DefaultWafPolicyList(ctx context.Context, id NginxDeploymentId) (result DefaultWafPolicyListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &DefaultWafPolicyListCustomPager{},
		Path:       fmt.Sprintf("%s/listDefaultWafPolicies", id.ID()),
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
		Values *[]NginxDeploymentDefaultWafPolicyProperties `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// DefaultWafPolicyListComplete retrieves all the results into a single object
func (c NginxDeploymentsClient) DefaultWafPolicyListComplete(ctx context.Context, id NginxDeploymentId) (DefaultWafPolicyListCompleteResult, error) {
	return c.DefaultWafPolicyListCompleteMatchingPredicate(ctx, id, NginxDeploymentDefaultWafPolicyPropertiesOperationPredicate{})
}

// DefaultWafPolicyListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c NginxDeploymentsClient) DefaultWafPolicyListCompleteMatchingPredicate(ctx context.Context, id NginxDeploymentId, predicate NginxDeploymentDefaultWafPolicyPropertiesOperationPredicate) (result DefaultWafPolicyListCompleteResult, err error) {
	items := make([]NginxDeploymentDefaultWafPolicyProperties, 0)

	resp, err := c.DefaultWafPolicyList(ctx, id)
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

	result = DefaultWafPolicyListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
