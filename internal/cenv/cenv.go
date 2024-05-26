package cenv

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"os"
	"strings"
)

func Init() {
	file, err := os.Open(".env")
	defer file.Close()

	if err != nil {
		slog.Error("No ENV file found")
		return
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(file)

	if err != nil {
        slog.Error("Reading ENV file failed: ", err)
		return
	}

	for {
		line, err := buf.ReadString('\n')

		if err != nil {
			if !errors.Is(err, io.EOF) {
                slog.Error("Reading line failed:", err)
			}
			break
		}

        // Delete new line char
        line = line[:len(line)-1]

		data := strings.Split(line, "=")

		if len(data) < 0 || len(data) > 2 {
			slog.Error("\"" + line + "\" " + "is not valid entry")
			continue
		}

		err = os.Setenv(data[0], data[1])

        if err != nil {
            slog.Error(err.Error())
        }
	}
}
