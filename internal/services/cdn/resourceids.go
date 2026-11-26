// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn

// CDN
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Endpoint -rewrite=true -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/endpoints/endpoint1

// CDN FrontDoor "Associations"
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FrontDoorCustomDomainAssociation -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/associations/assoc1
