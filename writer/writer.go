package writer

import (
	"github.com/dernise/venise/cache"
	"github.com/dernise/venise/database"
)

type OSM struct {
	cache *cache.OSM
	database *database.MongoDatabase
}
