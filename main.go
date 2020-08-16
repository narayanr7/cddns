/*
ccdns
Author:  robert@x1sec.com
License: see LICENCE
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	cddns "github.com/x1sec/cddns/pkg"
)

func main() {

	var (
		setupMenuFlag  bool
		debugFlag      bool
		configFilePath string
	)

	flag.Usage = func() {
		h := "CDDNS - Cloudflare DDNS with proxy support\n"
		h += "v0.1 by @x1sec - https://github.com/x1sec\n"
		h += "\n"
		h += "Usage:\n"
		h += "   cddns [OPTIONS]\n\n"
		h += "      -s, --setup        Interactive setup menu to create configuration file\n"
		h += "      -c, --config-file  Custom location for configuration file\n"
		h += "      -d  --debug        Enable debug output"
		h += "\n\n"
		fmt.Fprintf(os.Stderr, h)
	}

	flag.BoolVar(&setupMenuFlag, "setup", false, "")
	flag.BoolVar(&setupMenuFlag, "s", false, "")
	flag.BoolVar(&debugFlag, "debug", false, "")
	flag.BoolVar(&debugFlag, "d", false, "")
	flag.StringVar(&configFilePath, "config-file", "", "")
	flag.StringVar(&configFilePath, "c", "", "")
	flag.Parse()

	var config cddns.Configuration
	cddns.DebugEnabled = debugFlag

	if setupMenuFlag {
		if cddns.CreateConfig(&config) {
			if cddns.SaveConfig(&config, configFilePath) {
				fmt.Printf("Successfully saved configuration to: %s\n", configFilePath)
			}
		} else {
			fmt.Fprintf(os.Stderr, "Not creating configuration.")
		}
		os.Exit(0)
	}

	if err := cddns.LoadConfig(&config, configFilePath); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load configuration. Try creating with --setup flag\n")
		os.Exit(1)
	}

	ipAddress := ""

	done := make(chan bool)
	publicIP := cddns.PublicIP{}

	cf := cddns.ZoneUpdater{Config: &config}
	log.Printf("Cloudflare DDNS running for zone: %s\n", config.ZoneName)
	go func() {
		for {
			newIpAddress, err := publicIP.GetIP()
			if err != nil {
				log.Fatal(err)
			}
			if ipAddress != newIpAddress {
				msg := fmt.Sprintf("Changing IP address from %s to %s", ipAddress, newIpAddress)
				cddns.DebugPrint(msg)

				cf.UpdateAddress(newIpAddress)
				ipAddress = newIpAddress
			}
			time.Sleep(time.Second * time.Duration(config.PollInterval))
		}
	}()

	<-done
}
