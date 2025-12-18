package databaselink

import "slices"

type ForceLinkInvocation struct {
	Id                string
	LinkedDatabaseIds []string
}

// POST /databases/{databaseId}/forceLinkToReplicationGroup API can only link one database at a time into a GR
// (geo-replication) group. The API has to be called on the database that's not yet in the GR group
//
// Given an input fromIds: ['a'] and idsToLink: ['b', 'c'], following is the sequence of API invocations that should be made:
//   - POST /databases/b/forceLinkToReplicationGroup with linkedDatabaseIds: ['a', 'b'] (POST /databases/a also works since neither are in the GR group yet)
//   - POST /databases/c/forceLinkToReplicationGroup with linkedDatabaseIds: ['a', 'b', 'c'] (API call has to be made against 'c' because 'a' and 'b' are already in the GR group)
//
// This function returns a slice of invocation parameters that will be passed as API parameters
func ForceLinkInvocations(fromIds, idsToLink []string) (res []ForceLinkInvocation) {
	linkedIds := slices.Clone(fromIds)

	for _, id := range idsToLink {
		linkedIds = append(linkedIds, id)
		res = append(res, ForceLinkInvocation{
			Id:                id,
			LinkedDatabaseIds: slices.Clone(linkedIds),
		})
	}

	return
}
