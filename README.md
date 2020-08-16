# CDDNS
A DDNS (Dynamic DNS) agent for Cloudflare

## Introduction
CDDNS is an agent that continuously monitors changes in the public Internet IP address of the connection for a host and updates Cloudflare's DNS when there is a change. 
Cloudflare's reverse proxy functionality is supported. A useful use case is when you would like to host a webserver on an Internet connection where the public IP address may change (e.g. home Internet connection) and you would like to hide/protect your home Internet IP address when resolving and accessing the website. 

An example setup illustrated below has a Raspberry Pi running a webserver with CDDNS on a home Internet connection. The website's domain name (let's use example.com) resolves to Cloudflare's DNS servers. When a web client access example.com, the HTTP/S connection terminates on Cloudflare's infrastructure. Cloudflare then proxies the connection to the home Internet connection where the home router port forwards to the Raspberry Pi. At any time if the home Internet connection IP address changes, CDDNS sends an update informing Cloudflare of the new home IP address for it to proxy to.

![](doc/example.png)

# Usage
## Installation

### Binaries
Compiled arm and X86_64 executable files for Linux are available [here](https://github.com/x1sec/cddns/releases/).

An example to download and install on a Raspberry Pi:
```
$ wget https://github.com/x1sec/commit-stream/releases/latest/download/cddns-linux-amd64.gz
$ tar xf cddns-linux-arm.tar.gz
$ cd cddns-linux-arm.tar.gz
$ sudo ./install.sh
```

### Go get
If you would prefer to build yourself (and Go is setup [correctly](https://golang.org/doc/install)):
```
$ go get -u github.com/x1sec/cddns
```
### Building from source
```
$ git clone https://github.com/x1sec/cddns
$ cd cddns
$ make
```
To compile for arm (e.g. to run on a Raspberry pi):
```
$ make build-arm
```
To install as a system service:
```
$ sudo make install
```

## Running
### Usage options
```
Usage:
   cddns [OPTIONS]

      -s, --setup        Interactive setup menu to create configuration file
      -c, --config-file  Custom location for configuration file
      -d  --debug        Enable debug output
```
`ccdns` requires a configuration file which includes the zone/domain name and a Cloudflare token. This configuration file is generated either by selecting `y` in the system service installation or by running with the `--setup` parameter.
```
$ cddns --setup

Creating configuration
Enter cloudflare token: xxxxxxxxxxxxxx-xxxxxx
Verifying token ...OK!
Enter domain name: example.com
Do you wish to use cloudflare as a proxy (y/n)? y
Poll interval in seconds (default: 120) ?
```

By default the configuration file is written to `$HOME/.config/cddns/config.json`. Specify a custom location with `--config-file`.
The configuration file format:
```json
{
        "ZoneName": "domain_name_here",
        "Token": "token_here",
        "Proxied": true,
        "PollInterval": 120
}
```

### Running as a system service / daemon
An installation script is provided which will install `cddns` as a system service and start on reboot. The installation script can optionally invoke the configuration setup,

```
make install
``` 
or alternatively 
```./scripts/install.sh```

The installation script creates a new user named `cddns` and copies both the configuration and binary files to the directory `/opt/cddns/`.

## Setting up Cloudflare
`cddns` requires a Cloudflare token. After [creating an account with cloudflare](https://support.cloudflare.com/hc/en-us/articles/201720164-Creating-a-Cloudflare-account-and-adding-a-website) and [changing the nameservers in your domain registrar to to Cloudflare](https://support.cloudflare.com/hc/en-us/articles/205195708), a token needs to generated for `cddns`. 

Tokens can be generated under `My Profile` / `API Tokens`. Select 'Edit Zone: Use Template`.

![](doc/create_token_1.png)
