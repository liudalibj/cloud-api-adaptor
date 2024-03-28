//go:build cgo

// (C) Copyright Confidential Containers Contributors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"

	providers "github.com/confidential-containers/cloud-api-adaptor/src/cloud-providers"
	"github.com/confidential-containers/cloud-api-adaptor/src/cloud-providers/libvirt"
	"github.com/confidential-containers/cloud-api-adaptor/src/cloud-providers/util/cloudinit"
)

var logger = log.New(log.Writer(), "[adaptor/cloud/libvirt-ext] ", log.LstdFlags|log.Lmsgprefix)

type libvirtext struct {
	libvirtProvider providers.Provider
	serviceConfig   *libvirt.Config
}

func NewProvider(config *libvirt.Config) (providers.Provider, error) {

	libvirtProvider, err := libvirt.NewProvider(config)

	if err != nil {
		return nil, err
	}

	provider := &libvirtext{
		libvirtProvider: libvirtProvider,
		serviceConfig:   config,
	}

	return provider, nil
}

func (p *libvirtext) CreateInstance(ctx context.Context, podName, sandboxID string, cloudConfig cloudinit.CloudConfigGenerator, spec providers.InstanceTypeSpec) (*providers.Instance, error) {
	cloudInitCloudConfigData, ok := cloudConfig.(*cloudinit.CloudConfig)
	// Debug print cloudInitCloudConfigData
	if !ok {
		return nil, fmt.Errorf("User Data generator did not use the cloud-init Cloud Config data format")
	}
	userData, err := cloudInitCloudConfigData.Generate()
	if err != nil {
		return nil, err
	}
	logger.Printf("===CreateInstance: userData from libvirt-ext: %s", userData)

	return p.libvirtProvider.CreateInstance(ctx, podName, sandboxID, cloudConfig, spec)
}

func (p *libvirtext) DeleteInstance(ctx context.Context, instanceID string) error {
	return p.libvirtProvider.DeleteInstance(ctx, instanceID)
}

func (p *libvirtext) Teardown() error {
	return nil
}

func (p *libvirtext) ConfigVerifier() error {
	// Debug print p.serviceConfig
	logger.Printf("===p.serviceConfig from libvirt-ext: %v", p.serviceConfig)
	VolName := p.serviceConfig.VolName
	if len(VolName) == 0 {
		return fmt.Errorf("VolName is empty")
	}
	return nil
}
