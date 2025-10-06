package wrengo

// #cgo CFLAGS: -I${SRCDIR}/deps/wren/src/include -I${SRCDIR}/deps/wren/src/vm -I${SRCDIR}/deps/wren/src/optional
// #cgo LDFLAGS: -lm
// #include <stdlib.h>
// #include <string.h>
// #include "wren.h"
// #include "wren_vm.c"
// #include "wren_compiler.c"
// #include "wren_core.c"
// #include "wren_debug.c"
// #include "wren_primitive.c"
// #include "wren_utils.c"
// #include "wren_value.c"
// #include "wren_opt_meta.c"
// #include "wren_opt_random.c"
//
// #include <stdio.h>
//
// static void wrenWriteStdout(WrenVM* vm, const char* text) {
//   printf("%s", text);
// }
//
// static void wrenErrorStderr(WrenVM* vm, WrenErrorType type, const char* module, int line, const char* message) {
//   switch (type) {
//     case WREN_ERROR_COMPILE:
//       fprintf(stderr, "[%s line %d] [Error] %s", module, line, message);
//       break;
//     case WREN_ERROR_STACK_TRACE:
//       fprintf(stderr, "[%s line %d] in %s", module, line, message);
//       break;
//     case WREN_ERROR_RUNTIME:
//       fprintf(stderr, "[Runtime Error] %s", message);
//       break;
//   }
// }
import "C"
import (
	"errors"
	"runtime"
	"unsafe"
)

// InterpretResult represents the result of interpreting Wren code.
type InterpretResult int

const (
	ResultSuccess InterpretResult = iota
	ResultCompileError
	ResultRuntimeError
)

// WrenVM represents a Wren virtual machine instance.
type WrenVM struct {
	vm *C.WrenVM
}

// NewVM creates a new Wren virtual machine with default configuration.
func NewVM() *WrenVM {
	var config C.WrenConfiguration
	C.wrenInitConfiguration(&config)

	// Set default write and error functions
	config.writeFn = C.WrenWriteFn(C.wrenWriteStdout)
	config.errorFn = C.WrenErrorFn(C.wrenErrorStderr)

	vm := &WrenVM{
		vm: C.wrenNewVM(&config),
	}

	runtime.SetFinalizer(vm, (*WrenVM).Free)
	return vm
}

// NewVMWithConfig creates a new Wren virtual machine with custom configuration.
func NewVMWithConfig(config *Configuration) *WrenVM {
	cConfig := config.toCConfig()

	vm := &WrenVM{
		vm: C.wrenNewVM(&cConfig),
	}

	runtime.SetFinalizer(vm, (*WrenVM).Free)
	return vm
}

// Free disposes of all resources used by the VM.
func (vm *WrenVM) Free() {
	if vm.vm != nil {
		C.wrenFreeVM(vm.vm)
		vm.vm = nil
	}
}

// CollectGarbage immediately runs the garbage collector to free unused memory.
func (vm *WrenVM) CollectGarbage() {
	C.wrenCollectGarbage(vm.vm)
}

// Interpret runs Wren source code in the context of the specified module.
func (vm *WrenVM) Interpret(module, source string) (InterpretResult, error) {
	if vm.vm == nil {
		return ResultRuntimeError, errors.New("VM is not initialized")
	}

	cModule := C.CString(module)
	defer C.free(unsafe.Pointer(cModule))

	cSource := C.CString(source)
	defer C.free(unsafe.Pointer(cSource))

	result := C.wrenInterpret(vm.vm, cModule, cSource)

	return InterpretResult(result), nil
}

// Configuration holds the configuration options for a Wren VM.
type Configuration struct {
	InitialHeapSize   uint64
	MinHeapSize       uint64
	HeapGrowthPercent int
}

// DefaultConfiguration returns a Configuration with default values.
func DefaultConfiguration() *Configuration {
	return &Configuration{
		InitialHeapSize:   10 * 1024 * 1024, // 10MB
		MinHeapSize:       1 * 1024 * 1024,  // 1MB
		HeapGrowthPercent: 50,
	}
}

func (config *Configuration) toCConfig() C.WrenConfiguration {
	var cConfig C.WrenConfiguration
	C.wrenInitConfiguration(&cConfig)

	if config != nil {
		if config.InitialHeapSize > 0 {
			cConfig.initialHeapSize = C.size_t(config.InitialHeapSize)
		}
		if config.MinHeapSize > 0 {
			cConfig.minHeapSize = C.size_t(config.MinHeapSize)
		}
		if config.HeapGrowthPercent > 0 {
			cConfig.heapGrowthPercent = C.int(config.HeapGrowthPercent)
		}
	}

	return cConfig
}

// GetVersionNumber returns the Wren version number.
func GetVersionNumber() int {
	return int(C.wrenGetVersionNumber())
}
