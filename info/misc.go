package info

import "context"

// MaxBuilderFee returns the maximum builder fee (in tenths of a basis point)
// that user has approved for the given builder address.
func (c *Client) MaxBuilderFee(ctx context.Context, user, builder string) (int, error) {
	var out int
	if err := c.post(ctx, map[string]any{"type": "maxBuilderFee", "user": user, "builder": builder}, &out); err != nil {
		return 0, err
	}
	return out, nil
}

// LegalCheckResult reports the user's legal-check status.
type LegalCheckResult struct {
	Accepted             bool     `json:"accepted"`
	IpFromCountry        string   `json:"ipFromCountry,omitempty"`
	UserFlaggedAddresses []string `json:"userFlaggedAddresses,omitempty"`
	IsCountryAllowed     bool     `json:"isCountryAllowed,omitempty"`
}

// LegalCheck queries whether the user has accepted terms / passed geo checks.
func (c *Client) LegalCheck(ctx context.Context, user string) (*LegalCheckResult, error) {
	var out LegalCheckResult
	if err := c.post(ctx, map[string]any{"type": "legalCheck", "user": user}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ExtraAgent is one delegated-agent address authorized for a user.
type ExtraAgent struct {
	Address    string `json:"address"`
	Name       string `json:"name"`
	ValidUntil int64  `json:"validUntil"`
}

// ExtraAgents returns delegated agent wallets approved by user.
func (c *Client) ExtraAgents(ctx context.Context, user string) ([]ExtraAgent, error) {
	var out []ExtraAgent
	if err := c.post(ctx, map[string]any{"type": "extraAgents", "user": user}, &out); err != nil {
		return nil, err
	}
	return out, nil
}
