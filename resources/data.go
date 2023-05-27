package resources

type Data map[string]interface{}

func (d Data) GetId() string {
	if id, ok := d["id"]; ok {
		return id.(string)
	} else {
		return ""
	}
}

func (d Data) GetString(key string) string {
	if str, ok := d[key]; ok {
		return str.(string)
	} else {
		return ""
	}
}

func (d Data) GetInt64(key string) int64 {
	if val, ok := d[key]; ok {
		return val.(int64)
	} else {
		return 0
	}
}

func (d Data) GetData(key string) map[string]interface{} {
	if data, ok := d[key]; ok {
		return data.(map[string]interface{})
	} else {
		return make(map[string]interface{})
	}
}

func (d Data) Get(key string) interface{} {
	if val, ok := d[key]; ok {
		return val
	} else {
		return nil
	}
}
