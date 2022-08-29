// Licensed to SkyAPM org under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. SkyAPM org licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package idgen

import (
	"bytes"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/767829413/normal-frame/fork/SkyAPM/go2sky/internal/tool"
)

var ipAddr string
var number int64
var traceType = "A"
var pid string

func init() {
	// C0A88B79 ipv4 hex 8位
	ipv4 := tool.IPV4()
	res := strings.Split(ipv4, ".")
	for _, v := range res {
		i, _ := strconv.Atoi(v)
		if i < 16 {
			ipAddr += "0"
		}
		ipAddr += strings.ToUpper(strconv.FormatInt(int64(i), 16))
	}
	ipAddr += "-"

	s := os.Getenv("TRACE_TYPE")
	if len(s) == 1 && InArray(s, []string{"A", "B", "C", "D"}) {
		traceType = s
	}
	pid = strings.ToUpper(strconv.FormatInt(int64(os.Getpid()), 16))
	if len(pid) < 4 {
		pid = strings.Repeat("0", 4-len(pid)) + pid
	} else if len(pid) > 4 {
		pid = pid[:4]
	}
}

func InArray(s string, arr []string) bool {
	for k := range arr {
		if s == arr[k] {
			return true
		}
	}
	return false
}

// UUID generate UUID
/*func UUID() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(id.String(), "-", ""), nil
}

// GenerateGlobalID generates global unique id
func GenerateGlobalID() (globalID string, err error) {
	return TraceId()
}*/

// GenerateGlobalID generates global unique id
func GenerateGlobalID() (string, error) {
	return UUID()
}

func UUID() (string, error) {
	s := bytes.NewBufferString(ipAddr)
	// 13位时间戳
	s.WriteString(strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	s.WriteString("-")
	n := strconv.FormatInt(atomic.AddInt64(&number, 1)%1000000, 10)
	if len(n) < 6 {
		n = strings.Repeat("0", 6-len(n)) + n
	}
	s.WriteString(n)
	s.WriteString("-")
	s.WriteString(traceType)
	s.WriteString("-")
	s.WriteString(pid)
	return s.String(), nil
}
