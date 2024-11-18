package leads

import (
	"time"
)

type Leads struct {
	ID                  string  `json:"ID"`
	Title               string  `json:"TITLE"`
	Honorific           string  `json:"HONORIFIC"`
	Name                string  `json:"NAME"`
	SecondName          string  `json:"SECOND_NAME"`
	LastName            string  `json:"LAST_NAME"`
	CompanyTitle        string  `json:"COMPANY_TITLE"`
	CompanyID           string  `json:"COMPANY_ID"`
	ContactID           *string `json:"CONTACT_ID"`
	IsReturnCustomer    string  `json:"IS_RETURN_CUSTOMER"`
	Birthdate           string  `json:"BIRTHDATE"`
	SourceID            string  `json:"SOURCE_ID"`
	SourceDescription   string  `json:"SOURCE_DESCRIPTION"`
	StatusID            string  `json:"STATUS_ID"`
	StatusDescription   *string `json:"STATUS_DESCRIPTION"`
	Post                string  `json:"POST"`
	Comments            string  `json:"COMMENTS"`
	CurrencyID          string  `json:"CURRENCY_ID"`
	Opportunity         string  `json:"OPPORTUNITY"`
	IsManualOpportunity string  `json:"IS_MANUAL_OPPORTUNITY"`
	HasPhone            string  `json:"HAS_PHONE"`
	HasEmail            string  `json:"HAS_EMAIL"`
	HasImol             string  `json:"HAS_IMOL"`
	AssignedByID        string  `json:"ASSIGNED_BY_ID"`
	CreatedByID         string  `json:"CREATED_BY_ID"`
	ModifyByID          string  `json:"MODIFY_BY_ID"`
	DateCreate          string  `json:"DATE_CREATE"`
	DateModify          string  `json:"DATE_MODIFY"`
	DateClosed          string  `json:"DATE_CLOSED"`
	StatusSemanticID    string  `json:"STATUS_SEMANTIC_ID"`
	Opened              string  `json:"OPENED"`
	OriginatorID        *string `json:"ORIGINATOR_ID"`
	OriginID            *string `json:"ORIGIN_ID"`
	MovedByID           string  `json:"MOVED_BY_ID"`
	MovedTime           string  `json:"MOVED_TIME"`
	Address             *string `json:"ADDRESS"`
	Address2            *string `json:"ADDRESS_2"`
	AddressCity         *string `json:"ADDRESS_CITY"`
	AddressPostalCode   *string `json:"ADDRESS_POSTAL_CODE"`
	AddressRegion       *string `json:"ADDRESS_REGION"`
	AddressProvince     *string `json:"ADDRESS_PROVINCE"`
	AddressCountry      *string `json:"ADDRESS_COUNTRY"`
	AddressCountryCode  *string `json:"ADDRESS_COUNTRY_CODE"`
	AddressLocAddrID    *string `json:"ADDRESS_LOC_ADDR_ID"`
	UtmSource           *string `json:"UTM_SOURCE"`
	UtmMedium           *string `json:"UTM_MEDIUM"`
	UtmCampaign         *string `json:"UTM_CAMPAIGN"`
	UtmContent          *string `json:"UTM_CONTENT"`
	UtmTerm             *string `json:"UTM_TERM"`
	LastActivityBy      string  `json:"LAST_ACTIVITY_BY"`
	LastActivityTime    string  `json:"LAST_ACTIVITY_TIME"`
}

type RequestTime struct {
	Start            float64   `json:"start"`
	Finish           float64   `json:"finish"`
	Duration         float64   `json:"duration"`
	Processing       float64   `json:"processing"`
	DateStart        time.Time `json:"date_start"`
	DateFinish       time.Time `json:"date_finish"`
	OperatingResetAt int64     `json:"operating_reset_at"`
	Operating        float64   `json:"operating"`
}

type ApiResponseLeads struct {
	Result []Leads     `json:"result"`
	Total  int         `json:"total"`
	Time   RequestTime `json:"time"`
}
