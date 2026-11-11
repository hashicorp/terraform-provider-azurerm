// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package releases

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/httpclient"
	"github.com/hashicorp/hc-install/internal/pubkey"
	rjson "github.com/hashicorp/hc-install/internal/releasesjson"
	isrc "github.com/hashicorp/hc-install/internal/src"
	"github.com/hashicorp/hc-install/internal/validators"
	"github.com/hashicorp/hc-install/product"
)

type LatestVersion struct {
	Product            product.Product
	Constraints        version.Constraints
	InstallDir         string
	Timeout            time.Duration
	IncludePrereleases bool

	// LicenseDir represents directory path where to install license files
	// (required for enterprise versions, optional for Community editions).
	LicenseDir string

	// Enterprise indicates installation of enterprise version (leave nil for Community editions)
	Enterprise *EnterpriseOptions

	SkipChecksumVerification bool

	// ArmoredPublicKey is a public PGP key in ASCII/armor format to use
	// instead of built-in pubkey to verify signature of downloaded checksums
	ArmoredPublicKey string

	// ApiBaseURL is an optional field that specifies a custom URL to download the product from.
	// If ApiBaseURL is set, the product will be downloaded from this base URL instead of the default site.
	// Note: The directory structure of the custom URL must match the HashiCorp releases site (including the index.json files).
	ApiBaseURL string

	// HTTPClient represents the client to use for making
	// all round trips between the library (client) and the server.
	//
	// Defaults to [httpclient.New] with logger passed through if set via [SetLogger] earlier.
	//
	// Caller is responsible for passing logger to the client via [httpclient.WithLogger]
	// when overriding defaults.
	HTTPClient *http.Client

	logger        *log.Logger
	pathsToRemove []string
}

func (*LatestVersion) IsSourceImpl() isrc.InstallSrcSigil {
	return isrc.InstallSrcSigil{}
}

// SetLogger sets [log.Logger] to log internal debug messages.
//
// If you override HTTPClient you may also need to pass
// logger there via [httpclient.WithLogger].
func (lv *LatestVersion) SetLogger(logger *log.Logger) {
	lv.logger = logger
}

func (lv *LatestVersion) log() *log.Logger {
	if lv.logger == nil {
		return discardLogger
	}
	return lv.logger
}

func (lv *LatestVersion) httpClient() *http.Client {
	if lv.HTTPClient == nil {
		return httpclient.New(httpclient.WithLogger(lv.log()))
	}
	return lv.HTTPClient
}

func (lv *LatestVersion) Validate() error {
	if !validators.IsProductNameValid(lv.Product.Name) {
		return fmt.Errorf("invalid product name: %q", lv.Product.Name)
	}

	if !validators.IsBinaryNameValid(lv.Product.BinaryName()) {
		return fmt.Errorf("invalid binary name: %q", lv.Product.BinaryName())
	}

	if err := validateEnterpriseOptions(lv.Enterprise, lv.LicenseDir); err != nil {
		return err
	}

	return nil
}

func (lv *LatestVersion) Install(ctx context.Context) (string, error) {
	timeout := defaultInstallTimeout
	if lv.Timeout > 0 {
		timeout = lv.Timeout
	}
	ctx, cancelFunc := context.WithTimeout(ctx, timeout)
	defer cancelFunc()

	if lv.pathsToRemove == nil {
		lv.pathsToRemove = make([]string, 0)
	}

	dstDir := lv.InstallDir
	if dstDir == "" {
		var err error
		dirName := fmt.Sprintf("%s_*", lv.Product.Name)
		dstDir, err = os.MkdirTemp("", dirName)
		if err != nil {
			return "", err
		}
		lv.pathsToRemove = append(lv.pathsToRemove, dstDir)
		lv.log().Printf("created new temp dir at %s", dstDir)
	}
	lv.log().Printf("will install into dir at %s", dstDir)

	client := lv.httpClient()

	rels := rjson.NewReleases()
	if lv.ApiBaseURL != "" {
		rels.BaseURL = lv.ApiBaseURL
	}
	rels.SetLogger(lv.log())
	rels.SetHTTPClient(client)
	versions, err := rels.ListProductVersions(ctx, lv.Product.Name)
	if err != nil {
		return "", err
	}

	if len(versions) == 0 {
		return "", fmt.Errorf("no versions found for %q", lv.Product.Name)
	}

	versionToInstall, ok := lv.findLatestMatchingVersion(versions, lv.Constraints)
	if !ok {
		return "", fmt.Errorf("no matching version found for %q", lv.Constraints)
	}

	d := &rjson.Downloader{
		Logger:           lv.log(),
		VerifyChecksum:   !lv.SkipChecksumVerification,
		ArmoredPublicKey: pubkey.DefaultPublicKey,
		BaseURL:          rels.BaseURL,
		HTTPClient:       client,
	}
	if lv.ArmoredPublicKey != "" {
		d.ArmoredPublicKey = lv.ArmoredPublicKey
	}
	if lv.ApiBaseURL != "" {
		d.BaseURL = lv.ApiBaseURL
	}
	licenseDir := lv.LicenseDir
	up, err := d.DownloadAndUnpack(ctx, versionToInstall, dstDir, licenseDir)
	if up != nil {
		lv.pathsToRemove = append(lv.pathsToRemove, up.PathsToRemove...)
	}
	if err != nil {
		return "", err
	}

	execPath := filepath.Join(dstDir, lv.Product.BinaryName())

	lv.pathsToRemove = append(lv.pathsToRemove, execPath)

	lv.log().Printf("changing perms of %s", execPath)
	err = os.Chmod(execPath, 0o700)
	if err != nil {
		return "", err
	}

	return execPath, nil
}

func (lv *LatestVersion) Remove(ctx context.Context) error {
	if lv.pathsToRemove != nil {
		for _, path := range lv.pathsToRemove {
			err := os.RemoveAll(path)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (lv *LatestVersion) findLatestMatchingVersion(pvs rjson.ProductVersionsMap, vc version.Constraints) (*rjson.ProductVersion, bool) {
	expectedMetadata := enterpriseVersionMetadata(lv.Enterprise)
	versions := make(version.Collection, 0)
	for _, pv := range pvs.AsSlice() {
		if !lv.IncludePrereleases && pv.Version.Prerelease() != "" {
			// skip prereleases if desired
			continue
		}

		if pv.Version.Metadata() != expectedMetadata {
			continue
		}

		if vc.Check(pv.Version) {
			versions = append(versions, pv.Version)
		}
	}

	if len(versions) == 0 {
		return nil, false
	}

	sort.Stable(versions)
	latestVersion := versions[len(versions)-1]

	return pvs[latestVersion.Original()], true
}
