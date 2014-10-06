// +build windows

/*
   Copyright 2014 Franc[e]sco (lolisamurai@tfwno.gf)
   This file is part of kagami.
   kagami is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   kagami is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
   GNU General Public License for more details.
   You should have received a copy of the GNU General Public License
   along with kagami. If not, see <http://www.gnu.org/licenses/>.
*/

package utils

import (
	"syscall"
	"unsafe"
)

var (
	modkernel32          = syscall.NewLazyDLL("kernel32.dll")
	procGetConScrBufInfo = modkernel32.NewProc("GetConsoleScreenBufferInfo")
)

type coord struct {
	x int16
	y int16
}

type smallRect struct {
	left   int16
	top    int16
	right  int16
	bottom int16
}

type consoleScreenBuffer struct {
	size       coord
	cursorPos  coord
	attrs      int32
	window     smallRect
	maxWinSize coord
}

func getConsoleScreenBufferInfo(hCon syscall.Handle) (sb consoleScreenBuffer, err error) {
	rc, _, ec := syscall.Syscall(procGetConScrBufInfo.Addr(), 2,
		uintptr(hCon), uintptr(unsafe.Pointer(&sb)), 0)
	if rc == 0 {
		err = syscall.Errno(ec)
	}
	return
}

// GetConsoleWidth returns the terminal width in characters
func GetConsoleWidth() int {
	hCon, err := syscall.Open("CONOUT$", syscall.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	defer syscall.Close(hCon)

	sb, err := getConsoleScreenBufferInfo(hCon)
	if err != nil {
		panic(err)
	}
	return int(sb.size.x)
}
