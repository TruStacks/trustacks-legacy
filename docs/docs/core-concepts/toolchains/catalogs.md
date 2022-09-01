---
sidebar_position: 4
slug: /toolchains/catalogs
---

# Catalogs

Components are included in toolchains as dependencies through component catalogs.

The component is the single source of truth for component configuration, values template parameters, orchestration hooks, and all other component related functions and metadata, including the list of available components in the catalog.

## Manifest

The catalog manifest contains metadata about the catalog and the available components.

:::info

The manifest is exposed through a webserver at the `./well-known/catalog-manifest` path of the catalog host.

View the [public catalog](http://trustacks-catalog.eastus2.azurecontainer.io/.well-known/catalog-manifest)

:::

### Hook Source

The hook source contains the url and tag of the container image used to run component orchestration hooks.

#### Versioning

The hook source is the only versioned asset in the manifest. The manifest is intended to only be read once during the toolchain installation. Once the manifest is read, TruStacks generates helm assets and stores them as deployable helm charts.

All hooks will be pinned to the manifest version specified in the hook source during installation, while any operations on the toolchain after installation will be completed against the generated helm assets.
