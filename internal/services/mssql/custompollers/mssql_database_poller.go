// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/databases"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

const (
	msSqlDatabasePollInterval            = 1 * time.Minute
	msSqlDatabaseOnlineTargetOccurrences = 2
)

var _ pollers.PollerType = &MsSqlDatabaseOnlinePoller{}

type MsSqlDatabaseOnlinePoller struct {
	client                   *databases.DatabasesClient
	id                       commonids.SqlDatabaseId
	consecutiveOnlineResults int
}

func NewMsSqlDatabaseOnlinePoller(client *databases.DatabasesClient, id commonids.SqlDatabaseId) *MsSqlDatabaseOnlinePoller {
	return &MsSqlDatabaseOnlinePoller{
		client: client,
		id:     id,
	}
}

func (p *MsSqlDatabaseOnlinePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	operationsInProgress, err := hasInProgressDatabaseOperations(ctx, p.client, p.id)
	if err != nil {
		return nil, fmt.Errorf("checking database operations for %s: %+v", p.id, err)
	}

	if operationsInProgress {
		p.consecutiveOnlineResults = 0
		return &pollers.PollResult{
			PollInterval: msSqlDatabasePollInterval,
			Status:       pollers.PollingStatusInProgress,
		}, nil
	}

	resp, err := p.client.Get(ctx, p.id, databases.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	if resp.Model == nil {
		return nil, fmt.Errorf("retrieving %s: `model` was nil", p.id)
	}

	if resp.Model.Properties == nil {
		return nil, fmt.Errorf("retrieving %s: `properties` was nil", p.id)
	}

	if resp.Model.Properties.Status == nil {
		return nil, fmt.Errorf("retrieving %s: `status` was nil", p.id)
	}

	if pointer.From(resp.Model.Properties.Status) == databases.DatabaseStatusOnline {
		p.consecutiveOnlineResults++
	} else {
		p.consecutiveOnlineResults = 0
	}

	if p.consecutiveOnlineResults >= msSqlDatabaseOnlineTargetOccurrences {
		return &pollers.PollResult{
			PollInterval: msSqlDatabasePollInterval,
			Status:       pollers.PollingStatusSucceeded,
		}, nil
	}

	return &pollers.PollResult{
		PollInterval: msSqlDatabasePollInterval,
		Status:       pollers.PollingStatusInProgress,
	}, nil
}

type databaseOperation struct {
	Properties *databaseOperationProperties `json:"properties,omitempty"`
}

type databaseOperationProperties struct {
	State *string `json:"state,omitempty"`
}

type databaseOperationsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *databaseOperationsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

func hasInProgressDatabaseOperations(ctx context.Context, client *databases.DatabasesClient, id commonids.SqlDatabaseId) (bool, error) {
	operations, err := listDatabaseOperations(ctx, client, id)
	if err != nil {
		return false, err
	}

	for _, operation := range operations {
		if operation.Properties == nil || operation.Properties.State == nil {
			continue
		}

		switch *operation.Properties.State {
		case "CancelInProgress", "InProgress", "Pending":
			return true, nil
		}
	}

	return false, nil
}

func listDatabaseOperations(ctx context.Context, databaseClient *databases.DatabasesClient, id commonids.SqlDatabaseId) ([]databaseOperation, error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &databaseOperationsListCustomPager{},
		Path:       fmt.Sprintf("%s/operations", id.ID()),
	}

	req, err := databaseClient.Client.NewRequest(ctx, opts)
	if err != nil {
		return nil, err
	}

	resp, err := req.ExecutePaged(ctx)
	if err != nil {
		return nil, err
	}

	var values struct {
		Values *[]databaseOperation `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return nil, err
	}

	if values.Values == nil {
		return []databaseOperation{}, nil
	}

	return *values.Values, nil
}
