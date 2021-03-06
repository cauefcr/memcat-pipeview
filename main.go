package main

import (
	"fmt"
	"image/png"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/mozilla/masche/memaccess"
	"github.com/mozilla/masche/process"
)

var stride = 256

var turbo_srgb_bytes = [256][3]uint8{{48, 18, 59}, {50, 21, 67}, {51, 24, 74}, {52, 27, 81}, {53, 30, 88}, {54, 33, 95}, {55, 36, 102}, {56, 39, 109}, {57, 42, 115}, {58, 45, 121}, {59, 47, 128}, {60, 50, 134}, {61, 53, 139}, {62, 56, 145}, {63, 59, 151}, {63, 62, 156}, {64, 64, 162}, {65, 67, 167}, {65, 70, 172}, {66, 73, 177}, {66, 75, 181}, {67, 78, 186}, {68, 81, 191}, {68, 84, 195}, {68, 86, 199}, {69, 89, 203}, {69, 92, 207}, {69, 94, 211}, {70, 97, 214}, {70, 100, 218}, {70, 102, 221}, {70, 105, 224}, {70, 107, 227}, {71, 110, 230}, {71, 113, 233}, {71, 115, 235}, {71, 118, 238}, {71, 120, 240}, {71, 123, 242}, {70, 125, 244}, {70, 128, 246}, {70, 130, 248}, {70, 133, 250}, {70, 135, 251}, {69, 138, 252}, {69, 140, 253}, {68, 143, 254}, {67, 145, 254}, {66, 148, 255}, {65, 150, 255}, {64, 153, 255}, {62, 155, 254}, {61, 158, 254}, {59, 160, 253}, {58, 163, 252}, {56, 165, 251}, {55, 168, 250}, {53, 171, 248}, {51, 173, 247}, {49, 175, 245}, {47, 178, 244}, {46, 180, 242}, {44, 183, 240}, {42, 185, 238}, {40, 188, 235}, {39, 190, 233}, {37, 192, 231}, {35, 195, 228}, {34, 197, 226}, {32, 199, 223}, {31, 201, 221}, {30, 203, 218}, {28, 205, 216}, {27, 208, 213}, {26, 210, 210}, {26, 212, 208}, {25, 213, 205}, {24, 215, 202}, {24, 217, 200}, {24, 219, 197}, {24, 221, 194}, {24, 222, 192}, {24, 224, 189}, {25, 226, 187}, {25, 227, 185}, {26, 228, 182}, {28, 230, 180}, {29, 231, 178}, {31, 233, 175}, {32, 234, 172}, {34, 235, 170}, {37, 236, 167}, {39, 238, 164}, {42, 239, 161}, {44, 240, 158}, {47, 241, 155}, {50, 242, 152}, {53, 243, 148}, {56, 244, 145}, {60, 245, 142}, {63, 246, 138}, {67, 247, 135}, {70, 248, 132}, {74, 248, 128}, {78, 249, 125}, {82, 250, 122}, {85, 250, 118}, {89, 251, 115}, {93, 252, 111}, {97, 252, 108}, {101, 253, 105}, {105, 253, 102}, {109, 254, 98}, {113, 254, 95}, {117, 254, 92}, {121, 254, 89}, {125, 255, 86}, {128, 255, 83}, {132, 255, 81}, {136, 255, 78}, {139, 255, 75}, {143, 255, 73}, {146, 255, 71}, {150, 254, 68}, {153, 254, 66}, {156, 254, 64}, {159, 253, 63}, {161, 253, 61}, {164, 252, 60}, {167, 252, 58}, {169, 251, 57}, {172, 251, 56}, {175, 250, 55}, {177, 249, 54}, {180, 248, 54}, {183, 247, 53}, {185, 246, 53}, {188, 245, 52}, {190, 244, 52}, {193, 243, 52}, {195, 241, 52}, {198, 240, 52}, {200, 239, 52}, {203, 237, 52}, {205, 236, 52}, {208, 234, 52}, {210, 233, 53}, {212, 231, 53}, {215, 229, 53}, {217, 228, 54}, {219, 226, 54}, {221, 224, 55}, {223, 223, 55}, {225, 221, 55}, {227, 219, 56}, {229, 217, 56}, {231, 215, 57}, {233, 213, 57}, {235, 211, 57}, {236, 209, 58}, {238, 207, 58}, {239, 205, 58}, {241, 203, 58}, {242, 201, 58}, {244, 199, 58}, {245, 197, 58}, {246, 195, 58}, {247, 193, 58}, {248, 190, 57}, {249, 188, 57}, {250, 186, 57}, {251, 184, 56}, {251, 182, 55}, {252, 179, 54}, {252, 177, 54}, {253, 174, 53}, {253, 172, 52}, {254, 169, 51}, {254, 167, 50}, {254, 164, 49}, {254, 161, 48}, {254, 158, 47}, {254, 155, 45}, {254, 153, 44}, {254, 150, 43}, {254, 147, 42}, {254, 144, 41}, {253, 141, 39}, {253, 138, 38}, {252, 135, 37}, {252, 132, 35}, {251, 129, 34}, {251, 126, 33}, {250, 123, 31}, {249, 120, 30}, {249, 117, 29}, {248, 114, 28}, {247, 111, 26}, {246, 108, 25}, {245, 105, 24}, {244, 102, 23}, {243, 99, 21}, {242, 96, 20}, {241, 93, 19}, {240, 91, 18}, {239, 88, 17}, {237, 85, 16}, {236, 83, 15}, {235, 80, 14}, {234, 78, 13}, {232, 75, 12}, {231, 73, 12}, {229, 71, 11}, {228, 69, 10}, {226, 67, 10}, {225, 65, 9}, {223, 63, 8}, {221, 61, 8}, {220, 59, 7}, {218, 57, 7}, {216, 55, 6}, {214, 53, 6}, {212, 51, 5}, {210, 49, 5}, {208, 47, 5}, {206, 45, 4}, {204, 43, 4}, {202, 42, 4}, {200, 40, 3}, {197, 38, 3}, {195, 37, 3}, {193, 35, 2}, {190, 33, 2}, {188, 32, 2}, {185, 30, 2}, {183, 29, 2}, {180, 27, 1}, {178, 26, 1}, {175, 24, 1}, {172, 23, 1}, {169, 22, 1}, {167, 20, 1}, {164, 19, 1}, {161, 18, 1}, {158, 16, 1}, {155, 15, 1}, {152, 14, 1}, {149, 13, 1}, {146, 11, 1}, {142, 10, 1}, {139, 9, 2}, {136, 8, 2}, {133, 7, 2}, {129, 6, 2}, {126, 5, 2}, {122, 4, 3}}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func stdCheck(softs []error, hard error) {
	for _, soft := range softs {
		check(soft)
	}
	check(hard)
}

func muxErr(softs []error, hard error) error {
	for _, soft := range softs {
		if soft != nil {
			return soft
		}
	}
	return hard
}

func availableProgs(quiet bool) (names []string, pids []uint) {
	allpids, softs, hard := process.GetAllPids()
	names = make([]string, 0, len(allpids))
	pids = make([]uint, 0, len(allpids))
	stdCheck(softs, hard)
	for _, pid := range allpids {
		if !quiet {
			fmt.Printf("%d\n", pid)
		}
		proc, softs, hard := process.OpenFromPid(pid)
		err := muxErr(softs, hard)
		if err != nil {
			if !quiet {
				fmt.Printf("%s\n", err)
			}
			continue
		}
		name, softs, hard := proc.Name()
		err = muxErr(softs, hard)
		if err != nil {
			if !quiet {
				fmt.Printf("%s\n", err)
			}
			continue
		}
		if !quiet {
			fmt.Println(name)
		}
		names = append(names, name)
		pids = append(pids, pid)
		proc.Close()
	}
	return
}

type WrappedProcess struct {
	process.Process
}

func (p WrappedProcess) DumpMem() (memoryScan chan []byte, err error) {
	memoryScan = make(chan []byte, 1000)
	var last uintptr = 0
	go func() {
		count := 0
		err = muxErr(memaccess.WalkMemory(p, last, 4096<<2, func(addr uintptr, buf []byte) bool {
			if len(buf) == 0 || last == addr {
				return false
			}
			last = addr
			memoryScan <- append([]byte{}, buf...)
			count++
			if count%1000 == 0 {
				fmt.Println(addr)
			}
			return true
		}))
		for len(memoryScan) > 0 {
			time.Sleep(100 * time.Millisecond)
		}
		close(memoryScan)
	}()
	return
}

var (
	screenWidth  = 0
	screenHeight = 0
)

type Game struct {
	dump []byte
}

func (g *Game) Update() error {
	return nil
}

var height = int(0)

func clamp(val, min, max int) int {
	if val < min {
		return min
	} else if val > max {
		return max
	}
	return val
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	if len(g.dump) < screenHeight*screenWidth {
		height = screenHeight * screenWidth
		return
	} else {
		screen.ReplacePixels(g.dump[height : height+screenHeight*screenWidth*4])
	}
	height = clamp(height+screenHeight*screenWidth*4/60, 0, len(g.dump)-screenHeight*screenWidth*4)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ebiten.ScreenSizeInFullscreen()
}

func main() {
	names, pids := availableProgs(false)
	chosen := pids[len(pids)-1]
	for i, name := range names {
		if strings.Contains(name, os.Args[1]) {
			chosen = pids[i]
		}
	}
	p, softs, hard := process.OpenFromPid(chosen)
	if muxErr(softs, hard) != nil {
		fmt.Println("Error opening process")
		return
	}
	defer p.Close()
	wp := WrappedProcess{p}
	mem, err := wp.DumpMem()
	if err != nil {
		fmt.Println("Error dumping memory")
		return
	}
	screenWidth, screenHeight = ebiten.ScreenSizeInFullscreen()
	stride = screenHeight * screenWidth / 4
	g := &Game{
		dump: []byte{},
	}
	go func(mem chan []byte, g *Game) {
		for data := range mem {
			clrs := make([]byte, len(data)*4)
			for i := range data {
				clr := turbo_srgb_bytes[data[i]]
				clrs[4*i+0] = clr[0]
				clrs[4*i+1] = clr[1]
				clrs[4*i+2] = clr[2]
				clrs[4*i+3] = 255
			}
			g.dump = append(g.dump, clrs...)
		}
		// TODO make an image of the dump
		height := int(math.Ceil(float64(len(g.dump)/4)/float64(screenWidth)) * float64(screenWidth))
		image := ebiten.NewImage(screenWidth, height)
		g.dump = append(g.dump, make([]byte, height*screenWidth*4-len(g.dump))...)
		image.ReplacePixels(g.dump)
		f, err := os.Create(os.Args[1] + "_dump.png")
		if err != nil {
			fmt.Println("Error creating image")
			return
		}
		defer f.Close()
		err = png.Encode(f, image)
	}(mem, g)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Memory Viewer")
	ebiten.SetFullscreen(true)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
