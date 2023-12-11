// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/gallerysharingupdate"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSharedImageGallery() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSharedImageGalleryCreate,
		Read:   resourceSharedImageGalleryRead,
		Update: resourceSharedImageGalleryUpdate,
		Delete: resourceSharedImageGalleryDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseSharedImageGalleryID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SharedImageGalleryName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"sharing": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"permission": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(galleries.GallerySharingPermissionTypesCommunity),
								string(galleries.GallerySharingPermissionTypesGroups),
								string(galleries.GallerySharingPermissionTypesPrivate),
							}, false),
						},

						"community_gallery": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"eula": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"prefix": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validate.SharedImageGalleryPrefix,
									},
									"publisher_email": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"publisher_uri": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"tags": commonschema.Tags(),

			"unique_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSharedImageGalleryCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleriesClient
	gallerySharingUpdateClient := meta.(*clients.Client).Compute.GallerySharingUpdateClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewSharedImageGalleryID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id, galleries.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_shared_image_gallery", id.ID())
	}

	sharing, permission, err := expandSharedImageGallerySharing(d.Get("sharing").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `sharing`: %+v", err)
	}

	payload := galleries.Gallery{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &galleries.GalleryProperties{
			Description:    pointer.To(d.Get("description").(string)),
			SharingProfile: sharing,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if permission == galleries.GallerySharingPermissionTypesCommunity {
		updatePayload := gallerysharingupdate.SharingUpdate{
			OperationType: gallerysharingupdate.SharingUpdateOperationTypesEnableCommunity,
		}
		if err = gallerySharingUpdateClient.GallerySharingProfileUpdateThenPoll(ctx, id, updatePayload); err != nil {
			return fmt.Errorf("enabling community sharing of %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())

	return resourceSharedImageGalleryRead(d, meta)
}

func resourceSharedImageGalleryRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleriesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSharedImageGalleryID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, galleries.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.GalleryName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			d.Set("description", props.Description)

			uniqueName := ""
			if props.Identifier != nil && props.Identifier.UniqueName != nil {
				uniqueName = *props.Identifier.UniqueName
			}
			d.Set("unique_name", uniqueName)

			d.Set("sharing", flattenSharedImageGallerySharing(props.SharingProfile))
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}

func resourceSharedImageGalleryUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleriesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSharedImageGalleryID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, galleries.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	payload := resp.Model
	if payload == nil {
		payload = &galleries.Gallery{}
	}

	if payload.Properties == nil {
		payload.Properties = &galleries.GalleryProperties{}
	}

	if d.HasChange("description") {
		payload.Properties.Description = pointer.To(d.Get("description").(string))
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}
	return resourceSharedImageGalleryRead(d, meta)
}

func resourceSharedImageGalleryDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleriesClient
	gallerySharingUpdateClient := meta.(*clients.Client).Compute.GallerySharingUpdateClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSharedImageGalleryID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, galleries.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if model := resp.Model; model != nil {
		if prop := model.Properties; prop != nil && prop.SharingProfile != nil && prop.SharingProfile.Permissions != nil {
			if pointer.From(prop.SharingProfile.Permissions) == galleries.GallerySharingPermissionTypesCommunity {
				updatePayload := gallerysharingupdate.SharingUpdate{
					OperationType: gallerysharingupdate.SharingUpdateOperationTypesReset,
				}
				if err = gallerySharingUpdateClient.GallerySharingProfileUpdateThenPoll(ctx, *id, updatePayload); err != nil {
					return fmt.Errorf("reseting community sharing of %s: %+v", id, err)
				}
			}
		}
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandSharedImageGallerySharing(input []interface{}) (*galleries.SharingProfile, galleries.GallerySharingPermissionTypes, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, "", nil
	}

	v := input[0].(map[string]interface{})
	permission := galleries.GallerySharingPermissionTypes(v["permission"].(string))
	communityGallery := v["community_gallery"].([]interface{})

	if permission == galleries.GallerySharingPermissionTypesCommunity {
		if len(communityGallery) == 0 || communityGallery[0] == nil {
			return nil, permission, fmt.Errorf("`community_gallery` must be set when `permission` is set to `Community`")
		}
	}

	return &galleries.SharingProfile{
		Permissions:          pointer.To(permission),
		CommunityGalleryInfo: expandSharedImageGalleryCommunityGallery(communityGallery),
	}, permission, nil
}

func flattenSharedImageGallerySharing(input *galleries.SharingProfile) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	permission := ""
	if v := input.Permissions; v != nil {
		permission = string(pointer.From(v))
	}

	return []interface{}{
		map[string]interface{}{
			"permission":        permission,
			"community_gallery": flattenSharedImageGalleryCommunityGallery(input.CommunityGalleryInfo),
		},
	}
}

func expandSharedImageGalleryCommunityGallery(input []interface{}) *galleries.CommunityGalleryInfo {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &galleries.CommunityGalleryInfo{
		Eula:             pointer.To(v["eula"].(string)),
		PublicNamePrefix: pointer.To(v["prefix"].(string)),
		PublisherContact: pointer.To(v["publisher_email"].(string)),
		PublisherUri:     pointer.To(v["publisher_uri"].(string)),
	}
}

func flattenSharedImageGalleryCommunityGallery(input *galleries.CommunityGalleryInfo) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	eula := ""
	if input.Eula != nil {
		eula = pointer.From(input.Eula)
	}

	publicName := ""
	if input.PublicNames != nil {
		if v := pointer.From(input.PublicNames); len(v) > 0 {
			publicName = v[0]
		}
	}

	publicNamePrefix := ""
	if input.PublicNamePrefix != nil {
		publicNamePrefix = pointer.From(input.PublicNamePrefix)
	}

	publisherEmail := ""
	if input.PublisherContact != nil {
		publisherEmail = pointer.From(input.PublisherContact)
	}

	publisherUri := ""
	if input.PublisherUri != nil {
		publisherUri = pointer.From(input.PublisherUri)
	}

	return []interface{}{
		map[string]interface{}{
			"eula":            eula,
			"name":            publicName,
			"prefix":          publicNamePrefix,
			"publisher_email": publisherEmail,
			"publisher_uri":   publisherUri,
		},
	}
}
