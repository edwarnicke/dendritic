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
	"strings"
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
		logrus.Errorf("host.Init() resulted in err: %v", err)
		return err
	}
	logrus.Infof("host.Init() resulted in Loaded drivers: %v", state.Loaded)

	listSPIPins()

	a.PWDN = gpioreg.ByName(PWDN)
	logrus.Infof("a.PWDN - %+v", a.PWDN)
	err = a.PWDN.Out(gpio.Low)
	if err != nil {
		logrus.Errorf("ad1299.PWDN.Out(gpio.Low) - returned err: %v", err)
		return err
	}
	logrus.Info("ad1299.PWDN.Out(gpio.Low)")

	a.RESET = gpioreg.ByName(RESET)
	logrus.Infof("a.RESET - %+v", a.RESET)
	err = a.RESET.Out(gpio.Low)
	if err != nil {
		logrus.Errorf("ad1299.Reset.Out(gpio.Low) - returned err: %v", err)
		return err
	}
	logrus.Info("ad1299.RESET.Out(gpio.Low)")

	a.SPISTART = gpioreg.ByName(SPISTART)
	logrus.Infof("a.SPISTART - %+v", a.SPISTART)
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

	logrus.Infof("spireg.Open(\"\")")
	p, err := spireg.Open("")
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

	if p, ok := p.(spi.Pins); ok {
		fmt.Printf("  Port: %+v\n", p)
		fmt.Printf("  CLK : %+v", p.CLK().Number())
		fmt.Printf("  MOSI: %+v", p.MOSI().Number())
		fmt.Printf("  MISO: %+v", p.MISO().Number())
		fmt.Printf("  CS  : %+v\n", p.CS().Number())
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
	regs, err := a.DumpRegs()
	if err != nil {
		logrus.Errorf("error reading dumping registers %s - err: %s", ID, err)
	}
	logrus.Infof("registers: %v", regs)
	return err
}

func (a *ads1299) ReadRegs(r Register, count byte) (value []byte, err error) {
	if count > (0x17 - byte(r)) {
		return nil, fmt.Errorf("count (%d) must be smaller than (23 (0x17) - register number 0x%x)", count, r)
	}
	rreg := byte(RREG) | byte(r)
	write := make([]byte, count+2)
	write[0] = rreg
	write[1] = count
	read := make([]byte, len(write))
	logrus.Infof("reading %d register %s (% x) on spi", count+1, r, rreg)
	if err := a.Conn.Tx(write, read); err != nil {
		log.Errorf("c.Tx(write, read) - returned err: %v", err)
		return nil, err
	}
	log.Infof("%s: %v", r, read)
	return read, nil
}

func (a *ads1299) ReadReg(r Register) (value byte, err error) {
	regs, err := a.ReadRegs(r, 0)
	return regs[0], err
}

func (a *ads1299) DumpRegs() ([]byte, error) {
	rv, err := a.ReadRegs(ID, 0x17)
	return rv, err
}

func (a *ads1299) Close() error {
	a.Port.Close()
	return nil
}

func listSPIPins() {
	// Enumerate all SPI ports available and the corresponding pins.
	fmt.Print("SPI ports available:\n")
	for _, ref := range spireg.All() {
		fmt.Printf("- %s\n", ref.Name)
		if ref.Number != -1 {
			fmt.Printf("  %d\n", ref.Number)
		}
		if len(ref.Aliases) != 0 {
			fmt.Printf("  %s\n", strings.Join(ref.Aliases, " "))
		}

		p, err := ref.Open()
		if err != nil {
			fmt.Printf("  Failed to open: %v", err)
		}
		if p, ok := p.(spi.Pins); ok {
			fmt.Printf("  CLK : %s", p.CLK())
			fmt.Printf("  MOSI: %s", p.MOSI())
			fmt.Printf("  MISO: %s", p.MISO())
			fmt.Printf("  CS  : %s\n", p.CS())
		}
		if err := p.Close(); err != nil {
			fmt.Printf("  Failed to close: %v", err)
		}
	}
}
