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
	client DatabasesClientInterface
	id     databases.DatabaseId
	toIds  []string
}

func NewGeoReplicationPoller(client DatabasesClientInterface, id databases.DatabaseId, toIds []string) *geoReplicationPoller {
	return &geoReplicationPoller{
		client: client,
		id:     id,
		toIds:  toIds,
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

	idFound := make(map[string]bool, len(p.toIds))
	for _, toId := range p.toIds {
		idFound[toId] = false
	}

	for _, ldb := range pointer.From(resp.Model.Properties.GeoReplication.LinkedDatabases) {
		if pointer.From(ldb.State) == databases.LinkStateLinked {
			id := pointer.From(ldb.Id)
			if _, ok := idFound[id]; ok {
				idFound[id] = true
			}
		}
	}

	for _, found := range idFound {
		if !found {
			return &pollingInProgress, nil
		}
	}

	return &pollingSuccess, nil
}
