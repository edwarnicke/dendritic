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

import (
	"time"

	"github.com/kubernetes/kubernetes/pkg/kubelet/kubeletconfig/util/log"
	"github.com/sirupsen/logrus"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/host"
)

type Dentrite struct {
	Conn spi.Conn
	Port spi.Port
}

func (d *Dentrite) Start() {
	state, err := host.Init()
	if err != nil {
		logrus.Error("host.Init() resulted in err: %v", err)
	}
	logrus.Info("host.Init() resulted in state: %v", state)
	reset := gpioreg.ByName(RESET)
	pwdn := gpioreg.ByName(PWDN)
	spiStart := gpioreg.ByName(SPI_START)
	logrus.Infof("setting reset.Out(gpio.Low)")
	err = reset.Out(gpio.Low)
	if err != nil {
		logrus.Errorf("reset.Out(gpio.Low) - returned err: %v", err)
	}
	logrus.Infof("setting pwdn.Out(gpio.Low)")
	err = pwdn.Out(gpio.Low)
	if err != nil {
		logrus.Errorf("pwdn.Out(gpio.Low) - returned err: %v", err)
	}
	logrus.Info("setting spiStart.Out(gpio.Low)")
	err = spiStart.Out(gpio.Low)
	if err != nil {
		logrus.Errorf("spiStart.Out(gpio.Low) - returned err: %v", err)
	}
	logrus.Infof("Sleeping 500 ms")
	time.Sleep(500 * time.Millisecond)
	logrus.Infof("setting pwdn.Out(gpio.High)")
	err = pwdn.Out(gpio.High)
	if err != nil {
		logrus.Errorf("pwdn.Out(gpio.High) - returned err: %v", err)
	}
	log.Infof("Sleeping 500 ms")

	log.Infof("spireg.Open(\"0\")")
	p, err := spireg.Open("0")
	d.Port = p
	if err != nil {
		log.Errorf("spireg.Open(\"\") - returned err: %v", err)
	}
	log.Infof("p.Connect(200*physic.KiloHertz, spi.Mode1, 8)")
	c, err := p.Connect(200*physic.KiloHertz, spi.Mode1, 8)
	if err != nil {
		log.Errorf("p.Connect(200*physic.KiloHertz, spi.Mode1, 8) - returned err", err)
	}
	d.Conn = c
}

func (d *Dentrite) ChipId() (byte, error) {
	write := []byte{0x32, 0x00}
	read := make([]byte, len(write))
	log.Infof("transmitting % x on spi", write)
	if err := d.Conn.Tx(write, read); err != nil {
		log.Errorf("c.Tx(write, read) - returned err: %v", err)
		return 0, nil
	}
	log.Infof("ChipID: % x", read)
	return read[0], nil
}

// func (d *Dentrite) SPICommand(cmd []byte) ([]byte, error) {
// 	read := make([]byte, len(cmd))
// 	return nil, nil
// }
