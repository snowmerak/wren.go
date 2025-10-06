package wrengo

// #include <stdlib.h>
// #include "wren.h"
// #include "wren_callbacks.h"
import "C"
import (
	"sync"
	"unsafe"
)

const maxForeignMethods = 300 // Maximum number of foreign methods supported

// vmRegistry stores VM instances for callback access
var vmRegistry = make(map[*C.WrenVM]*WrenVM)
var vmMutex sync.RWMutex

// registerVM associates a C VM pointer with a Go WrenVM instance
func registerVM(vm *WrenVM) {
	vmMutex.Lock()
	defer vmMutex.Unlock()
	vmRegistry[vm.vm] = vm
}

// unregisterVM removes the association
func unregisterVM(vm *WrenVM) {
	vmMutex.Lock()
	defer vmMutex.Unlock()
	delete(vmRegistry, vm.vm)

	// Also clean up foreign data for this VM
	foreignDataMutex.Lock()
	defer foreignDataMutex.Unlock()
	delete(vmForeignDataStore, vm.vm)
}

// getVM retrieves the Go WrenVM instance from a C VM pointer
func getVM(cvm *C.WrenVM) *WrenVM {
	vmMutex.RLock()
	defer vmMutex.RUnlock()
	return vmRegistry[cvm]
}

// vmForeignData holds foreign function data for a single VM
type vmForeignData struct {
	wrapperMethods map[int]ForeignMethodFn  // wrapperID → function
	nextWrapperID  int                      // Next available wrapper ID (0-98)
	allocators     map[string]*ForeignClass // "module:className" → ForeignClass
}

// Foreign function callback registry - per VM
var (
	foreignDataMutex   sync.RWMutex
	vmForeignDataStore = make(map[*C.WrenVM]*vmForeignData)
)

//export wrengoBindForeignMethod
func wrengoBindForeignMethod(cvm *C.WrenVM, cModule, cClassName *C.char, isStatic C.bool, cSignature *C.char) C.WrenForeignMethodFn {
	module := C.GoString(cModule)
	className := C.GoString(cClassName)
	signature := C.GoString(cSignature)

	fn := lookupForeignMethod(module, className, bool(isStatic), signature)
	if fn == nil {
		return nil
	}

	foreignDataMutex.Lock()
	defer foreignDataMutex.Unlock()

	// Get or create foreign data for this VM
	data := vmForeignDataStore[cvm]
	if data == nil {
		data = &vmForeignData{
			wrapperMethods: make(map[int]ForeignMethodFn),
			nextWrapperID:  0,
			allocators:     make(map[string]*ForeignClass),
		}
		vmForeignDataStore[cvm] = data
	}

	// Check if we've exceeded the wrapper limit
	if data.nextWrapperID >= maxForeignMethods {
		panic("Maximum foreign methods (300) exceeded for this VM")
	}

	// Assign the next available wrapper ID
	wrapperID := data.nextWrapperID
	data.wrapperMethods[wrapperID] = fn
	data.nextWrapperID++

	// Return the appropriate C function pointer based on wrapper ID
	switch wrapperID {
	case 0:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_0)
	case 1:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_1)
	case 2:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_2)
	case 3:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_3)
	case 4:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_4)
	case 5:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_5)
	case 6:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_6)
	case 7:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_7)
	case 8:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_8)
	case 9:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_9)
	case 10:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_10)
	case 11:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_11)
	case 12:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_12)
	case 13:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_13)
	case 14:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_14)
	case 15:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_15)
	case 16:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_16)
	case 17:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_17)
	case 18:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_18)
	case 19:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_19)
	case 20:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_20)
	case 21:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_21)
	case 22:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_22)
	case 23:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_23)
	case 24:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_24)
	case 25:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_25)
	case 26:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_26)
	case 27:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_27)
	case 28:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_28)
	case 29:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_29)
	case 30:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_30)
	case 31:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_31)
	case 32:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_32)
	case 33:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_33)
	case 34:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_34)
	case 35:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_35)
	case 36:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_36)
	case 37:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_37)
	case 38:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_38)
	case 39:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_39)
	case 40:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_40)
	case 41:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_41)
	case 42:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_42)
	case 43:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_43)
	case 44:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_44)
	case 45:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_45)
	case 46:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_46)
	case 47:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_47)
	case 48:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_48)
	case 49:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_49)
	case 50:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_50)
	case 51:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_51)
	case 52:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_52)
	case 53:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_53)
	case 54:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_54)
	case 55:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_55)
	case 56:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_56)
	case 57:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_57)
	case 58:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_58)
	case 59:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_59)
	case 60:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_60)
	case 61:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_61)
	case 62:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_62)
	case 63:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_63)
	case 64:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_64)
	case 65:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_65)
	case 66:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_66)
	case 67:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_67)
	case 68:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_68)
	case 69:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_69)
	case 70:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_70)
	case 71:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_71)
	case 72:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_72)
	case 73:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_73)
	case 74:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_74)
	case 75:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_75)
	case 76:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_76)
	case 77:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_77)
	case 78:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_78)
	case 79:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_79)
	case 80:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_80)
	case 81:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_81)
	case 82:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_82)
	case 83:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_83)
	case 84:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_84)
	case 85:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_85)
	case 86:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_86)
	case 87:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_87)
	case 88:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_88)
	case 89:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_89)
	case 90:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_90)
	case 91:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_91)
	case 92:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_92)
	case 93:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_93)
	case 94:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_94)
	case 95:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_95)
	case 96:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_96)
	case 97:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_97)
	case 98:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_98)
	case 99:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_99)
	case 100:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_100)
	case 101:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_101)
	case 102:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_102)
	case 103:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_103)
	case 104:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_104)
	case 105:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_105)
	case 106:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_106)
	case 107:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_107)
	case 108:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_108)
	case 109:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_109)
	case 110:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_110)
	case 111:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_111)
	case 112:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_112)
	case 113:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_113)
	case 114:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_114)
	case 115:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_115)
	case 116:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_116)
	case 117:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_117)
	case 118:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_118)
	case 119:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_119)
	case 120:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_120)
	case 121:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_121)
	case 122:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_122)
	case 123:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_123)
	case 124:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_124)
	case 125:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_125)
	case 126:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_126)
	case 127:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_127)
	case 128:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_128)
	case 129:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_129)
	case 130:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_130)
	case 131:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_131)
	case 132:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_132)
	case 133:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_133)
	case 134:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_134)
	case 135:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_135)
	case 136:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_136)
	case 137:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_137)
	case 138:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_138)
	case 139:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_139)
	case 140:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_140)
	case 141:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_141)
	case 142:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_142)
	case 143:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_143)
	case 144:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_144)
	case 145:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_145)
	case 146:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_146)
	case 147:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_147)
	case 148:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_148)
	case 149:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_149)
	case 150:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_150)
	case 151:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_151)
	case 152:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_152)
	case 153:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_153)
	case 154:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_154)
	case 155:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_155)
	case 156:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_156)
	case 157:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_157)
	case 158:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_158)
	case 159:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_159)
	case 160:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_160)
	case 161:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_161)
	case 162:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_162)
	case 163:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_163)
	case 164:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_164)
	case 165:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_165)
	case 166:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_166)
	case 167:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_167)
	case 168:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_168)
	case 169:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_169)
	case 170:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_170)
	case 171:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_171)
	case 172:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_172)
	case 173:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_173)
	case 174:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_174)
	case 175:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_175)
	case 176:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_176)
	case 177:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_177)
	case 178:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_178)
	case 179:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_179)
	case 180:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_180)
	case 181:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_181)
	case 182:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_182)
	case 183:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_183)
	case 184:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_184)
	case 185:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_185)
	case 186:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_186)
	case 187:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_187)
	case 188:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_188)
	case 189:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_189)
	case 190:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_190)
	case 191:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_191)
	case 192:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_192)
	case 193:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_193)
	case 194:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_194)
	case 195:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_195)
	case 196:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_196)
	case 197:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_197)
	case 198:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_198)
	case 199:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_199)
	case 200:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_200)
	case 201:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_201)
	case 202:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_202)
	case 203:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_203)
	case 204:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_204)
	case 205:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_205)
	case 206:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_206)
	case 207:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_207)
	case 208:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_208)
	case 209:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_209)
	case 210:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_210)
	case 211:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_211)
	case 212:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_212)
	case 213:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_213)
	case 214:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_214)
	case 215:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_215)
	case 216:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_216)
	case 217:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_217)
	case 218:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_218)
	case 219:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_219)
	case 220:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_220)
	case 221:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_221)
	case 222:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_222)
	case 223:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_223)
	case 224:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_224)
	case 225:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_225)
	case 226:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_226)
	case 227:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_227)
	case 228:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_228)
	case 229:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_229)
	case 230:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_230)
	case 231:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_231)
	case 232:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_232)
	case 233:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_233)
	case 234:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_234)
	case 235:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_235)
	case 236:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_236)
	case 237:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_237)
	case 238:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_238)
	case 239:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_239)
	case 240:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_240)
	case 241:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_241)
	case 242:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_242)
	case 243:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_243)
	case 244:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_244)
	case 245:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_245)
	case 246:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_246)
	case 247:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_247)
	case 248:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_248)
	case 249:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_249)
	case 250:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_250)
	case 251:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_251)
	case 252:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_252)
	case 253:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_253)
	case 254:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_254)
	case 255:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_255)
	case 256:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_256)
	case 257:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_257)
	case 258:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_258)
	case 259:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_259)
	case 260:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_260)
	case 261:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_261)
	case 262:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_262)
	case 263:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_263)
	case 264:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_264)
	case 265:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_265)
	case 266:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_266)
	case 267:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_267)
	case 268:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_268)
	case 269:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_269)
	case 270:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_270)
	case 271:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_271)
	case 272:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_272)
	case 273:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_273)
	case 274:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_274)
	case 275:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_275)
	case 276:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_276)
	case 277:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_277)
	case 278:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_278)
	case 279:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_279)
	case 280:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_280)
	case 281:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_281)
	case 282:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_282)
	case 283:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_283)
	case 284:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_284)
	case 285:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_285)
	case 286:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_286)
	case 287:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_287)
	case 288:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_288)
	case 289:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_289)
	case 290:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_290)
	case 291:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_291)
	case 292:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_292)
	case 293:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_293)
	case 294:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_294)
	case 295:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_295)
	case 296:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_296)
	case 297:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_297)
	case 298:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_298)
	case 299:
		return C.WrenForeignMethodFn(C.wrengoForeignMethod_299)
	}

	return nil
}

//export wrengoBindForeignClass
func wrengoBindForeignClass(cvm *C.WrenVM, cModule, cClassName *C.char) C.WrenForeignClassMethods {
	var methods C.WrenForeignClassMethods

	module := C.GoString(cModule)
	className := C.GoString(cClassName)

	class := lookupForeignClass(module, className)
	if class == nil {
		return methods
	}

	foreignDataMutex.Lock()
	defer foreignDataMutex.Unlock()

	// Get or create foreign data for this VM
	data := vmForeignDataStore[cvm]
	if data == nil {
		data = &vmForeignData{
			wrapperMethods: make(map[int]ForeignMethodFn),
			nextWrapperID:  0,
			allocators:     make(map[string]*ForeignClass),
		}
		vmForeignDataStore[cvm] = data
	}

	// Store the class in VM-specific registry
	key := module + ":" + className
	data.allocators[key] = class

	if class.Allocate != nil {
		methods.allocate = C.WrenForeignMethodFn(C.wrengoForeignAllocateCallback)
	}

	if class.Finalize != nil {
		methods.finalize = C.WrenFinalizerFn(C.wrengoForeignFinalizeCallback)
	}

	return methods
}

//export goForeignMethodCallback
func goForeignMethodCallback(cvm *C.WrenVM, wrapperId C.int) {
	vm := getVM(cvm)
	if vm == nil {
		return
	}

	foreignDataMutex.RLock()
	data := vmForeignDataStore[cvm]
	foreignDataMutex.RUnlock()

	if data == nil {
		return
	}

	// Look up the function for this wrapper ID
	if fn, ok := data.wrapperMethods[int(wrapperId)]; ok {
		fn(vm)
	}
}

//export wrengoForeignAllocateCallback
func wrengoForeignAllocateCallback(cvm *C.WrenVM) {
	vm := getVM(cvm)
	if vm == nil {
		return
	}

	foreignDataMutex.RLock()
	data := vmForeignDataStore[cvm]
	foreignDataMutex.RUnlock()

	if data == nil {
		return
	}

	// Get the class being allocated from slot 0
	// Note: Wren passes the class in slot 0 during allocation
	// We need to find which class is being allocated
	// For now, call the first available allocator
	// This is a limitation - ideally we'd track which class is being allocated
	for _, class := range data.allocators {
		if class.Allocate != nil {
			class.Allocate(vm)
			return
		}
	}
}

//export wrengoForeignFinalizeCallback
func wrengoForeignFinalizeCallback(data unsafe.Pointer) {
	// Note: We don't have VM context here, so we can't call VM-specific finalizers
	// This is a limitation of the current design
	// Finalizers would need to be stored differently to work properly
}
