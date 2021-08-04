package frame

import "go_demoParser/bitbuffer"

//const (
//	NetMsg0			uint8 = 0
//	NetMsg1			uint8 = 1
//	DemoStart      	uint8 = 2
//	ConsoleCommand 	uint8 = 3
//	ClientData     	uint8 = 4
//	NextSection    	uint8 = 5
//	Event          	uint8 = 6
//	WeaponAnim     	uint8 = 7
//	Sound          	uint8 = 8
//	DemoBuffer     	uint8 = 9
//)

type Header struct {
	frameType uint8
	time      float32
	number    uint32
}

func (h *Header) Read(buffer *bitbuffer.BitBuffer) (err error) {
	h.frameType, err = buffer.ReadUint8(8)
	h.time, err = buffer.ReadFloat32()
	h.number, err = buffer.ReadUint32(32)
	return
}

func (h *Header) Type() uint8 {
	return h.frameType
}

func (h *Header) Time() float32 {
	return h.time
}

func (h *Header) Number() uint32 {
	return h.number
}
