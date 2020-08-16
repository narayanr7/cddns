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
Compiled arm and X86_64 executable files for Linux are available [here](https://github.com/x1sec/cddns/releases/)

### Go get
If you would prefer to build yourself (and Go is setup [correctly](https://golang.org/doc/install)):
```
go get -u github.com/x1sec/cddns
```

### Building from source
```
make
```

### Installation
An installation script is provided which will install `cddns` as a system service and start on reboot.
```
make install
```
