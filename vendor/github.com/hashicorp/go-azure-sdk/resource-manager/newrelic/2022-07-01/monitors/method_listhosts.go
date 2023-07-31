package monitors

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListHostsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VMInfo
}

type ListHostsCompleteResult struct {
	Items []VMInfo
}

// ListHosts ...
func (c MonitorsClient) ListHosts(ctx context.Context, id MonitorId, input HostsGetRequest) (result ListHostsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/listHosts", id.ID()),
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
		Values *[]VMInfo `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListHostsComplete retrieves all the results into a single object
func (c MonitorsClient) ListHostsComplete(ctx context.Context, id MonitorId, input HostsGetRequest) (ListHostsCompleteResult, error) {
	return c.ListHostsCompleteMatchingPredicate(ctx, id, input, VMInfoOperationPredicate{})
}

// ListHostsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c MonitorsClient) ListHostsCompleteMatchingPredicate(ctx context.Context, id MonitorId, input HostsGetRequest, predicate VMInfoOperationPredicate) (result ListHostsCompleteResult, err error) {
	items := make([]VMInfo, 0)

	resp, err := c.ListHosts(ctx, id, input)
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

	result = ListHostsCompleteResult{
		Items: items,
	}
	return
}
