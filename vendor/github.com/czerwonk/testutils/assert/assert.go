package assert

import "testing"

const equalError string = "%s: expected %v but got %v"

type Assertion func() bool

func StringEqual(name string, expected, current string, t *testing.T) {
	if current != expected {
		t.Fatalf(equalError, name, expected, current)
	}
}

func IntEqual(name string, expected, current int, t *testing.T) {
	if current != expected {
		t.Fatalf(equalError, name, expected, current)
	}
}

func ByteEqual(name string, expected, current byte, t *testing.T) {
	if current != expected {
		t.Fatalf(equalError, name, expected, current)
	}
}

func Int32Equal(name string, expected, current int32, t *testing.T) {
	if current != expected {
		t.Fatalf(equalError, name, expected, current)
	}
}

func Int64Equal(name string, expected, current int64, t *testing.T) {
	if current != expected {
		t.Fatalf(equalError, name, expected, current)
	}
}

func Float32Equal(name string, expected, current float32, t *testing.T) {
	if current != expected {
		t.Fatalf(equalError, name, expected, current)
	}
}

func Float64Equal(name string, expected, current float64, t *testing.T) {
	if current != expected {
		t.Fatalf(equalError, name, expected, current)
	}
}

func Complex64Equal(name string, expected, current complex64, t *testing.T) {
	if current != expected {
		t.Fatalf(equalError, name, expected, current)
	}
}

func Complex128Equal(name string, expected, current complex128, t *testing.T) {
	if current != expected {
		t.Fatalf(equalError, name, expected, current)
	}
}

func True(name string, current bool, t *testing.T) {
	if !current {
		t.Fatalf(equalError, name, true, false)
	}
}

func False(name string, current bool, t *testing.T) {
	if current {
		t.Fatalf(equalError, name, false, true)
	}
}

func That(name, text string, assertion Assertion, t *testing.T) {
	if !assertion() {
		t.Fatalf("%s: %s", name, text)
	}
}
