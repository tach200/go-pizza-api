package dominos

// func TestFormatProductData(t *testing.T) {
// 	type args struct {
// 		dealContent []DealContent
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want []ranking.Product
// 	}{
// 		{
// 			name: "test lol",
// 			args: args{
// 				dealContent: []DealContent{
// 					{
// 						Product: "/Content/Images/Site/icons/pizza.png",
// 					},
// 					{
// 						Product: "/Content/Images/Site/icons/pizza.png",
// 					},
// 					{
// 						Product: "/Content/Images/Site/icons/garlic-bread.png",
// 					},
// 					{
// 						Product: "/Content/Images/Site/icons/garlic-bread.png",
// 					},
// 				},
// 			},
// 			want: []ranking.Product{
// 				{
// 					ProductType:  "pizza",
// 					ProductCount: 2,
// 				},
// 				{
// 					ProductType:  "garlic-bread",
// 					ProductCount: 2,
// 				},
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := FormatProductData(tt.args.dealContent); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("FormatProductData() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
