# :memo: Adding support for a new provider

### Step 1: Initialize and register the cloud provider manager

The provider-specific cloud manager should be placed under cloud providers repo: `https://github.com/confidential-containers/cloud-providers/blob/main/<provider>/`.

:information_source:[Example code](https://github.com/confidential-containers/cloud-providers/blob/main/aws)

### Step 2: Add provider specific code

Under `https://github.com/confidential-containers/cloud-providers/<provider>/`, start by adding a new file called `types.go`. This file defines a configuration struct that contains the required parameters for a cloud provider.

:information_source:[Example code](https://github.com/confidential-containers/cloud-providers/blob/main/aws/types.go)

#### Step 2.1: Implement the Cloud interface

Create a provider-specific manager file called `manager.go`, which implements the following methods for parsing command-line flags, loading environment variables, and creating a new provider.

- ParseCmd
- LoadEnv
- NewProvider

Create an `init` function to add your manager to the cloud provider table.

```go
func init() {
	provider.AddCloudProvider("aws", &Manager{})
}
```

:information_source:[Example code](https://github.com/confidential-containers/cloud-providers/blob/main/aws/manager.go)

#### Step 2.2: Implement the Provider interface

The Provider interface defines a set of methods that need to be implemented by the cloud provider for managing virtual instances. Add the required methods:

 - CreateInstance
 - DeleteInstance
 - Teardown

:information_source:[Example code](https://github.com/confidential-containers/cloud-providers/blob/main/aws/provider.go#L76-L175)

Also, consider adding additional files to modularize the code. You can refer to existing providers such as `aws`, `azure`, `ibmcloud`, and `libvirt` for guidance. Adding unit tests wherever necessary is good practice.

- Make a new tag for cloud-providers repo
```
git tag v0.8.x
git push origin v0.8.x
```

#### Step 2.3: Include Provider package for peerpod-ctrl manager
- Get the new tag verions of cloud provider for [peerpod-ctrl](https://github.com/confidential-containers/peerpod-ctrl)
```
go get github.com/confidential-containers/cloud-providers@v0.8.x
go mod tidy
```
- Use your provider in peerpod-ctrl operator
To include your provider you need reference it from the operator. Go build tags are used to selectively include different providers.

:information_source:[Example code](https://github.com/confidential-containers/peerpod-ctrl/blob/main/controllers/aws.go)

```go
//go:build aws
```
Note the comment at the top of the file, when building ensure `-tags=` is set to include your new provider. See the [Makefile](https://github.com/confidential-containers/peerpod-ctrl/blob/main/Makefile#L66) for more context and usage.
- Make a new tag for peerpod-ctrl repo
```
git tag v0.8.x
git push origin v0.8.x
```

#### Step 2.4: Include peerpod-ctrl and provider package from CAA main
- Get the new tag verions of peerpod-ctrl in CAA
```
go get github.com/confidential-containers/peerpod-ctrl@v0.8.x
go mod tidy
```
`go mod tidy` will update `cloud-providers` in CAA to `v0.8.x`, and we better keep use same verion for go mod `cloud-providers` and `peerpod-ctrl`.

- Use your provider in CAA main package
To include your provider you need reference it from the main package. Go build tags are used to selectively include different providers.

:information_source:[Example code](https://github.com/confidential-containers/cloud-api-adaptor/blob/main/cmd/cloud-api-adaptor/aws.go)

```go
//go:build aws
```
Note the comment at the top of the file, when building ensure `-tags=` is set to include your new provider. See the [Makefile](https://github.com/confidential-containers/cloud-api-adaptor/blob/main/Makefile#L26) for more context and usage.

#### Step 3: Add documentation on how to build a Pod VM image

For using the provider, a pod VM image needs to be created in order to create the peer pod instances. Add the instructions for building the peer pod VM image at the root directory similar to the other providers.

#### Step 4: Add E2E tests for the new provider

For more information, please refer to the section on [adding support for a new cloud provider](../test/e2e/README.md#adding-support-for-a-new-cloud-provider) in the E2E testing documentation.
