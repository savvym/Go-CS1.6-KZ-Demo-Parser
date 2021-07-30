package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"go_demoParser/bitbuffer"
	"io/ioutil"
)


type DemoHeader struct {
	magic			string
	demoVersion 	uint32
	networkVersion	uint32
	mapName			string
	gameDll			string
	mapCRC			uint32
	dirOffset		uint32
}

type DemoDirEntry struct {
	number    uint32
	entryList []DirEntry
}

type DirEntry struct {
	dType		uint32
	title		string
	flags 		uint32
	play		int32
	time 		float32
	frames		uint32
	offset		uint32
	length		uint32
}

type FrameHeader struct {
	frameType uint8
	time      float32
	number    uint32
}

type GameDataFrameHeader struct {
	resolutionWidth		uint32
	resolutionHeight	uint32
	length				uint32
}

func main() {
	file, err := ioutil.ReadFile("/Users/zhanghaodong/Downloads/2.dem")
	if err != nil  {
		panic(err)
	}
	buffer := bitbuffer.NewBitBuffer(binary.LittleEndian)
	buffer.Feed(file)

	// 开始解析文件
	// 解析demo头
	demoHeader := ReadDemoHeader(buffer)
	fmt.Println(demoHeader)

	// 解析dirEntry
	dirEntry := ReadDirEntry(buffer)
	fmt.Println(dirEntry)

	off := dirEntry.entryList[1].offset
	buffer.Seek(int32(off),0)
	for {
		header := ReadFrameHeader(buffer)
		fmt.Println(header)
		if header.frameType == 0 || header.frameType == 1 {
			gameDataFrameHeader := ReadGameDataFrameHeader(buffer)
			//fmt.Println(gameDataFrameHeader)
			buffer.Seek(int32(gameDataFrameHeader.length),1)
		}else if header.frameType == 5 {
			break
		} else {
			length, err := GetFrameLength(header.frameType, buffer)
			if err != nil {
				panic(err)
			}
			//fmt.Println(length)
			buffer.Seek(length,1)
		}
	}


}



func ReadDemoHeader(buffer *bitbuffer.BitBuffer) *DemoHeader{
	buffer.Seek(0,0)
	magic, _ := buffer.ReadString(8)
	demoVersion, _ := buffer.ReadUint32(32)
	networkVersion, _ := buffer.ReadUint32(32)
	mapName, _ := buffer.ReadString(260)
	gameDll, _ := buffer.ReadString(260)
	mapCRC, _ := buffer.ReadUint32(32)
	dirOffset, _ := buffer.ReadUint32(32)
	return &DemoHeader{
		magic,
		demoVersion,
		networkVersion,
		mapName,
		gameDll,
		mapCRC,
		dirOffset,
	}
}

func ReadDirEntry(buffer *bitbuffer.BitBuffer) *DemoDirEntry {
	buffer.Seek(540,0)
	dirOffset, _ := buffer.ReadUint32(32)
	buffer.Seek(int32(dirOffset), 0)
	number, _ := buffer.ReadUint32(32)
	fmt.Println(number)
	demoDirEntry := &DemoDirEntry{
		number:    number,
		entryList: nil,
	}
	for i := uint32(0); i < number; i++ {

		dType, _ := buffer.ReadUint32(32)
		title, _ := buffer.ReadString(64)
		flags, _ := buffer.ReadUint32(32)
		play, _ := buffer.ReadInt32(32)
		time, _ := buffer.ReadFloat32()
		frames, _ := buffer.ReadUint32(32)
		offset, _ := buffer.ReadUint32(32)
		length, _ := buffer.ReadUint32(32)
		dirEntry := DirEntry{
			dType,
			title,
			flags,
			play,
			time,
			frames,
			offset,
			length,
		}
		demoDirEntry.entryList = append(demoDirEntry.entryList, dirEntry)
	}
	return demoDirEntry
}

func ReadFrameHeader(buffer *bitbuffer.BitBuffer) *FrameHeader {
	frameType, _ := buffer.ReadUint8(8)
	time, _ := buffer.ReadFloat32()
	number, _ := buffer.ReadUint32(32)
	header := &FrameHeader{
		frameType,
		time,
		number,
	}
	return header
}

func ReadGameDataFrameHeader(buffer *bitbuffer.BitBuffer) *GameDataFrameHeader {
	buffer.Seek(220, 1)
	resolutionWidth, _ := buffer.ReadUint32(32)
	resolutionHeight, _ := buffer.ReadUint32(32)
	buffer.Seek(236,1)
	length, _ := buffer.ReadUint32(32)
	return &GameDataFrameHeader{resolutionWidth, resolutionHeight, length}
}

func GetFrameLength(frameType uint8, buffer *bitbuffer.BitBuffer) (length int32, err error) {
	switch frameType {
	case 2:	// ???
		err = nil
		break
	case 3:
		length = 64
		err = nil
		break
	case 4:
		length = 32
		err = nil
		break
	case 5:		// end of segment
		break
	case 6:
		length = 84
		err = nil
		break
	case 7:
		length = 8
		err = nil
		break
	case 8:
		buffer.Seek(4,1)
		l, _ := buffer.ReadInt32(32)
		length = l
		buffer.Seek(-8,1)
		length += 24
		err = nil
		break
	case 9:
		l, _ := buffer.ReadInt32(32)
		length = l + 4
		buffer.Seek(-4,1)
		err = nil
		break
	default:
		err = errors.New("Unknown frame type ")
	}
	return
}
