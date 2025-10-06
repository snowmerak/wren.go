#ifndef WREN_CALLBACKS_H
#define WREN_CALLBACKS_H

#include "wren.h"

void wrengoWriteFn(WrenVM* vm, const char* text);
void wrengoErrorFn(WrenVM* vm, WrenErrorType type, const char* module, int line, const char* message);

#endif
