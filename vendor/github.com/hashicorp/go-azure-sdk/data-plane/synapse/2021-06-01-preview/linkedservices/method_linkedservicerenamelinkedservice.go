package linkedservices

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinkedServiceRenameLinkedServiceOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// LinkedServiceRenameLinkedService ...
func (c LinkedServicesClient) LinkedServiceRenameLinkedService(ctx context.Context, id LinkedServiceId, input ArtifactRenameRequest) (result LinkedServiceRenameLinkedServiceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/rename", id.Path()),
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

	result.Poller, err = dataplane.PollerFromResponse(resp, c.Client)
	if err != nil {
		return
	}

	return
}

// LinkedServiceRenameLinkedServiceThenPoll performs LinkedServiceRenameLinkedService then polls until it's completed
func (c LinkedServicesClient) LinkedServiceRenameLinkedServiceThenPoll(ctx context.Context, id LinkedServiceId, input ArtifactRenameRequest) error {
	return c.LinkedServiceRenameLinkedServiceCallbackThenPoll(ctx, id, input, nil)
}

// LinkedServiceRenameLinkedServiceCallbackThenPoll performs LinkedServiceRenameLinkedService, runs the optional callback function, then polls until it's completed
func (c LinkedServicesClient) LinkedServiceRenameLinkedServiceCallbackThenPoll(ctx context.Context, id LinkedServiceId, input ArtifactRenameRequest, callback func() error) error {
	result, err := c.LinkedServiceRenameLinkedService(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing LinkedServiceRenameLinkedService: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after LinkedServiceRenameLinkedService: %+v", err)
	}

	return nil
}
