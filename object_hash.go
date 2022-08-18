package labs

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/rogpeppe/go-internal/fmtsort"
	"reflect"
	"unsafe"
)

type rtype struct {
	Size       uintptr
	ptrdata    uintptr // number of bytes in the type that can contain pointers
	hash       uint32  // hash of type; avoids computation in hash tables
	tflag      uint8   // extra type information flags
	Align      uint8   // alignment of variable with this type
	FieldAlign uint8   // alignment of struct field with this type
	kind       uint8   // enumeration for C
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	equal     func(unsafe.Pointer, unsafe.Pointer) bool
	Gcdata    *byte // garbage collection data
	str       int32 // string form
	ptrToThis int32 // type for pointer to this type, may be zero
}

type FValue struct {
	Rtype *rtype
	Ptr   unsafe.Pointer
	flag  uintptr
}

func (v *FValue) IsPureData() bool {
	if v.Rtype.FieldAlign < 8 {
		return true
	}
	typeSize := v.Rtype.Size
	nFlag := int(typeSize) / 64
	if typeSize%64 > 0 {
		nFlag++
	}
	ptr := unsafe.Pointer(v.Rtype.Gcdata)
	for nFlag > 0 {
		flag := (*byte)(ptr)
		if *flag > 0 {
			return false
		}
		ptr = unsafe.Add(ptr, 1)
		nFlag--
	}
	//fmt.Println("true")
	return true
}

var hash = md5.New()

func ObjectHash(a any) string {
	value := reflect.ValueOf(a)
	value = value.Elem()
	hash.Reset()
	var pHead uintptr = 0
	var pTail uintptr = 0
	var hashStruct func(a reflect.Value)
	var hashMem func()
	hashStruct = func(a reflect.Value) {
		if a.Kind() < reflect.Array {
			if pHead == 0 {
				pHead = a.UnsafeAddr()
			}
			pTail = a.UnsafeAddr() + a.Type().Size()
			return
		} else if a.Kind() == reflect.Struct {
			if pHead == 0 {
				pHead = a.UnsafeAddr()
			}
			fv := (*FValue)(unsafe.Pointer(&a))
			if fv.IsPureData() {
				pTail = a.UnsafeAddr() + fv.Rtype.Size
			} else {
				for i := 0; i < a.NumField(); i++ {
					field := a.Field(i)
					if field.Kind() >= reflect.Array {
						pTail = field.UnsafeAddr()
						hashStruct(field)
					}
				}
			}
			return
		}
		hashMem()
		switch a.Kind() {
		case reflect.Ptr:
			if a.IsNil() {
				return
			}
			hashStruct(a.Elem())
		case reflect.Array:
			fallthrough
		case reflect.Slice:
			if a.Len() == 0 {
				return
			}
			i := 0
			v := a.Index(i)
			s := v.Type().Size()
			if v.Kind() < reflect.Array {
				for ; i < a.Len(); i++ {
					v = a.Index(i)
					us := unsafe.Slice((*uintptr)(unsafe.Pointer(&v)), 2)
					bytes := unsafe.Slice((*byte)(unsafe.Pointer(us[1])), s)
					hash.Write(bytes)
				}
			} else {
				for ; i < a.Len(); i++ {
					v = a.Index(i)
					hashStruct(a.Index(i))
					hashMem()
				}
			}
		case reflect.Map:
			if a.Len() == 0 {
				return
			}
			sorted := fmtsort.Sort(a)
			////指针类型Key忽略
			//if sorted.Key[0].Type().Kind() >= reflect.Array {
			//	return
			//}
			if sorted.Value[0].Kind() < reflect.Array {
				s := sorted.Value[0].Type().Size()
				for _, v := range sorted.Value {
					//fmt.Println("hash map size(", s, "):::", i, ":", sorted.Key[i], "->", v)
					us := unsafe.Slice((*uintptr)(unsafe.Pointer(&v)), 2)
					bytes := unsafe.Slice((*byte)(unsafe.Pointer(us[1])), s)
					hash.Write(bytes)
				}
			} else {
				for _, v := range sorted.Value {
					hashStruct(v)
					hashMem()
				}
			}
		case reflect.String:
			hash.Sum([]byte(a.String()))
		case reflect.Chan:
			fallthrough
		case reflect.Func:
			return
		default:
		}
	}
	hashMem = func() {
		if pHead == 0 || pTail <= pHead {
			return
		}
		bytes := unsafe.Slice((*byte)(unsafe.Pointer(pHead)), pTail-pHead)
		hash.Write(bytes)
		pHead = 0
	}
	hashStruct(value)
	hashMem()
	return hex.EncodeToString(hash.Sum(nil))
}
