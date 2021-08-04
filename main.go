package main

import (
	"encoding/binary"
	"fmt"
	"go_demoParser/bitbuffer"
	"go_demoParser/parse"
	"go_demoParser/parse/frame"
	"io/ioutil"
)

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
		//fmt.Println(header)
		if header.Type() == 0 || header.Type() == 1 {
			var gameData frame.GameDataFrame
			err := gameData.Read(buffer)
			if err != nil {
				panic(err)
			}
			//fmt.Println(gameData.GetMoveVars())
			// 解析server message
			// serverMessage := gameData.GetServerMessage()
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

		} else if header.Type() == 9 {
			var demoBuffer frame.DemoBufferFrame
			err := demoBuffer.Read(buffer)
			if err != nil {
				panic(err)
			}
			//fmt.Println(demoBuffer)
		} else {
			length, err := header.GetFrameLength(buffer)
			if err != nil {
				panic(err)
			}
			//fmt.Println(length)
			buffer.Seek(length, 1)
		}
	}

}
