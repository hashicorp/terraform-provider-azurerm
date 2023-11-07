// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helper

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/replicationlinks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
)

// FindDatabaseReplicationPartners looks for partner databases having one of the specified replication roles, by
// reading any replication links then attempting to discover and match the corresponding server/database resources for
// the other end of the link.
func FindDatabaseReplicationPartners(ctx context.Context, databasesClient *databases.DatabasesClient, replicationLinksClient *replicationlinks.ReplicationLinksClient, resourcesClient *resources.Client, id commonids.SqlDatabaseId, rolesToFind []replicationlinks.ReplicationRole) ([]databases.Database, error) {
	var partnerDatabases []databases.Database

	matchesRole := func(role replicationlinks.ReplicationRole) bool {
		for _, r := range rolesToFind {
			if r == role {
				return true
			}
		}
		return false
	}

	replicationDatabaseId := commonids.NewSqlDatabaseID(id.SubscriptionId, id.ResourceGroupName, id.ServerName, id.DatabaseName)

	results, err := replicationLinksClient.ListByDatabaseComplete(ctx, replicationDatabaseId)
	if err != nil {
		// Not sure if this is really an error anymore given the change in the API...
		if strings.Contains(err.Error(), "ResourceNotFound") {
			log.Printf("[INFO] %s does not exist, skipping lookup for Replication Links", id)
			return partnerDatabases, nil
		}
		return nil, fmt.Errorf("reading Replication Links for %s: %+v", id, err)
	}

	// Not sure if this is really an error anymore given the change in the API...
	if len(results.Items) == 0 {
		// return nil, fmt.Errorf("reading Replication Links for %s: replicationlinks.ListByDatabaseCompleteResult was empty", id)
		log.Printf("[INFO] reading Replication Links for %s: replicationlinks.ListByDatabaseCompleteResult was empty", id)
		return partnerDatabases, nil
	}

	var linkProps *replicationlinks.ReplicationLinkProperties

	// loop over all results that matched the replicationDatabaseId...
	for _, v := range results.Items {
		if linkProps = v.Properties; linkProps != nil {
			if linkProps.PartnerLocation == nil || linkProps.PartnerServer == nil || linkProps.PartnerDatabase == nil {
				log.Printf("[INFO] Replication Link Properties were invalid for %s", id)
				continue
			}

			log.Printf("[INFO] Replication Link found for %s", id)

			// Look for candidate partner SQL servers
			filter := fmt.Sprintf("(resourceType eq 'Microsoft.Sql/servers') and ((name eq '%s'))", pointer.From(linkProps.PartnerServer))
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

				partnerServerId, err := parse.ServerID(pointer.From(server.ID))
				if err != nil {
					return nil, fmt.Errorf("parsing Partner SQL Server ID %q: %+v", pointer.From(server.ID), err)
				}

				partnerDatabaseId := commonids.NewSqlDatabaseID(partnerServerId.SubscriptionId, partnerServerId.ResourceGroup, partnerServerId.Name, pointer.From(linkProps.PartnerDatabase))

				// Check if like-named server has a database named like the partner database, also with a replication link
				linksPossiblePartnerIterator, err := replicationLinksClient.ListByDatabaseComplete(ctx, partnerDatabaseId)
				if err != nil {
					return nil, fmt.Errorf("reading Replication Links for Database %q (%s): %+v", partnerDatabaseId.DatabaseName, partnerServerId, err)
				}

				for _, linkPossiblePartner := range linksPossiblePartnerIterator.Items {

					if linkPossiblePartner.Properties == nil {
						log.Printf("[INFO] Replication Link Properties was nil for Database %q (%s)", pointer.From(linkProps.PartnerDatabase), partnerServerId)
						continue
					}

					linkPropsPossiblePartner := pointer.From(linkPossiblePartner.Properties)

					// If the database has a replication link for the specified role, we'll consider it a partner of this database if the location is the same as expected partner
					if matchesRole(pointer.From(linkPropsPossiblePartner.Role)) {
						partnerDatabaseId := commonids.NewSqlDatabaseID(partnerServerId.SubscriptionId, partnerServerId.ResourceGroup, partnerServerId.Name, pointer.From(linkProps.PartnerDatabase))

						partnerDatabase, err := databasesClient.Get(ctx, partnerDatabaseId, databases.DefaultGetOperationOptions())
						if err != nil {
							return nil, fmt.Errorf("retrieving Partner %q: %+v", partnerDatabaseId, err)
						}

						if partnerDatabase.Model == nil {
							log.Printf("[INFO] Partner Database Model is nil for %s", partnerDatabaseId)
							continue
						}

						if location.Normalize(partnerDatabase.Model.Location) != location.NormalizeNilable(linkProps.PartnerLocation) {
							log.Printf("[INFO] Mismatch of possible Partner Database based on location (%q vs %q) for %s", location.Normalize(partnerDatabase.Model.Location), location.NormalizeNilable(linkProps.PartnerLocation), id)
							continue
						}

						if partnerDatabase.Model.Id != nil {
							log.Printf("[INFO] Found Partner %s", partnerDatabaseId)
							partnerDatabases = append(partnerDatabases, pointer.From(partnerDatabase.Model))
						}
					}
				}
			}
		} else {
			log.Printf("[INFO] Replication Link Properties was nil for %s", id)
			continue
		}
	}

	log.Printf("[INFO] Replication Link found for %s", id)

	return partnerDatabases, nil
}
