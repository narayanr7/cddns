/*
ccdns
Author:  robert@x1sec.com
License: see LICENCE
*/

package cddns

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type CfApi struct {
	token    string
	zoneName string
	zoneId   string
	recordID string
	proxied  bool
}

type Response struct {
	Success bool `json:"success"`
	Errors  []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`

	Messages []string `json:"messages"`
}

type Zone struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Record struct {
	ID         string `json:"id"`
	RecordType string `json:"type"`
	Content    string `json:"content"`
	Proxiable  string `json:"proxiable"`
}

type ZoneResult struct {
	Response
	Result []Zone `json:"result"`
}

type RecordResult struct {
	Response
	Result []Record `json:"result"`
}

type UpdateResult struct {
	Response
}
type ZoneUpdate struct {
	ID         string `json:"id"`
	RecordType string `json:"type"`
	Name       string `json:"name"`
	Content    string `json:"content"`
	Proxied    bool   `json:"proxied"`
}

type CreateZoneRequest struct {
	Name string `json:"name"`
}

func (api CfApi) GetZone(zoneName string) (string, error) {
	body := api.HttpRequest(http.MethodGet, "zones?name="+zoneName, nil)
	var res ZoneResult
	json.Unmarshal(body, &res)
	var msg string

	if res.Success == false {
		for _, errMsg := range res.Errors {
			msg = fmt.Sprintf("Unable to fetch zone:\n", zoneName, errMsg.Code, errMsg.Message)
			log.Println(msg)
		}
		return "", errors.New(msg)
	}

	return res.Result[0].ID, nil
}

func (api CfApi) GetRecords(zoneId string) []Record {
	path := fmt.Sprintf("zones/%s/dns_records?name=%s", zoneId, api.zoneName)

	body := api.HttpRequest(http.MethodGet, path, nil)
	var res RecordResult
	json.Unmarshal(body, &res)
	//fmt.Println(string(body))

	return res.Result
	//return res.Result[0].ID
}

func (api CfApi) AddRecord(zoneID string, ip string, proxied bool) {
	fmt.Println("adding record")
	path := fmt.Sprintf("zones/%s/dns_records", zoneID)
	zoneUpdate := ZoneUpdate{ID: zoneID, RecordType: "A", Name: api.zoneName, Content: ip, Proxied: proxied}
	postData, err := json.Marshal(zoneUpdate)
	if err != nil {
		log.Fatal(err)
	}
	body := api.HttpRequest(http.MethodPost, path, bytes.NewBuffer(postData))
	var res UpdateResult
	var msg string
	json.Unmarshal(body, &res)
	if res.Success == false {
		for _, errMsg := range res.Errors {
			msg = fmt.Sprintf("Unable to add new record to zone %s. Error: [%d] %s\n", api.zoneName, errMsg.Code, errMsg.Message)
			log.Println(msg)
		}
	} else {
		msg = fmt.Sprintf("Updated zone %s with new IP address %s\n", api.zoneName, ip)
		DebugPrint(msg)
	}
}

func (api CfApi) SetRecord(zoneId string, recordId string, ip string, proxied bool) {

	path := fmt.Sprintf("zones/%s/dns_records/%s", zoneId, recordId)
	zoneUpdate := ZoneUpdate{ID: zoneId, RecordType: "A", Name: api.zoneName, Content: ip, Proxied: proxied}
	//fmt.Println(body)
	putData, err := json.Marshal(zoneUpdate)

	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(string(putData))

	body := api.HttpRequest(http.MethodPut, path, bytes.NewBuffer(putData))
	var res UpdateResult
	var msg string
	json.Unmarshal(body, &res)

	if res.Success == false {
		for _, errMsg := range res.Errors {
			msg = fmt.Sprintf("Unable to update zone %s. Error: [%d] %s\n", api.zoneName, errMsg.Code, errMsg.Message)
			log.Println(msg)
		}
	} else {
		msg = fmt.Sprintf("Updated zone %s with IP address %s\n", api.zoneName, ip)
		DebugPrint(msg)
	}
}

func (api CfApi) VerifyToken() bool {
	path := "user/tokens/verify"
	body := api.HttpRequest(http.MethodGet, path, nil)
	var res Response
	json.Unmarshal(body, &res)
	if res.Success == false {
		for _, errMsg := range res.Errors {
			msg := fmt.Sprintf("Token verification failed. Error: [%d] %s\n", errMsg.Code, errMsg.Message)
			log.Println(msg)
		}
		return false
	}
	return true
}
func (api CfApi) CreateZone(zoneName string) bool {
	path := "zones"
	createZone := CreateZoneRequest{Name: zoneName}
	postData, err := json.Marshal(createZone)

	if err != nil {
		log.Fatal(err)
	}

	body := api.HttpRequest(http.MethodPost, path, bytes.NewBuffer(postData))
	var res Response
	json.Unmarshal(body, &res)

	if res.Success == false {
		for _, errMsg := range res.Errors {
			msg := fmt.Sprintf("Zone creation failed. Error: [%d] %s\n", errMsg.Code, errMsg.Message)
			log.Println(msg)
		}
		return false
	}
	return true
}

func (api CfApi) HttpRequest(method string, endpoint string, body io.Reader) []byte {

	baseUrl := "https://api.cloudflare.com/client/v4"
	url := fmt.Sprintf("%s/%s", baseUrl, endpoint)

	client := &http.Client{}

	req, err := http.NewRequest(method, url, body)

	if err != nil {
		log.Fatal(err)
	}

	authHeader := fmt.Sprintf("Bearer %s", api.token)
	req.Header.Add("Authorization", authHeader)
	req.Header.Add("Content-Type", "application/json")

	//requestDump, err := httputil.DumpRequest(req, true)
	//fmt.Println(string(requestDump))

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(resBody))
	if err != nil {
		log.Fatal(err)
	}
	return resBody
}
