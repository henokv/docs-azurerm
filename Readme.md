[![goreleaser](https://github.com/henokv/docs-azurerm/actions/workflows/release.yml/badge.svg)](https://github.com/henokv/docs-azurerm/actions/workflows/release.yml)

# docs-azurerm

This is a tool meant to document resources in the azure cloud

## Installation
To install download the latest version from the [releases](https://github.com/henokv/docs-azurerm/releases) page or if you have go installed run the command
```shell
go install github.com/henokv/docs-azurerm@latest
```


## Authentication
To authenticate check the azure-sdk-for-go [authentication](https://learn.microsoft.com/en-gb/azure/developer/go/azure-sdk-authentication) docs

## Supported resources:
- Virtual Networks

### Virtual Networks
The below command will generate the new docs in a directory called 'docs' relative to the current working directory
```shell
docs-azurerm vnet
```

Features:
- Subnets with ip space, route table & nsgs
- Peered VNETs & link to the vnet if the user account has access to that VNET
- Size of the IP space still free in the vnets