// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helper

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql"           // nolint: staticcheck
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/databases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// FindDatabaseReplicationPartners looks for partner databases having one of the specified replication roles, by
// reading any replication links then attempting to discover and match the corresponding server/database resources for
// the other end of the link.
func FindDatabaseReplicationPartners(ctx context.Context, databasesClient *databases.DatabasesClient, replicationLinksClient *sql.ReplicationLinksClient, resourcesClient *resources.Client, id commonids.SqlDatabaseId, primaryEnclaveType databases.AlwaysEncryptedEnclaveType, rolesToFind []sql.ReplicationRole) ([]databases.Database, error) {
	var partnerDatabases []databases.Database

	matchesRole := func(role sql.ReplicationRole) bool {
		for _, r := range rolesToFind {
			if r == role {
				return true
			}
		}
		return false
	}

	for linksIterator, err := replicationLinksClient.ListByDatabaseComplete(ctx, id.ResourceGroupName, id.ServerName, id.DatabaseName); linksIterator.NotDone(); err = linksIterator.NextWithContext(ctx) {
		if err != nil {
			return nil, fmt.Errorf("reading Replication Links for %s: %+v", id, err)
		}

		if linksIterator.Response().IsEmpty() {
			return nil, fmt.Errorf("reading Replication Links for %s: response was empty", id)
		}

		linkProps := linksIterator.Value().ReplicationLinkProperties
		if linkProps == nil {
			log.Printf("[INFO] Replication Link Properties was nil for %s", id)
			continue
		}

		if linkProps.PartnerLocation == nil || linkProps.PartnerServer == nil || linkProps.PartnerDatabase == nil {
			log.Printf("[INFO] Replication Link Properties was invalid for %s", id)
			continue
		}

		log.Printf("[INFO] Replication Link found for %s", id)

		// Look for candidate partner SQL servers
		filter := fmt.Sprintf("(resourceType eq 'Microsoft.Sql/servers') and ((name eq '%s'))", *linkProps.PartnerServer)
		var resourceList []resources.GenericResourceExpanded
		for resourcesIterator, err := resourcesClient.ListComplete(ctx, filter, "", nil); resourcesIterator.NotDone(); err = resourcesIterator.NextWithContext(ctx) {
			if err != nil {
				return nil, fmt.Errorf("retrieving Partner SQL Servers with filter %q for %s: %+v", filter, id, err)
			}
			resourceList = append(resourceList, resourcesIterator.Value())
		}

		for _, server := range resourceList {
			if server.ID == nil {
				log.Printf("[INFO] Partner SQL Server ID was nil for %s", id)
				continue
			}

			partnerServerId, err := parse.ServerID(*server.ID)
			if err != nil {
				return nil, fmt.Errorf("parsing Partner SQL Server ID %q: %+v", *server.ID, err)
			}

			// Check if like-named server has a database named like the partner database, also with a replication link
			for linksPossiblePartnerIterator, err := replicationLinksClient.ListByDatabaseComplete(ctx, partnerServerId.ResourceGroup, partnerServerId.Name, *linkProps.PartnerDatabase); linksPossiblePartnerIterator.NotDone(); err = linksPossiblePartnerIterator.NextWithContext(ctx) {
				if err != nil {
					if utils.ResponseWasNotFound(linksPossiblePartnerIterator.Response().Response) {
						log.Printf("[INFO] no replication link found for Database %q (%s)", *linkProps.PartnerDatabase, partnerServerId)
						continue
					}
					return nil, fmt.Errorf("reading Replication Links for Database %s (%s): %+v", *linkProps.PartnerDatabase, partnerServerId, err)
				}

				linkPossiblePartner := linksPossiblePartnerIterator.Value()
				if linkPossiblePartner.ReplicationLinkProperties == nil {
					log.Printf("[INFO] Replication Link Properties was nil for Database %s (%s)", *linkProps.PartnerDatabase, partnerServerId)
					continue
				}

				linkPropsPossiblePartner := *linkPossiblePartner.ReplicationLinkProperties

				// If the database has a replication link for the specified role, we'll consider it a partner of this database if the location is the same as expected partner
				if matchesRole(linkPropsPossiblePartner.Role) {
					partnerDatabaseId := commonids.NewSqlDatabaseID(partnerServerId.SubscriptionId, partnerServerId.ResourceGroup, partnerServerId.Name, *linkProps.PartnerDatabase)
					partnerDatabase, err := databasesClient.Get(ctx, partnerDatabaseId, databases.DefaultGetOperationOptions())
					if err != nil {
						return nil, fmt.Errorf("retrieving Partner %s: %+v", partnerDatabaseId, err)
					}

					if partnerDatabase := partnerDatabase.Model; partnerDatabase != nil {
						if location.Normalize(partnerDatabase.Location) != location.Normalize(*linkProps.PartnerLocation) {
							log.Printf("[INFO] Mismatch of possible Partner Database based on location (%s vs %s) for %s", location.Normalize(partnerDatabase.Location), location.Normalize(*linkProps.PartnerLocation), id)
							continue
						}

						if partnerDatabase.Id != nil && partnerDatabase.Properties != nil && partnerDatabase.Properties.PreferredEnclaveType != nil {
							if primaryEnclaveType != "" && primaryEnclaveType == *partnerDatabase.Properties.PreferredEnclaveType {
								log.Printf("[INFO] Found Partner %s", partnerDatabaseId)
								partnerDatabases = append(partnerDatabases, *partnerDatabase)
							} else {
								log.Printf("[INFO] Mismatch of possible Partner Database based on enclave type (%q vs %q) for %s", primaryEnclaveType, string(*partnerDatabase.Properties.PreferredEnclaveType), id)
							}
						}
					}
				}
			}
		}
	}

	return partnerDatabases, nil
}
