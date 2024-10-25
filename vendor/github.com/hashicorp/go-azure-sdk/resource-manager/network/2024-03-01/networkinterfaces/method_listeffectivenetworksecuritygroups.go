package networkinterfaces

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListEffectiveNetworkSecurityGroupsOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]EffectiveNetworkSecurityGroup
}

type ListEffectiveNetworkSecurityGroupsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []EffectiveNetworkSecurityGroup
}

type ListEffectiveNetworkSecurityGroupsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListEffectiveNetworkSecurityGroupsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListEffectiveNetworkSecurityGroups ...
func (c NetworkInterfacesClient) ListEffectiveNetworkSecurityGroups(ctx context.Context, id commonids.NetworkInterfaceId) (result ListEffectiveNetworkSecurityGroupsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &ListEffectiveNetworkSecurityGroupsCustomPager{},
		Path:       fmt.Sprintf("%s/effectiveNetworkSecurityGroups", id.ID()),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	result.Poller, err = resourcemanager.PollerFromResponse(resp, c.Client)
	if err != nil {
		return
	}

	return
}

// ListEffectiveNetworkSecurityGroupsThenPoll performs ListEffectiveNetworkSecurityGroups then polls until it's completed
func (c NetworkInterfacesClient) ListEffectiveNetworkSecurityGroupsThenPoll(ctx context.Context, id commonids.NetworkInterfaceId) error {
	result, err := c.ListEffectiveNetworkSecurityGroups(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ListEffectiveNetworkSecurityGroups: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after ListEffectiveNetworkSecurityGroups: %+v", err)
	}

	return nil
}
