package kv

import (
	"bytes"
	"context"

	"github.com/elliotcourant/noachis/pkg/datums"
	"github.com/elliotcourant/noachis/pkg/types"
	"github.com/pkg/errors"
)

type (
	RowEncoder struct {
		valueTypes []types.Type
	}
)

func NewRowEncoder(valueTypes []types.Type) RowEncoder {
	return RowEncoder{
		valueTypes: valueTypes,
	}
}

func (e RowEncoder) EncodeRow(
	ctx context.Context, values datums.Datums,
) ([]byte, error) {
	return EncodeRow(ctx, values, e.valueTypes)
}

func EncodeRow(
	ctx context.Context,
	values datums.Datums,
	valueTypes []types.Type,
) ([]byte, error) {
	buf := make([]byte, 0)
	buf = append(buf, datums.NewNullMap(values)...)
	for i, value := range values {
		encoded, err := value.Encode(ctx, valueTypes[i])
		if err != nil {
			return nil, err
		}

		buf = append(buf, encoded...)
	}

	return buf, nil
}

func DecodeRow(ctx context.Context, data []byte, valueTypes []types.Type) (datums.Datums, error) {
	buf := bytes.NewBuffer(data)

	nullMap, err := datums.DecodeNullMap(buf)
	if err != nil {
		return nil, err
	}

	result := make(datums.Datums, len(valueTypes), len(valueTypes))

	for i, valueType := range valueTypes {
		if nullMap.FieldIsNull(i) {
			result[i] = datums.Null
			continue
		}

		switch valueType.Family {
		case types.OIDFamily:
			result[i], err = datums.DecodeOid(ctx, buf, valueType)
		case types.DescriptorFamily:
		case types.TextFamily:
			result[i], err = datums.DecodeText(ctx, buf, valueType)
		case types.IntegerFamily:
		case types.BooleanFamily:
		case types.ArrayFamily:
		case types.UnknownFamily:
			fallthrough
		default:
			return nil, errors.Errorf("cannot decode type family %s", valueType.Family)
		}
		if err != nil {
			return nil, err
		}

		panic("not implemented")
	}

	return result, nil
}
