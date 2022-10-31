package bash

import (
	"log"
	"os"
	"os/exec"
	"strings"

	serv "github.com/DocLivesey/stubber_service/gen_service"
)

func between(value string, a string, b string) string {
	// Get substring between two strings.
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}

func Populate() ([]serv.Stub, error) { //[]main.Stub {
	out, err := exec.Command("ps", "-e", "-o", "pid,pcpu,pmem,command").CombinedOutput()
	if err != nil {
		log.Fatalf("Error on ps call: %s", err)
	}

	lines := strings.Split(string(out), "\n")

	out, err = exec.Command("ss", "-panlt").CombinedOutput()
	if err != nil {
		log.Fatalf("Error on ss call: %s", err)
	}

	sslines := strings.Split(string(out), "\n")

	var ports = make(map[string]string, 20)

	for i, l := range sslines {
		if i == 0 {
			continue
		}
		t := strings.Fields(l)
		if len(t) < 6 {
			continue
		}
		tmp := strings.Split(t[3], ":")
		port := tmp[1]
		pid := between(t[5], "pid=", ",f")
		ports[pid] = port
	}

	var stubs []serv.Stub //[]main.Stub

	dir := "/home/kuro/dev/tmp" //, err := os.Getwd()
	if err != nil {
		log.Fatalf("Cannot get workdir:  %s", err)
	}

	var dirRun func(string)

	dirRun = func(path string) {
		fileInfo, err := os.Stat(path)
		if err != nil {
			log.Fatalf("FileStat error\n\t%s", err)
		}
		if fileInfo.IsDir() {
			dirs, err := os.ReadDir(path)
			if err != nil {
				log.Fatalf("Error on reading directory: %s", err)
			}

			for _, d := range dirs {
				p := path + "/" + d.Name()
				dirRun(p)
				// println(d.Name())
			}
		} else {
			if strings.Contains(fileInfo.Name(), ".jar") {
				n := fileInfo.Name()
				state := false
				s := serv.Stub{Jar: &n, Path: path, State: &state}
				for _, l := range lines {
					if strings.Contains(l, path) {
						fl := strings.Fields(l)
						s.Pid = &fl[0]
						state = true
						s.State = &state
						s.Cpu = &fl[1]
						s.Mem = &fl[2]
						i := ports[fl[0]]
						s.Port = &i
						break
					}
				}
				stubs = append(stubs, s)
			}
		}
	}
	dirRun(dir)
	return stubs, nil
}
