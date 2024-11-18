package companies

type Company struct {
	ID                               string  `json:"ID"`
	Title                            string  `json:"TITLE"`
	HasEmail                         string  `json:"HAS_EMAIL"`
	Emails                           []Email `json:"EMAIL"`
	AssignedByID                     string  `json:"ASSIGNED_BY_ID"`
	INN                              int     `json:"UF_CRM_1726674632"`
	NISHA                            string  `json:"UF_CRM_1726587588"`
	GEOGRAPHIA                       string  `json:"UF_CRM_1726587646"`
	KONECHNIY_PRODUCT                string  `json:"UF_CRM_1726587741"`
	OBOROT_KOMPANII                  int     `json:"UF_CRM_1726587790"`
	SREDNEMESYACHNAYA_VIRYCHKA       int     `json:"UF_CRM_1726587900"`
	KOLICHESTVO_SOTRUDNIKOV          int     `json:"UF_CRM_EMPLOYEES"`
	KOLICHESTVO_SOTRUDNIKOV_OP       int     `json:"UF_CRM_1726587932"`
	EST_ROP                          string  `json:"UF_CRM_1726588037"`
	EST_HR                           string  `json:"UF_CRM_1726588054"`
	OSNOVNIE_KANALI_PRODAJ           string  `json:"UF_CRM_1726588155"`
	SEZONNOST                        string  `json:"UF_CRM_1726588171"`
	ZAPROS                           string  `json:"UF_CRM_1726588242"`
	KAKIE_CELI_PERED_KOMPANIEI       string  `json:"UF_CRM_1726588284"`
	CHTO_OJIDAETE_OT_SOTRUDNICHESTVA string  `json:"UF_CRM_1726588358"`
}

type Email struct {
	ID    string `json:"ID"`
	Value string `json:"VALUE"`
}
