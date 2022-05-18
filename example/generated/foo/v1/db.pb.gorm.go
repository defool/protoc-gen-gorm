package v1

// UserOrm is the ORM type for User
type UserOrm struct {
	Id        uint64
	Name      string `gorm:"size:32"`
	UserEmail string
	CompanyId uint64
	Company   *CompanyOrm
	Groups    []*GroupOrm
}

// ToOrm converts the pb message to orm message
func (m *User) ToOrm() *UserOrm {
	if m == nil {
		return nil
	}
	to := &UserOrm{}
	to.Id = m.Id
	to.Name = m.Name
	to.UserEmail = m.UserEmail
	to.CompanyId = m.CompanyId
	to.Company = m.Company.ToOrm()
	for _, v := range m.Groups {
		to.Groups = append(to.Groups, v.ToOrm())
	}
	return to
}

// ToPb converts the orm message to pb message
func (m *UserOrm) ToPb() *User {
	if m == nil {
		return nil
	}
	to := &User{}
	to.Id = m.Id
	to.Name = m.Name
	to.UserEmail = m.UserEmail
	to.CompanyId = m.CompanyId
	to.Company = m.Company.ToPb()
	for _, v := range m.Groups {
		to.Groups = append(to.Groups, v.ToPb())
	}
	return to
}

// CompanyOrm is the ORM type for Company
type CompanyOrm struct {
	Id   uint64
	Name string
}

// ToOrm converts the pb message to orm message
func (m *Company) ToOrm() *CompanyOrm {
	if m == nil {
		return nil
	}
	to := &CompanyOrm{}
	to.Id = m.Id
	to.Name = m.Name
	return to
}

// ToPb converts the orm message to pb message
func (m *CompanyOrm) ToPb() *Company {
	if m == nil {
		return nil
	}
	to := &Company{}
	to.Id = m.Id
	to.Name = m.Name
	return to
}

// GroupOrm is the ORM type for Group
type GroupOrm struct {
	Id   uint64
	Name string
}

// ToOrm converts the pb message to orm message
func (m *Group) ToOrm() *GroupOrm {
	if m == nil {
		return nil
	}
	to := &GroupOrm{}
	to.Id = m.Id
	to.Name = m.Name
	return to
}

// ToPb converts the orm message to pb message
func (m *GroupOrm) ToPb() *Group {
	if m == nil {
		return nil
	}
	to := &Group{}
	to.Id = m.Id
	to.Name = m.Name
	return to
}
