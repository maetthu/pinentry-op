package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/foxcpp/go-assuan/common"
	"github.com/foxcpp/go-assuan/pinentry"
)

func main() {
	envVar := "OP_PINENTRY_SECRET"

	if os.Getenv(envVar) == "" {
		log.Fatal(fmt.Sprintf("%s environment variable not set", envVar))
	}

	cmd := exec.Command("op", "read", "-n", os.Getenv(envVar))

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("1Password CLI exited non-zero (%s) [%s]", err, stderr.String())
	}

	callbacks := pinentry.Callbacks{
		GetPIN: func(settings pinentry.Settings) (string, *common.Error) {
			return stdout.String(), nil
		},
		Confirm: func(settings pinentry.Settings) (bool, *common.Error) {
			return true, nil
		},
		Msg: func(settings pinentry.Settings) *common.Error {
			return nil
		},
	}

	// Add SETKEYINFO handler for openfortivpn compatibility
	pinentry.ProtoInfo.Handlers["SETKEYINFO"] = func(pipe io.ReadWriter, state interface{}, params string) *common.Error {
		return nil
	}

	err := pinentry.Serve(callbacks, "Hello")

	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
}
