package assert

import "testing"

func TestStringEqual(t *testing.T) {
	var x string = "x"
	var y string = "x"

	StringEqual("test", x, y, t)
}

func TestIntEqual(t *testing.T) {
	var x int = 123
	var y int = 123

	IntEqual("test", x, y, t)
}

func TestByteEqual(t *testing.T) {
	var x int = 123
	var y int = 123

	IntEqual("test", x, y, t)
}

func TestInt32Equal(t *testing.T) {
	var x int32 = 123
	var y int32 = 123

	Int32Equal("test", x, y, t)
}

func TestInt64Equal(t *testing.T) {
	var x int64 = 123
	var y int64 = 123

	Int64Equal("test", x, y, t)
}

func TestFloat32Equal(t *testing.T) {
	var x float32 = 1.23
	var y float32 = 1.23

	Float32Equal("test", x, y, t)
}

func TestFloat64Equal(t *testing.T) {
	var x float64 = 1.23
	var y float64 = 1.23

	Float64Equal("test", x, y, t)
}

func TestComplex64Equal(t *testing.T) {
	var x complex64 = 1.23
	var y complex64 = 1.23

	Complex64Equal("test", x, y, t)
}

func TestComplex128Equal(t *testing.T) {
	var x complex128 = 1.23
	var y complex128 = 1.23

	Complex128Equal("test", x, y, t)
}

func TestTrue(t *testing.T) {
	True("test", true, t)
}

func TestFalse(t *testing.T) {
	False("test", false, t)
}

func TestThat(t *testing.T) {
	That("test", "foo", func() bool { return true }, t)
}
