package frame

import "go_demoParser/bitbuffer"

// frame type = 7
// 8byte
type WeaponAnimFrame struct {
	anim uint32
	body uint32
}

func (w *WeaponAnimFrame) Read(buffer *bitbuffer.BitBuffer) (err error) {
	w.anim, err = buffer.ReadUint32(32)
	w.body, err = buffer.ReadUint32(32)
	return
}
