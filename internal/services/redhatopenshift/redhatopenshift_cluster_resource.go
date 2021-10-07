package redhatopenshift

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/parse"
	openShiftValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/redhatopenshift/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourceOpenShiftCluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		// Create: resourceOpenShiftClusterCreate,
		// Read:   resourceOpenShiftClusterRead,
		// Update: resourceOpenShiftClusterUpdate,
		// Delete: resourceOpenShiftClusterDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ClusterID(id)
			return err
		}),

		CustomizeDiff: pluginsdk.CustomDiffInSequence(
			pluginsdk.ForceNewIfChange("service_principal_profile.client_id", func(ctx context.Context, old, new, meta interface{}) bool {
				return old.(string) != new.(string)
			}),
		),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(90 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(90 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"cluster_profile": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"pull_secret": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"domain": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Computed: true,
							// TODO: Random ID as default
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"openshift_version": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"service_principal_profile": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"client_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: openShiftValidate.ClientID,
						},
						"client_secret": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"network_profile": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"pod_cidr": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							Default:      "10.128.0.0/14",
							ValidateFunc: validate.CIDR,
						},
						"service_cidr": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							Default:      "172.30.0.0/16",
							ValidateFunc: validate.CIDR,
						},
					},
				},
			},

			"master_profile": {
				Type:     pluginsdk.TypeList,
				Required: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"subnet_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"vm_size": {
							Type:         pluginsdk.TypeString,
							Required:     false,
							ForceNew:     true,
							Default:      "Standard_D8s_v3",
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"worker_profile": {
				Type:     pluginsdk.TypeList,
				Required: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     false,
							ForceNew:     true,
							Default:      "worker",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"vm_size": {
							Type:         pluginsdk.TypeString,
							Required:     false,
							ForceNew:     true,
							Default:      "Standard_D4s_v3",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"disk_size_gb": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ForceNew:     true,
							Computed:     true,
							Default:      "128",
							ValidateFunc: validation.IntAtLeast(1),
						},
						"node_count": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Computed:     true,
							Default:      "3",
							ValidateFunc: validation.IntBetween(0, 1000),
						},
					},
				},
			},

			"api_server_profile": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"visibility": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							Default:      Public,
							ValidateFunc: validate.CIDR,
						},
					},
				},
			},

			"ingress_profile": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"visibility": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							Default:      Public,
							ValidateFunc: validate.CIDR,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}
