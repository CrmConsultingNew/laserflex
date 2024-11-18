package description

// Deal represents the structure of a deal in the response
type Deal struct {
	ID                  string  `json:"ID"`
	Title               string  `json:"TITLE"`
	TypeID              string  `json:"TYPE_ID"`
	StageID             string  `json:"STAGE_ID"`
	Probability         string  `json:"PROBABILITY"` // Use pointer to represent null
	CurrencyID          string  `json:"CURRENCY_ID"`
	Opportunity         string  `json:"OPPORTUNITY"`
	IsManualOpportunity string  `json:"IS_MANUAL_OPPORTUNITY"`
	TaxValue            string  `json:"TAX_VALUE"`
	LeadID              *string `json:"LEAD_ID"` // Use pointer to represent null
	CompanyID           string  `json:"COMPANY_ID"`
	ContactID           *string `json:"CONTACT_ID"` // Use pointer to represent null
	QuoteID             *string `json:"QUOTE_ID"`   // Use pointer to represent null
	BeginDate           string  `json:"BEGINDATE"`
	CloseDate           string  `json:"CLOSEDATE"`
	AssignedByID        string  `json:"ASSIGNED_BY_ID"`
	CreatedByID         string  `json:"CREATED_BY_ID"`
	ModifyByID          string  `json:"MODIFY_BY_ID"`
	DateCreate          string  `json:"DATE_CREATE"`
	DateModify          string  `json:"DATE_MODIFY"`
	Opened              string  `json:"OPENED"`
	Closed              string  `json:"CLOSED"`
	Comments            string  `json:"COMMENTS"`
	AdditionalInfo      *string `json:"ADDITIONAL_INFO"` // Use pointer to represent null
	LocationID          *string `json:"LOCATION_ID"`     // Use pointer to represent null
	CategoryID          string  `json:"CATEGORY_ID"`
	StageSemanticID     string  `json:"STAGE_SEMANTIC_ID"`
	IsNew               string  `json:"IS_NEW"`
	IsRecurring         string  `json:"IS_RECURRING"`
	IsReturnCustomer    string  `json:"IS_RETURN_CUSTOMER"`
	IsRepeatedApproach  string  `json:"IS_REPEATED_APPROACH"`
	SourceID            string  `json:"SOURCE_ID"`
	SourceDescription   *string `json:"SOURCE_DESCRIPTION"` // Use pointer to represent null
	OriginatorID        *string `json:"ORIGINATOR_ID"`      // Use pointer to represent null
	OriginID            *string `json:"ORIGIN_ID"`          // Use pointer to represent null
	MovedByID           string  `json:"MOVED_BY_ID"`
	MovedTime           string  `json:"MOVED_TIME"`
	LastActivityTime    string  `json:"LAST_ACTIVITY_TIME"`
	UTMSource           *string `json:"UTM_SOURCE"`   // Use pointer to represent null
	UTMMedium           *string `json:"UTM_MEDIUM"`   // Use pointer to represent null
	UTMCampaign         *string `json:"UTM_CAMPAIGN"` // Use pointer to represent null
	UTMContent          *string `json:"UTM_CONTENT"`  // Use pointer to represent null
	UTMTerm             *string `json:"UTM_TERM"`     // Use pointer to represent null
	LastActivityBy      string  `json:"LAST_ACTIVITY_BY"`
	// Add custom fields as needed
	UFCRM1710943349299 string `json:"UF_CRM_1710943349299"`
	UFCRM1710943898    string `json:"UF_CRM_1710943898"`
	UFCRM1710943969    string `json:"UF_CRM_1710943969"`
	UFCRM1711748170    string `json:"UF_CRM_1711748170"`
}

// ResponseTime represents the timing information in the response
type ResponseTime struct {
	Start            float64 `json:"start"`
	Finish           float64 `json:"finish"`
	Duration         float64 `json:"duration"`
	Processing       float64 `json:"processing"`
	DateStart        string  `json:"date_start"`
	DateFinish       string  `json:"date_finish"`
	OperatingResetAt int64   `json:"operating_reset_at"`
	Operating        float64 `json:"operating"` // This field should be a float64
}

// DealResponse represents the top-level structure of the JSON response
type DealResponse struct {
	Result Deal         `json:"result"`
	Time   ResponseTime `json:"time"`
}
