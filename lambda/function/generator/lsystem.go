package generator

type Lsystem struct {
	Variables []rune
	Constants []rune
	Axiom     string
	Rules     []Rule
	Angle     float64
}

type Rule struct {
	In  string
	Out string
}

func Iterate(lsystem *Lsystem, limit int) string {
	output := lsystem.Axiom
	for i := 0; i < limit; i++ {
		output = IterateOnce(lsystem, output)
	}
	return output
}

func IterateOnce(lsystem *Lsystem, axiom string) string {
	var output string
	for _, c := range axiom {
		tmp := string(c)
		isConstant := false
		for _, constant := range lsystem.Constants {
			if tmp == string(constant) {
				isConstant = true
				break
			}
		}
		if isConstant {
			output += string(c)
		} else {
			for _, rule := range lsystem.Rules {
				if rule.In == tmp {
					output += rule.Out
					break
				}
			}
		}
	}
	return output
}

func Process(lsystem string, operations map[rune]func()) error {
	for _, c := range []byte(lsystem) {
		match := operations[rune(c)]
		if match == nil {
			continue
		}
		match()
	}
	return nil
}
