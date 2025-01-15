package gis

import (
	"regexp"
	"testing"
	"time"
)

func TestGeoJson(t *testing.T) {
	geo := NewGeos()
	if err := geo.LoadGeoJson("polygon.json"); err != nil {
		t.Fatal(err)
	}
	// t.Logf("geo: %v", geo)

	t1 := time.Now()
	region, ok := geo.Contains(&Point{Lon: 116.71, Lat: 41.19}, 2)
	if !ok {
		t.Fatal("找不到区域")
	}

	t.Logf("%s, geowkt: %v", region.ExtPath, region.GeoWkt)
	t.Logf("time: %v", time.Since(t1))
}

func BenchmarkGeoJson(b *testing.B) {
	geo := NewGeos()
	if err := geo.LoadGeoJson("polygon.json"); err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		geo.Contains(&Point{116.71, 41.19}, 1)
	}
}

func TestIntersections(t *testing.T) {
	ok := Intersections(116.40, 39.91, 115.474294, 39.934587, 115.516785, 39.89779)
	t.Logf("%v", ok)
}

func TestFindString(t *testing.T) {
	complie := regexp.MustCompile(`POINT.+\((.*)\ (.*)\)`)
	res := complie.FindStringSubmatch("POINT (104.034114 34.439662)")
	t.Logf("%v", res)
}

func TestStatis(t *testing.T) {
	g := NewGeos()
	if err := g.LoadGeoJson("polygon.json"); err != nil {
		t.Fatal(err)
	}

	var regions = []int{50702, 450703, 450108, 450105, 450107, 450123, 451082, 451022, 451082, 450123, 450110, 450126, 451302, 450126, 450804, 450181, 450721, 450703, 450702, 450603, 450602, 450603, 450703, 450621, 451421, 451402, 451424, 451425, 451024, 451081, 451024, 451003, 451002, 451029, 451031, 451029, 451002, 450123, 450110, 450126}

	var routers = make([][]string, 0)
	for i := 0; i < len(regions) && i >= 0; i++ {
		var repeated = map[int]struct{}{}

		for j := i; j < len(regions); j++ {

			if _, ok := repeated[regions[j]]; ok {
				var router = make([]string, 0, j-i)
				for k := i; k < j; k++ {
					r, _ := g.FindByID(regions[k])
					router = append(router, r.Name)
				}
				routers = append(routers, router)

				i = j - 2
				break
			}
			repeated[regions[j]] = struct{}{}
		}
	}

	t.Logf("router: %v", routers)
}
