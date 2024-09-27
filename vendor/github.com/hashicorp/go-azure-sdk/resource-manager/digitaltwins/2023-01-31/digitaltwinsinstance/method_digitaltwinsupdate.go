package digitaltwinsinstance

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DigitalTwinsUpdateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *DigitalTwinsDescription
}

// DigitalTwinsUpdate ...
func (c DigitalTwinsInstanceClient) DigitalTwinsUpdate(ctx context.Context, id DigitalTwinsInstanceId, input DigitalTwinsPatchDescription) (result DigitalTwinsUpdateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPatch,
		Path:       id.ID(),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	if err = req.Marshal(input); err != nil {
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

// DigitalTwinsUpdateThenPoll performs DigitalTwinsUpdate then polls until it's completed
func (c DigitalTwinsInstanceClient) DigitalTwinsUpdateThenPoll(ctx context.Context, id DigitalTwinsInstanceId, input DigitalTwinsPatchDescription) error {
	result, err := c.DigitalTwinsUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing DigitalTwinsUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after DigitalTwinsUpdate: %+v", err)
	}

	return nil
}
