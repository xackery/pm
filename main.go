package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/shirou/gopsutil/process"
)

type ProcessInfo struct {
	PID        int32
	Username   string
	CPUPercent float64
	MemPercent float32
	MemKB      uint64
	Command    string
}

type ByMemUsageDesc []ProcessInfo

func (a ByMemUsageDesc) Len() int           { return len(a) }
func (a ByMemUsageDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByMemUsageDesc) Less(i, j int) bool { return a[i].MemKB > a[j].MemKB }

func main() {
	processes, err := process.Processes()
	if err != nil {
		log.Fatal(err)
	}

	// Collect process information
	var processInfoList []ProcessInfo

	for _, p := range processes {
		pid := p.Pid
		username, err := p.Username()
		if err != nil {
			continue
		}

		cpuPercent, err := p.CPUPercent()
		if err != nil {
			continue
		}

		memPercent, err := p.MemoryPercent()
		if err != nil {
			continue
		}

		memInfo, err := p.MemoryInfo()
		if err != nil {
			continue
		}

		cmdline, err := p.Cmdline()
		if err != nil {
			continue
		}

		cmd := cmdline
		if idx := strings.Index(cmdline, " "); idx != -1 {
			cmd = cmdline[:idx]
		}

		processInfoList = append(processInfoList, ProcessInfo{
			PID:        pid,
			Username:   username,
			CPUPercent: cpuPercent,
			MemPercent: memPercent,
			MemKB:      memInfo.RSS / 1024,
			Command:    cmd,
		})
	}

	// Sort by memory usage in descending order
	sort.Sort(ByMemUsageDesc(processInfoList))

	pidLen := 0
	memLen := 0
	memPctLen := 0
	userLen := 0

	if len("PID") > pidLen {
		pidLen = len("PID")
	}

	if len("MEM(MB)") > memLen {
		memLen = len("MEM(MB)")
	}

	if len("%MEM") > memPctLen {
		memPctLen = len("%MEM")
	}

	if len("USER") > userLen {
		userLen = len("USER")
	}

	for _, p := range processInfoList {
		if len(fmt.Sprintf("%d", p.PID)) > pidLen {
			pidLen = len(fmt.Sprintf("%d", p.PID))
		}

		if len(fmt.Sprintf("%dMB", p.MemKB/1024)) > memLen {
			memLen = len(fmt.Sprintf("%dMB", p.MemKB/1024))
		}

		if len(fmt.Sprintf("%.1f", p.MemPercent)) > memPctLen {
			memPctLen = len(fmt.Sprintf("%.1f", p.MemPercent))
		}

		if len(p.Username) > userLen {
			userLen = len(p.Username)
		}

	}

	format := fmt.Sprintf("%%-%dd %%-%ds %%-%d.1f %%-%ds %%s\n", pidLen, memLen, memPctLen, userLen)
	headerFormat := fmt.Sprintf("%%-%ds %%-%ds %%-%ds %%-%ds %%s\n", pidLen, memLen, memPctLen, userLen)

	fmt.Printf(headerFormat, "PID", "MEM(MB)", "%MEM", "USER", "COMMAND")
	for _, p := range processInfoList {
		if p.Command == "" {
			continue
		}
		if p.Command == "(sd-pam)" {
			continue
		}

		if p.MemKB/1024 == 0 {
			continue
		}

		fmt.Printf(format, p.PID, fmt.Sprintf("%dMB", p.MemKB/1024), p.MemPercent, p.Username, p.Command)
	}
}
