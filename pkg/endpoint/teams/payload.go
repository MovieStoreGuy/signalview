package teams

type BundledPayload struct {
	Count   int64
	Results []Model
}

type Model struct {
	Description       string
	ID                string
	Members           []string
	Name              string
	NotificationLists interface{} // To implement later
}
