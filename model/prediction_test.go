package model

import (
	"fmt"
	"testing"
)

func TestSEIRPrediction(t *testing.T) {
	shanghai := []int{10, 14, 22, 36, 41, 68, 80, 91, 111, 114, 139, 168, 191, 212, 228, 253, 274, 297, 315, 326, 337, 342, 352, 366}
	var alpha0 = 0.8
	var m_list = []float64{6.0, 3.0, 2.0}
	var k = 0.5
	cityData := []float64{}
	for _, data := range shanghai {
		cityData = append(cityData, float64(data))
	}
	rs, delta := SEIRPrediction(cityData, 60, k, m_list, alpha0)
	fmt.Println(rs)

	t.Error(rs, delta)
}
