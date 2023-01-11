package healthcare

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/healthcareapis/mgmt/2021-11-01/healthcareapis" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceHealthcareApisDicomService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceHealthcareApisDicomServiceCreate,
		Read:   resourceHealthcareApisDicomServiceRead,
		Update: resourceHealthcareApisDicomServiceUpdate,
		Delete: resourceHealthcareApisDicomServiceDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.HealthCareDicomServiceV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DicomServiceID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DicomServiceName(),
			},

			"workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WorkspaceID,
			},

			"location": commonschema.Location(),

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"authentication": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"authority": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"audience": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
					},
				},
			},

			"private_endpoint": {
				Type:     pluginsdk.TypeSet,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"service_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceHealthcareApisDicomServiceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceDicomServiceClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM Healthcare Dicom Service creation.")

	workspace, err := parse.WorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return fmt.Errorf("parsing healthcare workspace error: %+v", err)
	}

	dicomServiceId := parse.NewDicomServiceID(workspace.SubscriptionId, workspace.ResourceGroup, workspace.Name, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, dicomServiceId.ResourceGroup, dicomServiceId.WorkspaceName, dicomServiceId.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", dicomServiceId, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_healthcare_dicom_service", dicomServiceId.ID())
		}
	}

	identity, err := expandDicomManagedIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	t := d.Get("tags").(map[string]interface{})

	parameters := healthcareapis.DicomService{
		Identity: identity,
		DicomServiceProperties: &healthcareapis.DicomServiceProperties{
			PublicNetworkAccess: healthcareapis.PublicNetworkAccessEnabled,
		},
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Tags:     tags.Expand(t),
	}
	if enabled := d.Get("public_network_access_enabled").(bool); !enabled {
		parameters.DicomServiceProperties.PublicNetworkAccess = healthcareapis.PublicNetworkAccessDisabled
	}

	future, err := client.CreateOrUpdate(ctx, dicomServiceId.ResourceGroup, dicomServiceId.WorkspaceName, dicomServiceId.Name, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", dicomServiceId, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation %s: %+v", dicomServiceId, err)
	}

	d.SetId(dicomServiceId.ID())

	return resourceHealthcareApisDicomServiceRead(d, meta)
}

func resourceHealthcareApisDicomServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceDicomServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DicomServiceID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Dicom service error: %+v", err)
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	workspaceId := parse.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)
	d.Set("workspace_id", workspaceId.ID())

	if resp.Location != nil {
		d.Set("location", location.Normalize(*resp.Location))
	}

	if props := resp.DicomServiceProperties; props != nil {
		d.Set("authentication", flattenDicomAuthentication(props.AuthenticationConfiguration))
		d.Set("private_endpoint", flattenDicomServicePrivateEndpoint(props.PrivateEndpointConnections))
		d.Set("service_url", props.ServiceURL)

		if props.PublicNetworkAccess != "" {
			d.Set("public_network_access_enabled", props.PublicNetworkAccess == healthcareapis.PublicNetworkAccessEnabled)
		}
	}

	identity, _ := flattenDicomManagedIdentity(resp.Identity)
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if err := tags.FlattenAndSet(d, resp.Tags); err != nil {
		return err
	}
	return nil
}

func resourceHealthcareApisDicomServiceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceDicomServiceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DicomServiceID(d.Id())
	if err != nil {
		return err
	}

	expandedIdentity, err := expandDicomManagedIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	parameters := healthcareapis.DicomService{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		DicomServiceProperties: &healthcareapis.DicomServiceProperties{
			PublicNetworkAccess: healthcareapis.PublicNetworkAccessEnabled,
		},
		Identity: expandedIdentity,
	}

	if enabled := d.Get("public_network_access_enabled").(bool); !enabled {
		parameters.DicomServiceProperties.PublicNetworkAccess = healthcareapis.PublicNetworkAccessDisabled
	}

	if d.HasChange("tags") {
		if err := updateTags(d, meta); err != nil {
			return fmt.Errorf("updating tags error: %+v", err)
		}
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceHealthcareApisDicomServiceRead(d, meta)
}

func resourceHealthcareApisDicomServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceDicomServiceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DicomServiceID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Dicom service error: %+v", err)
	}
	future, err := client.Delete(ctx, id.ResourceGroup, id.Name, id.WorkspaceName)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	log.Printf("[DEBUG] Waiting for %s to be deleted..", id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Pending"},
		Target:                    []string{"Deleted"},
		Refresh:                   dicomServiceStateStatusCodeRefreshFunc(ctx, client, *id),
		Timeout:                   d.Timeout(pluginsdk.TimeoutDelete),
		ContinuousTargetOccurence: 3,
		PollInterval:              10 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}
	return nil
}

func dicomServiceStateStatusCodeRefreshFunc(ctx context.Context, client *healthcareapis.DicomServicesClient, id parse.DicomServiceId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)

		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return res, "Deleted", nil
			}
			return nil, "Error", fmt.Errorf("polling for the status of %s: %+v", id, err)
		}

		return res, "Pending", nil
	}
}

func updateTags(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceDicomServiceClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DicomServiceID(d.Id())
	if err != nil {
		return err
	}

	update := healthcareapis.DicomServicePatchResource{
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.Name, id.WorkspaceName, update)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update %s: %+v", id, err)
	}

	return nil
}

func flattenDicomAuthentication(input *healthcareapis.DicomServiceAuthenticationConfiguration) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	authBlock := make(map[string]interface{})
	if input.Authority != nil {
		authBlock["authority"] = *input.Authority
	}

	audience := make([]interface{}, 0)
	if input.Audiences != nil {
		for _, data := range *input.Audiences {
			audience = append(audience, data)
		}
	}
	authBlock["audience"] = audience

	return []interface{}{authBlock}
}

func flattenDicomServicePrivateEndpoint(input *[]healthcareapis.PrivateEndpointConnection) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, endpoint := range *input {
		result := map[string]interface{}{}
		if endpoint.Name != nil {
			result["name"] = *endpoint.Name
		}

		if endpoint.ID != nil {
			result["id"] = *endpoint.ID
		}
	}
	return results
}

func expandDicomManagedIdentity(input []interface{}) (*healthcareapis.ServiceManagedIdentityIdentity, error) {
	expanded, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	out := healthcareapis.ServiceManagedIdentityIdentity{
		Type: healthcareapis.ServiceManagedIdentityType(string(expanded.Type)),
	}

	if expanded.Type == identity.TypeUserAssigned || expanded.Type == identity.TypeSystemAssignedUserAssigned {
		out.UserAssignedIdentities = make(map[string]*healthcareapis.UserAssignedIdentity)
		for k := range expanded.IdentityIds {
			out.UserAssignedIdentities[k] = &healthcareapis.UserAssignedIdentity{
				// intentionally empty
			}
		}
	}
	return &out, nil
}

func flattenDicomManagedIdentity(input *healthcareapis.ServiceManagedIdentityIdentity) (*[]interface{}, error) {
	var config *identity.SystemAndUserAssignedList

	if input != nil {
		config = &identity.SystemAndUserAssignedList{
			Type: identity.Type(string(input.Type)),
		}

		if input.PrincipalID != nil {
			id := *input.PrincipalID
			config.PrincipalId = id.String()
		}
		if input.TenantID != nil {
			tenantID := *input.TenantID
			config.TenantId = tenantID.String()
		}

		userAssignedIdentityIds := make([]string, 0)
		for k := range input.UserAssignedIdentities {
			userAssignedIdentityIds = append(userAssignedIdentityIds, k)
		}
		config.IdentityIds = userAssignedIdentityIds
	}

	return identity.FlattenSystemAndUserAssignedList(config)
}
