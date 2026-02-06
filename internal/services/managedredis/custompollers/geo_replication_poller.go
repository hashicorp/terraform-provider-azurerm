// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/databases"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &geoReplicationPoller{}

type geoReplicationPoller struct {
	client              DatabasesClientInterface
	id                  databases.DatabaseId
	expectedLinkedDbIds []string
}

func NewGeoReplicationPoller(client DatabasesClientInterface, id databases.DatabaseId, toIds []string) *geoReplicationPoller {
	return &geoReplicationPoller{
		client:              client,
		id:                  id,
		expectedLinkedDbIds: toIds,
	}
}

func (p geoReplicationPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil {
		return nil, fmt.Errorf("polling for %s: properties were empty", p.id)
	}

	if resp.Model.Properties.GeoReplication == nil {
		return nil, fmt.Errorf("polling for %s: properties.geoReplication were empty", p.id)
	}

	linkedDatabaseIds := make(map[string]struct{}, len(p.expectedLinkedDbIds))

	for _, ldb := range pointer.From(resp.Model.Properties.GeoReplication.LinkedDatabases) {
		if pointer.From(ldb.State) == databases.LinkStateLinked {
			linkedDatabaseIds[pointer.From(ldb.Id)] = struct{}{}
		}
	}

	for _, ldb := range p.expectedLinkedDbIds {
		if _, ok := linkedDatabaseIds[ldb]; !ok {
			return &pollingInProgress, nil
		}
	}

	return &pollingSuccess, nil
}
