package main

import (
	"bufio"
	"fmt"
	"image"
	"image/png"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// TODO:
// - CLI options
//  - dump file
//  - pass over pixels
//  - pass over data
//  - silent
//  - window size
//  - useful window title
//  - window position
//  - window border
//  - fullscreen
//  - no-rendering
//  - file-width
//  - file-height
//  - file-order (hilbert, horizontal, vertical, ???)
//  - scroll speed

var (
	turbo_srgb_bytes = [256][3]uint8{{48, 18, 59}, {50, 21, 67}, {51, 24, 74}, {52, 27, 81}, {53, 30, 88}, {54, 33, 95}, {55, 36, 102}, {56, 39, 109}, {57, 42, 115}, {58, 45, 121}, {59, 47, 128}, {60, 50, 134}, {61, 53, 139}, {62, 56, 145}, {63, 59, 151}, {63, 62, 156}, {64, 64, 162}, {65, 67, 167}, {65, 70, 172}, {66, 73, 177}, {66, 75, 181}, {67, 78, 186}, {68, 81, 191}, {68, 84, 195}, {68, 86, 199}, {69, 89, 203}, {69, 92, 207}, {69, 94, 211}, {70, 97, 214}, {70, 100, 218}, {70, 102, 221}, {70, 105, 224}, {70, 107, 227}, {71, 110, 230}, {71, 113, 233}, {71, 115, 235}, {71, 118, 238}, {71, 120, 240}, {71, 123, 242}, {70, 125, 244}, {70, 128, 246}, {70, 130, 248}, {70, 133, 250}, {70, 135, 251}, {69, 138, 252}, {69, 140, 253}, {68, 143, 254}, {67, 145, 254}, {66, 148, 255}, {65, 150, 255}, {64, 153, 255}, {62, 155, 254}, {61, 158, 254}, {59, 160, 253}, {58, 163, 252}, {56, 165, 251}, {55, 168, 250}, {53, 171, 248}, {51, 173, 247}, {49, 175, 245}, {47, 178, 244}, {46, 180, 242}, {44, 183, 240}, {42, 185, 238}, {40, 188, 235}, {39, 190, 233}, {37, 192, 231}, {35, 195, 228}, {34, 197, 226}, {32, 199, 223}, {31, 201, 221}, {30, 203, 218}, {28, 205, 216}, {27, 208, 213}, {26, 210, 210}, {26, 212, 208}, {25, 213, 205}, {24, 215, 202}, {24, 217, 200}, {24, 219, 197}, {24, 221, 194}, {24, 222, 192}, {24, 224, 189}, {25, 226, 187}, {25, 227, 185}, {26, 228, 182}, {28, 230, 180}, {29, 231, 178}, {31, 233, 175}, {32, 234, 172}, {34, 235, 170}, {37, 236, 167}, {39, 238, 164}, {42, 239, 161}, {44, 240, 158}, {47, 241, 155}, {50, 242, 152}, {53, 243, 148}, {56, 244, 145}, {60, 245, 142}, {63, 246, 138}, {67, 247, 135}, {70, 248, 132}, {74, 248, 128}, {78, 249, 125}, {82, 250, 122}, {85, 250, 118}, {89, 251, 115}, {93, 252, 111}, {97, 252, 108}, {101, 253, 105}, {105, 253, 102}, {109, 254, 98}, {113, 254, 95}, {117, 254, 92}, {121, 254, 89}, {125, 255, 86}, {128, 255, 83}, {132, 255, 81}, {136, 255, 78}, {139, 255, 75}, {143, 255, 73}, {146, 255, 71}, {150, 254, 68}, {153, 254, 66}, {156, 254, 64}, {159, 253, 63}, {161, 253, 61}, {164, 252, 60}, {167, 252, 58}, {169, 251, 57}, {172, 251, 56}, {175, 250, 55}, {177, 249, 54}, {180, 248, 54}, {183, 247, 53}, {185, 246, 53}, {188, 245, 52}, {190, 244, 52}, {193, 243, 52}, {195, 241, 52}, {198, 240, 52}, {200, 239, 52}, {203, 237, 52}, {205, 236, 52}, {208, 234, 52}, {210, 233, 53}, {212, 231, 53}, {215, 229, 53}, {217, 228, 54}, {219, 226, 54}, {221, 224, 55}, {223, 223, 55}, {225, 221, 55}, {227, 219, 56}, {229, 217, 56}, {231, 215, 57}, {233, 213, 57}, {235, 211, 57}, {236, 209, 58}, {238, 207, 58}, {239, 205, 58}, {241, 203, 58}, {242, 201, 58}, {244, 199, 58}, {245, 197, 58}, {246, 195, 58}, {247, 193, 58}, {248, 190, 57}, {249, 188, 57}, {250, 186, 57}, {251, 184, 56}, {251, 182, 55}, {252, 179, 54}, {252, 177, 54}, {253, 174, 53}, {253, 172, 52}, {254, 169, 51}, {254, 167, 50}, {254, 164, 49}, {254, 161, 48}, {254, 158, 47}, {254, 155, 45}, {254, 153, 44}, {254, 150, 43}, {254, 147, 42}, {254, 144, 41}, {253, 141, 39}, {253, 138, 38}, {252, 135, 37}, {252, 132, 35}, {251, 129, 34}, {251, 126, 33}, {250, 123, 31}, {249, 120, 30}, {249, 117, 29}, {248, 114, 28}, {247, 111, 26}, {246, 108, 25}, {245, 105, 24}, {244, 102, 23}, {243, 99, 21}, {242, 96, 20}, {241, 93, 19}, {240, 91, 18}, {239, 88, 17}, {237, 85, 16}, {236, 83, 15}, {235, 80, 14}, {234, 78, 13}, {232, 75, 12}, {231, 73, 12}, {229, 71, 11}, {228, 69, 10}, {226, 67, 10}, {225, 65, 9}, {223, 63, 8}, {221, 61, 8}, {220, 59, 7}, {218, 57, 7}, {216, 55, 6}, {214, 53, 6}, {212, 51, 5}, {210, 49, 5}, {208, 47, 5}, {206, 45, 4}, {204, 43, 4}, {202, 42, 4}, {200, 40, 3}, {197, 38, 3}, {195, 37, 3}, {193, 35, 2}, {190, 33, 2}, {188, 32, 2}, {185, 30, 2}, {183, 29, 2}, {180, 27, 1}, {178, 26, 1}, {175, 24, 1}, {172, 23, 1}, {169, 22, 1}, {167, 20, 1}, {164, 19, 1}, {161, 18, 1}, {158, 16, 1}, {155, 15, 1}, {152, 14, 1}, {149, 13, 1}, {146, 11, 1}, {142, 10, 1}, {139, 9, 2}, {136, 8, 2}, {133, 7, 2}, {129, 6, 2}, {126, 5, 2}, {122, 4, 3}}
	stride           int
	screenWidth      = 0
	screenHeight     = 0
	numReads         = 0
	scale            = 2
)

type Game struct {
	dump       []byte
	height     int
	count      int
	vx, vy     float64
	scratchPad [3]struct {
		img    *ebiten.Image
		opt    *ebiten.DrawImageOptions
		x, y   float64
		filled bool
	}
	refreshed bool
}

func clamp(val, min, max int) int {
	if val < min {
		return min
	} else if val > max {
		return max
	}
	return val
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func pad(a []byte, n int) []byte {
	if len(a) >= n {
		return a[:n]
	}
	return append(a, make([]byte, n-len(a))...)
}

func (g *Game) Update() error {
	// if g.height < g.scratchPad[0].img.Bounds().Dx()*g.scratchPad[0].img.Bounds().Dy()*4 {
	// 	g.height = 0
	if g.height < g.scratchPad[0].img.Bounds().Dx()*g.scratchPad[0].img.Bounds().Dy()*4 {
		g.InitFrames()
		// g.RefreshPixels(0)
		// return nil
	}
	// 	if err := g.InitFrames(); err != nil {
	// 		return err
	// 	}
	// }
	// return nil
	for i := 0; i < 3; i++ {
		if !g.HaveEnoughBuffer(i) || !g.scratchPad[i].filled {
			g.scratchPad[i].filled = g.RefreshPixels(i)
			// if g.scratchPad[i].y < 0 {
			// 	g.scratchPad[i].y = 0
			// }
		}
		g.scratchPad[i].y -= g.vy / ebiten.DefaultTPS
		g.scratchPad[i].x -= g.vx / ebiten.DefaultTPS
		if g.scratchPad[i].y+float64(g.scratchPad[i].img.Bounds().Dy()) < 0 {
			// g.count = i
			g.scratchPad[i].y += float64(2 * g.scratchPad[i].img.Bounds().Dy())
			// }
			// g.refreshed = true
			g.scratchPad[i].filled = g.RefreshPixels(i)
		}
		// g.refreshed = false
	}
	return nil
}

func (g *Game) HaveEnoughBuffer(i int) bool {
	return len(g.dump) > g.height+g.scratchPad[i].img.Bounds().Dx()*g.scratchPad[i].img.Bounds().Dy()*4
}

func (g *Game) RefreshPixels(i int) bool {
	if g.HaveEnoughBuffer(i) {
		g.scratchPad[i].img.ReplacePixels(
			g.dump[g.height : g.height+g.scratchPad[i].img.Bounds().Dx()*g.scratchPad[i].img.Bounds().Dy()*4],
		)
		g.height += g.scratchPad[i].img.Bounds().Dx() * g.scratchPad[i].img.Bounds().Dy() * 4
		return true
	} else {
		g.scratchPad[i].img.ReplacePixels(
			pad(g.dump[g.height:], (g.scratchPad[i].img.Bounds().Dx() * g.scratchPad[i].img.Bounds().Dy() * 4)),
		)
		return false
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	// count := g.count
	// renderFrames := func(num int) {
	for num := 0; num < 3; num++ {
		// if num == count {
		// 	continue
		// }
		g.scratchPad[num].opt = &ebiten.DrawImageOptions{}
		// g.scratchPad[num].opt.GeoM.Translate(-float64(screenWidth/scale/2), -float64(screenHeight/scale/2))
		g.scratchPad[num].opt.GeoM.Scale(float64(scale), float64(scale))
		g.scratchPad[num].opt.GeoM.Translate(g.scratchPad[num].x*float64(scale), g.scratchPad[num].y*float64(scale))
		screen.DrawImage(g.scratchPad[num].img, g.scratchPad[num].opt)
	}
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f, %s", ebiten.CurrentFPS(), fmt.Sprint(g.scratchPad[0].img.At(0, 0), g.scratchPad[0].img.Bounds(), "\n", g.scratchPad[0], "\n", g.scratchPad[1], "\n", g.scratchPad[2])))
}

func (g *Game) InitFrames() error {
	for i := 0; i < 3; i++ {
		for len(g.dump) == 0 {
			time.Sleep(time.Millisecond * 100)
		}
		g.scratchPad[i].filled = g.RefreshPixels(i)
		if g.scratchPad[i].filled == false {
			break
		}
	}
	return nil
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	screenWidth, screenHeight = ebiten.ScreenSizeInFullscreen()
	side := RoundToLowerPowerOfTwo(min(screenWidth, screenHeight))
	screenWidth = side
	// screenHeight = side
	stride = 4 * screenWidth
	g := &Game{
		dump:  make([]byte, 0, (screenHeight/scale)*(screenWidth/scale)*4),
		vx:    0,
		vy:    math.Exp2(float64(scale + 1)),
		count: 3,
	}
	for i := 0; i < 3; i++ {
		g.scratchPad[i].img = ebiten.NewImage(screenWidth/scale, screenHeight/scale)
		// g.scratchPad[i].x = 0
		g.scratchPad[i].y = float64(screenHeight / scale * i)
	}
	data := make([]byte, 4096)
	bufIO := bufio.NewReader(os.Stdin)
	go func(bufIO *bufio.Reader) {
		for {
			read, err := bufIO.Read(data)
			if err != nil {
				break
			}
			os.Stdout.Write(data[:read])
			clrs := make([]byte, read*4)
			for i := 0; i < read; i++ {
				clr := turbo_srgb_bytes[data[i]]
				clrs[4*i+0] = clr[0]
				clrs[4*i+1] = clr[1]
				clrs[4*i+2] = clr[2]
				clrs[4*i+3] = 255
			}
			g.dump = append(g.dump, clrs...)
		}
		time.Sleep(time.Second)
		side := int(math.Ceil(math.Sqrt(float64((len(g.dump)) / 4))))
		side = RoundToLowerPowerOfTwo(side)
		length := int(math.Ceil((float64(len(g.dump)) / 4) / float64(side)))
		bounds := image.Rect(0, 0, side, length)
		img := image.NewRGBA(bounds)
		pad := func(a []byte, n int) []byte {
			return append(a, make([]byte, n-len(a))...)
		}
		img.Pix = pad(g.dump, length*side*4)
		f, err := os.Create(strconv.FormatInt((time.Now().Unix()), 10) + "_pipeview.png")
		if err != nil {
			fmt.Println("Error creating image", err)
			return
		}
		defer f.Close()

		err = png.Encode(f, img)
		if err != nil {
			fmt.Println("Error encoding image", err)
			return
		}
		// os.Exit(0)
	}(bufIO)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Memory Viewer")
	ebiten.SetWindowPosition(0, 0)
	ebiten.SetWindowResizable(true)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func RoundToLowerPowerOfTwo(x int) int {
	return int(math.Pow(2, math.Floor(math.Log2(float64(x)))))
}
