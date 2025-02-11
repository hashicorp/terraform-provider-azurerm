package backupvaults

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

type GetInSubscriptionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]BackupVaultResource
}

type GetInSubscriptionCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []BackupVaultResource
}

type GetInSubscriptionCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetInSubscriptionCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetInSubscription ...
func (c BackupVaultsClient) GetInSubscription(ctx context.Context, id commonids.SubscriptionId) (result GetInSubscriptionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &GetInSubscriptionCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.DataProtection/backupVaults", id.ID()),
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
		Values *[]BackupVaultResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GetInSubscriptionComplete retrieves all the results into a single object
func (c BackupVaultsClient) GetInSubscriptionComplete(ctx context.Context, id commonids.SubscriptionId) (GetInSubscriptionCompleteResult, error) {
	return c.GetInSubscriptionCompleteMatchingPredicate(ctx, id, BackupVaultResourceOperationPredicate{})
}

// GetInSubscriptionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c BackupVaultsClient) GetInSubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate BackupVaultResourceOperationPredicate) (result GetInSubscriptionCompleteResult, err error) {
	items := make([]BackupVaultResource, 0)

	resp, err := c.GetInSubscription(ctx, id)
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

	result = GetInSubscriptionCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
