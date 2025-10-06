package wrengo

// #include <stdlib.h>
// #include "wren.h"
// #include "wren_callbacks.h"
import "C"
import "unsafe"

// vmRegistry stores VM instances for callback access
var vmRegistry = make(map[*C.WrenVM]*WrenVM)

// registerVM associates a C VM pointer with a Go WrenVM instance
func registerVM(vm *WrenVM) {
	vmRegistry[vm.vm] = vm
}

// unregisterVM removes the association
func unregisterVM(vm *WrenVM) {
	delete(vmRegistry, vm.vm)
}

// getVM retrieves the Go WrenVM instance from a C VM pointer
func getVM(cvm *C.WrenVM) *WrenVM {
	return vmRegistry[cvm]
}

//export wrengoBindForeignMethod
func wrengoBindForeignMethod(cvm *C.WrenVM, cModule, cClassName *C.char, isStatic C.bool, cSignature *C.char) C.WrenForeignMethodFn {
	module := C.GoString(cModule)
	className := C.GoString(cClassName)
	signature := C.GoString(cSignature)

	fn := lookupForeignMethod(module, className, bool(isStatic), signature)
	if fn == nil {
		return nil
	}

	// Store the function in a registry with a unique ID
	storeForeignMethod(fn)

	// Return a C function pointer that will call our Go function
	return C.WrenForeignMethodFn(C.wrengoForeignMethodCallback)
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

// Foreign function callback registry
var (
	foreignMethodID      uint64
	foreignMethods       = make(map[uint64]ForeignMethodFn)
	foreignAllocatorID   uint64
	foreignAllocators    = make(map[uint64]ForeignClassAllocator)
	foreignFinalizerID   uint64
	foreignFinalizers    = make(map[uint64]ForeignClassFinalizer)
	currentMethodID      uint64
	currentAllocatorID   uint64
	currentFinalizerID   uint64
)

func storeForeignMethod(fn ForeignMethodFn) uint64 {
	foreignMethodID++
	id := foreignMethodID
	foreignMethods[id] = fn
	currentMethodID = id
	return id
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

//export wrengoForeignMethodCallback
func wrengoForeignMethodCallback(cvm *C.WrenVM) {
	vm := getVM(cvm)
	if vm == nil {
		return
	}

	if fn, ok := foreignMethods[currentMethodID]; ok {
		fn(vm)
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
