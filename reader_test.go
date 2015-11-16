package iso9660

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/hooklift/assert"
)

func TestNewReader(t *testing.T) {
	image, err := os.Open("./fixtures/test.iso")
	defer image.Close()
	r, err := NewReader(image)
	assert.Ok(t, err)
	// Test first half of primary volume descriptor
	assert.Equals(t, "CD001", string(r.pvd.StandardID[:]))
	assert.Equals(t, 1, int(r.pvd.Type))
	assert.Equals(t, 1, int(r.pvd.Version))
	assert.Equals(t, "Mac OS X", strings.TrimSpace(string(r.pvd.SystemID[:])))
	assert.Equals(t, "my-vol-id", strings.TrimSpace(string(r.pvd.ID[:])))
	assert.Equals(t, 181, int(r.pvd.VolumeSpaceSizeBE))
	assert.Equals(t, 1, int(r.pvd.VolumeSetSizeBE))
	assert.Equals(t, 1, int(r.pvd.VolumeSeqNumberBE))
	assert.Equals(t, 2048, int(r.pvd.LogicalBlkSizeBE))
	assert.Equals(t, 46, int(r.pvd.PathTableSizeBE))
	assert.Equals(t, 21, int(r.pvd.LocPathTableBE))
	assert.Equals(t, 0, int(r.pvd.LocOptPathTableBE))
	// Test root directory record values
	assert.Equals(t, 0, int(r.pvd.DirectoryRecord.ExtendedAttrLen))
	assert.Equals(t, 23, int(r.pvd.DirectoryRecord.ExtentLocationBE))
	assert.Equals(t, 2048, int(r.pvd.DirectoryRecord.ExtentLengthBE))
	assert.Equals(t, 2, int(r.pvd.DirectoryRecord.FileFlags))
	assert.Equals(t, 0, int(r.pvd.DirectoryRecord.FileUnitSize))
	assert.Equals(t, 0, int(r.pvd.DirectoryRecord.InterleaveGapSize))
	assert.Equals(t, 1, int(r.pvd.DirectoryRecord.VolumeSeqNumberBE))
	assert.Equals(t, 1, int(r.pvd.DirectoryRecord.FileIDLength))
	// Test second half of primary volume descriptor
	assert.Equals(t, "my-vol-id", strings.TrimSpace(string(r.pvd.ID[:])))
	assert.Equals(t, "test-volset-id", strings.TrimSpace(string(r.pvd.VolumeSetID[:])))
	assert.Equals(t, "hooklift", strings.TrimSpace(string(r.pvd.PublisherID[:])))
	assert.Equals(t, "hooklift", strings.TrimSpace(string(r.pvd.DataPreparerID[:])))
	assert.Equals(t, "MKISOFS ISO9660/HFS/UDF FILESYSTEM BUILDER & CDRECORD CD/DVD/BluRay CREATOR (C) 1993 E.YOUNGDALE (C) 1997 J.PEARSON/J.SCHILLING", strings.TrimSpace(string(r.pvd.AppID[:])))
	assert.Equals(t, 1, int(r.pvd.FileStructVersion))
}

func TestUnpackChildren(t *testing.T) {
	image, err := os.Open("./fixtures/test.iso")
	defer image.Close()
	reader, err := NewReader(image)
	assert.Ok(t, err)

	count := 0
	for {
		count++
		fi, err := reader.Next()
		if err == io.EOF {
			break
		}
		assert.Ok(t, err)

		f := fi.(*FileStat)
		fmt.Printf("%d->%#v<-\n", count, f.Name())
	}
}
