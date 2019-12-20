package azurerm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmKubernetesServiceVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmKubernetesServiceVersionsRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"location": azure.SchemaLocation(),

			"version_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"versions": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},

			"latest_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmKubernetesServiceVersionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	location := azure.NormalizeLocation(d.Get("location").(string))

	listResp, err := client.ListOrchestrators(ctx, location, "managedClusters")
	if err != nil {
		if utils.ResponseWasNotFound(listResp.Response) {
			return fmt.Errorf("Error: No Kubernetes Service versions found for location %q", location)
		}
		return fmt.Errorf("Error retrieving Kubernetes Versions in %q: %+v", location, err)
	}

	lv, err := version.NewVersion("0.0.0")
	if err != nil {
		return fmt.Errorf("Cannot set version baseline (likely an issue in go-version): %+v", err)
	}

	var versions []string
	versionPrefix := d.Get("version_prefix").(string)

	if props := listResp.OrchestratorVersionProfileProperties; props != nil {
		if orchestrators := props.Orchestrators; orchestrators != nil {
			for _, rawV := range *orchestrators {
				if rawV.OrchestratorType == nil || rawV.OrchestratorVersion == nil {
					continue
				}

				orchestratorType := *rawV.OrchestratorType
				kubeVersion := *rawV.OrchestratorVersion

				if !strings.EqualFold(orchestratorType, "Kubernetes") {
					log.Printf("[DEBUG] Orchestrator %q was not Kubernetes", orchestratorType)
					continue
				}

				if versionPrefix != "" && !strings.HasPrefix(kubeVersion, versionPrefix) {
					log.Printf("[DEBUG] Version %q doesn't match the prefix %q", kubeVersion, versionPrefix)
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

	d.SetId(*listResp.ID)
	d.Set("versions", versions)
	d.Set("latest_version", lv.Original())

	return nil
}
