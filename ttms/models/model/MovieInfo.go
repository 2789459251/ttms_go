package model

type MovieInfo struct {
	Info string `json:"info"`
}

func (movie MovieInfo) Index() string {
	return "movie_index"
}
func (movie MovieInfo) Mapping() string {
	return `{
		"mapping":{
			"properties":{
				"info":{
					"type":"text"
				}
			}
		}
	}`
}
