#!/bin/bash
INSTALL_PATH="/opt/cddns"

if [ "$EUID" -ne 0 ]; then
        echo "Please run as root user.. exiting"
        exit
fi

echo "Stopping service ..."
systemctl stop cddns.service

echo "Removing ${INSTALL_PATH} ..."
rm -rf $INSTALL_PATH 

echo "Deleting cddns user ..."
deluser cddns

echo "Removing systemd file ..."
rm /lib/systemd/system/cddns.service

echo "Reloading systemd ..."
systemctl daemon-reload

echo "Done!"
