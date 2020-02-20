// Lute - A structured markdown engine.
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

package main

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/88250/lute"
	gm "github.com/gomarkdown/markdown"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"gitlab.com/golang-commonmark/markdown"
	"gopkg.in/russross/blackfriday.v2"
)

const spec = "commonmark-spec"

func BenchmarkLute(b *testing.B) {
	bytes, err := ioutil.ReadFile(spec + ".md")
	if nil != err {
		b.Fatalf("read spec text failed: " + err.Error())
	}

	luteEngine := lute.New()
	luteEngine.GFMTaskListItem = true
	luteEngine.GFMTable = true
	luteEngine.GFMAutoLink = true
	luteEngine.GFMStrikethrough = true
	luteEngine.SoftBreak2HardBreak = false
	luteEngine.CodeSyntaxHighlight = false
	luteEngine.Footnotes = false
	luteEngine.ToC = false
	luteEngine.AutoSpace = false
	luteEngine.FixTermTypo = false
	luteEngine.ChinesePunct = false
	luteEngine.Emoji = false
	ht, er1 := luteEngine.Markdown("spec text", bytes)
	if nil != er1 {
		b.Fatalf("unexpected: %s", er1)
	}
	if err := ioutil.WriteFile(spec+".html", ht, 0644); nil != err {
		b.Fatalf("write spec html failed: %s", err)
	}

	b.SetParallelism(8)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			luteEngine.Markdown("spec text", bytes)
		}
	})

}

func BenchmarkGolangCommonMark(b *testing.B) {
	buf, err := ioutil.ReadFile(spec + ".md")
	if nil != err {
		b.Fatalf("read spec text failed: " + err.Error())
	}

	md := markdown.New(markdown.XHTMLOutput(true),
		markdown.Tables(true),
		markdown.Linkify(true),
		markdown.Typographer(false))

	b.SetParallelism(8)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			md.RenderToString(buf)
		}
	})

}

func BenchmarkGoldMark(b *testing.B) {
	md, err := ioutil.ReadFile(spec + ".md")
	if nil != err {
		b.Fatalf("read spec text failed: " + err.Error())
	}

	ge := goldmark.New(
		goldmark.WithRendererOptions(html.WithXHTML()),
		goldmark.WithExtensions(
			extension.Table, extension.Strikethrough, extension.TaskList, extension.Linkify,
		),
	)

	var out bytes.Buffer

	b.SetParallelism(8)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			out.Reset()
			if err := ge.Convert(md, &out); nil != err {
				panic(err)
			}
		}
	})

}

func BenchmarkBlackFriday(b *testing.B) {
	md, err := ioutil.ReadFile(spec + ".md")
	if nil != err {
		b.Fatalf("read spec text failed: " + err.Error())
	}
	b.SetParallelism(8)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			blackfriday.Run(md)
		}
	})
}

func BenchmarkGoMarkdown(b *testing.B) {
	md, err := ioutil.ReadFile(spec + ".md")
	if nil != err {
		b.Fatalf("read spec text failed: " + err.Error())
	}
	b.SetParallelism(8)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			gm.ToHTML(md, nil, nil)
		}
	})

}
