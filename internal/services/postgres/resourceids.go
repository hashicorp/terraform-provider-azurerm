// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package postgres

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AzureActiveDirectoryAdministrator -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DBforPostgreSQL/servers/server1/administrators/activeDirectory

// This is being used by the state migration in postgresql_aad_administrator.go which I'm not sure is actually correct, but seeing as it's a state migration it can't be changed now
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SqlAzureActiveDirectoryAdministrator -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Sql/servers/server1/administrators/activeDirectory -rewrite=true
