package main

import (
	"io"
	"os"
	"strings"
)

type rot13Reader struct {
	r io.Reader
}

func (r rot13Reader) Read (b []byte) (int, error) {
	tmp := make([]byte, 1)
	result := []byte{}
	n, err := r.r.Read(tmp);
	if tmp[0] > 65 && 78 > tmp[0] || tmp[0] > 97 && tmp[0] < 110 {
		tmp[0] = tmp[0] + 13
	} else if 79 <= tmp[0] && tmp[0] <= 90 || 110 <= tmp[0] && tmp[0] <= 120 {
		tmp[0] = tmp[0] - 13
	}
	result = append(result, tmp[0])
	copy(b, result)
	return n, err
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
