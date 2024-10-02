package tty

import (
	"context"
	"distrogo/internal/logger"
	"distrogo/internal/tty/cancelable_reader"
	"github.com/pkg/errors"
	"io"
	"os"
	"strings"
)

const (
	Err            = "tty error"
	ErrReadInput   = "input read error"
	ErrWriteInput  = "input write error"
	ErrReadOutput  = "output read error"
	ErrWriteOutput = "output write error"
)

const (
	commandDetach = "detach"
)

type cancelableReader interface {
	Read(p []byte) (int, error)
}

type DetachCallback func(err error)

type TTY struct {
	ctx              context.Context
	inputWriter      io.Writer
	inputReader      cancelableReader
	outputWriter     io.Writer
	outputReader     cancelableReader
	onDetachCallback DetachCallback
}

func NewTTY(ctx context.Context, writer io.Writer, reader io.Reader, onDetachCallback DetachCallback) *TTY {
	tty := &TTY{
		ctx:              ctx,
		inputWriter:      writer,
		inputReader:      cancelable_reader.New(ctx, reader),
		outputWriter:     os.Stdout,
		outputReader:     cancelable_reader.New(ctx, os.Stdin),
		onDetachCallback: onDetachCallback,
	}
	go tty.readRoutine()
	go tty.writeRoutine()
	return tty
}

func (t *TTY) writeRoutine() {
	for {
		select {
		case <-t.ctx.Done():
			return
		default:
			data := make([]byte, cancelable_reader.MaxBufSize)
			_, errRead := t.inputReader.Read(data)
			if errRead != nil {
				logger.Debug(errors.Wrap(errors.Wrap(errRead, ErrReadInput), Err))
				t.onDetachCallback(errRead)
				return
			}
			dataString := string(data)
			_, err := io.WriteString(t.outputWriter, dataString)
			if err != nil {
				logger.Error(errors.Wrap(errors.Wrap(err, ErrWriteOutput), Err))
			}
		}
	}
}

func (t *TTY) readRoutine() {
	for {
		select {
		case <-t.ctx.Done():
			return
		default:
			data := make([]byte, cancelable_reader.MaxBufSize)
			_, errRead := t.outputReader.Read(data)
			if errRead != nil {
				logger.Debug(errors.Wrap(errors.Wrap(errRead, ErrReadOutput), Err))

				if errors.Is(errRead, io.EOF) {
					t.writeToOutput("Exiting container session (Ctrl+D)...")
				}

				t.onDetachCallback(errRead)
				return
			}

			dataString := strings.TrimSuffix(string(data), "\n")
			if t.parseCommand(dataString) {
				return
			}

			_, err := io.WriteString(t.inputWriter, dataString)
			if err != nil {
				logger.Error(errors.Wrap(errors.Wrap(err, ErrWriteInput), Err))
				return
			}
		}
	}
}

func (t *TTY) parseCommand(command string) bool {
	trimmedCommand := strings.TrimSpace(command)
	switch trimmedCommand {
	case commandDetach:
		t.writeToOutput("Exiting container session (detach)...")
		t.onDetachCallback(nil)
		return true
	default:
		return false
	}
}

func (t *TTY) writeToOutput(value string) {
	_, errWrite := t.outputWriter.Write([]byte(value))
	if errWrite != nil {
		logger.Error(errors.Wrap(errors.Wrap(errWrite, ErrWriteOutput), Err))
	}
}
