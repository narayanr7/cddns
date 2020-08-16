/*
ccdns
Author:  robert@x1sec.com
License: see LICENCE
*/

package cddns

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const CONFIG_FILENAME = "config.json"

type Configuration struct {
	ZoneName     string `json:"ZoneName"`
	Token        string `json:"Token"`
	Proxied      bool   `json:"Proxied"`
	PollInterval int    `json:"PollInterval"`
}

func defaultConfigPath() string {
	userConfigDir, err := os.UserConfigDir()

	if err != nil {
		log.Fatal(err)
	}
	appConfigDir := filepath.Join(userConfigDir, "cddns")
	configFilePath := filepath.Join(appConfigDir, CONFIG_FILENAME)
	return configFilePath
}

func LoadConfig(configuration *Configuration, configFilePath string) error {

	// No configuration file path specified, try defaults
	if configFilePath == "" {

		// First try: current directory
		workingDir, _ := os.Getwd()
		configFilePath = filepath.Join(workingDir, CONFIG_FILENAME)

		// Second try: in $HOME/.config/
		if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
			configFilePath = defaultConfigPath()
		}

	}

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return err
	}
	DebugPrint("Using configuration file: " + configFilePath)
	configData, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	DebugPrint("Loaded config data:")
	DebugPrint(string(configData))

	if err := json.Unmarshal(configData, &configuration); err != nil {
		log.Panic(err)
	}
	return nil
}

func SaveConfig(configuration *Configuration, customFilePath string) bool {
	var configFilePath string
	if customFilePath == "" {
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			log.Fatal(err)
		}
		appConfigDir := filepath.Join(userConfigDir, "cddns")
		err = os.MkdirAll(appConfigDir, 0700)
		if err != nil {
			log.Fatal(err)
		}
		configFilePath = filepath.Join(appConfigDir, CONFIG_FILENAME)
	} else {
		configFilePath = customFilePath
	}

	jsonOut, err := json.MarshalIndent(configuration, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile(configFilePath, jsonOut, 0700); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Configuration saved to '%s'\n", configFilePath)
	return true
}

func CreateConfig(configuraiton *Configuration) bool {
	api := CfApi{}

	fmt.Println("Creating configuration")
	reader := bufio.NewReader(os.Stdin)

	/* Token input and validation */
	for {
		fmt.Print("Enter cloudflare token: ")
		configuraiton.Token, _ = reader.ReadString('\n')
		configuraiton.Token = strings.TrimSuffix(configuraiton.Token, "\n")
		api.token = configuraiton.Token
		fmt.Print("Verifying token ...")

		if ok := api.VerifyToken(); ok {
			fmt.Println("OK!")
			break
		} else {
			fmt.Println("Invalid. Try again or hit ctrl-c to exit.")
		}
	}

	/* Zone name input and optional creation */
	fmt.Print("Enter domain name: ")
	configuraiton.ZoneName, _ = reader.ReadString('\n')
	configuraiton.ZoneName = strings.TrimSuffix(configuraiton.ZoneName, "\n")
	/*_, err := api.GetZone(configuraiton.ZoneName)

	for {
		if err != nil {
			fmt.Printf("%s has not been configured in Cloudflare. Would you like me to create it (y/n)? ", configuraiton.ZoneName)
			resp, _ := reader.ReadString('\n')
			resp = strings.TrimSuffix(resp, "\n")
			if resp == "y" || resp == "Y" {
				if err := api.CreateZone(configuraiton.ZoneName); err == false {
					fmt.Println("Unable to create.")
				}
			} else if resp == "n" || resp == "N" {
				break
			}
		}
	}
	*/
	/* Proxied setting */
	for {
		fmt.Print("Do you wish to use cloudflare as a proxy (y/n)? ")
		useProxy, _ := reader.ReadString('\n')
		useProxy = strings.TrimSuffix(useProxy, "\n")

		if useProxy == "y" || useProxy == "Y" {
			configuraiton.Proxied = true
			break
		} else if useProxy == "n" || useProxy == "N" {
			configuraiton.Proxied = false
			break
		}
	}

	/* Poll interval */
	for {
		fmt.Print("Poll interval in seconds (default: 120) ?")
		interval, _ := reader.ReadString('\n')
		interval = strings.TrimSuffix(interval, "\n")
		if len(interval) == 0 {
			configuraiton.PollInterval = 120
			break
		}
		var err error
		configuraiton.PollInterval, err = strconv.Atoi(interval)
		if err != nil {
			fmt.Println("enter a valid integer")
		} else {
			break
		}
	}

	fmt.Println("\n=====================================")
	fmt.Println("Domain name: " + configuraiton.ZoneName)
	fmt.Printf("Proxy via Cloudlfare: %v\n", configuraiton.Proxied)
	fmt.Println("Cloudflare token: " + configuraiton.Token)
	fmt.Printf("Poll interval: %d\n", configuraiton.PollInterval)
	fmt.Println("=====================================\n")
	for {
		fmt.Print("Are these settings correct (y/n) ?")
		correct, _ := reader.ReadString('\n')
		correct = strings.TrimSuffix(correct, "\n")
		if correct == "y" || correct == "Y" {
			return true
		} else if correct == "n" || correct == "N" {
			return false
		}
	}

}
