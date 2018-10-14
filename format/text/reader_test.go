package text

import (
	"bytes"
	"io"
	"testing"
)

func TestReaderRead(t *testing.T) {

	var tests = []struct {
		name string
		text string
		data string
		err  string
	}{
		{
			name: "single byte",
			text: `   1: 3ody C38D7E
   2: C38D7E`,
			data: "3ody",
		},

		{
			name: "16 bytes",
			text: `   1: fkzd3 i13sk 56rdk t5map y7ehp jzo 32A020
   2: 32A020`,
			data: "fkzd3i13sk56rdkt5mapy7ehpjzo",
		},

		{
			name: "48 bytes",
			text: `   1: tm46g i8g97 m9hib xmkdf xg6uj xafg7 cc4hg zdjnx 23D6E3
   2: 68qdm q8cc6 erg97 ku46g 39aak fdpc6 iaf3q gije 54EFC6
   3: 54EFC6`,
			data: "tm46gi8g97m9hibxmkdfxg6ujxafg7cc4hgzdjnx68qdmq8cc6erg97ku46g39aakfdpc6iaf3qgije",
		},

		{
			name: "64 bytes",
			text: `   1: psqhd 8k96p m8xrt pt55q fhad9 ywos6 qxo49 68jc3 D96719
   2: bi3uz gqm8q 46dtd qgt47 pdbsg 86mfe 7cymc sjxwf D88100
   3: jikcj b8pna xyzzw 9xgh4 eynq 1DA73D
   4: 1DA73D`,
			data: "psqhd8k96pm8xrtpt55qfhad9ywos6qxo4968jc3bi3uzgqm8q46dtdqgt47pdbsg86mfe7cymcsjxwfjikcjb8pnaxyzzw9xgh4eynq",
		},

		{
			name: "single byte with description",
			text: `Vestibulum id ligula porta felis euismod semper.
Etiam porta sem malesuada magna mollis euismod.

   1: 3ody C38D7E
   2: C38D7E`,
			data: "3ody",
		},

		{
			name: "64 bytes with description",
			text: `Donec id elit non mi porta gravida at eget metus. Vestibulum id ligula porta felis euismod semper. Maecenas faucibus mollis interdum. Vestibulum id ligula porta felis euismod semper.

Integer posuere erat a ante venenatis dapibus posuere velit aliquet. Nulla vitae elit libero, a pharetra augue. Sed posuere consectetur est at lobortis. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus.

   1: psqhd 8k96p m8xrt pt55q fhad9 ywos6 qxo49 68jc3 D96719
   2: bi3uz gqm8q 46dtd qgt47 pdbsg 86mfe 7cymc sjxwf D88100
   3: jikcj b8pna xyzzw 9xgh4 eynq 1DA73D
   4: 1DA73D`,
			data: "psqhd8k96pm8xrtpt55qfhad9ywos6qxo4968jc3bi3uzgqm8q46dtdqgt47pdbsg86mfe7cymcsjxwfjikcjb8pnaxyzzw9xgh4eynq",
		},

		// Errors
		{name: "no data", text: `foo`},
		{name: "invalid data", text: `1: Lorem impsum`, err: `found "Lorem", expected data or checksum`},
		{name: "missing checksum", text: `6: j7n7s by`, err: `found "", expected data or checksum`},
		{name: "none zbase32 chars", text: `42: j7n7s lorem`, err: `found "lorem", expected data or checksum`},
		{name: "multi line", text: "6 j7n7s by\n7: as8kg", err: `found "", expected data or checksum`},
	}
	for _, test := range tests {
		for bs := int64(4); bs < 128; bs += 4 {
			in := bytes.NewReader([]byte(test.text))
			r := NewReader(in)
			var out bytes.Buffer
			buf := make([]byte, bs)
			for {
				n, err := r.Read(buf)
				if nil != err && io.EOF != err {
					if test.err != err.Error() {
						t.Errorf(`expected error "%v" got "%v" for buffer size %v`, test.err, err, bs)
					}
					break
				}
				out.Write(buf[:n])
				if io.EOF == err {
					break
				}
			}
			if out.String() != test.data {
				t.Errorf(`unexpected data %v got %v for buffer size %v`, test.data, out.String(), bs)
			}
		}
	}

	for _, test := range tests {
		for bs := int64(4); bs < 128; bs += 4 {
			in := bytes.NewReader([]byte(test.text))
			r := NewReader(in)
			var out bytes.Buffer

			for {
				_, err := io.CopyN(&out, r, bs)
				if nil != err && io.EOF != err {
					if test.err != err.Error() {
						t.Errorf(`expected error "%v" got "%v" for buffer size %v`, test.err, err, bs)
					}
					break
				} else if io.EOF == err {
					break
				}
			}
			if out.String() != test.data {
				t.Errorf(`unexpected data %v got %v for buffer size %v`, test.data, out.String(), bs)
			}
		}
	}
}
