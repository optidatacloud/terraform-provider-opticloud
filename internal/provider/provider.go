// Copyright (c) Optidata Cloud.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &OpticloudProvider{}

type OpticloudProvider struct {
	version string
	client  *OpticloudClient
}

type OpticloudProviderModel struct {
	APIURL    types.String `tfsdk:"api_url"`
	APIKey    types.String `tfsdk:"api_key"`
	SecretKey types.String `tfsdk:"secret_key"`
}

func (p *OpticloudProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "opticloud"
	resp.Version = p.version
}

func (p *OpticloudProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `The Opticloud provider is used to configure your Opticloud infrastructure.
To learn the basics of Terraform using this provider, follow the hands-on get started tutorials.
The provider needs to be configured with the proper credentials before it can be used using the Opticloud Manager APIs.`,
		Attributes: map[string]schema.Attribute{
			"api_url": schema.StringAttribute{
				Required:    true,
				Description: "Opticloud API URL",
			},
			"api_key": schema.StringAttribute{
				Required:    true,
				Description: "Opticloud API Key",
			},
			"secret_key": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "Opticloud Secret Key",
			},
		},
	}
}

func (p *OpticloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config OpticloudProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	p.client = NewOpticloudClient(
		config.APIURL.ValueString(),
		config.APIKey.ValueString(),
		config.SecretKey.ValueString(),
		true,
	)

	resp.ResourceData = p.client
	resp.DataSourceData = p.client
}

func (p *OpticloudProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewVMResource,
	}
}

func (p *OpticloudProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &OpticloudProvider{
			version: version,
		}
	}
}
