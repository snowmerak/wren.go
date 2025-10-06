package wrengo

// #include "wren.h"
import "C"
import "unsafe"

// SlotType represents the type of a value in a slot.
type SlotType int

const (
	TypeBool    SlotType = C.WREN_TYPE_BOOL
	TypeNum     SlotType = C.WREN_TYPE_NUM
	TypeForeign SlotType = C.WREN_TYPE_FOREIGN
	TypeList    SlotType = C.WREN_TYPE_LIST
	TypeMap     SlotType = C.WREN_TYPE_MAP
	TypeNull    SlotType = C.WREN_TYPE_NULL
	TypeString  SlotType = C.WREN_TYPE_STRING
	TypeUnknown SlotType = C.WREN_TYPE_UNKNOWN
)

// GetSlotCount returns the number of slots available to the current foreign method.
func (vm *WrenVM) GetSlotCount() int {
	return int(C.wrenGetSlotCount(vm.vm))
}

// EnsureSlots ensures that the foreign method stack has at least numSlots available.
func (vm *WrenVM) EnsureSlots(numSlots int) {
	C.wrenEnsureSlots(vm.vm, C.int(numSlots))
}

// GetSlotType gets the type of the object in the given slot.
func (vm *WrenVM) GetSlotType(slot int) SlotType {
	return SlotType(C.wrenGetSlotType(vm.vm, C.int(slot)))
}

// GetSlotBool reads a boolean value from the given slot.
func (vm *WrenVM) GetSlotBool(slot int) bool {
	return bool(C.wrenGetSlotBool(vm.vm, C.int(slot)))
}

// GetSlotDouble reads a number from the given slot.
func (vm *WrenVM) GetSlotDouble(slot int) float64 {
	return float64(C.wrenGetSlotDouble(vm.vm, C.int(slot)))
}

// GetSlotString reads a string from the given slot.
func (vm *WrenVM) GetSlotString(slot int) string {
	cStr := C.wrenGetSlotString(vm.vm, C.int(slot))
	return C.GoString(cStr)
}

// GetSlotBytes reads a byte array from the given slot.
func (vm *WrenVM) GetSlotBytes(slot int) []byte {
	var length C.int
	cBytes := C.wrenGetSlotBytes(vm.vm, C.int(slot), &length)
	return C.GoBytes(unsafe.Pointer(cBytes), length)
}

// GetSlotForeign reads a foreign object from the given slot and returns a pointer to its data.
func (vm *WrenVM) GetSlotForeign(slot int) unsafe.Pointer {
	return C.wrenGetSlotForeign(vm.vm, C.int(slot))
}

// SetSlotBool stores a boolean value in the given slot.
func (vm *WrenVM) SetSlotBool(slot int, value bool) {
	C.wrenSetSlotBool(vm.vm, C.int(slot), C.bool(value))
}

// SetSlotDouble stores a numeric value in the given slot.
func (vm *WrenVM) SetSlotDouble(slot int, value float64) {
	C.wrenSetSlotDouble(vm.vm, C.int(slot), C.double(value))
}

// SetSlotString stores a string in the given slot.
func (vm *WrenVM) SetSlotString(slot int, text string) {
	cStr := C.CString(text)
	defer C.free(unsafe.Pointer(cStr))
	C.wrenSetSlotString(vm.vm, C.int(slot), cStr)
}

// SetSlotBytes stores a byte array in the given slot.
func (vm *WrenVM) SetSlotBytes(slot int, bytes []byte) {
	if len(bytes) == 0 {
		C.wrenSetSlotBytes(vm.vm, C.int(slot), nil, 0)
		return
	}
	C.wrenSetSlotBytes(vm.vm, C.int(slot), (*C.char)(unsafe.Pointer(&bytes[0])), C.size_t(len(bytes)))
}

// SetSlotNull stores null in the given slot.
func (vm *WrenVM) SetSlotNull(slot int) {
	C.wrenSetSlotNull(vm.vm, C.int(slot))
}

// SetSlotNewList stores a new empty list in the given slot.
func (vm *WrenVM) SetSlotNewList(slot int) {
	C.wrenSetSlotNewList(vm.vm, C.int(slot))
}

// SetSlotNewMap stores a new empty map in the given slot.
func (vm *WrenVM) SetSlotNewMap(slot int) {
	C.wrenSetSlotNewMap(vm.vm, C.int(slot))
}

// GetListCount returns the number of elements in the list stored in the given slot.
func (vm *WrenVM) GetListCount(slot int) int {
	return int(C.wrenGetListCount(vm.vm, C.int(slot)))
}

// GetListElement reads element at index from the list in listSlot and stores it in elementSlot.
func (vm *WrenVM) GetListElement(listSlot, index, elementSlot int) {
	C.wrenGetListElement(vm.vm, C.int(listSlot), C.int(index), C.int(elementSlot))
}

// SetListElement sets the value at index in the list at listSlot to the value from elementSlot.
func (vm *WrenVM) SetListElement(listSlot, index, elementSlot int) {
	C.wrenSetListElement(vm.vm, C.int(listSlot), C.int(index), C.int(elementSlot))
}

// InsertInList inserts the value from elementSlot into the list at listSlot at the given index.
// Negative indexes can be used to insert from the end. Use -1 to append.
func (vm *WrenVM) InsertInList(listSlot, index, elementSlot int) {
	C.wrenInsertInList(vm.vm, C.int(listSlot), C.int(index), C.int(elementSlot))
}

// GetMapCount returns the number of entries in the map stored in the given slot.
func (vm *WrenVM) GetMapCount(slot int) int {
	return int(C.wrenGetMapCount(vm.vm, C.int(slot)))
}

// GetMapContainsKey returns true if the key in keySlot is found in the map in mapSlot.
func (vm *WrenVM) GetMapContainsKey(mapSlot, keySlot int) bool {
	return bool(C.wrenGetMapContainsKey(vm.vm, C.int(mapSlot), C.int(keySlot)))
}

// GetMapValue retrieves a value with the key in keySlot from the map in mapSlot and stores it in valueSlot.
func (vm *WrenVM) GetMapValue(mapSlot, keySlot, valueSlot int) {
	C.wrenGetMapValue(vm.vm, C.int(mapSlot), C.int(keySlot), C.int(valueSlot))
}

// SetMapValue inserts the value from valueSlot into the map at mapSlot with the key from keySlot.
func (vm *WrenVM) SetMapValue(mapSlot, keySlot, valueSlot int) {
	C.wrenSetMapValue(vm.vm, C.int(mapSlot), C.int(keySlot), C.int(valueSlot))
}

// RemoveMapValue removes a value from the map in mapSlot with the key from keySlot,
// and places it in removedValueSlot. If not found, removedValueSlot is set to null.
func (vm *WrenVM) RemoveMapValue(mapSlot, keySlot, removedValueSlot int) {
	C.wrenRemoveMapValue(vm.vm, C.int(mapSlot), C.int(keySlot), C.int(removedValueSlot))
}

// GetVariable looks up the top level variable with name in the given module and stores it in slot.
func (vm *WrenVM) GetVariable(module, name string, slot int) {
	cModule := C.CString(module)
	defer C.free(unsafe.Pointer(cModule))
	
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	
	C.wrenGetVariable(vm.vm, cModule, cName, C.int(slot))
}

// HasVariable returns true if the variable exists in the module.
func (vm *WrenVM) HasVariable(module, name string) bool {
	cModule := C.CString(module)
	defer C.free(unsafe.Pointer(cModule))
	
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	
	return bool(C.wrenHasVariable(vm.vm, cModule, cName))
}

// HasModule returns true if the module has been imported/resolved before.
func (vm *WrenVM) HasModule(module string) bool {
	cModule := C.CString(module)
	defer C.free(unsafe.Pointer(cModule))
	
	return bool(C.wrenHasModule(vm.vm, cModule))
}

// AbortFiber sets the current fiber to be aborted and uses the value in slot as the error object.
func (vm *WrenVM) AbortFiber(slot int) {
	C.wrenAbortFiber(vm.vm, C.int(slot))
}

// GetUserData returns the user data associated with the VM.
func (vm *WrenVM) GetUserData() unsafe.Pointer {
	return C.wrenGetUserData(vm.vm)
}

// SetUserData sets user data associated with the VM.
func (vm *WrenVM) SetUserData(userData unsafe.Pointer) {
	C.wrenSetUserData(vm.vm, userData)
}
