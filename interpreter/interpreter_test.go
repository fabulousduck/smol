package interpreter

import "testing"

func TestStacksFind(T *testing.T) {
	var testVariableName = "someVariable"
	var testValue = "10"
	var testStack = Stacks{
		stack{
			{testVariableName, testValue},
		},
	}

	scopeLevel, index := testStack.find(testVariableName)
	if scopeLevel != 0 || index != 0 {
		T.Logf("\nTestStacksFind | failed to lookup variable on stack")
		T.Fail()
	}
}

func TestStacksSet(T *testing.T) {
	var testVariableName = "someVariable"
	var testValue = "10"
	var newTestValue = "11"
	var testStack = Stacks{
		stack{
			{testVariableName, testValue},
		},
	}

	testStack.set(0, 0, newTestValue)

	if testStack[0][0].value != newTestValue {
		T.Logf("\nTestStacksFind | failed to set value on stack")
		T.Fail()
	}
}

func TestStackContains(T *testing.T) {

}

func TestHeapFind(T *testing.T) {

}
