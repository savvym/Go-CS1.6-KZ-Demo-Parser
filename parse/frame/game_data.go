package frame

import "go_demoParser/bitbuffer"

type GameDataFrame struct {
	gameDataFrameHeader
	movVars       *MoveVars
	serverMessage *[]byte
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
	g.serverMessage = &data
	return
}

func (g *GameDataFrame) GetServerMessage() []byte {
	return *g.serverMessage
}

func (g *GameDataFrame) GetMoveVars() *MoveVars {
	return g.movVars
}
