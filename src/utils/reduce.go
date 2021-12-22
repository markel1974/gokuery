/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package utils

import (
	"errors"
	"reflect"
)

func Reduce(source, initialValue, reducer interface{}) (interface{}, error) {
	srcV := reflect.ValueOf(source)
	kind := srcV.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return nil, errors.New("source value is not an array")
	}
	if reducer == nil {
		return nil, errors.New("reducer function cannot be nil")
	}
	rv := reflect.ValueOf(reducer)
	if rv.Kind() != reflect.Func {
		return nil, errors.New("reducer argument must be a function")
	}
	accumulator := initialValue
	accV := reflect.ValueOf(accumulator)
	for i := 0; i < srcV.Len(); i++ {
		entry := srcV.Index(i)
		reduceResults := rv.Call([]reflect.Value{
			accV,
			entry,
			reflect.ValueOf(i),
		})
		accV = reduceResults[0]
	}
	return accV.Interface(), nil
}
