package types

type Account struct {
	Id                int        `json:"id"`
	Name              string     `json:"name"`
	Owner             Permission `json:"owner"`
	Active            Permission `json:"active"`
	Posting           Permission `json:"posting"`
	MemoKey           string     `json:"memo_key"`
	JsonMetadata      string     `json:"json_metadata"`
	Balance           string     `json:"balance"`
	SavingsBalance    string     `json:"savings_balance"`
	SbdBalance        string     `json:"sbd_balance"`
	SavingsSbdBalance string     `json:"savings_sbd_balance"`
	RewardSbdBalance  string     `json:"reward_sbd_balance"`
}

type Permission struct {
	WeightThreshold uint32        `json:"weight_threshold"`
	AccountAuths    []interface{} `json:"account_auths"`
	KeyAuths        []interface{} `json:"key_auths"`
	AddressAuths    []interface{} `json:"address_auths"`
}

type Options struct {
	MemoKey       string   `json:"memo_key"`
	VotingAccount ObjectID `json:"voting_account"`
}
