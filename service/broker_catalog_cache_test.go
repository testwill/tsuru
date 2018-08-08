// Copyright 2018 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package service

import (
	"encoding/json"
	"fmt"

	"github.com/tsuru/tsuru/types/cache"
	"github.com/tsuru/tsuru/types/service"
	"gopkg.in/check.v1"
)

func (s *S) TestCacheSave(c *check.C) {
	catalog := service.BrokerCatalog{
		Services: []service.BrokerService{{
			ID:          "123",
			Name:        "service1",
			Description: "my service",
			Plans: []service.BrokerPlan{{
				ID:          "456",
				Name:        "my-plan",
				Description: "plan description",
			}},
		}},
	}
	service := &serviceBrokerCatalogCacheService{
		storage: &cache.MockCacheStorage{
			OnPut: func(entry cache.CacheEntry) error {
				c.Assert(entry.Key, check.Equals, "my-catalog")
				var cat service.BrokerCatalog
				err := json.Unmarshal([]byte(entry.Value), &cat)
				c.Assert(err, check.IsNil)
				c.Assert(cat, check.DeepEquals, catalog)
				return nil
			},
		},
	}
	err := service.Save("my-catalog", catalog)
	c.Assert(err, check.IsNil)
}

func (s *S) TestCacheLoad(c *check.C) {
	catalog := service.BrokerCatalog{
		Services: []service.BrokerService{{
			ID:          "123",
			Name:        "service1",
			Description: "my service",
			Plans: []service.BrokerPlan{{
				ID:          "456",
				Name:        "my-plan",
				Description: "plan description",
			}},
		}},
	}
	service := &serviceBrokerCatalogCacheService{
		storage: &cache.MockCacheStorage{
			OnGet: func(key string) (cache.CacheEntry, error) {
				c.Assert(key, check.Equals, "my-catalog")
				b, err := json.Marshal(catalog)
				c.Assert(err, check.IsNil)
				return cache.CacheEntry{Key: key, Value: string(b)}, nil
			},
		},
	}
	cat, err := service.Load("my-catalog")
	c.Assert(err, check.IsNil)
	c.Assert(cat, check.NotNil)
	c.Assert(*cat, check.DeepEquals, catalog)
}

func (s *S) TestCacheLoadNotFound(c *check.C) {
	service := &serviceBrokerCatalogCacheService{
		storage: &cache.MockCacheStorage{
			OnGet: func(key string) (cache.CacheEntry, error) {
				c.Assert(key, check.Equals, "unknown-catalog")
				return cache.CacheEntry{}, fmt.Errorf("not found")
			},
		},
	}
	cat, err := service.Load("unknown-catalog")
	c.Assert(err, check.NotNil)
	c.Assert(cat, check.IsNil)
}
