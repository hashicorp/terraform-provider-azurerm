package runbookdraft

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

type ReplaceContentOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]byte
}

// ReplaceContent ...
func (c RunbookDraftClient) ReplaceContent(ctx context.Context, id RunbookId, input []byte) (result ReplaceContentOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "text/powershell",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPut,
		Path:       fmt.Sprintf("%s/draft/content", id.ID()),
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

// ReplaceContentThenPoll performs ReplaceContent then polls until it's completed
func (c RunbookDraftClient) ReplaceContentThenPoll(ctx context.Context, id RunbookId, input []byte) error {
	result, err := c.ReplaceContent(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ReplaceContent: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after ReplaceContent: %+v", err)
	}

	return nil
}
