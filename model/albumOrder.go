/**
 * @Author zhangfan
 * @create 2021/2/25 上午9:50
 * Description:
 */

package model

import (
	"errors"
	"github.com/Fiber-Man/funplugin/plugin"
	"github.com/graphql-go/graphql"
	"time"
)

type AlbumOrder struct {
	ID           uint                        `gorm:"primary_key" gqlschema:"update!;delete!;query!;querys" description:"订单id"`
	UserId       uint                        `gorm:"DEFAULT:0;NOT NULL;" gqlschema:"create!;querys" description:"创建用户id" funservice:"employee"`
	SinglePrice  uint                        `gorm:"DEFAULT:0;NOT NULL;" gqlschema:"create;querys" description:"单价"`
	Amount       uint                        `gorm:"DEFAULT:0;NOT NULL;" gqlschema:"create;querys" description:"数量"`
	TotalPrice   uint                        `gorm:"DEFAULT:0;NOT NULL;" gqlschema:"create;querys" description:"总价"`
	Specs        AlbumOrderSpecsEnumType     `gorm:"DEFAULT:1;NOT NULL;" gqlschema:"create;querys" description:"纸张规格"`
	Material     AlbumOrderMaterialEnumType  `gorm:"DEFAULT:1;NOT NULL;" gqlschema:"create;querys" description:"材质"`
	Template     string                      `gorm:"Type:varchar(255);DEFAULT:'';NOT NULL;" gqlschema:"create;querys" description:"使用的相册模板"`
	Address      string                      `gorm:"Type:varchar(1000);DEFAULT:'';NOT NULL;" gqlschema:"create;querys" description:"收货地址"`
	UsageType    AlbumOrderUsageTypeEnumType `gorm:"DEFAULT:1;NOT NULL;" gqlschema:"create;querys"  description:"使用类型枚举"`
	Remark       string                      `gorm:"Type:varchar(1000);DEFAULT:'';NOT NULL;" gqlschema:"create;querys" description:"备注"`
	FreightPrice uint                        `gorm:"DEFAULT:0;NOT NULL;" gqlschema:"create;querys" description:"运费"`
	PaymentId    string                      `gorm:"Type:varchar(1000);DEFAULT:'';NOT NULL;" gqlschema:"create;update;querys" description:"支付id"`
	PayWay       AlbumOrderPayWayEnumType    `gorm:"DEFAULT:1;NOT NULL;" gqlschema:"create;update;querys" description:"支付方式枚举类型"`
	PayTime      time.Time                   `gorm:"DEFAULT:'1970-1-1 00:00:00';" description:"支付时间" gqlschema:"querys"`
	DeliveryId   string                      `gorm:"Type:varchar(1000);DEFAULT:'';NOT NULL;" gqlschema:"update;querys" description:"快递id"`
	Status       AlbumOrderStatusEnumType    `gorm:"DEFAULT:1;NOT NULL;" gqlschema:"create;update;querys" description:"订单状态枚举类型"`
	CreatedAt    time.Time                   `description:"创建时间" gqlschema:"querys"`
	UpdatedAt    time.Time                   `description:"更新时间" gqlschema:"querys"`
	DeletedAt    *time.Time
	v2           int `gorm:"-" exclude:"true"`
}

type AlbumOrders struct {
	TotalCount int
	Edges      []AlbumOrder
}

// 跨接口创建订单
func createOrder(userId uint, referId uint, referType string, remark string) (err error) {
	mutation := `mutation ($userId: Int!, $childrenId: Int!, $childrenType: OrderInfoChildrenTypeEnum!, $remark: String) {
				  order {
					orderinfos {
					  action {
						create(userId: $userId, childrenId: $childrenId, childrenType: $childrenType, remark: $remark) {
						  id
						}
					  }
					}
				  }
				}`
	params := map[string]interface{}{
		"userId":       userId,
		"childrenId":   referId,
		"childrenType": referType,
		"remark":       remark,
	}
	_, err = plugin.Go(mutation, params, nil)
	if err != nil {
		return err
	}
	return
}

func (o AlbumOrder) Query(params graphql.ResolveParams) (AlbumOrder, error) {
	p := params.Args
	err := db.Where(p).First(&o).Error
	return o, err
}

func (o AlbumOrder) Querys(params graphql.ResolveParams) (AlbumOrders, error) {
	var result AlbumOrders

	dbselect := GenSelet(db, params)
	dbcount := GenWhere(db.Model(o), params)

	err := dbselect.Find(&result.Edges).Error
	if err != nil {
		return result, err
	}
	err = dbcount.Count(&result.TotalCount).Error
	return result, err
}

func (o AlbumOrder) Create(params graphql.ResolveParams) (AlbumOrder, error) {
	p := params.Args
	o.UserId = uint(p["userId"].(int))
	// template
	if p["singlePrice"] != nil {
		o.SinglePrice = uint(p["singlePrice"].(int))
	}
	if p["amount"] != nil {
		o.Amount = uint(p["amount"].(int))
	}
	if p["totalPrice"] != nil {
		o.TotalPrice = uint(p["totalPrice"].(int))
	}
	if p["specs"] != nil {
		o.Specs = p["specs"].(AlbumOrderSpecsEnumType)
	}
	if p["material"] != nil {
		o.Material = p["material"].(AlbumOrderMaterialEnumType)
	}
	if p["template"] != nil {
		o.Template = p["template"].(string)
	}
	if p["address"] != nil {
		o.Address = p["address"].(string)
	}
	if p["usageType"] != nil {
		o.UsageType = p["usageType"].(AlbumOrderUsageTypeEnumType)
	}
	if p["remark"] != nil {
		o.Remark = p["remark"].(string)
	}
	if p["freightPrice"] != nil {
		o.TotalPrice = uint(p["freightPrice"].(int))
	}
	if p["paymentId"] != nil {
		o.PaymentId = p["paymentId"].(string)
	}
	if p["payWay"] != nil {
		o.PayWay = p["payWay"].(AlbumOrderPayWayEnumType)
	}
	o.PayTime = time.Now()
	if p["status"] != nil {
		o.Status = p["status"].(AlbumOrderStatusEnumType)
	}
	err := db.Create(&o).Error
	err = createOrder(o.UserId, o.ID, "PHOTO_ALBUM", o.Remark)
	return o, err
}

func (o AlbumOrder) Update(params graphql.ResolveParams) (AlbumOrder, error) {
	v, ok := params.Source.(AlbumOrder)
	if !ok {
		return o, errors.New("update param")
	}
	p := params.Args
	if p["paymentId"] != nil {
		v.PaymentId = p["paymentId"].(string)
	}
	if p["payWay"] != nil {
		v.PayWay = p["payWay"].(AlbumOrderPayWayEnumType)
	}
	v.PayTime = time.Now()
	if p["deliveryId"] != nil {
		v.DeliveryId = p["deliveryId"].(string)
	}
	if p["status"] != nil {
		v.Status = p["status"].(AlbumOrderStatusEnumType)
	}
	err := db.Save(&v).Error
	return v, err
}

func (o AlbumOrder) Delete(params graphql.ResolveParams) (AlbumOrder, error) {
	v, ok := params.Source.(AlbumOrder)
	if !ok {
		return o, errors.New("delete param")
	}
	err := db.Delete(&v).Error
	return v, err
}
