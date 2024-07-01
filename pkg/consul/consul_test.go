package consul

import (
	"fmt"
	"strings"
	"testing"
)

func Test_figureOutListenOn(t *testing.T) {
	fmt.Println(len(strings.Split(":00", ":")))
	fmt.Println(len(strings.Split("", ":")))

}
