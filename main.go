// Copyright (c) Optidata Cloud.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/optidatacloud/terraform-provider-opticloud/internal/provider"
)

var (
	version string = "dev"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/optidatacloud/opticloud",
		Debug:   debug,
	}

	log.Printf("Iniciando provider Opticloud, versão: %s", version)
	if err := providerserver.Serve(context.Background(), provider.New(version), opts); err != nil {
		log.Fatalf("Erro ao iniciar provider: %s", err)
	}
}
