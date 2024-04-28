package model

type MovieInfo string

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
