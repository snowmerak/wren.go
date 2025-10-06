package wrengo

// #include <stdlib.h>
// #include "wren.h"
// #include "wren_callbacks.h"
//
// // Forward declarations for foreign function callbacks
// extern WrenForeignMethodFn wrengoBindForeignMethod(WrenVM* vm, const char* module, const char* className, bool isStatic, const char* signature);
// extern WrenForeignClassMethods wrengoBindForeignClass(WrenVM* vm, const char* module, const char* className);
import "C"
import (
	"sync"
	"unsafe"
)

// ForeignMethodFn is a Go function that can be called from Wren.
// It receives the VM and should use slot methods to interact with parameters and return values.
type ForeignMethodFn func(vm *WrenVM)

// ForeignClassAllocator is called when a foreign class is instantiated.
type ForeignClassAllocator func(vm *WrenVM)

// ForeignClassFinalizer is called when a foreign class instance is garbage collected.
type ForeignClassFinalizer func(data unsafe.Pointer)

// ForeignClass holds the allocator and optional finalizer for a foreign class.
type ForeignClass struct {
	Allocate ForeignClassAllocator
	Finalize ForeignClassFinalizer
}

// foreignRegistry holds all registered foreign methods and classes
type foreignRegistry struct {
	mu      sync.RWMutex
	methods map[string]map[string]map[string]ForeignMethodFn // module -> class -> signature -> func
	classes map[string]map[string]*ForeignClass               // module -> class -> ForeignClass
}

var registry = &foreignRegistry{
	methods: make(map[string]map[string]map[string]ForeignMethodFn),
	classes: make(map[string]map[string]*ForeignClass),
}

// RegisterForeignMethod registers a Go function as a foreign method for Wren.
// The signature should match Wren's method signature format, e.g., "add(_,_)" for a static method with 2 parameters.
func RegisterForeignMethod(module, className string, isStatic bool, signature string, fn ForeignMethodFn) {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	if registry.methods[module] == nil {
		registry.methods[module] = make(map[string]map[string]ForeignMethodFn)
	}
	if registry.methods[module][className] == nil {
		registry.methods[module][className] = make(map[string]ForeignMethodFn)
	}

	// Create full signature with static prefix if needed
	fullSig := signature
	if isStatic {
		fullSig = "static " + signature
	}

	registry.methods[module][className][fullSig] = fn
}

// RegisterForeignClass registers allocator and optional finalizer for a foreign class.
func RegisterForeignClass(module, className string, allocate ForeignClassAllocator, finalize ForeignClassFinalizer) {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	if registry.classes[module] == nil {
		registry.classes[module] = make(map[string]*ForeignClass)
	}

	registry.classes[module][className] = &ForeignClass{
		Allocate: allocate,
		Finalize: finalize,
	}
}

// lookupForeignMethod finds a registered foreign method
func lookupForeignMethod(module, className string, isStatic bool, signature string) ForeignMethodFn {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	fullSig := signature
	if isStatic {
		fullSig = "static " + signature
	}

	if modMethods, ok := registry.methods[module]; ok {
		if classMethods, ok := modMethods[className]; ok {
			if fn, ok := classMethods[fullSig]; ok {
				return fn
			}
		}
	}

	return nil
}

// lookupForeignClass finds a registered foreign class
func lookupForeignClass(module, className string) *ForeignClass {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	if modClasses, ok := registry.classes[module]; ok {
		if class, ok := modClasses[className]; ok {
			return class
		}
	}

	return nil
}

// NewVMWithForeign creates a new VM with foreign method and class support.
func NewVMWithForeign() *WrenVM {
	var config C.WrenConfiguration
	C.wrenInitConfiguration(&config)

	// Set default write and error callbacks
	config.writeFn = C.WrenWriteFn(C.wrengoWriteFn)
	config.errorFn = C.WrenErrorFn(C.wrengoErrorFn)

	// Set foreign method and class binders
	config.bindForeignMethodFn = C.WrenBindForeignMethodFn(C.wrengoBindForeignMethod)
	config.bindForeignClassFn = C.WrenBindForeignClassFn(C.wrengoBindForeignClass)

	vm := &WrenVM{
		vm: C.wrenNewVM(&config),
	}

	registerVM(vm)
	return vm
}

// SetSlotNewForeign creates a new instance of a foreign class.
// It allocates size bytes of memory and returns a pointer to it.
func (vm *WrenVM) SetSlotNewForeign(slot, classSlot int, size int) unsafe.Pointer {
	return C.wrenSetSlotNewForeign(vm.vm, C.int(slot), C.int(classSlot), C.size_t(size))
}
