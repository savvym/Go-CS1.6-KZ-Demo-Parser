package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"go_demoParser/bitbuffer"
	"go_demoParser/parse"
	"go_demoParser/parse/frame"
	"io/ioutil"
)

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
	//fmt.Println(dirEntry)

	off := dirEntry.EntryList()[1].Offset()
	buffer.Seek(int32(off), 0)
	for {
		var header frame.Header
		err := header.Read(buffer)
		if err != nil {
			panic(err)
		}
		fmt.Println(header)
		if header.Type() == 0 || header.Type() == 1 {
			gameDataFrameHeader := ReadGameDataFrameHeader(buffer)
			//fmt.Println(buffer.BytePos())
			// GameData
			_, _ = buffer.Read(uint64(gameDataFrameHeader.length) * 8)
			//fmt.Println(len(data),data)
			// 解析GameData...

			// ...
		} else if header.Type() == 3 {
			var consoleCommand frame.ConsoleCommandFrame
			err := consoleCommand.Read(buffer)
			if err != nil {
				panic(err)
			}
			//fmt.Println(consoleCommand)

		} else if header.Type() == 4 {
			var clientData frame.ClientDataFrame
			err := clientData.Read(buffer)
			if err != nil {
				panic(err)
			}
			//fmt.Println(clientData)

		} else if header.Type() == 5 {
			break
		} else if header.Type() == 6 {
			var event frame.EventFrame
			err := event.Read(buffer)
			if err != nil {
				panic(err)
			}
			//fmt.Println(event)

		} else if header.Type() == 7 {
			var weaponAnim frame.WeaponAnimFrame
			err := weaponAnim.Read(buffer)
			if err != nil {
				panic(err)
			}
			//fmt.Println(weaponAnim)

		} else if header.Type() == 8 {
			var sound frame.SoundFrame
			err := sound.Read(buffer)
			if err != nil {
				panic(err)
			}
			//fmt.Println(sound)

		} else {
			length, err := GetFrameLength(header.Type(), buffer)
			if err != nil {
				panic(err)
			}
			//fmt.Println(length)
			buffer.Seek(length, 1)
		}
	}

}

func ReadGameDataFrameHeader(buffer *bitbuffer.BitBuffer) *GameDataFrameHeader {
	buffer.Seek(220, 1)
	resolutionWidth, _ := buffer.ReadUint32(32)
	resolutionHeight, _ := buffer.ReadUint32(32)
	buffer.Seek(236, 1)
	length, _ := buffer.ReadUint32(32)
	return &GameDataFrameHeader{resolutionWidth, resolutionHeight, length}
}

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
