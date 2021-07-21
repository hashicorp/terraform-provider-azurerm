package videoanalyzer

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/videoanalyzer/mgmt/2021-05-01-preview/videoanalyzer"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/azurerm/utils"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	msiparse "github.com/hashicorp/terraform-provider-azurerm/internal/services/msi/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/msi/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/videoanalyzer/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVideoAnalyzer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVideoAnalyzerCreateUpdate,
		Read:   resourceVideoAnalyzerRead,
		Update: resourceVideoAnalyzerCreateUpdate,
		Delete: resourceVideoAnalyzerDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.VideoAnalyzerID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-z0-9]{3,24}$"),
					"Video Analyzer name must be 3 - 24 characters long, contain only lowercase letters and numbers.",
				),
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"storage_account": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},

						"identity_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.UserAssignedIdentityID,
						},
					},
				},
			},

			//lintignore:XS003
			"identity": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string("UserAssigned"),
							}, true),
						},
						"identity_ids": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							MinItems: 1,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validate.UserAssignedIdentityID,
							},
						},
					},
				},
			},
			"tags": tags.Schema(),
		},
	}
}

func resourceVideoAnalyzerCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).VideoAnalyzer.VideoAnalyzersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId := parse.NewVideoAnalyzerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceId.ResourceGroup, resourceId.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", resourceId, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_video_analyzer", resourceId.ID())
		}
	}

	identity, err := expandAzureRmVideoAnalyzerIdentity(d)
	if err != nil {
		return err
	}
	parameters := videoanalyzer.Model{
		PropertiesType: &videoanalyzer.PropertiesType{
			StorageAccounts: expandVideoAnalyzerStorageAccounts(d),
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Identity: identity,
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.CreateOrUpdate(ctx, resourceId.ResourceGroup, resourceId.Name, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", resourceId, err)
	}

	d.SetId(resourceId.ID())
	return resourceVideoAnalyzerRead(d, meta)
}

func resourceVideoAnalyzerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).VideoAnalyzer.VideoAnalyzersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VideoAnalyzerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Video Analyzer %q was not found in Resource Group %q - removing from state", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Video Analyzer %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	props := resp.PropertiesType
	if props != nil {
		accounts := flattenVideoAnalyzerStorageAccounts(props.StorageAccounts)
		if err := d.Set("storage_account", accounts); err != nil {
			return fmt.Errorf("flattening `storage_account`: %s", err)
		}
	}

	flattenedIdentity, err := flattenAzureRmVideoServiceIdentity(resp.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %s", err)
	}

	if err := d.Set("identity", flattenedIdentity); err != nil {
		return fmt.Errorf("flattening `identity`: %s", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceVideoAnalyzerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).VideoAnalyzer.VideoAnalyzersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VideoAnalyzerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if response.WasNotFound(resp.Response) {
			return nil
		}
		return fmt.Errorf("issuing AzureRM delete request for Video Analyzer '%s': %+v", id.Name, err)
	}

	return nil
}

func expandVideoAnalyzerStorageAccounts(d *pluginsdk.ResourceData) *[]videoanalyzer.StorageAccount {
	storageAccountRaw := d.Get("storage_account").([]interface{})[0].(map[string]interface{})

	results := []videoanalyzer.StorageAccount{
		{
			ID: utils.String(storageAccountRaw["id"].(string)),
			Identity: &videoanalyzer.ResourceIdentity{
				UserAssignedIdentity: utils.String(storageAccountRaw["identity_id"].(string)),
			},
		},
	}

	return &results
}

func flattenVideoAnalyzerStorageAccounts(input *[]videoanalyzer.StorageAccount) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, storageAccount := range *input {
		output := make(map[string]interface{})

		if storageAccount.ID != nil {
			output["id"] = *storageAccount.ID
		}

		if storageAccount.Identity != nil {
			output["identity_id"] = *storageAccount.Identity.UserAssignedIdentity
		}

		results = append(results, output)
	}

	return results
}

func expandAzureRmVideoAnalyzerIdentity(d *pluginsdk.ResourceData) (*videoanalyzer.Identity, error) {
	identityRaw := d.Get("identity").([]interface{})
	if identityRaw[0] == nil {
		return nil, nil
	}
	identity := identityRaw[0].(map[string]interface{})
	result := &videoanalyzer.Identity{
		Type: utils.String(identity["type"].(string)),
	}
	var identityIdSet []interface{}
	if identityIds, exists := identity["identity_ids"]; exists {
		identityIdSet = identityIds.(*pluginsdk.Set).List()
	}

	// If type contains `UserAssigned`, `identity_ids` must be specified and have at least 1 element
	if *result.Type == "UserAssigned" {
		if len(identityIdSet) == 0 {
			return nil, fmt.Errorf("`identity_ids` must have at least 1 element when `type` includes `UserAssigned`")
		}

		userAssignedIdentities := make(map[string]*videoanalyzer.UserAssignedManagedIdentity)
		for _, id := range identityIdSet {
			userAssignedIdentities[id.(string)] = &videoanalyzer.UserAssignedManagedIdentity{}
		}
		result.UserAssignedIdentities = userAssignedIdentities
	} else if len(identityIdSet) > 0 {
		// If type does _not_ contain `UserAssigned` (i.e. is set to `SystemAssigned` or defaulted to `None`), `identity_ids` is not allowed
		return nil, fmt.Errorf("`identity_ids` can only be specified when `type` includes `UserAssigned`; but `type` is currently %q", *result.Type)
	}

	return result, nil
}

func flattenAzureRmVideoServiceIdentity(identity *videoanalyzer.Identity) ([]interface{}, error) {
	if identity == nil {
		return make([]interface{}, 0), nil
	}

	result := make(map[string]interface{})
	result["type"] = identity.Type

	identityIds := make([]interface{}, 0)
	if identity.UserAssignedIdentities != nil {
		/*
		   "userAssignedIdentities": {
		     "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/id1": {
		       "principalId": "00000000-0000-0000-0000-000000000000",
		       "clientId": "00000000-0000-0000-0000-000000000000"
		     },
		   }
		*/
		for key := range identity.UserAssignedIdentities {
			parsedId, err := msiparse.UserAssignedIdentityID(key)
			if err != nil {
				return nil, err
			}
			identityIds = append(identityIds, parsedId.ID())
		}
		result["identity_ids"] = identityIds
	}

	return []interface{}{result}, nil
}
