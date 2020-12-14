package consistenthash

func NewRing() *Ring

func (r *Ring) AddNode(id string)

func (r *Ring) RemoveNode(id string) error

func (r *Ring) Get(key string) string
