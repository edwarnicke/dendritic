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

type SpiCmd byte

const (
	WAKEUP    SpiCmd = 0x02
	STANDBY   SpiCmd = 0x04
	SPI_RESET SpiCmd = 0x06
	START     SpiCmd = 0x08
	STOP      SpiCmd = 0x0A
	RDATAC    SpiCmd = 0x10
	SDATAC    SpiCmd = 0x11
	RREG      SpiCmd = 0x20
	WREG      SpiCmd = 0x40
)
