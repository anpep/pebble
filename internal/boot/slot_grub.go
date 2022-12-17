//go:build termus && grub
// +build termus,grub

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

package boot

import (
	"fmt"
	"path/filepath"

	"github.com/canonical/pebble/internal/boot/grubenv"
)

const grubenvPath string = "/EFI/termus/grubenv"

var currentBootSlot string = ""

func getESPMount() (error, *mount) {
	err, devnode := FindPartitionByLabel("esp")
	if err != nil {
		return err, nil
	}
	return nil, &mount{devnode, "/mnt/esp", "vfat", 0, ""}
}

func getBootSlotFromGrubenv(env *grubenv.Env) (error, string) {
	v := env.Get("default")
	switch v {
	case "0":
		return nil, "a"
	case "1":
		return nil, "b"
	default:
		return fmt.Errorf("no boot slot corresponds to GRUB entry %q", v), ""
	}
}

func setGrubenvBootSlot(env *grubenv.Env, slot string) error {
	var v string
	switch slot {
	case "a":
		v = "0"
		break
	case "b":
		v = "1"
		break
	default:
		return fmt.Errorf("invalid boot slot %q", slot)
	}
	env.Set("default", v)
	return nil
}

func GetBootSlot() (error, string) {
	if currentBootSlot == "" {
		err, m := getESPMount()
		if err != nil {
			return err, ""
		}
		if err := m.mount(); err != nil {
			return err, ""
		}
		defer m.unmount()

		// Get current boot slot from the GRUB env block
		env := grubenv.NewEnv(filepath.Join(m.target, grubenvPath))
		if err := env.Load(); err != nil {
			return err, ""
		}
		err, slot := getBootSlotFromGrubenv(env)
		if err != nil {
			return err, ""
		}
		currentBootSlot = slot
	}

	return nil, currentBootSlot
}

func SetBootSlot(slot string) error {
	err, currentBootSlot := GetBootSlot()
	if currentBootSlot == slot {
		return nil
	}

	err, m := getESPMount()
	if err != nil {
		return err
	}
	if err := m.mount(); err != nil {
		return err
	}
	defer m.unmount()

	// Set new boot slot
	env := grubenv.NewEnv(filepath.Join(m.target, grubenvPath))
	if err := env.Load(); err != nil {
		return err
	}
	if err := setGrubenvBootSlot(env, slot); err != nil {
		return err
	}
	if err := env.Save(); err != nil {
		return err
	}

	currentBootSlot = slot
	return nil
}
