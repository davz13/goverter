package builder

import (
	"go/types"

	"github.com/dave/jennifer/jen"
)

type Type struct {
	T            types.Type
	Struct       bool
	StructType   *types.Struct
	Named        bool
	NamedType    *types.Named
	Pointer      bool
	PointerType  *types.Pointer
	PointerInner *Type
	List         bool
	ListFixed    bool
	ListInner    *Type
	Map          bool
	MapType      *types.Map
	MapKey       *Type
	MapValue     *Type
	Basic        bool
	BasicType    *types.Basic
}

type JenID = *jen.Statement

func TypeOf(t types.Type) *Type {
	rt := &Type{}
	rt.T = t
	switch value := t.(type) {
	case *types.Pointer:
		rt.Pointer = true
		rt.PointerType = value
		rt.PointerInner = TypeOf(value.Elem())
	case *types.Basic:
		rt.Basic = true
		rt.BasicType = value
	case *types.Map:
		rt.Map = true
		rt.MapType = value
		rt.MapKey = TypeOf(value.Key())
		rt.MapValue = TypeOf(value.Elem())
	case *types.Slice:
		rt.List = true
		rt.ListInner = TypeOf(value.Elem())
	case *types.Array:
		rt.List = true
		rt.ListFixed = true
		rt.ListInner = TypeOf(value.Elem())
	case *types.Named:
		underlying := TypeOf(value.Underlying())
		underlying.T = value
		underlying.Named = true
		underlying.NamedType = value
		return underlying
	case *types.Struct:
		rt.Struct = true
		rt.StructType = value
	default:
		panic("unknown types.Type " + t.String())
	}
	return rt
}

func (t Type) TypeAsJen() *jen.Statement {
	if t.Named {
		return toCode(t.NamedType, &jen.Statement{})
	}
	return toCode(t.T, &jen.Statement{})
}

func toCode(t types.Type, st *jen.Statement) *jen.Statement {
	switch cast := t.(type) {
	case *types.Named:
		return st.Qual(cast.Obj().Pkg().Path(), cast.Obj().Name())
	case *types.Map:
		key := toCode(cast.Key(), &jen.Statement{})
		return toCode(cast.Elem(), st.Map(key))
	case *types.Slice:
		return toCode(cast.Elem(), st.Index())
	case *types.Array:
		return toCode(cast.Elem(), st.Index(jen.Lit(int(cast.Len()))))
	case *types.Pointer:
		return toCode(cast.Elem(), st.Op("*"))
	case *types.Basic:
		switch cast.Kind() {
		case types.String:
			return st.String()
		case types.Int:
			return st.Int()
		case types.Int8:
			return st.Int8()
		case types.Int16:
			return st.Int16()
		case types.Int32:
			return st.Int32()
		case types.Int64:
			return st.Int64()
		case types.Uint:
			return st.Uint()
		case types.Uint8:
			return st.Uint8()
		case types.Uint16:
			return st.Uint16()
		case types.Uint32:
			return st.Uint32()
		case types.Uint64:
			return st.Uint64()
		case types.Bool:
			return st.Bool()
		case types.Complex128:
			return st.Complex128()
		case types.Complex64:
			return st.Complex64()
		case types.Float32:
			return st.Float32()
		case types.Float64:
			return st.Float64()
		}
	}
	panic("unsupported type " + t.String())
}