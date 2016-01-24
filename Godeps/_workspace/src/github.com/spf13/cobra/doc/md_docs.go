//Copyright 2015 Red Hat Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package doc

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/csats/kubeops/Godeps/_workspace/src/github.com/spf13/cobra"
)

func printOptions(w io.Writer, cmd *cobra.Command, name string) error {
	flags := cmd.NonInheritedFlags()
	flags.SetOutput(w)
	if flags.HasFlags() {
		if _, err := fmt.Fprintf(w, "### Options\n\n```\n"); err != nil {
			return err
		}
		flags.PrintDefaults()
		if _, err := fmt.Fprintf(w, "```\n\n"); err != nil {
			return err
		}
	}

	parentFlags := cmd.InheritedFlags()
	parentFlags.SetOutput(w)
	if parentFlags.HasFlags() {
		if _, err := fmt.Fprintf(w, "### Options inherited from parent commands\n\n```\n"); err != nil {
			return err
		}
		parentFlags.PrintDefaults()
		if _, err := fmt.Fprintf(w, "```\n\n"); err != nil {
			return err
		}
	}
	return nil
}

func GenMarkdown(cmd *cobra.Command, w io.Writer) error {
	return GenMarkdownCustom(cmd, w, func(s string) string { return s })
}

func GenMarkdownCustom(cmd *cobra.Command, w io.Writer, linkHandler func(string) string) error {
	name := cmd.CommandPath()

	short := cmd.Short
	long := cmd.Long
	if len(long) == 0 {
		long = short
	}

	if _, err := fmt.Fprintf(w, "## %s\n\n", name); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "%s\n\n", short); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "### Synopsis\n\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "\n%s\n\n", long); err != nil {
		return err
	}

	if cmd.Runnable() {
		if _, err := fmt.Fprintf(w, "```\n%s\n```\n\n", cmd.UseLine()); err != nil {
			return err
		}
	}

	if len(cmd.Example) > 0 {
		if _, err := fmt.Fprintf(w, "### Examples\n\n"); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "```\n%s\n```\n\n", cmd.Example); err != nil {
			return err
		}
	}

	if err := printOptions(w, cmd, name); err != nil {
		return err
	}
	if hasSeeAlso(cmd) {
		if _, err := fmt.Fprintf(w, "### SEE ALSO\n"); err != nil {
			return err
		}
		if cmd.HasParent() {
			parent := cmd.Parent()
			pname := parent.CommandPath()
			link := pname + ".md"
			link = strings.Replace(link, " ", "_", -1)
			if _, err := fmt.Fprintf(w, "* [%s](%s)\t - %s\n", pname, linkHandler(link), parent.Short); err != nil {
				return err
			}
			cmd.VisitParents(func(c *cobra.Command) {
				if c.DisableAutoGenTag {
					cmd.DisableAutoGenTag = c.DisableAutoGenTag
				}
			})
		}

		children := cmd.Commands()
		sort.Sort(byName(children))

		for _, child := range children {
			if !child.IsAvailableCommand() || child.IsHelpCommand() {
				continue
			}
			cname := name + " " + child.Name()
			link := cname + ".md"
			link = strings.Replace(link, " ", "_", -1)
			if _, err := fmt.Fprintf(w, "* [%s](%s)\t - %s\n", cname, linkHandler(link), child.Short); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(w, "\n"); err != nil {
			return err
		}
	}
	if !cmd.DisableAutoGenTag {
		if _, err := fmt.Fprintf(w, "###### Auto generated by spf13/cobra on %s\n", time.Now().Format("2-Jan-2006")); err != nil {
			return err
		}
	}
	return nil
}

func GenMarkdownTree(cmd *cobra.Command, dir string) error {
	identity := func(s string) string { return s }
	emptyStr := func(s string) string { return "" }
	return GenMarkdownTreeCustom(cmd, dir, emptyStr, identity)
}

func GenMarkdownTreeCustom(cmd *cobra.Command, dir string, filePrepender, linkHandler func(string) string) error {
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsHelpCommand() {
			continue
		}
		if err := GenMarkdownTreeCustom(c, dir, filePrepender, linkHandler); err != nil {
			return err
		}
	}

	basename := strings.Replace(cmd.CommandPath(), " ", "_", -1) + ".md"
	filename := filepath.Join(dir, basename)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.WriteString(f, filePrepender(filename)); err != nil {
		return err
	}
	if err := GenMarkdownCustom(cmd, f, linkHandler); err != nil {
		return err
	}
	return nil
}
