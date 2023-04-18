package errshit

import (
	"errors"
	"fmt"
)

type X struct {
	Name string
}

func (x X) Error() string {
	return fmt.Sprintf("X Err: %s", x.Name)
}

func (x X) Shit() string {
	return "no " + x.Name
}

func Main() {
	var e1 error = X{"foo"}
	e2 := fmt.Errorf("Wrapper2  %w", e1)
	e3 := fmt.Errorf("Wrapper3  %w", e2)
	fmt.Println(e3)

	var eo error = e3
	for {
		tmp := errors.Unwrap(eo)
		if tmp == nil {
			break
		}
		eo = tmp
	}

	fmt.Println(eo)
	// if ee3, ok := eo.(interface{ Shit() string }); ok {
	// 	fmt.Println(ee3.Shit())
	// }
	var perr interface{ Shit() string }
	if errors.As(eo, &perr) {
		fmt.Println("-----", perr.Shit())
	}
}
