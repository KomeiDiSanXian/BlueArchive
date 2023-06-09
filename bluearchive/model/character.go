package model

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// Character 角色信息
type Character struct {
	Name         string       `gorm:"primary_key"` // 主键使用日文罗马
	NameDisplay  string       // 显示名
	Profile      Profile      `gorm:"foreignkey:CharaName"` // 个人信息
	Property     Property     `gorm:"foreignkey:CharaName"` // 数据
	Star         int          // 初始星级数
	Type         string       // 角色属性
	AttackType   string       // 攻击属性
	ArmorType    string       // 装甲属性
	Firearm      Weapon       `gorm:"foreignkey:CharaName"` // 枪械相关
	Position     string       // 站位
	School       string       // 所属学校
	Equip        Equipment    `gorm:"foreignkey:CharaName"` // 装备
	Adapat       Adapatation  `gorm:"foreignkey:CharaName"` // 场地适应性
	ExSkill      []CharaSkill `gorm:"foreignkey:CharaName"` // Ex技
	NormalSkill  []CharaSkill `gorm:"foreignkey:CharaName"` // 普通技
	PassiveSkill []CharaSkill `gorm:"foreignkey:CharaName"` // 被动技
	SupportSkill []CharaSkill `gorm:"foreignkey:CharaName"` // 支援技
}

// Profile 角色个人信息
type Profile struct {
	CharaName    string `gorm:"primary_key"` // 外键
	FullName     string // 全名
	TransCN      string // 繁中译名
	Age          int    // 年龄
	BirthDate    string // 生日
	Height       int    // 身高
	Hobby        string // 爱好
	Illustrator  string // 画师
	Vocal        string // 声优
	Introduction string // 简介
}

// Property 角色属性数据
type Property struct {
	CharaName     string `gorm:"primary_key"` // 外键
	Attack        int    // 攻击力
	Denfence      int    // 防御力
	HealthPoint   int    // 生命值
	HealPoint     int    // 治愈力
	Accuracy      int    // 命中率
	Avoidance     int    // 回避率
	CritialRate   int    // 暴击率
	CritialDamage int    // 暴击伤害
	Stability     int    // 稳定性
	Range         int    // 射程
	CCIntensity   int    // CC强化力
	CCResistance  int    // CC抵抗力
}

// Weapon 枪械属性
type Weapon struct {
	CharaName  string        `gorm:"primary_key"` // 外键
	Name       string        `gorm:"primary_key"` // 武器名
	Type       string        // 种类
	Desription string        // 描述
	Stars      []WeaponStar  `gorm:"foreignkey:WeaponName"` // 星级相关
	Skill      []WeaponSkill `gorm:"foreignkey:WeaponName"` // 技能相关
}

// WeaponStar 枪械星级增幅
type WeaponStar struct {
	WeaponName string `gorm:"primary_key"` // 外键
	Star       int    `gorm:"primary_key"` // 星级
	Instensity string // 星级描述
}

// WeaponSkill 枪械技能
type WeaponSkill struct {
	WeaponName string `gorm:"primary_key"` // 外键
	Name       string `gorm:"primary_key"` // 技能名
	Level      int    // 技能等级
	Desription string // 技能描述
}

// Equipment 角色装备
type Equipment struct {
	CharaName string `gorm:"primary_key"` // 外键
	Slot0     string // 装备槽
	Slot1     string
	Slot2     string
	Slot3     string // 爱用品
}

// Adapatation 角色场地适应性
type Adapatation struct {
	CharaName string `gorm:"primary_key"` // 外键
	Outdoor   int    // 室外
	Indoor    int    // 室内
	Street    int    // 巷战
}

// CharaSkill 角色技能
type CharaSkill struct {
	CharaName  string     `gorm:"primary_key"` // 外键
	Name       string     `gorm:"primary_key"` // 技能名
	Level      int        // 技能等级
	Cost       int        // 花费, 普通和被动为0
	Usage      SkillUsage `gorm:"foreignkey:SkillName"` // 升级花费
	Desription string     // 描述
}

// SkillUsage 升级使用材料花费
type SkillUsage struct {
	SkillName string `gorm:"primary_key"` // 外键
	Level     int    `gorm:"primary_key"` // 技能等级
	Materials string `gorm:"primary_key"` // 材料名
	Expense   int    // 花费材料数
}

// CharacterRepository 是 Character 的数据库仓库
type CharacterRepository struct {
	db *gorm.DB
}

// NewCharacterRepository 创建一个 CharacterRepository 实例
func NewCharacterRepository(db *gorm.DB) *CharacterRepository {
	return &CharacterRepository{
		db: db,
	}
}

// CreateCharacter 在数据库中创建一个角色
func (r *CharacterRepository) CreateCharacter(character *Character) error {
	err := r.db.Create(character).Error
	return err
}

// GetCharacterByID 根据角色ID从数据库中获取角色
func (r *CharacterRepository) GetCharacterByID(characterID string) (*Character, error) {
	var character Character
	err := r.db.Where("name = ?", characterID).Preload("Profile").Preload("Property").Preload("Firearm").
		Preload("Equip").Preload("Adapat").Preload("ExSkill").Preload("NormalSkill").Preload("PassiveSkill").
		Preload("SupportSkill").First(&character).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 角色不存在
		}
		return nil, err
	}

	return &character, nil
}

// UpdateCharacter 更新角色信息
func (r *CharacterRepository) UpdateCharacter(character *Character) error {
	err := r.db.Save(character).Error
	return err
}

// DeleteCharacterByID 根据角色ID从数据库中删除角色
func (r *CharacterRepository) DeleteCharacterByID(characterID string) error {
	err := r.db.Where("name = ?", characterID).Delete(Character{}).Error
	return err
}
