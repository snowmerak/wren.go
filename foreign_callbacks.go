package wrengo

// #include <stdlib.h>
// #include "wren.h"
// #include "wren_callbacks.h"
import "C"
import (
	"sync"
	"unsafe"
)

const maxForeignMethods = 99 // Maximum number of foreign methods supported

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
	
	// Also clean up foreign data for this VM
	foreignDataMutex.Lock()
	defer foreignDataMutex.Unlock()
	delete(vmForeignDataStore, vm.vm)
}

// getVM retrieves the Go WrenVM instance from a C VM pointer
func getVM(cvm *C.WrenVM) *WrenVM {
	vmMutex.RLock()
	defer vmMutex.RUnlock()
	return vmRegistry[cvm]
}

// vmForeignData holds foreign function data for a single VM
type vmForeignData struct {
	wrapperMethods   map[int]ForeignMethodFn       // wrapperID → function
	nextWrapperID    int                            // Next available wrapper ID (0-98)
	allocators       map[string]*ForeignClass       // "module:className" → ForeignClass
}

// Foreign function callback registry - per VM
var (
	foreignDataMutex   sync.RWMutex
	vmForeignDataStore = make(map[*C.WrenVM]*vmForeignData)
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

	foreignDataMutex.Lock()
	defer foreignDataMutex.Unlock()

	// Get or create foreign data for this VM
	data := vmForeignDataStore[cvm]
	if data == nil {
		data = &vmForeignData{
			wrapperMethods: make(map[int]ForeignMethodFn),
			nextWrapperID:  0,
			allocators:     make(map[string]*ForeignClass),
		}
		vmForeignDataStore[cvm] = data
	}

	// Check if we've exceeded the wrapper limit
	if data.nextWrapperID >= maxForeignMethods {
		panic("Maximum foreign methods (99) exceeded for this VM")
	}

	// Assign the next available wrapper ID
	wrapperID := data.nextWrapperID
	data.wrapperMethods[wrapperID] = fn
	data.nextWrapperID++

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
	case 32: return C.WrenForeignMethodFn(C.wrengoForeignMethod_32)
	case 33: return C.WrenForeignMethodFn(C.wrengoForeignMethod_33)
	case 34: return C.WrenForeignMethodFn(C.wrengoForeignMethod_34)
	case 35: return C.WrenForeignMethodFn(C.wrengoForeignMethod_35)
	case 36: return C.WrenForeignMethodFn(C.wrengoForeignMethod_36)
	case 37: return C.WrenForeignMethodFn(C.wrengoForeignMethod_37)
	case 38: return C.WrenForeignMethodFn(C.wrengoForeignMethod_38)
	case 39: return C.WrenForeignMethodFn(C.wrengoForeignMethod_39)
	case 40: return C.WrenForeignMethodFn(C.wrengoForeignMethod_40)
	case 41: return C.WrenForeignMethodFn(C.wrengoForeignMethod_41)
	case 42: return C.WrenForeignMethodFn(C.wrengoForeignMethod_42)
	case 43: return C.WrenForeignMethodFn(C.wrengoForeignMethod_43)
	case 44: return C.WrenForeignMethodFn(C.wrengoForeignMethod_44)
	case 45: return C.WrenForeignMethodFn(C.wrengoForeignMethod_45)
	case 46: return C.WrenForeignMethodFn(C.wrengoForeignMethod_46)
	case 47: return C.WrenForeignMethodFn(C.wrengoForeignMethod_47)
	case 48: return C.WrenForeignMethodFn(C.wrengoForeignMethod_48)
	case 49: return C.WrenForeignMethodFn(C.wrengoForeignMethod_49)
	case 50: return C.WrenForeignMethodFn(C.wrengoForeignMethod_50)
	case 51: return C.WrenForeignMethodFn(C.wrengoForeignMethod_51)
	case 52: return C.WrenForeignMethodFn(C.wrengoForeignMethod_52)
	case 53: return C.WrenForeignMethodFn(C.wrengoForeignMethod_53)
	case 54: return C.WrenForeignMethodFn(C.wrengoForeignMethod_54)
	case 55: return C.WrenForeignMethodFn(C.wrengoForeignMethod_55)
	case 56: return C.WrenForeignMethodFn(C.wrengoForeignMethod_56)
	case 57: return C.WrenForeignMethodFn(C.wrengoForeignMethod_57)
	case 58: return C.WrenForeignMethodFn(C.wrengoForeignMethod_58)
	case 59: return C.WrenForeignMethodFn(C.wrengoForeignMethod_59)
	case 60: return C.WrenForeignMethodFn(C.wrengoForeignMethod_60)
	case 61: return C.WrenForeignMethodFn(C.wrengoForeignMethod_61)
	case 62: return C.WrenForeignMethodFn(C.wrengoForeignMethod_62)
	case 63: return C.WrenForeignMethodFn(C.wrengoForeignMethod_63)
	case 64: return C.WrenForeignMethodFn(C.wrengoForeignMethod_64)
	case 65: return C.WrenForeignMethodFn(C.wrengoForeignMethod_65)
	case 66: return C.WrenForeignMethodFn(C.wrengoForeignMethod_66)
	case 67: return C.WrenForeignMethodFn(C.wrengoForeignMethod_67)
	case 68: return C.WrenForeignMethodFn(C.wrengoForeignMethod_68)
	case 69: return C.WrenForeignMethodFn(C.wrengoForeignMethod_69)
	case 70: return C.WrenForeignMethodFn(C.wrengoForeignMethod_70)
	case 71: return C.WrenForeignMethodFn(C.wrengoForeignMethod_71)
	case 72: return C.WrenForeignMethodFn(C.wrengoForeignMethod_72)
	case 73: return C.WrenForeignMethodFn(C.wrengoForeignMethod_73)
	case 74: return C.WrenForeignMethodFn(C.wrengoForeignMethod_74)
	case 75: return C.WrenForeignMethodFn(C.wrengoForeignMethod_75)
	case 76: return C.WrenForeignMethodFn(C.wrengoForeignMethod_76)
	case 77: return C.WrenForeignMethodFn(C.wrengoForeignMethod_77)
	case 78: return C.WrenForeignMethodFn(C.wrengoForeignMethod_78)
	case 79: return C.WrenForeignMethodFn(C.wrengoForeignMethod_79)
	case 80: return C.WrenForeignMethodFn(C.wrengoForeignMethod_80)
	case 81: return C.WrenForeignMethodFn(C.wrengoForeignMethod_81)
	case 82: return C.WrenForeignMethodFn(C.wrengoForeignMethod_82)
	case 83: return C.WrenForeignMethodFn(C.wrengoForeignMethod_83)
	case 84: return C.WrenForeignMethodFn(C.wrengoForeignMethod_84)
	case 85: return C.WrenForeignMethodFn(C.wrengoForeignMethod_85)
	case 86: return C.WrenForeignMethodFn(C.wrengoForeignMethod_86)
	case 87: return C.WrenForeignMethodFn(C.wrengoForeignMethod_87)
	case 88: return C.WrenForeignMethodFn(C.wrengoForeignMethod_88)
	case 89: return C.WrenForeignMethodFn(C.wrengoForeignMethod_89)
	case 90: return C.WrenForeignMethodFn(C.wrengoForeignMethod_90)
	case 91: return C.WrenForeignMethodFn(C.wrengoForeignMethod_91)
	case 92: return C.WrenForeignMethodFn(C.wrengoForeignMethod_92)
	case 93: return C.WrenForeignMethodFn(C.wrengoForeignMethod_93)
	case 94: return C.WrenForeignMethodFn(C.wrengoForeignMethod_94)
	case 95: return C.WrenForeignMethodFn(C.wrengoForeignMethod_95)
	case 96: return C.WrenForeignMethodFn(C.wrengoForeignMethod_96)
	case 97: return C.WrenForeignMethodFn(C.wrengoForeignMethod_97)
	case 98: return C.WrenForeignMethodFn(C.wrengoForeignMethod_98)
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

	foreignDataMutex.Lock()
	defer foreignDataMutex.Unlock()

	// Get or create foreign data for this VM
	data := vmForeignDataStore[cvm]
	if data == nil {
		data = &vmForeignData{
			wrapperMethods: make(map[int]ForeignMethodFn),
			nextWrapperID:  0,
			allocators:     make(map[string]*ForeignClass),
		}
		vmForeignDataStore[cvm] = data
	}

	// Store the class in VM-specific registry
	key := module + ":" + className
	data.allocators[key] = class

	if class.Allocate != nil {
		methods.allocate = C.WrenForeignMethodFn(C.wrengoForeignAllocateCallback)
	}

	if class.Finalize != nil {
		methods.finalize = C.WrenFinalizerFn(C.wrengoForeignFinalizeCallback)
	}

	return methods
}



//export goForeignMethodCallback
func goForeignMethodCallback(cvm *C.WrenVM, wrapperId C.int) {
	vm := getVM(cvm)
	if vm == nil {
		return
	}

	foreignDataMutex.RLock()
	data := vmForeignDataStore[cvm]
	foreignDataMutex.RUnlock()

	if data == nil {
		return
	}

	// Look up the function for this wrapper ID
	if fn, ok := data.wrapperMethods[int(wrapperId)]; ok {
		fn(vm)
	}
}

//export wrengoForeignAllocateCallback
func wrengoForeignAllocateCallback(cvm *C.WrenVM) {
	vm := getVM(cvm)
	if vm == nil {
		return
	}

	foreignDataMutex.RLock()
	data := vmForeignDataStore[cvm]
	foreignDataMutex.RUnlock()

	if data == nil {
		return
	}

	// Get the class being allocated from slot 0
	// Note: Wren passes the class in slot 0 during allocation
	// We need to find which class is being allocated
	// For now, call the first available allocator
	// This is a limitation - ideally we'd track which class is being allocated
	for _, class := range data.allocators {
		if class.Allocate != nil {
			class.Allocate(vm)
			return
		}
	}
}

//export wrengoForeignFinalizeCallback
func wrengoForeignFinalizeCallback(data unsafe.Pointer) {
	// Note: We don't have VM context here, so we can't call VM-specific finalizers
	// This is a limitation of the current design
	// Finalizers would need to be stored differently to work properly
}
