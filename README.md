# iso9660
[![GoDoc](https://godoc.org/github.com/hooklift/iso9660?status.svg)](https://godoc.org/github.com/hooklift/iso9660)
[![Build Status](https://travis-ci.org/hooklift/iso9660.svg?branch=master)](https://travis-ci.org/hooklift/iso9660)

Go library and CLI to extract data from ISO9660 images.

### CLI

### Library
An example on how to use the library to extract files and directories from an ISO image can be found in our CLI source code at:
https://github.com/hooklift/iso9660/blob/master/cmd/iso9660/main.go

### Not supported
* Reading files recorded in interleave mode
* Multi-extent or reading files larger than 4GB
* Joliet extension
* Multisession extension
* Rock Ridge extension
* El Torito extension
