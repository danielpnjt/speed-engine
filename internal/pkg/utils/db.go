package utils

import (
	"gorm.io/gorm"
)

type DBCond struct {
	Joins         string
	InnerJoin     string
	GroupBy       string
	JoinArgs      interface{}
	InnerJoinArgs interface{}
	Where         interface{}
	WhereArgs     interface{}
	WhereOr       interface{}
	WhereOrArgs   interface{}
	WhereAnd      interface{}
	WhereAndArgs  interface{}
	Preload       string
	PreloadArgs   interface{}
	Order         interface{}
	Select        interface{}
	Limit         uint
	Offset        uint
}

func CompileConds(db *gorm.DB, conds ...DBCond) *gorm.DB {
	for _, cond := range conds {
		if cond.InnerJoin != "" {
			if cond.InnerJoinArgs != nil {
				if nestedConds, ok := cond.InnerJoinArgs.([]DBCond); ok {
					cond.InnerJoinArgs = func(db *gorm.DB) *gorm.DB {
						return CompileConds(db, nestedConds...)
					}
				}
				db = db.InnerJoins(cond.InnerJoin, cond.InnerJoinArgs)
			} else {
				db = db.InnerJoins(cond.InnerJoin)
			}
			continue
		}
		if cond.Joins != "" {
			if cond.JoinArgs != nil {
				if nestedConds, ok := cond.JoinArgs.([]DBCond); ok {
					cond.JoinArgs = func(db *gorm.DB) *gorm.DB {
						return CompileConds(db, nestedConds...)
					}
				}
				db = db.Joins(cond.Joins, cond.JoinArgs)
			} else {
				db = db.Joins(cond.Joins)
			}
			continue
		}
		if cond.Where != nil {
			if cond.WhereArgs != nil {
				if nestedConds, ok := cond.WhereArgs.([]DBCond); ok {
					cond.WhereArgs = func(db *gorm.DB) *gorm.DB {
						return CompileConds(db, nestedConds...)
					}
				}
				db = db.Where(cond.Where, cond.WhereArgs)
			} else {
				db = db.Where(cond.Where)
			}
			continue
		}
		if cond.WhereOr != nil {
			if cond.WhereOrArgs != nil {
				if nestedConds, ok := cond.WhereOrArgs.([]DBCond); ok {
					cond.WhereOrArgs = func(db *gorm.DB) *gorm.DB {
						return CompileConds(db, nestedConds...)
					}
				}
				db = db.Or(cond.WhereOr, cond.WhereOrArgs)
			} else {
				db = db.Or(cond.WhereOr)
			}
			continue
		}
		if cond.WhereAnd != nil {
			if cond.WhereAndArgs != nil {
				if nestedConds, ok := cond.WhereAndArgs.([]DBCond); ok {
					cond.WhereAndArgs = func(db *gorm.DB) *gorm.DB {
						return CompileConds(db, nestedConds...)
					}
				}
				db = db.Where(cond.WhereAnd, cond.WhereAndArgs)
			} else {
				db = db.Where(cond.WhereAnd)
			}
			continue
		}
		if cond.Preload != "" {
			if cond.PreloadArgs != nil {
				if nestedConds, ok := cond.PreloadArgs.([]DBCond); ok {
					cond.PreloadArgs = func(db *gorm.DB) *gorm.DB {
						return CompileConds(db, nestedConds...)
					}
				}
				db = db.Preload(cond.Preload, cond.PreloadArgs)
			} else {
				db = db.Preload(cond.Preload)
			}
			continue
		}
		if cond.GroupBy != "" {
			db = db.Group(cond.GroupBy)
			continue
		}
		if cond.Order != nil {
			db = db.Order(cond.Order)
			continue
		}
		if cond.Select != nil {
			db = db.Select(cond.Select)
			continue
		}
		if cond.Limit > 0 {
			db = db.Limit(int(cond.Limit))
			continue
		}
		if cond.Offset > 0 {
			db = db.Offset(int(cond.Offset))
			continue
		}
	}
	return db
}
