package wrengo

/*
#include <stdlib.h>
#include "wren.h"
#include "wren_callbacks.h"
*/
import "C"

// moduleDefinitions maps module names to their Wren source code
var moduleDefinitions = map[string]string{
	"async": `foreign class Async {
  foreign static sleep(ms)
  foreign static delay(ms) 
  foreign static timer(ms, message)
}`,
	"math": `foreign class Math {
  foreign static sqrt(x)
  foreign static pow(x, y)
  foreign static sin(x)
  foreign static cos(x)
  foreign static abs(x)
  foreign static max(a, b)
  foreign static min(a, b)
  foreign static pi
}`,
	"strings": `foreign class StringUtils {
  foreign static upper(str)
  foreign static lower(str)
  foreign static trim(str)
  foreign static contains(str, substr)
  foreign static split(str, delimiter)
  foreign static join(elements, delimiter)
  foreign static parseInt(str)
  foreign static parseFloat(str)
}`,
}

//export wrengoLoadModule
func wrengoLoadModule(vm *C.WrenVM, name *C.char) C.WrenLoadModuleResult {
	moduleName := C.GoString(name)
	
	var result C.WrenLoadModuleResult
	
	// Check if module is registered in our definitions
	if source, exists := moduleDefinitions[moduleName]; exists {
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