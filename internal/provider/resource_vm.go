// Copyright (c) Optidata Cloud.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type OpticloudVMResource struct {
	client *OpticloudClient
}

func NewVMResource() resource.Resource {
	return &OpticloudVMResource{}
}

func (r *OpticloudVMResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "opticloud_instance"
}

func (r *OpticloudVMResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID da VM no CloudStack",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Nome da VM",
			},
			"service_offering": schema.StringAttribute{
				Required:    true,
				Description: "Nome do Service Offering",
			},
			"template": schema.StringAttribute{
				Required:    true,
				Description: "Nome do Template",
			},
			"zone": schema.StringAttribute{
				Required:    true,
				Description: "Nome da Zona",
			},
			"service_offering_id": schema.StringAttribute{
				Computed:    true,
				Description: "ID do Service Offering",
			},
			"template_id": schema.StringAttribute{
				Computed:    true,
				Description: "ID do Template",
			},
			"zone_id": schema.StringAttribute{
				Computed:    true,
				Description: "ID da Zona",
			},
		},
	}
}

func (r *OpticloudVMResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*OpticloudClient)
}

type VMModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	ServiceOffering   types.String `tfsdk:"service_offering"`
	Template          types.String `tfsdk:"template"`
	Zone              types.String `tfsdk:"zone"`
	ServiceOfferingID types.String `tfsdk:"service_offering_id"`
	ZoneID            types.String `tfsdk:"zone_id"`
	TemplateID        types.String `tfsdk:"template_id"`
}

func (r *OpticloudVMResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan VMModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneID, err := r.client.GetZoneIDByName(plan.Zone.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Zone não encontrada", err.Error())
		return
	}

	templateID, err := r.client.GetTemplateIDByName(plan.Template.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Template não encontrado", err.Error())
		return
	}

	serviceOfferingID, err := r.client.GetServiceOfferingIDByName(plan.ServiceOffering.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Service Offering não encontrado", err.Error())
		return
	}

	vmResp, err := r.client.CreateVM(plan.Name.ValueString(), serviceOfferingID, templateID, zoneID)
	if err != nil {
		resp.Diagnostics.AddError("Erro ao criar VM", err.Error())
		return
	}

	state := VMModel{
		ID:                types.StringValue(vmResp.Id),
		Name:              plan.Name,
		ServiceOffering:   plan.ServiceOffering,
		Template:          plan.Template,
		Zone:              plan.Zone,
		ServiceOfferingID: types.StringValue(vmResp.Serviceofferingid),
		TemplateID:        types.StringValue(vmResp.Templateid),
		ZoneID:            types.StringValue(vmResp.Zoneid),
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *OpticloudVMResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state VMModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	vm, err := r.client.GetVM(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Erro ao ler VM '"+state.ID.ValueString()+"'", err.Error())
		return
	}

	state.ID = types.StringValue(vm.Id)
	state.Name = types.StringValue(vm.Name)
	state.ServiceOffering = types.StringValue(vm.Serviceofferingname)
	state.Template = types.StringValue(vm.Templatename)
	state.Zone = types.StringValue(vm.Zonename)
	state.ServiceOfferingID = types.StringValue(vm.Serviceofferingid)
	state.TemplateID = types.StringValue(vm.Templateid)
	state.ZoneID = types.StringValue(vm.Zoneid)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *OpticloudVMResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan VMModel
	var originalstate VMModel

	diags := req.State.Get(ctx, &originalstate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateVM(originalstate.ID.ValueString(), plan.Name.ValueString(), "")
	if err != nil {
		resp.Diagnostics.AddError("Erro ao atualizar VM '"+originalstate.ID.ValueString()+"'", err.Error())
		return
	}

	vm, err := r.client.GetVM(originalstate.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Erro ao ler VM '"+originalstate.ID.ValueString()+"'", err.Error())
		return
	}

	state := VMModel{
		ID:                types.StringValue(originalstate.ID.ValueString()),
		Name:              types.StringValue(vm.Name),
		ServiceOffering:   types.StringValue(vm.Serviceofferingname),
		Template:          types.StringValue(vm.Templatename),
		Zone:              types.StringValue(vm.Zonename),
		ServiceOfferingID: types.StringValue(vm.Serviceofferingid),
		TemplateID:        types.StringValue(vm.Templateid),
		ZoneID:            types.StringValue(vm.Zoneid),
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *OpticloudVMResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddError("Não implementado", "Delete não é suportado neste momento.")
}
