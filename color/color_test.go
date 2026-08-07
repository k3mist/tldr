package color

import "testing"

func TestColor(t *testing.T) {
	cases := []struct {
		name  string
		code  int
		flags []int
		want  string
	}{
		{"no flags", Red, nil, "\033[;31m"},
		{"one flag", Red, []int{Bold}, "\033[1;31m"},
		{"multiple flags", Red, []int{Bold, Underline}, "\033[1;4;31m"},
	}

	for _, c := range cases {
		if got := Color(c.code, c.flags...); got != c.want {
			t.Errorf("%s: Color(%d, %v) = %q, want %q", c.name, c.code, c.flags, got, c.want)
		}
	}
}

func TestColorBold(t *testing.T) {
	cases := []struct {
		code int
		want string
	}{
		{Red, Color(Red, Bold)},
		{Green, Color(Green, Bold)},
	}

	for _, c := range cases {
		if got := ColorBold(c.code); got != c.want {
			t.Errorf("ColorBold(%d) = %q, want %q", c.code, got, c.want)
		}
	}
}
