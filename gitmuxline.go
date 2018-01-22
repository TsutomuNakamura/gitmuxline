package main

import (
    "fmt"
    "os"
    "os/exec"
    "bufio"
    "strings"
    "errors"
    "os/user"
    "strconv"

    "github.com/shirou/gopsutil/load"
)

type TmuxLine struct {
    num_indexed     int
    num_modified    int
    num_deleted     int
    num_conflict    int
    num_untracked   int
}

func (tmuxLine *TmuxLine) gitmuxline(path string) error {

    var resultLine string
    // num_indexed     := 0
    // num_modified    := 0
    // num_deleted     := 0
    // num_conflict    := 0
    // num_untracked   := 0
    // num_ignored     := 0

    hostname, err := os.Hostname()
    if err != nil {
        return err
    }
    gitline += " " + hostname

    user, err := user.Current()
    if err != nil {
        return err
    }
    gitline += "  " + user.Username

    loadave, err := load.Avg()
    if err != nil {
        return err
    }
    gitline += "  " + strconv.FormatFloat(loadave.Load1, 'f', 2, 64) + " " + strconv.FormatFloat(loadave.Load5, 'f', 2, 64) + " " + strconv.FormatFloat(loadave.Load15, 'f', 2, 64)

    gitstatline, err := getgitstat(path)
    if err != nil {
        return err
    }
    gitline += gitstatline

    fmt.Println(gitline)

    return nil
}

func getgitstat (path string) (string, error) {
    var result_line string = " ⭠ "

    num_indexed     := 0
    num_modified    := 0
    num_deleted     := 0
    num_conflict    := 0
    num_untracked   := 0

    out, err := exec.Command("git", "-C", path, "status", "--porcelain", "--branch").Output()
    if err != nil {
        fmt.Printf("%s\n", err)
        return result_line, errors.New("Failed to get git status in " + path + "\n")
    }

    scanner := bufio.NewScanner(strings.NewReader(string(out)))
    scanner.Scan()
    branch_info := scanner.Text()


    branch_line := strings.Split(branch_info[3:], ".")

    if strings.Contains(branch_line[1], "[ahead ") {
        result_line += " "
    } else if strings.Contains(branch_line[1], "[behind ") {
        result_line += " "
    } else if strings.Contains(branch_line[1], "[") {
        result_line += " "
    } else {
        remote_url_byte, _ := exec.Command("git", "-C", path, "config", "--get", "remote.origin.url").Output()
        remote_url := string(remote_url_byte)

        if strings.Contains(remote_url, "github.com") {
            result_line += " "
        } else if strings.Contains(remote_url, "bitbucket.org") {
            result_line += " "
        } else if strings.Contains(remote_url, "gitlab.com") {
            result_line += " "
        } else {
            result_line += " "
            // result_line += " "
        }
    }
    result_line += branch_line[0]

    var line string
    var stats string
    for scanner.Scan() {
        line = scanner.Text()

        // line[0]: stats of file in index
        // line[1]: stats of file at local
        /*
           X          Y     Meaning
           -------------------------------------------------
                     [MD]   not updated
           M        [ MD]   updated in index
           A        [ MD]   added to index
           D         [ M]   deleted from index
           R        [ MD]   renamed in index
           C        [ MD]   copied in index
           [MARC]           index and work tree matches
           [ MARC]     M    work tree changed since index
           [ MARC]     D    deleted in work tree
           -------------------------------------------------
           D           D    unmerged, both deleted
           A           U    unmerged, added by us
           U           D    unmerged, deleted by them
           U           A    unmerged, added by them
           D           U    unmerged, deleted by us
           A           A    unmerged, both added
           U           U    unmerged, both modified
           -------------------------------------------------
           ?           ?    untracked
           !           !    ignored
           -------------------------------------------------
        */
        stats = line[0:1]
        if line[0] != ' ' {
            num_indexed++
        } else if stats == "??" {
            num_untracked++
        } else {
            if stats == " M" {
                num_modified++
            } else if stats == " D" {
                num_deleted++
            } else if stats == "DD" || stats == "AU" || stats == "UD" || stats == "UA" || stats == "DU" || stats == "AA" || stats == "UU" {
                num_conflict++
            }
        }
    }

    if num_indexed > 99 {
        result_line += "  99+"
    } else {
        result_line += "  " + strconv.Itoa(num_indexed)
    }
    if num_modified > 99 {
        result_line += "  99+"
    } else {
        result_line += "  " + strconv.Itoa(num_modified)
    }
    if num_deleted > 99 {
        result_line += "  99+"
    } else {
        result_line += "  " + strconv.Itoa(num_deleted)
    }
    if num_conflict > 99 {
        result_line += "  99+"
    } else {
        result_line += "  " + strconv.Itoa(num_conflict)
    }
    if num_untracked > 99 {
        result_line += "  99+"
    } else {
        result_line += "  " + strconv.Itoa(num_untracked)
    }

    return result_line, nil
}

func main() {
    err := gitmuxline(os.Args[1])
    if err != nil {
        fmt.Print(err)
    }
}

