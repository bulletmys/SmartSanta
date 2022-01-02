package algorithm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCountPreferences(t *testing.T) {
	type args struct {
		nodes map[uint64][]uint64
	}
	tests := []struct {
		name string
		args args
		want map[uint64]uint64
	}{
		{
			name: "t1",
			args: args{nodes: map[uint64][]uint64{
				1: {2, 3},
				2: {1, 5},
				3: {1, 4},
				4: {3, 5},
				5: {1},
			}},
			want: map[uint64]uint64{
				1: 2,
				2: 5,
				5: 1,
				3: 4,
				4: 3,
			},
		},
		{
			name: "t2",
			args: args{nodes: map[uint64][]uint64{
				1: {2, 7},
				2: {1, 3},
				3: {2},
				4: {3},
				5: {4},
				6: {5},
				7: {6},
			}},
			want: map[uint64]uint64{
				1: 7,
				2: 1,
				3: 2,
				4: 3,
				5: 4,
				6: 5,
				7: 6,
			},
		},
		{
			name: "t3",
			args: args{nodes: map[uint64][]uint64{
				1: {2},
				2: {3},
				3: {4},
				4: {1, 5},
				5: {6},
				6: {1},
			}},
			want: map[uint64]uint64{
				1: 2,
				2: 3,
				3: 4,
				4: 5,
				5: 6,
				6: 1,
			},
		},
		{
			name: "t4",
			args: args{nodes: map[uint64][]uint64{
				1: {2},
				2: {3},
				3: {4},
				4: {1, 5},
				5: {6},
				6: {1},
			}},
			want: map[uint64]uint64{
				1: 2,
				2: 3,
				3: 4,
				4: 5,
				5: 6,
				6: 1,
			},
		},
		{
			name: "t5",
			args: args{nodes: map[uint64][]uint64{
				1: {2},
				2: {3},
				4: {5},
				5: {6},
				6: {1, 7},
			}},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CountPreferences(tt.args.nodes)
			assert.Equal(t, tt.want, got)
		})
	}
}
