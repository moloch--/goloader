// +build go1.12
// +build !go1.18

package goloader

import (
	"fmt"
	"unsafe"
)

func addr2stackObjectRecords(addr unsafe.Pointer) *[]stackObjectRecord {
	n := int(*(*uintptr)(addr))
	slice := sliceHeader{
		Data: uintptr(add(addr, uintptr(PtrSize))),
		Len:  n,
		Cap:  n,
	}
	return (*[]stackObjectRecord)(unsafe.Pointer(&slice))
}

func (linker *Linker) _addStackObject(funcname string, symbolMap map[string]uintptr) (err error) {
	Func := linker.symMap[funcname].Func
	if Func != nil && len(Func.FuncData) > _FUNCDATA_StackObjects &&
		Func.FuncData[_FUNCDATA_StackObjects] != 0 {
		objects := addr2stackObjectRecords(adduintptr(Func.FuncData[_FUNCDATA_StackObjects], 0))
		for i := range *objects {
			name := EmptyString
			stkobjName := funcname + StkobjSuffix
			if symbol := linker.symMap[stkobjName]; symbol != nil {
				name = symbol.Reloc[i].Sym.Name
			}
			if ptr, ok := symbolMap[name]; ok {
				setStackObjectPtr(&((*objects)[i]), adduintptr(ptr, 0))
			} else {
				return fmt.Errorf("unresolve external Var! Function name:%s index:%d, name:%s", funcname, i, name)

			}
		}
	}
	return nil
}
