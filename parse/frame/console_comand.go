package frame

import "go_demoParser/bitbuffer"

// frame type = 3
type ConsoleCommandFrame struct {
	command string // 64byte
}

func (c *ConsoleCommandFrame) Read(buffer *bitbuffer.BitBuffer) (err error) {
	c.command, err = buffer.ReadString(64)
	return
}
