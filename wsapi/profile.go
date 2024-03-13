package wsapi

import (
	"patrickkdev/Go-IQOption-API/tjson"
	"time"
)

type CoreProfile struct {
	ClientSessionID string `json:"client_session_id"`
	Msg             msg    `json:"msg"`              
	Name            string `json:"name"`             
	RequestID       string `json:"request_id"`       
	SessionID       string `json:"session_id"`       
	Status          int64  `json:"status"`           
}

type msg struct {
	IsSuccessful bool          `json:"isSuccessful"`
	Message      []interface{} `json:"message"`     
	Result       result        `json:"result"`      
}

type result struct {
	AccountStatus         string             `json:"account_status"`         
	Address               string             `json:"address"`                
	AuthTwoFactor         interface{}        `json:"auth_two_factor"`        
	Avatar                string             `json:"avatar"`                 
	Balance               int64              `json:"balance"`                
	BalanceID             int64              `json:"balance_id"`             
	BalanceType           int64              `json:"balance_type"`           
	Balances              []interface{}      `json:"balances"`               
	Birthdate             int64              `json:"birthdate"`              
	BonusTotalWager       int64              `json:"bonus_total_wager"`      
	BonusWager            int64              `json:"bonus_wager"`             
	City                  string             `json:"city"`                   
	ClientCategoryID      int64              `json:"client_category_id"`     
	CompanyID             int64              `json:"company_id"`             
	ConfirmationRequired  int64              `json:"confirmation_required"`  
	ConfirmedPhones       []string           `json:"confirmed_phones"`       
	CountryID             int64              `json:"country_id"`             
	Created               int64              `json:"created"`                
	Currency              string             `json:"currency"`               
	CurrencyChar          string             `json:"currency_char"`          
	CurrencyID            int64              `json:"currency_id"`            
	Demo                  int64              `json:"demo"`                   
	DepositCount          int64              `json:"deposit_count"`          
	DepositInOneClick     bool               `json:"deposit_in_one_click"`   
	Email                 string             `json:"email"`                  
	FinanceState          string             `json:"finance_state"`          
	FirstName             string             `json:"first_name"`             
	Flag                  string             `json:"flag"`                               
	Gender                string             `json:"gender"`                 
	GroupID               int64              `json:"group_id"`               
	ID                    int64              `json:"id"`                     
	Infeed                int64              `json:"infeed"`                 
	IsActivated           bool               `json:"is_activated"`           
	IsIslamic             bool               `json:"is_islamic"`             
	IsVipGroup            bool               `json:"is_vip_group"`                          
	KycConfirmed          bool               `json:"kyc_confirmed"`          
	LastName              string             `json:"last_name"`              
	LastVisit             bool               `json:"last_visit"`             
	Locale                string             `json:"locale"`                 
	Mask                  string             `json:"mask"`                   
	Messages              int64              `json:"messages"`               
	Money                 money              `json:"money"`                  
	Name                  string             `json:"name"`                   
	Nationality           string             `json:"nationality"`            
	NeedPhoneConfirmation bool               `json:"need_phone_confirmation"`
	NewEmail              string             `json:"new_email"`              
	Nickname              string             `json:"nickname"`               
	Phone                 string             `json:"phone"`                  
	Popup                 []interface{}      `json:"popup"`                  
	PostalIndex           string             `json:"postal_index"`           
	Public                int64              `json:"public"`                 
	RateInOneClick        bool               `json:"rate_in_one_click"`      
	SiteID                int64              `json:"site_id"`                
	Skey                  string             `json:"skey"`                                 
	SSID                  bool               `json:"ssid"`                   
	Tc                    bool               `json:"tc"`                     
	Timediff              int64              `json:"timediff"`               
	Tin                   string             `json:"tin"`                    
	TournamentsIDS        interface{}        `json:"tournaments_ids"`        
	TradeRestricted       bool               `json:"trade_restricted"`       
	Trial                 bool               `json:"trial"`                  
	Tz                    string             `json:"tz"`                     
	TzOffset              int64              `json:"tz_offset"`              
	UserCircle            string             `json:"user_circle"`            
	UserGroup             string             `json:"user_group"`             
	UserID                int64              `json:"user_id"`                
	WelcomeSplash         int64              `json:"welcome_splash"`         
}

type money struct {
	Deposit  deposit `json:"deposit"` 
	Withdraw deposit `json:"withdraw"`
}

type deposit struct {
	Max int64 `json:"max"`
	Min int64 `json:"min"`
}

func GetCoreProfile(ws *Socket, serverTimeStamp int64, timeout time.Time) (*CoreProfile, error) {
	eventMsg := map[string]interface{}{
		"name": "core.get-profile",
		"body": map[string]interface{}{},
		"version": "1.0",
	}

	requestEvent := &Event{
		Name:      "sendMessage",
		Msg:       eventMsg,
		RequestId: "2",
		LocalTime: int(serverTimeStamp),
	}

	resp, err := EmitWithResponse(ws, requestEvent, "profile", time.Now().Add(1 * time.Minute))
	if err != nil {
		return nil, err
	}

	responseEvent, err := tjson.Unmarshal[CoreProfile](resp)
	if err != nil {
		return nil, err
	}

	return &responseEvent, nil
}

type UserProfileClient struct {
	RequestID string `json:"request_id"`
	Name      string `json:"name"`
	Msg       struct {
		Balances            []interface{} `json:"balances"`
		ClientCategoryID    int           `json:"client_category_id"`
		CountryID           int           `json:"country_id"`
		Flag                string        `json:"flag"`
		ImgURL              string        `json:"img_url"`
		IsDemoAccount       bool          `json:"is_demo_account"`
		IsVip               bool          `json:"is_vip"`
		VipBadge            bool          `json:"vip_badge"`
		IsSuccessful        bool          `json:"isSuccessful"`
		RegistrationTime    int           `json:"registration_time"`
		SelectedAssetID     int           `json:"selected_asset_id"`
		SelectedBalanceType int           `json:"selected_balance_type"`
		SelectedOptionType  int           `json:"selected_option_type"`
		UserID              int           `json:"user_id"`
		UserName            string        `json:"user_name"`
		Status              string        `json:"status"`
		Gender              string        `json:"gender"`
	} `json:"msg"`
	Status int `json:"status"`
}

func GetUserProfileClient(ws *Socket, userId int, serverTimeStamp int64, timeout time.Time) (*UserProfileClient, error) {
	eventMsg := map[string]interface{}{
		"name": "get-user-profile-client",
		"body": map[string]interface{}{
			"user_id": userId,
		},
		"version": "1.0",
	}

	event := &Event{
		Name:      "sendMessage",
		Msg:       eventMsg,
		RequestId: "157",
	}

	resp, err := EmitWithResponse(ws, event, "user-profile-client", timeout)
	if err != nil {
		return nil, err
	}

	responseEvent, err := tjson.Unmarshal[UserProfileClient](resp)
	if err != nil {
		return nil, err
	}

	return &responseEvent, nil
}
