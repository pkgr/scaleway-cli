// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/units"

	types "github.com/scaleway/scaleway-cli/commands/types"
	utils "github.com/scaleway/scaleway-cli/utils"
)

var cmdHistory = &types.Command{
	Exec:        runHistory,
	UsageLine:   "history [OPTIONS] IMAGE",
	Description: "Show the history of an image",
	Help:        "Show the history of an image.",
}

func init() {
	cmdHistory.Flag.BoolVar(&historyNoTrunc, []string{"-no-trunc"}, false, "Don't truncate output")
	cmdHistory.Flag.BoolVar(&historyQuiet, []string{"q", "-quiet"}, false, "Only show numeric IDs")
	cmdHistory.Flag.BoolVar(&historyHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var historyNoTrunc bool // --no-trunc flag
var historyQuiet bool   // -q, --quiet flag
var historyHelp bool    // -h, --help flag

func runHistory(cmd *types.Command, args []string) {
	if historyHelp {
		cmd.PrintUsage()
	}
	if len(args) != 1 {
		cmd.PrintShortUsage()
	}

	imageID := cmd.API.GetImageID(args[0], true)
	image, err := cmd.API.GetImage(imageID)
	if err != nil {
		log.Fatalf("Cannot get image %s: %v", imageID, err)
	}

	if imagesQ {
		fmt.Println(imageID)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 10, 1, 3, ' ', 0)
	defer w.Flush()
	fmt.Fprintf(w, "IMAGE\tCREATED\tCREATED BY\tSIZE\n")

	identifier := utils.TruncIf(image.Identifier, 8, !historyNoTrunc)

	creationDate, err := time.Parse("2006-01-02T15:04:05.000000+00:00", image.CreationDate)
	if err != nil {
		log.Fatalf("Unable to parse creation date from the Scaleway API: %v", err)
	}
	creationDateStr := units.HumanDuration(time.Now().UTC().Sub(creationDate))

	volumeName := utils.TruncIf(image.RootVolume.Name, 25, !historyNoTrunc)
	size := units.HumanSize(float64(image.RootVolume.Size))

	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", identifier, creationDateStr, volumeName, size)
}
