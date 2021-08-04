package frame

import (
	"errors"
	"go_demoParser/bitbuffer"
)

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

func (h *Header) GetFrameLength(buffer *bitbuffer.BitBuffer) (length int32, err error) {
	switch h.frameType {
	case 2: // START
		err = nil
		break
	case 3: // command
		length = 64
		err = nil
		break
	case 4: // client data
		length = 32
		err = nil
		break
	case 5: // end of segment
		break
	case 6: // event
		length = 84
		err = nil
		break
	case 7:
		length = 8
		err = nil
		break
	case 8:
		buffer.Seek(4, 1)
		l, _ := buffer.ReadInt32(32)
		length = l
		buffer.Seek(-8, 1)
		length += 24
		err = nil
		break
	case 9:
		l, _ := buffer.ReadInt32(32)
		length = l + 4
		buffer.Seek(-4, 1)
		err = nil
		break
	default:
		err = errors.New("Unknown parse type ")
	}
	return
}
