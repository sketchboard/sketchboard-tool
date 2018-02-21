package convertcsv

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"

	flags "github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
)

type convertCSVCmdOpts struct {
	CSVFile    string `short:"f" long:"file" description:"Read first column from a CSV file" required:"true"`
	NoteWidth  uint   `short:"w" long:"note-width" description:"Specify a note width"`
	NoteHeight uint   `short:"t" long:"note-height" description:"Specify a note height"`
}

var convertCSVCmd convertCSVCmdOpts

func Init(parser *flags.Parser) error {
	_, err := parser.AddCommand(
		"convertcsv",
		"Convert CSV first column to note shapes",
		"",
		&convertCSVCmd,
	)
	if err != nil {
		return errors.Wrapf(err, "convertcsv: Init")
	}

	return nil
}

func (ccmd *convertCSVCmdOpts) Execute(args []string) error {

	file, err := os.Open(ccmd.CSVFile)
	if err != nil {
		return errors.Wrapf(err, "convertcsv: open %s", ccmd.CSVFile)
	}
	csvReader := csv.NewReader(file)
	rows, err := csvReader.ReadAll()
	if err != nil {
		return errors.Wrapf(err, "convertcsv: read all from csv")
	}

	width := func() uint {
		if ccmd.NoteWidth <= 100 {
			return 200
		}
		return ccmd.NoteWidth
	}()

	height := func() uint {
		if ccmd.NoteHeight <= 40 {
			return 50
		}
		return ccmd.NoteHeight
	}()

	items, err := ccmd.toSketchboardClipboardJSON(rows, width, height)
	if err != nil {
		return errors.Wrapf(err, "convertcmd: sketchboard clipboard json")
	}

	err = printJSON(items)
	if err != nil {
		return err
	}

	return nil
}

func (ccmd *convertCSVCmdOpts) toSketchboardClipboardJSON(rows [][]string, noteWidth, noteHeight uint) (*ClipboardDiagram, error) {
	items := make([]DiagramItem, 0)

	lastShapeDimension := &ShapeDimension{}

	var rowFirstDimension *ShapeDimension

	for index, row := range rows {
		if len(row) <= 0 {
			continue
		}

		firstCol := row[0]

		x := lastShapeDimension.X + int(lastShapeDimension.Width) + 30
		y := lastShapeDimension.Y

		di := ccmd.toDiagramItem(firstCol, x, y, noteWidth, noteHeight, index)

		dim, err := di.ParseShape()
		if err != nil {
			return nil, errors.Wrapf(err, "convertcmd: last shape dimension")
		}

		lastShapeDimension = dim

		if rowFirstDimension == nil {
			rowFirstDimension = dim
		}

		if lastShapeDimension.Height > rowFirstDimension.Height {
			// take heighest on row
			// NOTE currently not calculating dynamic height
			rowFirstDimension.Height = lastShapeDimension.Height
		}

		if (lastShapeDimension.X - rowFirstDimension.X) > 1000 {
			lastShapeDimension.X = 0
			lastShapeDimension.Y = rowFirstDimension.Y + int(rowFirstDimension.Height) + 30

			rowFirstDimension = lastShapeDimension
			rowFirstDimension.Width = 0
		}

		items = append(items, di)
	}

	return &ClipboardDiagram{
		BoardID: "NONE",
		Items:   items,
	}, nil
}

type ElementType string

const (
	ElementTypeNote = "noteitem"
)

func (ccmd *convertCSVCmdOpts) toDiagramItem(text string, x, y int, width, height uint, index int) DiagramItem {

	di := DiagramItem{
		Text:            text,
		ElementType:     ElementTypeNote,
		Shape:           fmt.Sprintf("%d,%d,%d,%d", x, y, width, height),
		BackgroundColor: "204,204,255,0:51,51,51,1",
		TextColor:       "51,51,51,1",
		Props:           512,
		Version:         0,
		ID:              0,
		ClientID:        fmt.Sprintf("1-%d", index),
		CustomData:      "",
	}

	return di
}

func JSONMarshal(v interface{}, pretty, unescape bool) ([]byte, error) {
	marhal := func() ([]byte, error) {
		if pretty {
			return json.MarshalIndent(v, "", "  ")
		} else {
			return json.Marshal(v)
		}

	}
	b, err := marhal()

	if unescape {
		b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
		b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
		b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	}
	return b, err
}

type ClipboardDiagram struct {
	BoardID string        `json:"boardId"`
	Items   []DiagramItem `json:"items"`
}

type DiagramItem struct {
	Text            string `json:"text"`
	ElementType     string `json:"elementType"`
	Shape           string `json:"shape"`
	BackgroundColor string `json:"backgroundColor"`
	TextColor       string `json:"textColor"`
	Props           uint64 `json:"props"`
	Version         uint64 `json:"version"`
	ID              uint64 `json:"id"`
	ClientID        string `json:"clientId"`
	CustomData      string `json:"cd"`
}

type ShapeDimension struct {
	X      int
	Y      int
	Width  uint
	Height uint
}

var shapeRe = regexp.MustCompile(`([-]*\d+),([-]*\d+),(\d+),(\d+)`)

func (di *DiagramItem) ParseShape() (*ShapeDimension, error) {
	match := shapeRe.FindAllStringSubmatch(di.Shape, -1)

	if match == nil || len(match) == 0 || len(match[0]) != 5 {
		return nil, errors.Errorf("convertcsv: invalid shape %s", di.Shape)
	}

	m := match[0]

	x, err := strconv.Atoi(m[1])
	if err != nil {
		return nil, errors.Wrapf(err, "convertcsv: x %d", m[1])
	}
	y, err := strconv.Atoi(m[2])
	if err != nil {
		return nil, errors.Wrapf(err, "convertcsv: y %d", m[2])
	}
	width, err := strconv.Atoi(m[3])
	if err != nil {
		return nil, errors.Wrapf(err, "convertcsv: width %d", m[3])
	}
	height, err := strconv.Atoi(m[4])
	if err != nil {
		return nil, errors.Wrapf(err, "convertcsv: height %d", m[4])
	}

	if x > 20000 || x < -20000 {
		return nil, errors.Errorf("convertcsv: invalid x %d", x)
	}
	if y > 20000 || y < -20000 {
		return nil, errors.Errorf("convertcsv: invalid y %d", y)
	}

	if width <= 10 || width > 500 {
		return nil, errors.Errorf("convertcsv: invalid width %d", width)
	}
	if height <= 10 || height > 500 {
		return nil, errors.Errorf("convertcsv: invalid height %d", height)
	}

	return &ShapeDimension{
		X:      x,
		Y:      y,
		Width:  uint(width),
		Height: uint(height),
	}, nil
}

func printJSON(jsonData *ClipboardDiagram) error {
	bjson, err := JSONMarshal(jsonData, true, false)
	if err != nil {
		return errors.Wrapf(err, "convertcsv: print JSON")
	}

	str := string(bjson)
	fmt.Println(str)

	// clipboard.CopyToClipboard(str)

	return nil
}
