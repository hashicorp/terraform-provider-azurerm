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
