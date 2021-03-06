package vm

import (
	"fmt"
)

// Object represents all objects in Goby, including Array, Integer or even Method and Error.
type Object interface {
	Class() *RClass
	toString() string
	toJSON() string
	instanceVariableGet(string) (Object, bool)
	instanceVariableSet(string, Object) Object
}

// Pointer is used to point to an object. Variables should hold pointer instead of holding a object directly.
type Pointer struct {
	Target Object
}

func (p *Pointer) returnClass() *RClass {
	return p.Target.(*RClass)
}

// RObject represents any non built-in class's instance.
type RObject struct {
	*baseObj
	InitializeMethod *MethodObject
}

type baseObj struct {
	class             *RClass
	InstanceVariables *environment
}

// Class will return object's class
func (b *baseObj) Class() *RClass {
	if b.class == nil {
		panic(fmt.Sprint("Object doesn't have class."))
	}
	return b.class
}

func (b *baseObj) instanceVariableGet(name string) (Object, bool) {
	v, ok := b.InstanceVariables.get(name)

	if !ok {
		return NULL, false
	}

	return v, true
}

func (b *baseObj) instanceVariableSet(name string, value Object) Object {
	b.InstanceVariables.set(name, value)

	return value
}

// Polymorphic helper functions -----------------------------------------

// toString tells which class it belongs to.
func (ro *RObject) toString() string {
	return "<Instance of: " + ro.class.Name + ">"
}

// toJSON converts the receiver into JSON string.
func (ro *RObject) toJSON() string {
	return ro.toString()
}

// Other helper functions -----------------------------------------------

func checkArgumentLen(args []Object, class *RClass, methodName string) *Error {
	if len(args) > 1 {
		return &Error{Message: fmt.Sprintf("Too many arguments for %s#%s", class.ReturnName(), methodName)}
	}

	return nil
}

func wrongTypeError(c *RClass) *Error {
	return &Error{Message: fmt.Sprintf("expect argument to be %s type", c.ReturnName())}
}
