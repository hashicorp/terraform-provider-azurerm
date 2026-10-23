// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthcare_dicom_service

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2024-03-31/dicomservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceHealthcareApisDicomServiceUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceDicomServiceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dicomservices.ParseDicomServiceID(d.Id())
	if err != nil {
		return err
	}

	// Retrieve the existing resource
	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("cors") {
		payload.Properties.CorsConfiguration = expandDicomServiceCorsConfiguration(d.Get("cors").([]any))
	}

	if d.HasChange("encryption_key_url") {
		payload.Properties.Encryption = &dicomservices.Encryption{
			CustomerManagedKeyEncryption: &dicomservices.EncryptionCustomerManagedKeyEncryption{
				KeyEncryptionKeyURL: pointer.To(d.Get("encryption_key_url").(string)),
			},
		}
	}

	if d.HasChange("public_network_access_enabled") {
		if enabled := d.Get("public_network_access_enabled").(bool); !enabled {
			payload.Properties.PublicNetworkAccess = pointer.To(dicomservices.PublicNetworkAccessDisabled)
		} else {
			payload.Properties.PublicNetworkAccess = pointer.To(dicomservices.PublicNetworkAccessEnabled)
		}
	}

	if d.HasChange("identity") {
		i, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		payload.Identity = i
	}

	if d.HasChange("location") {
		payload.Location = pointer.To(location.Normalize(d.Get("location").(string)))
	}

	if d.HasChange("tags") {
		if err := updateTags(d, meta); err != nil {
			return fmt.Errorf("updating tags error: %+v", err)
		}
	}

	if err = client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceHealthcareApisDicomServiceRead(d, meta)
}
