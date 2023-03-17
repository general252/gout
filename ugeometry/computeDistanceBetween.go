package ugeometry

import (
	"github.com/tidwall/geodesic"
	"math"
)

type LatLng struct {
	Lat float64
	Lng float64
}

// ComputeDistanceBetween 用于返回两个 LatLng 之间的距离（以米为单位）。您可以选择指定自定义半径。半径默认为地球的半径(6378137米)。
func ComputeDistanceBetween(from, to LatLng, radius ...float64) float64 {
	h := newComputeDistance()
	if len(radius) > 0 {
		h.SetRadius(radius[0])
	}

	return h.GetDistance(from, to, computeTypeGoogleMap)
}

type computeDistance struct {
	radius float64 // 地球半径，单位m
}

func newComputeDistance() *computeDistance {
	return &computeDistance{
		radius: 6378137,
	}
}

// SetRadius 指定地球半径, 单位米, 默认: 6378137 米
func (tis *computeDistance) SetRadius(radius float64) *computeDistance {
	tis.radius = radius
	return tis
}

func (tis *computeDistance) toRadians(d float64) float64 {
	return d * math.Pi / 180.0
}

func (tis *computeDistance) GetDistance(from, to LatLng, t computeType) float64 {
	switch t {
	case computeTypeAsin:
		return tis.getDistanceAsin(from.Lat, from.Lng, to.Lat, to.Lng)
	case computeTypeGoogleMap:
		return tis.getDistanceGoogleMap(from.Lat, from.Lng, to.Lat, to.Lng)
	case computeTypeAcos:
		return tis.getDistanceAcos(from.Lat, from.Lng, to.Lat, to.Lng)
	case computeTypeLeaflet:
		return tis.getDistanceLeaflet(from.Lat, from.Lng, to.Lat, to.Lng)
	case computeTypeGlobe:
		return tis.getDistanceGlobe(from.Lat, from.Lng, to.Lat, to.Lng)
	case computeTypeBaidu:
		return tis.getDistanceBaidu(from.Lat, from.Lng, to.Lat, to.Lng)
	case computeTypeGaoDe:
		return tis.getDistanceGaoDe(from.Lat, from.Lng, to.Lat, to.Lng)
	default:
		return -1
	}
}

// getDistanceAsin 反正弦计算方式
func (tis *computeDistance) getDistanceAsin(latitude1, longitude1 float64, latitude2, longitude2 float64) float64 {

	// 纬度
	var lat1 = tis.toRadians(latitude1)
	var lat2 = tis.toRadians(latitude2)
	// 经度
	var lng1 = tis.toRadians(longitude1)
	var lng2 = tis.toRadians(longitude2)
	// 纬度之差
	var a = lat1 - lat2
	// 经度之差
	var b = lng1 - lng2
	// 计算两点距离的公式
	var s = 2 * math.Asin(math.Sqrt(math.Pow(math.Sin(a/2), 2)+
		math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(b/2), 2)))
	// 弧长乘地球半径, 返回单位: 米
	s = s * tis.radius

	return s
}

// getDistanceGoogleMap 基于googleMap中的算法得到两经纬度之间的距离, 计算精度与谷歌地图的距离精度差不多。 反正弦
func (tis *computeDistance) getDistanceGoogleMap(lat1, lon1 float64, lat2, lon2 float64) float64 {

	var radLat1 = tis.toRadians(lat1)
	var radLat2 = tis.toRadians(lat2)
	var a = radLat1 - radLat2
	var b = tis.toRadians(lon1) - tis.toRadians(lon2)
	var s = 2 * math.Asin(math.Sqrt(math.Pow(math.Sin(a/2), 2)+math.Cos(radLat1)*math.Cos(radLat2)*math.Pow(math.Sin(b/2), 2)))
	s = s * tis.radius
	return s
}

// getDistanceAcos 反余弦计算方式
func (tis *computeDistance) getDistanceAcos(lat1, lng1 float64, lat2, lng2 float64) float64 {

	// 经纬度（角度）转弧度。弧度用作参数，以调用Math.cos和Math.sin
	var radiansAX = tis.toRadians(lng1) // A经弧度
	var radiansAY = tis.toRadians(lat1) // A纬弧度
	var radiansBX = tis.toRadians(lng2) // B经弧度
	var radiansBY = tis.toRadians(lat2) // B纬弧度

	// 公式中“cosβ1cosβ2cos（α1-α2）+sinβ1sinβ2”的部分，得到∠AOB的cos值
	var cos = math.Cos(radiansAY)*math.Cos(radiansBY)*math.Cos(radiansAX-radiansBX) + math.Sin(radiansAY)*math.Sin(radiansBY)
	//        System.out.println("cos = " + cos); // 值域[-1,1]
	var acos = math.Acos(cos) // 反余弦值
	//        System.out.println("acos = " + acos); // 值域[0,π]
	//        System.out.println("∠AOB = " + Math.toDegrees(acos)); // 球心角 值域[0,180]
	return tis.radius * acos // 最终结果
}

// getDistanceLeaflet Leaflet https://makinacorpus.github.io/Leaflet.MeasureControl/
func (tis *computeDistance) getDistanceLeaflet(lat1, lng1 float64, lat2, lng2 float64) float64 {
	const DEG_TO_RAD = math.Pi / 180.0
	var e float64 = tis.radius
	var i = DEG_TO_RAD
	var n = (lat2 - lat1) * i
	var s = (lng2 - lng1) * i
	var a = lat1 * i
	var r = lat2 * i
	var h = math.Sin(n / 2)
	var l = math.Sin(s / 2)
	var u = h*h + l*l*math.Cos(a)*math.Cos(r)

	return 2 * e * math.Atan2(math.Sqrt(u), math.Sqrt(1-u))
}

// getDistanceGlobe geodesic.Globe
func (tis *computeDistance) getDistanceGlobe(lat1, lng1 float64, lat2, lng2 float64) float64 {
	var dist float64 = 0
	geodesic.Globe.Inverse(lat1, lng1, lat2, lng2, &dist, nil, nil)

	return dist
}

// getDistanceBaidu 百度
func (tis *computeDistance) getDistanceBaidu(lat1, lng1 float64, lat2, lng2 float64) float64 {
	var i = tis.toRadians(lng1)
	var jP = tis.toRadians(lat1)

	var e = tis.toRadians(lng2)
	var T = tis.toRadians(lat2)

	var jO = jP
	return tis.radius * math.Acos(math.Sin(jO)*math.Sin(T)+math.Cos(jO)*math.Cos(T)*math.Cos(e-i))
}

// getDistanceGaoDe 高德
func (tis *computeDistance) getDistanceGaoDe(lat1, lng1 float64, lat2, lng2 float64) float64 {
	var rr = math.Pi / 180.0
	var n = math.Cos

	var i = lat1 * rr
	var t = lng1 * rr
	var a = lat2 * rr
	var e = lng2 * rr
	var r = 2 * tis.radius
	e = e - t
	t = (1 - n(a-i) + (1-n(e))*n(i)*n(a)) / 2

	return r * math.Asin(math.Sqrt(t))
}

type (
	computeType int
)

const (
	computeTypeAsin = iota
	computeTypeGoogleMap
	computeTypeAcos
	computeTypeLeaflet
	computeTypeGlobe
	computeTypeBaidu
	computeTypeGaoDe
)
