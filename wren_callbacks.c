#include <stdio.h>
#include "wren.h"

void wrengoWriteFn(WrenVM* vm, const char* text) {
    printf("%s", text);
}

void wrengoErrorFn(WrenVM* vm, WrenErrorType type, const char* module, int line, const char* message) {
    switch (type) {
        case WREN_ERROR_COMPILE:
            fprintf(stderr, "[%s line %d] [Error] %s\n", module, line, message);
            break;
        case WREN_ERROR_STACK_TRACE:
            fprintf(stderr, "[%s line %d] in %s\n", module, line, message);
            break;
        case WREN_ERROR_RUNTIME:
            fprintf(stderr, "[Runtime Error] %s\n", message);
            break;
    }
}


// Foreign method wrapper functions
// These call back to Go with their specific wrapper ID
extern void goForeignMethodCallback(WrenVM* vm, int wrapperId);

void wrengoForeignMethod_0(WrenVM* vm) { goForeignMethodCallback(vm, 0); }
void wrengoForeignMethod_1(WrenVM* vm) { goForeignMethodCallback(vm, 1); }
void wrengoForeignMethod_2(WrenVM* vm) { goForeignMethodCallback(vm, 2); }
void wrengoForeignMethod_3(WrenVM* vm) { goForeignMethodCallback(vm, 3); }
void wrengoForeignMethod_4(WrenVM* vm) { goForeignMethodCallback(vm, 4); }
void wrengoForeignMethod_5(WrenVM* vm) { goForeignMethodCallback(vm, 5); }
void wrengoForeignMethod_6(WrenVM* vm) { goForeignMethodCallback(vm, 6); }
void wrengoForeignMethod_7(WrenVM* vm) { goForeignMethodCallback(vm, 7); }
void wrengoForeignMethod_8(WrenVM* vm) { goForeignMethodCallback(vm, 8); }
void wrengoForeignMethod_9(WrenVM* vm) { goForeignMethodCallback(vm, 9); }
void wrengoForeignMethod_10(WrenVM* vm) { goForeignMethodCallback(vm, 10); }
void wrengoForeignMethod_11(WrenVM* vm) { goForeignMethodCallback(vm, 11); }
void wrengoForeignMethod_12(WrenVM* vm) { goForeignMethodCallback(vm, 12); }
void wrengoForeignMethod_13(WrenVM* vm) { goForeignMethodCallback(vm, 13); }
void wrengoForeignMethod_14(WrenVM* vm) { goForeignMethodCallback(vm, 14); }
void wrengoForeignMethod_15(WrenVM* vm) { goForeignMethodCallback(vm, 15); }
void wrengoForeignMethod_16(WrenVM* vm) { goForeignMethodCallback(vm, 16); }
void wrengoForeignMethod_17(WrenVM* vm) { goForeignMethodCallback(vm, 17); }
void wrengoForeignMethod_18(WrenVM* vm) { goForeignMethodCallback(vm, 18); }
void wrengoForeignMethod_19(WrenVM* vm) { goForeignMethodCallback(vm, 19); }
void wrengoForeignMethod_20(WrenVM* vm) { goForeignMethodCallback(vm, 20); }
void wrengoForeignMethod_21(WrenVM* vm) { goForeignMethodCallback(vm, 21); }
void wrengoForeignMethod_22(WrenVM* vm) { goForeignMethodCallback(vm, 22); }
void wrengoForeignMethod_23(WrenVM* vm) { goForeignMethodCallback(vm, 23); }
void wrengoForeignMethod_24(WrenVM* vm) { goForeignMethodCallback(vm, 24); }
void wrengoForeignMethod_25(WrenVM* vm) { goForeignMethodCallback(vm, 25); }
void wrengoForeignMethod_26(WrenVM* vm) { goForeignMethodCallback(vm, 26); }
void wrengoForeignMethod_27(WrenVM* vm) { goForeignMethodCallback(vm, 27); }
void wrengoForeignMethod_28(WrenVM* vm) { goForeignMethodCallback(vm, 28); }
void wrengoForeignMethod_29(WrenVM* vm) { goForeignMethodCallback(vm, 29); }
void wrengoForeignMethod_30(WrenVM* vm) { goForeignMethodCallback(vm, 30); }
void wrengoForeignMethod_31(WrenVM* vm) { goForeignMethodCallback(vm, 31); }
