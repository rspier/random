#!/bin/bash

#    Copyright 2020 Google LLC
#
#    Licensed under the Apache License, Version 2.0 (the "License");
#    you may not use this file except in compliance with the License.
#    You may obtain a copy of the License at
#
#        http://www.apache.org/licenses/LICENSE-2.0
#
#    Unless required by applicable law or agreed to in writing, software
#    distributed under the License is distributed on an "AS IS" BASIS,
#    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#    See the License for the specific language governing permissions and
#    limitations under the License.

# example.sh demonstrates how you might use update-ddns.sh.

# Set appropriate values here:
SECRET="hmac-sha256:tsig-keyname:abcdefghijk....etc"
ZONE="example.com"
HOST="dynamic"
SERVER="ns.example.com"

# If you're behind NAT, you'll have to use something fancier to figure out what
# your IP address is.  Replace these...
IPV4="$( ifconfig enp3s0  | grep inet | head -1  | awk '{print $2}' )"
IPV6="$( ifconfig enp3s0  | grep inet6.*global | head -1  | awk '{print $2}' )"

function update() {
    RRTYPE=$1
    VALUE=$2

    ./update-ddns.sh  --debug \
        --server "${SERVER}" \
        --secret "${SECRET}"  \
        --zone "${ZONE}" \
        --host "${HOST}"\
        --ttl 30 \
        --rrtype "${RRTYPE}" \
        --value "${VALUE}"
}

NOW="$( date +%s )"

update A "${IPV4}"
update AAAA "${IPV6}"


