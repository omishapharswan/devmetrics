package analysis

import "testing"

func TestCalculateComplexity(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   int
	}{
		{"empty function", "func f() {}", 1},
		{"single if", "if x { return 1 }", 2},
		{"if-else-if-else", "if x {\n} else if y {\n} else {\n}", 3},
		{"loop and boolean ops", "for i := 0; i < n && ok; i++ {}", 3},
		{"switch with cases", "switch x {\ncase 1:\ncase 2:\ncase 3:\n}", 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateComplexity(tt.source)
			if got != tt.want {
				t.Errorf("CalculateComplexity(%q) = %d, want %d", tt.source, got, tt.want)
			}
		})
	}
}

func TestDetectFunctionsGo(t *testing.T) {
	source := `package main

func simple() int {
	return 1
}

func branchy(x int) int {
	if x > 0 {
		return 1
	} else if x < 0 {
		return -1
	}
	return 0
}
`
	funcs := DetectFunctions(source, "Go")
	if len(funcs) != 2 {
		t.Fatalf("expected 2 functions, got %d: %+v", len(funcs), funcs)
	}
	if funcs[0].Name != "simple" || funcs[0].Complexity != 1 {
		t.Errorf("simple() = %+v, want name=simple complexity=1", funcs[0])
	}
	if funcs[1].Name != "branchy" || funcs[1].Complexity != 3 {
		t.Errorf("branchy() = %+v, want name=branchy complexity=3", funcs[1])
	}
}

func TestDetectFunctionsPython(t *testing.T) {
	source := `def greet(name):
    if name:
        return "hi " + name
    return "hi"

def other():
    return 42
`
	funcs := DetectFunctions(source, "Python")
	if len(funcs) != 2 {
		t.Fatalf("expected 2 functions, got %d: %+v", len(funcs), funcs)
	}
	if funcs[0].Name != "greet" || funcs[0].Complexity != 2 {
		t.Errorf("greet() = %+v, want name=greet complexity=2", funcs[0])
	}
	if funcs[1].Name != "other" || funcs[1].Complexity != 1 {
		t.Errorf("other() = %+v, want name=other complexity=1", funcs[1])
	}
}

func TestDetectFunctionsJSArrow(t *testing.T) {
	source := `const add = (a, b) => {
  return a + b;
};

const isPositive = (n) => {
  if (n > 0) {
    return true;
  }
  return false;
};
`
	funcs := DetectFunctions(source, "JavaScript")
	if len(funcs) != 2 {
		t.Fatalf("expected 2 functions, got %d: %+v", len(funcs), funcs)
	}
	if funcs[0].Name != "add" || funcs[0].Complexity != 1 {
		t.Errorf("add() = %+v, want name=add complexity=1", funcs[0])
	}
	if funcs[1].Name != "isPositive" || funcs[1].Complexity != 2 {
		t.Errorf("isPositive() = %+v, want name=isPositive complexity=2", funcs[1])
	}
}
