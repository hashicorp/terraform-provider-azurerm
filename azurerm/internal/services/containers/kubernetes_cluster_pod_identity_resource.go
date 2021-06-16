package containers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2021-03-01/containerservice"
	"github.com/Azure/azure-sdk-for-go/services/msi/mgmt/2018-11-30/msi"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/containers/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/containers/validate"
	msiParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/parse"
	msiValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceKubernetesClusterPodIdentity() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKubernetesClusterPodIdentityCreateUpdate,
		Read:   resourceKubernetesClusterPodIdentityRead,
		Update: resourceKubernetesClusterPodIdentityCreateUpdate,
		Delete: resourceKubernetesClusterPodIdentityDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ClusterID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"cluster_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ClusterID,
			},

			"pod_identity": {
				Type:         pluginsdk.TypeSet,
				Optional:     true,
				AtLeastOneOf: []string{"pod_identity", "exception"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"namespace": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"identity_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: msiValidate.UserAssignedIdentityID,
						},
					},
				},
			},

			"exception": {
				Type:         pluginsdk.TypeSet,
				Optional:     true,
				AtLeastOneOf: []string{"pod_identity", "exception"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"namespace": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"pod_labels": {
							Type:     pluginsdk.TypeMap,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func resourceKubernetesClusterPodIdentityCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.KubernetesClustersClient
	userAssignedIdentitiesClient := meta.(*clients.Client).MSI.UserAssignedIdentitiesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Managed Kubernetes Cluster Pod Identity.")

	clusterId, err := parse.ClusterID(d.Get("cluster_id").(string))
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, clusterId.ResourceGroup, clusterId.ManagedClusterName)
	if err != nil {
		return fmt.Errorf("checking for presence of existing %s: %s", clusterId, err)
	}

	if existing.ManagedClusterProperties == nil {
		return fmt.Errorf("`ManagedClusterProperties` is nil for %s: %s", clusterId, err)
	}

	if d.IsNewResource() {
		if existing.ManagedClusterProperties.PodIdentityProfile != nil &&
			existing.ManagedClusterProperties.PodIdentityProfile.Enabled != nil &&
			*existing.ManagedClusterProperties.PodIdentityProfile.Enabled {
			return tf.ImportAsExistsError("azurerm_kubernetes_cluster_pod_identity", clusterId.ID())
		}
	}

	podIdentities, err := expandKubernetesPodIdentities(ctx, userAssignedIdentitiesClient, d.Get("pod_identity").(*pluginsdk.Set).List())
	if err != nil {
		return err
	}
	existing.ManagedClusterProperties.PodIdentityProfile = &containerservice.ManagedClusterPodIdentityProfile{
		Enabled:                        utils.Bool(true),
		UserAssignedIdentities:         podIdentities,
		UserAssignedIdentityExceptions: expandKubernetesPodIdentityExceptions(d.Get("exception").(*pluginsdk.Set).List()),
	}

	if existing.NetworkProfile != nil && existing.NetworkProfile.NetworkPlugin == containerservice.NetworkPluginKubenet {
		existing.ManagedClusterProperties.PodIdentityProfile.AllowNetworkPluginKubenet = utils.Bool(true)
	}

	future, err := client.CreateOrUpdate(ctx, clusterId.ResourceGroup, clusterId.ManagedClusterName, existing)
	if err != nil {
		return fmt.Errorf("creating/updating Pod Identity within %s: %+v", clusterId, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/updation of Pod Identity within %s: %+v", clusterId, err)
	}

	d.SetId(clusterId.ID())

	return resourceKubernetesClusterPodIdentityRead(d, meta)
}

func resourceKubernetesClusterPodIdentityRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.KubernetesClustersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ManagedClusterName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Managed Kubernetes Cluster %q was not found in Resource Group %q - removing from state!", id.ManagedClusterName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if resp.ManagedClusterProperties == nil ||
		resp.ManagedClusterProperties.PodIdentityProfile == nil ||
		resp.ManagedClusterProperties.PodIdentityProfile.Enabled == nil ||
		!*resp.ManagedClusterProperties.PodIdentityProfile.Enabled {
		log.Printf("[DEBUG] Pod Identity Managed Kubernetes Cluster %q was not found in Resource Group %q - removing from state!", id.ManagedClusterName, id.ResourceGroup)
		d.SetId("")
		return nil
	}

	d.Set("cluster_id", resp.ID)
	if err := d.Set("pod_identity", flattenKubernetesPodIdentities(resp.ManagedClusterProperties.PodIdentityProfile.UserAssignedIdentities)); err != nil {
		return fmt.Errorf("setting `pod_identity`: %+v", err)
	}

	if err := d.Set("exception", flattenKubernetesPodIdentityException(resp.ManagedClusterProperties.PodIdentityProfile.UserAssignedIdentityExceptions)); err != nil {
		return fmt.Errorf("setting `exception`: %+v", err)
	}

	return nil
}

func resourceKubernetesClusterPodIdentityDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.KubernetesClustersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ClusterID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, id.ResourceGroup, id.ManagedClusterName)
	if err != nil {
		return fmt.Errorf("checking for presence of existing %s: %s", id, err)
	}

	if existing.ManagedClusterProperties == nil {
		return fmt.Errorf("`ManagedClusterProperties` is nil for %s: %s", id, err)
	}

	existing.ManagedClusterProperties.PodIdentityProfile = &containerservice.ManagedClusterPodIdentityProfile{
		Enabled: utils.Bool(false),
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ManagedClusterName, existing)
	if err != nil {
		return fmt.Errorf("deleting Pod Identity for %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Pod Identity for %s: %+v", id, err)
	}

	return nil
}

func expandKubernetesPodIdentities(ctx context.Context, userAssignedIdentitiesClient *msi.UserAssignedIdentitiesClient, input []interface{}) (*[]containerservice.ManagedClusterPodIdentity, error) {
	if len(input) == 0 {
		return nil, nil
	}
	result := make([]containerservice.ManagedClusterPodIdentity, 0)
	for _, raw := range input {
		v := raw.(map[string]interface{})

		identityId, err := msiParse.UserAssignedIdentityID(v["identity_id"].(string))
		if err != nil {
			return nil, err
		}
		resp, err := userAssignedIdentitiesClient.Get(ctx, identityId.ResourceGroup, identityId.Name)
		if err != nil {
			return nil, err
		}
		if resp.UserAssignedIdentityProperties == nil ||
			resp.UserAssignedIdentityProperties.ClientID == nil ||
			resp.UserAssignedIdentityProperties.PrincipalID == nil {
			return nil, fmt.Errorf("clientId or principalId is nil for %s", identityId)
		}
		result = append(result, containerservice.ManagedClusterPodIdentity{
			Name:      utils.String(v["name"].(string)),
			Namespace: utils.String(v["namespace"].(string)),
			Identity: &containerservice.UserAssignedIdentity{
				ResourceID: utils.String(identityId.ID()),
				ClientID:   utils.String(resp.UserAssignedIdentityProperties.ClientID.String()),
				ObjectID:   utils.String(resp.UserAssignedIdentityProperties.PrincipalID.String()),
			},
		})
	}
	return &result, nil
}

func expandKubernetesPodIdentityExceptions(input []interface{}) *[]containerservice.ManagedClusterPodIdentityException {
	if len(input) == 0 {
		return nil
	}
	result := make([]containerservice.ManagedClusterPodIdentityException, 0)
	for _, raw := range input {
		v := raw.(map[string]interface{})

		// issue https://github.com/hashicorp/terraform-plugin-sdk/issues/588
		// once it's resolved, we could remove the check empty logic
		name := v["name"].(string)
		namespace := v["namespace"].(string)
		if name == "" || namespace == "" {
			continue
		}

		result = append(result, containerservice.ManagedClusterPodIdentityException{
			Name:      utils.String(name),
			Namespace: utils.String(namespace),
			PodLabels: utils.ExpandMapStringPtrString(v["pod_labels"].(map[string]interface{})),
		})
	}
	return &result
}

func flattenKubernetesPodIdentities(inputs *[]containerservice.ManagedClusterPodIdentity) []interface{} {
	if inputs == nil {
		return []interface{}{}
	}

	result := make([]interface{}, 0)
	for _, v := range *inputs {
		var name, namespace, identityId string
		if v.Name != nil {
			name = *v.Name
		}
		if v.Namespace != nil {
			namespace = *v.Namespace
		}
		if v.Identity != nil && v.Identity.ResourceID != nil {
			identityId = *v.Identity.ResourceID
		}
		result = append(result, map[string]interface{}{
			"name":        name,
			"namespace":   namespace,
			"identity_id": identityId,
		})
	}
	return result
}

func flattenKubernetesPodIdentityException(inputs *[]containerservice.ManagedClusterPodIdentityException) []interface{} {
	if inputs == nil {
		return []interface{}{}
	}
	result := make([]interface{}, 0)
	for _, v := range *inputs {
		var name, namespace string
		if v.Name != nil {
			name = *v.Name
		}
		if v.Namespace != nil {
			namespace = *v.Namespace
		}

		result = append(result, map[string]interface{}{
			"name":       name,
			"namespace":  namespace,
			"pod_labels": utils.FlattenMapStringPtrString(v.PodLabels),
		})
	}
	return result
}
