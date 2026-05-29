package appplatform

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

type CustomizedAcceleratorsDeleteOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// CustomizedAcceleratorsDelete ...
func (c AppPlatformClient) CustomizedAcceleratorsDelete(ctx context.Context, id CustomizedAcceleratorId) (result CustomizedAcceleratorsDeleteOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusNoContent,
		},
		HttpMethod: http.MethodDelete,
		Path:       id.ID(),
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

// CustomizedAcceleratorsDeleteThenPoll performs CustomizedAcceleratorsDelete then polls until it's completed
func (c AppPlatformClient) CustomizedAcceleratorsDeleteThenPoll(ctx context.Context, id CustomizedAcceleratorId) error {
	result, err := c.CustomizedAcceleratorsDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing CustomizedAcceleratorsDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after CustomizedAcceleratorsDelete: %+v", err)
	}

	return nil
}
