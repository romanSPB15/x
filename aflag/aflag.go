package aflag

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var source []string = os.Args

func SetArgs(args []string) {
	source = args
}

func arg(n int) string {
	if n < 0 || n >= len(source) {
		return ""
	}
	return source[n+1]
}

type Flags map[string]Flag

type Flag struct {
	Short    rune
	Optional bool
}

type Topic struct {
	Name   string
	F      func(*Topic)
	Flags  Flags
	values map[string]string
}

func (t *Topic) String(name string) string {
	return t.values[name]
}

func (t *Topic) Int(name string) (int, error) {
	return strconv.Atoi(t.values[name])
}

func (t *Topic) Bool(name string) bool {
	return t.values[name] != ""
}

type App struct {
	Topics []Topic
}

func parseTopic(t *Topic) string {
	fmt.Printf("parse '%s'\r\n", t.Name)
	if t.F != nil {
		defer t.F(t)
	}
	if len(source) < 3 {
		return ""
	}
	t.values = make(map[string]string)
	a := source[2:]

	n := ""

	checkValueOneArg := func(parts []string, name string, v string) string {
		switch len(parts) {
		case 1:
			return ""
		case 2:
			t.values[name] = parts[1]
			return ""
		}
		return fmt.Sprintf("invalid flag: %s\r\n", v)
	}

	for _, v := range a {
		long := strings.TrimPrefix(v, "--")

		if long != v {
			parts := strings.Split(long, "=")

			r := false

			for name := range t.Flags {
				if name == parts[0] {
					n = name
					t.values[name] = "true"
					r = true
					break
				}
			}

			if !r {
				return fmt.Sprintf("unknown flag: %s\r\n", v)
			}

			if err := checkValueOneArg(parts, n, v); err != "" {
				return err
			}

			if r {
				continue
			}
		}
		short := strings.TrimPrefix(v, "-")
		if short != v {
			parts := strings.Split(short, "=")
			r := false
			runes := []rune(parts[0])
			if len(runes) > 1 {
				return fmt.Sprintf("invalid short form flag: %s\r\n", v)
			}

			var candidates []string
			for name, f := range t.Flags {
				if f.Short != 0 && f.Short == runes[0] {

					n = name
					t.values[name] = "true"
					r = true
					break
				}
				if f.Short == 0 && []rune(name)[0] == runes[0] {
					candidates = append(candidates, name)
				}
			}
			if !r {
				if len(candidates) == 1 {
					name := candidates[0]
					n = name
					t.values[name] = "true"
					r = true
				} else if len(candidates) > 1 {
					return fmt.Sprintf("ambiguous short flag -%c matches: %v\r\n", runes[0], candidates)
				} else {
					return fmt.Sprintf("unknown flag: %s\r\n", v)
				}
			}

			if !r {
				return fmt.Sprintf("unknown flag: %s\r\n", v)
			}

			if err := checkValueOneArg(parts, n, v); err != "" {
				return err
			}

			if r {
				continue
			}
		}

		if n != "" {
			t.values[n] = v
			n = ""
		}
	}
	return ""
}

func (u *App) Parse() {
	if len(u.Topics) == 0 {
		return
	}
	z := (*Topic)(nil)
	zFlag := true
	for i, t := range u.Topics {
		if t.Name == "" {
			z = &u.Topics[i]
		} else if t.Name == arg(0) {
			zFlag = false
			fmt.Fprint(os.Stderr, parseTopic(&u.Topics[i]))
		}
	}
	if z != nil && zFlag {
		fmt.Fprint(os.Stderr, parseTopic(z))
	}
}

func (u *App) T(name string) *Topic {
	for i, v := range u.Topics {
		if v.Name == name {
			return &u.Topics[i]
		}
	}
	return nil
}
