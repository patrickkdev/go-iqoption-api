package broker

import (
	"github.com/patrickkdev/Go-IQOption-API/internal/tjson"
)

type CoreProfile struct {
	ClientSessionID string `json:"client_session_id"`
	Msg             struct {
		IsSuccessful bool              `json:"isSuccessful"`
		Message      []interface{}     `json:"message"`
		Result       CoreProfileResult `json:"result"`
	} `json:"msg"`
	Name      string `json:"name"`
	RequestID string `json:"request_id"`
	SessionID string `json:"session_id"`
	Status    int64  `json:"status"`
}

type CoreProfileResult struct {
	AccountStatus        string        `json:"account_status"`
	Address              string        `json:"address"`
	AuthTwoFactor        interface{}   `json:"auth_two_factor"`
	Avatar               string        `json:"avatar"`
	Balance              int64         `json:"balance"`
	BalanceID            int64         `json:"balance_id"`
	BalanceType          int64         `json:"balance_type"`
	Balances             []interface{} `json:"balances"`
	Birthdate            int64         `json:"-"`
	BonusTotalWager      int64         `json:"bonus_total_wager"`
	BonusWager           int64         `json:"bonus_wager"`
	City                 string        `json:"city"`
	ClientCategoryID     int64         `json:"client_category_id"`
	CompanyID            int64         `json:"company_id"`
	ConfirmationRequired int64         `json:"confirmation_required"`
	ConfirmedPhones      []string      `json:"confirmed_phones"`
	CountryID            int64         `json:"country_id"`
	Created              int64         `json:"created"`
	Currency             string        `json:"currency"`
	CurrencyChar         string        `json:"currency_char"`
	CurrencyID           int64         `json:"currency_id"`
	Demo                 int64         `json:"demo"`
	DepositCount         int64         `json:"deposit_count"`
	DepositInOneClick    bool          `json:"deposit_in_one_click"`
	Email                string        `json:"email"`
	FinanceState         string        `json:"finance_state"`
	FirstName            string        `json:"first_name"`
	Flag                 string        `json:"flag"`
	Gender               string        `json:"gender"`
	GroupID              int64         `json:"group_id"`
	ID                   int64         `json:"id"`
	Infeed               int64         `json:"infeed"`
	IsActivated          bool          `json:"is_activated"`
	IsIslamic            bool          `json:"is_islamic"`
	IsVipGroup           bool          `json:"is_vip_group"`
	KycConfirmed         bool          `json:"kyc_confirmed"`
	LastName             string        `json:"last_name"`
	LastVisit            bool          `json:"last_visit"`
	Locale               string        `json:"locale"`
	Mask                 string        `json:"mask"`
	Messages             int64         `json:"messages"`
	Money                struct {
		Deposit struct {
			Max int64 `json:"max"`
			Min int64 `json:"min"`
		} `json:"deposit"`
		Withdraw struct {
			Max int64 `json:"max"`
			Min int64 `json:"min"`
		} `json:"withdraw"`
	} `json:"money"`
	Name                  string        `json:"name"`
	Nationality           string        `json:"nationality"`
	NeedPhoneConfirmation bool          `json:"need_phone_confirmation"`
	NewEmail              string        `json:"new_email"`
	Nickname              string        `json:"nickname"`
	Phone                 string        `json:"phone"`
	Popup                 []interface{} `json:"popup"`
	PostalIndex           string        `json:"postal_index"`
	Public                int64         `json:"public"`
	RateInOneClick        bool          `json:"rate_in_one_click"`
	SiteID                int64         `json:"site_id"`
	Skey                  string        `json:"skey"`
	SSID                  bool          `json:"ssid"`
	Tc                    bool          `json:"tc"`
	Timediff              int64         `json:"timediff"`
	Tin                   string        `json:"tin"`
	TournamentsIDS        interface{}   `json:"tournaments_ids"`
	TradeRestricted       bool          `json:"trade_restricted"`
	Trial                 bool          `json:"trial"`
	Tz                    string        `json:"tz"`
	TzOffset              int64         `json:"tz_offset"`
	UserCircle            string        `json:"user_circle"`
	UserGroup             string        `json:"user_group"`
	UserID                int64         `json:"user_id"`
	WelcomeSplash         int64         `json:"welcome_splash"`
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

// Get core profile data
// This is mostly used to get clientSessionID in order to call GetProfile
func (c *Client) GetCoreProfile() (CoreProfile, error) {
	requestEvent := requestEvent{
		Name: "sendMessage",
		Msg: map[string]interface{}{
			"name":    "core.get-profile",
			"body":    map[string]interface{}{},
			"version": "1.0",
		},
	}.WithRandomRequestId()

	resp, err := c.ws.EmitWithResponse(requestEvent, "profile", c.getTimeout())
	if err != nil {
		return CoreProfile{}, err
	}

	responseEvent, err := tjson.Unmarshal[CoreProfile](resp)
	if err != nil {
		return CoreProfile{}, err
	}

	return responseEvent, nil
}

// Get user profile client
// You need to pass clientSessionID from GetCoreProfile in order to call this function.
func (c *Client) GetUserProfileClient(userId int) (UserProfileClient, error) {
	requestEvent := requestEvent{
		Name: "sendMessage",
		Msg: map[string]interface{}{
			"name": "get-user-profile-client",
			"body": map[string]interface{}{
				"user_id": userId,
			},
			"version": "1.0",
		},
	}.WithRandomRequestId()

	resp, err := c.ws.EmitWithResponse(requestEvent, "user-profile-client", c.getTimeout())
	if err != nil {
		return UserProfileClient{}, err
	}

	responseEvent, err := tjson.Unmarshal[UserProfileClient](resp)
	if err != nil {
		return UserProfileClient{}, err
	}

	return responseEvent, nil
}
