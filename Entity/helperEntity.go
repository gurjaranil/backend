package Entity

type LoginDetails struct {
	Email    string `valid:"email"      json:"email"`
	Password string `valid:"length(6|20)" json:"password"`
}

// Read implements io.Reader.
func (l LoginDetails) Read(p []byte) (n int, err error) {
	panic("unimplemented")
}
