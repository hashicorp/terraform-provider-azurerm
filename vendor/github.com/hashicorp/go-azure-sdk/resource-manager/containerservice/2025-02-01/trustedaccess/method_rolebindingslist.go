package trustedaccess

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

type RoleBindingsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]TrustedAccessRoleBinding
}

type RoleBindingsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []TrustedAccessRoleBinding
}

type RoleBindingsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *RoleBindingsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// RoleBindingsList ...
func (c TrustedAccessClient) RoleBindingsList(ctx context.Context, id commonids.KubernetesClusterId) (result RoleBindingsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &RoleBindingsListCustomPager{},
		Path:       fmt.Sprintf("%s/trustedAccessRoleBindings", id.ID()),
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
		Values *[]TrustedAccessRoleBinding `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// RoleBindingsListComplete retrieves all the results into a single object
func (c TrustedAccessClient) RoleBindingsListComplete(ctx context.Context, id commonids.KubernetesClusterId) (RoleBindingsListCompleteResult, error) {
	return c.RoleBindingsListCompleteMatchingPredicate(ctx, id, TrustedAccessRoleBindingOperationPredicate{})
}

// RoleBindingsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c TrustedAccessClient) RoleBindingsListCompleteMatchingPredicate(ctx context.Context, id commonids.KubernetesClusterId, predicate TrustedAccessRoleBindingOperationPredicate) (result RoleBindingsListCompleteResult, err error) {
	items := make([]TrustedAccessRoleBinding, 0)

	resp, err := c.RoleBindingsList(ctx, id)
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

	result = RoleBindingsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
