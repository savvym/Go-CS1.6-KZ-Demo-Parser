package frame

import "go_demoParser/bitbuffer"

// frame type = 6
// 84byte
type EventFrame struct {
	flags uint32
	index uint32
	delay float32
	args  eventArgs
}

type eventArgs struct {
	flags       uint32
	entityIndex uint32
	origin      [3]float32
	angles      [3]float32
	velocity    [3]float32
	ducking     uint32
	// skip 48byte
	//fparam1     float32
	//fparam2     float32
	//iparam1     uint32
	//iparam2     uint32
	//bparam1     uint32
	//bparam2     uint32
}

func (ef *EventFrame) Read(buffer *bitbuffer.BitBuffer) (err error) {
	ef.flags, err = buffer.ReadUint32(32)
	ef.index, err = buffer.ReadUint32(32)
	ef.delay, err = buffer.ReadFloat32()
	var args eventArgs
	err = args.read(buffer)
	ef.args = args
	return
}

func (ea *eventArgs) read(buffer *bitbuffer.BitBuffer) (err error) {
	ea.flags, err = buffer.ReadUint32(32)
	ea.entityIndex, err = buffer.ReadUint32(32)

	o1, err := buffer.ReadFloat32()
	o2, err := buffer.ReadFloat32()
	o3, err := buffer.ReadFloat32()
	ea.origin = [3]float32{o1, o2, o3}

	a1, err := buffer.ReadFloat32()
	a2, err := buffer.ReadFloat32()
	a3, err := buffer.ReadFloat32()
	ea.angles = [3]float32{a1, a2, a3}

	v1, err := buffer.ReadFloat32()
	v2, err := buffer.ReadFloat32()
	v3, err := buffer.ReadFloat32()
	ea.velocity = [3]float32{v1, v2, v3}

	ea.ducking, err = buffer.ReadUint32(32)
	buffer.Seek(24, 1)
	return
}
