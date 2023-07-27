package virtualwans

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationPolicyGroupsListByVpnServerConfigurationOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VpnServerConfigurationPolicyGroup
}

type ConfigurationPolicyGroupsListByVpnServerConfigurationCompleteResult struct {
	Items []VpnServerConfigurationPolicyGroup
}

// ConfigurationPolicyGroupsListByVpnServerConfiguration ...
func (c VirtualWANsClient) ConfigurationPolicyGroupsListByVpnServerConfiguration(ctx context.Context, id VpnServerConfigurationId) (result ConfigurationPolicyGroupsListByVpnServerConfigurationOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/configurationPolicyGroups", id.ID()),
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
		Values *[]VpnServerConfigurationPolicyGroup `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ConfigurationPolicyGroupsListByVpnServerConfigurationComplete retrieves all the results into a single object
func (c VirtualWANsClient) ConfigurationPolicyGroupsListByVpnServerConfigurationComplete(ctx context.Context, id VpnServerConfigurationId) (ConfigurationPolicyGroupsListByVpnServerConfigurationCompleteResult, error) {
	return c.ConfigurationPolicyGroupsListByVpnServerConfigurationCompleteMatchingPredicate(ctx, id, VpnServerConfigurationPolicyGroupOperationPredicate{})
}

// ConfigurationPolicyGroupsListByVpnServerConfigurationCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualWANsClient) ConfigurationPolicyGroupsListByVpnServerConfigurationCompleteMatchingPredicate(ctx context.Context, id VpnServerConfigurationId, predicate VpnServerConfigurationPolicyGroupOperationPredicate) (result ConfigurationPolicyGroupsListByVpnServerConfigurationCompleteResult, err error) {
	items := make([]VpnServerConfigurationPolicyGroup, 0)

	resp, err := c.ConfigurationPolicyGroupsListByVpnServerConfiguration(ctx, id)
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

	result = ConfigurationPolicyGroupsListByVpnServerConfigurationCompleteResult{
		Items: items,
	}
	return
}
