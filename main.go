package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"go_demoParser/bitbuffer"
	"go_demoParser/parse"
	"io/ioutil"
)

type FrameHeader struct {
	frameType uint8
	time      float32
	number    uint32
}

type GameDataFrameHeader struct {
	resolutionWidth  uint32
	resolutionHeight uint32
	length           uint32
}

func main() {
	file, err := ioutil.ReadFile("/Users/zhanghaodong/Downloads/2.dem")
	if err != nil {
		panic(err)
	}
	buffer := bitbuffer.NewBitBuffer(binary.LittleEndian)
	buffer.Feed(file)

	// 开始解析文件
	// 解析demo头
	var demoHeader parse.DemoHeader
	demoHeader.ReadDemoHeader(buffer)
	fmt.Println(demoHeader)

	// 解析dirEntry
	var dirEntry parse.DirectoryEntry
	dirEntry.ReadDirEntry(buffer)
	fmt.Println(dirEntry)

	off := dirEntry.EntryList()[1].Offset()
	buffer.Seek(int32(off), 0)
	for {
		header := ReadFrameHeader(buffer)
		fmt.Println(header)
		if header.frameType == 0 || header.frameType == 1 {
			gameDataFrameHeader := ReadGameDataFrameHeader(buffer)
			//fmt.Println(buffer.BytePos())
			// GameData
			buffer.Read(uint64(gameDataFrameHeader.length) * 8)
			//fmt.Println(len(data),data)
			// 解析GameData...

			// ...
		} else if header.frameType == 5 {
			break
		} else {
			length, err := GetFrameLength(header.frameType, buffer)
			if err != nil {
				panic(err)
			}
			//fmt.Println(length)
			buffer.Seek(length, 1)
		}

		//else if header.frameType == 8{
		//	channel, _ := buffer.ReadUint32(32)
		//	fmt.Println("channel", channel)
		//	l, _ := buffer.ReadInt32(32)
		//	data, _ := buffer.ReadString(uint64(l))
		//	fmt.Println("data", data)
		//	attenuation, _ := buffer.ReadFloat32()
		//	fmt.Println("attenuation", attenuation)
		//	volume, _ := buffer.ReadFloat32()
		//	fmt.Println("volume", volume)
		//	flags, _ := buffer.ReadUint32(32)
		//	fmt.Println("flags", flags)
		//	pitch, _ := buffer.ReadUint32(32)
		//	fmt.Println("pitch", pitch)
		//}
	}

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
	buffer.Seek(236, 1)
	length, _ := buffer.ReadUint32(32)
	return &GameDataFrameHeader{resolutionWidth, resolutionHeight, length}
}

//func ParseGameDataMessages(frameData []byte) {
//	buffer := bitbuffer.NewBitBuffer(binary.LittleEndian)
//	buffer.Feed(frameData)
//
//}

func GetFrameLength(frameType uint8, buffer *bitbuffer.BitBuffer) (length int32, err error) {
	switch frameType {
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
