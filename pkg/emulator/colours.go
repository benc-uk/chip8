package emulator

import (
	"image/color"
	"strconv"

	"gopkg.in/yaml.v2"
)

type ColourMap struct {
	backgroundColour color.RGBA
	defaultColour    color.RGBA
	spriteColours    map[int]color.RGBA
}

type ColourMapYaml struct {
	BackgroundColour int            `yaml:"background"`
	DefaultColour    int            `yaml:"default"`
	SpriteColours    map[string]int `yaml:"sprites"`
	RangeColours     []colourRange  `yaml:"ranges"`
}

type colourRange struct {
	Colour int
	Start  string
	End    string
}

var pallettes map[string][]color.RGBA

const PalletteSpectrum = "spectrum"
const PalletteC64 = "c64"
const PalletteVaporWave = "vapourwave"

func init() {
	pallettes = make(map[string][]color.RGBA)

	spectrum := make([]color.RGBA, 9)
	spectrum[0] = color.RGBA{R: 0, G: 0, B: 0, A: 255}       // black
	spectrum[1] = color.RGBA{R: 255, G: 255, B: 255, A: 255} // white
	spectrum[2] = color.RGBA{R: 255, G: 0, B: 0, A: 255}     // red
	spectrum[3] = color.RGBA{R: 0, G: 255, B: 0, A: 255}     // green
	spectrum[4] = color.RGBA{R: 0, G: 0, B: 255, A: 255}     // blue
	spectrum[5] = color.RGBA{R: 255, G: 0, B: 255, A: 255}   // magenta
	spectrum[6] = color.RGBA{R: 255, G: 255, B: 0, A: 255}   // yellow
	spectrum[7] = color.RGBA{R: 0, G: 255, B: 255, A: 255}   // cyan
	spectrum[8] = color.RGBA{R: 255, G: 120, B: 0, A: 255}   // orange
	pallettes[PalletteSpectrum] = spectrum

	c64 := make([]color.RGBA, 9)
	c64[0] = color.RGBA{R: 0, G: 0, B: 0, A: 255}       // black
	c64[1] = color.RGBA{R: 255, G: 255, B: 255, A: 255} // white
	c64[2] = color.RGBA{R: 255, G: 119, B: 119, A: 255} // red
	c64[3] = color.RGBA{R: 0, G: 204, B: 85, A: 255}    // green
	c64[4] = color.RGBA{R: 0, G: 0, B: 170, A: 255}     // blue
	c64[5] = color.RGBA{R: 204, G: 68, B: 204, A: 255}  // magenta
	c64[6] = color.RGBA{R: 238, G: 238, B: 119, A: 255} // yellow
	c64[7] = color.RGBA{R: 170, G: 255, B: 238, A: 255} // cyan
	c64[8] = color.RGBA{R: 221, G: 136, B: 85, A: 255}  // orange
	pallettes[PalletteC64] = c64

	vaporwave := make([]color.RGBA, 9)
	vaporwave[0] = color.RGBA{R: 0, G: 0, B: 0, A: 255}       // black
	vaporwave[1] = color.RGBA{R: 255, G: 255, B: 255, A: 255} // white
	vaporwave[2] = color.RGBA{R: 255, G: 106, B: 138, A: 255} // red
	vaporwave[3] = color.RGBA{R: 33, G: 222, B: 138, A: 255}  // green
	vaporwave[4] = color.RGBA{R: 134, G: 149, B: 232, A: 255} // blue
	vaporwave[5] = color.RGBA{R: 255, G: 106, B: 213, A: 255} // magenta
	vaporwave[6] = color.RGBA{R: 254, G: 222, B: 139, A: 255} // yellow
	vaporwave[7] = color.RGBA{R: 147, G: 208, B: 255, A: 255} // cyan
	vaporwave[8] = color.RGBA{R: 255, G: 165, B: 139, A: 255} // orange
	pallettes[PalletteVaporWave] = vaporwave
}

func SimpleColourMap(fgIndex, bgIndex int, pallette string) *ColourMap {
	colourMap := &ColourMap{}
	colourMap.defaultColour = getPalletColour(fgIndex, pallette)
	colourMap.backgroundColour = getPalletColour(bgIndex, pallette)
	return colourMap
}

func LoadColourMap(yamlRaw []byte, pallette string) (*ColourMap, error) {
	mapYaml := &ColourMapYaml{}
	err := yaml.Unmarshal(yamlRaw, mapYaml)
	if err != nil {
		return nil, err
	}

	colourMap := &ColourMap{}
	colourMap.defaultColour = getPalletColour(mapYaml.DefaultColour, pallette)
	colourMap.backgroundColour = getPalletColour(mapYaml.BackgroundColour, pallette)
	colourMap.spriteColours = make(map[int]color.RGBA)

	for _, colourRange := range mapYaml.RangeColours {
		startAddr, err := strconv.ParseInt(colourRange.Start, 16, 16)
		if err != nil {
			continue
		}
		endAddr, err := strconv.ParseInt(colourRange.End, 16, 16)
		if err != nil {
			continue
		}
		for addr := startAddr; addr <= endAddr; addr += 1 {
			colourMap.spriteColours[int(addr)] = getPalletColour(colourRange.Colour, pallette)
		}
	}

	for hexKey, palIndex := range mapYaml.SpriteColours {
		addr, err := strconv.ParseInt(hexKey, 16, 64)
		if err != nil {
			continue
		}
		colourMap.spriteColours[int(addr)] = getPalletColour(palIndex, pallette)
	}

	return colourMap, nil
}

func (cmap ColourMap) getSpriteColour(spriteAddress int) *color.RGBA {
	if cmap.spriteColours == nil {
		return cmap.getDefaultColour()
	}
	if colour, ok := cmap.spriteColours[spriteAddress]; ok {
		return &colour
	}
	return cmap.getDefaultColour()
}

func (cmap ColourMap) getDefaultColour() *color.RGBA {
	return &cmap.defaultColour
}

func (cmap ColourMap) getBackgroundColour() *color.RGBA {
	return &cmap.backgroundColour
}

// Get colour from pallette for given index
func getPalletColour(index int, pallette string) color.RGBA {
	if index <= 8 && index > 0 {
		return pallettes[pallette][index]
	}

	return color.RGBA{R: 0, G: 0, B: 0, A: 255}
}
