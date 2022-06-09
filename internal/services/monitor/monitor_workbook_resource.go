package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	workbook "github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/sdk/2022-04-01/insights"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	tagsHelper "github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMonitorWorkbook() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMonitorWorkbookCreateOrUpdate,
		Read:   resourceMonitorWorkbookRead,
		Update: resourceMonitorWorkbookCreateOrUpdate,
		Delete: resourceMonitorWorkbookDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := workbook.ParseWorkbookID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validation.IsUUID,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"display_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"serialized_data": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
			},

			"source_id": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validation.StringIsNotEmpty,
			},

			"category": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      "workbook",
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptionalForceNew(),

			"storage_uri": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				RequiredWith: []string{
					"identity",
				},
			},

			"tags": tags.Schema(),
		},

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v interface{}) error {
			// d.GetOk cannot identify whether the property value is set by the user or returned by the service backend, use GetRawConfig instead
			if t := diff.GetRawConfig().AsValueMap()["tags"]; !t.IsNull() {
				if value := t.AsValueMap()["hidden-title"]; !value.IsNull() {
					return fmt.Errorf("a tag with the key `hidden-title` should not be used to set the display name. Please Use `display_name` instead")
				}
			}

			return nil
		}),
	}
}

func resourceMonitorWorkbookCreateOrUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.WorkbookClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := workbook.NewWorkbookID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.WorkbooksGet(ctx, id, workbook.WorkbooksGetOperationOptions{CanFetchContent: utils.Bool(true)})
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_monitor_workbook", id.ID())
		}
	}

	kindValue := workbook.WorkbookSharedTypeKindShared
	props := workbook.Workbook{
		Kind:     &kindValue,
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Properties: &workbook.WorkbookProperties{
			Category:       d.Get("category").(string),
			Description:    utils.String(d.Get("description").(string)),
			DisplayName:    d.Get("display_name").(string),
			SerializedData: d.Get("serialized_data").(string),
		},

		Tags: tagsHelper.Expand(d.Get("tags").(map[string]interface{})),
	}

	if storageURI, ok := d.GetOk("storage_uri"); ok {
		props.Properties.StorageUri = utils.String(storageURI.(string))
	}

	identity, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	props.Identity = identity

	if _, err := client.WorkbooksCreateOrUpdate(ctx, id, props, workbook.WorkbooksCreateOrUpdateOperationOptions{SourceId: utils.String(d.Get("source_id").(string))}); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceMonitorWorkbookRead(d, meta)
}

func resourceMonitorWorkbookRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.WorkbookClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workbook.ParseWorkbookID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.WorkbooksGet(ctx, *id, workbook.WorkbooksGetOperationOptions{CanFetchContent: utils.Bool(true)})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ResourceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(*model.Location))
		identity, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return err
		}

		if err = d.Set("identity", identity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if props := model.Properties; props != nil {
			d.Set("category", props.Category)
			d.Set("description", props.Description)
			d.Set("display_name", props.DisplayName)
			d.Set("serialized_data", props.SerializedData)
			d.Set("source_id", props.SourceId)
			d.Set("storage_uri", props.StorageUri)
		}

		// The backend returns a tags with key `hidden-title` by default. Since it has the same value with `display_name` and will cause inconsistency with user's configuration, remove it as a workaround.
		if model.Tags != nil {
			if _, ok := (*model.Tags)["hidden-title"]; ok {
				delete(*model.Tags, "hidden-title")
			}
		}

		if err = tagsHelper.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceMonitorWorkbookDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.WorkbookClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workbook.ParseWorkbookID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.WorkbooksDelete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
