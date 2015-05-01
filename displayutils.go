package main

import (
	//	"fmt"
	"container/list"
	"github.com/nsf/termbox-go"
	"github.com/simulatedsimian/go_sandbox/geom"
	"reflect"
	"unicode"
)

func printAt(x, y int, s string, fg, bg termbox.Attribute) {
	for _, r := range s {
		termbox.SetCell(x, y, r, fg, bg)
		x++
	}
}

func printAtDef(x, y int, s string) {
	printAt(x, y, s, termbox.ColorDefault, termbox.ColorDefault)
}

func clearRect(rect geom.Rectangle, c rune, fg, bg termbox.Attribute) {
	w, h := termbox.Size()
	sz := geom.RectangleFromSize(geom.Coord{w, h})

	toClear, ok := geom.RectangleIntersection(rect, sz)
	if ok {
		for y := toClear.Min.Y; y < toClear.Max.Y; y++ {
			for x := toClear.Min.X; x < toClear.Max.X; x++ {
				termbox.SetCell(x, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
			}
		}
	}
}

func clearRectDef(rect geom.Rectangle) {
	clearRect(rect, '.', termbox.ColorDefault, termbox.ColorDefault)
}

type DisplayElement interface {
	GiveFocus() bool
	HandleInput(k termbox.Key, r rune)
	Draw()
}

type DisplayList struct {
	list       []DisplayElement
	focusIndex int
}

func (dl *DisplayList) AddElement(elem DisplayElement) {
	dl.list = append(dl.list, elem)
}

func (dl *DisplayList) Draw() {
	w, h := termbox.Size()
	clearRectDef(geom.RectangleFromSize(geom.Coord{w, h}))

	for _, elem := range dl.list {
		elem.Draw()
	}
}

func (dl *DisplayList) NextFocus() {
	if dl.list != nil && len(dl.list) > 0 {
		for {
			dl.focusIndex++
			if dl.focusIndex >= len(dl.list) {
				dl.focusIndex = 0
			}

			if dl.list[dl.focusIndex].GiveFocus() {
				break
			}
		}
	}
}

func (dl *DisplayList) PrevFocus() {
	if dl.list != nil && len(dl.list) > 0 {
		for {
			dl.focusIndex--
			if dl.focusIndex < 0 {
				dl.focusIndex = len(dl.list)
			}

			if dl.list[dl.focusIndex].GiveFocus() {
				break
			}
		}
	}
}

func (dl *DisplayList) HandleInput(k termbox.Key, r rune) {

	if dl.list != nil && len(dl.list) > 0 {
		if k == termbox.KeyTab {
			dl.NextFocus()
		} else {
			dl.list[dl.focusIndex].HandleInput(k, r)
		}
	}
}

type InputHandler func(inp string)

type TextInputField struct {
	x, y       int
	inp        []rune
	cursorLoc  int
	inpHandler InputHandler
	hasFocus   bool
	history    *list.List
	histPos    *list.Element
}

func MakeTextInputField(x, y int, inpHandler InputHandler) *TextInputField {
	return &TextInputField{x, y, nil, 0, inpHandler, false, list.New(), nil}
}

func (t *TextInputField) HandleInput(k termbox.Key, r rune) {
	if k == termbox.KeyEnter {
		t.inp = TrimRunes(t.inp, unicode.IsSpace)
		if len(t.inp) > 0 {
			if t.histPos != nil && reflect.DeepEqual(t.histPos.Value, t.inp) {
				t.history.MoveToBack(t.histPos)
			} else {
				t.history.PushBack(t.inp)
			}
			t.histPos = nil

			t.inpHandler(string(t.inp))
			t.inp = nil
			t.cursorLoc = 0
		}
	}

	if k == termbox.KeyArrowUp {
		if t.history.Len() > 0 {
			if t.histPos == nil {
				t.histPos = t.history.Back()
			} else {
				t.histPos = t.histPos.Prev()
				if t.histPos == nil {
					t.histPos = t.history.Front()
				}
			}
			t.inp = CloneRuneSlice(t.histPos.Value.([]rune))
			t.cursorLoc = len(t.inp)
		}
	}

	if k == termbox.KeyArrowDown {
		if t.history.Len() > 0 {
			if t.histPos != nil {
				t.histPos = t.histPos.Next()
			}

			t.inp = nil
			if t.histPos != nil {
				t.inp = append(t.inp, t.histPos.Value.([]rune)...)
			}
			t.cursorLoc = len(t.inp)
		}
	}

	if r > ' ' {
		t.inp = append(t.inp, r)
		t.cursorLoc++
	}

	if k == 32 {
		t.inp = append(t.inp, ' ')
		t.cursorLoc++
	}

	if len(t.inp) > 0 && (k == termbox.KeyBackspace || k == termbox.KeyBackspace2) {
		t.inp = t.inp[:len(t.inp)-1]
		t.cursorLoc--
	}

	termbox.SetCursor(t.x+t.cursorLoc, t.y)
	//	printAtDef(t.x, t.y+1, fmt.Sprintf("%v, %v               ", k, r))
}

func (t *TextInputField) Draw() {
	printAtDef(t.x, t.y, string(t.inp)+" ")
}

func (t *TextInputField) GiveFocus() bool {
	termbox.SetCursor(t.x, t.y)
	return true
}

type StaticText struct {
	x, y int
	text string
}

func (t *StaticText) HandleInput(k termbox.Key, r rune) {
}

func (t *StaticText) Draw() {
	printAtDef(t.x, t.y, t.text)
}

func (t *StaticText) GiveFocus() bool {
	return false
}

type ScrollingTextOutput struct {
	x, y int
	w, h int
	text []string
}

func (t *ScrollingTextOutput) HandleInput(k termbox.Key, r rune) {
}

func (t *ScrollingTextOutput) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (t *ScrollingTextOutput) WriteLine(l string) {
	t.text = append(t.text, l)
}

func (t *ScrollingTextOutput) Draw() {
	clearRectDef(geom.RectangleFromPosSize(geom.Coord{t.x, t.y}, geom.Coord{t.w, t.h}))

	start := 0

	if len(t.text) > t.h {
		start = len(t.text) - t.h
	}

	y := t.y
	for l := start; l < len(t.text); l++ {
		printAtDef(t.x, y, t.text[l])
		y++
	}
}

func (t *ScrollingTextOutput) GiveFocus() bool {
	return false
}
