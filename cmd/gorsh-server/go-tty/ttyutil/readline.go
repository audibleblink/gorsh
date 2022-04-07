package ttyutil

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-runewidth"
	"go-tty"
)

type ctx struct {
	w        io.Writer
	input    []rune
	last     []rune
	prompt   string
	cursor_x int
	old_row  int
	old_crow int
	size     int
}

func (c *ctx) redraw(dirty bool, passwordChar rune) error {
	var buf bytes.Buffer

	buf.WriteString("\x1b[5>h")

	buf.WriteString("\x1b[1G")
	if dirty {
		buf.WriteString("\x1b[0K")
	}
	for i := 0; i < c.old_row-c.old_crow; i++ {
		buf.WriteString("\x1b[B")
	}
	for i := 0; i < c.old_row; i++ {
		if dirty {
			buf.WriteString("\x1b[2K")
		}
		buf.WriteString("\x1b[A")
	}

	var rs []rune
	if passwordChar != 0 {
		for i := 0; i < len(c.input); i++ {
			rs = append(rs, passwordChar)
		}
	} else {
		rs = c.input
	}

	ccol, crow, col, row := -1, 0, 0, 0
	plen := len([]rune(c.prompt))
	for i, r := range []rune(c.prompt + string(rs)) {
		if i == plen+c.cursor_x {
			ccol = col
			crow = row
		}
		rw := runewidth.RuneWidth(r)
		if col+rw > c.size {
			col = 0
			row++
			if dirty {
				buf.WriteString("\n\r\x1b[0K")
			}
		}
		if dirty {
			buf.WriteString(string(r))
		}
		col += rw
	}
	if dirty {
		buf.WriteString("\x1b[1G")
		for i := 0; i < row; i++ {
			buf.WriteString("\x1b[A")
		}
	}
	if ccol == -1 {
		ccol = col
		crow = row
	}
	for i := 0; i < crow; i++ {
		buf.WriteString("\x1b[B")
	}
	buf.WriteString(fmt.Sprintf("\x1b[%dG", ccol+1))

	buf.WriteString("\x1b[5>l")
	io.Copy(c.w, &buf)

	c.old_row = row
	c.old_crow = crow

	return nil
}

func ReadLine(tty *tty.TTY) (string, error) {
	c := new(ctx)
	c.w = colorable.NewColorableStdout()
	quit := false
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	go func() {
		<-sc
		c.input = nil
		quit = true
	}()
	c.size = 80

	dirty := true
loop:
	for !quit {
		err := c.redraw(dirty, 0)
		if err != nil {
			return "", err
		}
		dirty = false

		r, err := tty.ReadRune()
		if err != nil {
			break
		}
		switch r {
		case 0:
		case 1: // CTRL-A
			c.cursor_x = 0
		case 2: // CTRL-B
			if c.cursor_x > 0 {
				c.cursor_x--
			}
		case 3: // BREAK
			return "", nil
		case 4: // CTRL-D
			if len(c.input) > 0 {
				continue
			}
			return "", io.EOF
		case 5: // CTRL-E
			c.cursor_x = len(c.input)
		case 6: // CTRL-F
			if c.cursor_x < len(c.input) {
				c.cursor_x++
			}
		case 8, 0x7F: // BS
			if c.cursor_x > 0 {
				c.input = append(c.input[0:c.cursor_x-1], c.input[c.cursor_x:len(c.input)]...)
				c.cursor_x--
				dirty = true
			}
		case 27:
			if !tty.Buffered() {
				return "", io.EOF
			}
			r, err = tty.ReadRune()
			if err == nil && r == 0x5b {
				r, err = tty.ReadRune()
				if err != nil {
					panic(err)
				}
				switch r {
				case 'C':
					if c.cursor_x < len(c.input) {
						c.cursor_x++
					}
				case 'D':
					if c.cursor_x > 0 {
						c.cursor_x--
					}
				}
			}
		case 10: // LF
			break loop
		case 11: // CTRL-K
			c.input = c.input[:c.cursor_x]
			dirty = true
		case 12: // CTRL-L
			dirty = true
		case 13: // CR
			break loop
		case 21: // CTRL-U
			c.input = c.input[c.cursor_x:]
			c.cursor_x = 0
			dirty = true
		case 23: // CTRL-W
			for i := len(c.input) - 1; i >= 0; i-- {
				if i == 0 || c.input[i] == ' ' || c.input[i] == '\t' {
					c.input = append(c.input[:i], c.input[c.cursor_x:]...)
					c.cursor_x = i
					dirty = true
					break
				}
			}
		default:
			tmp := []rune{}
			tmp = append(tmp, c.input[0:c.cursor_x]...)
			tmp = append(tmp, r)
			c.input = append(tmp, c.input[c.cursor_x:len(c.input)]...)
			c.cursor_x++
			dirty = true
		}
	}
	os.Stdout.WriteString("\n")

	if c.input == nil {
		return "", io.EOF
	}

	return string(c.input), nil
}
