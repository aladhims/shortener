package base

import "testing"

func TestReverse(t *testing.T) {
	origin := "a"
	expected := "a"

	if actual := reverse(origin); actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}

	origin = "test1streverse"
	expected = "esreverts1tset"

	if actual := reverse(origin); actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestEncode(t *testing.T) {
	var origin uint64
	origin = 0
	expected := "A"

	if actual := Encode(origin); actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}

	origin = 26
	expected = "a"

	if actual := Encode(origin); actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}

	origin = 7912
	expected = "CDm"

	if actual := Encode(origin); actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}

	origin = 12313923077881
	expected = "DexMSUCd"

	if actual := Encode(origin); actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}

}

func TestDecode(t *testing.T) {
	var expected uint64
	expected = 0
	encoded := "A"

	if actual := Decode(encoded); actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}

	expected = 26
	encoded = "a"

	if actual := Decode(encoded); actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}

	expected = 7912
	encoded = "CDm"

	if actual := Decode(encoded); actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}

	expected = 12313923077881
	encoded = "DexMSUCd"

	if actual := Decode(encoded); actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}
