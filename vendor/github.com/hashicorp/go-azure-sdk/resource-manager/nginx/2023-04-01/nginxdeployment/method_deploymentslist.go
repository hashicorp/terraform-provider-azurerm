package nginxdeployment

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

type DeploymentsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]NginxDeployment
}

type DeploymentsListCompleteResult struct {
	Items []NginxDeployment
}

// DeploymentsList ...
func (c NginxDeploymentClient) DeploymentsList(ctx context.Context, id commonids.SubscriptionId) (result DeploymentsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/providers/Nginx.NginxPlus/nginxDeployments", id.ID()),
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
		Values *[]NginxDeployment `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// DeploymentsListComplete retrieves all the results into a single object
func (c NginxDeploymentClient) DeploymentsListComplete(ctx context.Context, id commonids.SubscriptionId) (DeploymentsListCompleteResult, error) {
	return c.DeploymentsListCompleteMatchingPredicate(ctx, id, NginxDeploymentOperationPredicate{})
}

// DeploymentsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c NginxDeploymentClient) DeploymentsListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate NginxDeploymentOperationPredicate) (result DeploymentsListCompleteResult, err error) {
	items := make([]NginxDeployment, 0)

	resp, err := c.DeploymentsList(ctx, id)
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

	result = DeploymentsListCompleteResult{
		Items: items,
	}
	return
}
