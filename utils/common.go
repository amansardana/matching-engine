package utils

import (
	"fmt"
)

func UintToPaddedString(num uint64) string {
	return fmt.Sprintf("%019d", num)
}
