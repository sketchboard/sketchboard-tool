package convertcsv

import (
	"reflect"
	"testing"
)

func TestDiagramItem_ParseShape(t *testing.T) {
	type fields struct {
		Text            string
		ElementType     string
		Shape           string
		BackgroundColor string
		TextColor       string
		Props           uint64
		Version         uint64
		ID              uint64
		ClientID        string
		CustomData      string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *ShapeDimension
		wantErr bool
	}{
		{
			name: "t1",
			fields: fields{
				Shape: "1,2,50,10",
			},
			want: &ShapeDimension{
				X:      1,
				Y:      2,
				Width:  50,
				Height: 10,
			},
			wantErr: false,
		},
		{
			name: "negative x and y",
			fields: fields{
				Shape: "-1,-20,50,10",
			},
			want: &ShapeDimension{
				X:      -1,
				Y:      -20,
				Width:  50,
				Height: 10,
			},
			wantErr: false,
		},
		{
			name: "negative y",
			fields: fields{
				Shape: "10,-20,50,10",
			},
			want: &ShapeDimension{
				X:      10,
				Y:      -20,
				Width:  50,
				Height: 10,
			},
			wantErr: false,
		},
		{
			name: "negative width and height",
			fields: fields{
				Shape: "-1,-20,-50,-10",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			di := &DiagramItem{
				Text:            tt.fields.Text,
				ElementType:     tt.fields.ElementType,
				Shape:           tt.fields.Shape,
				BackgroundColor: tt.fields.BackgroundColor,
				TextColor:       tt.fields.TextColor,
				Props:           tt.fields.Props,
				Version:         tt.fields.Version,
				ID:              tt.fields.ID,
				ClientID:        tt.fields.ClientID,
				CustomData:      tt.fields.CustomData,
			}
			got, err := di.ParseShape()
			if (err != nil) != tt.wantErr {
				t.Errorf("DiagramItem.ParseShape() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DiagramItem.ParseShape() = %v, want %v", got, tt.want)
			}
		})
	}
}
