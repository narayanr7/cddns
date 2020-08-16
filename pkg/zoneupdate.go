/*
ccdns
Author:  robert@x1sec.com
License: see LICENCE
*/

package cddns

import (
	"fmt"
	"log"
)

type ZoneUpdater struct {
	api      *CfApi
	zoneID   string
	recordID string
	Config   *Configuration
	debug    bool
}

func (cf *ZoneUpdater) Init() {

	cf.api = &CfApi{token: cf.Config.Token, zoneName: cf.Config.ZoneName, proxied: cf.Config.Proxied}
	if cf.api.VerifyToken() == false {
		log.Fatal("Unable to verify token.. Exiting")
	}
	cf.zoneID, _ = cf.api.GetZone(cf.Config.ZoneName)
}

func (cf ZoneUpdater) UpdateAddress(ip string) {
	if cf.api == nil {
		cf.Init()
	}

	records := cf.api.GetRecords(cf.zoneID)

	for _, record := range records {
		if record.RecordType == "A" {
			cf.recordID = record.ID
			break
		}
	}

	if cf.recordID == "" {
		DebugPrint(fmt.Sprintf("Adding record. ZoneID=%s, IP=ip Proxied=%t\n", cf.zoneID, ip, cf.Config.Proxied))
		cf.api.AddRecord(cf.zoneID, ip, cf.Config.Proxied)
	} else {
		DebugPrint(fmt.Sprintf("Updating record. ZoneID=%s, RecordID=%s, IP=%s Proxied=%t\n", cf.zoneID, cf.recordID, ip, cf.Config.Proxied))
		cf.api.SetRecord(cf.zoneID, cf.recordID, ip, cf.Config.Proxied)
	}
}
