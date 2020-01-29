package json

import (
	"encoding/json"
	"github.com/deltrinos/tpl21/log"
)

func PrettyPrintJson(o interface{}) string {
	prettyJSON, err := json.MarshalIndent(o, "", "    ")
	if err != nil {
		log.Error().Err(err).Msgf("failed to json.MarshalIndent: %v", err)
		return ""
	}
	return string(prettyJSON)
}
