// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package interpreter

import (
	"fmt"

	"github.com/google/cql/model"
	"github.com/google/cql/result"
	"github.com/google/cql/types"
)

// AGGREGATE FUNCTIONS - https://cql.hl7.org/09-b-cqlreference.html#aggregate-functions

// AllTrue(argument List<Boolean>) Boolean
// https://cql.hl7.org/09-b-cqlreference.html#alltrue
func (i *interpreter) evalAllTrue(m model.IUnaryExpression, operand result.Value) (result.Value, error) {
	if result.IsNull(operand) {
		return result.New(true)
	}
	l, err := result.ToSlice(operand)
	if err != nil {
		return result.Value{}, err
	}
	for _, elem := range l {
		if result.IsNull(elem) {
			continue
		}
		bv, err := result.ToBool(elem)
		if err != nil {
			return result.Value{}, err
		}
		if !bv {
			return result.New(false)
		}
	}
	return result.New(true)
}

// AnyTrue(argument List<Boolean>) Boolean
// https://cql.hl7.org/09-b-cqlreference.html#anytrue
func (i *interpreter) evalAnyTrue(m model.IUnaryExpression, operand result.Value) (result.Value, error) {
	if result.IsNull(operand) {
		return result.New(false)
	}
	l, err := result.ToSlice(operand)
	if err != nil {
		return result.Value{}, err
	}
	for _, elem := range l {
		if result.IsNull(elem) {
			continue
		}
		bv, err := result.ToBool(elem)
		if err != nil {
			return result.Value{}, err
		}
		if bv {
			return result.New(true)
		}
	}
	return result.New(false)
}

// Count(argument List<T>) Integer
// https://cql.hl7.org/09-b-cqlreference.html#count
func (i *interpreter) evalCount(m model.IUnaryExpression, operand result.Value) (result.Value, error) {
	if result.IsNull(operand) {
		return result.New(0)
	}
	l, err := result.ToSlice(operand)
	if err != nil {
		return result.Value{}, err
	}
	count := 0
	for _, elem := range l {
		if !result.IsNull(elem) {
			count++
		}
	}
	return result.New(count)
}

// Sum(argument List<Decimal>) Decimal
// Sum(argument List<Integer>) Integer
// Sum(argument List<Long>) Long
// Sum(argument List<Quantity>) Quantity
// https://cql.hl7.org/09-b-cqlreference.html#sum
func (i *interpreter) evalSum(m model.IUnaryExpression, operand result.Value) (result.Value, error) {
	if result.IsNull(operand) {
		return result.New(nil)
	}
	l, err := result.ToSlice(operand)
	if err != nil {
		return result.Value{}, err
	}
	lType, ok := operand.RuntimeType().(*types.List)
	if !ok {
		return result.Value{}, fmt.Errorf("Sum(%v) operand is not a list", m.GetName())
	}
	switch lType.ElementType {
	case types.Any:
		// Special case for handling lists that contain only null runtime values.
		return result.New(nil)
	case types.Decimal:
		var sum float64
		var foundValue bool
		for _, elem := range l {
			if result.IsNull(elem) {
				continue
			}
			foundValue = true
			v, err := result.ToFloat64(elem)
			if err != nil {
				return result.Value{}, err
			}
			sum += v
		}
		if !foundValue {
			return result.New(nil)
		}
		return result.New(sum)
	case types.Integer:
		var sum int32
		var foundValue bool
		for _, elem := range l {
			if result.IsNull(elem) {
				continue
			}
			foundValue = true
			v, err := result.ToInt32(elem)
			if err != nil {
				return result.Value{}, err
			}
			sum += v
		}
		if !foundValue {
			return result.New(nil)
		}
		return result.New(sum)
	case types.Long:
		var sum int64
		var foundValue bool
		for _, elem := range l {
			if result.IsNull(elem) {
				continue
			}
			foundValue = true
			v, err := result.ToInt64(elem)
			if err != nil {
				return result.Value{}, err
			}
			sum += v
		}
		if !foundValue {
			return result.New(nil)
		}
		return result.New(sum)
	case types.Quantity:
		var sum result.Quantity
		var foundValue bool
		for _, elem := range l {
			if result.IsNull(elem) {
				continue
			}
			v, err := result.ToQuantity(elem)
			if err != nil {
				return result.Value{}, err
			}
			if !foundValue {
				foundValue = true
				sum = result.Quantity{Value: 0, Unit: v.Unit}
			}
			if sum.Unit != v.Unit {
				return result.Value{}, fmt.Errorf("Sum(%v) got List of Quantity values with different units which is not supported, got %v and %v", m.GetName(), sum.Unit, v.Unit)
			}
			sum.Value += v.Value
		}
		if !foundValue {
			return result.New(nil)
		}
		return result.New(sum)
	default:
		return result.Value{}, fmt.Errorf("Sum(%v) operand is not a list of Decimal, Integer, Long, or Quantity", m.GetName())
	}
}
