package compute

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-07-01/compute"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceManagedDiskExport() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceManagedDiskExportRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"managed_disk_id": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     azure.ValidateResourceID,
			},

			"duration_in_seconds": {
				Type:     pluginsdk.TypeInt,
				Required: true,
			},

			"access": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "Read",
			},

			"sas": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceManagedDiskExportRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	managedDiskId := d.Get("managed_disk_id").(string)
	durationInSeconds := int32(d.Get("duration_in_seconds").(int))
	access := d.Get("access").(string)

	parsedManagedDiskId, err := parse.ManagedDiskID(managedDiskId)
	if err != nil {
		return fmt.Errorf("parsing Managed Disk ID %q: %+v", parsedManagedDiskId.ID(), err)
	}

	diskName := parsedManagedDiskId.DiskName
	resourceGroupName := parsedManagedDiskId.ResourceGroup

	grantAccessData := compute.GrantAccessData{
		Access:            compute.AccessLevel(access),
		DurationInSeconds: &durationInSeconds,
	}

	// Request to Grant Access for Disk
	diskGrantFuture, err := client.GrantAccess(ctx, resourceGroupName, diskName, grantAccessData)
	if err != nil {
		return fmt.Errorf("Error while granting access for disk export %q: %+v", parsedManagedDiskId.ID(), err)
	}

	// Wait until the Grant Request is complete
	err = diskGrantFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Grant access operation failed %q (Resource Group %q): %+v", diskName, resourceGroupName, err)
	}

	// Fetch the SAS token from the response
	diskGrantResponse, err := diskGrantFuture.Result(*client)
	if err != nil {
		return fmt.Errorf("Error while fetching the response %q: %+v", parsedManagedDiskId.ID(), err)
	}

	sasToken := diskGrantResponse.AccessSAS

	d.Set("sas", *sasToken)
	tokenHash := sha256.Sum256([]byte(*sasToken))
	d.SetId(hex.EncodeToString(tokenHash[:]))

	return nil

}
