package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ghodss/yaml"
	"github.com/yargevad/filepathx"
)

var (
	yaml1  bool
	yaml2  bool
	spaced bool
	csv    bool
	ilp    bool
)

func main() {
	flag.BoolVar(&yaml1, "yh", false, "yaml hash using PID as keys")
	flag.BoolVar(&yaml2, "ya", false, "yaml array")
	flag.BoolVar(&spaced, "s", false, "spaced format")
	flag.BoolVar(&csv, "c", false, "csv format")
	flag.BoolVar(&ilp, "i", false, "influx line protocol")
	flag.Parse()

	if !yaml1 && !yaml2 && !spaced && !csv && !ilp {
		csv = true
	}

	files, _ := filepathx.Glob(filepath.Join("/proc", "*", "status"))

	switch {
	case spaced:
		fmt.Println("PID", "NAME", "STATE", "VMRSS", "VMHWM", "VMSWAP")
	case csv:
		fmt.Println("Pid,Name,State,VmRSS,VmHWM,VmSwap")
	}

	wg := &sync.WaitGroup{}
	tokens := make(chan struct{}, 10)

	for _, file := range files {
		wg.Add(1)
		tokens <- struct{}{}

		go func(file string) {
			defer func() {
				wg.Done()
				<-tokens
			}()
			m := make(map[string]interface{})
			data, err := os.ReadFile(file)
			if err != nil {
				return
			}

			if err := yaml.Unmarshal(data, &m); err != nil {
				return
			}

			if m["VmRSS"] == nil {
				m["VmRSS"] = "0"
			}

			if m["VmHWM"] == nil {
				m["VmHWM"] = "0"
			}

			if m["VmSwap"] == nil {
				m["VmSwap"] = "0"
			}

			m["VmRSS"] = strings.ReplaceAll(m["VmRSS"].(string), " kB", "")
			m["VmHWM"] = strings.ReplaceAll(m["VmHWM"].(string), " kB", "")
			m["VmSwap"] = strings.ReplaceAll(m["VmSwap"].(string), " kB", "")
			state := string(m["State"].(string)[0])

			switch {
			case yaml1:
				fmt.Printf("%0.f: {name: \"%s\", state: %s, rss: %s, rss_hwm: %s, swap: %s }\n", m["Pid"], m["Name"], state, m["VmRSS"], m["VmHWM"], m["VmSwap"])
			case yaml2:
				fmt.Printf("- { pid: %0.f, name: \"%s\", state: %s, rss: %s, rss_hwm: %s, swap: %s }\n", m["Pid"], m["Name"], state, m["VmRSS"], m["VmHWM"], m["VmSwap"])
			case spaced:
				fmt.Println(m["Pid"], strings.ReplaceAll(m["Name"].(string), " ", "_"), state, m["VmRSS"], m["VmHWM"], m["VmSwap"])
			case csv:
				fmt.Printf("%0.f,%s,%s,%s,%s,%s\n", m["Pid"], m["Name"], state, m["VmRSS"], m["VmHWM"], m["VmSwap"])
			case ilp:
				fmt.Printf("stat,name=\"%s\",state=%s rss=%s,rss_hwm=%s,swap=%s %d\n", m["Name"], state, m["VmRSS"], m["VmHWM"], m["VmSwap"], time.Now().Unix())
			}
		}(file)
	}
}
