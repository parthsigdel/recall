// A place to save commands that you think you may use in future.

// Short glosssary:
// Record: Record refers to the combination of a title (for a command), and the command itself.
// Example: "Run test in python":"pytest"

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/CoderParth/recall/shell"
)

const recallDirName = "recall"
const recallFileName = "data"

var (
	save      = flag.Bool("s", false, "Save a record.")
	initShell = flag.String("init", "", "Initialize shell for keybindings.")
)

type Record struct {
	Title   string
	Command string
}

func main() {
	flag.Parse()
	if *initShell != "" {
		config, err := shell.Run(*initShell)
		if err != nil {
			panic(err)
		}
		fmt.Print(config)
		return
	}

	recallFile := readOrCreateFile(recallDirName, recallFileName)
	defer recallFile.Close()

	if *save {
		record := parseStdInput()
		saveToFile(recallFile, record)
		return
	}

	var buf bytes.Buffer
	readFile(recallFile, &buf)

	p := tea.NewProgram(initialModel(buf.Bytes(), recallFile),
		tea.WithOutput(os.Stderr))
	m, err := p.Run()
	if err != nil {
		panic(err)
	}

	fmt.Print(m.(model).choice)
}

func deleteRecord(recallFile *os.File, buf *[]byte, record string) {
	record = strings.TrimSpace(record)
	reader := bufio.NewReader(bytes.NewReader(*buf))
	var newBuf bytes.Buffer

	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		if line == record {
			continue
		}
		newBuf.WriteString(line)
		newBuf.WriteByte('\n')
	}

	// Overwrite the file with new buffer.
	recallFile.Seek(0, 0)
	recallFile.Truncate(0)
	_, err := recallFile.WriteString(newBuf.String())
	if err != nil {
		panic(err)
	}

	*buf = newBuf.Bytes()
}

func parseStdInput() Record {
	var title, command string
	var scanner *bufio.Scanner
	for {
		fmt.Print("Title: ")
		scanner = bufio.NewScanner(os.Stdin)
		scanner.Scan()
		title = scanner.Text()
		if len(strings.TrimSpace(title)) == 0 {
			fmt.Println("Title cannot be empty.")
			continue
		}
		break
	}

	for {
		fmt.Print("Command: ")
		scanner = bufio.NewScanner(os.Stdin)
		scanner.Scan()
		command = scanner.Text()
		if len(strings.TrimSpace(command)) == 0 {
			fmt.Println("Command cannot be empty.")
			continue
		}
		break
	}

	return Record{
		Title:   title,
		Command: command,
	}
}

func searchInBuffer(reader *bufio.Reader, searchString string) []string {
	match := []string{}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		if strings.Contains(strings.ToLower(line), strings.ToLower(searchString)) {
			match = append(match, line)
		}

	}
	return match
}

func decodeRecord(s string) (string, string) {
	var title, command strings.Builder
	i := 1 // start i = 1 to skip the first double quote (") from the encoded title

	// Build title
	for i < len(s) {
		ch := s[i]
		if ch == '\\' && i < len(s)-1 {
			if s[i+1] == '"' {
				i += 1
				continue
			}
		}
		if ch == '"' && i < len(s)-2 {
			if s[i+1] == ':' && s[i+2] == '"' {
				// Found delimeter (":")
				i += 3 // skip the delimeter of length 3
				break
			}
		}
		i += 1
		title.WriteByte(ch)
	}

	// Build command
	for i < len(s)-2 { // skip \n and " present on the last two indexes.
		ch := s[i]
		if ch == '\\' && i < len(s)-1 {
			if s[i+1] == '"' {
				i += 1
				continue
			}
		}
		i += 1
		command.WriteByte(ch)
	}

	return title.String(), command.String()
}

func saveToFile(recallFile *os.File, r Record) {
	encodedRecord := fmt.Sprintf(`"%s":"%s"%c`, encodeRecord(r.Title), encodeRecord(r.Command), '\n')
	_, err := recallFile.WriteString(encodedRecord)
	if err != nil {
		panic(err)
	}
}

func encodeRecord(s string) string {
	var b strings.Builder
	for _, ch := range s {
		if ch == '"' {
			b.WriteRune('\\')
			b.WriteRune('"')
			continue
		}
		b.WriteRune(ch)
	}
	return b.String()
}

// Open a recall/data file in user config dir, or create one if it does not exist.
func readOrCreateFile(dir, file string) *os.File {
	// Get the default root directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	// Create full file path.
	path := filepath.Join(configDir, dir, file)

	// Create a dir, if it doesnot exist yet.
	// .Dir strips the last file name "data" from filepath
	err = os.MkdirAll(filepath.Dir(path), 0755) // 0755 grants full read, write, and execute persmissions to the owner.
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	return f
}

func readFile(f *os.File, buf *bytes.Buffer) {
	_, err := f.Stat()
	if err != nil {
		panic(err)
	}

	_, err = buf.ReadFrom(f)
	if err != nil {
		panic(err)
	}
}
