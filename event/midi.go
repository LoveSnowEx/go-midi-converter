package event

import "go-midi-converter/util"

func NoteOn(deltaTime uint32, c, noteNum, velocity uint8) (out []byte) {
	out = append(util.VlqEncode(deltaTime), 0x90|c, noteNum, velocity)
	return
}

func NoteOff(deltaTime uint32, c, noteNum, velocity uint8) (out []byte) {
	out = append(util.VlqEncode(deltaTime), 0x80|c, noteNum, velocity)
	return
}

func ProgramChange(deltaTime uint32, c, programNum uint8) (out []byte) {
	out = append(util.VlqEncode(deltaTime), 0xC0|c, programNum)
	return
}
