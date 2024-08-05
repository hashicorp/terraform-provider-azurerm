// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/provider"
)

type resource struct {
	name   string
	schema interface{}
}

type Resources struct {
	resources []resource
}

type set map[string]struct{}

func (s set) Exists(key string) bool {
	key = strings.ReplaceAll(key, " ", "")
	_, ok := s[strings.ToLower(key)]
	return ok
}

func newSet(val string) set {
	if val == "" {
		return nil
	}
	vals := strings.Split(val, ",")
	res := make(set, len(vals))
	for _, name := range vals {
		name = strings.ReplaceAll(name, " ", "")
		res[strings.ToLower(name)] = struct{}{}
	}
	return res
}

func AzurermAllResources(service, skipService string, resources, skipResources string) Resources {
	var (
		rps              = newSet(service)
		skipRPs          = newSet(skipService)
		resourcesMap     = newSet(resources)
		skipResourcesMap = newSet(skipResources)
	)

	shouldSkipRP := func(name string) bool {
		if len(rps) > 0 && !rps.Exists(name) {
			return true
		}
		if skipRPs.Exists(name) {
			return true
		}
		return false
	}

	shouldSKipResource := func(name string) bool {
		if len(resourcesMap) > 0 && !resourcesMap.Exists(name) {
			return true
		}
		if skipResourcesMap.Exists(name) {
			return true
		}
		return false
	}

	var res Resources
	for _, r := range provider.SupportedTypedServices() {
		name := r.Name()
		if shouldSkipRP(name) {
			continue
		}
		for _, svc := range r.Resources() {
			if shouldSKipResource(svc.ResourceType()) {
				continue
			}
			res.resources = append(res.resources, resource{
				name:   svc.ResourceType(),
				schema: svc,
			})
		}
	}

	for _, r := range provider.SupportedUntypedServices() {
		name := r.Name()
		if shouldSkipRP(name) {
			continue
		}
		for name, svc := range r.SupportedResources() {
			if shouldSKipResource(name) {
				continue
			}
			res.resources = append(res.resources, resource{
				name:   name,
				schema: svc,
			})
		}
	}
	return res
}
