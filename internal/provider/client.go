// Copyright (c) Optidata Cloud.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"

	"github.com/apache/cloudstack-go/v2/cloudstack"
)

type OpticloudClient struct {
	cs *cloudstack.CloudStackClient
}

func NewOpticloudClient(endpoint, apiKey, secretKey string, async bool) *OpticloudClient {
	cs := cloudstack.NewAsyncClient(endpoint, apiKey, secretKey, async)
	return &OpticloudClient{cs: cs}
}

func (c *OpticloudClient) CreateVM(name, serviceOfferingID, templateID, zoneID string) (*cloudstack.DeployVirtualMachineResponse, error) {
	p := c.cs.VirtualMachine.NewDeployVirtualMachineParams(serviceOfferingID, templateID, zoneID)
	p.SetName(name)
	return c.cs.VirtualMachine.DeployVirtualMachine(p)
}

func (c *OpticloudClient) GetVM(id string) (*cloudstack.VirtualMachine, error) {
	vm, _, err := c.cs.VirtualMachine.GetVirtualMachineByID(id)
	return vm, err
}

func (c *OpticloudClient) ListVMs(name string) ([]*cloudstack.VirtualMachine, error) {
	p := c.cs.VirtualMachine.NewListVirtualMachinesParams()
	if name != "" {
		p.SetName(name)
	}
	resp, err := c.cs.VirtualMachine.ListVirtualMachines(p)
	if err != nil {
		return nil, err
	}
	return resp.VirtualMachines, nil
}

func (c *OpticloudClient) GetZoneIDByName(name string) (string, error) {
	params := c.cs.Zone.NewListZonesParams()
	params.SetName(name)
	resp, err := c.cs.Zone.ListZones(params)
	if err != nil {
		return "", fmt.Errorf("error on getting zone '%s': %w", name, err)
	}
	if len(resp.Zones) == 0 {
		return "", fmt.Errorf("zone '%s' not found", name)
	}
	return resp.Zones[0].Id, nil
}

func (c *OpticloudClient) GetTemplateIDByName(name string) (string, error) {
	params := c.cs.Template.NewListTemplatesParams("all")
	params.SetName(name)
	resp, err := c.cs.Template.ListTemplates(params)
	if err != nil {
		return "", fmt.Errorf("error on getting template '%s': %w", name, err)
	}
	if len(resp.Templates) == 0 {
		return "", fmt.Errorf("template '%s' not found", name)
	}
	return resp.Templates[0].Id, nil
}

func (c *OpticloudClient) GetServiceOfferingIDByName(name string) (string, error) {
	params := c.cs.ServiceOffering.NewListServiceOfferingsParams()
	params.SetName(name)
	resp, err := c.cs.ServiceOffering.ListServiceOfferings(params)
	if err != nil {
		return "", fmt.Errorf("error on getting service offering '%s': %w", name, err)
	}
	if len(resp.ServiceOfferings) == 0 {
		return "", fmt.Errorf("service offering '%s' not found", name)
	}
	return resp.ServiceOfferings[0].Id, nil
}

func (c *OpticloudClient) UpdateVM(id, name, serviceOfferingID string) (*cloudstack.UpdateVirtualMachineResponse, error) {
	if id == "" {
		return nil, fmt.Errorf("instance ID required")
	}

	params := c.cs.VirtualMachine.NewUpdateVirtualMachineParams(id)

	if name != "" {
		params.SetName(name)
	}

	resp, err := c.cs.VirtualMachine.UpdateVirtualMachine(params)
	if err != nil {
		return nil, fmt.Errorf("update error: %w", err)
	}

	return resp, nil
}

func (c *OpticloudClient) DeleteVM(id string) error {
	params := c.cs.VirtualMachine.NewDestroyVirtualMachineParams(id)
	_, err := c.cs.VirtualMachine.DestroyVirtualMachine(params)
	if err != nil {
		return fmt.Errorf("delete error %s: %w", id, err)
	}
	return nil
}
