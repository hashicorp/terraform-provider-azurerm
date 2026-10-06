package workspaces

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

type QlAadAdminsDeleteOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// QlAadAdminsDelete ...
func (c WorkspacesClient) QlAadAdminsDelete(ctx context.Context, id WorkspaceId) (result QlAadAdminsDeleteOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusNoContent,
			http.StatusOK,
		},
		HttpMethod: http.MethodDelete,
		Path:       fmt.Sprintf("%s/sqlAdministrators/activeDirectory", id.ID()),
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

// QlAadAdminsDeleteThenPoll performs QlAadAdminsDelete then polls until it's completed
func (c WorkspacesClient) QlAadAdminsDeleteThenPoll(ctx context.Context, id WorkspaceId) error {
	result, err := c.QlAadAdminsDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing QlAadAdminsDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after QlAadAdminsDelete: %+v", err)
	}

	return nil
}
