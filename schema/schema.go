package schema

import (
	"errors"

	"github.com/Fiber-Man/albumEditorPlugin/model"
	"github.com/Fiber-Man/funplugin"
	"github.com/Fiber-Man/funplugin/plugin"
	"github.com/graphql-go/graphql"
)

var deviceSchema *funplugin.ObjectSchema
var applyformSchema *funplugin.ObjectSchema
var load = false

func Init() {
	// InitAccount()
	if employee, ok := plugin.GetObject("employee"); !ok {
		panic(errors.New("not have object type"))
	} else {
		albumEditor := graphql.NewObject(graphql.ObjectConfig{
			Name:        "employee_albumEditor",
			Description: "employee_albumEditor plugin",
			Fields: graphql.Fields{
				"applyforms": applyformSchema.Query["applyforms"],
			},
		})

		employee.AddFieldConfig("albumEditor", &graphql.Field{
			Name:        "queryemployee_albumEditor",
			Description: "queryemployee_albumEditor",
			Type:        albumEditor,
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return params.Source, nil
			},
		})
	}

	deviceSchema.GraphQLType.AddFieldConfig("applyforms", applyformSchema.Query["applyforms"])

	if field, err := plugin.AutoField("EmployeeID:employee"); err != nil {
		panic(errors.New("not have object type"))
	} else {
		applyformSchema.GraphQLType.AddFieldConfig("employee", field)
	}

	if field, err := plugin.AutoField("DistributerID:employee"); err != nil {
		panic(errors.New("not have object type"))
	} else {
		applyformSchema.GraphQLType.AddFieldConfig("distributer", field)
	}

	if field, err := plugin.AutoField("device"); err != nil {
		panic(errors.New("not have object type"))
	} else {
		applyformSchema.GraphQLType.AddFieldConfig("device", field)
	}
}

func marge(oc *funplugin.ObjectSchema) {
	for k, v := range oc.Query {
		queryFields[k] = v
	}
	for k, v := range oc.Mutation {
		mutationFields[k] = v
	}
}

var queryFields = graphql.Fields{
	// "account":  &queryAccount,
	// "accounts": &queryAccountList,
	// "authority":  &queryAuthority,
	// "authoritys": &queryAuthorityList,
}

var mutationFields = graphql.Fields{
	// "createAccount": &createAccount,
	// "updateAccount": &updateAccount,
}

// NewSchema 用于插件主程序调用
func NewPlugSchema(pls funplugin.PluginManger) funplugin.Schema {
	if load != true {
		deviceSchema, _ = pls.NewSchemaBuilder(model.Device{})
		marge(deviceSchema)

		applyformSchema, _ = pls.NewSchemaBuilder(model.ApplyForm{})
		marge(applyformSchema)
		load = true
	}

	// roleSchema, _ := pls.NewSchemaBuilder(model.Role{})
	// marge(roleSchema)

	// roleAccountSchema, _ := pls.NewSchemaBuilder(model.RoleAccount{})
	// marge(roleAccountSchema)

	return funplugin.Schema{
		Object: map[string]*graphql.Object{
			// "account": accountType,
			"device":    deviceSchema.GraphQLType,
			"applyform": applyformSchema.GraphQLType,
			// "role":        roleSchema.GraphQLType,
			// "roleaccount": roleAccountSchema.GraphQLType,
		},
		Query:    queryFields,
		Mutation: mutationFields,
	}
}
