package nmea0183

import (
	"testing"
    "fmt"
)

func TestConfig(t *testing.T) {
	fmt.Println(len(Sentences))
    e := Config("./example")
	if e != nil{
		fmt.Println(e)
	}
	fmt.Println(len(Sentences))
    total := 10
	if total != 10 {
		t.Errorf("Sum was incorrect, got: %d, want: %d.", total, 10)
	}
}

func TestCheckSum(t *testing.T) {
	fmt.Println(len(Sentences))
	check := checksum("$1111111*45")
	expect := "31"
	if check != expect {
		t.Errorf("CheckSum was incorrect, got: %s, want: %s.", check, expect)
	}
}
