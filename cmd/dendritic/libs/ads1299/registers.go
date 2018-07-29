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

package ads1299

import (
	"github.com/edwarnicke/dendritic/cmd/dendritic/libs/spicmds"
	"github.com/kubernetes/kubernetes/pkg/kubelet/kubeletconfig/util/log"
	"periph.io/x/periph/conn/spi"
)

type Register byte

const (
	ID         Register = 0x00
	CONFIG1    Register = 0x01
	CONFIG2    Register = 0x02
	CONFIG3    Register = 0x03
	LOFF       Register = 0x04
	CH1SET     Register = 0x05
	CH2SET     Register = 0x06
	CH3SET     Register = 0x07
	CH4SET     Register = 0x08
	CH5SET     Register = 0x09
	CH6SET     Register = 0x0A
	CH7SET     Register = 0x0B
	CH8SET     Register = 0x0C
	BIAS_SENSP Register = 0x0D
	BIAS_SENSN Register = 0x0E
	LOFF_SENSP Register = 0x0F
	LOFF_SENSN Register = 0x10
	LOFF_FLIP  Register = 0x11
	LOFF_STATP Register = 0x12
	LOFF_STATN Register = 0x13
	GPIO       Register = 0x14
	MISC1      Register = 0x15
	MISC2      Register = 0x16
	CONFIG4    Register = 0x17
)

func (r *Register) Read(c spi.Conn) (byte, error) {
	rreg := byte(spicmds.RREG) | byte(*r)
	write := []byte{rreg, 0x0}
	read := make([]byte, len(write))
	log.Infof("reading register %s (% x) on spi", r, rreg)
	if err := c.Tx(write, read); err != nil {
		log.Errorf("c.Tx(write, read) - returned err: %v", err)
		return 0x0, err
	}
	log.Infof("%s: % x", r, read)
	return read[0], nil
}

func (r *Register) Write(c spi.Conn, value byte) error {
	rreg := byte(spicmds.WREG) | byte(*r)
	write := []byte{rreg, 0x0, value}
	read := make([]byte, len(write))
	log.Infof("writing % x to register %s (% x) on spi", value, r, rreg)
	if err := c.Tx(write, read); err != nil {
		log.Errorf("c.Tx(write, read) - returned err: %v", err)
		return err
	}
	log.Infof("%s: % x", r, read)
	return nil
}
