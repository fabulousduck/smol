package registertable

/*
RegisterTable is a simple collection of registers so they can be indexed
*/
type RegisterTable map[int]Register

/*
Register simulates a basic CPU register
*/
type Register struct {
	Value int
	Name  string
}

/*
Find finds a variable on the memory table
Returns the index at which it is found
Returns -1 if the value cannot be found
*/
func (table RegisterTable) Find(name string) int {
	for i := 0; i < len(table); i++ {
		region := table[i]
		if region.Name == name {
			return i
		}
	}
	return -1
}
