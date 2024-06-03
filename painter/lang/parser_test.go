package lang

import (
	"github.com/KPI-IM-24/go_lab3/painter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name      string
		command   string
		op        painter.Operation
		wantErr   bool
		operation bool
	}{
		{"Test Background Rectangle Command", "bgrect 0.25 0.25 0.75 0.75", &painter.BgRect{X1: 200, Y1: 200, X2: 600, Y2: 600}, false, true},
		{"Test Figure Command", "figure 0.25 0.25", &painter.TFigure{X: 200, Y: 200}, false, true},
		{"Test Move Command", "move 0.125 0.125", &painter.MoveFigures{X: 100, Y: 100}, false, true},
		{"Test Update Command", "update", painter.UpdateOp, false, true},
		{"Test Invalid Command", "invalidcommand", nil, true, true},
		{"Test White Fill Command", "white", painter.OperationFunc(painter.WhiteFill), false, false},
		{"Test Green Fill Command", "green", painter.OperationFunc(painter.GreenFill), false, false},
		{"Test Reset Screen Command", "reset", painter.OperationFunc(painter.Reset), false, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := &Parser{}
			ops, err := parser.Parse(strings.NewReader(tc.command))

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				if tc.operation {
					require.Len(t, ops, 2)
					assert.IsType(t, tc.op, ops[1])
					assert.Equal(t, tc.op, ops[1])
				} else {
					require.Len(t, ops, 1)
					assert.IsType(t, tc.op, ops[0])
				}
			}
		})
	}
}
