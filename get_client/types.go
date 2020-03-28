package main

import "time"

type portalClient struct {
	ID                        int         `json:"id"`
	Name                      string      `json:"name"`
	FirstName                 string      `json:"first_name"`
	LastName                  string      `json:"last_name"`
	Company                   interface{} `json:"company"`
	City                      string      `json:"city"`
	Country                   string      `json:"country"`
	State                     string      `json:"state"`
	DateCreated               time.Time   `json:"date_created"`
	Currency                  string      `json:"currency"`
	Phone                     string      `json:"phone"`
	CountryName               string      `json:"country_name"`
	LongName                  string      `json:"long_name"`
	Fax                       interface{} `json:"fax"`
	Users                     []int       `json:"users"`
	Address1                  string      `json:"address1"`
	Address2                  interface{} `json:"address2"`
	Email                     string      `json:"email"`
	ZipCode                   string      `json:"zip_code"`
	VatID                     interface{} `json:"vat_id"`
	SuspendInsteadOfTerminate bool        `json:"suspend_instead_of_terminate"`
	Credits                   []struct {
		Client   int    `json:"client"`
		Currency string `json:"currency"`
		Amount   string `json:"amount"`
	} `json:"credits"`
	CustomFields          []interface{} `json:"custom_fields"`
	UptodateCredit        string        `json:"uptodate_credit"`
	HasOpenstackServices  bool          `json:"has_openstack_services"`
	Status                string        `json:"status"`
	TaxExempt             bool          `json:"tax_exempt"`
	OutofcreditDatetime   interface{}   `json:"outofcredit_datetime"`
	OpenstackBillingPlans []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"openstack_billing_plans"`
	GroupName             string      `json:"group_name"`
	ResellerClient        interface{} `json:"reseller_client"`
	BelongsToReseller     bool        `json:"belongs_to_reseller"`
	ResellerClientDetails interface{} `json:"reseller_client_details"`
}

type getUserResp struct {
	PortalID    string `json:"portal_id"`
	PortalName  string `json:"portal_name"`
	PortalEmail string `json:"portal_email"`
}
