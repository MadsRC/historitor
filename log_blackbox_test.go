package historitor_test

import (
	"bufio"
	"github.com/MadsRC/historitor"
	"github.com/stretchr/testify/require"
	"os"
	"sync"
	"testing"
)

type TestingT interface {
	Helper()
	Errorf(format string, args ...interface{})
	FailNow()
}

// forLine reads the entire "On the Origin of Species" book by Charles Darwin and calls the provided function for each
// line in the file. It returns the number of lines in the file.
// This file was grabbed from Open Library: https://openlibrary.org/books/OL7101861M/On_the_Origin_of_Species_by_Means_of_Natural_Selection
func forLine(t TestingT, f func(string)) int {
	t.Helper()
	file, err := os.Open("testdata/onoriginofspecie00darw_djvu.txt")
	defer func() {
		_ = file.Close()
	}()
	require.NoError(t, err)

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	noLines := 0

	for fileScanner.Scan() {
		noLines++
		f(fileScanner.Text())
	}

	return noLines
}

// TestLog_Write_lines tests that the log writes the expected number of lines.
// We are reading the entire "On the Origin of Species" book by Charles Darwin and writing each line to the log.
// We then compare the number of lines in the file with the number of lines in the log.
func TestLog_Write_lines(t *testing.T) {
	l, err := historitor.NewLog(historitor.WithName(t.Name()))
	require.NoError(t, err)

	n := forLine(t, func(line string) {
		l.Write(line)
	})

	require.Equal(t, n, l.Size())

}

// TestLog_Write_concurrent tests that the log writes the expected number of lines concurrently.
// We are partitioning the file into n parts and writing each part concurrently to the log.
//
// This test should be run with Go's `-race` flag to check for race conditions.
func TestLog_Write_concurrent(t *testing.T) {
	l, err := historitor.NewLog(historitor.WithName(t.Name()))
	require.NoError(t, err)

	var n = 10

	bufs := make([][]string, n)

	for i := 0; i < n; i++ {
		bufs[i] = make([]string, 0)
	}

	var iter = 0

	// split the file into n parts
	tot := forLine(t, func(line string) {
		bufs[iter] = append(bufs[iter], line)
		if iter == n-1 {
			iter = 0
		} else {
			iter++
		}
	})

	wg := sync.WaitGroup{}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for _, line := range bufs[i] {
				l.Write(line)
			}
		}(i)
	}

	wg.Wait()

	require.Equal(t, tot, l.Size())
}
