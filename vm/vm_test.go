package vm

import (
	"github.com/goby-lang/goby/compiler"
	"testing"
)

func TestVM_REPLExec(t *testing.T) {
	tests := []struct {
		inputs   []string
		expected interface{}
	}{
		{
			[]string{`

def foo
  123
end
`, `
foo
`},
			123},
		{
			[]string{
				`
class Foo
  def bar(x)
    x + 10
  end
end
`,
				`
Foo.new.bar(90)
`},
			100},
		{
			[]string{
				`
def foo
  123
end
`,
				`
foo
`,
				`
def foo
  345
end
`,

				`
foo
`,
			}, 345},
	}

	for i, test := range tests {
		v := initTestVM()
		v.InitForREPL()

		for _, input := range test.inputs {
			sets, err := compiler.CompileToInstructions(input)

			if err != nil {
				t.Fatalf(err.Error())
			}

			v.REPLExec(sets)
		}

		evaluated := v.GetExecResult()
		checkExpected(t, i, evaluated, test.expected)
		// Because REPL should maintain a base call frame so that the whole program won't exit
		v.checkCFP(t, i, 1)
	}
}

func initTestVM() *VM {
	return New("./", []string{})
}

func (v *VM) testEval(t *testing.T, input string) Object {
	iss, err := compiler.CompileToInstructions(input)

	if err != nil {
		t.Fatal(err.Error())
	}

	v.ExecInstructions(iss, "./")

	return v.mainThread.stack.top().Target
}

func (v *VM) checkCFP(t *testing.T, index, expectedCFP int) {
	if v.mainThread.cfp != expectedCFP {
		t.Fatalf("At case %d expect main thread's cfp to be %d. got: %d", index, expectedCFP, v.mainThread.cfp)
	}
}

func (v *VM) checkSP(t *testing.T, index, expectedSp int) {
	if v.mainThread.sp != expectedSp {
		t.Fatalf("At case %d expect main thread's sp to be %d. got: %d", index, expectedSp, v.mainThread.sp)
	}

}

func testIntegerObject(t *testing.T, i int, obj Object, expected int) bool {
	switch result := obj.(type) {
	case *IntegerObject:
		if result.Value != expected {
			t.Fatalf("At test case %d: object has wrong value. expect=%d, got=%d", i, expected, result.Value)
			return false
		}

		return true
	case *Error:
		t.Fatalf("At test case %d: %s", i, result.Message)
		return false
	default:
		t.Fatalf("At test case %d: object is not Integer. got=%T (%+v).", i, obj, obj)
		return false
	}
}

func testNullObject(t *testing.T, i int, obj Object) bool {
	switch result := obj.(type) {
	case *NullObject:
		return true
	case *Error:
		t.Fatalf("At test case %d: %s", i, result.Message)
		return false
	default:
		t.Fatalf("At test case %d: object is not NULL. got=%T (%+v)", i, obj, obj)
		return false
	}
}

func testStringObject(t *testing.T, i int, obj Object, expected string) bool {
	switch result := obj.(type) {
	case *StringObject:
		if result.Value != expected {
			t.Fatalf("At test case %d: object has wrong value. expect=%s, got=%s", i, expected, result.Value)
			return false
		}

		return true
	case *Error:
		t.Fatalf(result.Message)
		return false
	default:
		t.Fatalf("At test case %d: object is not String. got=%T (%+v).", i, obj, obj)
		return false
	}
}

func testBooleanObject(t *testing.T, i int, obj Object, expected bool) bool {
	switch result := obj.(type) {
	case *BooleanObject:
		if result.Value != expected {
			t.Fatalf("At test case %d: object has wrong value. expect=%d, got=%d", i, expected, result.Value)
			return false
		}

		return true
	case *Error:
		t.Fatalf(result.Message)
		return false
	default:
		t.Fatalf("At test case %d: object is not Boolean. got=%T (%+v).", i, obj, obj)
		return false
	}
}

func testArrayObject(t *testing.T, index int, obj Object, expected []interface{}) bool {
	result, ok := obj.(*ArrayObject)
	if !ok {
		t.Fatalf("At test case %d: object is not Array. got=%T (%+v)", index, obj, obj)
		return false
	}

	if len(result.Elements) != len(expected) {
		t.Fatalf("Don't equals length of array. Expect %d, got=%d", len(expected), len(result.Elements))
	}

	for i := 0; i < len(result.Elements); i++ {
		checkExpected(t, i, result.Elements[i], expected[i])
	}

	return true
}

func checkExpected(t *testing.T, i int, evaluated Object, expected interface{}) {
	if isError(evaluated) {
		t.Fatalf("At test case %d: %s", i, evaluated.toString())
		return
	}

	switch expected := expected.(type) {
	case int:
		testIntegerObject(t, i, evaluated, expected)
	case string:
		testStringObject(t, i, evaluated, expected)
	case bool:
		testBooleanObject(t, i, evaluated, expected)
	case nil:
		testNullObject(t, i, evaluated)
	default:
		t.Fatalf("Unknown type %T at case %d", expected, i)
	}
}

func isError(obj Object) bool {
	if obj != nil {
		_, ok := obj.(*Error)
		return ok
	}
	return false
}
