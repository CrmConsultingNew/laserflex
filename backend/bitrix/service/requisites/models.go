package requisites

type Requisites struct {
	ID                 string      `json:"ID"`
	ENTITYTYPEID       string      `json:"ENTITY_TYPE_ID"`
	ENTITYID           string      `json:"ENTITY_ID"`
	PRESETID           string      `json:"PRESET_ID"`
	DATECREATE         string      `json:"DATE_CREATE"`
	DATEMODIFY         string      `json:"DATE_MODIFY"`
	CREATEDBYID        string      `json:"CREATED_BY_ID"`
	MODIFYBYID         string      `json:"MODIFY_BY_ID"`
	NAME               string      `json:"NAME"`
	CODE               interface{} `json:"CODE"`
	XMLID              interface{} `json:"XML_ID"`
	ORIGINATORID       interface{} `json:"ORIGINATOR_ID"`
	ACTIVE             string      `json:"ACTIVE"`
	ADDRESSONLY        string      `json:"ADDRESS_ONLY"`
	SORT               string      `json:"SORT"`
	RQNAME             string      `json:"RQ_NAME"`
	RQFIRSTNAME        string      `json:"RQ_FIRST_NAME"`
	RQLASTNAME         string      `json:"RQ_LAST_NAME"`
	RQSECONDNAME       string      `json:"RQ_SECOND_NAME"`
	RQCOMPANYID        string      `json:"RQ_COMPANY_ID"`
	RQCOMPANYNAME      string      `json:"RQ_COMPANY_NAME"`
	RQCOMPANYFULLNAME  string      `json:"RQ_COMPANY_FULL_NAME"`
	RQCOMPANYREGDATE   string      `json:"RQ_COMPANY_REG_DATE"`
	RQDIRECTOR         string      `json:"RQ_DIRECTOR"`
	RQACCOUNTANT       string      `json:"RQ_ACCOUNTANT"`
	RQCEONAME          string      `json:"RQ_CEO_NAME"`
	RQCEOWORKPOS       string      `json:"RQ_CEO_WORK_POS"`
	RQCONTACT          string      `json:"RQ_CONTACT"`
	RQEMAIL            string      `json:"RQ_EMAIL"`
	RQPHONE            string      `json:"RQ_PHONE"`
	RQFAX              string      `json:"RQ_FAX"`
	RQIDENTTYPE        string      `json:"RQ_IDENT_TYPE"`
	RQIDENTDOC         string      `json:"RQ_IDENT_DOC"`
	RQIDENTDOCSER      string      `json:"RQ_IDENT_DOC_SER"`
	RQIDENTDOCNUM      string      `json:"RQ_IDENT_DOC_NUM"`
	RQIDENTDOCPERSNUM  string      `json:"RQ_IDENT_DOC_PERS_NUM"`
	RQIDENTDOCDATE     string      `json:"RQ_IDENT_DOC_DATE"`
	RQIDENTDOCISSUEDBY string      `json:"RQ_IDENT_DOC_ISSUED_BY"`
	RQIDENTDOCDEPCODE  string      `json:"RQ_IDENT_DOC_DEP_CODE"`
	RQINN              string      `json:"RQ_INN"`
	RQKPP              string      `json:"RQ_KPP"`
	RQUSRLE            string      `json:"RQ_USRLE"`
	RQIFNS             string      `json:"RQ_IFNS"`
	RQOGRN             string      `json:"RQ_OGRN"`
	RQOGRNIP           string      `json:"RQ_OGRNIP"`
	RQOKPO             string      `json:"RQ_OKPO"`
	RQOKTMO            string      `json:"RQ_OKTMO"`
	RQOKVED            string      `json:"RQ_OKVED"`
	RQEDRPOU           string      `json:"RQ_EDRPOU"`
	RQDRFO             string      `json:"RQ_DRFO"`
	RQKBE              string      `json:"RQ_KBE"`
	RQIIN              string      `json:"RQ_IIN"`
	RQBIN              string      `json:"RQ_BIN"`
	RQSTCERTSER        string      `json:"RQ_ST_CERT_SER"`
	RQSTCERTNUM        string      `json:"RQ_ST_CERT_NUM"`
	RQSTCERTDATE       string      `json:"RQ_ST_CERT_DATE"`
	RQVATPAYER         string      `json:"RQ_VAT_PAYER"`
	RQVATID            string      `json:"RQ_VAT_ID"`
	RQVATCERTSER       string      `json:"RQ_VAT_CERT_SER"`
	RQVATCERTNUM       string      `json:"RQ_VAT_CERT_NUM"`
	RQVATCERTDATE      string      `json:"RQ_VAT_CERT_DATE"`
	RQRESIDENCECOUNTRY string      `json:"RQ_RESIDENCE_COUNTRY"`
	RQBASEDOC          string      `json:"RQ_BASE_DOC"`
	RQREGON            string      `json:"RQ_REGON"`
	RQKRS              string      `json:"RQ_KRS"`
	RQPESEL            string      `json:"RQ_PESEL"`
	RQLEGALFORM        string      `json:"RQ_LEGAL_FORM"`
	RQSIRET            string      `json:"RQ_SIRET"`
	RQSIREN            string      `json:"RQ_SIREN"`
	RQCAPITAL          string      `json:"RQ_CAPITAL"`
	RQRCS              string      `json:"RQ_RCS"`
	RQCNPJ             string      `json:"RQ_CNPJ"`
	RQSTATEREG         string      `json:"RQ_STATE_REG"`
	RQMNPLREG          string      `json:"RQ_MNPL_REG"`
	RQCPF              string      `json:"RQ_CPF"`
}
