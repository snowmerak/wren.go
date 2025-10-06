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
void wrengoForeignMethod_32(WrenVM* vm) { goForeignMethodCallback(vm, 32); }
void wrengoForeignMethod_33(WrenVM* vm) { goForeignMethodCallback(vm, 33); }
void wrengoForeignMethod_34(WrenVM* vm) { goForeignMethodCallback(vm, 34); }
void wrengoForeignMethod_35(WrenVM* vm) { goForeignMethodCallback(vm, 35); }
void wrengoForeignMethod_36(WrenVM* vm) { goForeignMethodCallback(vm, 36); }
void wrengoForeignMethod_37(WrenVM* vm) { goForeignMethodCallback(vm, 37); }
void wrengoForeignMethod_38(WrenVM* vm) { goForeignMethodCallback(vm, 38); }
void wrengoForeignMethod_39(WrenVM* vm) { goForeignMethodCallback(vm, 39); }
void wrengoForeignMethod_40(WrenVM* vm) { goForeignMethodCallback(vm, 40); }
void wrengoForeignMethod_41(WrenVM* vm) { goForeignMethodCallback(vm, 41); }
void wrengoForeignMethod_42(WrenVM* vm) { goForeignMethodCallback(vm, 42); }
void wrengoForeignMethod_43(WrenVM* vm) { goForeignMethodCallback(vm, 43); }
void wrengoForeignMethod_44(WrenVM* vm) { goForeignMethodCallback(vm, 44); }
void wrengoForeignMethod_45(WrenVM* vm) { goForeignMethodCallback(vm, 45); }
void wrengoForeignMethod_46(WrenVM* vm) { goForeignMethodCallback(vm, 46); }
void wrengoForeignMethod_47(WrenVM* vm) { goForeignMethodCallback(vm, 47); }
void wrengoForeignMethod_48(WrenVM* vm) { goForeignMethodCallback(vm, 48); }
void wrengoForeignMethod_49(WrenVM* vm) { goForeignMethodCallback(vm, 49); }
void wrengoForeignMethod_50(WrenVM* vm) { goForeignMethodCallback(vm, 50); }
void wrengoForeignMethod_51(WrenVM* vm) { goForeignMethodCallback(vm, 51); }
void wrengoForeignMethod_52(WrenVM* vm) { goForeignMethodCallback(vm, 52); }
void wrengoForeignMethod_53(WrenVM* vm) { goForeignMethodCallback(vm, 53); }
void wrengoForeignMethod_54(WrenVM* vm) { goForeignMethodCallback(vm, 54); }
void wrengoForeignMethod_55(WrenVM* vm) { goForeignMethodCallback(vm, 55); }
void wrengoForeignMethod_56(WrenVM* vm) { goForeignMethodCallback(vm, 56); }
void wrengoForeignMethod_57(WrenVM* vm) { goForeignMethodCallback(vm, 57); }
void wrengoForeignMethod_58(WrenVM* vm) { goForeignMethodCallback(vm, 58); }
void wrengoForeignMethod_59(WrenVM* vm) { goForeignMethodCallback(vm, 59); }
void wrengoForeignMethod_60(WrenVM* vm) { goForeignMethodCallback(vm, 60); }
void wrengoForeignMethod_61(WrenVM* vm) { goForeignMethodCallback(vm, 61); }
void wrengoForeignMethod_62(WrenVM* vm) { goForeignMethodCallback(vm, 62); }
void wrengoForeignMethod_63(WrenVM* vm) { goForeignMethodCallback(vm, 63); }
void wrengoForeignMethod_64(WrenVM* vm) { goForeignMethodCallback(vm, 64); }
void wrengoForeignMethod_65(WrenVM* vm) { goForeignMethodCallback(vm, 65); }
void wrengoForeignMethod_66(WrenVM* vm) { goForeignMethodCallback(vm, 66); }
void wrengoForeignMethod_67(WrenVM* vm) { goForeignMethodCallback(vm, 67); }
void wrengoForeignMethod_68(WrenVM* vm) { goForeignMethodCallback(vm, 68); }
void wrengoForeignMethod_69(WrenVM* vm) { goForeignMethodCallback(vm, 69); }
void wrengoForeignMethod_70(WrenVM* vm) { goForeignMethodCallback(vm, 70); }
void wrengoForeignMethod_71(WrenVM* vm) { goForeignMethodCallback(vm, 71); }
void wrengoForeignMethod_72(WrenVM* vm) { goForeignMethodCallback(vm, 72); }
void wrengoForeignMethod_73(WrenVM* vm) { goForeignMethodCallback(vm, 73); }
void wrengoForeignMethod_74(WrenVM* vm) { goForeignMethodCallback(vm, 74); }
void wrengoForeignMethod_75(WrenVM* vm) { goForeignMethodCallback(vm, 75); }
void wrengoForeignMethod_76(WrenVM* vm) { goForeignMethodCallback(vm, 76); }
void wrengoForeignMethod_77(WrenVM* vm) { goForeignMethodCallback(vm, 77); }
void wrengoForeignMethod_78(WrenVM* vm) { goForeignMethodCallback(vm, 78); }
void wrengoForeignMethod_79(WrenVM* vm) { goForeignMethodCallback(vm, 79); }
void wrengoForeignMethod_80(WrenVM* vm) { goForeignMethodCallback(vm, 80); }
void wrengoForeignMethod_81(WrenVM* vm) { goForeignMethodCallback(vm, 81); }
void wrengoForeignMethod_82(WrenVM* vm) { goForeignMethodCallback(vm, 82); }
void wrengoForeignMethod_83(WrenVM* vm) { goForeignMethodCallback(vm, 83); }
void wrengoForeignMethod_84(WrenVM* vm) { goForeignMethodCallback(vm, 84); }
void wrengoForeignMethod_85(WrenVM* vm) { goForeignMethodCallback(vm, 85); }
void wrengoForeignMethod_86(WrenVM* vm) { goForeignMethodCallback(vm, 86); }
void wrengoForeignMethod_87(WrenVM* vm) { goForeignMethodCallback(vm, 87); }
void wrengoForeignMethod_88(WrenVM* vm) { goForeignMethodCallback(vm, 88); }
void wrengoForeignMethod_89(WrenVM* vm) { goForeignMethodCallback(vm, 89); }
void wrengoForeignMethod_90(WrenVM* vm) { goForeignMethodCallback(vm, 90); }
void wrengoForeignMethod_91(WrenVM* vm) { goForeignMethodCallback(vm, 91); }
void wrengoForeignMethod_92(WrenVM* vm) { goForeignMethodCallback(vm, 92); }
void wrengoForeignMethod_93(WrenVM* vm) { goForeignMethodCallback(vm, 93); }
void wrengoForeignMethod_94(WrenVM* vm) { goForeignMethodCallback(vm, 94); }
void wrengoForeignMethod_95(WrenVM* vm) { goForeignMethodCallback(vm, 95); }
void wrengoForeignMethod_96(WrenVM* vm) { goForeignMethodCallback(vm, 96); }
void wrengoForeignMethod_97(WrenVM* vm) { goForeignMethodCallback(vm, 97); }
void wrengoForeignMethod_98(WrenVM* vm) { goForeignMethodCallback(vm, 98); }
