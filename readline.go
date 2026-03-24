package main

import (
	"os"
	"path/filepath"

	"github.com/chzyer/readline"
	"robpike.io/ivy/config"
)

// ReadlineReader wraps readline.Instance to implement io.Reader.
type ReadlineReader struct {
	rl   *readline.Instance
	conf *config.Config
	buf  []byte
}

// NewReadlineReader returns a new ReadlineReader if input is a terminal.
// It returns nil, nil if input is not a terminal.
func NewReadlineReader(conf *config.Config) (*ReadlineReader, error) {
	if !isTerminal() {
		return nil, nil
	}

	historyFile := ""
	homeDir, err := os.UserHomeDir()
	if err == nil {
		historyFile = filepath.Join(homeDir, ".ivy_history")
	}

	rlConf := &readline.Config{
		Prompt:          conf.Prompt(),
		HistoryFile:     historyFile,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	}

	rl, err := readline.NewEx(rlConf)
	if err != nil {
		return nil, err
	}

	return &ReadlineReader{rl: rl, conf: conf}, nil
}

func (r *ReadlineReader) Read(p []byte) (n int, err error) {
	for len(r.buf) == 0 {
		r.rl.SetPrompt(r.conf.Prompt())
		line, err := r.rl.Readline()
		if err != nil {
			return 0, err
		}
		// readline strips the newline, but ivy needs it.
		r.buf = []byte(line + "\n")
	}

	n = copy(p, r.buf)
	r.buf = r.buf[n:]
	return n, nil
}

func (r *ReadlineReader) Close() error {
	return r.rl.Close()
}

// isTerminal reports whether stdin is a terminal.
func isTerminal() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
