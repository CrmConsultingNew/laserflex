package models

import (
	"encoding/json"
	"strconv"
	"time"
)

type Float64Str float64

func (fs *Float64Str) UnmarshalJSON(data []byte) error {
	// First try to unmarshal to a float64
	var f float64
	if err := json.Unmarshal(data, &f); err == nil {
		*fs = Float64Str(f)
		return nil
	}

	// If that fails, try to unmarshal to a string and then convert
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	// Convert string to float64
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return err
	}

	*fs = Float64Str(f)
	return nil
}

type Deal struct {
	ID                  string      `json:"ID"`
	Title               string      `json:"TITLE"`
	TypeID              string      `json:"TYPE_ID"`
	StageID             string      `json:"STAGE_ID"`
	Probability         *Float64Str `json:"PROBABILITY"` // Nullable
	CurrencyID          string      `json:"CURRENCY_ID"`
	Opportunity         string      `json:"OPPORTUNITY"`
	IsManualOpportunity string      `json:"IS_MANUAL_OPPORTUNITY"`
	TaxValue            string      `json:"TAX_VALUE"`
	LeadID              *string     `json:"LEAD_ID"` // Nullable
	CompanyID           string      `json:"COMPANY_ID"`
	ContactID           string      `json:"CONTACT_ID"` // Nullable
	QuoteID             *string     `json:"QUOTE_ID"`   // Nullable
	BeginDate           time.Time   `json:"BEGINDATE"`
	CloseDate           time.Time   `json:"CLOSEDATE"`
	AssignedByID        string      `json:"ASSIGNED_BY_ID"`
	CreatedByID         string      `json:"CREATED_BY_ID"`
	ModifyByID          string      `json:"MODIFY_BY_ID"`
	DateCreate          string      `json:"DATE_CREATE"`
	DateModify          time.Time   `json:"DATE_MODIFY"`
	Opened              string      `json:"OPENED"`
	Closed              string      `json:"CLOSED"`
	Comments            string      `json:"COMMENTS"`
	AdditionalInfo      *string     `json:"ADDITIONAL_INFO"` // Nullable
	LocationID          *string     `json:"LOCATION_ID"`     // Nullable
	CategoryID          string      `json:"CATEGORY_ID"`
	StageSemanticID     string      `json:"STAGE_SEMANTIC_ID"`
	IsNew               string      `json:"IS_NEW"`
	IsRecurring         string      `json:"IS_RECURRING"`
	IsReturnCustomer    string      `json:"IS_RETURN_CUSTOMER"`
	IsRepeatedApproach  string      `json:"IS_REPEATED_APPROACH"`
	SourceID            string      `json:"SOURCE_ID"`
	SourceDescription   *string     `json:"SOURCE_DESCRIPTION"` // Nullable
	OriginatorID        *string     `json:"ORIGINATOR_ID"`      // Nullable
	OriginID            *string     `json:"ORIGIN_ID"`          // Nullable
	MovedByID           string      `json:"MOVED_BY_ID"`
	MovedTime           time.Time   `json:"MOVED_TIME"`
	LastActivityTime    time.Time   `json:"LAST_ACTIVITY_TIME"`
	UtmSource           *string     `json:"UTM_SOURCE"`   // Nullable
	UtmMedium           *string     `json:"UTM_MEDIUM"`   // Nullable
	UtmCampaign         *string     `json:"UTM_CAMPAIGN"` // Nullable
	UtmContent          *string     `json:"UTM_CONTENT"`  // Nullable
	UtmTerm             *string     `json:"UTM_TERM"`     // Nullable
	LastActivityBy      string      `json:"LAST_ACTIVITY_BY"`
	LinkOnParentDealNps string      `json:"UF_CRM_1712927864"`
	Branch              string      `json:"UF_CRM_1690982742603"`
	RatingNps           string      `json:"UF_CRM_1712927909"`
	VisitDate           time.Time   `json:"UF_CRM_1690209734961"`
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

type ApiResponseHarizma struct {
	Result Deal        `json:"result"`
	Total  int         `json:"total"`
	Time   RequestTime `json:"time"`
}
