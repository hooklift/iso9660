package iso9660

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	// ErrInvalidImage is returned when an attempt to unpack the image Primary Volume Descriptor failed or
	// when the end of the image was reached without finding a primary volume descriptor
	ErrInvalidImage = func(err error) error { return fmt.Errorf("invalid-iso9660-image: %s", err) }
	// ErrCorruptedImage is returned when a seek operation, on the image, failed.
	ErrCorruptedImage = func(err error) error { return fmt.Errorf("corrupted-image: %s", err) }
)

// Reader defines the state of the ISO9660 image reader. It needs to be instantiated
// from its constructor.
type Reader struct {
	image io.ReadSeeker
	pvd   PrimaryVolume
}

// NewReader creates a new ISO9660 reader.
func NewReader(rs io.ReadSeeker) (*Reader, error) {
	// Starts reading from image data area
	sector := dataAreaSector
	// Iterates over volume descriptors until it finds the primary volume descriptor
	// or an error condition.
	for {
		offset, err := rs.Seek(int64(sector*sectorSize), os.SEEK_SET)
		if err != nil {
			return nil, ErrCorruptedImage(err)
		}

		var volDesc VolumeDescriptor
		if err := binary.Read(rs, binary.BigEndian, &volDesc); err != nil {
			return nil, ErrCorruptedImage(err)
		}

		if volDesc.Type == primaryVol {
			if _, err := rs.Seek(offset, os.SEEK_SET); err != nil {
				return nil, ErrCorruptedImage(err)
			}

			reader := new(Reader)
			reader.image = rs

			if err := reader.unpackPVD(); err != nil {
				return nil, ErrCorruptedImage(err)
			}

			return reader, nil
		}

		if volDesc.Type == volSetTerminator {
			return nil, ErrInvalidImage(errors.New("Volume Set Terminator reached. A Primary Volume Descriptor was not found."))
		}
		sector++
	}
}

// Next moves onto the next directory record if there is any
func (r *Reader) Next() (os.FileInfo, error) {
	if r == nil {
		panic("missing reader instance")
	}

	drecord := new(DirectoryRecord)
	if err := r.unpackDRecord(drecord); err != nil {
		return nil, err
	}

	fi := &FileStat{image: r.image, DirectoryRecord: *drecord}

	return fi, nil
}

// unpackPVD unpacks Primary Volume Descriptor in three phases. This is
// because the root directory record is a variable-length record and Go's binary
// package doesn't support unpacking variable-length structs easily.
func (r *Reader) unpackPVD() error {
	// Unpack first half
	var pvd1 PrimaryVolumePart1
	if err := binary.Read(r.image, binary.BigEndian, &pvd1); err != nil {
		return ErrCorruptedImage(err)
	}
	r.pvd.PrimaryVolumePart1 = pvd1

	// Unpack root directory record
	var drecord DirectoryRecord
	if err := r.unpackDRecord(&drecord); err != nil {
		return ErrCorruptedImage(err)
	}
	r.pvd.DirectoryRecord = drecord

	// Unpack second half
	var pvd2 PrimaryVolumePart2
	if err := binary.Read(r.image, binary.BigEndian, &pvd2); err != nil {
		return ErrCorruptedImage(err)
	}
	r.pvd.PrimaryVolumePart2 = pvd2

	return nil
}

func (r *Reader) unpackDRecord(drecord *DirectoryRecord) error {
	var len byte
	if err := binary.Read(r.image, binary.BigEndian, &len); err != nil {
		return ErrCorruptedImage(err)
	}

	if len == 0 {
		return nil
	}

	if err := binary.Read(r.image, binary.BigEndian, drecord); err != nil {
		return ErrCorruptedImage(err)
	}

	name := make([]byte, drecord.FileIDLength)
	if err := binary.Read(r.image, binary.BigEndian, name); err != nil {
		return ErrCorruptedImage(err)
	}
	//drecord.FileID = string(name)
	//fmt.Printf("name: ->%s<-\n", name)

	//Padding field as per section 9.1.12 in ECMA-119
	if (drecord.FileIDLength % 2) == 0 {
		var zero byte
		if err := binary.Read(r.image, binary.BigEndian, &zero); err != nil {
			return ErrCorruptedImage(err)
		}
	}

	// System use field as per section 9.1.13 in ECMA-119
	totalLen := 34 + drecord.FileIDLength - (drecord.FileIDLength % 2)
	sysUseLen := int64(len - totalLen)
	if sysUseLen > 0 {
		sysData := make([]byte, sysUseLen)
		if err := binary.Read(r.image, binary.BigEndian, sysData); err != nil {
			return ErrCorruptedImage(err)
		}
	}
	return nil
}
