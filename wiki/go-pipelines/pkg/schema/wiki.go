package schema

import (
	// "encoding/json"
	"time"
)

// Original schema returned from wiki event stream for page create
type WikiPageCreateOrigin struct {
	Schema string `json:"$schema"`
	Meta   struct {
		URI       string    `json:"uri"`
		RequestID string    `json:"request_id"`
		ID        string    `json:"id"`
		Dt        time.Time `json:"dt"`
		Domain    string    `json:"domain"`
		Stream    string    `json:"stream"`
		Topic     string    `json:"topic"`
		Partition int       `json:"partition"`
		Offset    int64     `json:"offset"`
	} `json:"meta"`
	Database         string `json:"database"`
	PageID           int64  `json:"page_id"`
	PageTitle        string `json:"page_title"`
	PageNamespace    int    `json:"page_namespace"`
	RevID            int64  `json:"rev_id"`
	RevTimestamp     string `json:"rev_timestamp"`
	RevSHA1          string `json:"rev_sha1"`
	RevMinorEdit     bool   `json:"rev_minor_edit"`
	RevLen           int    `json:"rev_len"`
	RevContentModel  string `json:"rev_content_model"`
	RevContentFormat string `json:"rev_content_format"`
	Performer        struct {
		UserText           string    `json:"user_text"`
		UserGroups         []string  `json:"user_groups"`
		UserIsBot          bool      `json:"user_is_bot"`
		UserID             int       `json:"user_id"`
		UserRegistrationDt time.Time `json:"user_registration_dt"`
		UserEditCount      int       `json:"user_edit_count"`
	} `json:"performer"`
	PageIsRedirect bool   `json:"page_is_redirect"`
	Comment        string `json:"comment"`
	ParsedComment  string `json:"parsedcomment"`
	Dt             string `json:"dt"`
	RevSlots       struct {
		Main struct {
			RevSlotSHA1        string `json:"rev_slot_sha1"`
			RevSlotSize        int    `json:"rev_slot_size"`
			RevSlotOriginRevID int64  `json:"rev_slot_origin_rev_id"`
		} `json:"main"`
	} `json:"rev_slots"`
}

// Flattened and filtered version of page create
type WikiPageCreateLocal struct {
	Uri                string    `json:"uri"`
	Domain             string    `json:"domain"`
	Stream             string    `json:"stream"`
	Database           string    `json:"database"`
	PageID             int64     `json:"page_id"`
	PageTitle          string    `json:"page_title"`
	PageNamespace      int       `json:"page_namespace"`
	UserIsBot          bool      `json:"user_is_bot"`
	UserID             int       `json:"user_id"`
	UserRegistrationDt time.Time `json:"user_registration_dt"`
	UserEditCount      int       `json:"user_edit_count"`
	Comment            string    `json:"comment"`
	ParsedComment      string    `json:"parsedcomment"`
	Dt                 string    `json:"dt"`
	Page               string    `json:page`
}

// // Parse the wiki page create json
// func WikiPageCreate(jsonData string) (PageCreate, error) {

// 	// Create struct
// 	var origin originPageCreate

// 	// Unpack json
// 	err := json.Unmarshal([]byte(jsonData), &origin)

// 	// Handle error
// 	if err != nil {
// 		return PageCreate{}, err
// 	}

// 	// Create the new structure
// 	novel := PageCreate{
// 		Uri:                origin.Meta.URI,
// 		Domain:             origin.Meta.Domain,
// 		Database:           origin.Database,
// 		PageID:             origin.PageID,
// 		PageTitle:          origin.PageTitle,
// 		PageNamespace:      origin.PageNamespace,
// 		UserIsBot:          origin.Performer.UserIsBot,
// 		UserID:             origin.Performer.UserID,
// 		UserRegistrationDt: origin.Performer.UserRegistrationDt,
// 		UserEditCount:      origin.Performer.UserEditCount,
// 		Comment:            origin.Comment,
// 		ParsedComment:      origin.ParsedComment,
// 		Dt:                 origin.Dt,
// 	}

// 	// Return
// 	return novel, nil
// }
