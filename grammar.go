package grammar

import (
	"bufio"
	"bytes"
	"io"
	"math/rand"
	"time"
	"unicode"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type Grammar struct {
	data  map[string][][]string
	start string
}

func (g Grammar) Execute(w io.Writer) (int, error) {
	return g.execOne(w, g.start)
}

func (g Grammar) execOne(w io.Writer, s string) (int, error) {
	result, ok := g.data[s]
	if !ok {
		toWrite := make([]byte, len(s) + 1)
		copy(toWrite, s)
		toWrite[len(toWrite) - 1] = ' '
		return w.Write(toWrite)
	}
	count := 0
	for _, piece := range result[rand.Intn(len(result))] {
		c, err := g.execOne(w, piece)
		count += c
		if err != nil {
			return count, err
		}
	}
	return count, nil
}

func New(r io.Reader) (*Grammar, error) {
	buf := bufio.NewReader(r)
	g := &Grammar{make(map[string][][]string), "Sentence"}
	line, err := buf.ReadSlice('.')
	for ; err == nil; line, err = buf.ReadSlice('.') {
		splat := bytes.Fields(line) // First field is left side, last is ".".
		stringified := make([]string, len(splat)-2)
		for i, word := range splat[1 : len(splat)-1] {
			stringified[i] = string(word)
		}
		key := string(splat[0])
		g.data[key] = append(g.data[key], stringified)
	}
	if err != io.EOF {
		return g, err
	}
	for _, c := range line { // leftovers
		if !unicode.IsSpace(rune(c)) { // i.e. there's something after the last '.':
			return g, io.ErrUnexpectedEOF
		}
	}
	return g, nil
}
