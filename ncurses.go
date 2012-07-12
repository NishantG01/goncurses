// goncurses - ncurses library for Go.
//
// Copyright (c) 2011, Rob Thornton 
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without 
// modification, are permitted provided that the following conditions are met:
//
//   * Redistributions of source code must retain the above copyright notice, 
//     this list of conditions and the following disclaimer.
//
//   * Redistributions in binary form must reproduce the above copyright 
//     notice, this list of conditions and the following disclaimer in the 
//     documentation and/or other materials provided with the distribution.
//  
//   * Neither the name of the copyright holder nor the names of its 
//     contributors may be used to endorse or promote products derived from 
//     this software without specific prior written permission.
//      
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" 
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE 
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE 
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE 
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR 
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF 
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS 
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN 
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) 
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE 
// POSSIBILITY OF SUCH DAMAGE.

/* ncurses library

   1. No functions which operate only on stdscr have been implemented because 
   it makes little sense to do so in a Go implementation. Stdscr is treated the
   same as any other window.

   2. Whenever possible, versions of ncurses functions which could potentially
   have a buffer overflow, like the getstr() family of functions, have not been
   implemented. Instead, only the mvwgetnstr() and wgetnstr() can be used. */
package goncurses

/* 
#cgo LDFLAGS: -lncurses
#include <ncurses.h>
#include <stdlib.h>

bool ncurses_has_mouse(void) {
	return has_mouse();
}
*/
import "C"

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

// Synconize options for Sync() function
const (
	SYNC_NONE   = iota
	SYNC_CURSOR // Sync cursor in all sub/derived windows
	SYNC_DOWN   // Sync changes in all parent windows
	SYNC_UP     // Sync change in all child windows
)

// Definitions for printed characters not found on most keyboards. Ideally, 
// these would not be hard-coded as they are potentially different on
// different systems. However, some ncurses implementations seem to be
// heavily reliant on macros which prevent these definitions from being
// handled by cgo properly. If they don't work for you, you won't be able
// to use them until either a) the Go team works out a way to overcome this
// limitation in godefs/cgo or b) an alternative method is found. Work is
// being done to find a solution from the ncurses source code.
const (
	ACS_DEGREE = iota + 4194406
	ACS_PLMINUS
	ACS_BOARD
	ACS_LANTERN
	ACS_LRCORNER
	ACS_URCORNER
	ACS_LLCORNER
	ACS_ULCORNER
	ACS_PLUS
	ACS_S1
	ACS_S3
	ACS_HLINE
	ACS_S7
	ACS_S9
	ACS_LTEE
	ACS_RTEE
	ACS_BTEE
	ACS_TTEE
	ACS_VLINE
	ACS_LEQUAL
	ACS_GEQUAL
	ACS_PI
	ACS_NEQUAL
	ACS_STERLING
	ACS_BULLET
	ACS_LARROW  = 4194347
	ACS_RARROW  = 4194348
	ACS_DARROW  = 4194349
	ACS_UARROW  = 4194350
	ACS_BLOCK   = 4194352
	ACS_DIAMOND = 4194400
	ACS_CKBOARD = 4194401
)

// Text attributes
const (
	A_NORMAL     = C.A_NORMAL
	A_STANDOUT   = C.A_STANDOUT
	A_UNDERLINE  = C.A_UNDERLINE
	A_REVERSE    = C.A_REVERSE
	A_BLINK      = C.A_BLINK
	A_DIM        = C.A_DIM
	A_BOLD       = C.A_BOLD
	A_PROTECT    = C.A_PROTECT
	A_INVIS      = C.A_INVIS
	A_ALTCHARSET = C.A_ALTCHARSET
	A_CHARTEXT   = C.A_CHARTEXT
)

var attrList = map[C.int]string{
	C.A_NORMAL:     "normal",
	C.A_STANDOUT:   "standout",
	C.A_UNDERLINE:  "underline",
	C.A_REVERSE:    "reverse",
	C.A_BLINK:      "blink",
	C.A_DIM:        "dim",
	C.A_BOLD:       "bold",
	C.A_PROTECT:    "protect",
	C.A_INVIS:      "invis",
	C.A_ALTCHARSET: "altcharset",
	C.A_CHARTEXT:   "chartext",
}

// Colors available to ncurses. Combine these with the dim/bold attributes
// for bright/dark versions of each color. These colors can be used for
// both background and foreground colors.
const (
	C_BLACK   = C.COLOR_BLACK
	C_BLUE    = C.COLOR_BLUE
	C_CYAN    = C.COLOR_CYAN
	C_GREEN   = C.COLOR_GREEN
	C_MAGENTA = C.COLOR_MAGENTA
	C_RED     = C.COLOR_RED
	C_WHITE   = C.COLOR_WHITE
	C_YELLOW  = C.COLOR_YELLOW
)

const (
	KEY_TAB       = 9               // tab
	KEY_RETURN    = 10              // enter key vs. KEY_ENTER
	KEY_DOWN      = C.KEY_DOWN      // down arrow key
	KEY_UP        = C.KEY_UP        // up arrow key
	KEY_LEFT      = C.KEY_LEFT      // left arrow key
	KEY_RIGHT     = C.KEY_RIGHT     // right arrow key
	KEY_HOME      = C.KEY_HOME      // home key
	KEY_BACKSPACE = C.KEY_BACKSPACE // backpace
	KEY_F1        = C.KEY_F0 + 1    // F1 key
	KEY_F2        = C.KEY_F0 + 2    // F2 key
	KEY_F3        = C.KEY_F0 + 3    // F3 key
	KEY_F4        = C.KEY_F0 + 4    // F4 key
	KEY_F5        = C.KEY_F0 + 5    // F5 key
	KEY_F6        = C.KEY_F0 + 6    // F6 key
	KEY_F7        = C.KEY_F0 + 7    // F7 key
	KEY_F8        = C.KEY_F0 + 8    // F8 key
	KEY_F9        = C.KEY_F0 + 9    // F9 key
	KEY_F10       = C.KEY_F0 + 10   // F10 key 
	KEY_F11       = C.KEY_F0 + 11   // F11 key
	KEY_F12       = C.KEY_F0 + 12   // F12 key
	KEY_DL        = C.KEY_DL        // delete-line key
	KEY_IL        = C.KEY_IL        // insert-line key
	KEY_DC        = C.KEY_DC        // delete-character key
	KEY_IC        = C.KEY_IC        // insert-character key
	KEY_EIC       = C.KEY_EIC       // sent by rmir or smir in insert mode
	KEY_CLEAR     = C.KEY_CLEAR     // clear-screen or erase key 
	KEY_EOS       = C.KEY_EOS       // clear-to-end-of-screen key 
	KEY_EOL       = C.KEY_EOL       // clear-to-end-of-line key 
	KEY_SF        = C.KEY_SF        // scroll-forward key 
	KEY_SR        = C.KEY_SR        // scroll-backward key 
	KEY_PAGEDOWN  = C.KEY_NPAGE     // page-down key (next-page)
	KEY_PAGEUP    = C.KEY_PPAGE     // page-up key (prev-page)
	KEY_STAB      = C.KEY_STAB      // set-tab key 
	KEY_CTAB      = C.KEY_CTAB      // clear-tab key 
	KEY_CATAB     = C.KEY_CATAB     // clear-all-tabs key 
	KEY_ENTER     = C.KEY_ENTER     // enter/send key
	KEY_PRINT     = C.KEY_PRINT     // print key 
	KEY_LL        = C.KEY_LL        // lower-left key (home down) 
	KEY_A1        = C.KEY_A1        // upper left of keypad 
	KEY_A3        = C.KEY_A3        // upper right of keypad 
	KEY_B2        = C.KEY_B2        // center of keypad 
	KEY_C1        = C.KEY_C1        // lower left of keypad 
	KEY_C3        = C.KEY_C3        // lower right of keypad 
	KEY_BTAB      = C.KEY_BTAB      // back-tab key 
	KEY_BEG       = C.KEY_BEG       // begin key 
	KEY_CANCEL    = C.KEY_CANCEL    // cancel key 
	KEY_CLOSE     = C.KEY_CLOSE     // close key 
	KEY_COMMAND   = C.KEY_COMMAND   // command key 
	KEY_COPY      = C.KEY_COPY      // copy key 
	KEY_CREATE    = C.KEY_CREATE    // create key 
	KEY_END       = C.KEY_END       // end key 
	KEY_EXIT      = C.KEY_EXIT      // exit key 
	KEY_FIND      = C.KEY_FIND      // find key 
	KEY_HELP      = C.KEY_HELP      // help key 
	KEY_MARK      = C.KEY_MARK      // mark key 
	KEY_MESSAGE   = C.KEY_MESSAGE   // message key 
	KEY_MOVE      = C.KEY_MOVE      // move key 
	KEY_NEXT      = C.KEY_NEXT      // next key 
	KEY_OPEN      = C.KEY_OPEN      // open key 
	KEY_OPTIONS   = C.KEY_OPTIONS   // options key 
	KEY_PREVIOUS  = C.KEY_PREVIOUS  // previous key 
	KEY_REDO      = C.KEY_REDO      // redo key 
	KEY_REFERENCE = C.KEY_REFERENCE // reference key 
	KEY_REFRESH   = C.KEY_REFRESH   // refresh key 
	KEY_REPLACE   = C.KEY_REPLACE   // replace key 
	KEY_RESTART   = C.KEY_RESTART   // restart key 
	KEY_RESUME    = C.KEY_RESUME    // resume key
	KEY_SAVE      = C.KEY_SAVE      // save key 
	KEY_SBEG      = C.KEY_SBEG      // shifted begin key 
	KEY_SCANCEL   = C.KEY_SCANCEL   // shifted cancel key 
	KEY_SCOMMAND  = C.KEY_SCOMMAND  // shifted command key 
	KEY_SCOPY     = C.KEY_SCOPY     // shifted copy key 
	KEY_SCREATE   = C.KEY_SCREATE   // shifted create key 
	KEY_SDC       = C.KEY_SDC       // shifted delete-character key 
	KEY_SDL       = C.KEY_SDL       // shifted delete-line key 
	KEY_SELECT    = C.KEY_SELECT    // select key 
	KEY_SEND      = C.KEY_SEND      // shifted end key 
	KEY_SEOL      = C.KEY_SEOL      // shifted clear-to-end-of-line key 
	KEY_SEXIT     = C.KEY_SEXIT     // shifted exit key 
	KEY_SFIND     = C.KEY_SFIND     // shifted find key 
	KEY_SHELP     = C.KEY_SHELP     // shifted help key 
	KEY_SHOME     = C.KEY_SHOME     // shifted home key 
	KEY_SIC       = C.KEY_SIC       // shifted insert-character key 
	KEY_SLEFT     = C.KEY_SLEFT     // shifted left-arrow key 
	KEY_SMESSAGE  = C.KEY_SMESSAGE  // shifted message key 
	KEY_SMOVE     = C.KEY_SMOVE     // shifted move key 
	KEY_SNEXT     = C.KEY_SNEXT     // shifted next key 
	KEY_SOPTIONS  = C.KEY_SOPTIONS  // shifted options key 
	KEY_SPREVIOUS = C.KEY_SPREVIOUS // shifted previous key 
	KEY_SPRINT    = C.KEY_SPRINT    // shifted print key 
	KEY_SREDO     = C.KEY_SREDO     // shifted redo key 
	KEY_SREPLACE  = C.KEY_SREPLACE  // shifted replace key 
	KEY_SRIGHT    = C.KEY_SRIGHT    // shifted right-arrow key 
	KEY_SRSUME    = C.KEY_SRSUME    // shifted resume key 
	KEY_SSAVE     = C.KEY_SSAVE     // shifted save key 
	KEY_SSUSPEND  = C.KEY_SSUSPEND  // shifted suspend key 
	KEY_SUNDO     = C.KEY_SUNDO     // shifted undo key 
	KEY_SUSPEND   = C.KEY_SUSPEND   // suspend key 
	KEY_UNDO      = C.KEY_UNDO      // undo key 
	KEY_MOUSE     = C.KEY_MOUSE     // any mouse event
	KEY_RESIZE    = C.KEY_RESIZE    // Terminal resize event 
	KEY_EVENT     = C.KEY_EVENT     // We were interrupted by an event 
	KEY_MAX       = C.KEY_MAX       // Maximum key value is KEY_EVENT (0633)
)

var keyList = map[C.int]string{
	9:               "tab",
	10:              "enter", // On some keyboards?
	C.KEY_DOWN:      "down",
	C.KEY_UP:        "up",
	C.KEY_LEFT:      "left",
	C.KEY_RIGHT:     "right",
	C.KEY_HOME:      "home",
	C.KEY_BACKSPACE: "backspace",
	C.KEY_ENTER:     "enter", // And not others?
	C.KEY_F0:        "F0",
	C.KEY_F0 + 1:    "F1",
	C.KEY_F0 + 2:    "F2",
	C.KEY_F0 + 3:    "F3",
	C.KEY_F0 + 4:    "F4",
	C.KEY_F0 + 5:    "F5",
	C.KEY_F0 + 6:    "F6",
	C.KEY_F0 + 7:    "F7",
	C.KEY_F0 + 8:    "F8",
	C.KEY_F0 + 9:    "F9",
	C.KEY_F0 + 10:   "F10",
	C.KEY_F0 + 11:   "F11",
	C.KEY_F0 + 12:   "F12",
	C.KEY_MOUSE:     "mouse",
	C.KEY_NPAGE:     "page down",
	C.KEY_PPAGE:     "page up",
}

// Mouse button events
const (
	M_ALL            = C.ALL_MOUSE_EVENTS
	M_ALT            = C.BUTTON_ALT      // alt-click
	M_B1_PRESSED     = C.BUTTON1_PRESSED // button 1
	M_B1_RELEASED    = C.BUTTON1_RELEASED
	M_B1_CLICKED     = C.BUTTON1_CLICKED
	M_B1_DBL_CLICKED = C.BUTTON1_DOUBLE_CLICKED
	M_B1_TPL_CLICKED = C.BUTTON1_TRIPLE_CLICKED
	M_B2_PRESSED     = C.BUTTON2_PRESSED // button 2
	M_B2_RELEASED    = C.BUTTON2_RELEASED
	M_B2_CLICKED     = C.BUTTON2_CLICKED
	M_B2_DBL_CLICKED = C.BUTTON2_DOUBLE_CLICKED
	M_B2_TPL_CLICKED = C.BUTTON2_TRIPLE_CLICKED
	M_B3_PRESSED     = C.BUTTON3_PRESSED // button 3
	M_B3_RELEASED    = C.BUTTON3_RELEASED
	M_B3_CLICKED     = C.BUTTON3_CLICKED
	M_B3_DBL_CLICKED = C.BUTTON3_DOUBLE_CLICKED
	M_B3_TPL_CLICKED = C.BUTTON3_TRIPLE_CLICKED
	M_B4_PRESSED     = C.BUTTON4_PRESSED // button 4
	M_B4_RELEASED    = C.BUTTON4_RELEASED
	M_B4_CLICKED     = C.BUTTON4_CLICKED
	M_B4_DBL_CLICKED = C.BUTTON4_DOUBLE_CLICKED
	M_B4_TPL_CLICKED = C.BUTTON4_TRIPLE_CLICKED
	M_CTRL           = C.BUTTON_CTRL           // ctrl-click
	M_SHIFT          = C.BUTTON_SHIFT          // shift-click
	M_POSITION       = C.REPORT_MOUSE_POSITION // mouse moved
)

type Window struct {
	win *C.WINDOW
}

type Pad Window

// BaudRate returns the speed of the terminal in bits per second
func BaudRate() int {
	return int(C.baudrate())
}

// Beep requests the terminal make an audible bell or, if not available,
// flashes the screen. Note that screen flashing doesn't work on all
// terminals
func Beep() {
	C.beep()
}

// Turn on/off buffering; raw user signals are passed to the program for
// handling. Overrides raw mode
func CBreak(on bool) {
	if on {
		C.cbreak()
		return
	}
	C.nocbreak()
}

// Test whether colour values can be changed
func CanChangeColor() bool {
	return bool(C.bool(C.can_change_color()))
}

// Get RGB values for specified colour
func ColorContent(col int) (int, int, int) {
	var r, g, b C.short
	C.color_content(C.short(col), (*C.short)(&r), (*C.short)(&g),
		(*C.short)(&b))
	return int(r), int(g), int(b)
}

// Return the value of a color pair which can be passed to functions which
// accept attributes like AddChar or AttrOn/Off.
func ColorPair(pair int) int {
	return int(C.COLOR_PAIR(C.int(pair)))
}

// CursesVersion returns the version of the ncurses library currently linked to
func CursesVersion() string {
	return C.GoString(C.curses_version())
}

// Set the cursor visibility. Options are: 0 (invisible/hidden), 1 (normal)
// and 2 (extra-visible)
func Cursor(vis byte) error {
	if C.curs_set(C.int(vis)) == C.ERR {
		return errors.New("Failed to enable ")
	}
	return nil
}

// Echo turns on/off the printing of typed characters
func Echo(on bool) {
	if on {
		C.echo()
		return
	}
	C.noecho()
}

// Must be called prior to exiting the program in order to make sure the
// terminal returns to normal operation
func End() {
	C.endwin()
}

// Flash requests the terminal flashes the screen or, if not available,
// make an audible bell. Note that screen flashing doesn't work on all
// terminals
func Flash() {
	C.flash()
}

// Returns an array of integers representing the following, in order:
// x, y and z coordinates, id of the device, and a bit masked state of
// the devices buttons
func GetMouse() ([]int, error) {
	if bool(C.ncurses_has_mouse()) != true {
		return nil, errors.New("Mouse support not enabled")
	}
	var event C.MEVENT
	if C.getmouse(&event) != C.OK {
		return nil, errors.New("Failed to get mouse event")
	}
	return []int{int(event.x), int(event.y), int(event.z), int(event.id),
		int(event.bstate)}, nil
}

// Behaves like cbreak() but also adds a timeout for input. If timeout is
// exceeded after a call to Getch() has been made then GetChar will return
// with an error.
func HalfDelay(delay int) error {
	var cerr C.int
	if delay > 0 {
		cerr = C.halfdelay(C.int(delay))
	}
	if cerr == C.ERR {
		return errors.New("Unable to set delay mode")
	}
	return nil
}

// HasColors returns true if terminal can display colors
func HasColors() bool {
	return bool(C.has_colors())
}

// HasKey returns true if terminal recognized the given character
func HasKey(ch int) bool {
	if C.has_key(C.int(ch)) == 1 {
		return true
	}
	return false
}

// InitColor is used to set 'color' to the specified RGB values. Values may
// be between 0 and 1000.
func InitColor(col int, r, g, b int) error {
	if C.init_color(C.short(col), C.short(r), C.short(g),
		C.short(b)) == C.ERR {
		return errors.New("Failed to set new color definition")
	}
	return nil
}

// InitPair sets a colour pair designated by 'pair' to fg and bg colors
func InitPair(pair byte, fg, bg int) error {
	if pair == 0 || C.int(pair) > (C.COLOR_PAIRS-1) {
		return errors.New("Invalid color pair selected")
	}
	if C.init_pair(C.short(pair), C.short(fg), C.short(bg)) == C.ERR {
		return errors.New("Failed to init color pair")
	}
	return nil
}

// Initialize the ncurses library. You must run this function prior to any 
// other goncurses function in order for the library to work
func Init() (stdscr Window, err error) {
	stdscr = Window{C.initscr()}
	if unsafe.Pointer(stdscr.win) == nil {
		err = errors.New("An error occurred initializing ncurses")
	}
	return
}

// IsEnd returns true if End() has been called, otherwise false
func IsEnd() bool {
	return bool(C.isendwin())
}

// IsTermResized returns true if ResizeTerm would modify any current Windows 
// if called with the given parameters
func IsTermResized(nlines, ncols int) bool {
	return bool(C.is_term_resized(C.int(nlines), C.int(ncols)))
}

// Returns a string representing the value of input returned by Getch
func Key(k int) string {
	key, ok := keyList[C.int(k)]
	if !ok {
		key = fmt.Sprintf("%c", k)
	}
	return key
}

func Mouse() bool {
	return bool(C.ncurses_has_mouse())
}

func MouseInterval() {
}

// MouseMask accepts a single int of OR'd mouse events. If a mouse event
// is triggered, GetChar() will return KEY_MOUSE. To retrieve the actual
// event use GetMouse() to pop it off the queue. Pass a pointer as the 
// second argument to store the prior events being monitored or nil.
func MouseMask(mask int, old *int) (m int) {
	if bool(C.ncurses_has_mouse()) {
		m = int(C.mousemask((C.mmask_t)(mask),
			(*C.mmask_t)(unsafe.Pointer(old))))
	}
	return
}

// NewWindow creates a window of size h(eight) and w(idth) at y, x
func NewWindow(h, w, y, x int) (window Window, err error) {
	window = Window{C.newwin(C.int(h), C.int(w), C.int(y), C.int(x))}
	if window.win == nil {
		err = errors.New("Failed to create a new window")
	}
	return
}

// NL turns newline translation on/off.
func NL(on bool) {
	if on {
		C.nl()
		return
	}
	C.nonl()
}

// Raw turns on input buffering; user signals are disabled and the key strokes 
// are passed directly to input. Set to false if you wish to turn this mode
// off
func Raw(on bool) {
	if on {
		C.raw()
		return
	}
	C.noraw()
}

// ResizeTerm will attempt to resize the terminal. This only has an effect if
// the terminal is in an XWindows (GUI) environment.
func ResizeTerm(nlines, ncols int) error {
	if C.resizeterm(C.int(nlines), C.int(ncols)) == C.ERR {
		return errors.New("Failed to resize terminal")
	}
	return nil
}

// Enables colors to be displayed. Will return an error if terminal is not
// capable of displaying colors
func StartColor() error {
	if C.has_colors() == C.bool(false) {
		return errors.New("Terminal does not support colors")
	}
	if C.start_color() == C.ERR {
		return errors.New("Failed to enable color mode")
	}
	return nil
}

// Update the screen, refreshing all windows
func Update() error {
	if C.doupdate() == C.ERR {
		return errors.New("Failed to update")
	}
	return nil
}

// NewPad creates a window which is not restricted by the terminal's 
// dimentions (unlike a Window)
func NewPad(lines, cols int) Pad {
	return Pad{C.newpad(C.int(lines), C.int(cols))}
}

// Echo prints a single character to the pad immediately. This has the
// same effect of calling AddChar() + Refresh() but has a significant
// speed advantage
func (p *Pad) Echo(ch int) {
	C.pechochar(p.win, C.chtype(ch))
}

func (p *Pad) NoutRefresh(py, px, ty, tx, by, bx int) {
	C.pnoutrefresh(p.win, C.int(py), C.int(px), C.int(ty),
		C.int(tx), C.int(by), C.int(bx))
}

// Refresh the pad at location py, px using the rectangle specified by
// ty, tx, by, bx (bottom/top y/x)
func (p *Pad) Refresh(py, px, ty, tx, by, bx int) {
	C.prefresh(p.win, C.int(py), C.int(px), C.int(ty), C.int(tx),
		C.int(by), C.int(bx))
}

// Sub creates a sub-pad lines by columns in size
func (p *Pad) Sub(y, x, h, w int) Pad {
	return Pad{C.subpad(p.win, C.int(h), C.int(w), C.int(y),
		C.int(x))}
}

// Window is a helper function for calling Window functions on a pad like
// Print(). Convention would be to use Pad.Window().Print().
func (p *Pad) Window() *Window {
	return (*Window)(p)
}

// AddChar prints a single character to the window. The character can be
// OR'd together with attributes and colors. If optional first or second
// arguments are given they are the y and x coordinates on the screen
// respectively. If only y is given, x is assumed to be zero.
func (w *Window) AddChar(args ...int) {
	var cattr C.int
	var count, y, x int

	if len(args) > 1 {
		y = args[0]
		count += 1
	}
	if len(args) > 2 {
		x = args[1]
		count += 1
	}
	cattr |= C.int(args[count])
	if count > 0 {
		C.mvwaddch(w.win, C.int(y), C.int(x), C.chtype(cattr))
		return
	}
	C.waddch(w.win, C.chtype(cattr))
}

// Turn off character attribute.
func (w *Window) AttrOff(attr int) (err error) {
	if C.wattroff(w.win, C.int(attr)) == C.ERR {
		err = errors.New(fmt.Sprintf("Failed to unset attribute: %s",
			attrList[C.int(attr)]))
	}
	return
}

// Turn on character attribute
func (w *Window) AttrOn(attr int) (err error) {
	if C.wattron(w.win, C.int(attr)) == C.ERR {
		err = errors.New(fmt.Sprintf("Failed to set attribute: %s",
			attrList[C.int(attr)]))
	}
	return
}

func (w *Window) Background(attr int) {
	C.wbkgd(w.win, C.chtype(attr))
}

// Border uses the characters supplied to draw a border around the window.
// t, b, r, l, s correspond to top, bottom, right, left and side respectively.
func (w *Window) Border(ls, rs, ts, bs, tl, tr, bl, br int) error {
	res := C.wborder(w.win, C.chtype(ls), C.chtype(rs), C.chtype(ts),
		C.chtype(bs), C.chtype(tl), C.chtype(tr), C.chtype(bl),
		C.chtype(br))
	if res == C.ERR {
		return errors.New("Failed to draw box around window")
	}
	return nil
}

// Box draws a border around the given window. For complete control over the
// characters used to draw the border use Border()
func (w *Window) Box(vch, hch int) error {
	if C.box(w.win, C.chtype(vch), C.chtype(hch)) == C.ERR {
		return errors.New("Failed to draw box around window")
	}
	return nil
}

// Clear the screen
func (w *Window) Clear() error {
	if C.wclear(w.win) == C.ERR {
		return errors.New("Failed to clear screen")
	}
	return nil
}

// ClearOk clears the window completely prior to redrawing it. If called
// on stdscr then the whole screen is redrawn no matter which window has
// Refresh() called on it. Defaults to False.
func (w *Window) ClearOk(ok bool) {
	C.clearok(w.win, C.bool(ok))
}

// Clear starting at the current cursor position, moving to the right, to the 
// bottom of window
func (w *Window) ClearToBottom() error {
	if C.wclrtobot(w.win) == C.ERR {
		return errors.New("Failed to clear bottom of window")
	}
	return nil
}

// Clear from the current cursor position, moving to the right, to the end 
// of the line
func (w *Window) ClearToEOL() error {
	if C.wclrtoeol(w.win) == C.ERR {
		return errors.New("Failed to clear to end of line")
	}
	return nil
}

// Color sets the forground/background color pair for the entire window
func (w *Window) Color(pair byte) {
	C.wcolor_set(w.win, C.short(C.COLOR_PAIR(C.int(pair))), nil)
}

// ColorOff turns the specified color pair off
func (w *Window) ColorOff(pair byte) error {
	if C.wattroff(w.win, C.COLOR_PAIR(C.int(pair))) == C.ERR {
		return errors.New("Failed to enable color pair")
	}
	return nil
}

// Normally color pairs are turned on via attron() in ncurses but this
// implementation chose to make it seperate
func (w *Window) ColorOn(pair byte) error {
	if C.wattron(w.win, C.COLOR_PAIR(C.int(pair))) == C.ERR {
		return errors.New("Failed to enable color pair")
	}
	return nil
}

// Copy is similar to Overlay and Overwrite but provides a finer grain of
// control. 
func (w *Window) Copy(src *Window, sy, sx, dtr, dtc, dbr, dbc int,
	overlay bool) error {
	var ol int
	if overlay {
		ol = 1
	}
	if C.copywin(src.win, w.win, C.int(sy), C.int(sx),
		C.int(dtr), C.int(dtc), C.int(dbr), C.int(dbc), C.int(ol)) ==
		C.ERR {
		return errors.New("Failed to copy window")
	}
	return nil
}

// DelChar
func (w *Window) DelChar(coord ...int) error {
	if len(coord) > 2 {
		return errors.New(fmt.Sprintf("Invalid number of arguments, "+
			"expected 2, got %d", len(coord)))
	}
	var err C.int
	if len(coord) > 1 {
		var x int
		y := coord[0]
		if len(coord) > 2 {
			x = coord[1]
		}
		err = C.mvwdelch(w.win, C.int(y), C.int(x))
	} else {
		err = C.wdelch(w.win)
	}
	if err != C.OK {
		return errors.New("An error occurred when trying to delete " +
			"character")
	}
	return nil
}

// Delete the window
func (w *Window) Delete() error {
	if C.delwin(w.win) == C.ERR {
		return errors.New("Failed to delete window")
	}
	w = nil
	return nil
}

// Derived creates a new window of height and width at the coordinates
// y, x.  These coordinates are relative to the original window thereby 
// confining the derived window to the area of original window. See the
// SubWindow function for additional notes.
func (w *Window) Derived(height, width, y, x int) Window {
	return Window{C.derwin(w.win, C.int(height), C.int(width), C.int(y),
		C.int(x))}
}

// Duplicate the window, creating an exact copy.
func (w *Window) Duplicate() Window {
	return Window{C.dupwin(w.win)}
}

// Test whether the given mouse coordinates are within the window or not
func (w *Window) Enclose(y, x int) bool {
	return bool(C.wenclose(w.win, C.int(y), C.int(x)))
}

// Erase the contents of the window, effectively clearing it
func (w *Window) Erase() {
	C.werase(w.win)
}

// Get a character from standard input
func (w *Window) GetChar(coords ...int) int {
	var y, x, count int
	if len(coords) > 1 {
		y = coords[0]
		count++
	}
	if len(coords) > 2 {
		x = coords[1]
		count++
	}
	if count > 0 {
		return int(C.mvwgetch(w.win, C.int(y), C.int(x)))
	}
	return int(C.wgetch(w.win))
}

// Reads at most 'n' characters entered by the user from the Window. Attempts
// to enter greater than 'n' characters will elicit a 'beep'
func (w *Window) GetString(n int) (string, error) {
	// TODO: add move portion of code...
	cstr := make([]C.char, n)
	if C.wgetnstr(w.win, (*C.char)(&cstr[0]), C.int(n)) == C.ERR {
		return "", errors.New("Failed to retrieve string from input stream")
	}
	return C.GoString(&cstr[0]), nil
}

// Getyx returns the current cursor location in the Window. Note that it uses 
// ncurses idiom of returning y then x.
func (w *Window) Getyx() (int, int) {
	// In some cases, getxy() and family are macros which don't play well with
	// cgo
	return int(w.win._cury), int(w.win._curx)
}

// HLine draws a horizontal line starting at y, x and ending at width using 
// the specified character
func (w *Window) HLine(y, x, ch, wid int) {
	// TODO: move portion	
	C.mvwhline(w.win, C.int(y), C.int(x), C.chtype(ch), C.int(wid))
	return
}

// IsCleared returns the value set in ClearOk
func (w *Window) IsCleared() bool {
	return bool(w.win._clear)
}

// IsKeypad returns the value set in Keypad
func (w *Window) IsKeypad() bool {
	return bool(w.win._use_keypad)
}

// Keypad turns on/off the keypad characters, including those like the F1-F12 
// keys and the arrow keys
func (w *Window) Keypad(keypad bool) error {
	var err C.int
	if err = C.keypad(w.win, C.bool(keypad)); err == C.ERR {
		return errors.New("Unable to set keypad mode")
	}
	return nil
}

// Returns the maximum size of the Window. Note that it uses ncurses idiom
// of returning y then x.
func (w *Window) Maxyx() (int, int) {
	// This hack is necessary to make cgo happy
	return int(w.win._maxy + 1), int(w.win._maxx + 1)
}

// Move the cursor to the specified coordinates within the window
func (w *Window) Move(y, x int) {
	C.wmove(w.win, C.int(y), C.int(x))
	return
}

// NoutRefresh flags the window for redrawing. In order to actually perform
// the changes, Update() must be called. This function when coupled with
// Update() provides a speed increase over using Refresh() on each window.
func (w *Window) NoutRefresh() {
	C.wnoutrefresh(w.win)
	return
}

// Overlay copies overlapping sections of src window onto the destination
// window. Non-blank elements are not overwritten.
func (w *Window) Overlay(src *Window) error {
	if C.overlay(src.win, w.win) == C.ERR {
		return errors.New("Failed to overlay window")
	}
	return nil
}

// Overwrite copies overlapping sections of src window onto the destination
// window. This function is considered "destructive" by copying all
// elements of src onto the destination window.
func (w *Window) Overwrite(src *Window) error {
	if C.overwrite(src.win, w.win) == C.ERR {
		return errors.New("Failed to overwrite window")
	}
	return nil
}

func (w *Window) Parent() *Window {
	return &Window{w.win._parent}
}

// Print a string to the given window. The first two arguments may be
// coordinates to print to. If only one integer is supplied, it is assumed to
// be the y coordinate, x therefore defaults to 0. In order to simulate the 'n' 
// versions of functions like addnstr use a string slice.
// Examples:
// goncurses.Print("hello!")
// goncurses.Print("hello %s!", "world")
// goncurses.Print(23, "hello!") // moves to 23, 0 and prints "hello!"
// goncurses.Print(5, 10, "hello %s!", "world") // move to 5, 10 and print
//                                              // "hello world!"
func (w *Window) Print(args ...interface{}) {
	var count, y, x int

	if len(args) > 1 {
		if reflect.TypeOf(args[0]).String() == "int" {
			y = args[0].(int)
			count += 1
		}
	}
	if len(args) > 2 {
		if reflect.TypeOf(args[1]).String() == "int" {
			x = args[1].(int)
			count += 1
		}
	}

	cstr := C.CString(fmt.Sprintf(args[count].(string), args[count+1:]...))
	defer C.free(unsafe.Pointer(cstr))

	if count > 0 {
		C.mvwaddstr(w.win, C.int(y), C.int(x), cstr)
		return
	}
	C.waddstr(w.win, cstr)
}

// Refresh the window so it's contents will be displayed
func (w *Window) Refresh() {
	C.wrefresh(w.win)
}

// Resize the window to new height, width
func (w *Window) Resize(height, width int) {
	C.wresize(w.win, C.int(height), C.int(width))
}

// Scroll the contents of the window. Use a negative number to scroll up,
// a positive number to scroll down. ScrollOk Must have been called prior.
func (w *Window) Scroll(n int) {
	C.wscrl(w.win, C.int(n))
}

// ScrollOk sets whether scrolling will work
func (w *Window) ScrollOk(ok bool) {
	C.scrollok(w.win, C.bool(ok))
}

// SubWindow creates a new window of height and width at the coordinates
// y, x.  This window shares memory with the original window so changes
// made to one window are reflected in the other. It is necessary to call
// Touch() on this window prior to calling Refresh in order for it to be
// displayed.
func (w *Window) Sub(height, width, y, x int) Window {
	return Window{C.subwin(w.win, C.int(height), C.int(width), C.int(y),
		C.int(x))}
}

// Sync updates all parent or child windows which were created via
// SubWindow() or DerivedWindow(). Argument can be one of: SYNC_DOWN, which
// syncronizes all parent windows (done by Refresh() by default so should
// rarely, if ever, need to be called); SYNC_UP, which updates all child
// windows to match any updates made to the parent; and, SYNC_CURSOR, which
// updates the cursor position only for all windows to match the parent window
func (w *Window) Sync(sync int) {
	switch sync {
	case SYNC_DOWN:
		C.wsyncdown(w.win)
	case SYNC_CURSOR:
		C.wcursyncup(w.win)
	case SYNC_UP:
		C.wsyncup(w.win)
	}
}

// Touch indicates that the window contains changes which should be updated
// on the next call to Refresh
func (w *Window) Touch() {
	// may not use touchwin() directly. cgo does not handle macros well.
	y, _ := w.Maxyx()
	C.wtouchln(w.win, 0, C.int(y), 1)
}

// VLine draws a verticle line starting at y, x and ending at height using 
// the specified character
func (w *Window) VLine(y, x, ch, h int) {
	// TODO: move portion
	C.mvwvline(w.win, C.int(y), C.int(x), C.chtype(ch), C.int(h))
	return
}
