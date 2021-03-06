---
layout: "guides"
page_title: "Nomad Deployment Guide"
sidebar_current: "guides-install-production-deployment-guide"
description: |-
  This deployment guide covers the steps required to install and
  configure a single HashiCorp Nomad cluster as defined in the
  Nomad Reference Architecture
ea_version: 0.9
---

# Nomad Reference Install Guide

This deployment guide covers the steps required to install and configure a single HashiCorp Nomad cluster as defined in the [Nomad Reference Architecture](/guides/install/production/reference-architecture.html).

These instructions are for installing and configuring Nomad on Linux hosts running the systemd system and service manager.

## Reference Material

This deployment guide is designed to work in combination with the [Nomad Reference Architecture](/guides/install/production/reference-architecture.html) and [Consul Deployment Guide](https://www.consul.io/docs/guides/deployment-guide.html). Although it is not a strict requirement to follow the Nomad Reference Architecture, please ensure you are familiar with the overall architecture design. For example, installing Nomad server agents on multiple physical or virtual (with correct anti-affinity) hosts for high-availability.

## Overview

To provide a highly-available single cluster architecture, we recommend Nomad server agents be deployed to more than one host, as shown in the [Nomad Reference Architecture](/guides/install/production/reference-architecture.html).

![Reference diagram](/assets/images/nomad_reference_diagram.png)

These setup steps should be completed on all Nomad hosts:

- [Download Nomad](#download-nomad)
- [Install Nomad](#install-nomad)
- [Configure systemd](#configure-systemd)
- [Configure Nomad](#configure-nomad)
- [Start Nomad](#start-nomad)

## Download Nomad

Precompiled Nomad binaries are available for download at [https://releases.hashicorp.com/nomad/](https://releases.hashicorp.com/nomad/) and Nomad Enterprise binaries are available for download by following the instructions made available to HashiCorp Enterprise customers.

```text
export NOMAD_VERSION="0.9.0"
curl --silent --remote-name https://releases.hashicorp.com/nomad/${NOMAD_VERSION}/nomad_${NOMAD_VERSION}_linux_amd64.zip
```

You may perform checksum verification of the zip packages using the SHA256SUMS and SHA256SUMS.sig files available for the specific release version. HashiCorp provides [a guide on checksum verification](https://www.hashicorp.com/security.html) for precompiled binaries.

## Install Nomad

Unzip the downloaded package and move the `nomad` binary to `/usr/local/bin/`. Check `nomad` is available on the system path.

```text
unzip nomad_${NOMAD_VERSION}_linux_amd64.zip
sudo chown root:root nomad
sudo mv nomad /usr/local/bin/
nomad version
```

The `nomad` command features opt-in autocompletion for flags, subcommands, and arguments (where supported). Enable autocompletion.

```text
nomad -autocomplete-install
complete -C /usr/local/bin/nomad nomad
```

Create a data directory for Nomad.

```text
sudo mkdir --parents /opt/nomad
```

## Configure systemd

Systemd uses [documented sane defaults](https://www.freedesktop.org/software/systemd/man/systemd.directives.html) so only non-default values must be set in the configuration file.

Create a Nomad service file at `/etc/systemd/system/nomad.service`.

```text
sudo touch /etc/systemd/system/nomad.service
```

Add this configuration to the Nomad service file:

```text
[Unit]
Description=Nomad
Documentation=https://nomadproject.io/docs/
Wants=network-online.target
After=network-online.target

[Service]
ExecReload=/bin/kill -HUP $MAINPID
ExecStart=/usr/local/bin/nomad agent -config /etc/nomad.d
KillMode=process
KillSignal=SIGINT
LimitNOFILE=infinity
LimitNPROC=infinity
Restart=on-failure
RestartSec=2
StartLimitBurst=3
StartLimitIntervalSec=10
TasksMax=infinity

[Install]
WantedBy=multi-user.target
```

The following parameters are set for the `[Unit]` stanza:

- [`Description`](https://www.freedesktop.org/software/systemd/man/systemd.unit.html#Description=) - Free-form string describing the nomad service
- [`Documentation`](https://www.freedesktop.org/software/systemd/man/systemd.unit.html#Documentation=) - Link to the nomad documentation
- [`Wants`](https://www.freedesktop.org/software/systemd/man/systemd.unit.html#Wants=) - Configure a dependency on the network service
- [`After`](https://www.freedesktop.org/software/systemd/man/systemd.unit.html#After=) - Configure an ordering dependency on the network service being started before the nomad service

The following parameters are set for the `[Service]` stanza:

- [`ExecReload`](https://www.freedesktop.org/software/systemd/man/systemd.service.html#ExecReload=) - Send Nomad a `SIGHUP` signal to trigger a configuration reload
- [`ExecStart`](https://www.freedesktop.org/software/systemd/man/systemd.service.html#ExecStart=) - Start Nomad with the `agent` argument and path to a directory of configuration files
- [`KillMode`](https://www.freedesktop.org/software/systemd/man/systemd.kill.html#KillMode=) - Treat nomad as a single process
- [`LimitNOFILE`, `LimitNPROC`](https://www.freedesktop.org/software/systemd/man/systemd.exec.html#Process%20Properties) - Disable limits for file descriptors and processes
- [`RestartSec`](https://www.freedesktop.org/software/systemd/man/systemd.service.html#RestartSec=) - Restart nomad after 2 seconds of it being considered 'failed'
- [`Restart`](https://www.freedesktop.org/software/systemd/man/systemd.service.html#Restart=) - Restart nomad unless it returned a clean exit code
- [`StartLimitBurst`, `StartLimitIntervalSec`](https://www.freedesktop.org/software/systemd/man/systemd.unit.html#StartLimitIntervalSec=interval) - Configure unit start rate limiting
- [`TasksMax`](https://www.freedesktop.org/software/systemd/man/systemd.resource-control.html#TasksMax=N) - Disable task limits (only available in systemd >= 226)

The following parameters are set for the `[Install]` stanza:

- [`WantedBy`](https://www.freedesktop.org/software/systemd/man/systemd.unit.html#WantedBy=) - Creates a weak dependency on nomad being started by the multi-user run level

## Configure Nomad

Nomad uses [documented sane defaults](/docs/configuration/index.html) so only non-default values must be set in the configuration file. Configuration can be read from multiple files and is loaded in lexical order. See the [full description](/docs/configuration/index.html) for more information about configuration loading and merge semantics.

Some configuration settings are common to both server and client Nomad agents, while some configuration settings must only exist on one or the other. Follow the [common configuration](#common-configuration) guidance on all hosts and then the specific guidance depending on whether you are configuring a Nomad [server](#server-configuration) or [client](#client-configuration).

- [Common Nomad configuration](#common-configuration)
- [Configure a Nomad server](#server-configuration)
- [Configure a Nomad client](#client-configuration)

### Common configuration

Create a configuration file at `/etc/nomad.d/nomad.hcl`:

```text
sudo mkdir --parents /etc/nomad.d
sudo chmod 700 /etc/nomad.d
sudo touch /etc/nomad.d/nomad.hcl
```

Add this configuration to the `nomad.hcl` configuration file:

~> **Note:** Replace the `datacenter` parameter value with the identifier you will use for the datacenter this Nomad cluster is deployed in.

```hcl
datacenter = "dc1"
data_dir = "/opt/nomad"
```

- [`datacenter`](/docs/configuration/index.html#datacenter) - The datacenter in which the agent is running.
- [`data_dir`](/docs/configuration/index.html#data_dir) - The data directory for the agent to store state.

### Server configuration

Create a configuration file at `/etc/nomad.d/server.hcl`:

```text
sudo touch /etc/nomad.d/server.hcl
```

Add this configuration to the `server.hcl` configuration file:

~> **NOTE** Replace the `bootstrap_expect` value with the number of Nomad servers you will use; three or five [is recommended](/docs/internals/consensus.html#deployment-table).

```hcl
server {
  enabled = true
  bootstrap_expect = 3
}
```

- [`server`](/docs/configuration/server.html#enabled) - Specifies if this agent should run in server mode. All other server options depend on this value being set.
- [`bootstrap_expect`](/docs/configuration/server.html#bootstrap_expect) - The number of expected servers in the cluster. Either this value should not be provided or the value must agree with other servers in the cluster.

### Client configuration

Create a configuration file at `/etc/nomad.d/client.hcl`:

```text
sudo touch /etc/nomad.d/client.hcl
```

Add this configuration to the `client.hcl` configuration file:

```hcl
client {
  enabled = true
}
```

- [`client`](/docs/configuration/client.html#enabled) - Specifies if this agent should run in client mode. All other client options depend on this value being set.

~> **NOTE** The [`options`](/docs/configuration/client.html#options-parameters) parameter can be used to enable or disable specific configurations on Nomad clients, unique to your use case requirements.

### ACL configuration

The [Access Control](/guides/security/acl.html) guide provides instructions on configuring and enabling ACLs.

### TLS configuration

Securing Nomad's cluster communication with mutual TLS (mTLS) is recommended for production deployments and can even ease operations by preventing mistakes and misconfigurations. Nomad clients and servers should not be publicly accessible without mTLS enabled.

The [Securing Nomad with TLS](/guides/security/securing-nomad.html) guide provides instructions on configuring and enabling TLS.

## Start Nomad

Enable and start Nomad using the systemctl command responsible for controlling systemd managed services. Check the status of the nomad service using systemctl.

```text
sudo systemctl enable nomad
sudo systemctl start nomad
sudo systemctl status nomad
```

## Next Steps

- Read [Outage Recovery](/guides/operations/outage.html) to learn
  the steps required to recover from a Nomad cluster outage.
- Read [Autopilot](/guides/operations/autopilot.html) to learn about
  features in Nomad 0.8 to allow for automatic operator-friendly
  management of Nomad servers.
