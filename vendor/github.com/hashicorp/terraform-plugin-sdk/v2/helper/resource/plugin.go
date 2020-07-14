package resource

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	grpcplugin "github.com/hashicorp/terraform-plugin-sdk/v2/internal/helper/plugin"
	proto "github.com/hashicorp/terraform-plugin-sdk/v2/internal/tfplugin5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	tftest "github.com/hashicorp/terraform-plugin-test"
)

func runProviderCommand(f func() error, wd *tftest.WorkingDir, opts *plugin.ServeOpts) error {
	// Run the provider in the same process as the test runner using the
	// reattach behavior in Terraform. This ensures we get test coverage
	// and enables the use of delve as a debugger.

	// the provider name is technically supposed to be specified in the
	// format returned by addrs.Provider.GetDisplay(), but 1. I'm not
	// importing the entire addrs package for this and 2. we only get the
	// provider name here. Fortunately, when only a provider name is
	// specified in a provider block--which is how the config file we
	// generate does things--Terraform just automatically assumes it's in
	// the legacy namespace and the default registry.terraform.io host,
	// so we can just construct the output of GetDisplay() ourselves, based
	// on the provider name.
	providerName := wd.GetHelper().GetPluginName()

	// providerName gets returned as terraform-provider-foo, and we need
	// just foo. So let's fix that.
	providerName = strings.TrimPrefix(providerName, "terraform-provider-")

	// if we didn't override the logger, let's set a default one.
	if opts.Logger == nil {
		opts.Logger = hclog.New(&hclog.LoggerOptions{
			Name:   "plugintest",
			Level:  hclog.Trace,
			Output: ioutil.Discard,
		})
	}

	// this is needed so Terraform doesn't default to expecting protocol 4;
	// we're skipping the handshake because Terraform didn't launch the
	// plugin.
	os.Setenv("PLUGIN_PROTOCOL_VERSIONS", "5")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	config, closeCh, err := plugin.DebugServe(ctx, opts)
	if err != nil {
		return err
	}

	// plugin.DebugServe hijacks our log output location, so let's reset it
	logging.SetOutput()

	reattachInfo := map[string]plugin.ReattachConfig{}
	var namespaces []string
	host := "registry.terraform.io"
	if v := os.Getenv("TF_ACC_PROVIDER_NAMESPACE"); v != "" {
		namespaces = append(namespaces, v)
	} else {
		// unfortunately, we need to populate both of them
		// Terraform 0.12.26 and higher uses the legacy mode ("-")
		// Terraform 0.13.0 and higher uses the default mode ("hashicorp")
		// because of the change in how providers are addressed in 0.13
		namespaces = append(namespaces, "-", "hashicorp")
	}
	if v := os.Getenv("TF_ACC_PROVIDER_HOST"); v != "" {
		host = v
	}

	for _, ns := range namespaces {
		reattachInfo[strings.TrimSuffix(host, "/")+"/"+
			strings.TrimSuffix(ns, "/")+"/"+
			providerName] = config
	}

	reattachStr, err := json.Marshal(reattachInfo)
	if err != nil {
		return err
	}
	wd.Setenv("TF_REATTACH_PROVIDERS", string(reattachStr))

	// ok, let's call whatever Terraform command the test was trying to
	// call, now that we know it'll attach back to that server we just
	// started.
	err = f()
	if err != nil {
		log.Printf("[WARN] Got error running Terraform: %s", err)
	}

	// cancel the server so it'll return. Otherwise, this closeCh won't get
	// closed, and we'll hang here.
	cancel()

	// wait for the server to actually shut down; it may take a moment for
	// it to clean up, or whatever.
	<-closeCh

	// once we've run the Terraform command, let's remove the reattach
	// information from the WorkingDir's environment. The WorkingDir will
	// persist until the next call, but the server in the reattach info
	// doesn't exist anymore at this point, so the reattach info is no
	// longer valid. In theory it should be overwritten in the next call,
	// but just to avoid any confusing bug reports, let's just unset the
	// environment variable altogether.
	wd.Unsetenv("TF_REATTACH_PROVIDERS")

	// return any error returned from the orchestration code running
	// Terraform commands
	return err
}

// defaultPluginServeOpts builds ths *plugin.ServeOpts that you usually want to
// use when running runProviderCommand. It just sets the ProviderFunc to return
// the provider under test.
func defaultPluginServeOpts(wd *tftest.WorkingDir, providers map[string]*schema.Provider) *plugin.ServeOpts {
	var provider *schema.Provider
	for _, p := range providers {
		provider = p
	}
	return &plugin.ServeOpts{
		GRPCProviderFunc: func() proto.ProviderServer {
			return grpcplugin.NewGRPCProviderServer(provider)
		},
	}
}
