package deals

import (
	"reflect"
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

func Test_convertToScoreArr(t *testing.T) {
	type args struct {
		keywords   []string
		pizzaSizes map[string]float64
	}
	tests := []struct {
		name    string
		args    args
		want    []float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToScoreArr(tt.args.keywords, tt.args.pizzaSizes)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToScoreArr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertToScoreArr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getDealKeywords(t *testing.T) {
	type args struct {
		dealDesc string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "extract correct keywords from the deal",
			args: args{
				dealDesc: "X1 Large Stuffed Crust Pizzas And A Classic Side",
			},
			want:    []string{"1", "Large"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getDealKeywords(tt.args.dealDesc)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDealKeywords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDealKeywords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_rankScore(t *testing.T) {
	type args struct {
		dealTitle  string
		dealDesc   string
		pizzaSizes map[string]float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
		{
			name: "test result for pizza hut",
			args: args{
				dealTitle:  "£19.99 Stuffed Crust Deal",
				dealDesc:   "X1 Large Stuffed Crust Deal And A Classic Side",
				pizzaSizes: pizzahutSizes,
			},
			want: 0.7003501750875438,
		},
		{
			name: "pizzahut",
			args: args{
				dealTitle:  "£25.99 Deal",
				dealDesc:   "2 Large Pizzas, 2 Classic Sides And A 1.5 Litre Drink",
				pizzaSizes: pizzahutSizes,
			},
			want: 1.0773374374759523,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rankScore(tt.args.dealTitle, tt.args.dealDesc, tt.args.pizzaSizes); got != tt.want {
				t.Errorf("rankScore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculateDiscount(t *testing.T) {
	type args struct {
		dealPercentage float64
		dealCost       float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
		{
			name: "calculate discount",
			args: args{
				dealPercentage: 50,
				dealCost:       100,
			},
			want: 50,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateDiscount(tt.args.dealPercentage, tt.args.dealCost); got != tt.want {
				t.Errorf("calculateDiscount() = %v, want %v", got, tt.want)
			}
		})
	}
}
