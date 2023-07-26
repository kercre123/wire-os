#!/bin/bash

ssh -i ~/ssh_root_key root@$1 "systemctl stop wired"

scp -i ~/ssh_root_key -O ./build/wired root@$1:/usr/bin/

scp -i ~/ssh_root_key -Or ./webroot/* root@$1:/etc/wired/webroot/

rsync -e "ssh -i ~/ssh_root_key" -avr ./modfiles/* root@$1:/etc/wired/mods/

ssh -i ~/ssh_root_key root@$1 "systemctl start wired"
