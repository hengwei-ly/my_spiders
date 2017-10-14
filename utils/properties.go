package utils

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"unicode"
)

func readQouteString(txt string) (string, string) {
	var buf bytes.Buffer

	isEscape := false
	for idx, c := range txt {
		switch c {
		case '\\':
			if isEscape {
				isEscape = false
				buf.WriteRune(c)
			} else {
				isEscape = true
			}
		case 't':
			if isEscape {
				isEscape = false
				buf.WriteRune('\t')
			} else {
				buf.WriteRune(c)
			}
		case 'r':
			if isEscape {
				isEscape = false
				buf.WriteRune('\r')
			} else {
				buf.WriteRune(c)
			}
		case 'n':
			if isEscape {
				isEscape = false
				buf.WriteRune('\n')
			} else {
				buf.WriteRune(c)
			}
		case '"':
			if !isEscape {
				return buf.String(), txt[idx+1:]
			}

			isEscape = false
			buf.WriteRune('"')
		default:
			buf.WriteRune(c)
		}
	}
	return buf.String(), ""
}

func skipWhitespace(txt string) string {
	for idx, c := range txt {
		if !unicode.IsSpace(c) {
			return txt[idx:]
		}
	}
	return ""
}

func readString(txt string, breakIfEqualChar bool) (string, string) {
	var buf bytes.Buffer
	for idx, c := range skipWhitespace(txt) {
		switch c {
		case '"':
			if buf.Len() == 0 {
				return readQouteString(txt[idx+1:])
			}
			buf.WriteRune(c)
		case '#':
			if buf.Len() == 0 {
				return "", ""
			}
			buf.WriteRune(c)
		case '=':
			if breakIfEqualChar {
				return buf.String(), txt[idx:]
			}
			buf.WriteRune(c)
		case '\\':
			return "", ""
		default:
			if unicode.IsSpace(c) {
				return buf.String(), txt[idx+1:]
			}

			buf.WriteRune(c)
		}
	}

	return buf.String(), ""
}

func readEqualChar(txt string) (bool, string) {
	for idx, c := range txt {
		if unicode.IsSpace(c) {
			continue
		}
		if c == '=' {
			return true, txt[idx+1:]
		}
		return false, ""
	}
	return false, ""
}

func ReadProperties(nm string) (map[string]string, error) {
	f, e := os.Open(nm)
	if nil != e {
		return nil, e
	}
	defer f.Close()
	return Read(f)
}

func Read(r io.Reader) (map[string]string, error) {
	cfg := map[string]string{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		txt := scanner.Text()
		key, retain := readString(txt, true)
		if len(key) == 0 {
			continue
		}

		hasEqualChar, retain := readEqualChar(retain)
		if !hasEqualChar {
			continue
		}

		value, retain := readString(retain, false)
		if len(value) == 0 {
			continue
		}
		if skipWhitespace(retain) != "" {
			continue
		}
		cfg[key] = os.ExpandEnv(value)
	}

	return expandAll(cfg), nil
}

func expandAll(cfg map[string]string) map[string]string {
	remain := 0
	expend := func(key string) string {
		if value, ok := cfg[key]; ok {
			return value
		}
		remain++
		return key
	}

	for i := 0; i < 100; i++ {
		for k, v := range cfg {
			cfg[k] = os.Expand(v, expend)
		}
		if 0 == remain {
			break
		}
	}
	return cfg
}

func WriteWith(w io.Writer, values map[string]string) error {
	var err error
	for k, v := range values {
		io.WriteString(w, k)
		io.WriteString(w, "=")
		io.WriteString(w, v)
		_, err = io.WriteString(w, "\r\n")
	}
	return err
}

func WriteProperties(nm string, values map[string]string) error {
	if len(values) == 0 {
		return nil
	}
	f, e := os.Create(nm)
	if nil != e {
		return e
	}
	defer f.Close()
	return WriteWith(f, values)
}

func UpdateWith(r io.Reader, w io.Writer, updated map[string]string) error {
	updatedCopy := map[string]string{}
	for k, v := range updated {
		updatedCopy[k] = v
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		txt := scanner.Text()

		for k, v := range updated {
			if strings.Contains(txt, k) {
				ss := strings.SplitN(txt, "=", 2)
				if 2 == len(ss) {
					key := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(ss[0]), "#"))
					if key == k {
						if ss = strings.SplitN(ss[1], "#", 2); 2 == len(ss) {
							txt = k + "=" + v + " #" + ss[1]
						} else {
							txt = k + "=" + v
						}
						delete(updatedCopy, k)
						break
					}
				}
			}
		}
		io.WriteString(w, txt)
		io.WriteString(w, "\r\n")
	}

	if scanner.Err() != nil {
		return scanner.Err()
	}

	var err error
	for k, v := range updatedCopy {
		io.WriteString(w, k)
		io.WriteString(w, "=")
		io.WriteString(w, v)
		_, err = io.WriteString(w, "\r\n")
	}
	return err
}

func UpdateProperties(nm string, updated map[string]string) error {
	if len(updated) == 0 {
		return nil
	}
	f, e := os.Open(nm)
	if nil != e {
		return e
	}
	defer f.Close()

	out, e := os.Create(nm + ".tmp")
	if nil != e {
		return e
	}
	defer out.Close()

	if e := UpdateWith(f, out, updated); nil != e {
		return e
	}

	if e := out.Close(); nil != e {
		return e
	}
	if e := f.Close(); nil != e {
		return e
	}
	if e := os.Remove(nm); nil != e {
		return e
	}
	if e := os.Rename(nm+".tmp", nm); nil != e {
		return e
	}
	return nil
}
