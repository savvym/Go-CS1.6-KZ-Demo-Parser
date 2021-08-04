package frame

import "go_demoParser/bitbuffer"

// frame type = 4
// 32byte
type ClientDataFrame struct {
	origin     [3]float32
	viewAngles [3]float32
	weaponBits int32
	fov        float32
}

func (c *ClientDataFrame) Read(buffer *bitbuffer.BitBuffer) (err error) {
	o1, err := buffer.ReadFloat32()
	o2, err := buffer.ReadFloat32()
	o3, err := buffer.ReadFloat32()
	c.origin = [3]float32{o1, o2, o3}

	v1, err := buffer.ReadFloat32()
	v2, err := buffer.ReadFloat32()
	v3, err := buffer.ReadFloat32()
	c.viewAngles = [3]float32{v1, v2, v3}

	w, err := buffer.ReadInt32(32)
	c.weaponBits = w

	fov, err := buffer.ReadFloat32()
	c.fov = fov

	return
}
