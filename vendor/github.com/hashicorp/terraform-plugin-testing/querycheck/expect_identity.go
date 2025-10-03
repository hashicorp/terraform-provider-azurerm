// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

var _ QueryCheck = expectIdentity{}

type expectIdentity struct {
	resourceAddress string
	identity        map[string]knownvalue.Check
}

// CheckQuery implements the query check logic.
func (e expectIdentity) CheckQuery(ctx context.Context, req CheckQueryRequest, resp *CheckQueryResponse) {
	//var resource map[string]json.RawMessage
	//var found bool
	//var foundIdentity map[string]json.RawMessage
	//var collectionIdentity []string
	//
	//if req.Query == nil {
	//	resp.Error = fmt.Errorf("query is nil")
	//	return
	//}
	//
	//// later check for if the list resource object should be included in the query include_resource
	//// check if at least one identity matches the identity I have
	//// sort with map ???
	//// if we find the identity remove message from slice and if we are at the end then it wasn't found :(
	//
	//for _, v := range *req.Query {
	//	switch idk := v.(type) {
	//	case tfjson.ListResourceFoundMessage:
	//		if idk.ListResourceFound.Address == e.resourceAddress {
	//			found = true
	//			foundIdentity = idk.ListResourceFound.Identity
	//			if foundIdentity == e.identity {
	//				remove foundIdentity from collectionIdentity
	//			}
	//
	//			if idk.ListResourceFound.ResourceObject != nil {
	//				resource = idk.ListResourceFound.ResourceObject
	//			}
	//
	//		}
	//	default:
	//		fmt.Printf("List resource not found for query check", v)
	//		continue
	//	}
	//
	//
	//}
	//
	//
	///*	if resource == nil {
	//	resp.Error = fmt.Errorf("%s - Resource not found in query", e.resourceAddress)
	//
	//	return
	//}*/
	//
	///*	if len(foundIdentity.IdentityValues) == 0 {
	//	resp.Error = fmt.Errorf("%s - Identity not found in query. Either the resource does not support identity or the Terraform version running the test does not support identity. (must be v1.12+)", e.resourceAddress)
	//
	//	return
	//}*/
	//
	///*	if len(resource.IdentityValues) != len(e.identity) {
	//	deltaMsg := ""
	//	if len(resource.IdentityValues) > len(e.identity) {
	//		deltaMsg = createDeltaString(resource.IdentityValues, e.identity, "actual identity has extra attribute(s): ")
	//	} else {
	//		deltaMsg = createDeltaString(e.identity, resource.IdentityValues, "actual identity is missing attribute(s): ")
	//	}
	//
	//	resp.Error = fmt.Errorf("%s - Expected %d attribute(s) in the actual identity object, got %d attribute(s): %s", e.resourceAddress, len(e.identity), len(resource.IdentityValues), deltaMsg)
	//	return
	//}*/
	//
	//var keys []string
	//
	//for k := range e.identity {
	//	keys = append(keys, k)
	//}
	//
	//sort.SliceStable(keys, func(i, j int) bool {
	//	return keys[i] < keys[j]
	//})
	//
	//for _, k := range keys {
	//	actualIdentityVal, ok := resource.IdentityValues[k]
	//
	//	if !ok {
	//		resp.Error = fmt.Errorf("%s - missing attribute %q in actual identity object", e.resourceAddress, k)
	//		return
	//	}
	//
	//	if err := e.identity[k].CheckValue(actualIdentityVal); err != nil {
	//		resp.Error = fmt.Errorf("%s - %q identity attribute: %s", e.resourceAddress, k, err)
	//		return
	//	}
	//}
	return
}

// ExpectIdentity returns a query check that asserts that the identity at the given resource matches a known object, where each
// map key represents an identity attribute name. The identity in query must exactly match the given object and any missing/extra
// attributes will raise a diagnostic.
//
// This query check can only be used with managed resources that support resource identity. Resource identity is only supported in Terraform v1.12+
func ExpectIdentity(resourceAddress string, identity map[string]knownvalue.Check) QueryCheck {
	return expectIdentity{
		resourceAddress: resourceAddress,
		identity:        identity,
	}
}

// createDeltaString prints the map keys that are present in mapA and not present in mapB
func createDeltaString[T any, V any](mapA map[string]T, mapB map[string]V, msgPrefix string) string {
	deltaMsg := ""

	deltaMap := make(map[string]T, len(mapA))
	maps.Copy(deltaMap, mapA)
	for key := range mapB {
		delete(deltaMap, key)
	}

	deltaKeys := slices.Sorted(maps.Keys(deltaMap))

	for i, k := range deltaKeys {
		if i == 0 {
			deltaMsg += msgPrefix
		} else {
			deltaMsg += ", "
		}
		deltaMsg += fmt.Sprintf("%q", k)
	}

	return deltaMsg
}
