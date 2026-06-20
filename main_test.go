package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/CoderParth/recall/shell"
)

const testRecallFileDir = "test_recall_dir"
const testRecallFileName = "test_recall_data"

func TestReadOrCreateFile(t *testing.T) {
	recallTestFile := readOrCreateFile(testRecallFileDir, testRecallFileName)
	_, err := recallTestFile.Stat()
	if err != nil {
		t.Errorf("Error opening or creating recall test file: %v \n", err)
	}
	recallTestFile.Close()
	deleteTestDir()
}

func TestEncode(t *testing.T) {
	// Double quotes within command or titles
	// are handled by encoding '\' infront of the double quote.
	input := `A title with added "double quotes".`
	received := encodeRecord(input)
	expected := `A title with added \"double quotes\".`
	if received != expected {
		t.Errorf("Expected: %s \n Received: %s \n", expected, received)
	}
}

func TestDecode(t *testing.T) {
	input := fmt.Sprintf(`"Run test in python":"pytest"%c`, '\n')
	receivedTitle, receivedCommand := decodeRecord(input)
	expectedTitle, expectedCommand := "Run test in python", "pytest"
	if receivedTitle != expectedTitle {
		t.Errorf("Expected Title: %s \n Received Title: %s \n", expectedTitle, receivedTitle)
	}

	if receivedCommand != expectedCommand {
		t.Errorf("Expected Command: %s \n Received Command: %s \n", expectedCommand, receivedCommand)
	}
}

// Test shell run.
var configs = [3]string{
	`
_recall() {
    local cmd
    cmd=$(recall) || return
    READLINE_LINE="$cmd"
    READLINE_POINT=${#READLINE_LINE}
}
bind -x '"\C-f":"_recall"'
`,
	`
_recall() {
    local cmd
    cmd=$(recall) || return
    BUFFER="$cmd"
    CURSOR=${#BUFFER}
    zle reset-prompt
}
zle -N _recall
bindkey '^F' _recall
`,
	`
function _recall
    set cmd (recall)
    or return
    commandline --replace -- $cmd
    commandline --function repaint
end
bind \cf _recall
`,
}

func TestShellRun(t *testing.T) {
	shells := []string{"bash", "zsh", "fish"}
	for i, s := range shells {
		config, err := shell.Run(s)
		if err != nil {
			panic(err)
		}
		if config != configs[i] {
			t.Errorf("Expected: %s \n Received: %s \n", configs[i], config)
		}
	}
}

func TestSaveToFile(t *testing.T) {
	recallTestFile := readOrCreateFile(testRecallFileDir, testRecallFileName)
	r := Record{
		Title:   `A title with added "double quotes".`,
		Command: "npm double quote command",
	}

	saveToFile(recallTestFile, r)
	recallTestFile.Sync()

	// // reset cursor to beginning
	recallTestFile.Seek(0, io.SeekStart)

	// Read from file, and check if the content exists.
	var buf bytes.Buffer
	readFile(recallTestFile, &buf)

	reader := bufio.NewReader(bytes.NewReader(buf.Bytes()))
	line, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	title, command := decodeRecord(line)

	if title != r.Title {
		t.Errorf("Expected Title: %s \n Received Title: %s \n", r.Title, title)
	}

	if command != r.Command {
		t.Errorf("Expected Command: %s \n Received Command: %s \n", r.Command, command)
	}

	recallTestFile.Close()
	deleteTestDir()
}

func TestDeleteCommand(t *testing.T) {
	recallTestFile := readOrCreateFile(testRecallFileDir, testRecallFileName)

	r := Record{
		Title:   `Run test in python`,
		Command: "pytest",
	}
	saveToFile(recallTestFile, r)

	var buf bytes.Buffer
	readFile(recallTestFile, &buf)
	bufBytes := buf.Bytes()
	input := fmt.Sprintf(`"Run test in python":"pytest"%c`, '\n') // "title":"command"
	deleteRecord(recallTestFile, &bufBytes, input)

	if len(bufBytes) > 0 {
		t.Errorf("Buffer still has some content left in it. Length: %v ", len(bufBytes))
	}

	recallTestFile.Close()
	deleteTestDir()
}

func deleteTestDir() {
	// Get the default root directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	// Create full dir path.
	path := filepath.Join(configDir, testRecallFileDir)

	err = os.RemoveAll(path)
	if err != nil {
		panic(err)

	}
}
