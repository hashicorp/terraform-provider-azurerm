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
)

// FindDatabaseReplicationPartners looks for partner databases having one of the specified replication roles, by
// reading any replication links then attempting to discover and match the corresponding server/database resources for
// the other end of the link.
func FindDatabaseReplicationPartners(ctx context.Context, databasesClient *databases.DatabasesClient, replicationLinksClient *replicationlinks.ReplicationLinksClient, resourcesClient *resources.Client, id commonids.SqlDatabaseId, rolesToFind []replicationlinks.ReplicationRole) ([]databases.Database, error) {
	var partnerDatabases []databases.Database
	var linkProps *replicationlinks.ReplicationLinkProperties

	matchesRole := func(role replicationlinks.ReplicationRole) bool {
		for _, r := range rolesToFind {
			if r == role {
				return true
			}
		}
		return false
	}

	log.Printf("[INFO] FindDatabaseReplicationPartners lookup for %s", id)

	// Bug 2805551 ReplicationLink API ListByDatabase missed subsubcriptionId in partnerDatabaseId in response body
	results, err := replicationLinksClient.ListByDatabaseComplete(ctx, id)
	if err != nil {
		// Not sure if this is really an error anymore given the change in the API...
		if strings.Contains(err.Error(), "ResourceNotFound") {
			log.Printf("[INFO] %s returned no result records, skipping lookup for Replication Links", id)
			return partnerDatabases, nil
		}
		return nil, fmt.Errorf("reading Replication Links for %s: %+v", id, err)
	}

	// loop over all results that matches this DatabaseId...
	for _, v := range results.Items {
		if linkProps = v.Properties; linkProps != nil {
			if linkProps.PartnerLocation == nil || linkProps.PartnerServer == nil || linkProps.PartnerDatabase == nil {
				log.Printf("[INFO] Replication Link Properties were invalid for %s", id)
				continue
			}

			log.Printf("[INFO] Replication Link found for %s", id)
			var partnerTargetDatabaseId *commonids.SqlDatabaseId

			if linkProps.PartnerDatabaseId != nil {
				partnerTargetDatabaseId, err = commonids.ParseSqlDatabaseIDInsensitively(pointer.From(linkProps.PartnerDatabaseId))
				if err != nil {
					return nil, fmt.Errorf("parsing Partner SQL Server ID %s: %+v", pointer.From(linkProps.PartnerDatabaseId), err)
				}

				log.Printf("[INFO] Partner Database ID from ReplicationLinkProperties: %q", pointer.From(partnerTargetDatabaseId))

				// Check if like-named server has a database named like the partner database, also with a replication link
				partnerResults, err := replicationLinksClient.ListByDatabaseComplete(ctx, pointer.From(partnerTargetDatabaseId))
				if err != nil {
					// Not sure if this is really an error anymore given the change in the API...
					if strings.Contains(err.Error(), "ResourceNotFound") || len(partnerResults.Items) == 0 {
						log.Printf("[INFO] %s returned no result records, skipping lookup of Partner Replication Links", partnerTargetDatabaseId)
						continue
					}
					return nil, fmt.Errorf("reading Partner Replication Links for Database %s: %+v", partnerTargetDatabaseId, err)
				}

				for _, partnerResult := range partnerResults.Items {
					if partnerProps := partnerResult.Properties; partnerProps != nil {
						// If the database has a replication link for the specified role, we'll consider
						// it a partner of this database if the location is the same as the expected partner
						if matchesRole(pointer.From(partnerProps.Role)) {
							partnerDatabase, err := databasesClient.Get(ctx, pointer.From(partnerTargetDatabaseId), databases.DefaultGetOperationOptions())
							if err != nil {
								return nil, fmt.Errorf("retrieving Partner %s: %+v", pointer.From(partnerTargetDatabaseId), err)
							}

							if model := partnerDatabase.Model; model != nil {
								if location.Normalize(model.Location) != location.NormalizeNilable(linkProps.PartnerLocation) {
									log.Printf("[INFO] Mismatch of possible Partner Database based on location (%q vs %q) for %s", location.Normalize(model.Location), location.NormalizeNilable(linkProps.PartnerLocation), pointer.From(model.Id))
									continue
								}

								log.Printf("[INFO] Found Partner %s", pointer.From(partnerTargetDatabaseId))
								partnerDatabases = append(partnerDatabases, pointer.From(model))
							} else {
								log.Printf("[INFO] Partner Database %s: Model is nil", pointer.From(partnerTargetDatabaseId))
								continue
							}
						} else {
							log.Printf("[INFO] Role Mismatch: wanted %+v, got %q", rolesToFind, pointer.From(partnerProps.Role))
						}
					} else {
						log.Printf("[INFO] Partner Replication Link Properties was nil for %s", pointer.From(partnerTargetDatabaseId))
						continue
					}
				}
			}
		} else {
			log.Printf("[INFO] Replication Link Properties was nil for %s", id)
			continue
		}
	}

	log.Printf("[INFO] %d Replication Link(s) found for %s", len(partnerDatabases), id)
	return partnerDatabases, nil
}
