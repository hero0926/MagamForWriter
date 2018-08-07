// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"time"
)

import (
	"strconv"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	var mw *walk.MainWindow

	if err := (MainWindow{
		AssignTo: &mw,
		Title:    "Walk LogView Example",
		MinSize:  Size{320, 240},
		Size:     Size{400, 600},
		Layout:   VBox{MarginsZero: true},
	}.Create()); err != nil {
		log.Fatal(err)
	}

	lv, err := NewLogView(mw)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(lv)

	go func() {
		for i := 0; i < 10000; i++ {
			time.Sleep(1 * time.Second)

			str := TxtFileOpen("D:\\새 텍스트 문서.txt")
			log.Println(strconv.Itoa(CountAll(str)) + "\r\n")
		}
	}()

	mw.Run()
}