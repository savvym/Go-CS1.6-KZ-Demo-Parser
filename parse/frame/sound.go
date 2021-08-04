package frame

import "go_demoParser/bitbuffer"

// frame type = 8
type SoundFrame struct {
	channel         uint32
	soundNameLength int32  // sound_name length
	soundName       string // sound_name
	attenuation     float32
	volume          float32
	flags           uint32
	pitch           uint32
}

func (s *SoundFrame) Read(buffer *bitbuffer.BitBuffer) (err error) {
	s.channel, err = buffer.ReadUint32(32)
	s.soundNameLength, err = buffer.ReadInt32(32)
	s.soundName, err = buffer.ReadString(uint64(s.soundNameLength))
	s.attenuation, err = buffer.ReadFloat32()
	s.volume, err = buffer.ReadFloat32()
	s.flags, err = buffer.ReadUint32(32)
	s.pitch, err = buffer.ReadUint32(32)
	return
}
