package frame

import (
	"encoding/binary"
	"errors"
	"fmt"
	"go_demoParser/bitbuffer"
	"reflect"
)

var message = map[uint8]string{
	1:  "svc_nop",
	2:  "svc_disconnect",
	3:  "svc_event",
	4:  "svc_version",
	5:  "svc_setview",
	6:  "svc_sound",
	7:  "svc_time",
	8:  "svc_print",
	9:  "svc_stufftext",
	10: "svc_setangle",
	11: "svc_serverinfo",
	12: "svc_lightstyle",
	13: "svc_updateuserinfo",
	14: "svc_deltadescription",
	15: "svc_clientdata",
	16: "svc_stopsound",
	17: "svc_pings",
	18: "svc_particle",
	19: "svc_damage",
	20: "svc_spawnstatic",
	21: "svc_event_reliable",
	22: "svc_spawnbaseline",
	23: "svc_tempentity",
	24: "svc_setpause",
	25: "svc_signonnum",
	26: "svc_centerprint",
	27: "svc_killedmonster",
	28: "svc_foundsecret",
	29: "svc_spawnstaticsound",
	30: "svc_intermission",
	31: "svc_finale",
	32: "svc_cdtrack",
	33: "svc_restore",
	34: "svc_cutscene",
	35: "svc_weaponanim",
	36: "svc_decalname",
	37: "svc_roomtype",
	38: "svc_addangle",
	39: "svc_newusermsg",
	40: "svc_packetentities",
	41: "svc_deltapacketentities",
	42: "svc_choke",
	43: "svc_resourcelist",
	44: "svc_newmovevars",
	45: "svc_resourcerequest",
	46: "svc_customization",
	47: "svc_crosshairangle",
	48: "svc_soundfade",
	49: "svc_filetxferfailed",
	50: "svc_hltv",
	51: "svc_director",
	52: "svc_voiceinit",
	53: "svc_voicedata",
	54: "svc_sendextrainfo",
	55: "svc_timescale",
	56: "svc_resourcelocation",
	57: "svc_sendcvarvalue",
	58: "svc_sendcvarvalue2",
}

type GameDataFrame struct {
	gameDataFrameHeader
	movVars          *MoveVars
	serverMessageRaw *[]byte
	serverMessage    *map[string]interface{}
}

type gameDataFrameHeader struct {
	resolutionWidth  uint32
	resolutionHeight uint32
	length           uint32
}

type MoveVars struct {
	gravity           float32
	stopSpeed         float32
	maxSpeed          float32
	spectatorMaxSpeed float32
	accelerate        float32
	airAccelerate     float32
	waterAccelerate   float32
	friction          float32
	edgeFriction      float32
	waterFriction     float32
	entGravity        float32
	bounce            float32
	stepSize          float32
	maxVelocity       float32
	zMax              float32
	waveHeight        float32
	footsteps         uint32
	skyName           string
	rollAngle         float32
	rollSpeed         float32
	skyColorR         float32
	skyColorG         float32
	skyColorB         float32
	skyVecX           float32
	skyVecY           float32
	skyVecZ           float32
}

func (g *GameDataFrame) Read(buffer *bitbuffer.BitBuffer) (err error) {
	buffer.Seek(220, 1)
	// 220
	g.resolutionWidth, err = buffer.ReadUint32(32)
	g.resolutionHeight, err = buffer.ReadUint32(32)
	// 236
	buffer.Seek(60, 1)
	gravity, err := buffer.ReadFloat32()
	stopSpeed, err := buffer.ReadFloat32()
	maxSpeed, err := buffer.ReadFloat32()
	spectatorMaxSpeed, err := buffer.ReadFloat32()
	accelerate, err := buffer.ReadFloat32()
	airAccelerate, err := buffer.ReadFloat32()
	waterAccelerate, err := buffer.ReadFloat32()
	friction, err := buffer.ReadFloat32()
	edgeFriction, err := buffer.ReadFloat32()
	waterFriction, err := buffer.ReadFloat32()
	entGravity, err := buffer.ReadFloat32()
	bounce, err := buffer.ReadFloat32()
	stepSize, err := buffer.ReadFloat32()
	maxVelocity, err := buffer.ReadFloat32()
	zMax, err := buffer.ReadFloat32()
	waveHeight, err := buffer.ReadFloat32()
	footsteps, err := buffer.ReadUint32(32)
	skyName, err := buffer.ReadString(32)
	rollAngle, err := buffer.ReadFloat32()
	rollSpeed, err := buffer.ReadFloat32()
	skyColorR, err := buffer.ReadFloat32()
	skyColorG, err := buffer.ReadFloat32()
	skyColorB, err := buffer.ReadFloat32()
	skyVecX, err := buffer.ReadFloat32()
	skyVecY, err := buffer.ReadFloat32()
	skyVecZ, err := buffer.ReadFloat32()
	g.movVars = &MoveVars{
		gravity,
		stopSpeed,
		maxSpeed,
		spectatorMaxSpeed,
		accelerate,
		airAccelerate,
		waterAccelerate,
		friction,
		edgeFriction,
		waterFriction,
		entGravity,
		bounce,
		stepSize,
		maxVelocity,
		zMax,
		waveHeight,
		footsteps,
		skyName,
		rollAngle,
		rollSpeed,
		skyColorR,
		skyColorG,
		skyColorB,
		skyVecX,
		skyVecY,
		skyVecZ,
	}
	buffer.Seek(44, 1)
	g.length, err = buffer.ReadUint32(32)
	data, err := buffer.Read(uint64(g.length) * 8)
	g.serverMessageRaw = &data
	return
}

func (g *GameDataFrame) GetServerMessage() []byte {
	return *g.serverMessageRaw
}

func (g *GameDataFrame) GetMoveVars() *MoveVars {
	return g.movVars
}

var messageHandler = map[string]interface{}{
	"svc_sendcvarvalue2": messageSendCvarValue2,
	"svc_time":           messageSvcTime,
	"svc_clientdata":     messageClientData,
}

func (g *GameDataFrame) ParseServerMessage() {
	buffer := bitbuffer.NewBitBuffer(binary.LittleEndian)
	buffer.Feed(*g.serverMessageRaw)
	//readingGameData := true
	for {
		messageId, err := buffer.ReadUint8(8)
		if err != nil {
			return
		}
		messageName := message[messageId]
		fmt.Println(messageId, messageName)
		if messageHandler[messageName] == nil {
			panic("no such messageHandler")
		}
		call(messageName, buffer)
	}
	return
}

func messageClientData(buffer *bitbuffer.BitBuffer) (err error) {
	deltaSequence, err := buffer.ReadBoolean()
	if deltaSequence {
		buffer.Read(8)
	}
	// 持续施工中...
	return
}

func messageSendCvarValue2(buffer *bitbuffer.BitBuffer) (err error) {
	buffer.Seek(4, 1)
	str, err := buffer.ReadStringToEnd()
	fmt.Println(str)
	return
}

func messageSvcTime(buffer *bitbuffer.BitBuffer) (err error) {
	buffer.Seek(4, 1)
	return
}

func call(funcName string, params ...interface{}) (result interface{}, err error) {
	f := reflect.ValueOf(messageHandler[funcName])
	if len(params) != f.Type().NumIn() {
		err = errors.New("The number of params is out of index.")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	var res []reflect.Value
	res = f.Call(in)
	result = res[0].Interface()
	return
}
