package services

// MobileCategory represents a single mobile phone category.
type MobileCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// GetMobileCategories returns the supported mobile categories.
func GetMobileCategories() []MobileCategory {
	return []MobileCategory{
		{ID: 1, Name: "Smartphone", Slug: "smartphone"},
		{ID: 2, Name: "Cellphone", Slug: "cellphone"},
		{ID: 3, Name: "Keypad Phone", Slug: "keypad-phone"},
		{ID: 4, Name: "Wireless Phone", Slug: "wireless-phone"},
	}
}
