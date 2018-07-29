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
	"fmt"
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

const (
	RESET    = "6"
	PWDN     = "13"
	SPISTART = "26"
)

type ADS1299 interface {
	Init() error
	Close() error
}

func New() ADS1299 {
	return &ads1299{}
}

type ads1299 struct {
	PWDN     gpio.PinIO
	RESET    gpio.PinIO
	SPISTART gpio.PinIO
	Port     spi.PortCloser
	Conn     spi.Conn
}

func (a *ads1299) Init() error {
	state, err := host.Init()
	if err != nil {
		logrus.Error("host.Init() resulted in err: %v", err)
		return err
	}
	logrus.Info("host.Init() resulted in state: %v", state)

	a.PWDN = gpioreg.ByName(PWDN)
	err = a.PWDN.Out(gpio.Low)
	if err != nil {
		logrus.Errorf("ad1299.PWDN.Out(gpio.Low) - returned err: %v", err)
		return err
	}
	logrus.Info("ad1299.PWDN.Out(gpio.Low)")

	a.RESET = gpioreg.ByName(RESET)
	err = a.RESET.Out(gpio.Low)
	if err != nil {
		logrus.Errorf("ad1299.Reset.Out(gpio.Low) - returned err: %v", err)
		return err
	}
	logrus.Info("ad1299.RESET.Out(gpio.Low)")

	a.SPISTART = gpioreg.ByName(SPISTART)
	err = a.SPISTART.Out(gpio.Low)
	if err != nil {
		logrus.Errorf("ad1299.SPISTART.Out(gpio.Low) - returned err: %v", err)
		return err
	}
	logrus.Info("ad1299.SPISTART.Out(gpio.Low)")

	logrus.Infof("Sleeping 500 ms")
	time.Sleep(500 * time.Millisecond)
	logrus.Infof("setting ad1299.PWDN.Out(gpio.High)")
	err = a.PWDN.Out(gpio.High)
	if err != nil {
		logrus.Errorf("ad1299.PWDN.Out(gpio.High) - returned err: %v", err)
	}
	logrus.Infof("Sleeping 500 ms")
	time.Sleep(500 * time.Millisecond)
	logrus.Infof("setting ads1299.RESET.Out(gpio.High)")
	a.RESET.Out(gpio.High)
	if err != nil {
		logrus.Errorf("ads1299.RESET.Out(gpio.High) - returned err: %v", err)
	}

	logrus.Infof("spireg.Open(\"0\")")
	p, err := spireg.Open("0")
	a.Port = p
	if err != nil {
		logrus.Errorf("spireg.Open(\"\") - returned err: %v", err)
	}

	logrus.Infof("p.Connect(200*physic.KiloHertz, spi.Mode1, 8)")
	c, err := p.Connect(200*physic.KiloHertz, spi.Mode1, 8)
	a.Conn = c
	if err != nil {
		logrus.Errorf("p.Connect(200*physic.KiloHertz, spi.Mode1, 8) - returned err", err)
	}

	ticker := time.NewTicker(500 * time.Millisecond)
	var reg byte
	go func() {
		for _ = range ticker.C {
			reg, err = a.ReadReg(ID)
			if err != nil {
				logrus.Errorf("error reading register %s - err: %s", ID, err)
				break
			}
			logrus.Infof("register %s - % x", ID, reg)
		}
	}()
	time.Sleep(5 * time.Second)
	if reg == 0x0 {
		err = fmt.Errorf("failed to read non-zero ID register")
	}
	return err
}

func (a *ads1299) ReadReg(r Register) (value byte, err error) {
	rreg := byte(RREG) | byte(r)
	write := []byte{rreg, 0x0}
	read := make([]byte, len(write))
	logrus.Infof("reading register %s (% x) on spi", r, rreg)
	if err := a.Conn.Tx(write, read); err != nil {
		log.Errorf("c.Tx(write, read) - returned err: %v", err)
		return 0x0, err
	}
	log.Infof("%s: % x", r, read)
	return read[0], nil
}

func (a *ads1299) Close() error {
	a.Port.Close()
	return nil
}
