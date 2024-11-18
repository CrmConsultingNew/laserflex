package deals

type ApiResponse struct {
	Result DealInfo `json:"result"`
}

type DealInfo struct {
	ID                               string `json:"ID"`
	Title                            string `json:"TITLE"`
	Opportunity                      string `json:"OPPORTUNITY"`
	CompanyID                        int    `json:"COMPANY_ID"`
	NISHA                            string `json:"UF_CRM_1726587543"`
	GEOGRAPHIA                       string `json:"UF_CRM_1726587678"`
	KONECHNIY_PRODUCT                string `json:"UF_CRM_1726587712"`
	KOLICHESTVO_SOTRUDNIKOV          int    `json:"UF_CRM_1706503695"`
	KOLICHESTVO_SOTRUDNIKOV_OP       int    `json:"UF_CRM_1726587982"`
	EST_ROP                          string `json:"UF_CRM_1726588010"`
	EST_HR                           string `json:"UF_CRM_1726588094"`
	OSNOVNIE_KANALI_PRODAJ           string `json:"UF_CRM_1726588113"`
	SEZONNOST                        string `json:"UF_CRM_1726588197"`
	ZAPROS                           string `json:"UF_CRM_1726588214"`
	KAKIE_CELI_PERED_KOMPANIEI       string `json:"UF_CRM_1726588301"`
	CHTO_OJIDAETE_OT_SOTRUDNICHESTVA string `json:"UF_CRM_1726588336"`
}
