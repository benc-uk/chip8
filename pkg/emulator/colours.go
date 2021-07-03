package emulator

import (
	"image/color"
	"io/ioutil"
	"log"
	"path/filepath"
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

var pallette []color.RGBA

func init() {
	pallette = make([]color.RGBA, 9)
	pallette[0] = color.RGBA{R: 0, G: 0, B: 0, A: 255}       // black
	pallette[1] = color.RGBA{R: 255, G: 255, B: 255, A: 255} // white
	pallette[2] = color.RGBA{R: 255, G: 0, B: 0, A: 255}     // red
	pallette[3] = color.RGBA{R: 0, G: 255, B: 0, A: 255}     // green
	pallette[4] = color.RGBA{R: 0, G: 0, B: 255, A: 255}     // blue
	pallette[5] = color.RGBA{R: 255, G: 0, B: 255, A: 255}   // magenta
	pallette[6] = color.RGBA{R: 255, G: 255, B: 0, A: 255}   // yellow
	pallette[7] = color.RGBA{R: 0, G: 255, B: 255, A: 255}   // cyan
	pallette[8] = color.RGBA{R: 255, G: 120, B: 0, A: 255}   // orange
}

// var joustMap = ColourMap{
// 	0xA6A: 2,
// 	0xA62: 2,
// }

// var carMap = ColourMap{
// 	0x32c: 1, // roadside
// 	0x330: 3, // grass
// 	0x338: 8, // cars
// }

// var invadersMap = ColourMap{
// 	0x03B7: 2,
// 	0x03CF: 1,
// 	0x03BD: 4,
// 	0x03C3: 5,
// }

// var brixMap = ColourMap{
// 	0x030C: 5,
// 	0x0312: 1,
// 	0x0310: 4,
// 	0x030E: 6,
// }

// var spacejamMap = ColourMap{
// 	0x0396: 1,
// 	0x03CC: 5,
// 	0x0397: 7,
// }

func LoadColourMap(pgmPath string) *ColourMap {
	path := filepath.Dir(pgmPath)
	romName := filepath.Base(pgmPath)
	mapFilePath := filepath.Join(path, romName+".colours.yaml")
	yamlRaw, err := ioutil.ReadFile(mapFilePath)
	if err != nil {
		return nil
	}

	mapYaml := &ColourMapYaml{}
	log.Println(string(yamlRaw))
	err = yaml.Unmarshal(yamlRaw, mapYaml)
	if err != nil {
		return nil
	}

	colourMap := &ColourMap{}
	colourMap.defaultColour = getPalletColour(mapYaml.DefaultColour)
	colourMap.backgroundColour = getPalletColour(mapYaml.BackgroundColour)
	colourMap.spriteColours = make(map[int]color.RGBA)
	for hexKey, palIndex := range mapYaml.SpriteColours {
		addr, err := strconv.ParseInt(hexKey, 16, 64)
		if err != nil {
			continue
		}
		colourMap.spriteColours[int(addr)] = getPalletColour(palIndex)
	}

	return colourMap
}

func (cmap ColourMap) getSpriteColour(spriteAddress int) *color.RGBA {
	if colour, ok := cmap.spriteColours[spriteAddress]; ok {
		return &colour
	}
	return nil
}

func (cmap ColourMap) getDefaultColour() *color.RGBA {
	return &cmap.defaultColour
}

func (cmap ColourMap) getBackgroundColour() *color.RGBA {
	return &cmap.backgroundColour
}

// Get colour from pallette for given index
func getPalletColour(index int) color.RGBA {
	if index < len(pallette) {
		return pallette[index]
	}

	return color.RGBA{R: 0, G: 0, B: 0, A: 255}
}
