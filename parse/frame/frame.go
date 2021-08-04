package frame

//const (
//	NetMsg0			uint8 = 0
//	NetMsg1			uint8 = 1
//	DemoStart      	uint8 = 2
//	ConsoleCommand 	uint8 = 3
//	ClientData     	uint8 = 4
//	NextSection    	uint8 = 5
//	Event          	uint8 = 6
//	WeaponAnim     	uint8 = 7
//	Sound          	uint8 = 8
//	DemoBuffer     	uint8 = 9
//)

type Header struct {
	frameType uint8
	time      float32
	number    uint32
}

type ConsoleCommandFrame struct {
	command string // 64byte
}

// 32byte
type ClientDataFrame struct {
	origin     [3]float32
	viewAngles [3]float32
	weaponBits int32
	fov        float32
}

// 84byte
type EventFrame struct {
	flags uint32
	index uint32
	delay float32
	args  EventArgs
}

type EventArgs struct {
	flags       uint32
	entityIndex uint32
	origin      [3]float32
	angles      [3]float32
	velocity    [3]float32
	ducking     uint32
	fparam1     float32
	fparam2     float32
	iparam1     uint32
	iparam2     uint32
	bparam1     uint32
	bparam2     uint32
}

// 8byte
type WeaponAnimFrame struct {
	anim uint32
	body uint32
}

//
type SoundFrame struct {
	channel         uint32
	soundNameLength int32  // sound_name length
	soundName       string // sound_name
	attenuation     float32
	volume          float32
	flags           uint32
	pitch           uint32
}
