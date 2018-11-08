package phonedata

import (
	"testing"
	"time"
	"os"
	"github.com/stretchr/testify/require"
)

func since(t time.Time) int64 {
	return time.Since(t).Nanoseconds()
}


func TestParser(t *testing.T) {
	a := require.New(t)
	f, err := os.Open("./phone.dat")
	a.Nil(err, "fail")
	defer f.Close()

	p, err := NewParser(f)
	a.Nil(err, "fail")

	phone, err := p.Find("13626143333")
	a.Nil(err, "fail")
	t.Log(phone.String())
}
