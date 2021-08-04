package frame

import "go_demoParser/bitbuffer"

// frame type = 9
type DemoBufferFrame struct {
	length int32
	data   string
}

func (d *DemoBufferFrame) Read(buffer *bitbuffer.BitBuffer) (err error) {
	d.length, err = buffer.ReadInt32(32)
	d.data, err = buffer.ReadString(uint64(d.length))
	return
}
