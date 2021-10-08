// ---------------------------------------------------------------------------
// LN-ELECTRONIC PROJECT LN-16C10R
// wwww.ln-electronic.com  ProjectManager : @Kimera
// FileData: 8/10/2564 7:37 2564  FileName : main.go
// ---------------------------------------------------------------------------

package main

import (
	"fmt"
	"github.com/lnelectronic/algorithmId/algorithmln"
)

func main() {
	ln_id := algorithmln.ID()

	fmt.Println("----------ID uint64-------------")
	fmt.Println(ln_id)
	fmt.Println("--------------------------------")
}
