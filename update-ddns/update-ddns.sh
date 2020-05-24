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
#

# shflags is from https://github.com/kward/shflags
source shflags

DEFINE_string host "" "hostname to update (do not include zone)" "n"
DEFINE_string zone "" "zone to update" "z"
DEFINE_string server "" "server to send updates to" "s"
DEFINE_string rrtype "A" "dns record type" "t"
DEFINE_string value "" "new value" "v"
DEFINE_integer ttl "30" "ttl" "l"
DEFINE_string secret "" "TSIG secret" "e"
DEFINE_boolean debug false "debug output" "d" 
DEFINE_boolean remove false "remove record" "rm"

FLAGS "$@" || exit $?
eval set -- "${FLAGS_ARGV}"

function message() {
    echo "server ${FLAGS_server}"
    echo "zone ${FLAGS_zone}"
    echo "update delete ${FLAGS_host}.${FLAGS_zone} ${FLAGS_rrtype}"
    if [[ ${FLAGS_remove} != ${FLAGS_TRUE} ]]; then
      echo "update add ${FLAGS_host}.${FLAGS_zone} ${FLAGS_ttl} ${FLAGS_rrtype} ${FLAGS_value}"
    fi
    echo send
}

if [[ -z "$FLAGS_host" || -z "${FLAGS_zone}" || -z "${FLAGS_server}" ]]; then 
    echo "Required args --host or --zone or --server missing."
    exit 1
fi

if [[ -z "${FLAGS_value}" &&  ${FLAGS_remove} == ${FLAGS_false} ]]; then
    echo "No --value specified."
    exit 1
fi

oldval="$(dig +short "${FLAGS_rrtype}" "${FLAGS_host}.${FLAGS_zone}")"

if [[ "${oldval}" == "${FLAGS_value}" ]]; then 
    # echo "old value matches new value, skipping"
    exit 0
fi

nsupdate_args=""
if [[ ${FLAGS_debug} == ${FLAGS_TRUE} ]]; then
  nsupdate_args+="-v -d "
  message
fi 
if [[ -n "${FLAGS_secret}" ]]; then
    nsupdate_args="-y ${FLAGS_secret}"
fi
# TODO(handle keyfile variant)

message | nsupdate ${nsupdate_args}
