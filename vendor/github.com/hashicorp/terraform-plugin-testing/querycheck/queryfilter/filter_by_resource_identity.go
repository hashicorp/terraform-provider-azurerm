// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package queryfilter

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

type filterByResourceIdentity struct {
	identity map[string]knownvalue.Check
}

func (f filterByResourceIdentity) Filter(ctx context.Context, req FilterQueryRequest, resp *FilterQueryResponse) {
	if len(req.QueryItem.Identity) != len(f.identity) {
		resp.Include = false
		return
	}

	var keys []string

	for k := range f.identity {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, k := range keys {
		actualIdentityVal, ok := req.QueryItem.Identity[k]

		if !ok {
			resp.Include = false
			return
		}

		if err := f.identity[k].CheckValue(actualIdentityVal); err != nil {
			resp.Include = false
			return
		}
	}

	resp.Include = true
}

// ByResourceIdentity returns a query filter that only includes query items that match
// the given resource identity.
//
// Errors thrown by the given known value checks are only used to filter out non-matching query
// items and are otherwise ignored.
func ByResourceIdentity(identity map[string]knownvalue.Check) QueryFilter {
	return filterByResourceIdentity{
		identity: identity,
	}
}
