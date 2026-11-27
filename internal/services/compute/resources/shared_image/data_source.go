// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package shared_image

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryimages"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	computeClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Shared Image within a Shared Image Gallery.",
		Read:        dataSourceSharedImageRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of the Shared Image.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.SharedImageName,
			},

			"gallery_name": {
				Description:  "The name of the Shared Image Gallery in which the Shared Image exists.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.SharedImageGalleryName,
			},

			"location": commonschema.LocationComputed(),

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"architecture": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"os_type": {
				Description: "The type of Operating System present in this Shared Image.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"specialized": {
				Description: "Whether the Operating System used inside this Image has been specialized (for example, `sysprep` on Windows has been run).",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"hyper_v_generation": {
				Description: "The generation of HyperV that the Virtual Machine used to create the Shared Image is based on.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"identifier": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"publisher": {
							Description: "The Publisher Name for this Gallery Image.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"offer": {
							Description: "The Offer Name for this Shared Image.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"sku": {
							Description: "The Name of the SKU for this Gallery Image.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},

			"description": {
				Description: "The description of this Shared Image.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"eula": {
				Description: "The End User Licence Agreement for the Shared Image.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"privacy_statement_uri": {
				Description: "The URI containing the Privacy Statement for this Shared Image.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"purchase_plan": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Description: "The Purchase Plan Name for this Shared Image.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"publisher": {
							Description: "The Purchase Plan Publisher for this Gallery Image.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"product": {
							Description: "The Purchase Plan Product for this Gallery Image.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},

			"release_note_uri": {
				Description: "The URI containing the Release Notes for this Shared Image.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"trusted_launch_supported": {
				Description: "Whether supports creation of both Trusted Launch Virtual Machines and Gen2 Virtual Machines with standard security created from the Shared Image.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"trusted_launch_enabled": {
				Description: "Whether Trusted Launch has to be enabled for the Virtual Machine created from the Shared Image.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"confidential_vm_supported": {
				Description: "Whether supports creation of both Confidential Virtual Machines and Gen2 Virtual Machines with standard security from a compatible Gen2 OS disk VHD or Gen2 Managed image.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"confidential_vm_enabled": {
				Description: "Whether Confidential Virtual Machines enabled. It will enable all the features of trusted, with higher confidentiality features for isolate machines or encrypted data. Available for Gen2 machines.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"accelerated_network_support_enabled": {
				Description: "Whether the Shared Image supports Accelerated Network.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"hibernation_enabled": {
				Description: "Whether the Shared Image supports hibernation.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceSharedImageRead(d *pluginsdk.ResourceData, meta any) error {
	client := computeClient.From(meta).GalleryImagesClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id := galleryimages.NewGalleryImageID(subscriptionId, d.Get("resource_group_name").(string), d.Get("gallery_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.ImageName)
	d.Set("gallery_name", id.GalleryName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			d.Set("description", props.Description)
			d.Set("eula", props.Eula)
			d.Set("os_type", string(props.OsType))
			d.Set("architecture", pointer.From(props.Architecture))
			d.Set("specialized", props.OsState == galleryimages.OperatingSystemStateTypesSpecialized)
			d.Set("hyper_v_generation", pointer.From(props.HyperVGeneration))
			d.Set("privacy_statement_uri", props.PrivacyStatementUri)
			d.Set("release_note_uri", props.ReleaseNoteUri)

			if err := d.Set("identifier", flattenGalleryImageDataSourceIdentifier(&props.Identifier)); err != nil {
				return fmt.Errorf("setting `identifier`: %+v", err)
			}

			if err := d.Set("purchase_plan", flattenGalleryImageDataSourcePurchasePlan(props.PurchasePlan)); err != nil {
				return fmt.Errorf("setting `purchase_plan`: %+v", err)
			}

			trustedLaunchSupported := false
			trustedLaunchEnabled := false
			cvmEnabled := false
			cvmSupported := false
			acceleratedNetworkSupportEnabled := false
			hibernationEnabled := false
			if model.Properties.Features != nil {
				for _, feature := range *model.Properties.Features {
					if feature.Name == nil || feature.Value == nil {
						continue
					}

					if strings.EqualFold(*feature.Name, "SecurityType") {
						trustedLaunchSupported = strings.EqualFold(*feature.Value, "TrustedLaunchSupported")
						trustedLaunchEnabled = strings.EqualFold(*feature.Value, "TrustedLaunch")
						cvmSupported = strings.EqualFold(*feature.Value, "ConfidentialVmSupported")
						cvmEnabled = strings.EqualFold(*feature.Value, "ConfidentialVm")
					}

					if strings.EqualFold(*feature.Name, "IsAcceleratedNetworkSupported") {
						acceleratedNetworkSupportEnabled = strings.EqualFold(*feature.Value, "true")
					}

					if strings.EqualFold(*feature.Name, "IsHibernateSupported") {
						hibernationEnabled = strings.EqualFold(*feature.Value, "true")
					}
				}
			}

			d.Set("confidential_vm_supported", cvmSupported)
			d.Set("confidential_vm_enabled", cvmEnabled)
			d.Set("trusted_launch_supported", trustedLaunchSupported)
			d.Set("trusted_launch_enabled", trustedLaunchEnabled)
			d.Set("accelerated_network_support_enabled", acceleratedNetworkSupportEnabled)
			d.Set("hibernation_enabled", hibernationEnabled)
		}

		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func flattenGalleryImageDataSourceIdentifier(input *galleryimages.GalleryImageIdentifier) []any {
	if input == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"offer":     input.Offer,
			"publisher": input.Publisher,
			"sku":       input.Sku,
		},
	}
}

func flattenGalleryImageDataSourcePurchasePlan(input *galleryimages.ImagePurchasePlan) []any {
	if input == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"name":      pointer.From(input.Name),
			"publisher": pointer.From(input.Publisher),
			"product":   pointer.From(input.Product),
		},
	}
}
