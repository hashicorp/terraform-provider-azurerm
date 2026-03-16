// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package databaselink

import "slices"

type LinkUnlinkInvocation struct {
	Id  string
	Ids []string
}

// POST /databases/{databaseId}/forceLinkToReplicationGroup API can only link one database at a time into a GR
// (geo-replication) group. The API has to be called on the database that's not yet in the GR group
//
// Given an input fromIds: ['a'] and idsToLink: ['b', 'c'], following is the sequence of API invocations that should be made:
//   - POST /databases/b/forceLinkToReplicationGroup with linkedDatabaseIds: ['a', 'b'] (POST /databases/a also works since neither are in the GR group yet)
//   - POST /databases/c/forceLinkToReplicationGroup with linkedDatabaseIds: ['a', 'b', 'c'] (API call has to be made against 'c' because 'a' and 'b' are already in the GR group)
//
// This function returns a slice of invocation parameters that will be passed as API parameters
func ForceLinkInvocations(fromIds, idsToLink []string) (res []LinkUnlinkInvocation) {
	linkedIds := slices.Clone(fromIds)

	for _, id := range idsToLink {
		linkedIds = append(linkedIds, id)
		res = append(res, LinkUnlinkInvocation{
			Id:  id,
			Ids: slices.Clone(linkedIds),
		})
	}

	return
}

// POST /databases/{databaseId}/forceUnlink can only unlink one database at a time from a GR (geo-replication) group.
// The API has to be called on the database that's currently in the GR group (intermediateIds).
//
// Given an input intermediateIds: ['a', 'b'] and idsToUnlink: ['c', 'd'], following is the sequence of API invocations that should be made:
//   - POST /databases/a/forceUnlink with ids: ['c'] (POST /databases/{id}/forceUnlink cannot be called against itself)
//   - POST /databases/a/forceUnlink with ids: ['d']
//
// This function returns a slice of invocation parameters that will be passed as API parameters
func ForceUnlinkInvocations(intermediateIds, idsToUnlink []string) (res []LinkUnlinkInvocation) {
	if len(intermediateIds) > 0 {
		id := intermediateIds[0]

		for _, idToUnlink := range idsToUnlink {
			res = append(res, LinkUnlinkInvocation{
				Id:  id,
				Ids: []string{idToUnlink},
			})
		}
	}
	return
}
