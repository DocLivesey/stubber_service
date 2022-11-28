package bash

import (
	"os/exec"

	"github.com/DocLivesey/stubber_service/gen_service"
)

var java = "java"

func StopStub(stub gen_service.Stub) error {
	if err := exec.Command("kill", *stub.Pid).Run(); err != nil {
		return err
	}
	return nil
}

func StartStub(stub gen_service.Stub) error {
	cmd := exec.Command(java, "-Xmx1G", "-jar", stub.Path, "&> /dev/null")

	if err := cmd.Start(); err != nil {
		return err
	}
	return nil
}
