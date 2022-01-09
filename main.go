package main

import (
	_ "embed"
	"fmt"

	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	screenWidth  = 240
	screenHeight = 240
	tileSize     = 16
	tileXNum     = 25
)

var (
	mfont font.Face

	tileImg   *ebiten.Image
	gopherImg *ebiten.Image
)

type Game struct {
	layers [][]int
	posX   int
	posY   int
}

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
	tileImg, _, err = ebitenutil.NewImageFromFile("resources/gopher_front.png")
	if err != nil {
		log.Fatal(err)
	}
	gopherImg, _, err = ebitenutil.NewImageFromFile("resources/gopher_front.png")
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.posX--
	} else if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.posX++
	} else if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.posY--
	} else if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.posY++
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	text.Draw(screen, fmt.Sprintf("posX: %d, posY: %d", g.posX, g.posY), mfont, 10, 20, color.White)
	const xNum = screenWidth / tileSize
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.1, 0.1)
	op.GeoM.Translate(50, 50)
	screen.DrawImage(gopherImg, op)
	// for _, row := range g.layers {
	// 	for i, t := range row {
	// 		op := &ebiten.DrawImageOptions{}
	// 		op.GeoM.Translate(float64(i%xNum)*tileSize, float64((i/xNum)*tileSize))
	// 		sx := (t % tileXNum) * tileSize
	// 		sy := (t / tileXNum) * tileSize
	// 		screen.DrawImage(tileImg.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
	// 	}

	// }
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	game := &Game{
		layers: [][]int{
			{3, 3, 3},
			{3, 3, 3},
			{3, 3, 3},
		},
	}
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("ice_puzzle")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
