// +build go1.8
// +build !go1.13

package goloader

import (
	"unsafe"
)

const (
	_MainInitTask   = "main.init"
	_InitTaskSuffix = ".init"
)

func getInitFuncName(packagename string) string {
	return packagename + _InitTaskSuffix
}

func isMainInitFunc(name string) bool {
	if name == _MainInitTask {
		return true
	}
	return false
}

func (linker *Linker) doInitialize(codeModule *CodeModule, symbolMap map[string]uintptr) error {
	for _, name := range linker.initFuncs {
		if funcPtr, ok := codeModule.Syms[name]; ok {
			funcPtrContainer := (uintptr)(unsafe.Pointer(&funcPtr))
			runFunc := *(*func())(unsafe.Pointer(&funcPtrContainer))
			runFunc()
		}
	}
	return nil
}