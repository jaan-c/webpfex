package canvas

import "testing"

func TestInverse(t *testing.T) {
	var red uint16 = 0xaaaa
	var green uint16 = 0xbbbb
	var blue uint16 = 0xcccc
	var alpha uint16 = 0xdddd
	var value uint64 = 0xaaaabbbbccccdddd

	if v := MakeColor(value).Value(); v != value {
		t.Errorf("Expecting value to be %q, got %q", value, v)
	}

	color := MakeColorRgba(red, green, blue, alpha)
	if r := color.R(); r != red {
		t.Errorf("Expecting R to be %q, got %q", red, r)
	}
	if g := color.G(); g != green {
		t.Errorf("Expecting G to be %q, got %q", green, g)
	}
	if b := color.B(); b != blue {
		t.Errorf("Expecting B to be %q, got %q", blue, b)
	}
	if a := color.A(); a != alpha {
		t.Errorf("Expecting A to be %q, got %q", alpha, a)
	}

	r, g, b, a := color.Rgba()
	if r != red {
		t.Errorf("Expecting R to be %q, got %q", red, r)
	}
	if g != green {
		t.Errorf("Expecting G to be %q, got %q", green, g)
	}
	if b != blue {
		t.Errorf("Expecting B to be %q, got %q", blue, b)
	}
	if a != alpha {
		t.Errorf("Expecting A to be %q, got %q", alpha, a)
	}
}
