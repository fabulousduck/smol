package memtable

import (
	"os"

	"github.com/fabulousduck/smol/errors"
)

/*
MemTable is a simple collection of memory regions in use
*/
type MemTable map[string]*MemRegion

/*
MemRegion represents a region of memory in the IR
*/
type MemRegion struct {
	Addr, Size, Value int
}

/*
Put places a variable on a memoryTable

notes

chip-8's blocks are 8 bit, so 1 byte.
with a total of 4096 bytes
*/
func (table MemTable) Put(name string, value int, size int) *MemRegion {
	region := new(MemRegion)
	//check if there is any memory left for our variable
	currentMemSize := table.getSize()
	if currentMemSize >= 95 {
		errors.OutOfMemoryError()
	}

	region.Addr = table.FindNextEmptyAddr()
	region.Size = 1
	region.Value = value
	table[name] = region

	return region
}

/*
IsValidMemRegion allows the caller to check if an address given
is a memory region or not
*/
func IsValidMemRegion(regionAddr int) bool {
	return regionAddr > 0xEA0-0x200 && regionAddr < 0xEFF-0x200
}

/*
LookupVariable looks up if a variable has been defined on the memory table
internal lookups are silent since sometimes we need to do a lookup if an
internal variable that is not user defined has been set
*/
func (table *MemTable) LookupVariable(name string, internalLookup bool) *MemRegion {
	if val, ok := (*table)[name]; ok {
		return val
	}
	if !internalLookup {
		errors.UndefinedVariableError(name)
	}
	return nil
}

func (table MemTable) FindNextEmptyAddr() int {

	varAddrSpaceStart := 0xEA0 - 0x200
	varAddrSpaceEnd := 0xEFF - 0x200

	currentSpaceUsed := 0

	for i := 0; i < len(table); i++ {
		currentSpaceUsed++
	}
	if varAddrSpaceStart+currentSpaceUsed+0x2 > varAddrSpaceEnd {
		errors.OutOfMemoryError()
	}

	return varAddrSpaceStart + currentSpaceUsed
}

func (table MemTable) getSize() int {
	return len(table)
}

/*
Move moves a variable on the memory table.
Mostly used for compression of variable space
*/
func (table *MemTable) Move(name string, to int, internal bool) {
	if val, ok := (*table)[name]; ok {
		val.Addr = to
		return
	}

	if !internal {
		errors.UndefinedVariableError(name)
	}
}

/*
FindByAddr finds a given memory region by its addres instead of name
*/
func (table *MemTable) FindByAddr(addr int) string {
	for k, v := range *table {
		if v.Addr == addr {
			return k
		}
	}

	errors.UnAssignedMemoryLookupError()
	os.Exit(65)
	return ""
}
