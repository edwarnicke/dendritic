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
	"io"
	"sync"
	"time"

	"github.com/go-errors/errors"
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
	CLKSEL   = "GPIO17"
	DRDY     = "GPIO24"
	TCLK     = 500 * time.Nanosecond
	TPOR     = 262144 * TCLK // 2^18 * t_clk
)

type ADS1299 interface {
	Init() error
	Close() error
	ReadReg(r Register) (value byte, err error)
	DumpRegs() ([]byte, error)
}

func New() ADS1299 {
	return &ads1299{}
}

type conn interface {
	spi.Conn
	io.ReadWriter
}

type ads1299 struct {
	PWDN     gpio.PinOut
	RESET    gpio.PinOut
	CLKSEL   gpio.PinOut
	SPISTART gpio.PinOut
	DRDY     gpio.PinIn
	Port     spi.PortCloser
	Conn     conn
	sync.Mutex
}

func (a *ads1299) Reset() error {
	a.Lock()
	defer a.Unlock()
	cmd := []byte{byte(SPI_RESET)}
	if _, err := a.Conn.Write(cmd); err != nil {
		return errors.Wrap(err, 0)
	}
	time.Sleep(18 * TCLK)
	return nil
}

func (a *ads1299) Sdatac() error {
	a.Lock()
	defer a.Unlock()
	cmd := []byte{byte(SDATAC)}
	if _, err := a.Conn.Write(cmd); err != nil {
		return errors.Wrap(err, 0)
	}
	time.Sleep(4 * TCLK)
	return nil
}

func (a *ads1299) Rdatac() error {
	a.Lock()
	defer a.Unlock()
	cmd := []byte{byte(RDATAC)}
	if _, err := a.Conn.Write(cmd); err != nil {
		return errors.Wrap(err, 0)
	}
	time.Sleep(4 * TCLK)
	return nil
}

func (a *ads1299) Standy() error {
	a.Lock()
	defer a.Unlock()
	cmd := []byte{byte(STANDBY)}
	if _, err := a.Conn.Write(cmd); err != nil {
		return errors.Wrap(err, 0)
	}
	time.Sleep(4 * TCLK)
	return nil
}

func (a *ads1299) Wakeup() error {
	a.Lock()
	defer a.Unlock()
	cmd := []byte{byte(WAKEUP)}
	if _, err := a.Conn.Write(cmd); err != nil {
		return errors.Wrap(err, 0)
	}
	time.Sleep(4 * TCLK)
	return nil
}

func (a *ads1299) Start() error {
	a.Lock()
	defer a.Unlock()
	cmd := []byte{byte(START)}
	if _, err := a.Conn.Write(cmd); err != nil {
		return errors.Wrap(err, 0)
	}
	time.Sleep(4 * TCLK)
	return nil
}

func (a *ads1299) Init() error {
	_, err := host.Init()
	if err != nil {
		return errors.Wrap(err, 0)
	}
	if err := a.setup(); err != nil {
		return err
	}
	if err := a.PowerUp(); err != nil {
		return err
	}

	if err := a.Reset(); err != nil {
		return err
	}
	if err := a.Sdatac(); err != nil {
		return err
	}
	if err := a.WriteReg(CONFIG3, 0xE0); err != nil {
		return err
	}
	if err := a.WriteReg(CONFIG1, 0x96); err != nil {
		return err
	}
	if err := a.WriteReg(CONFIG2, 0xC0); err != nil {
		return err
	}
	for chset := CH1SET; chset <= CH8SET; chset++ {
		if err := a.WriteReg(chset, 0x01); err != nil {
			return err
		}
	}
	if err := a.Start(); err != nil {
		return err
	}

	// if err := a.Rdatac(); err != nil {
	// 	return err
	// }

	return nil
}

func (a *ads1299) setup() error {
	if err := a.setupPins(); err != nil {
		return err
	}
	if err := a.setupSPI(); err != nil {
		return errors.Wrap(err, 0)
	}
	return nil
}

func (a *ads1299) setupSPI() error {
	p, err := spireg.Open("SPI0.0")
	a.Port = p
	if err != nil {
		return errors.Wrap(err, 0)
	}

	c, err := p.Connect(8*physic.MegaHertz, spi.Mode1, 8)
	if err != nil {
		return errors.Errorf("p.Connect(8*physic.MegaHertz, spi.Mode1, 8) - returned err: %s", err)
	}
	con, ok := c.(conn)
	if !ok {
		return errors.Errorf("error could not convert spi.Conn to io.ReadWriter")
	}
	a.Conn = con
	return nil
}

func (a *ads1299) setupPins() error {
	a.PWDN = gpioreg.ByName(PWDN)
	a.RESET = gpioreg.ByName(RESET)
	a.CLKSEL = gpioreg.ByName(CLKSEL)
	a.SPISTART = gpioreg.ByName(SPISTART)
	a.DRDY = gpioreg.ByName(DRDY)
	if err := a.PWDN.Out(gpio.Low); err != nil {
		return err
	}
	if err := a.RESET.Out(gpio.Low); err != nil {
		return err
	}
	if err := a.DRDY.In(gpio.PullDown, gpio.FallingEdge); err != nil {
		return err
	}
	if err := a.CLKSEL.Out(gpio.High); err != nil {
		return err
	}
	if err := a.SPISTART.Out(gpio.Low); err != nil {
		return err
	}
	// Wait for oscilator to wake up
	time.Sleep(4 * TPOR)
	return nil
}

func (a *ads1299) PowerUp() error {
	a.Lock()
	defer a.Unlock()
	if err := a.PWDN.Out(gpio.High); err != nil {
		return err
	}
	time.Sleep(2 * TPOR)
	return nil
}

func (a *ads1299) PowerDown() error {
	a.Lock()
	defer a.Unlock()
	if err := a.PWDN.Out(gpio.Low); err != nil {
		return err
	}
	time.Sleep(2 * TPOR)
	return nil
}

func (a *ads1299) ReadReg(r Register) (value byte, err error) {
	if err := a.Sdatac(); err != nil {
		return 0, err
	}
	rreg := byte(RREG) | byte(r)
	write := []byte{rreg, 0x0}
	read := []byte{0x0}
	if _, err := a.Conn.Write(write); err != nil {
		return 0, err
	}
	if _, err := a.Conn.Read(read); err != nil {
		return 0, err
	}
	return read[0], nil
}

func (a *ads1299) DumpRegs() ([]byte, error) {
	regcount := 17
	rv := make([]byte, regcount)
	for reg := 0; reg < regcount; reg++ {
		r, err := a.ReadReg(Register(reg))
		if err != nil {
			return rv, err
		}
		rv[reg] = r
	}
	return rv, nil
}

func (a *ads1299) WriteReg(r Register, value byte) error {
	for v, _ := a.ReadReg(r); v != value; v, _ = a.ReadReg(r) {
		wreg := byte(WREG) | byte(r)
		write := []byte{wreg, 0x0, value}
		if _, err := a.Conn.Write(write); err != nil {
			return err
		}
		time.Sleep(4 * TCLK)
	}
	return nil
}

func (a *ads1299) Close() error {
	if a.Conn != nil {
		if err := a.PowerDown(); err != nil {
			return err
		}
	}
	if a.Port != nil {
		if err := a.Port.Close(); err != nil {
			return err
		}
	}
	return nil
}
