//go:build cgo

// (C) Copyright Confidential Containers Contributors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"

	provider "github.com/confidential-containers/cloud-api-adaptor/cloud-providers"
	libvirt "github.com/confidential-containers/cloud-api-adaptor/cloud-providers/libvirt"
)

type Manager struct {
	libvirtManager *libvirt.Manager
}

func init() {
	libvirtManager := &libvirt.Manager{}
	manager := &Manager{
		libvirtManager: libvirtManager,
	}
	provider.AddCloudProvider("libvirt-ext", manager)
}

func (m *Manager) ParseCmd(flags *flag.FlagSet) {
	m.libvirtManager.ParseCmd(flags)
}

func (m *Manager) LoadEnv() {

	m.libvirtManager.LoadEnv()
}

func (m *Manager) NewProvider() (provider.Provider, error) {
	return NewProvider(m.libvirtManager.GetConfig())
}
