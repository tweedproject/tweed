package tweed

import (
	"fmt"
	"io"
	"os/exec"
)

type Exec struct {
	Run    string
	Env    []string
	Stdout chan string
	Stderr chan string
	Done   chan int
}

func (e Exec) drain(out chan string, in io.Reader, buf []byte) {
	for {
		n, err := in.Read(buf)
		if err != nil {
			if err == io.EOF {
				return
			}
			out <- fmt.Sprintf("---\nERROR: %s\n", err)
			return
		}
		out <- string(buf[:n])
	}
}

func run1(e Exec) ([]byte, error) {
	return e.run1()
}

func (e Exec) run1() ([]byte, error) {
	cmd := exec.Command(e.Run)
	cmd.Dir = "/"
	cmd.Env = e.Env
	cmd.Stderr = nil

	return cmd.Output()
}

func run(e Exec) error {
	return e.run()
}

// FIXME this elaborateness is not needed anymore
func (e Exec) run() error {
	cmd := exec.Command(e.Run)
	cmd.Dir = "/"
	cmd.Env = e.Env
	cmd.Stderr = nil

	fd1, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	fd2, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	go e.drain(e.Stdout, fd1, make([]byte, 8192))
	go e.drain(e.Stderr, fd2, make([]byte, 8192))

	fmt.Printf("running `%s'\n", e.Run)
	cmd.Start()
	err = cmd.Wait()
	if cmd.ProcessState.Exited() {
		e.Done <- cmd.ProcessState.ExitCode()
	}
	close(e.Done)

	return err
}
