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

	"github.com/sirupsen/logrus"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/host"
)

const (
	RESET    = "GPIO6"
	PWDN     = "GPIO13"
	SPISTART = "GPIO26"
	CE0      = "GPIO8"
	CE1      = "GPIO7"
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
	CE0      gpio.PinIO
	CE1      gpio.PinIO
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

	a.CE0 = gpioreg.ByName(CE0)
	logrus.Infof("a.CE0 - %+v", a.CE0)
	err = a.CE0.Out(gpio.Low)
	if err != nil {
		logrus.Errorf("ad1299.CE0.Out(gpio.Low) - returned err: %v", err)
		return err
	}
	logrus.Info("ad1299.CE0.Out(gpio.Low)")

	a.CE1 = gpioreg.ByName(CE1)
	logrus.Infof("a.CE1 - %+v", a.CE1)
	err = a.CE1.Out(gpio.Low)
	if err != nil {
		logrus.Errorf("ad1299.CE1.Out(gpio.Low) - returned err: %v", err)
		return err
	}
	logrus.Info("ad1299.CE1.Out(gpio.Low)")

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
	p, err := spireg.Open("SPI0.0")
	a.Port = p
	if err != nil {
		logrus.Errorf("spireg.Open(\"\") - returned err: %v", err)
	}

	logrus.Infof("p.Connect(200*physic.KiloHertz, spi.Mode1, 8)")
	c, err := p.Connect(244*physic.KiloHertz, spi.Mode1|spi.NoCS, 8)
	a.Conn = c
	if err != nil {
		logrus.Errorf("p.Connect(244*physic.KiloHertz, spi.Mode1, 8) - returned err", err)
	}

	reg, err := a.ReadReg(ID)
	if err != nil {
		logrus.Errorf("error reading register %s - err: %s", ID, err)
	}
	logrus.Infof("register %s - % x", ID, reg)
	time.Sleep(500 * time.Millisecond)
	regs, err := a.DumpRegs()
	if err != nil {
		logrus.Errorf("error reading dumping registers %s - err: %s", ID, err)
	}

	logrus.Infof("registers: %v", regs)
	time.Sleep(500 * time.Millisecond)
	err = a.WriteReg(CH1SET, 0x60)
	if err != nil {
		logrus.Errorf("error writing registers %s - err: %s", CH1SET, err)
	}

	logrus.Infof("registers: %v", regs)
	time.Sleep(500 * time.Millisecond)
	err = a.WriteReg(CH1SET, 0x60)
	if err != nil {
		logrus.Errorf("error writing registers %s - err: %s", CH1SET, err)
	}

	logrus.Infof("registers: %v", regs)
	time.Sleep(500 * time.Millisecond)
	err = a.WriteReg(CH1SET, 0x60)
	if err != nil {
		logrus.Errorf("error writing registers %s - err: %s", CH1SET, err)
	}

	reg, err = a.ReadReg(ID)
	if err != nil {
		logrus.Errorf("error reading register %s - err: %s", ID, err)
	}
	logrus.Infof("register %s - % x", ID, reg)
	time.Sleep(500 * time.Millisecond)

	reg, err = a.ReadReg(CH1SET)
	if err != nil {
		logrus.Errorf("error reading register %s - err: %s", CH1SET, err)
	}
	logrus.Infof("register %s - % x", CH1SET, reg)
	time.Sleep(500 * time.Millisecond)

	return err
}

func (a *ads1299) ReadRegs(r Register, count byte) (value []byte, err error) {
	if count > (0x17 - byte(r)) {
		return nil, fmt.Errorf("count (%d) must be smaller than (23 (0x17) - register number 0x%x)", count, r)
	}
	rreg := byte(RREG) | byte(r)
	write := make([]byte, count+3)
	write[0] = rreg
	write[1] = count
	read := make([]byte, len(write))
	logrus.Infof("reading %d register %s (% x) on spi", count+1, r, rreg)
	if err := a.Conn.Tx(write, read); err != nil {
		logrus.Errorf("c.Tx(write, read) - returned err: %v", err)
		return nil, err
	}
	logrus.Infof("reading register %s: len(read): %d : %v", r, len(read), read)
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

func (a *ads1299) WriteReg(r Register, value byte) error {
	wreg := byte(WREG) | byte(r)
	write := []byte{wreg, 0x0, value, 0x0}
	read := make([]byte, len(write))
	logrus.Infof("writing value 0x%x to register %s (0x%x)", value, r, byte(r))
	if err := a.Conn.Tx(write, read); err != nil {
		logrus.Errorf("c.Tx(write, read) - returned err: %v", err)
		return err
	}
	return nil
}

func (a *ads1299) Close() error {
	if a.Port != nil {
		a.Port.Close()
	}
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
