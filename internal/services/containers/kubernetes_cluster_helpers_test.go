// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/managedclusters"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/testclient"
)

var aksMetadataCache sync.Map

// AKSTestMetadata stores supported Kubernetes versions and Service Mesh revisions
// retrieved dynamically from Azure for a specific location.
type AKSTestMetadata struct {
	OlderKubernetesVersion        string
	CurrentKubernetesVersion      string
	OlderKubernetesVersionAlias   string
	CurrentKubernetesVersionAlias string
	PreviousMeshRevision          string
	LatestMeshRevision            string
}

// getAKSTestMetadata returns Kubernetes versions and Service Mesh revisions for an Azure location.
// Use this in test functions (*testing.T). It fails the test immediately if an error happens.
func getAKSTestMetadata(t *testing.T, loc string) AKSTestMetadata {
	t.Helper()
	meta, err := obtainAKSMetadata(loc)
	if err != nil {
		t.Fatalf("obtaining AKS test metadata for location %q: %+v", loc, err)
	}
	return meta
}

// obtainAKSMetadata returns Kubernetes versions and Service Mesh revisions for an Azure location.
// Use this in test config helpers where *testing.T is not available. It returns an error if retrieval fails.
func obtainAKSMetadata(loc string) (AKSTestMetadata, error) {
	norm := location.Normalize(loc)
	if val, ok := aksMetadataCache.Load(norm); ok {
		return val.(AKSTestMetadata), nil
	}

	client, err := testclient.Build()
	if err != nil {
		return AKSTestMetadata{}, fmt.Errorf("building test client: %+v", err)
	}

	locId := managedclusters.NewLocationID(client.Account.SubscriptionId, norm)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// 1. Query Kubernetes versions for location
	k8sResp, err := client.Containers.KubernetesClustersClient.ListKubernetesVersions(ctx, locId)
	if err != nil {
		return AKSTestMetadata{}, fmt.Errorf("listing kubernetes versions for location %q: %+v", norm, err)
	}

	var parsedVersions []*version.Version
	if k8sResp.Model != nil && k8sResp.Model.Values != nil {
		for _, v := range *k8sResp.Model.Values {
			if v.IsPreview != nil && *v.IsPreview {
				continue
			}
			if v.PatchVersions != nil {
				for patchStr := range *v.PatchVersions {
					if pv, err := version.NewVersion(patchStr); err == nil {
						parsedVersions = append(parsedVersions, pv)
					}
				}
			}
		}
	}

	sort.Sort(version.Collection(parsedVersions))
	if len(parsedVersions) < 2 {
		return AKSTestMetadata{}, fmt.Errorf("expected at least 2 non-preview Kubernetes versions in %q, got %d", norm, len(parsedVersions))
	}

	olderK8s := parsedVersions[len(parsedVersions)-2]
	currK8s := parsedVersions[len(parsedVersions)-1]

	var olderAlias, currAlias string
	if len(olderK8s.Segments()) >= 2 {
		olderAlias = fmt.Sprintf("%d.%d", olderK8s.Segments()[0], olderK8s.Segments()[1])
	} else {
		olderAlias = olderK8s.String()
	}
	if len(currK8s.Segments()) >= 2 {
		currAlias = fmt.Sprintf("%d.%d", currK8s.Segments()[0], currK8s.Segments()[1])
	} else {
		currAlias = currK8s.String()
	}

	// 2. Query Service Mesh revisions for location
	meshResult, err := client.Containers.KubernetesClustersClient.ListMeshRevisionProfilesComplete(ctx, locId)
	if err != nil {
		return AKSTestMetadata{}, fmt.Errorf("listing mesh revision profiles for location %q: %+v", norm, err)
	}

	type parsedMeshRevision struct {
		raw string
		ver *version.Version
	}
	var parsedMeshRevisions []parsedMeshRevision
	seenMesh := make(map[string]struct{})
	for _, profile := range meshResult.Items {
		if profile.Properties != nil && profile.Properties.MeshRevisions != nil {
			for _, mr := range *profile.Properties.MeshRevisions {
				if mr.Revision != nil {
					rev := *mr.Revision
					if _, exists := seenMesh[rev]; !exists {
						seenMesh[rev] = struct{}{}
						clean := strings.TrimPrefix(rev, "asm-")
						clean = strings.ReplaceAll(clean, "-", ".")
						v, _ := version.NewVersion(clean)
						parsedMeshRevisions = append(parsedMeshRevisions, parsedMeshRevision{
							raw: rev,
							ver: v,
						})
					}
				}
			}
		}
	}

	sort.Slice(parsedMeshRevisions, func(i, j int) bool {
		vi := parsedMeshRevisions[i].ver
		vj := parsedMeshRevisions[j].ver
		if vi != nil && vj != nil {
			return vi.LessThan(vj)
		}
		if vi != nil {
			return false
		}
		if vj != nil {
			return true
		}
		return parsedMeshRevisions[i].raw < parsedMeshRevisions[j].raw
	})

	if len(parsedMeshRevisions) < 2 {
		return AKSTestMetadata{}, fmt.Errorf("expected at least 2 service mesh revisions in %q, got %d", norm, len(parsedMeshRevisions))
	}

	metadata := AKSTestMetadata{
		OlderKubernetesVersion:        olderK8s.String(),
		CurrentKubernetesVersion:      currK8s.String(),
		OlderKubernetesVersionAlias:   olderAlias,
		CurrentKubernetesVersionAlias: currAlias,
		PreviousMeshRevision:          parsedMeshRevisions[len(parsedMeshRevisions)-2].raw,
		LatestMeshRevision:            parsedMeshRevisions[len(parsedMeshRevisions)-1].raw,
	}

	aksMetadataCache.Store(norm, metadata)
	return metadata, nil
}
