#!/bin/bash

# Copyright 2016 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

env GOOS=linux GOARCH=arm GOARM=7 go build -o sign main.go struct.go
sleep 5s

scp sign pi@192.168.1.2:/home/pi/smzdm/
scp config.yaml pi@192.168.1.2:/home/pi/smzdm/
scp mail.ghtml pi@192.168.1.2:/home/pi/smzdm/

