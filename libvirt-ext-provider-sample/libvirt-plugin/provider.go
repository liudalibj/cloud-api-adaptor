//go:build cgo

// (C) Copyright Confidential Containers Contributors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"

	provider "github.com/confidential-containers/cloud-api-adaptor/cloud-providers"
	libvirt "github.com/confidential-containers/cloud-api-adaptor/cloud-providers/libvirt"
	"github.com/confidential-containers/cloud-api-adaptor/cloud-providers/util/cloudinit"
)

var logger = log.New(log.Writer(), "[adaptor/cloud/libvirt-ext] ", log.LstdFlags|log.Lmsgprefix)

type libvirtExtProvider struct {
	libvirtProvider provider.Provider
	serviceConfig   *libvirt.Config
}

func NewProvider(config *libvirt.Config) (provider.Provider, error) {

	libvirtProvider, err := libvirt.NewProvider(config)

	if err != nil {
		return nil, err
	}

	provider := &libvirtExtProvider{
		libvirtProvider: libvirtProvider,
		serviceConfig:   config,
	}

	return provider, nil
}

func (p *libvirtExtProvider) CreateInstance(ctx context.Context, podName, sandboxID string, cloudConfig cloudinit.CloudConfigGenerator, spec provider.InstanceTypeSpec) (*provider.Instance, error) {
	cloudInitCloudConfigData, _ := cloudConfig.(*cloudinit.CloudConfig)

	logger.Printf("===CreateInstance: cloudInitCloudConfigData: %s", cloudInitCloudConfigData)

	return p.libvirtProvider.CreateInstance(ctx, podName, sandboxID, cloudInitCloudConfigData, spec)
}

func (p *libvirtExtProvider) DeleteInstance(ctx context.Context, instanceID string) error {
	return p.libvirtProvider.DeleteInstance(ctx, instanceID)
}

func (p *libvirtExtProvider) Teardown() error {
	return nil
}

func (p *libvirtExtProvider) ConfigVerifier() error {
	VolName := p.serviceConfig.VolName
	if len(VolName) == 0 {
		return fmt.Errorf("VolName is empty")
	}

	return nil
}
