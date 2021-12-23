package preprocess

type Condition int
type ConditionStack struct {
	pos int
	s   []Condition
}

const (
	GLOBAL  Condition = iota // 全局
	IN_THEN                  // 执行到同级  #else #elif #endif
	IN_ELSE                  // 执行到 #endif
)

func NewConditionStack() *ConditionStack {
	return &ConditionStack{
		pos: 0,
		s:   []Condition{GLOBAL},
	}
}

// 栈顶
func (c *ConditionStack) Top() Condition {
	if c.pos >= 0 && c.pos < len(c.s) {
		return c.s[c.pos]
	}
	return GLOBAL
}

// 压入栈
func (c *ConditionStack) Push(cdt Condition) {
	if c.pos+1 < len(c.s) {
		c.pos++
		c.s[c.pos] = cdt
		return
	}
	c.pos++
	c.s = append(c.s, cdt)
}

// 弹出栈
func (c *ConditionStack) Pop() Condition {
	if c.pos >= 0 && c.pos < len(c.s) {
		p := c.s[c.pos]
		c.pos--
		return p
	}
	return GLOBAL
}
