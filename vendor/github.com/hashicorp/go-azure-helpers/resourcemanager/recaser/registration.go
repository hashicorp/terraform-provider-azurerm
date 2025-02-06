// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recaser

import (
	"strings"
	"sync"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var knownResourceIds = make(map[string]resourceids.ResourceId)

// KnownResourceIds returns the map of resource IDs that have been registered by each API imported via the
// RegisterResourceId function. This is the case for all APIs generated via the Pandora project via init().
// The keys for the map are the lower-cased ID strings with the user-specified segments
// stripped out, leaving the path intact. Example:
// "/subscriptions//resourceGroups//providers/Microsoft.BotService/botServices/"
func KnownResourceIds() map[string]resourceids.ResourceId {
	return knownResourceIds
}

var resourceIdsWriteLock = &sync.Mutex{}

func init() {
	//register common ids
	for _, id := range commonids.CommonIds() {
		RegisterResourceId(id)
	}
}

// RegisterResourceId adds ResourceIds to a list of known ids
func RegisterResourceId(id resourceids.ResourceId) {
	key := strings.ToLower(id.ID())

	resourceIdsWriteLock.Lock()
	if _, ok := knownResourceIds[key]; !ok {
		knownResourceIds[key] = id
	}
	resourceIdsWriteLock.Unlock()
}
