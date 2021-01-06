// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"flag"
	"fmt"
	"sort"

	"golang.org/x/tools/internal/lsp/protocol"
	"golang.org/x/tools/internal/span"
	"golang.org/x/tools/internal/tool"
)

type completion struct {
	app *Application
}

func (c completion) Name() string      { return "completion" }
func (c completion) Usage() string     { return "<position>" }
func (c completion) ShortHelp() string { return "display selected identifier completions" }
func (c completion) DetailedHelp(f *flag.FlagSet) {
	fmt.Fprintf(f.Output(), `
Example:

   $ gopls completion helper/helper.go:21:13
   $ gopls completion helper/helper.go:#292

   gopls completion flags are:

`)
}

func (c completion) Run(ctx context.Context, args ...string) error {
	if len(args) != 1 {
		return tool.CommandLineErrorf("completion expects 1 argument (file)")
	}

	conn, err := c.app.connect(ctx)
	if err != nil {
		return err
	}
	defer conn.terminate(ctx)

	from := span.Parse(args[0])
	file := conn.AddFile(ctx, from.URI())

	if file.err != nil {
		return file.err
	}

	loc, err := file.mapper.Location(from)
	if err != nil {
		return err
	}

	p := protocol.CompletionParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{
				URI: loc.URI,
			},
			Position: loc.Range.Start,
		},
	}

	completion, err := conn.Completion(ctx, &p)
	if err != nil {
		return err
	}

	completions := make([]string, 0, len(completion.Items))
	for _, comp := range completion.Items {
		completions = append(completions, fmt.Sprintf("%s %s", comp.Label, comp.Detail))
	}
	sort.Strings(completions)

	for _, c := range completions {
		fmt.Println(c)
	}

	return nil
}
