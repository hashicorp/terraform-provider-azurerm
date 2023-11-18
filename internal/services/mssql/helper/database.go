// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helper

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql"           // nolint: staticcheck
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// FindDatabaseReplicationPartners looks for partner databases having one of the specified replication roles, by
// reading any replication links then attempting to discover and match the corresponding server/database resources for
// the other end of the link.
func FindDatabaseReplicationPartners(ctx context.Context, databasesClient *sql.DatabasesClient, replicationLinksClient *sql.ReplicationLinksClient, resourcesClient *resources.Client, id parse.DatabaseId, rolesToFind []sql.ReplicationRole) ([]sql.Database, error) {
	var partnerDatabases []sql.Database

	matchesRole := func(role sql.ReplicationRole) bool {
		for _, r := range rolesToFind {
			if r == role {
				return true
			}
		}
		return false
	}

	for linksIterator, err := replicationLinksClient.ListByDatabaseComplete(ctx, id.ResourceGroup, id.ServerName, id.Name); linksIterator.NotDone(); err = linksIterator.NextWithContext(ctx) {
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
					partnerDatabaseId := parse.NewDatabaseID(partnerServerId.SubscriptionId, partnerServerId.ResourceGroup, partnerServerId.Name, *linkProps.PartnerDatabase)
					partnerDatabase, err := databasesClient.Get(ctx, partnerDatabaseId.ResourceGroup, partnerDatabaseId.ServerName, partnerDatabaseId.Name)
					if err != nil {
						return nil, fmt.Errorf("retrieving Partner %s: %+v", partnerDatabaseId, err)
					}
					if location.NormalizeNilable(partnerDatabase.Location) != location.Normalize(*linkProps.PartnerLocation) {
						log.Printf("[INFO] Mismatch of possible Partner Database based on location (%s vs %s) for %s", location.NormalizeNilable(partnerDatabase.Location), location.Normalize(*linkProps.PartnerLocation), id)
						continue
					}
					if partnerDatabase.ID != nil {
						log.Printf("[INFO] Found Partner %s", partnerDatabaseId)
						partnerDatabases = append(partnerDatabases, partnerDatabase)
					}
				}
			}
		}
	}

	return partnerDatabases, nil
}
