package address

type AddressIndex interface {
	Lookup(addr string) (userID string, ok bool)
}
