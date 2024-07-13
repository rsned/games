package mnkgame

import "testing"

func TestOutcome(t *testing.T) {
	tests := []struct {
		have Outcome
		want string
	}{
		{
			have: OutcomeIncomplete,
			want: "Game Unfinished",
		},
		{
			have: 173,
			want: "Game Unfinished",
		},
	}

	for _, test := range tests {
		if got := test.have.String(); got != test.want {
			t.Errorf("%v.String() = %v, want %s", test.have, got, test.want)
		}
	}
}
