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
// Each registered method gets assigned one of these wrappers (0-98, total 99)
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
void wrengoForeignMethod_32(WrenVM* vm);
void wrengoForeignMethod_33(WrenVM* vm);
void wrengoForeignMethod_34(WrenVM* vm);
void wrengoForeignMethod_35(WrenVM* vm);
void wrengoForeignMethod_36(WrenVM* vm);
void wrengoForeignMethod_37(WrenVM* vm);
void wrengoForeignMethod_38(WrenVM* vm);
void wrengoForeignMethod_39(WrenVM* vm);
void wrengoForeignMethod_40(WrenVM* vm);
void wrengoForeignMethod_41(WrenVM* vm);
void wrengoForeignMethod_42(WrenVM* vm);
void wrengoForeignMethod_43(WrenVM* vm);
void wrengoForeignMethod_44(WrenVM* vm);
void wrengoForeignMethod_45(WrenVM* vm);
void wrengoForeignMethod_46(WrenVM* vm);
void wrengoForeignMethod_47(WrenVM* vm);
void wrengoForeignMethod_48(WrenVM* vm);
void wrengoForeignMethod_49(WrenVM* vm);
void wrengoForeignMethod_50(WrenVM* vm);
void wrengoForeignMethod_51(WrenVM* vm);
void wrengoForeignMethod_52(WrenVM* vm);
void wrengoForeignMethod_53(WrenVM* vm);
void wrengoForeignMethod_54(WrenVM* vm);
void wrengoForeignMethod_55(WrenVM* vm);
void wrengoForeignMethod_56(WrenVM* vm);
void wrengoForeignMethod_57(WrenVM* vm);
void wrengoForeignMethod_58(WrenVM* vm);
void wrengoForeignMethod_59(WrenVM* vm);
void wrengoForeignMethod_60(WrenVM* vm);
void wrengoForeignMethod_61(WrenVM* vm);
void wrengoForeignMethod_62(WrenVM* vm);
void wrengoForeignMethod_63(WrenVM* vm);
void wrengoForeignMethod_64(WrenVM* vm);
void wrengoForeignMethod_65(WrenVM* vm);
void wrengoForeignMethod_66(WrenVM* vm);
void wrengoForeignMethod_67(WrenVM* vm);
void wrengoForeignMethod_68(WrenVM* vm);
void wrengoForeignMethod_69(WrenVM* vm);
void wrengoForeignMethod_70(WrenVM* vm);
void wrengoForeignMethod_71(WrenVM* vm);
void wrengoForeignMethod_72(WrenVM* vm);
void wrengoForeignMethod_73(WrenVM* vm);
void wrengoForeignMethod_74(WrenVM* vm);
void wrengoForeignMethod_75(WrenVM* vm);
void wrengoForeignMethod_76(WrenVM* vm);
void wrengoForeignMethod_77(WrenVM* vm);
void wrengoForeignMethod_78(WrenVM* vm);
void wrengoForeignMethod_79(WrenVM* vm);
void wrengoForeignMethod_80(WrenVM* vm);
void wrengoForeignMethod_81(WrenVM* vm);
void wrengoForeignMethod_82(WrenVM* vm);
void wrengoForeignMethod_83(WrenVM* vm);
void wrengoForeignMethod_84(WrenVM* vm);
void wrengoForeignMethod_85(WrenVM* vm);
void wrengoForeignMethod_86(WrenVM* vm);
void wrengoForeignMethod_87(WrenVM* vm);
void wrengoForeignMethod_88(WrenVM* vm);
void wrengoForeignMethod_89(WrenVM* vm);
void wrengoForeignMethod_90(WrenVM* vm);
void wrengoForeignMethod_91(WrenVM* vm);
void wrengoForeignMethod_92(WrenVM* vm);
void wrengoForeignMethod_93(WrenVM* vm);
void wrengoForeignMethod_94(WrenVM* vm);
void wrengoForeignMethod_95(WrenVM* vm);
void wrengoForeignMethod_96(WrenVM* vm);
void wrengoForeignMethod_97(WrenVM* vm);
void wrengoForeignMethod_98(WrenVM* vm);

#endif
