// Copyright 2018 Ed Warnicke

//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at

//      http://www.apache.org/licenses/LICENSE-2.0

//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package libs

// type SpiCmd struct {
// 	Name string
// 	Byte [2]byte
// }

// func (s *SpiCmd) String() string {
// 	return s.Name
// }

// func (s *SpiCmd) Byte() [2]byte {
// 	return s.Byte
// }

// type Cmd [2]byte

var spiCmds = map[string]([]byte){
	"WAKEUP":  []byte{0x02, 0x00},
	"STANDBY": []byte{0x04, 0x00},
	"RESET":   []byte{0x06, 0x00},
	"START":   []byte{0x08, 0x00},
	"STOP":    []byte{0x0A, 0x00},
	"RDATAC":  []byte{0x10, 0x00},
	"SDATAC":  []byte{0x11, 0x00},
	"RDATA":   []byte{0x12, 0x00},
	"RREG":    []byte{0x20, 0x00},
	"WREG":    []byte{0x40, 0x00},
}

type SpiCmd [2]byte
