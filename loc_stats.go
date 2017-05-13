// Usage:
//   go run loc_stats.go path/to/dir

package main

import (
  "bufio"
  "flag"
  "fmt"
  "os"
  "path/filepath"
  "log"
  "strings"
)

type LineType int
type Counts struct {
  code, comment, void int
}
const (
  LineOfCode    = 1 << 0
  LineOfComment = 1 << 1
  LineOfVoid    = 1 << 2
)

func rubyLineType(line string) LineType {
  trimmedLine := strings.TrimLeft(line, " \t")
  if strings.HasPrefix(trimmedLine, "#") {
    return LineOfComment
  }
  if len(trimmedLine) == 0 {
    return LineOfVoid
  }
  return LineOfCode
}

// Calculates the Ruby stats for the file at the specified path.
// It returns 3 integer values representing, in this order, the
// number of code, comment, void.
func rubyFileStats(path string) Counts {
  var counts Counts
  file, err := os.Open(path)
  if err != nil {
    log.Fatal(err)
  }
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lineType := rubyLineType(scanner.Text())
    switch lineType {
    case LineOfCode: counts.code++
    case LineOfComment: counts.comment++
    case LineOfVoid: counts.void++
    }
  }
  if err := scanner.Err(); err != nil {
    fmt.Fprintln(os.Stderr, "reading standard input:", err)
  }
  defer file.Close()
  return counts
}

func rubyWalk(counts map[string]*Counts, dir string) filepath.WalkFunc {
  return func(path string, info os.FileInfo, err error) error {
    localPath := strings.Replace(path, fmt.Sprintf("%s/", dir), "", -1)
    pathItems := strings.Split(localPath, "/")
    pathKey := pathItems[0]
    if strings.HasSuffix(path, ".rb") {
      newCounts := rubyFileStats(path)
      if counts[pathKey] == nil {
        counts[pathKey] = &newCounts
      } else {
        counts[pathKey].code += newCounts.code
        counts[pathKey].comment += newCounts.comment
        counts[pathKey].void += newCounts.void
      }
    }
    return nil
  }
}

func main() {
  flag.Parse()
  var dir = flag.Arg(0)
  var counts = make(map[string]*Counts)
  filepath.Walk(dir, rubyWalk(counts, dir))
  for k, v := range counts {
    fmt.Printf("%s: loc=%d comments=%d void=%d\n", k, v.code, v.comment, v.void)
  }
}
