package item

type Item struct {
	key      string
	contents *string
	name     *string
	version  *string
	hmac     *string
}

func NewItem(key string, contents *string, name *string, version *string, hmac *string) *Item {
	return &Item{key, contents, name, version, hmac}
}

func (i *Item) GetKey() string {
	return i.key
}

func (i *Item) GetContents() *string {
	return i.contents
}

func (i *Item) GetName() *string {
	return i.name
}

func (i *Item) GetVersion() *string {
	return i.version
}

func (i *Item) GetHMAC() *string {
	return i.hmac
}
