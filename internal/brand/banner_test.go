package brand

import (
	"os"
	"testing"
)

func TestRootBannerMatchesEmbeddedBanner(t *testing.T) {
	rootBanner, err := os.ReadFile("../../HHHQ")
	if err != nil {
		t.Fatal(err)
	}
	if string(rootBanner) != Banner() {
		t.Fatal("root HHHQ and embedded brand banner differ")
	}
}
