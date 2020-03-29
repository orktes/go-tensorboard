package events

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/Applifier/go-tensorflow/types/tensorflow/core/framework"
	"github.com/pkg/errors"
)

// TensorContentToGoType converts TensorContent to a go type
func TensorContentToGoType(tensor *framework.TensorProto) (interface{}, error) {
	if len(tensor.TensorContent) == 0 {
		return nil, errors.New("tensor has no content defined")
	}

	if tensor.Dtype != framework.DataType_DT_STRING {
		typ := typeOf(tensor.Dtype, tensor.TensorShape.Dim)
		val := reflect.New(typ)
		if err := decodeTensor(bytes.NewReader(tensor.TensorContent), tensor.TensorShape.Dim, typ, val); err != nil {
			return nil, err
		}
		return reflect.Indirect(val).Interface(), nil
	}

	return nil, errors.New("string tensors not supported")
}

func decodeTensor(r *bytes.Reader, shape []*framework.TensorShapeProto_Dim, typ reflect.Type, ptr reflect.Value) error {
	switch typ.Kind() {
	case reflect.Bool:
		b, err := r.ReadByte()
		if err != nil {
			return err
		}
		ptr.Elem().SetBool(b == 1)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		if err := binary.Read(r, nativeEndian, ptr.Interface()); err != nil {
			return err
		}

	case reflect.Slice:
		val := reflect.Indirect(ptr)
		val.Set(reflect.MakeSlice(typ, int(shape[0].Size_), int(shape[0].Size_)))

		// Optimization: if only one dimension is left we can use binary.Read() directly for this slice
		if len(shape) == 1 && val.Len() > 0 {
			switch val.Index(0).Kind() {
			case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
				return binary.Read(r, nativeEndian, val.Interface())
			}
		}

		for i := 0; i < val.Len(); i++ {
			if err := decodeTensor(r, shape[1:], typ.Elem(), val.Index(i).Addr()); err != nil {
				return err
			}
		}

	default:
		return fmt.Errorf("unsupported type %v", typ)
	}
	return nil
}

var types = []struct {
	typ      reflect.Type
	dataType framework.DataType
}{
	{reflect.TypeOf(float32(0)), framework.DataType_DT_FLOAT},
	{reflect.TypeOf(float64(0)), framework.DataType_DT_DOUBLE},
	{reflect.TypeOf(int32(0)), framework.DataType_DT_INT32},
	{reflect.TypeOf(uint32(0)), framework.DataType_DT_UINT32},
	{reflect.TypeOf(int16(0)), framework.DataType_DT_INT16},
	{reflect.TypeOf(int8(0)), framework.DataType_DT_INT8},
	{reflect.TypeOf(uint8(0)), framework.DataType_DT_UINT8},
	{reflect.TypeOf(""), framework.DataType_DT_STRING},
	{reflect.TypeOf(complex(float32(0), float32(0))), framework.DataType_DT_COMPLEX64},
	{reflect.TypeOf(int64(0)), framework.DataType_DT_INT64},
	{reflect.TypeOf(uint64(0)), framework.DataType_DT_UINT64},
	{reflect.TypeOf(false), framework.DataType_DT_BOOL},
	{reflect.TypeOf(complex(float64(0), float64(0))), framework.DataType_DT_COMPLEX128},
	// TODO: support DT_RESOURCE representation in go.
	// TODO: support DT_VARIANT representation in go.
}

// typeOf converts from a DataType and Shape to the equivalent Go type.
func typeOf(dt framework.DataType, shape []*framework.TensorShapeProto_Dim) reflect.Type {
	var ret reflect.Type
	for _, t := range types {
		if dt == t.dataType {
			ret = t.typ
			break
		}
	}
	if ret == nil {
		panic(fmt.Sprintf("DataType %v is not supported", dt))
	}
	for range shape {
		ret = reflect.SliceOf(ret)
	}
	return ret
}

var nativeEndian binary.ByteOrder

func init() {
	buf := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)

	switch buf {
	case [2]byte{0xCD, 0xAB}:
		nativeEndian = binary.LittleEndian
	case [2]byte{0xAB, 0xCD}:
		nativeEndian = binary.BigEndian
	default:
		panic("Could not determine native endianness.")
	}
}
