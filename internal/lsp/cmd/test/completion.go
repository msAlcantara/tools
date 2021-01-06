// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmdtest

import (
	"fmt"
	"sort"
	"testing"

	"golang.org/x/tools/internal/lsp/tests"
	"golang.org/x/tools/internal/span"
)

func (r *runner) Completion(t *testing.T, spn span.Span, test tests.Completion, items tests.CompletionItems) {
	itemString := make([]string, 0, len(test.CompletionItems))
	for _, pos := range test.CompletionItems {
		item := tests.ToProtocolCompletionItem(*items[pos])
		itemString = append(itemString, fmt.Sprintf("%s %s", item.Label, item.Detail))
	}
	sort.Strings(itemString)
	var expect string
	for _, i := range itemString {
		expect += i + "\n"
	}
	expect = r.Normalize(expect)

	uri := spn.URI()
	filename := uri.Filename()
	target := fmt.Sprintf("%s:%v:%v", filename, spn.Start().Line(), spn.Start().Column())

	got, stderr := r.NormalizeGoplsCmd(t, "completion", target)

	if stderr != "" {
		// t.Errorf("implementation failed for %s: %s", target, stderr)
		// fmt.Printf("STDERR: %s\n\ngot: %s\n", stderr, got)
	} else if expect != got {
		t.Errorf("implementation failed for %s expected:\n%s\ngot:\n%s", target, expect, got)
	}
}
