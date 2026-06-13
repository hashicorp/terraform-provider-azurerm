// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotecteditems"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &SiteRecoveryReplicatedVMDiskTypePoller{}

const siteRecoveryReplicatedVMDiskTypeRequiredSuccessfulReads = 10

type SiteRecoveryReplicatedVMDiskTypeUpdate struct {
	DiskId                  string
	RecoveryTargetDiskType  string
	RecoveryReplicaDiskType string
}

type SiteRecoveryReplicatedVMDiskTypePoller struct {
	client                     *replicationprotecteditems.ReplicationProtectedItemsClient
	id                         replicationprotecteditems.ReplicationProtectedItemId
	updates                    map[string]SiteRecoveryReplicatedVMDiskTypeUpdate
	consecutiveSuccessfulReads int
}

func NewSiteRecoveryReplicatedVMDiskTypePoller(client *replicationprotecteditems.ReplicationProtectedItemsClient, id replicationprotecteditems.ReplicationProtectedItemId, updates []SiteRecoveryReplicatedVMDiskTypeUpdate) *SiteRecoveryReplicatedVMDiskTypePoller {
	updatesByDiskId := make(map[string]SiteRecoveryReplicatedVMDiskTypeUpdate, len(updates))
	for _, update := range updates {
		updatesByDiskId[normalizeSiteRecoveryReplicatedVMManagedDiskID(update.DiskId)] = update
	}

	return &SiteRecoveryReplicatedVMDiskTypePoller{
		client:  client,
		id:      id,
		updates: updatesByDiskId,
	}
}

func (p *SiteRecoveryReplicatedVMDiskTypePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	if resp.Model == nil {
		return nil, fmt.Errorf("retrieving %s: `model` was nil", p.id)
	}

	if resp.Model.Properties == nil {
		return nil, fmt.Errorf("retrieving %s: `properties` was nil", p.id)
	}

	result := &pollers.PollResult{
		HttpResponse: &client.Response{
			OData:    resp.OData,
			Response: resp.HttpResponse,
		},
		PollInterval: 30 * time.Second,
	}

	if siteRecoveryReplicatedVMStateIndicatesFailure(pointer.From(resp.Model.Properties.ProtectionState)) {
		return nil, pollers.PollingFailedError{
			HttpResponse: result.HttpResponse,
			Message:      fmt.Sprintf("%s entered protection state %q", p.id, pointer.From(resp.Model.Properties.ProtectionState)),
		}
	}

	a2aDetails, ok := resp.Model.Properties.ProviderSpecificDetails.(replicationprotecteditems.A2AReplicationDetails)
	if !ok {
		return nil, pollers.PollingFailedError{
			HttpResponse: result.HttpResponse,
			Message:      fmt.Sprintf("%s returned provider specific details that were not A2A", p.id),
		}
	}

	if siteRecoveryReplicatedVMStateIndicatesFailure(pointer.From(a2aDetails.VMProtectionState)) {
		return nil, pollers.PollingFailedError{
			HttpResponse: result.HttpResponse,
			Message:      fmt.Sprintf("%s entered VM protection state %q", p.id, pointer.From(a2aDetails.VMProtectionState)),
		}
	}

	pending, err := p.pendingDiskTypeUpdates(a2aDetails, result.HttpResponse)
	if err != nil {
		return nil, err
	}

	if len(pending) == 0 {
		p.consecutiveSuccessfulReads++
		if p.consecutiveSuccessfulReads < siteRecoveryReplicatedVMDiskTypeRequiredSuccessfulReads {
			log.Printf("[DEBUG] waiting for managed disk type updates for %s: matched %d/%d consecutive reads", p.id, p.consecutiveSuccessfulReads, siteRecoveryReplicatedVMDiskTypeRequiredSuccessfulReads)
			result.Status = pollers.PollingStatusInProgress
			return result, nil
		}

		result.Status = pollers.PollingStatusSucceeded
		return result, nil
	}

	p.consecutiveSuccessfulReads = 0
	sort.Strings(pending)
	log.Printf("[DEBUG] waiting for managed disk type updates for %s: %s", p.id, strings.Join(pending, "; "))
	result.Status = pollers.PollingStatusInProgress
	return result, nil
}

func (p *SiteRecoveryReplicatedVMDiskTypePoller) pendingDiskTypeUpdates(details replicationprotecteditems.A2AReplicationDetails, resp *client.Response) ([]string, error) {
	pending := make([]string, 0)

	protectedDisks := make(map[string]replicationprotecteditems.A2AProtectedManagedDiskDetails)
	if details.ProtectedManagedDisks != nil {
		for _, disk := range *details.ProtectedManagedDisks {
			diskId := normalizeSiteRecoveryReplicatedVMManagedDiskID(pointer.From(disk.DiskId))
			if _, ok := p.updates[diskId]; ok {
				protectedDisks[diskId] = disk
			}
		}
	}

	for diskId, update := range p.updates {
		disk, ok := protectedDisks[diskId]
		if !ok {
			pending = append(pending, fmt.Sprintf("%s is not present in protected managed disks", update.DiskId))
			continue
		}

		if siteRecoveryReplicatedVMStateIndicatesFailure(pointer.From(disk.DiskState)) {
			return nil, pollers.PollingFailedError{
				HttpResponse: resp,
				Message:      fmt.Sprintf("%s entered disk state %q", update.DiskId, pointer.From(disk.DiskState)),
			}
		}

		actualTargetDiskType := pointer.From(disk.RecoveryTargetDiskAccountType)
		actualTargetReplicaDiskType := pointer.From(disk.RecoveryReplicaDiskAccountType)

		mismatches := make([]string, 0)
		if actualTargetDiskType != update.RecoveryTargetDiskType {
			mismatches = append(mismatches, fmt.Sprintf("target_disk_type is %q, expected %q", actualTargetDiskType, update.RecoveryTargetDiskType))
		}
		if actualTargetReplicaDiskType != update.RecoveryReplicaDiskType {
			mismatches = append(mismatches, fmt.Sprintf("target_replica_disk_type is %q, expected %q", actualTargetReplicaDiskType, update.RecoveryReplicaDiskType))
		}
		if len(mismatches) > 0 {
			pending = append(pending, fmt.Sprintf("%s: %s", update.DiskId, strings.Join(mismatches, ", ")))
		}
	}

	return pending, nil
}

func normalizeSiteRecoveryReplicatedVMManagedDiskID(input string) string {
	if parsed, err := commonids.ParseManagedDiskIDInsensitively(input); err == nil {
		return strings.ToLower(parsed.ID())
	}

	return strings.ToLower(input)
}

func siteRecoveryReplicatedVMStateIndicatesFailure(input string) bool {
	input = strings.ToLower(input)
	if input == "" || input == "noerror" {
		return false
	}

	return strings.Contains(input, "failed") || strings.Contains(input, "error")
}
