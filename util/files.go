package util 
import (
    "os"
    "bufio"
    "bytes"
    "io"
    "strings"
)
func ReadLines(path string) (lines []string, err error) {
    var (
        file *os.File
        part []byte
        prefix bool
    )
    if file, err = os.Open(path); err != nil {
        return
    }
    defer file.Close()

    reader := bufio.NewReader(file)
    buffer := bytes.NewBuffer(make([]byte, 0))
    for {
        if part, prefix, err = reader.ReadLine(); err != nil {
            break
        }
        buffer.Write(part)
        if !prefix {
            lines = append(lines, buffer.String())
            buffer.Reset()
        }
    }
    if err == io.EOF {
        err = nil
    }
    return
}

func WriteLines(lines []string, path string) (err error) {
    var (
        file *os.File
    )
    if file, err = os.Create(path); err != nil {
        return
    }
    defer file.Close()
    for _,item := range lines {
        file.WriteString(strings.TrimSpace(item) + "\n"); 
    }
    return
}
