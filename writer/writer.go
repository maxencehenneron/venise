package writer

import (
	"github.com/dernise/venise/cache"
	"github.com/dernise/venise/database"
)

type OSM struct {
	cache           *cache.OSM
	database        *database.MongoDatabase
	interestingTags map[string][]string
}

func (o *OSM) HasInterrestingTag(tags map[string]string) bool {
	for key, values := range o.interestingTags {
		if val, ok := tags[key]; ok {
			for _, value := range values {
				if val == value {
					return true
				}
			}
		}
	}
	return false
}
