#ifndef WREN_CALLBACKS_H
#define WREN_CALLBACKS_H

#include "wren.h"

void wrengoWriteFn(WrenVM* vm, const char* text);
void wrengoErrorFn(WrenVM* vm, WrenErrorType type, const char* module, int line, const char* message);

// Foreign function callbacks
void wrengoForeignMethodCallback(WrenVM* vm);
void wrengoForeignAllocateCallback(WrenVM* vm);
void wrengoForeignFinalizeCallback(void* data);

#endif
