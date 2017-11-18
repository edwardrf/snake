package snake

import (
	"math/rand"
	"time"
)

type Dir int

const (
	DirNone Dir = iota
	DirUp
	DirRight
	DirDown
	DirLeft
)

type Status int

const (
	StatusPlaying = iota
	StatusAte
	StatusWon
	StatusLost
)

type vec2d struct {
	x, y int
}

type pos vec2d

type mov vec2d

type snake struct {
	body []pos
	dir  Dir
	grow int
}

type food struct {
	pos
	val int
}

type Game struct {
	width, height int
	snake         snake
	food          *food
	wall          []byte
}

var movMap map[Dir]mov

func init() {
	movMap = make(map[Dir]mov)
	movMap[DirUp] = mov{0, -1}
	movMap[DirRight] = mov{1, 0}
	movMap[DirDown] = mov{0, 1}
	movMap[DirLeft] = mov{-1, 0}
	rand.Seed(time.Now().Unix())
}

func New(w, h int) *Game {
	g := Game{width: w, height: h}
	g.wall = make([]byte, g.width*g.height)
	for i := 0; i < g.height; i++ {
		for j := 0; j < g.width; j++ {
			if i == 0 || j == 0 || i == g.height-1 || j == g.width-1 {
				g.wall[i*g.width+j] = '#'
			} else {
				g.wall[i*g.width+j] = ' '
			}
		}
	}
	g.initSnake(3)
	g.addFood(1)
	return &g
}

func (g *Game) Step(op Dir) Status {

	if op == DirNone {
		op = g.snake.dir
	} else {
		// Invalid op check, OpDown while moving in DirUp
		om := movMap[op]
		dm := movMap[g.snake.dir]
		if om.x+dm.x == 0 && om.y+dm.y == 0 {
			op = g.snake.dir
		}
	}

	g.snake.move(op)

	head := g.snake.head()

	// Wraping
	if head.x < 0 {
		head.x = g.width - 1
	}
	if head.x >= g.width {
		head.x = 0
	}
	if head.y < 0 {
		head.y = g.height - 1
	}
	if head.y >= g.height {
		head.y = 0
	}

	if *head == g.food.pos {
		g.snake.grow = g.food.val
		g.addFood(1)
		return StatusAte
	}

	if g.hasWall(*head) {
		return StatusLost
	}

	if g.snake.isCollided() {
		return StatusLost
	}

	if len(g.snake.body) == g.width*g.height {
		return StatusWon
	}

	return StatusPlaying

}

func (g *Game) addFood(val int) {
	var p pos
	for {
		p.x = rand.Intn(g.width)
		p.y = rand.Intn(g.height)

		if g.snake.isOnBody(p) {
			continue
		}

		if g.hasWall(p) {
			continue
		}
		break
	}
	g.food = &food{p, val}
}

func (g Game) hasWall(p pos) bool {
	return g.wall[p.y*g.width+p.x] == '#'
}

func (g Game) String() string {
	w := g.width + 2
	buf := make([]byte, w*g.height)
	for i := 0; i < g.height; i++ {
		copy(buf[i*w:i*w+g.width], g.wall[i*g.width:])
		buf[i*w+g.width] = '\r'
		buf[i*w+g.width+1] = '\n'
	}

	buf[g.food.y*w+g.food.x] = '*'
	for _, b := range g.snake.body {
		buf[b.y*w+b.x] = '@'
	}
	return string(buf)
}

func (g *Game) initSnake(l int) {
	g.snake.body = make([]pos, l)
	wgap := g.width / 5
	hgap := g.height / 5
	x := wgap + rand.Intn(wgap*3)
	y := hgap + rand.Intn(hgap*3)
	p := pos{x, y}

	dx, dy := 1, 0
	for i := 0; i < l; i++ {
		g.snake.body[i] = p

		np := pos{p.x + dx, p.y + dy}

		if g.hasWall(np) {
			dx = 0
			dy = 1
			np.x = p.x + dx
			np.y = p.y + dy
		}
		p = np
	}

	g.snake.dir = DirLeft
}

func (s snake) isOnBody(p pos) bool {
	for _, bp := range s.body {
		if bp.x == p.x && bp.y == p.y {
			return true
		}
	}
	return false
}

func (s snake) isCollided() bool {
	for i := 1; i < len(s.body); i++ {
		if s.body[0] == s.body[i] {
			return true
		}
	}
	return false
}

func (s *snake) move(d Dir) {
	m := movMap[d]
	head := pos{s.body[0].x + m.x, s.body[0].y + m.y}
	b := make([]pos, 0, len(s.body)+1)
	b = append(b, head)
	s.body = append(b, s.body...)
	if s.grow == 0 {
		s.body = s.body[:len(s.body)-1]
	} else {
		s.grow--
	}
	s.dir = d
}

func (s snake) head() *pos {
	return &s.body[0]
}
