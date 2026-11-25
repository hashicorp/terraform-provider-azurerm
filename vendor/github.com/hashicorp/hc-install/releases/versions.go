// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package releases

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/httpclient"
	rjson "github.com/hashicorp/hc-install/internal/releasesjson"
	"github.com/hashicorp/hc-install/internal/validators"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/src"
)

// Versions allows listing all versions of a product
// which match Constraints
type Versions struct {
	Product     product.Product
	Constraints version.Constraints
	Enterprise  *EnterpriseOptions // require enterprise version if set (leave nil for OSS)

	ListTimeout time.Duration

	// Install represents configuration for installation of any listed version
	Install InstallationOptions

	// ApiBaseURL is an optional field that specifies a custom URL to list the product versions from.
	// If ApiBaseURL is set, versions will be listed from this base URL instead of the default site,
	// and each listed ExactVersion will also be installed from it.
	// Note: The directory structure of the custom URL must match the HashiCorp releases site (including the index.json files).
	ApiBaseURL string

	// HTTPClient represents the client to use for making
	// all round trips between the library (client) and the server.
	//
	// Defaults to [httpclient.New] if not set.
	//
	// If set, it is also passed through to each listed ExactVersion,
	// so that installation uses the same client.
	HTTPClient *http.Client
}

type InstallationOptions struct {
	Timeout    time.Duration
	Dir        string
	LicenseDir string

	SkipChecksumVerification bool

	// ArmoredPublicKey is a public PGP key in ASCII/armor format to use
	// instead of built-in pubkey to verify signature of downloaded checksums
	// during installation
	ArmoredPublicKey string
}

func (v *Versions) httpClient() *http.Client {
	if v.HTTPClient == nil {
		return httpclient.New()
	}
	return v.HTTPClient
}

func (v *Versions) List(ctx context.Context) ([]src.Source, error) {
	if !validators.IsProductNameValid(v.Product.Name) {
		return nil, fmt.Errorf("invalid product name: %q", v.Product.Name)
	}

	if err := validateEnterpriseOptions(v.Enterprise, v.Install.LicenseDir); err != nil {
		return nil, err
	}

	timeout := defaultListTimeout
	if v.ListTimeout > 0 {
		timeout = v.ListTimeout
	}
	ctx, cancelFunc := context.WithTimeout(ctx, timeout)
	defer cancelFunc()

	r := rjson.NewReleases()
	if v.ApiBaseURL != "" {
		r.BaseURL = v.ApiBaseURL
	}
	r.SetHTTPClient(v.httpClient())
	pvs, err := r.ListProductVersions(ctx, v.Product.Name)
	if err != nil {
		return nil, err
	}

	versions := pvs.AsSlice()
	sort.Stable(versions)

	expectedMetadata := enterpriseVersionMetadata(v.Enterprise)

	installables := make([]src.Source, 0)
	for _, pv := range versions {
		if !v.Constraints.Check(pv.Version) {
			// skip version which doesn't match constraint
			continue
		}

		if pv.Version.Metadata() != expectedMetadata {
			// skip version which doesn't match required metadata for enterprise or OSS versions
			continue
		}

		ev := &ExactVersion{
			Product:    v.Product,
			Version:    pv.Version,
			InstallDir: v.Install.Dir,
			Timeout:    v.Install.Timeout,
			LicenseDir: v.Install.LicenseDir,

			ApiBaseURL: v.ApiBaseURL,
			HTTPClient: v.HTTPClient,

			ArmoredPublicKey:         v.Install.ArmoredPublicKey,
			SkipChecksumVerification: v.Install.SkipChecksumVerification,
		}

		if v.Enterprise != nil {
			ev.Enterprise = &EnterpriseOptions{
				Meta: v.Enterprise.Meta,
			}
		}

		installables = append(installables, ev)
	}

	return installables, nil
}
