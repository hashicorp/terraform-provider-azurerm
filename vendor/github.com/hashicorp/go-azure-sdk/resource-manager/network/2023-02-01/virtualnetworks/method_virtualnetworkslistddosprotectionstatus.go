package virtualnetworks

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

type VirtualNetworksListDdosProtectionStatusOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PublicIPDdosProtectionStatusResult
}

type VirtualNetworksListDdosProtectionStatusCompleteResult struct {
	Items []PublicIPDdosProtectionStatusResult
}

type VirtualNetworksListDdosProtectionStatusOperationOptions struct {
	SkipToken *string
	Top       *int64
}

func DefaultVirtualNetworksListDdosProtectionStatusOperationOptions() VirtualNetworksListDdosProtectionStatusOperationOptions {
	return VirtualNetworksListDdosProtectionStatusOperationOptions{}
}

func (o VirtualNetworksListDdosProtectionStatusOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o VirtualNetworksListDdosProtectionStatusOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o VirtualNetworksListDdosProtectionStatusOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.SkipToken != nil {
		out.Append("skipToken", fmt.Sprintf("%v", *o.SkipToken))
	}
	if o.Top != nil {
		out.Append("top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// VirtualNetworksListDdosProtectionStatus ...
func (c VirtualNetworksClient) VirtualNetworksListDdosProtectionStatus(ctx context.Context, id commonids.VirtualNetworkId, options VirtualNetworksListDdosProtectionStatusOperationOptions) (result VirtualNetworksListDdosProtectionStatusOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		Path:          fmt.Sprintf("%s/ddosProtectionStatus", id.ID()),
		OptionsObject: options,
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

	var values struct {
		Values *[]PublicIPDdosProtectionStatusResult `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	result.Poller, err = resourcemanager.PollerFromResponse(resp, c.Client)
	if err != nil {
		return
	}

	return
}

// VirtualNetworksListDdosProtectionStatusThenPoll performs VirtualNetworksListDdosProtectionStatus then polls until it's completed
func (c VirtualNetworksClient) VirtualNetworksListDdosProtectionStatusThenPoll(ctx context.Context, id commonids.VirtualNetworkId, options VirtualNetworksListDdosProtectionStatusOperationOptions) error {
	result, err := c.VirtualNetworksListDdosProtectionStatus(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing VirtualNetworksListDdosProtectionStatus: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after VirtualNetworksListDdosProtectionStatus: %+v", err)
	}

	return nil
}
