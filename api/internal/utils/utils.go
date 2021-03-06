/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package utils

import (
	"github.com/sony/sonyflake"
	"net"
	"os"
	"strconv"
)

var _sf *sonyflake.Sonyflake

func init() {
	saltStr, ok := os.LookupEnv("FLAKE_SALT")
	var salt uint16
	if ok {
		i, err := strconv.Atoi(saltStr)
		if err != nil {
			panic(err)
		}
		salt = uint16(i)
	}

	_sf = sonyflake.NewSonyflake(sonyflake.Settings{
		MachineID: func() (u uint16, e error) {
			return sumIP(GetOutboundIP()) + salt, nil
		},
	})
	if _sf == nil {
		panic("sonyflake init failed")
	}
}

func sumIP(ip net.IP) uint16 {
	total := 0
	for i := range ip {
		total += int(ip[i])
	}
	return uint16(total)
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

func GetFlakeUid() uint64 {
	uid, err := _sf.NextID()
	if err != nil {
		panic("get sony flake uid failed:" + err.Error())
	}
	return uid
}

func GetFlakeUidStr() string {
	return strconv.FormatUint(GetFlakeUid(), 10)
}
