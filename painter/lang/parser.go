package lang

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/KPI-IM-24/go_lab3/painter"
	"io"
)

const screenSize = 800

type Parser struct {
	bgColor  painter.Operation
	bgRect   *painter.BgRect
	figures  []*painter.TFigure
	moveOps  []painter.Operation
	updateOp painter.Operation
}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	p.initialize()
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		commandLine := scanner.Text()
		if err := p.parse(commandLine); err != nil {
			return nil, err
		}
	}
	return p.finalResult(), nil
}

func (p *Parser) parse(commandLine string) error {
	tokens := strings.Split(commandLine, " ")
	if len(tokens) < 1 {
		return fmt.Errorf("invalid command format: %s", commandLine)
	}

	instruction := tokens[0]
	args, err := toIntArgs(tokens, screenSize)
	if err != nil {
		return fmt.Errorf("invalid argument format for %s: %s", instruction, commandLine)
	}

	switch instruction {
	case "white":
		p.bgColor = painter.OperationFunc(painter.WhiteFill)
	case "green":
		p.bgColor = painter.OperationFunc(painter.GreenFill)
	case "bgrect":
		if len(args) != 4 {
			return fmt.Errorf("invalid number of arguments for bgrect: %s", commandLine)
		}
		p.bgRect = &painter.BgRect{X1: args[0], Y1: args[1], X2: args[2], Y2: args[3]}
	case "figure":
		if len(args) != 2 {
			return fmt.Errorf("invalid number of arguments for figure: %s", commandLine)
		}
		p.figures = append(p.figures, &painter.TFigure{X: args[0], Y: args[1]})
	case "move":
		if len(args) != 2 {
			return fmt.Errorf("invalid number of arguments for move: %s", commandLine)
		}
		p.moveOps = append(p.moveOps, &painter.MoveFigures{X: args[0], Y: args[1], Figures: p.figures})
	case "reset":
		p.resetState()
		p.bgColor = painter.OperationFunc(painter.Reset)
	case "update":
		p.updateOp = painter.UpdateOp
	default:
		return fmt.Errorf("unknown command: %s", instruction)
	}
	return nil
}

func toIntArgs(tokens []string, screenSize int) ([]int, error) {
	args := make([]int, 0, len(tokens)-1)
	for _, arg := range tokens[1:] {
		val, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			return nil, err
		}
		args = append(args, int(val*float64(screenSize)))
	}
	return args, nil
}

func (p *Parser) finalResult() []painter.Operation {
	var result []painter.Operation
	result = append(result, p.bgColor)
	if p.bgRect != nil {
		result = append(result, p.bgRect)
	}
	result = append(result, p.moveOps...)
	p.moveOps = nil // Clear accumulated move operations
	if len(p.figures) > 0 {
		for _, figure := range p.figures {
			result = append(result, figure)
		}
	}
	if p.updateOp != nil {
		result = append(result, p.updateOp)
	}
	return result
}

func (p *Parser) initialize() {
	if p.bgColor == nil {
		p.bgColor = painter.OperationFunc(painter.Reset)
	}
	p.updateOp = nil
}

func (p *Parser) resetState() {
	p.bgRect = nil
	p.bgRect = nil
	p.figures = nil
	p.moveOps = nil
	p.updateOp = nil
}
