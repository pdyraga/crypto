// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.7,amd64,!gccgo,!appengine

package blake2b

import (
	"encoding/hex"
	"fmt"

	"golang.org/x/sys/cpu"
)

func init() {
	useAVX2 = cpu.X86.HasAVX2
	useAVX = cpu.X86.HasAVX
	useSSE4 = cpu.X86.HasSSE41
}

var Counter = 0

//go:noescape
func hashBlocksAVX2(h *[8]uint64, c *[2]uint64, flag uint64, blocks []byte)

//go:noescape
func hashBlocksAVX(h *[8]uint64, c *[2]uint64, flag uint64, blocks []byte)

//go:noescape
func hashBlocksSSE4(h *[8]uint64, c *[2]uint64, flag uint64, blocks []byte)

func hashBlocks(h *[8]uint64, c *[2]uint64, flag uint64, blocks []byte) {
	// here

	fmt.Printf("// %v\n", Counter)
	Counter = Counter + 1
	fmt.Printf("{\n")
	fmt.Printf("mIn: \"%v\",\n", hex.EncodeToString(blocks))
	fmt.Printf("hIn: [8]uint64{")
	for hi, hh := range h {
		if hi == len(h)-1 {
			fmt.Printf("0x%x", hh)
		} else {
			fmt.Printf("0x%x, ", hh)
		}
	}
	fmt.Printf("},\n")
	fmt.Printf("c: [2]uint64{0x%x, 0x%x},\n", c[0], c[1])
	fmt.Printf("f: 0x%x,\n", flag)
	fmt.Printf("rounds: 12,\n")

	switch {
	case useAVX2:
		hashBlocksAVX2(h, c, flag, blocks)
	case useAVX:
		hashBlocksAVX(h, c, flag, blocks)
	case useSSE4:
		hashBlocksSSE4(h, c, flag, blocks)
	default:
		hashBlocksGeneric(h, c, flag, blocks)
	}

	fmt.Printf("hOut: [8]uint64{")
	for hi, hh := range h {
		if hi == len(h)-1 {
			fmt.Printf("0x%x", hh)
		} else {
			fmt.Printf("0x%x, ", hh)
		}
	}
	fmt.Printf("},\n")
	fmt.Printf("},\n")
}
