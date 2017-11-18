package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/edwardrf/snake"
)

var keyDirMap map[byte]snake.Dir

func init() {
	keyDirMap = make(map[byte]snake.Dir)
	keyDirMap['w'] = snake.DirUp
	keyDirMap['s'] = snake.DirDown
	keyDirMap['a'] = snake.DirLeft
	keyDirMap['d'] = snake.DirRight
	keyDirMap['k'] = snake.DirUp
	keyDirMap['j'] = snake.DirDown
	keyDirMap['h'] = snake.DirLeft
	keyDirMap['l'] = snake.DirRight
}

func main() {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	// restore the echoing state when exiting
	defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()

	input := make(chan byte, 100)

	go func() {
		var b []byte = make([]byte, 1)
		for {
			os.Stdin.Read(b)
			input <- b[0]
		}
	}()

	g := snake.New(40, 20)
	t := time.Tick(10 * time.Millisecond)

	dirs := make([]snake.Dir, 0, 100)
	cnt := 0
	lmt := 25
	for {
		select {
		case k := <-input:
			if k == 'q' {
				fmt.Printf("Exiting\n")
				return
			}
			d, ok := keyDirMap[k]
			if ok {
				dirs = append(dirs, d)
			} else {
				fmt.Printf("UNKNOWN KEY:%v", k)
			}
		case <-t:
			cnt++
			if cnt > lmt {
				dir := snake.DirNone
				if len(dirs) > 0 {
					dir = dirs[0]
					dirs = dirs[1:]
				}
				s := g.Step(dir)
				if s == snake.StatusLost {
					fmt.Println("You lost!")
					return
				}

				if s == snake.StatusWon {
					fmt.Println("You WON!")
					return
				}

				if s == snake.StatusAte && lmt > 6 {
					lmt--
				}
				fmt.Printf("\033[2J\033[1;1H%s", g.String())
				cnt = 0
			}
		}
	}
}
