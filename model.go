package main

import "github.com/hazelcast/hazelcast-go-client/serialization"

type Reading struct {
	Name      string
	Timestamp int64
}

const (
	dataClassID      = 1
	portableFactorID = 1
)

func (d *Reading) ClassID() int32 {
	return dataClassID
}

func (d *Reading) FactoryID() int32 {
	return portableFactorID
}

func (d *Reading) WritePortable(writer serialization.PortableWriter) error {
	writer.WriteUTF("Name", d.Name)
	writer.WriteInt64("Timestamp", d.Timestamp)
	return nil
}

func (d *Reading) ReadPortable(reader serialization.PortableReader) error {
	d.Name = reader.ReadUTF("Name")
	d.Timestamp = reader.ReadInt64("Timestamp")
	return reader.Error()
}

type DataPortableFactory struct{}

func (pf *DataPortableFactory) Create(classID int32) serialization.Portable {
	if classID == dataClassID {
		return &Reading{}
	}
	return nil
}
