// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_group

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerinstance/2025-09-01/containerinstance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceContainerGroupUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Containers.ContainerInstanceClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := containerinstance.ParseContainerGroupID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.ContainerGroupsGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("reading %s: %v", id, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("reading %s: `model` was nil", id)
	}

	model := *existing.Model

	if d.HasChange("identity") {
		expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		model.Identity = expandedIdentity

		// As API doesn't return the value of StorageAccountKey, so it has to get the value from tf config and set it to request payload. Otherwise, the Update API call would fail
		addedEmptyDirs := map[string]bool{}
		_, initContainerVolumes, err := expandContainerGroupInitContainers(d, addedEmptyDirs)
		if err != nil {
			return err
		}
		_, _, containerVolumes, err := expandContainerGroupContainers(d, addedEmptyDirs)
		if err != nil {
			return err
		}
		var containerGroupVolumes []containerinstance.Volume
		if initContainerVolumes != nil {
			containerGroupVolumes = initContainerVolumes
		}
		if containerGroupVolumes != nil {
			containerGroupVolumes = append(containerGroupVolumes, containerVolumes...)
		}
		model.Properties.Volumes = pointer.To(containerGroupVolumes)

		// As API doesn't return the value of WorkspaceKey, so it has to get the value from tf config and set it to request payload. Otherwise, the Update API call would fail
		if diagnostics := expandContainerGroupDiagnostics(d.Get("diagnostics").([]any)); diagnostics != nil && diagnostics.LogAnalytics != nil {
			model.Properties.Diagnostics.LogAnalytics.WorkspaceKey = diagnostics.LogAnalytics.WorkspaceKey
		}

		// As Update API doesn't support to update identity, so it has to use CreateOrUpdate API to update identity
		if err := client.ContainerGroupsCreateOrUpdateThenPoll(ctx, *id, model); err != nil {
			return fmt.Errorf("updating %s: %+v", *id, err)
		}
	}

	if d.HasChange("tags") {
		updateParameters := containerinstance.ContainerGroupUpdate{
			Tags: tags.Expand(d.Get("tags").(map[string]any)),
		}

		// As CreateOrUpdate API doesn't support to update tags, so it has to use Update API to update tags
		if _, err := client.ContainerGroupsUpdate(ctx, *id, updateParameters); err != nil {
			return fmt.Errorf("updating tags %s: %+v", *id, err)
		}
	}

	return resourceContainerGroupRead(d, meta)
}
