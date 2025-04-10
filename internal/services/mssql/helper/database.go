// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helper

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/resources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/replicationlinks"
)

// FindDatabaseReplicationPartners looks for partner databases having one of the specified replication roles, by
// reading any replication links then attempting to discover and match the corresponding server/database resources for
// the other end of the link.
func FindDatabaseReplicationPartners(ctx context.Context, databasesClient *databases.DatabasesClient, replicationLinksClient *replicationlinks.ReplicationLinksClient, resourcesClient *resources.ResourcesClient, id commonids.SqlDatabaseId, primaryEnclaveType databases.AlwaysEncryptedEnclaveType, rolesToFind []replicationlinks.ReplicationRole) ([]databases.Database, error) {
	var partnerDatabases []databases.Database

	log.Printf("[INFO] Looking For Replication Link Partners For: %s", id)

	matchesRole := func(role replicationlinks.ReplicationRole) bool {
		for _, r := range rolesToFind {
			if r == role {
				return true
			}
		}
		return false
	}

	listByDatabaseResult, err := replicationLinksClient.ListByDatabaseComplete(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("reading Replication Links for %s: %+v", id, err)
	}

	// Not 200?
	if !response.WasStatusCode(listByDatabaseResult.LatestHttpResponse, 200) {
		return nil, fmt.Errorf("reading Replication Links for %s: response was empty", id)
	}

	log.Printf("[INFO] Found %d Replication Link Partners", len(listByDatabaseResult.Items))

	for _, replicationLink := range listByDatabaseResult.Items {
		linkProps := replicationLink.Properties
		if linkProps == nil {
			log.Printf("[INFO] Replication Link Properties were nil for %s", id)
			continue
		}

		if linkProps.PartnerLocation == nil || linkProps.PartnerServer == nil || linkProps.PartnerDatabase == nil {
			log.Printf("[INFO] Replication Link Properties were invalid for %s", id)
			continue
		}

		log.Printf("[INFO] Replication Link found for %s", id)
		// Look for candidate partner SQL servers
		filter := fmt.Sprintf("(resourceType eq 'Microsoft.Sql/servers') and ((name eq '%s'))", *linkProps.PartnerServer)

		partner, err := commonids.ParseSqlDatabaseIDInsensitively(*linkProps.PartnerDatabaseId)
		if err != nil {
			return nil, fmt.Errorf("parsing Partner SQL Database ID %q: %+v", *linkProps.PartnerDatabaseId, err)
		}

		fmtString := fmt.Sprintf("/subscriptions/%s", partner.SubscriptionId)
		subscription, err := commonids.ParseSubscriptionID(fmtString)
		if err != nil {
			return nil, fmt.Errorf("parsing Partner SQL Server Subscription %q: %+v", fmtString, err)
		}

		listOptions := resources.ListOperationOptions{
			Expand: nil,
			Filter: pointer.To(filter),
			Top:    pointer.FromInt64(100),
		}

		resourceListResult, err := resourcesClient.ListComplete(ctx, *subscription, listOptions)
		if err != nil {
			return nil, fmt.Errorf("retrieving Partner SQL Servers with filter %q for %s: %+v", filter, *linkProps.PartnerDatabaseId, err)
		}

		for _, server := range resourceListResult.Items {
			if server.Id == nil {
				log.Printf("[INFO] Partner SQL Server ID was nil for %s", *server.Id)
				continue
			}

			partnerServerId, err := commonids.ParseSqlServerIDInsensitively(*server.Id)
			if err != nil {
				return nil, fmt.Errorf("parsing Partner SQL Server ID %q: %+v", *server.Id, err) // error
			}

			// Check if like-named server has a database named like the partner database, also with a replication link
			linksPossiblePartnerIterator, err := replicationLinksClient.ListByServerComplete(ctx, *partnerServerId)
			if err != nil {
				if response.WasNotFound(linksPossiblePartnerIterator.LatestHttpResponse) {
					log.Printf("[INFO] no replication link found for Database %q (%s)", *linkProps.PartnerDatabase, partnerServerId)
					continue
				}
				return nil, fmt.Errorf("reading Replication Links for Database %s (%s): %+v", *linkProps.PartnerDatabase, partnerServerId, err)
			}

			for _, v := range linksPossiblePartnerIterator.Items {
				linkPossiblePartner := v

				if linkPossiblePartner.Properties == nil {
					log.Printf("[INFO] Replication Link Properties was nil for Database %s (%s)", *linkProps.PartnerDatabase, partnerServerId)
					continue
				}

				linkPropsPossiblePartner := *linkPossiblePartner.Properties

				// If the database has a replication link for the specified role, we'll consider it a partner of this database if the location is the same as expected partner
				if matchesRole(*linkPropsPossiblePartner.PartnerRole) {
					partnerServerDatabase, err := commonids.ParseSqlDatabaseIDInsensitively(*linkProps.PartnerDatabaseId)
					if err != nil {
						return nil, fmt.Errorf("parsing Partner SQL Database ID %q: %+v", *linkProps.PartnerDatabaseId, err)
					}

					partnerDatabaseId := commonids.NewSqlDatabaseID(partnerServerId.SubscriptionId, partnerServerId.ResourceGroupName, partnerServerId.ServerName, partnerServerDatabase.DatabaseName)
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
					} else {
						log.Printf("[INFO] PartnerDatabase.Model was 'nil'")
					}
				}
			}
		}
	}

	return partnerDatabases, nil
}
