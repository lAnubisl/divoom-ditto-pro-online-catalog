package gallery

func (c *Client) ListCategories() []Category {
	// Android 3.8.40: CloudClassifyModel's standard gallery tabs, in app order,
	// using its English resources. Device-specific and moderation tabs are omitted.
	return []Category{
		{ID: 20, Name: "Follow"}, {ID: 18, Name: "Recommend"}, {ID: 0, Name: "NEW"},
		{ID: 3, Name: "Character"}, {ID: 17, Name: "Season"}, {ID: 4, Name: "Emoji"},
		{ID: 8, Name: "Pattern"}, {ID: 9, Name: "Creative"}, {ID: 6, Name: "Nature"},
		{ID: 5, Name: "Daily"}, {ID: 15, Name: "Gadget"}, {ID: 7, Name: "Icon"},
		{ID: 16, Name: "Business"}, {ID: 1, Name: "Default"}, {ID: 40, Name: "AI"},
		{ID: 12, Name: "Photo"}, {ID: 19, Name: "Planet"},
	}
}
