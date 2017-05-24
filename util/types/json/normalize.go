// Copyright 2017 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package json

import (
	"encoding/json"
	"unsafe"
)

func normalize(in interface{}) (j JSON) {
	switch t := in.(type) {
	case nil:
		j.typeCode = typeCodeLiteral
		j.bit64 = int64(jsonLiteralNil)
	case bool:
		j.typeCode = typeCodeLiteral
		if t {
			j.bit64 = int64(jsonLiteralTrue)
		} else {
			j.bit64 = int64(jsonLiteralFalse)
		}
	case int64:
		j.typeCode = typeCodeInt64
		j.bit64 = t
	case float64:
		j.typeCode = typeCodeFloat64
		*(*float64)(unsafe.Pointer(&j.bit64)) = t
	case json.Number:
		if i64, err := t.Int64(); err == nil {
			j.typeCode = typeCodeInt64
			j.bit64 = i64
		} else {
			f64, _ := t.Float64()
			j.typeCode = typeCodeFloat64
			*(*float64)(unsafe.Pointer(&j.bit64)) = f64
		}
	case string:
		j.typeCode = typeCodeString
		j.str = t
	case map[string]interface{}:
		j.typeCode = typeCodeObject
		j.object = make(map[string]JSON, len(t))
		for key, value := range t {
			j.object[key] = normalize(value)
		}
	case []interface{}:
		j.typeCode = typeCodeArray
		j.array = make([]JSON, 0, len(t))
		for _, elem := range t {
			j.array = append(j.array, normalize(elem))
		}
	default:
		panic(internalErrorUnknownTypeCode)
	}
	return
}
