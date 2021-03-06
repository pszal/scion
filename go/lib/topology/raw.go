// Copyright 2017 ETH Zurich
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package topology

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"strings"

	"github.com/netsec-ethz/scion/go/lib/common"
	"github.com/netsec-ethz/scion/go/lib/overlay"
)

const CfgName = "topology.json"

const (
	ErrorOpen    = "Unable to open topology"
	ErrorParse   = "Unable to parse topology from JSON"
	ErrorConvert = "Unable to convert RawTopo to Topo"
)

var RawCurr *RawTopo

// Structures directly filled from JSON

// RawTopo is used to un/marshal from/to JSON and should usually not be used by
// Go code directly. Use Topo (from lib/topology/topology.go) instead.
type RawTopo struct {
	Timestamp          int64
	TimestampHuman     string
	ISD_AS             string
	Overlay            string
	MTU                int
	BorderRouters      map[string]RawBRInfo   `json:",omitempty"`
	ZookeeperService   map[int]RawAddrPort    `json:",omitempty"`
	BeaconService      map[string]RawAddrInfo `json:",omitempty"`
	CertificateService map[string]RawAddrInfo `json:",omitempty"`
	PathService        map[string]RawAddrInfo `json:",omitempty"`
	SibraService       map[string]RawAddrInfo `json:",omitempty"`
	RainsService       map[string]RawAddrInfo `json:",omitempty"`
	DiscoveryService   map[string]RawAddrInfo `json:",omitempty"`
}

type RawBRInfo struct {
	InternalAddrs []RawAddrInfo
	Interfaces    map[common.IFIDType]RawBRIntf
}

func (b RawBRInfo) String() string {
	var s []string
	s = append(s, fmt.Sprintf("Loc addrs:\n  %s\nInterfaces:", b.InternalAddrs))
	for ifid, intf := range b.Interfaces {
		s = append(s, fmt.Sprintf("%d: %+v", ifid, intf))
	}
	return strings.Join(s, "\n")
}

type RawBRIntf struct {
	InternalAddrIdx int
	Overlay         string
	Bind            *RawAddrPort
	Public          RawAddrPort
	Remote          RawAddrPort
	Bandwidth       int
	ISD_AS          string
	LinkType        string
	MTU             int
}

// Convert a RawBRIntf struct (filled from JSON) to a TopoAddr (used by Go code)
func (b RawBRIntf) localTopoAddr(o overlay.Type) (*TopoAddr, *common.Error) {
	s := &RawAddrInfo{
		Public: []RawAddrPortOverlay{
			{RawAddrPort: RawAddrPort{Addr: b.Public.Addr, L4Port: b.Public.L4Port}},
		},
	}
	if o.IsUDP() {
		s.Public[0].OverlayPort = b.Public.L4Port
	}
	if b.Bind != nil {
		s.Bind = []RawAddrPort{{Addr: b.Bind.Addr, L4Port: b.Bind.L4Port}}
	}
	return s.ToTopoAddr(o)
}

// make an AddrInfo object from a BR interface Remote entry
func (b RawBRIntf) remoteAddrInfo(o overlay.Type) (*AddrInfo, *common.Error) {
	ip := net.ParseIP(b.Remote.Addr)
	if ip == nil {
		return nil, common.NewError("Could not parse remote IP from string", "ip", b.Remote.Addr)
	}
	ai := &AddrInfo{Overlay: o, IP: ip, L4Port: b.Remote.L4Port}
	if o.IsUDP() {
		ai.OverlayPort = b.Remote.L4Port
	}
	return ai, nil
}

type RawAddrInfo struct {
	Public []RawAddrPortOverlay
	Bind   []RawAddrPort `json:",omitempty"`
}

func (s *RawAddrInfo) ToTopoAddr(ot overlay.Type) (t *TopoAddr, err *common.Error) {
	return TopoAddrFromRAI(s, ot)
}

func (rai RawAddrInfo) String() string {
	var s []string
	s = append(s, fmt.Sprintf("Public: %s", rai.Public))
	if len(rai.Bind) > 0 {
		s = append(s, fmt.Sprintf("Bind: %s", rai.Bind))
	}
	return strings.Join(s, "\n")
}

type RawAddrPort struct {
	Addr   string
	L4Port int
}

func (a RawAddrPort) String() string {
	return fmt.Sprintf("%s:%d", a.Addr, a.L4Port)
}

// Since Public addresses may be associated with an Overlay port, extend the
// structure used for Bind addresses.
type RawAddrPortOverlay struct {
	RawAddrPort
	OverlayPort int `json:",omitempty"`
}

func (a RawAddrPortOverlay) String() string {
	return fmt.Sprintf("%s:%d/%d", a.Addr, a.L4Port, a.OverlayPort)
}

// Load returns a new TopoMeta object loaded from 'path'.
func Load(path string) (*Topo, *common.Error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, common.NewError(ErrorOpen, "err", err)
	}
	return parse(b, path)
}

func parse(data []byte, path string) (*Topo, *common.Error) {
	rt := &RawTopo{}
	if err := json.Unmarshal(data, rt); err != nil {
		return nil, common.NewError(ErrorParse, "err", err, "path", path)
	}
	RawCurr = rt
	ct, err := TopoFromRaw(rt)
	if err != nil {
		return nil, common.NewError(ErrorConvert, "err", err, "path", path)
	}
	return ct, nil
}
