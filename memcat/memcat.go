package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mozilla/masche/memaccess"
	"github.com/mozilla/masche/process"
)

// TODO:
// - Add a flag to dump memory in hex
// - Add a flag to dump memory in ascii
// - Add a flag to dump memory in binary
// - Add a flag to dump memory in a more human readable format,
// 		like braile dots encoding of a byte
// - Add a flag to dump memory in a more ordered representation of the memory location,
// 		so as to be able to find the memory location of a given byte in the dump
// - ???

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func stdCheck(softs []error, hard error) {
	for _, soft := range softs {
		check(soft)
	}
	check(hard)
}

func muxErr(softs []error, hard error) error {
	for _, soft := range softs {
		if soft != nil {
			return soft
		}
	}
	return hard
}

func availableProgs(quiet bool) (names []string, pids []uint) {
	allpids, softs, hard := process.GetAllPids()
	names = make([]string, 0, len(allpids))
	pids = make([]uint, 0, len(allpids))
	stdCheck(softs, hard)
	for _, pid := range allpids {
		if !quiet {
			fmt.Printf("%d\n", pid)
		}
		proc, softs, hard := process.OpenFromPid(pid)
		err := muxErr(softs, hard)
		if err != nil {
			if !quiet {
				fmt.Printf("%s\n", err)
			}
			continue
		}
		name, softs, hard := proc.Name()
		err = muxErr(softs, hard)
		if err != nil {
			if !quiet {
				fmt.Printf("%s\n", err)
			}
			continue
		}
		if !quiet {
			fmt.Println(name)
		}
		names = append(names, name)
		pids = append(pids, pid)
		proc.Close()
	}
	return
}

type WrappedProcess struct {
	process.Process
}

func (p WrappedProcess) DumpMem() (memoryScan chan []byte, err error) {
	memoryScan = make(chan []byte, 1000)
	var last uintptr = 0
	go func() {
		count := 0
		err = muxErr(memaccess.WalkMemory(p, last, 4096<<2, func(addr uintptr, buf []byte) bool {
			if len(buf) == 0 || last == addr {
				return false
			}
			last = addr
			memoryScan <- append([]byte{}, buf...)
			count++
			return true
		}))
		for len(memoryScan) > 0 {
			time.Sleep(100 * time.Millisecond)
		}
		close(memoryScan)
		fmt.Println("done scanning")
	}()
	return
}

func main() {
	names, pids := availableProgs(false)
	chosen := pids[len(pids)-1]
	for i, name := range names {
		if strings.Contains(name, os.Args[1]) {
			chosen = pids[i]
		}
	}
	p, softs, hard := process.OpenFromPid(chosen)
	if muxErr(softs, hard) != nil {
		fmt.Println("Error opening process")
		return
	}
	defer p.Close()
	wp := WrappedProcess{p}
	mem, err := wp.DumpMem()
	if err != nil {
		fmt.Println("Error dumping memory")
		return
	}
	for buf := range mem {
		os.Stdout.Write(buf)
	}
}
