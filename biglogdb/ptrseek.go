package biglogdb

import (
	"encoding/binary"
	"fmt"
)

const (
	LogFilePtrSeekSize = 4 + 4 + 4
)

type LogFilePtrSeek struct {
	Filenum  uint32
	Fileseek uint32
	Valsize  uint32
}

func (lp *LogFilePtrSeek) Copy() *LogFilePtrSeek {
	return &LogFilePtrSeek{
		Filenum:  lp.Filenum,
		Fileseek: lp.Fileseek,
		Valsize:  lp.Valsize,
	}
}

//////////////////////////////////////

// assembling datas
func (lp *LogFilePtrSeek) Serialize() ([]byte, error) {
	data := make([]byte, LogFilePtrSeekSize)
	binary.BigEndian.PutUint32(data[0:4], lp.Filenum)
	binary.BigEndian.PutUint32(data[4:8], lp.Fileseek)
	binary.BigEndian.PutUint32(data[8:12], lp.Valsize)
	return data, nil
}

func (lp *LogFilePtrSeek) Parse(data []byte, seek uint32) (uint32, error) {
	if len(data) < int(seek)+LogFilePtrSeekSize {
		return 0, fmt.Errorf("data size error")
	}
	lp.Filenum = binary.BigEndian.Uint32(data[seek : seek+4])
	lp.Fileseek = binary.BigEndian.Uint32(data[seek+4 : seek+8])
	lp.Valsize = binary.BigEndian.Uint32(data[seek+8 : seek+12])
	return seek + LogFilePtrSeekSize, nil
}

func (lp *LogFilePtrSeek) Size() uint32 {
	return LogFilePtrSeekSize
}
