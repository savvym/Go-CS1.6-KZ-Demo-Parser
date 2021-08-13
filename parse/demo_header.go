package parse

import "go_demoParser/bitbuffer"

type DemoHeader struct {
	magic          string
	demoVersion    uint32
	networkVersion uint32
	mapName        string
	gameDll        string
	mapCRC         uint32
	dirOffset      uint32
}

func (dh *DemoHeader) ReadDemoHeader(buffer *bitbuffer.BitBuffer) {
	buffer.Seek(0, 0)
	dh.magic, _ = buffer.ReadString(8)
	dh.demoVersion, _ = buffer.ReadUint32(32)
	dh.networkVersion, _ = buffer.ReadUint32(32)
	dh.mapName, _ = buffer.ReadString(260)
	dh.gameDll, _ = buffer.ReadString(260)
	dh.mapCRC, _ = buffer.ReadUint32(32)
	dh.dirOffset, _ = buffer.ReadUint32(32)
}

func (dh *DemoHeader) Magic() string {
	return dh.magic
}

func (dh *DemoHeader) DemoVersion() uint32 {
	return dh.demoVersion
}

func (dh *DemoHeader) NerworkVersion() uint32 {
	return dh.networkVersion
}

func (dh *DemoHeader) MapName() string {
	return dh.mapName
}

func (dh *DemoHeader) GameDll() string {
	return dh.gameDll
}

func (dh *DemoHeader) MapCRC() uint32 {
	return dh.mapCRC
}

func (dh *DemoHeader) DirOffset() uint32 {
	return dh.dirOffset
}
