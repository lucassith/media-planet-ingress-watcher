package kubectl

import (
	"bytes"
	"io"
	"os/exec"
	"strings"
)

// MakeIngressFile creates new ingress yaml template and writes it to the output.
func MakeIngressFile(hostname string, output io.Writer) (int, error) {
	template := GetTemplate()
	template = strings.Replace(template, TemplateValue, hostname, 3)
	return output.Write([]byte(template))
}

// ExecuteKubectl runs 'kubectl apply -f' command on provided file.
func ExecuteKubectl(file string) (string, error) {
	command := exec.Command("kc", "apply", "-f", file)
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	if err != nil {
		return string(stderr.Bytes()), err
	} else {
		return string(stdout.Bytes()), nil
	}
}
