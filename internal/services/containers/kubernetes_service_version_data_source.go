package containers

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceKubernetesServiceVersions() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceKubernetesServiceVersionsRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"location": azure.SchemaLocation(),

			"version_prefix": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"versions": {
				Type:     pluginsdk.TypeList,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				Computed: true,
			},

			"latest_version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"include_preview": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func dataSourceKubernetesServiceVersionsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ServicesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewServiceVersionID(subscriptionId, azure.NormalizeLocation(d.Get("location").(string)))

	listResp, err := client.ListOrchestrators(ctx, id.LocationName, "managedClusters")
	if err != nil {
		if utils.ResponseWasNotFound(listResp.Response) {
			return fmt.Errorf("Error: No Kubernetes Service versions found for location %q", id.LocationName)
		}
		return fmt.Errorf("retrieving Kubernetes Versions in %q: %+v", id.LocationName, err)
	}

	lv, err := version.NewVersion("0.0.0")
	if err != nil {
		return fmt.Errorf("Cannot set version baseline (likely an issue in go-version): %+v", err)
	}

	var versions []string
	versionPrefix := d.Get("version_prefix").(string)
	includePreview := d.Get("include_preview").(bool)

	if props := listResp.OrchestratorVersionProfileProperties; props != nil {
		if orchestrators := props.Orchestrators; orchestrators != nil {
			for _, rawV := range *orchestrators {
				if rawV.OrchestratorType == nil || rawV.OrchestratorVersion == nil {
					continue
				}

				isPreview := false
				orchestratorType := *rawV.OrchestratorType
				kubeVersion := *rawV.OrchestratorVersion

				if rawV.IsPreview == nil {
					log.Printf("[DEBUG] IsPreview is nil")
				} else {
					isPreview = *rawV.IsPreview
					log.Printf("[DEBUG] IsPreview is %t", isPreview)
				}

				if !strings.EqualFold(orchestratorType, "Kubernetes") {
					log.Printf("[DEBUG] Orchestrator %q was not Kubernetes", orchestratorType)
					continue
				}

				if versionPrefix != "" && !strings.HasPrefix(kubeVersion, versionPrefix) {
					log.Printf("[DEBUG] Version %q doesn't match the prefix %q", kubeVersion, versionPrefix)
					continue
				}

				if isPreview && !includePreview {
					log.Printf("[DEBUG] Orchestrator %q is a preview release, ignoring", kubeVersion)
					continue
				}

				versions = append(versions, kubeVersion)
				v, err := version.NewVersion(kubeVersion)
				if err != nil {
					log.Printf("[WARN] Cannot parse orchestrator version %q - skipping: %s", kubeVersion, err)
					continue
				}

				if v.GreaterThan(lv) {
					lv = v
				}
			}
		}
	}

	d.SetId(id.ID())
	d.Set("versions", versions)
	d.Set("latest_version", lv.Original())

	return nil
}
