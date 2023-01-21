package tgbot

// These are one of the possible events Handle() can deal with.
//
// For convenience, all Telebot-provided endpoints start with
// an "alert" character \a.
const (
	// Basic message handlers.
	OnText              = "\atext"
	OnEdited            = "\aedited"
	OnPhoto             = "\aphoto"
	OnAudio             = "\aaudio"
	OnAnimation         = "\aanimation"
	OnDocument          = "\adocument"
	OnSticker           = "\asticker"
	OnVideo             = "\avideo"
	OnVoice             = "\avoice"
	OnVideoNote         = "\avideo_note"
	OnContact           = "\acontact"
	OnLocation          = "\alocation"
	OnVenue             = "\avenue"
	OnDice              = "\adice"
	OnInvoice           = "\ainvoice"
	OnPayment           = "\apayment"
	OnGame              = "\agame"
	OnPoll              = "\apoll"
	OnPollAnswer        = "\apoll_answer"
	OnPinned            = "\apinned"
	OnChannelPost       = "\achannel_post"
	OnEditedChannelPost = "\aedited_channel_post"

	OnAddedToGroup      = "\aadded_to_group"
	OnUserJoined        = "\auser_joined"
	OnUserLeft          = "\auser_left"
	OnNewGroupTitle     = "\anew_chat_title"
	OnNewGroupPhoto     = "\anew_chat_photo"
	OnGroupPhotoDeleted = "\achat_photo_deleted"
	OnGroupCreated      = "\agroup_created"
	OnSuperGroupCreated = "\asupergroup_created"
	OnChannelCreated    = "\achannel_created"

	// OnMigration happens when group switches to
	// a supergroup. You might want to update
	// your internal references to this chat
	// upon switching as its ID will change.
	OnMigration = "\amigration"

	OnMedia           = "\amedia"
	OnCallback        = "\acallback"
	OnQuery           = "\aquery"
	OnInlineResult    = "\ainline_result"
	OnShipping        = "\ashipping_query"
	OnCheckout        = "\apre_checkout_query"
	OnMyChatMember    = "\amy_chat_member"
	OnChatMember      = "\achat_member"
	OnChatJoinRequest = "\achat_join_request"
	OnProximityAlert  = "\aproximity_alert_triggered"
	OnAutoDeleteTimer = "\amessage_auto_delete_timer_changed"
	OnWebApp          = "\aweb_app"

	OnVideoChatStarted      = "\avideo_chat_started"
	OnVideoChatEnded        = "\avideo_chat_ended"
	OnVideoChatParticipants = "\avideo_chat_participants_invited"
	OnVideoChatScheduled    = "\avideo_chat_scheduled"
)

// ChatAction is a client-side status indicating bot activity.
type ChatAction string

const (
	Typing            ChatAction = "typing"
	UploadingPhoto    ChatAction = "upload_photo"
	UploadingVideo    ChatAction = "upload_video"
	UploadingAudio    ChatAction = "upload_audio"
	UploadingDocument ChatAction = "upload_document"
	UploadingVNote    ChatAction = "upload_video_note"
	RecordingVideo    ChatAction = "record_video"
	RecordingAudio    ChatAction = "record_audio"
	RecordingVNote    ChatAction = "record_video_note"
	FindingLocation   ChatAction = "find_location"
	ChoosingSticker   ChatAction = "choose_sticker"
)

// ParseMode determines the way client applications treat the text of the message
type ParseMode = string

const (
	ModeDefault    ParseMode = ""
	ModeMarkdown   ParseMode = "Markdown"
	ModeMarkdownV2 ParseMode = "MarkdownV2"
	ModeHTML       ParseMode = "HTML"
)
