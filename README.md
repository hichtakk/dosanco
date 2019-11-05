Dosanco
=======

[![Build Status](https://github.com/hichikaw/dosanco/workflows/Build/badge.svg)](https://github.com/hichikaw/dosanco/actions?workflow=Build)
[![container repository](https://img.shields.io/badge/docker-v0.0.3-blue)](https://hub.docker.com/r/hichtakk/dosanco/)

**Dosanco**: Simple DataCenter infrastructure database.

<img src="./docs/image/logo.png" width="200">

## Description
Dosanco provides API server and command-line client for managing DC and IT infrastructure information.
It helps facility operators to find out the facility failure impact to IT equipment and IT operators to design physical location of IT equipment for high availability.

Dosanco manages these entities as follows.

### Network
- Subnets
- IP address allocations
- VLANs

### Node
- Node
- Node Group

### DataCenter
- Data Center
  - Floor
  - Data Hall
- UPS
  - Row PDUs
- Rack PDUs

## Architecture
For understanding relation between each resouces, please refer the architecture diagram.

[architecture](docs/architecture.md)

## Getting Started
Most of dosanco user will work with dosanco command-line client.

### Install CLI client
Download latest `dosanco` binary from [release](https://github.com/actapio/dosanco/releases) or you can build by yourself.

`$ go get -u github.com/hichikaw/dosanco/cli/dosanco`

### Configure CLI client
Before start using dosanco, you need to configure API server endpoint to cli client.

`$ dosanco config set-endpoint ${DOSANCO_ENDPOINT}`

Now you are ready to use dosanco!

### CLI operation
You will operate DC/IT information with dosanco subcommands.
Each subcommands require resource name argument after that.

```
$ dosanco ${SUBCOMMAND} ${RESOURCE} [arguments/options]
```

* **show**:    display registered resource information
* **create**:  register new resource to dosanco
* **update**:  update registered resource information
* **delete**:  delete registered resource information


```
# example 1: display data center information
$ dosanco show datacenter
Name            Address
EMT             1640 Riverside Drive, Hill Valley, CA
AHQ             890 Fifth Avenue, Manhattan, NY
JKR             1165 Shakespeare Ave, The Bronx, NY
```

```
# example 2: create new subnetwork
# Dosanco recognize the parent-child relation between exist network and new network.
# It automatically register with proper relationship along the prefix length.

$ dosanco create network 192.168.1.0/24 --description "my home network"
network created. ID: 20,  CIDR: 192.168.1.0/24,  Description: my home network

# Dosanco is capable of showing subnets in relation tree style. try with `-t/--tree` option.

$ dosanco show network --tree
0.0.0.0/0
   192.168.0.0/16
      192.168.1.0/24
```


### More Dosanco
Dosanco CLI client provides well explained help messages.
You can see them with option `-h/--help` for subcommands and resources.

## License

## Contribution
Please feel free to send PR to dosanco.
