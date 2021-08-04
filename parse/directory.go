package parse

import (
	"go_demoParser/bitbuffer"
)

type DirectoryEntry struct {
	number    uint32
	entryList []*Entry
}

type Entry struct {
	dType  uint32
	title  string
	flags  uint32
	play   int32
	time   float32
	frames uint32
	offset uint32
	length uint32
}

func (dir *DirectoryEntry) ReadDirEntry(buffer *bitbuffer.BitBuffer) {
	buffer.Seek(540, 0)
	dirOffset, _ := buffer.ReadUint32(32)
	buffer.Seek(int32(dirOffset), 0)
	number, _ := buffer.ReadUint32(32)
	for i := uint32(0); i < number; i++ {
		var e Entry
		e.readEntry(buffer)
		dir.entryList = append(dir.entryList, &e)
	}
}

func (e *Entry) readEntry(buffer *bitbuffer.BitBuffer) {
	e.dType, _ = buffer.ReadUint32(32)
	e.title, _ = buffer.ReadString(64)
	e.flags, _ = buffer.ReadUint32(32)
	e.play, _ = buffer.ReadInt32(32)
	e.time, _ = buffer.ReadFloat32()
	e.frames, _ = buffer.ReadUint32(32)
	e.offset, _ = buffer.ReadUint32(32)
	e.length, _ = buffer.ReadUint32(32)
}

func (dir *DirectoryEntry) Count() uint32 {
	return dir.number
}

func (dir *DirectoryEntry) EntryList() []*Entry {
	return dir.entryList
}

func (e *Entry) Type() uint32 {
	return e.dType
}

func (e *Entry) Description() string {
	return e.title
}

func (e *Entry) Flags() uint32 {
	return e.flags
}

func (e *Entry) Play() int32 {
	return e.play
}

func (e *Entry) Time() float32 {
	return e.time
}

func (e *Entry) Frames() uint32 {
	return e.frames
}

func (e *Entry) Offset() uint32 {
	return e.offset
}

func (e *Entry) Length() uint32 {
	return e.length
}
