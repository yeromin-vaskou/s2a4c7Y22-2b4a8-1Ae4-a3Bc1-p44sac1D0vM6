var res string = GetJson(&general) // general is a struct
	jq := gojsonq.New().FromString(res) 

	cool_json := jq.From("Countries").Where("CountryCode", "=", "AF").Get()