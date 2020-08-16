/*
ccdns
Author:  robert@x1sec.com
License: see LICENCE
*/

package cddns

import (
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"time"
)

type PublicIP struct{}

func (i PublicIP) Extract(ip string) (string, error) {
	re := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
	match := re.FindString(ip)
	if match == "" {
		return "", errors.New("IP address not found")
	} else {
		return match, nil
	}
}

func (i *PublicIP) Urls() []string {
	urls := []string{"https://ifconfig.me", "https://icanhazip.com/", "http://dns.loopia.se/checkip/checkip.php", "https://upbox.net"}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(urls), func(i, j int) { urls[i], urls[j] = urls[j], urls[i] })
	return urls
}

func (i PublicIP) GetIP() (string, error) {

	for _, url := range i.Urls() {
		ip, err := i.Try(url)
		if err == nil {
			DebugPrint("Public IP address is " + ip)
			return ip, nil
		}
	}
	return "", errors.New("Can't identify public IP address")
}

func (i PublicIP) Try(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	ip, err := i.Extract(string(body))
	return ip, err
}
