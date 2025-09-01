package reflect

import (
	"reflect"
)

func Implements(impl any, iface any) bool {
	return reflect.TypeOf(impl).Implements(reflect.TypeOf(iface).Elem())
}

func CallMethod(impl any, methodName string, args []reflect.Value) []reflect.Value {
	return reflect.ValueOf(impl).MethodByName(methodName).Call(args)
}
