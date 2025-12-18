// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package databaselink

// Given current state of geoReplication database ids (fromIds) and intended state (toIds), compute the idsToUnlink,
// intermediateIds (ids after unlinking is done), and idsToLink.
//
// Example 1:
//   - fromIds: ['a']
//   - toIds: ['a', 'b', 'c']
//
// Returns ([], ['a'], ['b', 'c'])
//
// Example 2:
//   - fromIds: ['a', 'b', 'c']
//   - toIds: ['a', 'b']
//
// Returns (['c'], ['a', 'b'], [])
//
// Example 3:
//   - fromIds: ['a', 'b', 'c']
//   - toIds: ['b', 'c', 'd']
//
// Returns (['a'], ['b', 'c'], ['d'])
func LinkUnlink(fromIds, toIds []string) (idsToUnlink, intermediateIds, idsToLink []string) {
	fromMap, toMap := make(map[string]bool, len(fromIds)), make(map[string]bool, len(toIds))
	for _, id := range fromIds {
		fromMap[id] = true
	}
	for _, id := range toIds {
		toMap[id] = true
	}

	for _, id := range fromIds {
		if !toMap[id] {
			idsToUnlink = append(idsToUnlink, id)
		} else {
			intermediateIds = append(intermediateIds, id)
		}
	}

	for _, id := range toIds {
		if !fromMap[id] {
			idsToLink = append(idsToLink, id)
		}
	}

	return
}
