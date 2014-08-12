package main

import (
    "fmt"
    "os"
    "path/filepath"
    "bufio"
    "log"
    "regexp"
    "strconv"
    "sort"
)

type Req struct {
    No      int
    Dep     int
    Name    string
    Desc    string
}

type ByReq []Req

func (r ByReq) Len() int {return len(r)}
func (r ByReq) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r ByReq) Less(i, j int) bool { return r[i].No < r[j].No }

var Reqs ByReq

func ReadReqs(fp string) {
    re := regexp.MustCompile(`{(?P<req>\d{1,5})}{(?P<bel>\d{1,5})}{(?P<name>.*?)}`)

    // Open an input file, exit on error.
    inputFile, err := os.Open(fp)
    if err != nil {
        log.Fatal("Error opening input file:", err)
    }

    // Closes the file when we leave the scope of the current function,
    // this makes sure we never forget to close the file if the
    // function can exit in multiple places.
    defer inputFile.Close()

    scanner := bufio.NewScanner(inputFile)


    for scanner.Scan() {
        scanner_string := string(scanner.Text())

        matches := re.FindStringSubmatch(scanner_string)

        if len(matches) >= 4 {
            var CReq Req;
            value, _ := strconv.Atoi(matches[1])
            CReq.No = value
            for _, req := range Reqs {
                if req.No == value {
                    fmt.Println("      Error Duplicate Req: ", req)
                }
            }
            value, _ = strconv.Atoi(matches[2])
            CReq.Dep = value
            CReq.Name = matches[3]

            Reqs = append(Reqs, CReq)
            // fmt.Printf("Req # %d => %s\n", CReq.Req, CReq.Name)
        }
    }

    // When finished scanning if any error other than io.EOF occured
    // it will be returned by scanner.Err().
    if err := scanner.Err(); err != nil {
        log.Fatal(scanner.Err())
    }
}

func VisitFile(fp string, fi os.FileInfo, err error) error {
    if err != nil {
        fmt.Println(err) // can't walk here,
        return nil       // but continue walking elsewhere
    }
    if !!fi.IsDir() {
        return nil // not a file.  ignore.
    }
    matched, err := filepath.Match("*.tex", fi.Name())
    if err != nil {
        fmt.Println(err) // malformed pattern
        return err       // this is fatal.
    }
    if matched {
        // fmt.Println(fp, )
        ReadReqs(fp)
    }
    return nil
}

func main() {
    filepath.Walk("./", VisitFile)

    if len(Reqs) > 0 {
        sort.Sort(ByReq(Reqs))
        for _, req := range Reqs {
            fmt.Println(req.No, " => ", req.Dep, "  ", req.Name)
        }
    }
}
