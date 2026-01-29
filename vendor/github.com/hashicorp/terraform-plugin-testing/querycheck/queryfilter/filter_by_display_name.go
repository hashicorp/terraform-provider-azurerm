// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package queryfilter

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

type filterByDisplayName struct {
	displayNameCheck knownvalue.Check
}

func (f filterByDisplayName) Filter(ctx context.Context, req FilterQueryRequest, resp *FilterQueryResponse) {
	if err := f.displayNameCheck.CheckValue(req.QueryItem.DisplayName); err == nil {
		resp.Include = true
		return
	}
}

// ByDisplayNameExact returns a query filter that only includes query items that match
// the specified display name.
func ByDisplayName(displayNameCheck knownvalue.Check) QueryFilter {
	return filterByDisplayName{
		displayNameCheck: displayNameCheck,
	}
}
