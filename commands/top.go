// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"

	types "github.com/scaleway/scaleway-cli/commands/types"
	"github.com/scaleway/scaleway-cli/utils"
)

var cmdTop = &types.Command{
	Exec:        runTop,
	UsageLine:   "top [OPTIONS] SERVER", // FIXME: add ps options
	Description: "Lookup the running processes of a server",
	Help:        "Lookup the running processes of a server.",
}

func init() {
	cmdTop.Flag.BoolVar(&topHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var topHelp bool // -h, --help flag

func runTop(cmd *types.Command, args []string) {
	if topHelp {
		cmd.PrintUsage()
	}
	if len(args) != 2 {
		cmd.PrintShortUsage()
	}

	serverID := cmd.API.GetServerID(args[0])
	command := "ps"
	server, err := cmd.API.GetServer(serverID)
	if err != nil {
		log.Fatalf("Failed to get server information for %s: %v", serverID, err)
	}

	execCmd := append(utils.NewSSHExecCmd(server.PublicAddress.IP, true, []string{command}))

	log.Debugf("Executing: ssh %s", strings.Join(execCmd, " "))
	out, err := exec.Command("ssh", execCmd...).CombinedOutput()
	fmt.Printf("%s", out)
	if err != nil {
		log.Fatal(err)
	}
}
