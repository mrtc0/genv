package dotenv

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

func WriteFile(filename string, envMap map[string]string) error {
	content, err := Marshal(envMap)
	if err != nil {
		return err
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content + "\n")
	if err != nil {
		return err
	}
	return file.Sync()
}

func ReadFile(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		return nil, err
	}

	envMap, err := Unmarshal(buf.Bytes())
	if err != nil {
		return nil, err
	}

	return envMap, nil
}

// Marshal is converts a map of environment variables to a dotenv format string.
func Marshal(envMap map[string]string) (string, error) {
	lines := make([]string, 0, len(envMap))

	for k, v := range envMap {
		if d, err := strconv.Atoi(v); err == nil {
			lines = append(lines, fmt.Sprintf(`%s=%d`, k, d))
		} else {
			lines = append(lines, fmt.Sprintf(`%s="%s"`, k, backslashEscape(v)))
		}
	}

	sort.Strings(lines)
	return strings.Join(lines, "\n"), nil
}

func Unmarshal(contnet []byte) (map[string]string, error) {
	lines := strings.Split(string(contnet), "\n")
	envMap := make(map[string]string)

	for _, line := range lines {
		if skipLine(line) {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line: %s", line)
		}

		envMap[parts[0]] = extractValue(parts[1])
	}
	return envMap, nil
}

func extractValue(part string) string {
	if len(part) == 0 {
		return ""
	}

	if part[0] == '"' && part[len(part)-1] == '"' {
		part = part[1 : len(part)-1]
		return backslashUnescape(part)
	}

	if part[0] == '\'' && part[len(part)-1] == '\'' {
		part = part[1 : len(part)-1]
		return backslashUnescape(part)
	}

	return backslashUnescape(part)
}

func backslashEscape(line string) string {
	line = strings.Replace(line, `\`, `\\`, -1)
	line = strings.Replace(line, "\n", `\n`, -1)
	line = strings.Replace(line, "\r", `\r`, -1)
	line = strings.Replace(line, `"`, `\"`, -1)

	return line
}

func backslashUnescape(part string) string {
	part = strings.Replace(part, `\\`, `\`, -1)
	part = strings.Replace(part, `\n`, "\n", -1)
	part = strings.Replace(part, `\r`, "\r", -1)
	part = strings.Replace(part, `\"`, `"`, -1)

	return part
}

func skipLine(line string) bool {
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return true
	}

	if line[0] == '#' || line[0] == '\n' {
		return true
	}

	return false
}
