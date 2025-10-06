package wrengo

// #include <stdlib.h>
// #include "wren.h"
// #include "wren_callbacks.h"
import "C"
import (
	"sync"
	"unsafe"
)

const maxForeignMethods = 32 // Maximum number of foreign methods supported

// vmRegistry stores VM instances for callback access
var vmRegistry = make(map[*C.WrenVM]*WrenVM)
var vmMutex sync.RWMutex

// registerVM associates a C VM pointer with a Go WrenVM instance
func registerVM(vm *WrenVM) {
	vmMutex.Lock()
	defer vmMutex.Unlock()
	vmRegistry[vm.vm] = vm
}

// unregisterVM removes the association
func unregisterVM(vm *WrenVM) {
	vmMutex.Lock()
	defer vmMutex.Unlock()
	delete(vmRegistry, vm.vm)
}

// getVM retrieves the Go WrenVM instance from a C VM pointer
func getVM(cvm *C.WrenVM) *WrenVM {
	vmMutex.RLock()
	defer vmMutex.RUnlock()
	return vmRegistry[cvm]
}

// Foreign function callback registry
var (
	foreignMethodsMutex sync.RWMutex
	foreignMethodID     uint64
	foreignMethods      = make(map[uint64]ForeignMethodFn)
	foreignWrapperID    = make(map[uint64]int) // Maps method ID to wrapper ID (0-31)
	
	foreignAllocatorID uint64
	foreignAllocators  = make(map[uint64]ForeignClassAllocator)
	foreignFinalizerID uint64
	foreignFinalizers  = make(map[uint64]ForeignClassFinalizer)
	currentAllocatorID uint64
	currentFinalizerID uint64
)

//export wrengoBindForeignMethod
func wrengoBindForeignMethod(cvm *C.WrenVM, cModule, cClassName *C.char, isStatic C.bool, cSignature *C.char) C.WrenForeignMethodFn {
	module := C.GoString(cModule)
	className := C.GoString(cClassName)
	signature := C.GoString(cSignature)

	fn := lookupForeignMethod(module, className, bool(isStatic), signature)
	if fn == nil {
		return nil
	}

	foreignMethodsMutex.Lock()
	defer foreignMethodsMutex.Unlock()

	// Store the function in a registry with a unique ID
	foreignMethodID++
	id := foreignMethodID
	foreignMethods[id] = fn
	
	// Assign a wrapper ID (0-31) by cycling through available wrappers
	wrapperID := int((id - 1) % maxForeignMethods)
	foreignWrapperID[id] = wrapperID

	// Return the appropriate C function pointer based on wrapper ID
	switch wrapperID {
	case 0: return C.WrenForeignMethodFn(C.wrengoForeignMethod_0)
	case 1: return C.WrenForeignMethodFn(C.wrengoForeignMethod_1)
	case 2: return C.WrenForeignMethodFn(C.wrengoForeignMethod_2)
	case 3: return C.WrenForeignMethodFn(C.wrengoForeignMethod_3)
	case 4: return C.WrenForeignMethodFn(C.wrengoForeignMethod_4)
	case 5: return C.WrenForeignMethodFn(C.wrengoForeignMethod_5)
	case 6: return C.WrenForeignMethodFn(C.wrengoForeignMethod_6)
	case 7: return C.WrenForeignMethodFn(C.wrengoForeignMethod_7)
	case 8: return C.WrenForeignMethodFn(C.wrengoForeignMethod_8)
	case 9: return C.WrenForeignMethodFn(C.wrengoForeignMethod_9)
	case 10: return C.WrenForeignMethodFn(C.wrengoForeignMethod_10)
	case 11: return C.WrenForeignMethodFn(C.wrengoForeignMethod_11)
	case 12: return C.WrenForeignMethodFn(C.wrengoForeignMethod_12)
	case 13: return C.WrenForeignMethodFn(C.wrengoForeignMethod_13)
	case 14: return C.WrenForeignMethodFn(C.wrengoForeignMethod_14)
	case 15: return C.WrenForeignMethodFn(C.wrengoForeignMethod_15)
	case 16: return C.WrenForeignMethodFn(C.wrengoForeignMethod_16)
	case 17: return C.WrenForeignMethodFn(C.wrengoForeignMethod_17)
	case 18: return C.WrenForeignMethodFn(C.wrengoForeignMethod_18)
	case 19: return C.WrenForeignMethodFn(C.wrengoForeignMethod_19)
	case 20: return C.WrenForeignMethodFn(C.wrengoForeignMethod_20)
	case 21: return C.WrenForeignMethodFn(C.wrengoForeignMethod_21)
	case 22: return C.WrenForeignMethodFn(C.wrengoForeignMethod_22)
	case 23: return C.WrenForeignMethodFn(C.wrengoForeignMethod_23)
	case 24: return C.WrenForeignMethodFn(C.wrengoForeignMethod_24)
	case 25: return C.WrenForeignMethodFn(C.wrengoForeignMethod_25)
	case 26: return C.WrenForeignMethodFn(C.wrengoForeignMethod_26)
	case 27: return C.WrenForeignMethodFn(C.wrengoForeignMethod_27)
	case 28: return C.WrenForeignMethodFn(C.wrengoForeignMethod_28)
	case 29: return C.WrenForeignMethodFn(C.wrengoForeignMethod_29)
	case 30: return C.WrenForeignMethodFn(C.wrengoForeignMethod_30)
	case 31: return C.WrenForeignMethodFn(C.wrengoForeignMethod_31)
	}

	return nil
}

//export wrengoBindForeignClass
func wrengoBindForeignClass(cvm *C.WrenVM, cModule, cClassName *C.char) C.WrenForeignClassMethods {
	var methods C.WrenForeignClassMethods

	module := C.GoString(cModule)
	className := C.GoString(cClassName)

	class := lookupForeignClass(module, className)
	if class == nil {
		return methods
	}

	if class.Allocate != nil {
		storeForeignAllocator(class.Allocate)
		methods.allocate = C.WrenForeignMethodFn(C.wrengoForeignAllocateCallback)
	}

	if class.Finalize != nil {
		storeForeignFinalizer(class.Finalize)
		methods.finalize = C.WrenFinalizerFn(C.wrengoForeignFinalizeCallback)
	}

	return methods
}

func storeForeignAllocator(fn ForeignClassAllocator) uint64 {
	foreignAllocatorID++
	id := foreignAllocatorID
	foreignAllocators[id] = fn
	currentAllocatorID = id
	return id
}

func storeForeignFinalizer(fn ForeignClassFinalizer) uint64 {
	foreignFinalizerID++
	id := foreignFinalizerID
	foreignFinalizers[id] = fn
	currentFinalizerID = id
	return id
}

//export goForeignMethodCallback
func goForeignMethodCallback(cvm *C.WrenVM, wrapperId C.int) {
	vm := getVM(cvm)
	if vm == nil {
		return
	}

	foreignMethodsMutex.RLock()
	defer foreignMethodsMutex.RUnlock()

	// Find the method ID that matches this wrapper ID
	wrapperIDInt := int(wrapperId)
	for methodID, wid := range foreignWrapperID {
		if wid == wrapperIDInt {
			if fn, ok := foreignMethods[methodID]; ok {
				fn(vm)
				return
			}
		}
	}
}

//export wrengoForeignAllocateCallback
func wrengoForeignAllocateCallback(cvm *C.WrenVM) {
	vm := getVM(cvm)
	if vm == nil {
		return
	}

	if fn, ok := foreignAllocators[currentAllocatorID]; ok {
		fn(vm)
	}
}

//export wrengoForeignFinalizeCallback
func wrengoForeignFinalizeCallback(data unsafe.Pointer) {
	if fn, ok := foreignFinalizers[currentFinalizerID]; ok {
		fn(data)
	}
}
