package reflect

import "reflect"

// CopyType returns a pointer to newly allocated instance of
// the type of obj, obj is expected to be a pointer
func CopyType(obj interface{}) interface{} {
	objPVal := reflect.ValueOf(obj)
	objVal := reflect.Indirect(objPVal)

	objType := objVal.Type()
	newObjVal := reflect.New(objType)
	return newObjVal.Interface()
}
