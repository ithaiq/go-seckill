package gnet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"ithaiq/gtcp/giface"
	"ithaiq/gtcp/utils"
)

type DataPack struct {
}

func NewDataPack() *DataPack {
	return new(DataPack)
}

func (d *DataPack) GetHeadLen() uint32 {
	//DataLen + ID
	return 8
}

func (d *DataPack) Pack(msg giface.IMessage) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	if err := binary.Write(buffer, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (d *DataPack) UnPack(data []byte) (giface.IMessage, error) {
	buffer := bytes.NewReader(data)
	msg := &Message{}
	if err := binary.Read(buffer, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	if err := binary.Read(buffer, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("recv too large msg data")
	}
	if err := binary.Read(buffer, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	return msg, nil
}
