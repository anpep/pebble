// Copyright (c) 2022 Canonical Ltd
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3 as
// published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"

	"github.com/jessevdk/go-flags"

	"github.com/canonical/pebble/internal/boot"
)

type cmdRevertFirmware struct {
	clientMixin
}

var shortRevertFirmwareHelp = `Switch to non-running boot slot`

var longRevertFirmwareHelp = `
The revert-firmware command sets the default boot slot to the non-running
slot with previous firmware.
`

func (cmd *cmdRevertFirmware) Execute(args []string) error {
	if len(args) > 1 {
		return ErrExtraArgs
	}

	err, currentBootSlot := boot.GetBootSlot()
	if err != nil {
		return err
	}

	mapping := map[string]string{"a": "b", "b": "a"}
	if err := boot.SetBootSlot(mapping[currentBootSlot]); err != nil {
		return err
	}

	fmt.Printf("reverted to boot slot %q", mapping[currentBootSlot])
	return nil
}

func init() {
	addCommand("revert-firmware", shortRevertFirmwareHelp, longRevertFirmwareHelp, func() flags.Commander { return &cmdRevertFirmware{} }, nil, nil)
}
