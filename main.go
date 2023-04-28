package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"math"
	"os"
	"strings"
	"unicode"

	"github.com/tidwall/gjson"
)

func main() {
	args := os.Args
	switch args[1] {
	case "manifest":
		printManifest()
	case "transform":
		printTransform()
	}
}

func printManifest() {
	manifest := `
    {
        "name": "fen",
        "version": "0.1",
        "description": "Create chess boards from FEN strings.",
        "transforms": [
            {
                "from": "fen",
                "to": ["html"],
                "arguments": [
					{
						"name": "width",
						"description": "Width of the SVG, given as a ratio to the surrounding figure tag (created automatically).",
						"default": 1.0,
						"type": "f64"
					},
					{
						"name": "save",
						"description": "The name of the SVG file that is saved. No file is saved if this argument is left empty.",
						"default": ""
					}
				]
            }
        ]
    }
    `
	fmt.Println(manifest)
}

func printTransform() {
	scanner := bufio.NewScanner(os.Stdin)
	input := ""
	for scanner.Scan() {
		input += scanner.Text()
	}

	fen := gjson.Get(input, "data").String()
	width := gjson.Get(input, "arguments.width").Float()
	save := gjson.Get(input, "arguments.save").String()

	svg := getSVG(strings.TrimSpace(fen))
	html := getHTML(svg, width)

	if len(save) != 0 {
		if err := os.WriteFile(save, []byte(svg), 0644); err != nil {
			os.Stderr.WriteString("Could not save SVG to file.")
		}
	}

	println(`["`, strings.ReplaceAll(html, `"`, `\"`), `"]`)
}

func getHTML(svg string, width float64) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(svg))
	html := `<figure>`
	html += `<img src="data:image/svg+xml;base64,` + encoded + `"`
	if width >= 0.0 {
		percentage := math.Round(width * 100.0)
		html += ` style="width:` + fmt.Sprint(percentage) + `%"`
	}
	html += `/></figure>`
	return html
}

func getSVG(fen string) string {
	svg := `<svg xmlns="http://www.w3.org/2000/svg" width="400" height="400">`
	svg += getBoardBackground()
	svg += getPieces(fen)
	svg += `</svg>`
	return svg
}

func getBoardBackground() string {
	str := `<rect width="400" height="400" style="fill:#FFCD9F"/>`
	start_x := 50
	for y := 0; y <= 350; y += 50 {
		for x := start_x; x <= 350; x += 100 {
			str += `<rect x="`
			str += fmt.Sprint(x)
			str += `" y="`
			str += fmt.Sprint(y)
			str += `" width="50" height="50" style="fill:#D08C52"/>`
		}
		start_x = start_x ^ 50
	}
	return str
}

func getPieces(fen string) string {
	str := ""
	i := 0
	for _, r := range strings.ReplaceAll(fen, "/", "") {
		if unicode.IsDigit(r) {
			i += int(r - '0')
		} else {
			x := (i % 8) * 50
			y := (i / 8) * 50
			str += getPiece(r, x, y)
			i += 1
		}

		if i > 63 {
			break
		}
	}
	return str
}

func getPiece(p rune, x int, y int) string {
	var svg string
	switch p {
	case 'K':
		svg = WhiteKing
	case 'Q':
		svg = WhiteQueen
	case 'R':
		svg = WhiteRook
	case 'B':
		svg = WhiteBishop
	case 'N':
		svg = WhiteKnight
	case 'P':
		svg = WhitePawn
	case 'k':
		svg = BlackKing
	case 'q':
		svg = BlackQueen
	case 'r':
		svg = BlackRook
	case 'b':
		svg = BlackBishop
	case 'n':
		svg = BlackKnight
	case 'p':
		svg = BlackPawn
	}
	coords := `x="` + fmt.Sprint(x+2) + `" y="` + fmt.Sprint(y+2) + `" `
	svg = strings.ReplaceAll(svg, "<svg", "<svg "+coords)
	svg = strings.ReplaceAll(svg, "\n", "")
	return svg
}
