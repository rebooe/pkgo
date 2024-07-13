package gis

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
)

var pointRe = regexp.MustCompile(`POINT.+\((.*)\ (.*)\)`)

type Geo struct {
	Regions []Region    // 区域信息
	indexID map[int]int // 索引
}

func NewGeos() *Geo {
	return &Geo{
		Regions: make([]Region, 0),
		indexID: make(map[int]int),
	}
}

type Region struct {
	ID          int
	PID         int       // 父id
	Deep        int       // 层级
	Name        string    // 名称
	ExtPath     string    // 全称
	GeoWkt      Point     // 中心点
	Coordinates []Polygon // 边界区域
}

// 加载 geoJson 格式文件
func (g *Geo) LoadGeoJson(path string) error {
	body, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	collection := geojson.FeatureCollection{}
	if err := collection.UnmarshalJSON(body); err != nil {
		return err
	}

	for _, feature := range collection.Features {
		if err := g.addGeom(feature); err != nil {
			return err
		}
	}
	return nil
}

func (g *Geo) addGeom(feature *geojson.Feature) error {
	id, err := strconv.Atoi(fmt.Sprintf("%s", feature.Properties["id"]))
	if err != nil {
		return err
	}

	deep, err := strconv.Atoi(fmt.Sprintf("%s", feature.Properties["deep"]))
	if err != nil {
		return err
	}

	pid, err := strconv.Atoi(fmt.Sprintf("%s", feature.Properties["pid"]))
	if err != nil {
		return err
	}

	name := fmt.Sprintf("%s", feature.Properties["name"])

	extPath := fmt.Sprintf("%s", feature.Properties["ext_path"])

	geoWkt := fmt.Sprintf("%s", feature.Properties["geo_wkt"])
	match := pointRe.FindStringSubmatch(geoWkt)
	if len(match) != 3 {
		return nil
	}
	wktLon, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		return err
	}
	wktLat, err := strconv.ParseFloat(match[2], 64)
	if err != nil {
		return err
	}

	region := Region{
		ID:          id,
		Deep:        deep,
		PID:         pid,
		Name:        name,
		ExtPath:     extPath,
		GeoWkt:      Point{Lat: wktLat, Lon: wktLon},
		Coordinates: make([]Polygon, 0, 1),
	}

	switch v := feature.Geometry.(type) {
	case *geom.Polygon:
		polygons := conversionPolygon(v)
		region.Coordinates = append(region.Coordinates, polygons...)

	case *geom.MultiPolygon:
		for i := 0; i < v.NumPolygons(); i++ {
			polygons := conversionPolygon(v.Polygon(i))
			region.Coordinates = append(region.Coordinates, polygons...)
		}

	case nil:
		break
	default:
		return errors.New("geom type error")
	}

	g.Regions = append(g.Regions, region)
	g.indexID[region.ID] = len(g.Regions) - 1 // 添加索引
	return nil
}

// 判断点在哪个边界区内
// @return 返回区域索引
func (g *Geo) Contains(point *Point, deep int) (Region, bool) {

	for _, region := range g.Regions {
		if region.Deep != deep {
			continue
		}
		for _, coords := range region.Coordinates {
			if coords.Contains(point) {
				return region, true
			}
		}
	}
	return Region{}, false
}

func (g *Geo) FindByID(id int) (Region, bool) {
	if index, ok := g.indexID[id]; ok {
		return g.Regions[index], true
	}
	return Region{}, false
}

// 闭合多边形
type Polygon struct {
	Points []Point
}

type Point struct {
	Lat float64 // 经度(-90, 90)
	Lon float64 // 纬度(-180, 180)
}

// 判断点是否在平面内
func (p *Polygon) Contains(point *Point) bool {
	intersections := 0
	ps := p.Points

	if Intersections(point.Lon, point.Lat, ps[len(ps)-1].Lon, ps[len(ps)-1].Lat, ps[0].Lon, ps[0].Lat) {
		intersections++
	}
	for i := 0; i < len(ps)-1; i++ {
		if Intersections(point.Lon, point.Lat, ps[i].Lon, ps[i].Lat, ps[i+1].Lon, ps[i+1].Lat) {
			intersections++
		}
	}
	return intersections%2 == 1
}

// Geojson 转 Polygon
func conversionPolygon(p *geom.Polygon) (polygons []Polygon) {
	for i := 0; i < p.NumLinearRings(); i++ {
		linear := p.LinearRing(i)
		polygon := Polygon{
			Points: make([]Point, 0, linear.NumCoords()),
		}

		for _, coord := range linear.Coords() {
			point := Point{Lon: coord[0], Lat: coord[1]}
			polygon.Points = append(polygon.Points, point)
		}
		polygons = append(polygons, polygon)
	}
	return
}

// 射线法求交点
func Intersections(x, y, x1, y1, x2, y2 float64) bool {
	x0 := (x2-x1)*(y-y1)/(y2-y1) + x1 // 计算交点
	return ((x0 > x1 && x0 < x2) || (x0 > x2 && x0 < x1)) && (x0 > x)
}
