## Microsoft Azure Export for Terraform

A tool to bring your existing Azure resources under the management of Terraform.

### Video (v0.12.0)
[![Video](https://img.youtube.com/vi/LWk9SU7AmDA/hqdefault.jpg)](https://youtu.be/LWk9SU7AmDA)

## Goal

Azure Export for Terraform exports supported resources into Terraform state and generate the corresponding Terraform configuration. Both the Terraform state and configuration are expected to be consistent with the resources' remote state, i.e., `terraform plan` shows no diff. The user then is able to use Terraform to manage these resources.

It supports both the [Terraform AzureRM provider](https://github.com/hashicorp/terraform-provider-azurerm) and the [Terraform AzAPI provider](https://github.com/Azure/terraform-provider-azapi).

## Non Goal

The Terraform configurations generated by `aztfexport` are not meant to be comprehensive and do not ensure that the infrastructure can be fully reproduced from said generated configurations. For details, please see [limitations](#limitations)).

## Install

### From Release

Precompiled binaries and Window MSI are available at [Releases](https://github.com/Azure/aztfexport/releases).

### From Go toolchain

```bash
go install github.com/Azure/aztfexport@latest
```

### From Package Manager

#### Windows

```bash
winget install aztfexport
```

#### Homebrew (Linux/macOS)

```bash
brew install aztfexport
```

#### dnf (Linux)

Supported versions:

- RHEL 8 (amd64, arm64)
- RHEL 9 (amd64, arm64)

1. Import the Microsoft repository key:

    ```
    rpm --import https://packages.microsoft.com/keys/microsoft.asc
    ```

2. Add `packages-microsoft-com-prod` repository:

    ```
    ver=8 # or 9
    dnf install -y https://packages.microsoft.com/config/rhel/${ver}/packages-microsoft-prod.rpm
    ```

3. Install:

    ```
    dnf install aztfexport
    ```

#### apt (Linux)

Supported versions:

- Ubuntu 20.04 (amd64, arm64)
- Ubuntu 22.04 (amd64, arm64)

1. Import the Microsoft repository key:

    ```
    curl -sSL https://packages.microsoft.com/keys/microsoft.asc > /etc/apt/trusted.gpg.d/microsoft.asc
    ```

2. Add `packages-microsoft-com-prod` repository:

    ```
    ver=20.04 # or 22.04
    apt-add-repository https://packages.microsoft.com/ubuntu/${ver}/prod
    ```

3. Install:

    ```
    apt-get install aztfexport
    ```

#### AUR (Linux)

```bash
yay -S aztfexport
```

## Prerequisites

`aztfexport` requires a `terraform` executable installed in the `$PATH` with version `>= v0.12`.

## How it Works

`aztfexport` leverages [`aztft`](https://github.com/magodo/aztft) to identify the Terraform resource type corresponding to an Azure resource ID. Then it runs `terraform import` under the hood to import each resource. Afterwards, it runs [`tfadd`](https://github.com/magodo/tfadd) to generate the Terraform HCL code for each imported resource.

## Usage

Read the [Azure Export documentation](https://learn.microsoft.com/en-us/azure/developer/terraform/azure-export-for-terraform/export-terraform-overview) which covers scenarios and usage.

### Config

`aztfexport` will create a configuration file at `$HOME/.aztfexport/config.json`. This file is aim to be managed by command `aztfexport config [subcommand]`, which includes following subcommands:

- `get`: Get a config item
- `set`: Set a config item
- `show`: Show the full configuration

Currently, the following config items are supported:

- `installation_id`: A UUID created on first run. If there is Azure CLI or Azure Powershell installed on the current machine, the UUID will be the same value among these tools. Otherwise, a new one will be created. This is used as an identifier in the telemetry trace.
- `telemetry_enabled`: Enables telemetry. We use telemetry to identify issues and areas for improvement, in order to optimize this tool for better performance, reliability, and user experience. If you wish to disable our telemetry, set this to false.

## Limitations

Visit [this page](https://learn.microsoft.com/en-us/azure/developer/terraform/azure-export-for-terraform/export-terraform-concepts#limitations) on the Azure Export for Terraform documentation that discusses the currently known limitations of the tool.

## Additional Resources

- [The aztfexport Github Page](https://azure.github.io/aztfexport): Everything about aztfexport, including comparisons with other existing import solutions.
- [aztft](https://github.com/magodo/aztft): A Go program and library for identifying the correct Terraform AzureRM provider resource type on the Azure resource id.
- [tfadd](https://github.com/magodo/tfadd): A Go program and library for generating Terraform configuration from Terraform state.
