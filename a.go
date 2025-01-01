package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func main() {
    w := csv.NewWriter(os.Stdout)
    defer w.Flush()
    w.Comma = '\t'

    a := []string{"a", "b"}
    i, j := 0, 1
    x := float64(2.7)

    // Write row.
    err := w.Write(
        []string{
            a[i], a[j],
            strconv.FormatFloat(x, 'f', 4, 64),
        },
    )
    if err != nil {
        fmt.Println(err)
    }
}
