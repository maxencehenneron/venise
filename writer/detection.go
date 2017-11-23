package writer

import "github.com/dernise/venise/structures"

type PolygonDetection struct {
	Values  []string `json:"values"`
	Polygon string   `json:"polygon"`
	Key     string   `json:"key"`
}

type Detection []PolygonDetection

func (detections Detection) IsPolygon(way *structures.Way) bool {
	for _, detection := range detections {
		switch detection.Polygon {
		case "all":
			if _, ok := way.Tags[detection.Key]; ok {
				return true
			}
			break
		case "whitelist":
			if tagValue, ok := way.Tags[detection.Key]; ok {
				for _, value := range detection.Values {
					if tagValue == value {
						return true
					}
				}
			}
			break
		case "blacklist":
			if tagValue, ok := way.Tags[detection.Key]; ok {
				for _, value := range detection.Values {
					if tagValue == value {
						return false
					}
				}
				return true
			}
			break
		}
	}
	return false
}
