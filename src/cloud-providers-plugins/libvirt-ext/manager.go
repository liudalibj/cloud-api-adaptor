//go:build cgo

// (C) Copyright Confidential Containers Contributors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"

	providers "github.com/confidential-containers/cloud-api-adaptor/src/cloud-providers"
	"github.com/confidential-containers/cloud-api-adaptor/src/cloud-providers/libvirt"
)

type Manager struct {
	libvirtManager *libvirt.Manager
}

func InitCloud() {
	libvirtManager := &libvirt.Manager{}
	manager := &Manager{
		libvirtManager: libvirtManager,
	}
	providers.AddCloudProvider("libvirt-ext", manager)
}

func (m *Manager) ParseCmd(flags *flag.FlagSet) {
	m.libvirtManager.ParseCmd(flags)
}

func (m *Manager) LoadEnv() {

	m.libvirtManager.LoadEnv()
}

func (m *Manager) NewProvider() (providers.Provider, error) {
	return NewProvider(m.libvirtManager.GetConfig())
}
