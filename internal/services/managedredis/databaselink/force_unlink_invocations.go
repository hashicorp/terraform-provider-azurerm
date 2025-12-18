// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package databaselink

type ForceUnlinkInvocation struct {
	Id  string
	Ids []string
}

// POST /databases/{databaseId}/forceUnlink can only unlink one database at a time from a GR (geo-replication) group.
// The API has to be called on the database that's currently in the GR group (intermediateIds).
//
// Given an input intermediateIds: ['a', 'b'] and idsToUnlink: ['c', 'd'], following is the sequence of API invocations that should be made:
//   - POST /databases/a/forceUnlink with ids: ['c'] (POST /databases/{id} cannot be called against itself)
//   - POST /databases/a/forceUnlink with ids: ['d']
//
// This function returns a slice of invocation parameters that will be passed as API parameters
func ForceUnlinkInvocations(intermediateIds, idsToUnlink []string) (res []ForceUnlinkInvocation) {
	if len(intermediateIds) > 0 {
		id := intermediateIds[0]

		for _, idToUnlink := range idsToUnlink {
			res = append(res, ForceUnlinkInvocation{
				Id:  id,
				Ids: []string{idToUnlink},
			})
		}
	}
	return
}
