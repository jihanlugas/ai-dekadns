package helper

import (
	"bytes"
	"fmt"
	"os/exec"
)

func OpenSSL(args ...string) ([]byte, error) {
	cmd := exec.Command("openssl", args...)

	out := &bytes.Buffer{}
	errs := &bytes.Buffer{}

	cmd.Stdout, cmd.Stderr = out, errs

	if err := cmd.Run(); err != nil {
		if len(errs.Bytes()) > 0 {
			return nil, fmt.Errorf("error running %s (%s):\n %v", cmd.Args, err, errs.String())
		}
		return nil, err
	}

	return out.Bytes(), nil
}
