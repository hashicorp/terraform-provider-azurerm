package backupautomaticandondemands

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

type BackupsAutomaticAndOnDemandCreateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// BackupsAutomaticAndOnDemandCreate ...
func (c BackupAutomaticAndOnDemandsClient) BackupsAutomaticAndOnDemandCreate(ctx context.Context, id BackupId) (result BackupsAutomaticAndOnDemandCreateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
		},
		HttpMethod: http.MethodPut,
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

// BackupsAutomaticAndOnDemandCreateThenPoll performs BackupsAutomaticAndOnDemandCreate then polls until it's completed
func (c BackupAutomaticAndOnDemandsClient) BackupsAutomaticAndOnDemandCreateThenPoll(ctx context.Context, id BackupId) error {
	result, err := c.BackupsAutomaticAndOnDemandCreate(ctx, id)
	if err != nil {
		return fmt.Errorf("performing BackupsAutomaticAndOnDemandCreate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after BackupsAutomaticAndOnDemandCreate: %+v", err)
	}

	return nil
}
