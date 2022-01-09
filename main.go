package main

import (
	_ "embed"
	"fmt"
	"math/rand"
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
	screenWidth  = 240
	screenHeight = 240
	tileSize     = 16
	tileXNum     = 25
	gridWidth    = 9
	gridHeight   = 9
)

var (
	mfont font.Face

	whiteTileImg  *ebiten.Image
	cursorTileImg *ebiten.Image
	redTileImg    *ebiten.Image
	gopherImg     *ebiten.Image
)

type Game struct {
	layers [][]int
	posX   int
	posY   int
}

// posX,posY follow below axis
// +---->Y
// |
// |
// v
// X

func init() {
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

	cursorTileImg = ebiten.NewImage(tileSize, tileSize)
	cursorTileImg.Fill(color.RGBA{50, 50, 50, 50})

	redTileImg = ebiten.NewImage(tileSize, tileSize)
	redTileImg.Fill(color.RGBA{255, 0, 0, 90})

	gopherImg, _, err = ebitenutil.NewImageFromFile("resources/gopher_front.png")
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error {
	// handle key press
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		if 0 < g.posY && g.layers[g.posX][g.posY-1] != 1 {
			g.posY--
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		if g.posY < gridHeight-1 && g.layers[g.posX][g.posY+1] != 1 {
			g.posY++
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		if 0 < g.posX && g.layers[g.posX-1][g.posY] != 1 {
			g.posX--
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		if g.posX < gridWidth-1 && g.layers[g.posX+1][g.posY] != 1 {
			g.posX++
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		// reset
		g.layers = newLayers()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("posX: %d, posY: %d", g.posX, g.posY))
	ebitenutil.DebugPrintAt(screen, "Move by WASD, Reset by R", 150, 0)
	const xNum = screenWidth / tileSize
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.1, 0.1)
	op.GeoM.Translate(float64(10+5*g.posY), float64(10+5*g.posX))
	screen.DrawImage(gopherImg, op)
	for i, row := range g.layers {
		for j, t := range row {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(0.7, 0.7)
			op.GeoM.Translate(float64(tileSize*j+80), float64(tileSize*i+80))
			if g.posX == i && g.posY == j {
				screen.DrawImage(cursorTileImg, op)

			} else {
				switch t {
				case 1:
					screen.DrawImage(redTileImg, op)
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

//generate layers
func newLayers() [][]int {
	layers := make([][]int, gridHeight)
	for i := 0; i < gridHeight; i++ {
		layers[i] = make([]int, gridWidth)
	}
	for i := 0; i < 10; i++ {
		ri := rand.Intn(gridHeight * gridWidth)
		layers[ri/gridWidth][ri%gridHeight] = 1
	}
	return layers
}
func sampleLayers() [][]int {
	return [][]int{
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

func main() {
	rand.Seed(time.Now().UnixNano())
	layers := newLayers()
	game := &Game{
		layers: layers,
	}
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("ice_puzzle")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
