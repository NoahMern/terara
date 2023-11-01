package types

import (
	"encoding/binary"
	"errors"
	"math"
	"strconv"
)

// this represents data types that can be stored in a document
const (
	// NullType represents a null value
	NullType byte = iota
	BoolType
	Int64Type
	Int32Type
	FloatType
	StringType
	CharType

	BigIntType
	BigFloatType

	ArrayType
	DocumentType
	CollectionType

	DateType
	TimeStampType

	EmailType
	PhoneType
	MoneyType
	UUIDType

	BinaryType
	BlobType

	LongitudeType
	LatitudeType
	CurrencyCodeType
	CountryCodeType

	// internal types
	EOFType
	NameType

	// this is the last type
	LastType
)

const (
	NullTerm = '\x00'
)

func IsInternal(t byte) bool {
	return t == EOFType || t == NameType || t == LastType
}

var (
	// ErrInvalidLength is returned when the length of the byte slice is invalid
	ErrInvalidLength = errors.New("invalid length")
	// ErrInvalidType is returned when the type is invalid
	ErrInvalidType = errors.New("invalid type")
)

// value is an interface that represents values that the database can work with
type Object interface {
	// returns the type name
	Type() byte
	// returns the value
	Value() interface{}
	// returns the value as a string
	String() string
}

type Marshaler interface {
	MarshalObject() ([]byte, error)
}

// an interface that represents types that can be unmarshaled
// returns the number of bytes read and an error
type Unmarshaler interface {
	UnmarshalObject([]byte) (int, error)
}

type Int64 int64

func (i Int64) Type() byte {
	return Int64Type
}

func (i Int64) Value() interface{} {
	return int64(i)
}

func (i Int64) String() string {
	return strconv.FormatInt(int64(i), 10)
}

func (i Int64) MarshalObject() ([]byte, error) {
	// encode the int64 to a 8 two's complement byte slice
	b := make([]byte, 9)
	b[0] = i.Type()
	binary.BigEndian.PutUint64(b[1:], uint64(i))
	return b, nil
}

func (i *Int64) UnmarshalObject(b []byte) (int, error) {
	// check length
	if len(b) < 9 {
		return 0, ErrInvalidLength
	}
	// check type
	if b[0] != i.Type() {
		return 0, ErrInvalidType
	}
	// decode the byte slice into an int64
	*i = Int64(int64(binary.BigEndian.Uint64(b[1:9])))
	return 9, nil
}

type Int32 int32

func (i Int32) Type() byte {
	return Int32Type
}

func (i Int32) Value() interface{} {
	return int32(i)
}

func (i Int32) String() string {
	return strconv.FormatInt(int64(i), 10)
}

func (i Int32) MarshalObject() ([]byte, error) {
	// encode the int32 to a 4 two's complement byte slice
	b := make([]byte, 5)
	b[0] = i.Type()
	binary.BigEndian.PutUint32(b[1:], uint32(i))
	return b, nil
}

func (i *Int32) UnmarshalObject(b []byte) (int, error) {
	// check length
	if len(b) < 5 {
		return 0, ErrInvalidLength
	}
	// check type
	if b[0] != i.Type() {
		return 0, ErrInvalidType
	}
	// decode the byte slice into an int32
	*i = Int32(int32(binary.BigEndian.Uint32(b[1:5])))
	return 5, nil
}

type Float float64

func (f Float) Type() byte {
	return FloatType
}

func (f Float) Value() interface{} {
	return float64(f)
}

func (f Float) String() string {
	return strconv.FormatFloat(float64(f), 'f', -1, 64)
}

func (f Float) MarshalObject() ([]byte, error) {
	// encode the float64 to a 8 byte slice 64-bit IEEE 754-2008 binary floating point
	b := make([]byte, 9)
	b[0] = f.Type()
	binary.BigEndian.PutUint64(b[1:], math.Float64bits(float64(f)))
	return b, nil
}

func (f *Float) UnmarshalObject(b []byte) (int, error) {
	// check length
	if len(b) < 9 {
		return 0, ErrInvalidLength
	}
	// check type
	if b[0] != f.Type() {
		return 0, ErrInvalidType
	}
	// decode the byte slice into a float64
	*f = Float(math.Float64frombits(binary.BigEndian.Uint64(b[1:9])))
	return 9, nil
}

type String string

func (s String) Type() byte {
	return StringType
}

func (s String) Value() interface{} {
	return string(s)
}

func (s String) String() string {
	return string(s)
}

func (s String) MarshalObject() ([]byte, error) {
	// encode the string as a c string
	b := make([]byte, len(s)+2)
	b[0] = s.Type()
	copy(b[1:], []byte(s))
	b[len(b)-1] = NullTerm
	return b, nil
}

func (s *String) UnmarshalObject(b []byte) (int, error) {
	if len(b) < 2 {
		return 0, ErrInvalidLength
	}
	if b[0] != s.Type() {
		return 0, ErrInvalidType
	}
	// find the null terminator
	for i := 1; i < len(b); i++ {
		if b[i] == NullTerm {
			*s = String(string(b[1:i]))
			return i + 1, nil
		}
	}
	return 0, ErrInvalidLength
}

type Bool bool

func (b Bool) Type() byte {
	return BoolType
}

func (b Bool) Value() interface{} {
	return bool(b)
}

func (b Bool) String() string {
	return strconv.FormatBool(bool(b))
}

func (b Bool) MarshalObject() ([]byte, error) {
	// encode the bool to a 1 byte slice
	return []byte{b.Type(), boolToByte(bool(b))}, nil
}

func (b *Bool) UnmarshalObject(by []byte) (int, error) {
	// check length
	if len(by) < 1 {
		return 0, ErrInvalidLength
	}
	// check type
	if by[0] != b.Type() {
		return 0, ErrInvalidType
	}
	// decode the byte slice into a bool
	*b = Bool(byteToBool(by[1]))
	return 2, nil
}

func boolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

func byteToBool(b byte) bool {
	if b == 1 {
		return true
	}
	return false
}

type Null struct{}

func (n Null) Type() byte {
	return NullType
}

func (n Null) Value() interface{} {
	return nil
}

func (n Null) String() string {
	return "null"
}

func (n Null) MarshalObject() ([]byte, error) {
	return []byte{NullType}, nil
}

func (n *Null) UnmarshalObject(b []byte) (int, error) {
	// check length
	if len(b) < 1 {
		return 0, ErrInvalidLength
	}
	if b[0] != n.Type() {
		return 0, ErrInvalidType
	}
	return 1, nil
}

type Char byte

func (c Char) Type() byte {
	return CharType
}

func (c Char) Value() interface{} {
	return c
}

func (c Char) String() string {
	return string(c)
}

func (c Char) MarshalObject() ([]byte, error) {
	return []byte{c.Type(), byte(c)}, nil
}

func (c *Char) UnmarshalObject(b []byte) (int, error) {
	// check length
	if len(b) < 2 {
		return 0, ErrInvalidLength
	}
	if b[0] != c.Type() {
		return 0, ErrInvalidType
	}
	*c = Char(rune(b[1]))
	return 2, nil
}

type EOF struct{}

func (e EOF) Type() byte {
	return EOFType
}

func (e EOF) Value() interface{} {
	return nil
}

func (e EOF) String() string {
	return "EOF"
}

func (e EOF) MarshalObject() ([]byte, error) {
	return []byte{EOFType}, nil
}

func (e *EOF) UnmarshalObject(b []byte) (int, error) {
	if len(b) < 1 {
		return 0, ErrInvalidLength
	}
	if b[0] != e.Type() {
		return 0, ErrInvalidType
	}
	return 1, nil
}

type Name string

func (n Name) Type() byte {
	return NameType
}

func (n Name) Value() interface{} {
	return string(n)
}

func (n Name) String() string {
	return string(n)
}

func (n Name) MarshalObject() ([]byte, error) {
	// marshal name as a c string
	b := make([]byte, 1)
	b[0] = n.Type()
	b = append(b, []byte(n)...)
	b = append(b, NullTerm)
	return b, nil
}

func (n *Name) UnmarshalObject(b []byte) (int, error) {
	if len(b) < 2 {
		return 0, ErrInvalidLength
	}
	if b[0] != n.Type() {
		return 0, ErrInvalidType
	}
	// find the null terminator
	for i := 1; i < len(b); i++ {
		if b[i] == NullTerm {
			*n = Name(string(b[1:i]))
			return i + 1, nil
		}
	}
	return 0, ErrInvalidLength
}

func UnmarshalInt64(b []byte) (Int64, int, error) {
	var i Int64
	n, err := i.UnmarshalObject(b)
	return i, n, err
}

func UnmarshalInt32(b []byte) (Int32, int, error) {
	var i Int32
	n, err := i.UnmarshalObject(b)
	return i, n, err
}

func UnmarshalFloat(b []byte) (Float, int, error) {
	var f Float
	n, err := f.UnmarshalObject(b)
	return f, n, err
}

func UnmarshalString(b []byte) (String, int, error) {
	var s String
	n, err := s.UnmarshalObject(b)
	return s, n, err
}

func UnmarshalBool(b []byte) (Bool, int, error) {
	var bo Bool
	n, err := bo.UnmarshalObject(b)
	return bo, n, err
}

func UnmarshalNull(b []byte) (Null, int, error) {
	var null Null
	n, err := null.UnmarshalObject(b)
	return null, n, err
}

func UnmarshalChar(b []byte) (Char, int, error) {
	var c Char
	n, err := c.UnmarshalObject(b)
	return c, n, err
}

func UnmarshalEOF(b []byte) (EOF, int, error) {
	var e EOF
	n, err := e.UnmarshalObject(b)
	return e, n, err
}

func UnmarshalName(b []byte) (Name, int, error) {
	var name Name
	n, err := name.UnmarshalObject(b)
	return name, n, err
}

func UnmarshalObject(b []byte) (Object, int, error) {
	if len(b) < 1 {
		return nil, 0, ErrInvalidLength
	}
	switch b[0] {
	case Int64Type:
		return UnmarshalInt64(b)
	case Int32Type:
		return UnmarshalInt32(b)
	case FloatType:
		return UnmarshalFloat(b)
	case StringType:
		return UnmarshalString(b)
	case BoolType:
		return UnmarshalBool(b)
	case NullType:
		return UnmarshalNull(b)
	case CharType:
		return UnmarshalChar(b)
	case EOFType:
		return UnmarshalEOF(b)
	case NameType:
		return UnmarshalName(b)
	case ArrayType:
		return UnmarshalArray(b)
	}
	return nil, 0, ErrInvalidType
}

// change an object into a marshaler
func MarshalObject(o Object) ([]byte, error) {
	switch o.Type() {
	case Int64Type:
		return o.(Int64).MarshalObject()
	case Int32Type:
		return o.(Int32).MarshalObject()
	case FloatType:
		return o.(Float).MarshalObject()
	case StringType:
		return o.(String).MarshalObject()
	case BoolType:
		return o.(Bool).MarshalObject()
	case NullType:
		return o.(Null).MarshalObject()
	case CharType:
		return o.(Char).MarshalObject()
	case EOFType:
		return o.(EOF).MarshalObject()
	case NameType:
		return o.(Name).MarshalObject()
	case ArrayType:
		return o.(Array).MarshalObject()
	}
	return nil, ErrInvalidType
}

func NewInt32(i int32) (Int32, error) {
	return Int32(i), nil
}

func NewInt64(i int64) (Int64, error) {
	return Int64(i), nil
}

func NewFloat(f float64) (Float, error) {
	return Float(f), nil
}
