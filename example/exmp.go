package example

type Exmp struct {
	Public  int
	private int
	Other   string
}

func NewExmp() Exmp {
	return Exmp{
		Public:  1,
		private: 5,
		Other:   "3",
	}
}
