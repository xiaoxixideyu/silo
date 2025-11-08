# accessor
根据model struct生成对应column和setter方法

## 安装
```
go install .
```

## 使用方式
```sh
Usage:
  accessor [struct to parse] [flags]

Examples:
accessor User User2

Flags:
  -c, --cast string     camel or snake, used when update track is true (default "camel")
  -C, --column string   column middle name (default "Column")
  -f, --file string     set file to search, use filename without .go suffix (default ".")
  -g, --getter          generate getter (default false)
  -h, --help            help for accessor
  -m, --method string   update method name, used when update track is true (default "Update")
  -p, --path string     set path to search (default ".")
      --prefix string   method prefix
  -s, --setter          generate setter (default true)
  -u, --update          with update track (default true)
```

## 例子

```go
// user.go

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
	CreatedAt time.Time              `accessor:";eq:Equal"` // 用Equal方法判断是否相等
	changes   map[string]interface{} `accessor:"-"`         // 用-表示忽略这个属性
}

func (b *BaseModel) Update(name string, value interface{}) {
	if b.changes == nil {
		b.changes = make(map[string]interface{})
	}
	b.changes[name] = value
}
```

生成user_accessor.go
```go
// user_accessor.go

const (
	UserModelColumnCreatedAt = "createdAt"
	UserModelColumnId        = "_id"
	UserModelColumnStatus    = "st"
	UserModelColumnName      = "name"
)

func (u *UserModel) SetCreatedAt(v time.Time) {
	if !u.CreatedAt.Equal(v) {
		u.CreatedAt = v
		u.Update(UserModelColumnCreatedAt, v)
	}
}

func (u *UserModel) SetID(v int) {
	if u.ID != v {
		u.ID = v
		u.Update(UserModelColumnId, v)
	}
}

func (u *UserModel) SetStatus(v UserStatus) {
	if u.Status != v {
		u.Status = v
		u.Update(UserModelColumnStatus, v)
	}
}

func (u *UserModel) SetName(v string) {
	u.Name = v
	u.Update(UserModelColumnName, v)
}

```