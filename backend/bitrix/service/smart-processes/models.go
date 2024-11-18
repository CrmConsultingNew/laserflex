package smart_processes

type CRMItem struct {
	ID                  int           `json:"id"`
	XmlID               string        `json:"xmlId"`
	Title               string        `json:"title"`
	CreatedBy           int           `json:"createdBy"`
	UpdatedBy           int           `json:"updatedBy"`
	MovedBy             int           `json:"movedBy"`
	CreatedTime         string        `json:"createdTime"`
	UpdatedTime         string        `json:"updatedTime"`
	MovedTime           string        `json:"movedTime"`
	CategoryID          int           `json:"categoryId"`
	Opened              string        `json:"opened"`
	StageID             string        `json:"stageId"`
	PreviousStageID     string        `json:"previousStageId"`
	BeginDate           string        `json:"begindate"`
	CloseDate           string        `json:"closedate"`
	CompanyID           int           `json:"companyId"`
	ContactID           int           `json:"contactId"`
	Opportunity         float64       `json:"opportunity"`
	IsManualOpportunity string        `json:"isManualOpportunity"`
	TaxValue            int           `json:"taxValue"`
	CurrencyID          string        `json:"currencyId"`
	MyCompanyID         int           `json:"mycompanyId"`
	SourceID            string        `json:"sourceId"`
	SourceDescription   string        `json:"sourceDescription"`
	WebformID           int           `json:"webformId"`
	UfCrm26_1712127967  []interface{} `json:"ufCrm26_1712127967"`
	UfCrm26_1712128088  string        `json:"ufCrm26_1712128088"`
	UfCrm26_1713936681  struct {
		ID         int    `json:"id"`
		URL        string `json:"url"`
		URLMachine string `json:"urlMachine"`
	} `json:"ufCrm26_1713936681"`
	UfCrm26_1717229767 interface{} `json:"ufCrm26_1717229767"`
	AssignedByID       int         `json:"assignedById"`
	LastActivityBy     int         `json:"lastActivityBy"`
	LastActivityTime   string      `json:"lastActivityTime"`
	ParentID2          int         `json:"parentId2"`
	UtmSource          interface{} `json:"utmSource"`
	UtmMedium          interface{} `json:"utmMedium"`
	UtmCampaign        interface{} `json:"utmCampaign"`
	UtmContent         interface{} `json:"utmContent"`
	UtmTerm            interface{} `json:"utmTerm"`
	Observers          []int       `json:"observers"`
	ContactIDs         []int       `json:"contactIds"`
	EntityTypeID       int         `json:"entityTypeId"`
}

type CRMResponse struct {
	Result struct {
		Items []CRMItem `json:"items"`
	} `json:"result"`
}
