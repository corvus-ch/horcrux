package text

import "github.com/corvus-ch/horcrux/input"

type Data struct {
	Input input.Input
	Lines chan Line
}

func NewData(in input.Input, lines chan Line) *Data {
	return &Data{
		Input: in,
		Lines: lines,
	}
}
