package manageddatabases

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

type ReevaluateInaccessibleDatabaseStateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *ManagedDatabase
}

// ReevaluateInaccessibleDatabaseState ...
func (c ManagedDatabasesClient) ReevaluateInaccessibleDatabaseState(ctx context.Context, id commonids.SqlManagedInstanceDatabaseId) (result ReevaluateInaccessibleDatabaseStateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/reevaluateInaccessibleDatabaseState", id.ID()),
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

// ReevaluateInaccessibleDatabaseStateThenPoll performs ReevaluateInaccessibleDatabaseState then polls until it's completed
func (c ManagedDatabasesClient) ReevaluateInaccessibleDatabaseStateThenPoll(ctx context.Context, id commonids.SqlManagedInstanceDatabaseId) error {
	result, err := c.ReevaluateInaccessibleDatabaseState(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ReevaluateInaccessibleDatabaseState: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after ReevaluateInaccessibleDatabaseState: %+v", err)
	}

	return nil
}
