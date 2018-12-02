package text

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/corvus-ch/horcrux/output"
	"github.com/sebdah/goldie"
	"github.com/stretchr/testify/assert"
)

var writeTests = []struct {
	text   string
	data   string
	length int
}{
	{
		text: `   1: 3ody C38D7E
   2: C38D7E
`,
		data: "3ody",
	},

	{
		text: `   1: fkzd3 i13sk 56rdk t5map y7ehp C77D41
   2: jzo 32A020
   3: 32A020
`,
		data: "fkzd3i13sk56rdkt5mapy7ehpjzo",
	},

	{
		text: `   1: tm46g i8g97 m9hib xmkdf xg6uj AAE271
   2: xafg7 cc4hg zdjnx 68qdm q8cc6 44262B
   3: erg97 ku46g 39aak fdpc6 iaf3q E86D44
   4: gije 54EFC6
   5: 54EFC6
`,
		data: "tm46gi8g97m9hibxmkdfxg6ujxafg7cc4hgzdjnx68qdmq8cc6erg97ku46g39aakfdpc6iaf3qgije",
	},

	{
		text: `   1: psqhd 8k96p m8xrt pt55q fhad9 987C45
   2: ywos6 qxo49 68jc3 bi3uz gqm8q EE6133
   3: 46dtd qgt47 pdbsg 86mfe 7cymc F79384
   4: sjxwf jikcj b8pna xyzzw 9xgh4 FFECEA
   5: eynq 1DA73D
   6: 1DA73D
`,
		data: "psqhd8k96pm8xrtpt55qfhad9ywos6qxo4968jc3bi3uzgqm8q46dtdqgt47pdbsg86mfe7cymcsjxwfjikcjb8pnaxyzzw9xgh4eynq",
	},
}

func TestWriterWrite(t *testing.T) {
	var buf bytes.Buffer
	for _, test := range writeTests {
		t.Run(test.data, func(t *testing.T) {
			buf.Reset()
			w, err := newWriter(42, &buf)
			if err != nil {
				t.Fatal(err)
			}
			w.Write([]byte(test.data))
			w.Close()
			assert.Equal(t, test.text, buf.String())
		})
	}
}

func TestWriter_Write_LineLength(t *testing.T) {
	var buf bytes.Buffer

	data := []byte("dmig43jftat8sj71pq1j5cbu4hntdyt")
	for i := 14; i < len(data)+13; i++ {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			buf.Reset()
			w, err := newWriter(uint8(i), &buf)
			if err != nil {
				t.Fatal(err)
			}
			w.Write(data)
			w.Close()
			goldie.Assert(t, fmt.Sprintf("line_length_%d", i), buf.Bytes())
		})
	}
}

func newWriter(n uint8, buf *bytes.Buffer) (io.WriteCloser, error) {
	f := &Format{LineLength: n}
	return NewWriter(buf, f, NewData(f.input, 42, output.NewOutput()))
}
