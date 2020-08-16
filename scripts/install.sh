#!/bin/bash

INSTALL_PATH=/opt/cddns/

if [ "$EUID" -ne 0 ]; then
	echo "Please run as root user.. exiting"
	exit
fi

read -p "Run cddns setup (y/n)? " choice
if [[ $choice =~ ^(y|yes|Y|YES) ]]; then
	./bin/cddns -s -c ./config.json
fi

if [ ! -d $INSTALL_PATH ]; then
	echo "Creating directory ${INSTALL_PATH}"
	mkdir $INSTALL_PATH
fi

echo "Creating cddns user ..."
useradd -r cddns

echo "Copying files to ${INSTALL_PATH} ..."
cp ./bin/cddns $INSTALL_PATH
cp ./config.json $INSTALL_PATH

chown -R cddns:cddns /opt/cddns/

echo "Installing service script ..."
cp ./scripts/cddns.service /lib/systemd/system/

echo "Reloading systemd daemon ..."
systemctl daemon-reload

echo "Enabling on startup ..."
systemctl enable cddns

echo "Starting service ..."
systemctl start cddns

echo "Cleaning up ..."
rm ./config.json

