package memtable

import (
	"github.com/fabulousduck/smol/errors"
)

/*
MemTable is a simple collection of memory regions in use
*/
type MemTable []*MemRegion

/*
MemRegion represents a region of memory in the IR
*/
type MemRegion struct {
	Addr, Size, Value int
	Name              string
}

/*
Find finds a variable on the memory table
Returns the index at which it is found
Returns -1 if the value cannot be found
*/
func (table MemTable) Find(name string) int {
	for i := 0; i < len(table); i++ {
		region := table[i]
		if region.Name == name {
			return i
		}
	}
	return -1
}

/*
Put places a variable on a memoryTable

notes

chip-8's blocks are 8 bit, so 1 byte.
with a total of 4096 bytes
*/
func (table *MemTable) Put(name string, value int) *MemRegion {
	region := new(MemRegion)
	//check if there is any memory left for our variable
	currentMemSize := table.getSize()
	if currentMemSize >= 4096 {
		errors.OutOfMemoryError()
	}

	region.Addr = table.findNextEmptyAddr()
	region.Size = 2
	region.Value = value
	region.Name = name

	*table = append(*table, region)
	return region
}

func (table MemTable) findNextEmptyAddr() int {
	//Note: all ints are padded to by 2 bytes long

	varAddrSpaceStart := 0xEA0
	varAddrSpaceEnd := 0xEFF

	currentSpaceUsed := 0

	for i := 0; i < len(table); i++ {
		currentSpaceUsed += 0x2
	}

	if varAddrSpaceStart+currentSpaceUsed+0x2 > varAddrSpaceEnd {
		errors.OutOfMemoryError()
	}

	return varAddrSpaceStart + currentSpaceUsed
}

func (table MemTable) getSize() int {
	totalSize := 0
	for i := 0; i < len(table); i++ {
		totalSize += table[i].Size
	}

	return totalSize
}
