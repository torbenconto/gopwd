package qr

import (
	"io"
	"rsc.io/qr"
	"strings"
)

const (
	white        = "\033[47m  \033[0m"
	black        = "\033[40m  \033[0m"
	blackWhite   = "▄"
	blackBlack   = " "
	whiteBlack   = "▀"
	whiteWhite   = "█"
	H            = qr.H
	M            = qr.M
	L            = qr.L
	defaultQuiet = 0
)

type Config struct {
	Level          qr.Level
	Writer         io.Writer
	HalfBlocks     bool
	BlackChar      string
	BlackWhiteChar string
	WhiteChar      string
	WhiteBlackChar string
	QuietZone      int
}

func (c *Config) writeFullBlocks(w io.Writer, code *qr.Code) {
	ww := c.WhiteChar
	bb := c.BlackChar
	for i := 0; i < code.Size; i++ {
		for j := 0; j < code.Size; j++ {
			if code.Black(j, i) {
				w.Write([]byte(bb))
			} else {
				w.Write([]byte(ww))
			}
		}
		w.Write([]byte("\n"))
	}
}

func (c *Config) writeHalfBlocks(w io.Writer, code *qr.Code) {
	ww := c.WhiteChar
	bb := c.BlackChar
	wb := c.WhiteBlackChar
	bw := c.BlackWhiteChar
	for i := 0; i < code.Size; i += 2 {
		for j := 0; j < code.Size; j++ {
			nextBlack := false
			if i+1 < code.Size {
				nextBlack = code.Black(j, i+1)
			}
			currBlack := code.Black(j, i)
			if currBlack && nextBlack {
				w.Write([]byte(bb))
			} else if currBlack && !nextBlack {
				w.Write([]byte(bw))
			} else if !currBlack && !nextBlack {
				w.Write([]byte(ww))
			} else {
				w.Write([]byte(wb))
			}
		}
		w.Write([]byte("\n"))
	}
}

func stringRepeat(s string, count int) string {
	return strings.Repeat(s, count)
}

func GenerateWithConfig(text string, config Config) {
	if config.QuietZone < 1 {
		config.QuietZone = 1
	}
	w := config.Writer
	code, _ := qr.Encode(text, config.Level)
	if config.HalfBlocks {
		config.writeHalfBlocks(w, code)
	} else {
		config.writeFullBlocks(w, code)
	}
}

func Generate(text string, l qr.Level, w io.Writer) {
	config := Config{
		Level:     l,
		Writer:    w,
		BlackChar: black,
		WhiteChar: white,
		QuietZone: defaultQuiet,
	}
	GenerateWithConfig(text, config)
}

func GenerateHalfBlock(text string, l qr.Level, w io.Writer) {
	config := Config{
		Level:          l,
		Writer:         w,
		HalfBlocks:     true,
		BlackChar:      blackBlack,
		WhiteBlackChar: whiteBlack,
		WhiteChar:      whiteWhite,
		BlackWhiteChar: blackWhite,
		QuietZone:      defaultQuiet,
	}
	GenerateWithConfig(text, config)
}
