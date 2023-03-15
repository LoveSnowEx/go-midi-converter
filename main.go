package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"go-midi-converter/event"
	"go-midi-converter/note"
	"go-midi-converter/util"
	"io"
	"os"
)

func PrintBytes(s []byte) {
	fmt.Println(FormatBytes(s))
}

func FormatBytes(s []byte) string {
	return fmt.Sprintf("% x", s)
}

func VlqTest() {
	s := []uint32{
		0x00,
		0x40,
		0x7F,
		0x80,
		0x2000,
		0x3FFF,
		0x4000,
		0x100000,
		0x1FFFFF,
		0x200000,
		0x8000000,
		0xFFFFFFF,
	}
	for i := range s {
		PrintBytes(util.VlqEncode(s[i]))
	}
}

func NewHeader() []byte {
	headerHex := "4D546864000000060001000101E0"
	header, err := hex.DecodeString(headerHex)
	if err != nil {
		panic(err)
	}
	return header
}

func NewTrack(r io.Reader) []byte {
	var buf bytes.Buffer

	trackHex := "4D54726B"
	track, err := hex.DecodeString(trackHex)
	if err != nil {
		panic(err)
	}
	pc := event.ProgramChange(0, 0, 0x2a)
	notes := ReadNotation(r)

	size := util.Uint32ToBytes(uint32(len(pc) + len(notes)))

	buf.Write(track)
	buf.Write(size)
	buf.Write(pc)
	buf.Write(notes)
	return buf.Bytes()
}

func ReadNotation(r io.Reader) []byte {
	br := bufio.NewReader(r)
	type noteTigger struct {
		noteNum uint8
		beat    uint32
	}
	noteFuncs := [...]func(uint8) uint8{
		note.C,
		note.D,
		note.E,
		note.F,
		note.G,
		note.A,
		note.B,
	}
	s := []noteTigger{}
	for {
		b, err := br.ReadByte()
		if err != nil {
			break
		}
		switch {
		case b == '-':
			s[len(s)-1].beat++
		case b >= '1' && b <= '7':
			s = append(s, noteTigger{
				noteNum: noteFuncs[b-'1'](4),
				beat:    1,
			})
		}
	}
	bytes := []byte{}
	rem := uint32(0)
	const tpqn = 0x1E0
	const velocity = 0x20
	for i := range s {
		fmt.Printf("noteNum: %d, sec: %d\n", s[i].noteNum, s[i].beat)
		bytes = append(bytes,
			event.NoteOn(0, 0, s[i].noteNum, velocity)...)
		rem = s[i].beat
		bytes = append(bytes,
			event.NoteOff(rem*tpqn, 0, s[i].noteNum, velocity)...)
	}
	return bytes
}

func NewMidiWriter(fileName string) (w *bufio.Writer, close func() error) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	return bufio.NewWriter(file), file.Close
}

func main() {
	// VlqTest()

	in, err := os.OpenFile("in.txt", os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer in.Close()
	os.Stdin = in

	fileName := "out.mid"
	w, close := NewMidiWriter(fileName)
	if w == nil {
		return
	}
	defer close()

	bytes := append(NewHeader(), NewTrack(os.Stdin)...)
	PrintBytes(bytes)
	w.Write(bytes)
	w.Flush()
}
