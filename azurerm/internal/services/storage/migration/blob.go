package migration

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func BlobV0ToV1() schema.StateUpgrader {
	return schema.StateUpgrader{
		// this should have been applied from pre-0.12 migration system; backporting just in-case
		Type:    blobSchemaForV0().CoreConfigSchema().ImpliedType(),
		Upgrade: blobUpgradeV0ToV1,
		Version: 0,
	}
}

func blobSchemaForV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"storage_account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"storage_container_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  0,
			},
			"content_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "application/octet-stream",
				ConflictsWith: []string{"source_uri"},
			},
			"source": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"source_uri"},
			},
			"source_uri": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"source"},
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parallelism": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  8,
				ForceNew: true,
			},
			"attempts": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
				ForceNew: true,
			},
		},
	}
}

func blobUpgradeV0ToV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	environment := meta.(*clients.Client).Account.Environment

	blobName := rawState["name"]
	containerName := rawState["storage_container_name"]
	storageAccountName := rawState["storage_account_name"]
	newID := fmt.Sprintf("https://%s.blob.%s/%s/%s", storageAccountName, environment.StorageEndpointSuffix, containerName, blobName)
	rawState["id"] = newID

	return rawState, nil
}
