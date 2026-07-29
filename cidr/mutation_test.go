package cidr

import (
	"math/big"
	"net"
	"testing"
)

// SubnetBig and HostBig take a caller-owned *big.Int and must not modify it.

func TestSubnetBigDoesNotMutateNum(t *testing.T) {
	cases := []struct {
		Base   string
		Bits   int
		Num    int64
		Output string
	}{
		{"10.3.0.0/16", 8, 5, "10.3.5.0/24"},
		{"192.168.0.0/21", 4, 7, "192.168.3.128/25"},
		{"2001:db8::/32", 16, 7, "2001:db8:7::/48"},
		{"fe80::/49", 16, 7, "fe80:0:0:3:8000::/65"},
	}

	for _, tc := range cases {
		_, base, _ := net.ParseCIDR(tc.Base)
		num := big.NewInt(tc.Num)
		got, err := SubnetBig(base, tc.Bits, num)
		if err != nil {
			t.Fatalf("SubnetBig(%s, %d, %d) failed: %s", tc.Base, tc.Bits, tc.Num, err)
		}
		if got.String() != tc.Output {
			t.Errorf("SubnetBig(%s, %d, %d) = %s; want %s", tc.Base, tc.Bits, tc.Num, got, tc.Output)
		}
		if num.Int64() != tc.Num {
			t.Errorf("SubnetBig(%s, %d, %d) mutated num to %s", tc.Base, tc.Bits, tc.Num, num)
		}
	}
}

// Reusing one *big.Int across calls must be stable: before the fix the second
// call saw the shifted value and returned an out-of-range error.
func TestSubnetBigRepeatedCalls(t *testing.T) {
	cases := []struct {
		Base string
		Bits int
		Num  int64
	}{
		{"10.3.0.0/16", 8, 5},
		{"2001:db8::/32", 16, 7},
	}

	for _, tc := range cases {
		_, base, _ := net.ParseCIDR(tc.Base)
		num := big.NewInt(tc.Num)
		first, err := SubnetBig(base, tc.Bits, num)
		if err != nil {
			t.Fatalf("SubnetBig(%s, %d, %d) failed: %s", tc.Base, tc.Bits, tc.Num, err)
		}
		second, err := SubnetBig(base, tc.Bits, num)
		if err != nil {
			t.Fatalf("second SubnetBig(%s, %d, %d) failed: %s", tc.Base, tc.Bits, tc.Num, err)
		}
		if first.String() != second.String() {
			t.Errorf("SubnetBig(%s, %d, %d) not idempotent: %s then %s", tc.Base, tc.Bits, tc.Num, first, second)
		}
	}
}

func TestHostBigDoesNotMutateNum(t *testing.T) {
	cases := []struct {
		Base string
		Num  int64
	}{
		{"10.3.0.0/16", 5},
		{"10.3.0.0/16", -2},
		{"2001:db8::/32", 5},
		{"2001:db8::/32", -2},
	}

	for _, tc := range cases {
		_, base, _ := net.ParseCIDR(tc.Base)
		num := big.NewInt(tc.Num)
		if _, err := HostBig(base, num); err != nil {
			t.Fatalf("HostBig(%s, %d) failed: %s", tc.Base, tc.Num, err)
		}
		if num.Int64() != tc.Num {
			t.Errorf("HostBig(%s, %d) mutated num to %s", tc.Base, tc.Num, num)
		}
	}
}
