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
