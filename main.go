package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	tileWhite = iota
	tileStart
	tileGoal
	tileBlock
)

const (
	tileSize   = 16
	gridWidth  = 9
	gridHeight = 9
)

var (
	mfont font.Face

	whiteTileImg  *ebiten.Image
	cursorTileImg *ebiten.Image
	startTileImg  *ebiten.Image
	goalTileImg   *ebiten.Image
	redTileImg    *ebiten.Image
)

type layers = [][]int

type Puzzle struct {
	layers layers
	sx     int
	sy     int
	gx     int
	gy     int
}

type Game struct {
	posX   int
	posY   int
	puzzle Puzzle
}

// posX,posY follow below axis
// +---->Y
// |
// |
// v
// X

func init() {
	rand.Seed(time.Now().UnixNano())
	// setting font
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	mfont, err = opentype.NewFace(
		tt,
		&opentype.FaceOptions{Size: 12, DPI: 72, Hinting: font.HintingFull},
	)

	// load images
	whiteTileImg = ebiten.NewImage(tileSize, tileSize)
	whiteTileImg.Fill(color.White)

	startTileImg = ebiten.NewImage(tileSize, tileSize)
	startTileImg.Fill(color.RGBA{0, 0, 255, 50})

	goalTileImg = ebiten.NewImage(tileSize, tileSize)
	goalTileImg.Fill(color.RGBA{50, 255, 0, 50})

	redTileImg = ebiten.NewImage(tileSize, tileSize)
	redTileImg.Fill(color.RGBA{255, 0, 0, 90})

	cursorTileImg = ebiten.NewImage(tileSize, tileSize)
	cursorTileImg.Fill(color.RGBA{50, 50, 50, 50})

}

func (g *Game) moveLeft() bool {
	if 0 < g.posY && g.puzzle.layers[g.posX][g.posY-1] != tileBlock {
		g.posY--
		return true
	}
	return false
}
func (g *Game) moveRight() bool {
	if g.posY < gridHeight-1 && g.puzzle.layers[g.posX][g.posY+1] != tileBlock {
		g.posY++
		return true
	}
	return false
}
func (g *Game) moveUp() bool {
	if 0 < g.posX && g.puzzle.layers[g.posX-1][g.posY] != tileBlock {
		g.posX--
		return true
	}
	return false
}
func (g *Game) moveDown() bool {
	if g.posX < gridWidth-1 && g.puzzle.layers[g.posX+1][g.posY] != tileBlock {
		g.posX++
		return true
	}
	return false
}
func (g *Game) newStage() {
	g.puzzle = newPuzzle()
	g.posX = g.puzzle.sx
	g.posY = g.puzzle.sy
}

func (g *Game) Update() error {
	// handle key press
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.moveLeft()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.moveRight()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.moveUp()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.moveDown()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		// reset
		fmt.Println("reset")
		g.newStage()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyJ) {
		for g.moveDown() {
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		for g.moveLeft() {
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		for g.moveRight() {
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyK) {
		for g.moveUp() {
		}
	}
	// clear
	if g.posX == g.puzzle.gx && g.posY == g.puzzle.gy {
		fmt.Println("clear")
		g.newStage()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("posX: %d, posY: %d", g.posX, g.posY))
	ebitenutil.DebugPrintAt(screen, "Move by WASD, Warp by HJKL, Regenerate by R", 0, 30)
	for i, row := range g.puzzle.layers {
		for j, t := range row {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(0.7, 0.7)
			op.GeoM.Translate(float64(tileSize*j+80), float64(tileSize*i+80))
			if g.posX == i && g.posY == j {
				screen.DrawImage(cursorTileImg, op)

			} else {
				switch t {
				case tileBlock:
					screen.DrawImage(redTileImg, op)
				case tileStart:
					screen.DrawImage(startTileImg, op)
				case tileGoal:
					screen.DrawImage(goalTileImg, op)
				default:
					screen.DrawImage(whiteTileImg, op)
				}
			}

		}
	}
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

// generate layers not always reach goal
func newPuzzle() Puzzle {
	layers := make([][]int, gridHeight)
	ri := rand.Intn(gridHeight * gridWidth)
	sx, sy := ri/gridWidth, ri%gridWidth
	gri := rand.Intn(gridHeight * gridWidth)
	gx, gy := gri/gridWidth, gri%gridWidth
	for gx != sx && gy != sy {
		gri := rand.Intn(gridHeight * gridWidth)
		gx, gy = gri/gridWidth, gri%gridWidth
	}
	for i := 0; i < gridHeight; i++ {
		layers[i] = make([]int, gridWidth)
	}
	for i := 0; i < 10; i++ {
		ri := rand.Intn(gridHeight * gridWidth)
		rx, ry := ri/gridWidth, ri%gridWidth
		if (rx == sx && ry == sy) || (rx == gx && ry == gy) {
			continue
		}
		layers[rx][ry] = tileBlock
	}
	layers[sx][sy] = tileStart
	layers[gx][gy] = tileGoal

	return Puzzle{layers, sx, sy, gx, gy}
}

func sampleLayers() layers {
	return layers{
		{0, 0, 0, 0, 1, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 1, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 1, 0, 1, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
}

func getPuzzleFromFile(file string) Puzzle {
	fp, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	sc := bufio.NewScanner(fp)
	sc.Scan()
	hw := strings.Split(sc.Text(), " ")
	var h, w int
	h, err = strconv.Atoi(hw[0])
	w, err = strconv.Atoi(hw[1])
	if err != nil {
		panic(err)
	}
	var sx, sy, gx, gy int
	sc.Scan()
	sxsy := strings.Split(sc.Text(), " ")
	sc.Scan()
	gxgy := strings.Split(sc.Text(), " ")
	sx, err = strconv.Atoi(sxsy[0])
	sy, err = strconv.Atoi(sxsy[1])
	gx, err = strconv.Atoi(gxgy[0])
	gy, err = strconv.Atoi(gxgy[1])
	if err != nil {
		panic(err)
	}

	layers := make([][]int, h)
	for i := 0; i < h; i++ {
		layers[i] = make([]int, w)
	}
	layers[sx][sy] = tileStart
	layers[gx][gy] = tileGoal
	for i := 0; i < h; i++ {
		sc.Scan()
		row := sc.Text()
		for j, c := range row {
			switch c {
			case '.':
				continue
			case '#':
				layers[i][j] = tileBlock
			case 's':
				layers[i][j] = tileStart
			case 'g':
				layers[i][j] = tileGoal
			default:
				panic("unknown character appeared")
			}
		}
	}
	return Puzzle{layers, sx, sy, gx, gy}
}

func main() {
	puzzle := getPuzzleFromFile("resources/sample_layer.txt")
	// puzzle := newPuzzle()
	game := &Game{
		puzzle: puzzle,
	}
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("ice_puzzle")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
