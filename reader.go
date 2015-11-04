package isofs

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/lunixbochs/struc"
)

var (
	// ErrInvalidImage is returned when an attempt to unpack the image Primary Volume Descriptor failed or
	// when the end of the image was reached without finding a primary volume descriptor
	ErrInvalidImage = func(err error) error { return fmt.Errorf("invalid-iso9660-image: %s", err) }
)

// Reader defines the state of the ISO9660 image reader. It needs to be instantiated
// from its constructor.
type Reader struct {
	imageReader     io.ReadSeeker
	currentDirEntry *DirectoryRecord
	primaryVolume   *PrimaryVolume
	buffer          []byte
	// pathTracker allow us to track nested directory paths so that the user
	// doesn't need to do recursion.
	pathTracker string
}

// NewReader creates a new ISO9660 reader.
func NewReader(rs io.ReadSeeker) (*Reader, error) {
	// Starts reading from image data area
	sector := dataAreaSector
	buffer := make([]byte, sectorSize)
	for {
		var volType uint8
		_, err := rs.Seek(int64(sector*sectorSize), 0)
		if err != nil {
			panic(err)
		}

		_, err = rs.Read(buffer)
		// If EOF is reached, it means a primary volume descriptor was not found.
		if err == io.EOF {
			return nil, ErrInvalidImage(err)
		}

		err = binary.Read(bytes.NewReader(buffer), binary.BigEndian, &volType)
		if err != nil {
			panic(err)
		}

		if volType == volSetTerminator {
			return nil, ErrInvalidImage(errors.New("Volume set terminator reached. A primary volume descriptor was not found."))
		}

		if volType == primaryVol {
			pvd := new(PrimaryVolume)
			if err := struc.Unpack(bytes.NewReader(buffer), pvd); err != nil {
				return nil, ErrInvalidImage(err)
			}

			return &Reader{
				imageReader:   rs,
				primaryVolume: pvd,
			}, nil
		}
		sector++
	}
}

// Next moves onto the next directory record if there is any
func (r *Reader) Next() bool {
	// Load directory record if currentDirEntry is nil
	// If not, load next dir extent
	return false
}

// Value returns file's information. If file is not a directory, Sys()
// returns an io.Reader for you to read the file's content from.
func (r *Reader) Value() (os.FileInfo, error) {
	return nil, nil
}
