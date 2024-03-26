package provider

import (
	"flag"
	"os"
	"path/filepath"
	"plugin"
	"strings"
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
	logger.Printf("Loading cloud providers from %s", path)
	clouds, err := os.ReadDir(path)
	if err != nil {
		logger.Printf("failed to ReadDir %s", err)
	} else {
		for _, entry := range clouds {
			// Only check the ".so" files under PROVIDERS_PLUGIN_PATH, skip files in all sub directory
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".so") {
				fullpath := filepath.Join(path, entry.Name())
				logger.Printf("Found cloud provider file %s", fullpath)
				cloud, err := plugin.Open(fullpath)
				if err != nil {
					logger.Printf("Failed to open the cloud provider file %s", err)
				} else {
					// Every provider plugin must have "InitCloud" function
					// The code in the function should looks like: "provider.AddCloudProvider("xxx", &Manager{})"
					ifunc, err := cloud.Lookup("InitCloud")
					if err != nil {
						logger.Printf("Failed to find the InitCloud function %s", err)
					} else {
						// Run the InitCloud from cloud provider to register it to providerTable
						initCloudFunc := ifunc.(func())
						initCloudFunc()
					}
				}
			}
		}
	}
}

func Get(name string) CloudProvider {
	if os.Getenv("CLOUD_PROVIDER_PLUGIN_PATH") != "" {
		LoadCloudProviders(os.Getenv("CLOUD_PROVIDER_PLUGIN_PATH"))
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
