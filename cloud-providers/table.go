package provider

import (
	"flag"
	"os"
	"path/filepath"
	"plugin"
	"strings"
)

const (
	Version      = "0.0.0"
	ProviderPath = "/providers/"
)

type CloudProvider interface {
	ParseCmd(flags *flag.FlagSet)
	LoadEnv()
	NewProvider() (Provider, error)
}

var providerTable map[string]CloudProvider = make(map[string]CloudProvider)

// LoadCloudProviders loads cloud providers from the directory with the given path, looking for
// all .so files in there and call the InitCloud function from the found .so files
func LoadCloudProviders(path string) {
	clouds, err := os.ReadDir(path)
	if err != nil {
		logger.Printf("failed to ReadDir %s", err)
	} else {
		for _, entry := range clouds {
			// Only check the ".so" files under path, skip files in the sub directory
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".so") {
				fullpath := filepath.Join(path, entry.Name())
				logger.Printf("Found cloud provider file %s", fullpath)
				cloud, err := plugin.Open(fullpath)
				if err != nil {
					logger.Printf("Failed to open the cloud provider file %s", err)
				} else {
					ifunc, err := cloud.Lookup("InitCloud")
					if err != nil {
						logger.Printf("Failed to find the InitCloud function %s", err)
					} else {
						// Run the InitCloud from cloud provider to register it to cloudTable
						initCloudFunc := ifunc.(func())
						initCloudFunc()
					}
				}
			}
		}
	}
}

func Get(name string) CloudProvider {
	// Get the length of the cloudTable
	length := len(providerTable)
	logger.Printf("Length of the providerTable: %d", length)
	if length == 0 {
		logger.Printf("Loading CloudProviders from %s", ProviderPath)
		LoadCloudProviders(ProviderPath)
	}
	return providerTable[name]
}

func AddCloudProvider(name string, cloud CloudProvider) {
	providerTable[name] = cloud
}

func List() []string {

	var list []string

	for name := range providerTable {
		list = append(list, name)
	}

	return list
}
