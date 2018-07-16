#!/bin/bash

ServiceName=ninjaServer.service

# This script should be run while in the same directory as "ninjaServer.service"
# if you are intending to update it.
sudo cp "$ServiceName" /etc/systemd/system/"$ServiceName"

# Stop and disable the service (if it's there)
sudo systemctl stop "$ServiceName"
sudo systemctl disable "$ServiceName"

# Reloads systemd manager configuration and reruns the generators.
sudo systemctl daemon-reload

# Enables and Starts the service again.
sudo systemctl enable "$ServiceName"
sudo systemctl start "$ServiceName"

# Prints the service's status so you can see if everything worked correctly.
systemctl status "$ServiceName"
