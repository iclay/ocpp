package logwriter

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestSpliter(t *testing.T) {
	h := &HourlySplit{
		Dir:           "./tmp",
		FileFormat:    "log_2006-01-02T15",
		MaxFileNumber: 2,
		MaxDiskUsage:  100,
	}
	defer os.RemoveAll("./tmp")

	for i := 0; i < 30; i++ {
		time.Sleep(time.Duration(time.Second))
		n, err := h.Write([]byte("1234567890"))
		if err != nil {
			t.Error(err)
		}
		log.Println("write ", n)
	}
}
