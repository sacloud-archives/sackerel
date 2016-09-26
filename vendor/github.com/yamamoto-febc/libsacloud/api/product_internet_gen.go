package api

/************************************************
  generated by IDE. for [ProductInternetAPI]
************************************************/

import (
	"github.com/yamamoto-febc/libsacloud/sacloud"
)

/************************************************
   To support influent interface for Find()
************************************************/

func (api *ProductInternetAPI) Reset() *ProductInternetAPI {
	api.reset()
	return api
}

func (api *ProductInternetAPI) Offset(offset int) *ProductInternetAPI {
	api.offset(offset)
	return api
}

func (api *ProductInternetAPI) Limit(limit int) *ProductInternetAPI {
	api.limit(limit)
	return api
}

func (api *ProductInternetAPI) Include(key string) *ProductInternetAPI {
	api.include(key)
	return api
}

func (api *ProductInternetAPI) Exclude(key string) *ProductInternetAPI {
	api.exclude(key)
	return api
}

func (api *ProductInternetAPI) FilterBy(key string, value interface{}) *ProductInternetAPI {
	api.filterBy(key, value, false)
	return api
}

// func (api *ProductInternetAPI) FilterMultiBy(key string, value interface{}) *ProductInternetAPI {
// 	api.filterBy(key, value, true)
// 	return api
// }

func (api *ProductInternetAPI) WithNameLike(name string) *ProductInternetAPI {
	return api.FilterBy("Name", name)
}

func (api *ProductInternetAPI) WithTag(tag string) *ProductInternetAPI {
	return api.FilterBy("Tags.Name", tag)
}
func (api *ProductInternetAPI) WithTags(tags []string) *ProductInternetAPI {
	return api.FilterBy("Tags.Name", []interface{}{tags})
}

// func (api *ProductInternetAPI) WithSizeGib(size int) *ProductInternetAPI {
// 	api.FilterBy("SizeMB", size*1024)
// 	return api
// }

// func (api *ProductInternetAPI) WithSharedScope() *ProductInternetAPI {
// 	api.FilterBy("Scope", "shared")
// 	return api
// }

// func (api *ProductInternetAPI) WithUserScope() *ProductInternetAPI {
// 	api.FilterBy("Scope", "user")
// 	return api
// }

func (api *ProductInternetAPI) SortBy(key string, reverse bool) *ProductInternetAPI {
	api.sortBy(key, reverse)
	return api
}

func (api *ProductInternetAPI) SortByName(reverse bool) *ProductInternetAPI {
	api.sortByName(reverse)
	return api
}

// func (api *ProductInternetAPI) SortBySize(reverse bool) *ProductInternetAPI {
// 	api.sortBy("SizeMB", reverse)
// 	return api
// }

/************************************************
  To support CRUD(Create/Read/Update/Delete)
************************************************/

//func (api *ProductInternetAPI) New() *sacloud.ProductInternet {
// 	return &sacloud.ProductInternet{}
//}

// func (api *ProductInternetAPI) Create(value *sacloud.ProductInternet) (*sacloud.ProductInternet, error) {
// 	return api.request(func(res *sacloud.Response) error {
// 		return api.create(api.createRequest(value), res)
// 	})
// }

func (api *ProductInternetAPI) Read(id int64) (*sacloud.ProductInternet, error) {
	return api.request(func(res *sacloud.Response) error {
		return api.read(id, nil, res)
	})
}

// func (api *ProductInternetAPI) Update(id int64, value *sacloud.ProductInternet) (*sacloud.ProductInternet, error) {
// 	return api.request(func(res *sacloud.Response) error {
// 		return api.update(id, api.createRequest(value), res)
// 	})
// }

// func (api *ProductInternetAPI) Delete(id int64) (*sacloud.ProductInternet, error) {
// 	return api.request(func(res *sacloud.Response) error {
// 		return api.delete(id, nil, res)
// 	})
// }

/************************************************
  Inner functions
************************************************/

func (api *ProductInternetAPI) setStateValue(setFunc func(*sacloud.Request)) *ProductInternetAPI {
	api.baseAPI.setStateValue(setFunc)
	return api
}

func (api *ProductInternetAPI) request(f func(*sacloud.Response) error) (*sacloud.ProductInternet, error) {
	res := &sacloud.Response{}
	err := f(res)
	if err != nil {
		return nil, err
	}
	return res.InternetPlan, nil
}

func (api *ProductInternetAPI) createRequest(value *sacloud.ProductInternet) *sacloud.Request {
	req := &sacloud.Request{}
	req.InternetPlan = value
	return req
}
