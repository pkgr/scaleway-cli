// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"log"
	"sort"
	"strings"

	types "github.com/scaleway/scaleway-cli/commands/types"
	utils "github.com/scaleway/scaleway-cli/utils"
)

var cmdCompletion = &types.Command{
	Exec:        runCompletion,
	UsageLine:   "_completion [OPTIONS] CATEGORY",
	Description: "Completion helper",
	Help:        "Completion helper.",
	Hidden:      true,
	Examples: `
    $ scw _completion servers-all
    $ scw _completion images-all
    $ scw _completion snapshots-all
    $ scw _completion volumes-all
    $ scw _completion bootscripts-all
    $ scw _completion servers-names
    $ scw _completion images-names
    $ scw _completion snapshots-names
    $ scw _completion volumes-names
    $ scw _completion bootscripts-names
`,
}

func init() {
	cmdCompletion.Flag.BoolVar(&completionHelp, []string{"h", "-help"}, false, "Print usage")
	cmdCompletion.Flag.BoolVar(&completionPrefix, []string{"-prefix"}, false, "Prefix entries")
}

// Flags
var completionHelp bool   // -h, --help flag
var completionPrefix bool // --prefix flag

func wordifyName(name string, kind string) string {
	ret := ""
	if completionPrefix {
		ret += kind + "\\:"
	}
	ret += utils.Wordify(name)
	return ret
}

func runCompletion(cmd *types.Command, args []string) {
	if completionHelp {
		cmd.PrintUsage()
	}
	if len(args) != 1 {
		cmd.PrintShortUsage()
	}

	category := args[0]

	elements := []string{}

	switch category {
	case "servers-all":
		for identifier, name := range cmd.API.Cache.Servers {
			elements = append(elements, identifier, wordifyName(name, "server"))
		}
	case "servers-names":
		for _, name := range cmd.API.Cache.Servers {
			elements = append(elements, wordifyName(name, "server"))
		}
	case "images-all":
		for identifier, name := range cmd.API.Cache.Images {
			elements = append(elements, identifier, wordifyName(name, "image"))
		}
	case "images-names":
		for _, name := range cmd.API.Cache.Images {
			elements = append(elements, wordifyName(name, "image"))
		}
	case "volumes-all":
		for identifier, name := range cmd.API.Cache.Volumes {
			elements = append(elements, identifier, wordifyName(name, "volume"))
		}
	case "volumes-names":
		for _, name := range cmd.API.Cache.Volumes {
			elements = append(elements, wordifyName(name, "volume"))
		}
	case "snapshots-all":
		for identifier, name := range cmd.API.Cache.Snapshots {
			elements = append(elements, identifier, wordifyName(name, "snapshot"))
		}
	case "snapshots-names":
		for _, name := range cmd.API.Cache.Snapshots {
			elements = append(elements, wordifyName(name, "snapshot"))
		}
	case "bootscripts-all":
		for identifier, name := range cmd.API.Cache.Bootscripts {
			elements = append(elements, identifier, wordifyName(name, "bootscript"))
		}
	case "bootscripts-names":
		for _, name := range cmd.API.Cache.Bootscripts {
			elements = append(elements, wordifyName(name, "bootscript"))
		}
	default:
		log.Fatalf("Unhandled category of completion: %s", category)
	}

	sort.Strings(elements)
	fmt.Println(strings.Join(utils.RemoveDuplicates(elements), "\n"))
}
