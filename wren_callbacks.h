#ifndef WREN_CALLBACKS_H
#define WREN_CALLBACKS_H

#include "wren.h"

void wrengoWriteFn(WrenVM* vm, const char* text);
void wrengoErrorFn(WrenVM* vm, WrenErrorType type, const char* module, int line, const char* message);

// Foreign function callbacks
void wrengoForeignMethodCallback(WrenVM* vm);
void wrengoForeignAllocateCallback(WrenVM* vm);
void wrengoForeignFinalizeCallback(void* data);

// Foreign method wrapper functions (for multiple method support)
// Each registered method gets assigned one of these wrappers
void wrengoForeignMethod_0(WrenVM* vm);
void wrengoForeignMethod_1(WrenVM* vm);
void wrengoForeignMethod_2(WrenVM* vm);
void wrengoForeignMethod_3(WrenVM* vm);
void wrengoForeignMethod_4(WrenVM* vm);
void wrengoForeignMethod_5(WrenVM* vm);
void wrengoForeignMethod_6(WrenVM* vm);
void wrengoForeignMethod_7(WrenVM* vm);
void wrengoForeignMethod_8(WrenVM* vm);
void wrengoForeignMethod_9(WrenVM* vm);
void wrengoForeignMethod_10(WrenVM* vm);
void wrengoForeignMethod_11(WrenVM* vm);
void wrengoForeignMethod_12(WrenVM* vm);
void wrengoForeignMethod_13(WrenVM* vm);
void wrengoForeignMethod_14(WrenVM* vm);
void wrengoForeignMethod_15(WrenVM* vm);
void wrengoForeignMethod_16(WrenVM* vm);
void wrengoForeignMethod_17(WrenVM* vm);
void wrengoForeignMethod_18(WrenVM* vm);
void wrengoForeignMethod_19(WrenVM* vm);
void wrengoForeignMethod_20(WrenVM* vm);
void wrengoForeignMethod_21(WrenVM* vm);
void wrengoForeignMethod_22(WrenVM* vm);
void wrengoForeignMethod_23(WrenVM* vm);
void wrengoForeignMethod_24(WrenVM* vm);
void wrengoForeignMethod_25(WrenVM* vm);
void wrengoForeignMethod_26(WrenVM* vm);
void wrengoForeignMethod_27(WrenVM* vm);
void wrengoForeignMethod_28(WrenVM* vm);
void wrengoForeignMethod_29(WrenVM* vm);
void wrengoForeignMethod_30(WrenVM* vm);
void wrengoForeignMethod_31(WrenVM* vm);

#endif
