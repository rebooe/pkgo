// 坐标系相互转换方法
// - WGS84 世界坐标系
// - GCJ02 火星坐标系
// - BD-09 百度坐标系
package gis

import (
	"math"
	"strconv"
)

const (
	pi  = math.Pi                // 圆周率
	xpi = pi * 3000.0 / 180.0    // 圆周率对应的经纬度偏移
	a   = 6378245.0              // 长半轴
	ee  = 0.00669342162296594323 // 扁率
)

func transformLat(x, y float64) float64 {
	ret := -100.0 + 2.0*x + 3.0*y + 0.2*y*y + 0.1*x*y + 0.2*math.Sqrt(math.Abs(x))
	ret += (20.0*math.Sin(6.0*x*pi) + 20.0*math.Sin(2.0*x*pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(y*pi) + 40.0*math.Sin(y/3.0*pi)) * 2.0 / 3.0
	ret += (160.0*math.Sin(y/12.0*pi) + 320*math.Sin(y*pi/30.0)) * 2.0 / 3.0
	return ret
}

func transformlng(x, y float64) float64 {
	ret := 300.0 + x + 2.0*y + 0.1*x*x + 0.1*x*y + 0.1*math.Sqrt(math.Abs(x))
	ret += (20.0*math.Sin(6.0*x*pi) + 20.0*math.Sin(2.0*x*pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(x*pi) + 40.0*math.Sin(x/3.0*pi)) * 2.0 / 3.0
	ret += (150.0*math.Sin(x/12.0*pi) + 300.0*math.Sin(x/30.0*pi)) * 2.0 / 3.0
	return ret
}

func outOfChina(lat, lng float64) bool {
	if lng < 72.004 || lng > 137.8347 {
		return true
	}
	if lat < 0.8293 || lat > 55.8271 {
		return true
	}
	return false
}

func transform(lat, lng float64) []float64 {
	if outOfChina(lat, lng) {
		return []float64{lat, lng}
	}
	dLat := transformLat(lng-105.0, lat-35.0)
	dlng := transformlng(lng-105.0, lat-35.0)
	radLat := lat / 180.0 * pi
	magic := math.Sin(radLat)
	magic = 1 - ee*magic*magic
	SqrtMagic := math.Sqrt(magic)
	dLat = (dLat * 180.0) / ((a * (1 - ee)) / (magic * SqrtMagic) * pi)
	dlng = (dlng * 180.0) / (a / SqrtMagic * math.Cos(radLat) * pi)
	mgLat := lat + dLat
	mglng := lng + dlng
	return []float64{mgLat, mglng}
}

func WGS84ToGCJ02(lat, lng float64) (float64, float64) {
	if outOfChina(lat, lng) {
		return lat, lng
	}
	dLat := transformLat(lng-105.0, lat-35.0)
	dlng := transformlng(lng-105.0, lat-35.0)
	radLat := lat / 180.0 * pi
	magic := math.Sin(radLat)
	magic = 1 - ee*magic*magic
	SqrtMagic := math.Sqrt(magic)
	dLat = (dLat * 180.0) / ((a * (1 - ee)) / (magic * SqrtMagic) * pi)
	dlng = (dlng * 180.0) / (a / SqrtMagic * math.Cos(radLat) * pi)
	mgLat := lat + dLat
	mglng := lng + dlng
	return mgLat, mglng
}

func GCJ02ToWGS84(lat, lng float64) (float64, float64) {
	gps := transform(lat, lng)
	lngtitude := lng*2 - gps[1]
	latitude := lat*2 - gps[0]
	return latitude, lngtitude
}

func GCJ02ToBD09(lat, lng float64) (float64, float64) {
	x := lng
	y := lat
	z := math.Sqrt(x*x+y*y) + 0.00002*math.Sin(y*xpi)
	theta := math.Atan2(y, x) + 0.000003*math.Cos(x*xpi)
	templng := z*math.Cos(theta) + 0.0065
	tempLat := z*math.Sin(theta) + 0.006
	return tempLat, templng
}

func BD09ToGCJ02(lat, lng float64) (float64, float64) {
	x := lng - 0.0065
	y := lat - 0.006
	z := math.Sqrt(x*x+y*y) - 0.00002*math.Sin(y*xpi)
	theta := math.Atan2(y, x) - 0.000003*math.Cos(x*xpi)
	templng := z * math.Cos(theta)
	tempLat := z * math.Sin(theta)
	return tempLat, templng
}

func WGS84ToBD09(lat, lng float64) (float64, float64) {
	gcj02Lat, gcj02Lng := WGS84ToGCJ02(lat, lng)
	bd09Lat, bd09Lng := GCJ02ToBD09(gcj02Lat, gcj02Lng)
	return bd09Lat, bd09Lng
}

func BD09ToWGS84(lat, lng float64) (float64, float64) {
	gcj02Lat, gcj02Lng := BD09ToGCJ02(lat, lng)
	WGS84Lat, WGS84Lng := GCJ02ToWGS84(gcj02Lat, gcj02Lng)
	//保留小数点后六位
	WGS84Lat = retain6(WGS84Lat)
	WGS84Lng = retain6(WGS84Lng)
	return WGS84Lat, WGS84Lng
}

// 保留小数点后六位
func retain6(num float64) float64 {
	value, _ := strconv.ParseFloat(strconv.FormatFloat(num, 'f', 6, 64), 64)
	return value
}
