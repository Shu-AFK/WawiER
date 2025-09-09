package wawi_reqs

type Order struct {
	TotalItems         int          `json:"TotalItems"`
	PageNumber         int          `json:"PageNumber"`
	PageSize           int          `json:"PageSize"`
	Items              []SalesOrder `json:"Items"`
	TotalPages         int          `json:"TotalPages"`
	HasPreviousPage    bool         `json:"HasPreviousPage"`
	HasNextPage        bool         `json:"HasNextPage"`
	NextPageNumber     int          `json:"NextPageNumber"`
	PreviousPageNumber int          `json:"PreviousPageNumber"`
}

type SalesOrder struct {
	Id                       int             `json:"Id"`
	Number                   string          `json:"Number"`
	ExternalNumber           string          `json:"ExternalNumber"`
	BillingNumber            string          `json:"BillingNumber"`
	CompanyId                int             `json:"CompanyId"`
	DepartureCountry         Country         `json:"DepartureCountry"`
	CustomerId               int             `json:"CustomerId"`
	CustomerVatID            string          `json:"CustomerVatID"`
	MerchantVatID            string          `json:"MerchantVatID"`
	BillingAddress           Address         `json:"BillingAddress"`
	Shipmentaddress          Address         `json:"Shipmentaddress"`
	SalesOrderDate           string          `json:"SalesOrderDate"`
	SalesOrderPaymentDetails PaymentDetails  `json:"SalesOrderPaymentDetails"`
	SalesOrderShippingDetail ShippingDetails `json:"SalesOrderShippingDetail"`
	ColorcodeId              int             `json:"ColorcodeId"`
	IsExternalInvoice        bool            `json:"IsExternalInvoice"`
	Comment                  string          `json:"Comment"`
	CustomerComment          string          `json:"CustomerComment"`
	IsCancelled              bool            `json:"IsCancelled"`
	LanguageIso              string          `json:"LanguageIso"`
	CancellationDetails      Cancellation    `json:"CancellationDetails"`
	SalesChannelId           string          `json:"SalesChannelId"`
	UserCreatedId            int             `json:"UserCreatedId"`
	UserId                   int             `json:"UserId"`
}

type Country struct {
	CountryISO     string  `json:"CountryISO"`
	State          string  `json:"State"`
	CurrencyIso    string  `json:"CurrencyIso"`
	CurrencyFactor float64 `json:"CurrencyFactor"`
}

type Address struct {
	Id                int    `json:"Id"`
	Company           string `json:"Company"`
	Company2          string `json:"Company2"`
	FormOfAddress     string `json:"FormOfAddress"`
	Title             string `json:"Title"`
	FirstName         string `json:"FirstName"`
	LastName          string `json:"LastName"`
	Street            string `json:"Street"`
	Address2          string `json:"Address2"`
	PostalCode        string `json:"PostalCode"`
	City              string `json:"City"`
	State             string `json:"State"`
	CountryIso        string `json:"CountryIso"`
	VatID             string `json:"VatID"`
	PhoneNumber       string `json:"PhoneNumber"`
	MobilePhoneNumber string `json:"MobilePhoneNumber"`
	EmailAddress      string `json:"EmailAddress"`
	Fax               string `json:"Fax"`
}

type PaymentDetails struct {
	PaymentMethodId  int         `json:"PaymentMethodId"`
	PaymentStatus    interface{} `json:"PaymentStatus"`
	TotalGrossAmount float64     `json:"TotalGrossAmount"`
	CurrencyIso      string      `json:"CurrencyIso"`
	CurrencyFactor   float64     `json:"CurrencyFactor"`
	DateOfPayment    string      `json:"DateOfPayment"`
	StillToPay       float64     `json:"StillToPay"`
	PaymentTarget    float64     `json:"PaymentTarget"`
	CashDiscount     float64     `json:"CashDiscount"`
	CashDiscountDays int         `json:"CashDiscountDays"`
}

type ShippingDetails struct {
	ShippingMethodId       int     `json:"ShippingMethodId"`
	DeliveryCompleteStatus int     `json:"DeliveryCompleteStatus"`
	ShippingPriority       int     `json:"ShippingPriority"`
	ShippingDate           string  `json:"ShippingDate"`
	EstimatedDeliveryDate  string  `json:"EstimatedDeliveryDate"`
	DeliveredDate          string  `json:"DeliveredDate"`
	OnHoldReasonId         int     `json:"OnHoldReasonId"`
	ExtraWeight            float64 `json:"ExtraWeight"`
}

type Cancellation struct {
	CancellationReasonId int    `json:"CancellationReasonId"`
	CancellationComment  string `json:"CancellationComment"`
	Date                 string `json:"Date"`
}

type OrderItem struct {
	Id                int     `json:"Id"`
	SalesOrderId      int     `json:"SalesOrderId"`
	ItemId            int     `json:"ItemId"`
	Name              string  `json:"Name"`
	SKU               string  `json:"SKU"`
	Type              int     `json:"Type"`
	Quantity          float64 `json:"Quantity"`
	QuantityDelivered float64 `json:"QuantityDelivered"`
	QuantityReturned  float64 `json:"QuantityReturned"`
	SalesUnit         string  `json:"SalesUnit"`
	SalesPriceNet     float64 `json:"SalesPriceNet"`
	SalesPriceGross   float64 `json:"SalesPriceGross"`
	Discount          float64 `json:"Discount"`
	PurchasePriceNet  float64 `json:"PurchasePriceNet"`
	TaxRate           float64 `json:"TaxRate"`
	Notice            string  `json:"Notice"`
}

type StockResponse struct {
	TotalItems         int         `json:"TotalItems"`
	PageNumber         int         `json:"PageNumber"`
	PageSize           int         `json:"PageSize"`
	Items              []StockItem `json:"Items"`
	TotalPages         int         `json:"TotalPages"`
	HasPreviousPage    bool        `json:"HasPreviousPage"`
	HasNextPage        bool        `json:"HasNextPage"`
	NextPageNumber     int         `json:"NextPageNumber"`
	PreviousPageNumber int         `json:"PreviousPageNumber"`
}

type StockItem struct {
	WarehouseId                   int     `json:"WarehouseId"`
	StorageLocationId             int     `json:"StorageLocationId"`
	StorageLocationName           string  `json:"StorageLocationName"`
	ItemId                        int     `json:"ItemId"`
	ShelfLifeExpirationDate       string  `json:"ShelfLifeExpirationDate"`
	BatchNumber                   string  `json:"BatchNumber"`
	QuantityTotal                 float64 `json:"QuantityTotal"`
	QuantityLockedForShipment     float64 `json:"QuantityLockedForShipment"`
	QuantityLockedForAvailability float64 `json:"QuantityLockedForAvailability"`
	QuantityInPickingLists        float64 `json:"QuantityInPickingLists"`
}
