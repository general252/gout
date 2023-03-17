package ugeometry

import (
	"fmt"
	"testing"
)

func Test_computeDistance_GetDistance(t *testing.T) {

	var show = func(from, to LatLng) {
		h := newComputeDistance()
		fmt.Printf(">>>>>>>>>>>>>>>>>   : (%v, %v) (%v, %v)\n", from.Lat, from.Lng, to.Lat, to.Lng)
		fmt.Printf("computeTypeAsin     : %v\n", h.GetDistance(from, to, computeTypeAsin))
		fmt.Printf("computeTypeGoogleMap: %v\n", h.GetDistance(from, to, computeTypeGoogleMap))
		fmt.Printf("computeTypeAcos     : %v\n", h.GetDistance(from, to, computeTypeAcos))
		fmt.Printf("computeTypeLeaflet  : %v\n", h.GetDistance(from, to, computeTypeLeaflet))
		fmt.Printf("computeTypeGlobe    : %v\n", h.GetDistance(from, to, computeTypeGlobe))
		fmt.Printf("computeTypeBaidu    : %v\n", h.GetDistance(from, to, computeTypeBaidu))
		fmt.Printf("computeTypeGaoDe    : %v\n", h.GetDistance(from, to, computeTypeGaoDe))
		fmt.Println()
	}

	show(LatLng{
		Lat: 31.81981493276764,
		Lng: 117.20793509299229,
	}, LatLng{
		Lat: 31.856604164072078,
		Lng: 117.21616301031396,
	})

	show(LatLng{
		Lat: 33.4911,
		Lng: -112.4223,
	}, LatLng{
		Lat: 32.1189,
		Lng: -113.1123,
	})

	dis := ComputeDistanceBetween(LatLng{
		Lat: 33.4911,
		Lng: -112.4223,
	}, LatLng{
		Lat: 32.1189,
		Lng: -113.1123,
	})
	fmt.Printf("dis: %v\n", dis)
}
