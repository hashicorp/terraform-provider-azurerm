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
	log.Printf("[INFO] Looking For Replication Link Partners For: %q", id)

	matchesRole := func(role replicationlinks.ReplicationRole) bool {
		log.Printf("[INFO] Looking For Database Role(s):")

		for _, r := range rolesToFind {
			log.Printf("[INFO]    Looking for Database Role: %q", r)

			if r == role {
				log.Printf("[INFO]    Found Database Role: %q", role)
				return true
			}
		}

		log.Printf("[INFO] Expected Database Role(s) not found, invalid partner database")
		return false
	}

	// This will return a list of databases that have replication links for this database...
	linksIterator, err := replicationLinksClient.ListByDatabaseComplete(ctx, id)
	if err != nil {
		if response.WasNotFound(linksIterator.LatestHttpResponse) {
			return nil, nil
		}

		return nil, fmt.Errorf("reading Replication Links for %q: %+v", id, err)
	}

	for _, item := range linksIterator.Items {
		if item.Properties == nil {
			log.Printf("[INFO] Replication Link Properties was nil for %q", id)
			continue
		}

		linkProps := item.Properties

		if linkProps.PartnerLocation == nil || linkProps.PartnerServer == nil || linkProps.PartnerDatabase == nil {
			log.Printf("[INFO] Replication Link Properties were invalid for %q", id)
			continue
		}

		log.Printf("[INFO] Replication Link found for %q", id)

		// Look for candidate partner SQL servers
		linkDatabase, err := commonids.ParseSqlDatabaseIDInsensitively(*linkProps.PartnerDatabaseId)
		if err != nil {
			return nil, fmt.Errorf("parsing Replication Link SQL Database ID %q: %+v", *linkProps.PartnerDatabaseId, err)
		}

		fmtString := fmt.Sprintf("/subscriptions/%s", linkDatabase.SubscriptionId)
		linkSubscription, err := commonids.ParseSubscriptionID(fmtString)
		if err != nil {
			return nil, fmt.Errorf("parsing Replication Link SQL Database Subscription %q: %+v", fmtString, err)
		}

		log.Printf("[INFO] Replication Link SQL Databases Subscription: %q", *linkSubscription)

		filter := fmt.Sprintf("(resourceType eq 'Microsoft.Sql/servers') and ((name eq '%s'))", linkDatabase.ServerName)
		listOptions := resources.ListOperationOptions{
			Expand: nil,
			Filter: pointer.To(filter),
			Top:    pointer.FromInt64(100),
		}

		resourcesIterator, err := resourcesClient.ListComplete(ctx, *linkSubscription, listOptions)
		if err != nil {
			return nil, fmt.Errorf("retrieving Partner SQL Servers with filter %q for %q: %+v", filter, linkDatabase, err)
		}

		log.Printf("[INFO] Found %d Partner SQL Servers with filter: %q", len(resourcesIterator.Items), filter)

		for _, server := range resourcesIterator.Items {
			if server.Id == nil {
				log.Printf("[INFO] Partner SQL Server ID was nil for %q", id)
				continue
			}

			partnerServerId, err := commonids.ParseSqlServerIDInsensitively(*server.Id)
			if err != nil {
				return nil, fmt.Errorf("parsing Partner SQL Server ID %q: %+v", *server.Id, err)
			}

			log.Printf("[INFO] Parsed Partner SQL Server ID: %q", partnerServerId)

			// Check if like-named server has a database named like the partner database, also with a replication link
			partnerDatabase, err := commonids.ParseSqlDatabaseIDInsensitively(*linkProps.PartnerDatabaseId)
			if err != nil {
				return nil, fmt.Errorf("parsing Partner SQL Database ID %q: %+v", *server.Id, err)
			}

			linksPossiblePartnerIterator, err := replicationLinksClient.ListByDatabaseComplete(ctx, *partnerDatabase)
			if err != nil {
				if response.WasNotFound(linksPossiblePartnerIterator.LatestHttpResponse) {
					log.Printf("[INFO] no replication link found for SQL Database %q (%q)", *linkProps.PartnerDatabase, partnerServerId)
					continue
				}
				return nil, fmt.Errorf("reading Replication Links for SQL Database %s (%q): %+v", partnerDatabase, partnerServerId, err)
			}

			for _, linkPossiblePartnerItem := range linksPossiblePartnerIterator.Items {
				if linkPossiblePartnerItem.Properties == nil {
					log.Printf("[INFO] Replication Link Properties was nil for SQL Database %q (%q)", partnerDatabase, partnerServerId)
					continue
				}

				if linkPossiblePartnerItem.Properties == nil {
					log.Printf("[INFO] Replication Link Properties was nil for SQL Database %q (%q)", partnerDatabase, partnerServerId)
					continue
				}

				linkPropsPossiblePartner := *linkPossiblePartnerItem.Properties

				// If the database has a replication link for the specified role, we'll consider it a partner of this database if the location is the same as expected partner
				if matchesRole(pointer.From(linkPropsPossiblePartner.Role)) {
					partnerDatabaseId := commonids.NewSqlDatabaseID(partnerServerId.SubscriptionId, partnerServerId.ResourceGroupName, partnerServerId.ServerName, *linkProps.PartnerDatabase)
					partnerDatabase, err := databasesClient.Get(ctx, partnerDatabaseId, databases.DefaultGetOperationOptions())
					if err != nil {
						return nil, fmt.Errorf("retrieving SQL Partner Database %q: %+v", partnerDatabaseId, err)
					}

					log.Printf("[INFO] Replication Link Partner SQL Database ID: %q", partnerDatabaseId)

					if partnerDatabase := partnerDatabase.Model; partnerDatabase != nil {
						partnerDatabaseProps := partnerDatabase.Properties
						if partnerDatabaseProps == nil {
							return nil, fmt.Errorf("Partner SQL Database Properties were nil")
						}

						log.Printf("[INFO] Partner SQL Database Location: %q :: Replication Link Partner SQL Database Location: %q", location.Normalize(partnerDatabase.Location), location.Normalize(*linkProps.PartnerLocation))

						if location.Normalize(partnerDatabase.Location) != location.Normalize(*linkProps.PartnerLocation) {
							log.Printf("[INFO] Mismatch of possible Partner SQL Database on location (%q vs %q) for %q", location.Normalize(partnerDatabase.Location), location.Normalize(*linkProps.PartnerLocation), id)
							continue
						}

						// NOTE: nil or an empty string means that the PreferredEnclaveType is 'disabled', the
						// partner database must have the same 'PreferredEnclaveType' as the database we are
						// looking for else it is not a Partner Database...
						partnerDatabasePropsPreferredEnclaveType := ""
						if partnerDatabaseProps.PreferredEnclaveType != nil {
							partnerDatabasePropsPreferredEnclaveType = string(*partnerDatabaseProps.PreferredEnclaveType)
						}

						log.Printf("[INFO] SQL Database Preferred Enclave Type: %q :: Partner SQL Database Preferred Enclave Type: %q", primaryEnclaveType, partnerDatabasePropsPreferredEnclaveType)

						if partnerDatabase.Id != nil && partnerDatabaseProps != nil && partnerDatabasePropsPreferredEnclaveType == string(primaryEnclaveType) {
							log.Printf("[INFO] Found Partner SQL Database ID: %s", partnerDatabaseId)
							partnerDatabases = append(partnerDatabases, *partnerDatabase)
						} else {
							log.Printf("[INFO] Mismatch of possible Partner SQL Database on preferred enclave type (%q vs %q) for %q", primaryEnclaveType, partnerDatabasePropsPreferredEnclaveType, id)
						}
					} else {
						log.Printf("[INFO] partnerDatabase.Model %q was nil", partnerDatabaseId)
					}
				}
			}
		}
	}

	return partnerDatabases, nil
}
