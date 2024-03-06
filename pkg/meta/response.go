package meta

type BasicResponse struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type HybrisResponse struct {
	HybrisResponse interface{} `json:"hybris"`
}

type StrReturn struct {
	ListString []string `json:"lString"`
}
