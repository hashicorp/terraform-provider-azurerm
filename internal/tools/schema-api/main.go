// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/differ"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/providerjson"
)

func main() {
	f := flag.NewFlagSet("providerJson", flag.ExitOnError)

	apiPort := f.Int("api-port", 8080, "the port on which to run the Provider JSON server")
	dumpSchema := f.Bool("dump", false, "used to simply dump the entire provider schema")
	providerName := f.String("provider-name", "azurerm", "set the provider name, defaults to `azurerm`")
	exportSchema := f.String("export", "", "export the schema to the given path/filename. Intended for use in the release process")
	detectBreakingChanges := f.String("detect", "", "compare current schema to named dump.")
	errorOnBreakingChange := f.Bool("error-on-violation", false, "should the detect mode exit with a non-zero error code. Defaults to `false`")

	if err := f.Parse(os.Args[1:]); err != nil {
		fmt.Printf("error parsing args: %+v", err)
		os.Exit(1)
	}

	data := providerjson.LoadData()

	switch {
	case pointer.From(dumpSchema):
		{
			// Call the method to stdout
			log.Printf("dumping schema for '%s'", *providerName)
			wrappedProvider := &providerjson.ProviderWrapper{
				ProviderName:  *providerName,
				SchemaVersion: "1",
			}
			if err := providerjson.DumpWithWrapper(wrappedProvider, data); err != nil {
				log.Fatalf("error dumping provider: %+v", err)
			}

			os.Exit(0)
		}

	case pointer.From(detectBreakingChanges) != "":
		{
			d := differ.Differ{}
			if violations := d.Diff(*detectBreakingChanges, *providerName); violations != nil {
				for _, v := range violations {
					log.Println(v)
				}
				if pointer.From(errorOnBreakingChange) {
					os.Exit(1)
				}
			}

			os.Exit(0)
		}

	case pointer.From(exportSchema) != "":
		{
			log.Printf("dumping schema for '%s'", *providerName)
			wrappedProvider := &providerjson.ProviderWrapper{
				ProviderName:  *providerName,
				SchemaVersion: "1",
			}
			if err := providerjson.WriteWithWrapper(wrappedProvider, data, *exportSchema); err != nil {
				log.Fatalf("error writing provider schema for %q to %q: %+v", *providerName, *exportSchema, err)
			}

			os.Exit(0)
		}
	}

	if *apiPort < 1024 || *apiPort > 65534 {
		log.Fatal(fmt.Printf("invalid value for apiport, must be between 1024 and 65534, got %+v", apiPort))
	}

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sig
		log.Printf("%s signal received, closing provider API server on port %d", sig, apiPort)
		os.Exit(0)
	}()

	mux := http.NewServeMux()
	// paths
	mux.HandleFunc(providerjson.DataSourcesList, data.ListDataSources)
	mux.HandleFunc(providerjson.ResourcesList, data.ListResources)

	mux.HandleFunc(providerjson.DataSourcesPath, data.DataSourcesHandler)
	mux.HandleFunc(providerjson.ResourcesPath, data.ResourcesHandler)

	mux.HandleFunc(providerjson.DumpSchema, data.DumpAllSchema)

	log.Printf("starting api service on localhost:%d", *apiPort)
	log.Println(http.ListenAndServe(fmt.Sprintf(":%d", *apiPort), mux))
}
