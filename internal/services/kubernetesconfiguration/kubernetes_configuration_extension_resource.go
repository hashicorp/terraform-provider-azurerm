package kubernetesconfiguration

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kubernetesconfiguration/sdk/2022-03-01/extensions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceKubernetesConfigurationExtension() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKubernetesConfigurationExtensionCreate,
		Read:   resourceKubernetesConfigurationExtensionRead,
		Update: resourceKubernetesConfigurationExtensionUpdate,
		Delete: resourceKubernetesConfigurationExtensionDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := extensions.ParseExtensionID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9-.]{0,252}$"),
					"name must be between 1 and 253 characters in length and may contain only letters, numbers, periods (.), hyphens (-), and must begin with a letter or number.",
				),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"cluster_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"cluster_resource_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"managedClusters",
				}, false),
			},

			"configuration_protected_settings": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"configuration_settings": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"extension_type": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"release_train": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"version"},
				ValidateFunc:  validation.StringIsNotEmpty,
			},

			"version": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"release_train"},
				ValidateFunc:  validation.StringIsNotEmpty,
			},

			"release_namespace": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"target_namespace"},
				ValidateFunc:  validation.StringIsNotEmpty,
			},

			"target_namespace": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"release_namespace"},
				ValidateFunc:  validation.StringIsNotEmpty,
			},

			"auto_upgrade_minor_version": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"aks_assigned_identity": commonschema.SystemAssignedIdentityComputed(),
		},
	}
}

func resourceKubernetesConfigurationExtensionCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KubernetesConfiguration.ExtensionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	clusterResourceName := d.Get("cluster_resource_name").(string)
	clusterRp := ""
	if clusterResourceName == "managedClusters" {
		clusterRp = "Microsoft.ContainerService"
	} else {
		return fmt.Errorf("failed to get cluster RP for `cluster_resource_name` %s", clusterResourceName)
	}

	id := extensions.NewExtensionID(subscriptionId, d.Get("resource_group_name").(string), clusterRp, clusterResourceName, d.Get("cluster_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_kubernetes_configuration_extension", id.ID())
		}
	}

	props := extensions.Extension{
		Properties: &extensions.ExtensionProperties{
			ConfigurationProtectedSettings: expandMapStringString(d.Get("configuration_protected_settings").(map[string]interface{})),
			ConfigurationSettings:          expandMapStringString(d.Get("configuration_settings").(map[string]interface{})),
			ExtensionType:                  utils.String(d.Get("extension_type").(string)),
		},
	}

	if releaseTrain, ok := d.GetOk("release_train"); ok {
		props.Properties.ReleaseTrain = utils.String(releaseTrain.(string))
	}

	if version, ok := d.GetOk("version"); ok {
		props.Properties.Version = utils.String(version.(string))
	}

	if releaseNamespace, ok := d.GetOk("release_namespace"); ok {
		props.Properties.Scope = &extensions.Scope{
			Cluster: &extensions.ScopeCluster{
				ReleaseNamespace: utils.String(releaseNamespace.(string)),
			},
		}
	}

	if targetNamespace, ok := d.GetOk("target_namespace"); ok {
		props.Properties.Scope = &extensions.Scope{
			Namespace: &extensions.ScopeNamespace{
				TargetNamespace: utils.String(targetNamespace.(string)),
			},
		}
	}

	if err := client.CreateThenPoll(ctx, id, props); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceKubernetesConfigurationExtensionRead(d, meta)
}

func resourceKubernetesConfigurationExtensionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KubernetesConfiguration.ExtensionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := extensions.ParseExtensionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ExtensionName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("cluster_name", id.ClusterName)
	d.Set("cluster_resource_name", id.ClusterResourceName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if err := d.Set("aks_assigned_identity", flattenAksAssignedIdentity(props.AksAssignedIdentity)); err != nil {
				return fmt.Errorf("setting `aks_assigned_identity`: %+v", err)
			}

			d.Set("auto_upgrade_minor_version", props.AutoUpgradeMinorVersion)
			d.Set("configuration_protected_settings", d.Get("configuration_protected_settings"))
			d.Set("configuration_settings", props.ConfigurationSettings)
			d.Set("extension_type", props.ExtensionType)
			d.Set("release_train", props.ReleaseTrain)

			if props.Scope != nil {
				if props.Scope.Namespace != nil {
					d.Set("target_namespace", props.Scope.Namespace.TargetNamespace)
				}

				if props.Scope.Cluster != nil {
					d.Set("release_namespace", props.Scope.Cluster.ReleaseNamespace)
				}
			}

			d.Set("version", props.Version)
		}
	}

	return nil
}

func resourceKubernetesConfigurationExtensionUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KubernetesConfiguration.ExtensionsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := extensions.ParseExtensionID(d.Id())
	if err != nil {
		return err
	}

	props := extensions.PatchExtension{
		Properties: &extensions.PatchExtensionProperties{},
	}

	if d.HasChange("configuration_protected_settings") {
		props.Properties.ConfigurationProtectedSettings = expandMapStringString(d.Get("configuration_protected_settings").(map[string]interface{}))
	}

	if d.HasChange("configuration_settings") {
		props.Properties.ConfigurationSettings = expandMapStringString(d.Get("configuration_settings").(map[string]interface{}))
	}

	if d.HasChange("release_train") {
		if releaseTrain, ok := d.GetOk("release_train"); ok {
			props.Properties.ReleaseTrain = utils.String(releaseTrain.(string))
		}
	}

	if d.HasChange("version") {
		if version, ok := d.GetOk("version"); ok {
			props.Properties.Version = utils.String(version.(string))
		}
	}

	if err := client.UpdateThenPoll(ctx, *id, props); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceKubernetesConfigurationExtensionRead(d, meta)
}

func resourceKubernetesConfigurationExtensionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KubernetesConfiguration.ExtensionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := extensions.ParseExtensionID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id, extensions.DeleteOperationOptions{ForceDelete: utils.Bool(false)}); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandMapStringString(input map[string]interface{}) *map[string]string {
	result := make(map[string]string)
	for k, v := range input {
		result[k] = v.(string)
	}
	return &result
}

func flattenAksAssignedIdentity(input *extensions.ExtensionPropertiesAksAssignedIdentity) []interface{} {
	var transform *identity.SystemAssigned

	if input != nil {
		transform = &identity.SystemAssigned{}
		if input.Type != nil {
			transform.Type = identity.Type(*input.Type)
		}

		if input.PrincipalId != nil {
			transform.PrincipalId = *input.PrincipalId
		}

		if input.TenantId != nil {
			transform.TenantId = *input.TenantId
		}
	}

	return identity.FlattenSystemAssigned(transform)
}
