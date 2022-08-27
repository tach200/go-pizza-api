package deals

import (
	"testing"
)

func Test_getDealCost(t *testing.T) {
	type args struct {
		dealTitle string
		dealDesc  string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "extract cost from the desc",
			args: args{
				dealTitle: "2 Medium Pizzas & 1 Side For £20.99",
				dealDesc:  "2 Medium Pizzas And A Classic Side",
			},
			want:    20.99,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getDealCost(tt.args.dealTitle, tt.args.dealDesc)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDealCost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getDealCost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculateScoreArr(t *testing.T) {
	type args struct {
		scoreArr []float64
		dealCost float64
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "calculate scores",
			args: args{
				scoreArr: []float64{2, 9.5},
				dealCost: 10.00,
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "calculate scores",
			args: args{
				scoreArr: []float64{2, 13},
				dealCost: 10.00,
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculateScoreArr(tt.args.scoreArr, tt.args.dealCost)
			if (err != nil) != tt.wantErr {
				t.Errorf("calculateScoreArr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("calculateScoreArr() = %v, want %v", got, tt.want)
			}
		})
	}
}
