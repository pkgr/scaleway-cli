// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"os"
	"sort"
	"sync"
	"text/tabwriter"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/units"

	"github.com/scaleway/scaleway-cli/api"
	types "github.com/scaleway/scaleway-cli/commands/types"
	"github.com/scaleway/scaleway-cli/utils"
)

var cmdImages = &types.Command{
	Exec:        runImages,
	UsageLine:   "images [OPTIONS]",
	Description: "List images",
	Help:        "List images.",
}

func init() {
	cmdImages.Flag.BoolVar(&imagesA, []string{"a", "-all"}, false, "Show all iamges")
	cmdImages.Flag.BoolVar(&imagesNoTrunc, []string{"-no-trunc"}, false, "Don't truncate output")
	cmdImages.Flag.BoolVar(&imagesQ, []string{"q", "-quiet"}, false, "Only show numeric IDs")
	cmdImages.Flag.BoolVar(&imagesHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var imagesA bool       // -a flag
var imagesQ bool       // -q flag
var imagesNoTrunc bool // -no-trunc flag
var imagesHelp bool    // -h, --help flag

func runImages(cmd *types.Command, args []string) {
	if imagesHelp {
		cmd.PrintUsage()
	}
	if len(args) != 0 {
		cmd.PrintShortUsage()
	}

	wg := sync.WaitGroup{}
	chEntries := make(chan api.ScalewayImageInterface)
	var entries = []api.ScalewayImageInterface{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		images, err := cmd.API.GetImages()
		if err != nil {
			log.Fatalf("unable to fetch images from the Scaleway API: %v", err)
		}
		for _, val := range *images {
			creationDate, err := time.Parse("2006-01-02T15:04:05.000000+00:00", val.CreationDate)
			if err != nil {
				log.Fatalf("unable to parse creation date from the Scaleway API: %v", err)
			}
			chEntries <- api.ScalewayImageInterface{
				Type:         "image",
				CreationDate: creationDate,
				Identifier:   val.Identifier,
				Name:         val.Name,
				Public:       val.Public,
				Tag:          "latest",
				VirtualSize:  float64(val.RootVolume.Size),
			}
		}
	}()

	if imagesA {
		wg.Add(1)
		go func() {
			defer wg.Done()
			snapshots, err := cmd.API.GetSnapshots()
			if err != nil {
				log.Fatalf("unable to fetch snapshots from the Scaleway API: %v", err)
			}
			for _, val := range *snapshots {
				creationDate, err := time.Parse("2006-01-02T15:04:05.000000+00:00", val.CreationDate)
				if err != nil {
					log.Fatalf("unable to parse creation date from the Scaleway API: %v", err)
				}
				chEntries <- api.ScalewayImageInterface{
					Type:         "snapshot",
					CreationDate: creationDate,
					Identifier:   val.Identifier,
					Name:         val.Name,
					Tag:          "<snapshot>",
					VirtualSize:  float64(val.Size),
					Public:       false,
				}
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			bootscripts, err := cmd.API.GetBootscripts()
			if err != nil {
				log.Fatalf("unable to fetch bootscripts from the Scaleway API: %v", err)
			}
			for _, val := range *bootscripts {
				chEntries <- api.ScalewayImageInterface{
					Type:       "bootscript",
					Identifier: val.Identifier,
					Name:       val.Title,
					Tag:        "<bootscript>",
					Public:     false,
				}
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			volumes, err := cmd.API.GetVolumes()
			if err != nil {
				log.Fatalf("unable to fetch volumes from the Scaleway API: %v", err)
			}
			for _, val := range *volumes {
				creationDate, err := time.Parse("2006-01-02T15:04:05.000000+00:00", val.CreationDate)
				if err != nil {
					log.Fatalf("unable to parse creation date from the Scaleway API: %v", err)
				}
				chEntries <- api.ScalewayImageInterface{
					Type:         "volume",
					CreationDate: creationDate,
					Identifier:   val.Identifier,
					Name:         val.Name,
					Tag:          "<volume>",
					VirtualSize:  float64(val.Size),
					Public:       false,
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(chEntries)
	}()

	done := false
	for {
		select {
		case entry, ok := <-chEntries:
			if !ok {
				done = true
				break
			}
			entries = append(entries, entry)
		}
		if done {
			break
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
	defer w.Flush()
	if !imagesQ {
		fmt.Fprintf(w, "REPOSITORY\tTAG\tIMAGE ID\tCREATED\tVIRTUAL SIZE\n")
	}
	sort.Sort(api.ByCreationDate(entries))
	for _, image := range entries {
		if imagesQ {
			fmt.Fprintf(w, "%s\n", image.Identifier)
		} else {
			tag := image.Tag
			shortID := utils.TruncIf(image.Identifier, 8, !imagesNoTrunc)
			name := utils.Wordify(image.Name)
			if !image.Public && image.Type == "image" {
				name = "user/" + name
			}
			shortName := utils.TruncIf(name, 25, !imagesNoTrunc)
			var creationDate, virtualSize string
			if image.CreationDate.IsZero() {
				creationDate = "n/a"
			} else {
				creationDate = units.HumanDuration(time.Now().UTC().Sub(image.CreationDate))
			}
			if image.VirtualSize == 0 {
				virtualSize = "n/a"
			} else {
				virtualSize = units.HumanSize(image.VirtualSize)
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", shortName, tag, shortID, creationDate, virtualSize)
		}
	}
}
