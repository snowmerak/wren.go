package wrengo

/*
#include <stdlib.h>
#include "wren.h"
#include "wren_callbacks.h"
*/
import "C"

//export wrengoLoadModule
func wrengoLoadModule(vm *C.WrenVM, name *C.char) C.WrenLoadModuleResult {
	moduleName := C.GoString(name)
	
	var result C.WrenLoadModuleResult
	
	if moduleName == "async" {
		// Native module with minimal Wren class declarations
		source := `foreign class Async {
  foreign static sleep(ms)
  foreign static delay(ms) 
  foreign static timer(ms, message)
}`
		cSource := C.CString(source)
		result.source = cSource
		result.onComplete = nil
		result.userData = nil
	} else {
		// Module not found
		result.source = nil
		result.onComplete = nil
		result.userData = nil
	}
	
	return result
}