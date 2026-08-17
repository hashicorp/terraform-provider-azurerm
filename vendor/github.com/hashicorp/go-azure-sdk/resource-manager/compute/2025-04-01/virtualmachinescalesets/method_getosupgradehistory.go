package virtualmachinescalesets

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetOSUpgradeHistoryOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]UpgradeOperationHistoricalStatusInfo
}

type GetOSUpgradeHistoryCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []UpgradeOperationHistoricalStatusInfo
}

type GetOSUpgradeHistoryCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetOSUpgradeHistoryCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetOSUpgradeHistory ...
func (c VirtualMachineScaleSetsClient) GetOSUpgradeHistory(ctx context.Context, id VirtualMachineScaleSetId) (result GetOSUpgradeHistoryOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &GetOSUpgradeHistoryCustomPager{},
		Path:       fmt.Sprintf("%s/osUpgradeHistory", id.ID()),
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
		Values *[]UpgradeOperationHistoricalStatusInfo `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GetOSUpgradeHistoryComplete retrieves all the results into a single object
func (c VirtualMachineScaleSetsClient) GetOSUpgradeHistoryComplete(ctx context.Context, id VirtualMachineScaleSetId) (GetOSUpgradeHistoryCompleteResult, error) {
	return c.GetOSUpgradeHistoryCompleteMatchingPredicate(ctx, id, UpgradeOperationHistoricalStatusInfoOperationPredicate{})
}

// GetOSUpgradeHistoryCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualMachineScaleSetsClient) GetOSUpgradeHistoryCompleteMatchingPredicate(ctx context.Context, id VirtualMachineScaleSetId, predicate UpgradeOperationHistoricalStatusInfoOperationPredicate) (result GetOSUpgradeHistoryCompleteResult, err error) {
	items := make([]UpgradeOperationHistoricalStatusInfo, 0)

	resp, err := c.GetOSUpgradeHistory(ctx, id)
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

	result = GetOSUpgradeHistoryCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
