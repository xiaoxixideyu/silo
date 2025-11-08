package example

import "time"

//go:generate go run accessor UserModel -f user

type UserStatus int

type UserModel struct {
	BaseModel `accessor:",inline"` // 使用inline表示需要给BaseModel里面的属性也生成记录
	ID        int                  `accessor:"_id"`      // 用_id做为Column名字
	Status    UserStatus           `accessor:"st;eq:=="` // 用st做为Column名字, 並且使用==方法判断是否相等
	Name      string               `accessor:";eq:-"`    // 不判断是否相同，直接标记变更
}

// time.Time, isotime.ISOTIme, decimal.Decimal 默认会使用Equal
type BaseModel struct {
	CreatedAt time.Time      `accessor:";eq:Equal"` // 用Equal方法判断是否相等
	changes   map[string]any `accessor:"-"`         // 用-表示忽略这个属性
}

func (b *BaseModel) Update(name string, value any) {
	if b.changes == nil {
		b.changes = make(map[string]any)
	}
	b.changes[name] = value
}
