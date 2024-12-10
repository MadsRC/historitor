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
	l, err := historitor.NewLog(historitor.WithLogName(t.Name()))
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
	l, err := historitor.NewLog(historitor.WithLogName(t.Name()))
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

// TestLog_Read
func TestLog_Read(t *testing.T) {
	c := historitor.NewConsumer(historitor.WithConsumerName(t.Name()))
	cg := historitor.NewConsumerGroup(historitor.WithConsumerGroupName(t.Name()), historitor.WithConsumerGroupMember(c))
	l, err := historitor.NewLog(historitor.WithLogName(t.Name()))
	require.NoError(t, err)
	l.AddGroup(cg)

	_ = l.Write("value")

	entries, err := l.Read(cg.GetName(), c.GetName(), 1)
	require.NoError(t, err)
	err = l.Acknowledge(cg.GetName(), c.GetName(), entries[0].ID)
	require.NoError(t, err)

	require.Equal(t, 1, len(entries))
	require.Equal(t, "value", entries[0].Payload.(string))

}

// TestLog_Write_Read_ordered tests that the log writes and reads the expected number of lines in order of writing when
// a single consumer reads the log.
func TestLog_Write_Read_ordered(t *testing.T) {
	c := historitor.NewConsumer(historitor.WithConsumerName(t.Name()))
	cg := historitor.NewConsumerGroup(historitor.WithConsumerGroupName(t.Name()), historitor.WithConsumerGroupMember(c))
	l, err := historitor.NewLog(historitor.WithLogName(t.Name()))
	require.NoError(t, err)
	l.AddGroup(cg)

	written := make([]string, 0)
	out := make([]string, 0)

	n := forLine(t, func(line string) {
		written = append(written, line)
		l.Write(line)
	})

	for _, v := range written {
		l.Write(v)
	}

	entries, err := l.Read(cg.GetName(), c.GetName(), n)
	require.NoError(t, err)

	for _, entry := range entries {
		out = append(out, entry.Payload.(string))
		err = l.Acknowledge(cg.GetName(), c.GetName(), entry.ID)
		require.NoError(t, err)
	}

	require.Equal(t, n, len(entries))
	require.Equal(t, written, out)
}

// TestLog_UpdateEntry tests that UpdateEntry updates the entry with the provided ID.
func TestLog_UpdateEntry(t *testing.T) {
	c := historitor.NewConsumer(historitor.WithConsumerName(t.Name()))
	cg := historitor.NewConsumerGroup(historitor.WithConsumerGroupName(t.Name()), historitor.WithConsumerGroupMember(c))
	l, err := historitor.NewLog(historitor.WithLogName(t.Name()))
	require.NoError(t, err)
	l.AddGroup(cg)

	entryID1 := l.Write("valueOne")
	entryID2 := l.Write("valueTwo")
	entryID3 := l.Write("valueThree")

	ok := l.UpdateEntry(entryID2, "valueTwoUpdated")
	require.True(t, ok)

	entries, err := l.Read(cg.GetName(), c.GetName(), 3)
	require.NoError(t, err)

	require.Equal(t, 3, len(entries))
	require.Equal(t, "valueOne", entries[0].Payload.(string))
	require.Equal(t, entryID1, entries[0].ID)
	require.Equal(t, "valueTwoUpdated", entries[1].Payload.(string))
	require.Equal(t, entryID2, entries[1].ID)
	require.Equal(t, "valueThree", entries[2].Payload.(string))
	require.Equal(t, entryID3, entries[2].ID)
}
