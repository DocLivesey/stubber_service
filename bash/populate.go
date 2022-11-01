package bash

import (
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
		return nil, err
	}

	lines := strings.Split(string(out), "\n")

	out, err = exec.Command("ss", "-panlt").CombinedOutput()
	if err != nil {
		return nil, err
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
	// if err != nil {
	// 	log.Fatalf("Cannot get workdir:  %s", err)
	// }

	var dirRun func(string) error

	dirRun = func(path string) error {
		fileInfo, err := os.Stat(path)
		if err != nil {
			return err
		}
		if fileInfo.IsDir() {
			dirs, err := os.ReadDir(path)
			if err != nil {
				return err
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
		return nil
	}
	err = dirRun(dir)
	if err != nil {
		return nil, err
	}
	return stubs, nil
}

func StubStatus(s *serv.Stub) error {
	out, err := exec.Command("ps", "-e", "-o", "pid,pcpu,pmem,command").CombinedOutput()
	if err != nil {
		return err
	}

	lines := strings.Split(string(out), "\n")

	out, err = exec.Command("ss", "-panlt").CombinedOutput()
	if err != nil {
		return err
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

	fileInfo, err := os.Stat(s.Path)
	if err != nil {
		return err
	}
	n := fileInfo.Name()
	state := false
	*s = serv.Stub{Jar: &n, Path: s.Path, State: &state}
	for _, l := range lines {
		if strings.Contains(l, s.Path) {
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

	return nil
}
