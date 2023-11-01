package types

import (
	"errors"
)

var (
	ErrInvalidDocument = errors.New("invalid document")
)

// this represents the document type
type Document interface {
	Object

	ID() Object
	Del(key []byte) error
	Get(key []byte) (Object, error)
	Set(key []byte, value Object) error
	Keys() [][]byte
}

// takes a document and returns a byte slice
func GenericDocumentUnmarshaler(doc Document) ([]byte, error) {
	b := make([]byte, 1)
	b[0] = doc.Type()
	if doc.ID() == nil {
		return nil, ErrInvalidDocument
	}
	keys := doc.Keys()
	for _, key := range keys {
		// make key a Name type
		name := Name(key)
		marshaledName, err := MarshalObject(name)
		if err != nil {
			return nil, err
		}
		value, err := doc.Get(key)
		if err != nil {
			return nil, err
		}
		marshaledValue, err := MarshalObject(value)
		if err != nil {
			return nil, err
		}
		b = append(b, marshaledName...)
		b = append(b, marshaledValue...)
	}
	// add EOF at the end
	marshaledEOF, _ := MarshalObject(EOF{})
	b = append(b, marshaledEOF...)
	return b, nil
}

type Array []Object

func (a Array) Type() byte {
	return ArrayType
}

func (a Array) Value() interface{} {
	return a
}

func (a Array) String() string {
	// convert an array to string
	s := "["
	for i, value := range a {
		s += value.String()
		if i != len(a)-1 {
			s += ", "
		}
	}
	s += "]"
	return s
}

func (a Array) MarshalObject() ([]byte, error) {
	b := make([]byte, 1)
	b[0] = a.Type()
	for _, value := range a {
		marshaledValue, err := MarshalObject(value)
		if err != nil {
			return nil, err
		}
		b = append(b, marshaledValue...)
	}
	return b, nil
}

func (a *Array) UnmarshalObject(b []byte) (int, error) {
	if len(b) < 1 {
		return 0, ErrInvalidLength
	}
	if b[0] != a.Type() {
		return 0, ErrInvalidType
	}
	count := 1
	for {
		if len(b) < 1 {
			return 0, ErrInvalidLength
		}
		if b[count] == EOFType {
			break
		}
		value, n, err := UnmarshalObject(b[count:])
		if err != nil {
			return 0, err
		}
		count += n
		*a = append(*a, value)
	}
	return count, nil
}

func UnmarshalArray(b []byte) (Array, int, error) {
	var a Array
	count, err := a.UnmarshalObject(b)
	return a, count, err
}
