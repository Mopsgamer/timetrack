package main

import (
	"fmt"
	"os"

	"slices"

	"time"

	"github.com/gdamore/tcell/v2"
)

const help = "Welcome to timetrack!\n" +
	"This screen available through <?> and <h> keys.\n\n" +
	"Use ▲ and ▼ to move.\n" +
	"<C-c> - quit\n" +
	"<q>, <esc> - back\n" +
	"<a> - add new record\n" +
	"<A> - rename record\n" +
	"<r> - reset record time\n" +
	"<d> - delete current record\n" +
	"<C-f>, </> - search records\n" +
	"<D> - delete all found records"

func main() {
	if slices.Contains(os.Args, "--help") || slices.Contains(os.Args, "-h") {
		fmt.Println(help)
		os.Exit(1)
		return
	}
	screen, err := tcell.NewScreen()
	if err != nil {
		println(err)
		os.Exit(1)
		return
	}
	if err := screen.Init(); err != nil {
		println(err)
		os.Exit(1)
		return
	}
	defer screen.Fini()

	state := new(State)
	LoadState(state)
	redraw(screen, state)

	go func() {
		for {
			time.Sleep(50 * time.Millisecond)
			redraw(screen, state)
		}
	}()

	for {
		handleEvent(screen, state)
	}
}
